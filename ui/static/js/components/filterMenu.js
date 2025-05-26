export class FilterContent extends HTMLElement {
  constructor() {
    super();
    this.mobileBreakpoint = 1080;
    this._onResize = () => this.updateLayout();
    this._applyOnClick = () => this.applyFilters();
    this.desktopContainer = document.getElementById("desktopFilters");
    this.mobileContainer = document.getElementById("mobileFilters");
    this.filterContent = document.getElementById("filters");

    this.mobileFilterButton = document.querySelector("#mobileFilterMenu");
    this.mobileFilters = document.querySelector("#mobileFilters");
    this.applyButton = this.querySelector("#filterApply")

    this._onMobileMenuClick = () => this.mobileFilterToggle();
    this._mobileMenuExit = (e) => this.mobileFitlerExit(e);
    htmx.process(this.filterContent);
  }

  connectedCallback() {
    console.log("FilterContent connected")
    this.desktopContainer.innerHTML = "";
    this.mobileContainer.innerHTML = "";
    this.currentState = "";


    this.applyButton.addEventListener("click", this._applyOnClick);
    this.mobileFilterButton.addEventListener("click", this._onMobileMenuClick);
    document.addEventListener("click", this._mobileMenuExit);
    window.addEventListener("resize", this._onResize);
    this.updateLayout();
  }

  disconnectedCallback() { 
    console.log("FilterContent disconnected")
    this.applyButton.removeEventListener('click', this._applyOnClick);
    this.mobileFilterButton.removeEventListener("click", this._onMobileMenuClick);
    document.removeEventListener("click", this._mobileMenuExit);
    window.removeEventListener("resize", this._onResize);
  }

  updateLayout() {
    if (window.innerWidth > this.mobileBreakpoint) {
      if (this.currentState !== "desktop") {
        this.setDesktopLayout();
      }
    } else {
      if (this.currentState !== "mobile") {
        this.setMobileLayout();
      }
    }
  }

  setDesktopLayout() {
    if (this.desktopContainer && this.filterContent) {
      this.desktopContainer.appendChild(this.filterContent);
    }
    this.currentState = "desktop";
  }

  setMobileLayout() {
    if (this.mobileContainer && this.filterContent) {
      this.mobileContainer.appendChild(this.filterContent);
    }
    this.currentState = "mobile";
  }

  applyFilters() {
    // Get all filters and trigger the results htmx element
    // with newly created url
    const filters = document.querySelectorAll('[data-type]');
    const newURL = new URL(window.location.href);
    newURL.search = "";

    for (let f of filters) {
      let urlParam = f.dataset.type;
      let urlValue = filterFunctionMap[urlParam]?.(f);
      if (typeof urlValue === 'object') {
        urlValue.forEach(id => newURL.searchParams.append(urlParam, id));
      } else if (typeof urlValue === 'string') {
        newURL.searchParams.append(urlParam, urlValue);
      }
    }

    let resultsDiv = document.getElementById("results");
    resultsDiv.setAttribute("hx-get", newURL.toString());

    htmx.process(resultsDiv);
    htmx.trigger(resultsDiv, "filter-change");
  }

  mobileFilterToggle() {
    let isOpen = this.mobileFilters.classList.contains("opacity-100");
    if (isOpen) {
      this.mobileFilters.classList.remove("opacity-100", "scale-100", "pointer-events-auto");
      this.mobileFilters.classList.add("opacity-0", "scale-95", "pointer-events-none");
    } else {
      this.mobileFilters.classList.remove("opacity-0", "scale-95", "pointer-events-none");
      this.mobileFilters.classList.add("opacity-100", "scale-100", "pointer-events-auto");
    }
  }

  mobileFitlerExit(e) {
    let isOpen = this.mobileFilters.classList.contains("opacity-100");
    if (!isOpen) {
      return;
    }

    let clickMenu = this.mobileFilters.contains(e.target);
    let clickButton = this.mobileFilterButton.contains(e.target);
    let dropDowns = document.querySelectorAll(".dropdown");
    let clickDropDown = false;

    dropDowns.forEach((d) => {
      if (d.contains(e.target)) {
        clickDropDown = true;
      }
    })

    if (!(clickMenu || clickButton || clickDropDown)) {
      this.mobileFilters.classList.remove("opacity-100", "scale-100", "pointer-events-auto");
      this.mobileFilters.classList.add("opacity-0", "scale-95", "pointer-events-none");
    } }
}

let filterFunctionMap = {
  "sort" : (f) => {
    return f.value;
  },
  "person" : (f) => {
    return f.getFilterIds();
  },
  "creator" : (f) => {
    return f.getFilterIds();
  },
  "character" : (f) => {
    return f.getFilterIds();
  },
  "tag" : (f) => {
    return f.getFilterIds();
  }
}

if (!customElements.get("filter-content")) {
  customElements.define("filter-content", FilterContent);
}
