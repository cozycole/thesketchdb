import { CatalogFilter } from "../components/catalogFilter.js";
import { FilterContent } from "../components/filterMenu.js";

export function initViewCatalog() {
  // scroll up on pagination button click
  addPaginationListener();

  document.addEventListener("htmx:afterSwap", (e) => {
    addPaginationListener();
  });

  // style tab bar
  const tabContainer = document.getElementById("tabContainer");
  const tabUnderline = document.getElementById("tabUnderline");
  const tabs = document.querySelectorAll(".tab");

  // Set active tab
  const path = window.location.pathname;
  const match = path.match(/^\/catalog\/(\w+)/);
  const currentTab = match ? match[1] : null;

  const activeTab = document.querySelector(".tab.active");
  if (activeTab) updateUnderline(activeTab);

  // On tab click
  tabs.forEach((tab) => {
    if (tab.dataset.tab === currentTab) styleTabs(tab, tabs);

    tab.addEventListener("click", () => {
      styleTabs(tab, tabs);
    });
  });
}

function styleTabs(activeTab, tabs) {
  tabs.forEach((t) =>
    t.classList.remove("text-orange-600", "font-bold", "active"),
  );
  activeTab.classList.add("text-orange-600", "font-bold", "active");
  activeTab.classList.remove("hover:text-slate-700");
  updateUnderline(activeTab);

  if (!activeTab.dataset.tab) {
    const tabName = activeTab.textContent;
    showToast(`${tabName} catalog coming soon!`, "success");
  }
}

function updateUnderline(tab) {
  const containerRect = tabContainer.getBoundingClientRect();
  const tabRect = tab.getBoundingClientRect();
  const left = tabRect.left - containerRect.left + tabContainer.scrollLeft;

  tabUnderline.style.left = `${left}px`;
  tabUnderline.style.width = `${tab.offsetWidth}px`;
}

function addPaginationListener() {
  document.querySelectorAll(".htmxSearchPage").forEach((e) => {
    e.addEventListener("click", () => {
      window.scrollTo(0, 0);
    });
  });
}

function showToast(message, type = "info", duration = 3000) {
  const container = document.getElementById("toast-container");

  const toast = document.createElement("div");
  toast.className = `
    flex items-center max-w-xs w-full text-white px-4 py-3 rounded shadow-lg
    transition transform duration-300 ease-in-out translate-x-4 opacity-0
    ${
      type === "success"
        ? "bg-orange-500"
        : type === "error"
          ? "bg-red-500"
          : type === "warning"
            ? "bg-yellow-500 text-black"
            : "bg-slate-800"
    }
  `;
  toast.innerHTML = `
    <span class="flex-1">${message}</span>
  `;

  container.appendChild(toast);

  // Trigger animation
  requestAnimationFrame(() => {
    toast.classList.remove("translate-x-4", "opacity-0");
  });

  setTimeout(() => {
    toast.classList.add("opacity-0", "translate-x-4");
    toast.addEventListener("transitionend", () => toast.remove());
  }, duration);
}
