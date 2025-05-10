export class FilterContent extends HTMLElement {
  constructor() {
    super();

    this.filterContent = document.getElementById("filters")
    this.desktopContainer = document.getElementById("desktopFilters");
    this.mobileContainer = document.getElementById("mobileFilters");

    this.mobileBreakpoint = 1024;
    window.addEventListener("resize", () => this.updateLayout());
  }

  connectedCallback() {
    this.currentState = "";
    this.updateLayout();
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
}
