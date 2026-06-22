import htmx from "htmx.org";

let pendingPaginationScroll = false;

export function initViewShow() {
  setupNavbar();

  document.addEventListener("change", (e) => {
    if (e.target.id !== "sortDropdown") return;

    const url = new URL(window.location.href);

    url.searchParams.set("sort", e.target.value);
    url.searchParams.delete("page");

    htmx.ajax("GET", url.pathname + url.search, {
      target: "#results",
      swap: "innerHTML",
    });

    history.pushState({}, "", url.pathname + url.search);
  });

  // This is needed since we don't want to scroll to the top of
  // the page, swapped content can alter the height of the page, causing
  // incorrect scrolling distance.
  addPaginationListener();
  document.body.addEventListener("htmx:afterSettle", (e) => {
    if (!pendingPaginationScroll) return;

    if (e.target.id === "results") {
      scrollToShowNavigation();
      pendingPaginationScroll = false;
    }
  });
}

function addPaginationListener() {
  document.addEventListener("click", (e) => {
    if (e.target.classList.contains("htmxSearchPage")) {
      pendingPaginationScroll = true;
    }
  });
}

function scrollToShowNavigation() {
  const anchor = document.getElementById("showNavigationAnchor");
  if (!anchor) return;

  anchor.scrollIntoView({
    behavior: "smooth",
    block: "start",
  });
}

function setupNavbar() {
  const navLinks = document.querySelectorAll("#showNavigation a");
  const selectedStyles = ["border-slate-950", "font-bold", "text-slate-950"];
  const unselectedStyles = [
    "border-transparent",
    "font-medium",
    "text-slate-600",
    "hover:border-slate-300",
  ];
  navLinks.forEach((e) => {
    e.addEventListener("click", () => {
      // remove selected styling
      navLinks.forEach((e) => {
        e.classList.remove(...selectedStyles);
        e.classList.add(...unselectedStyles);
      });

      e.classList.remove(...unselectedStyles);
      e.classList.add(...selectedStyles);
    });
  });
}
