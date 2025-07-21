// the goal here is to create a dropdown that triggers an
// htmx hx-get of new content when new dropdown item is clicked
export class SelectDropdown extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {
    this.button = this.querySelector("[data-dropdown-button]");
    this.buttonText = this.button.querySelector("button");
    this.list = this.querySelector("[data-dropdown-list]");
    this.items = this.querySelectorAll("[data-dropdown-item]");
    this.defaultLabel = this.getAttribute("default-label") || "Select";
    this.targetSelector = this.getAttribute("hx-target");

    if (this.button) {
      this.button.addEventListener("click", (e) => {
        e.stopPropagation();
        this.list?.classList.toggle("hidden");
      });
    }

    this.items.forEach((item) => {
      item.addEventListener("click", () => this.selectItem(item));
    });

    document.addEventListener("click", (e) => {
      if (!this.list.contains(e.target)) {
        this.list?.classList.add("hidden");
      }
    });
  }

  selectItem(item) {
    const label = item.getAttribute("data-label");
    const url = item.getAttribute("data-url");

    if (this.buttonText && label) {
      this.buttonText.textContent = label;
    }

    this.list?.classList.add("hidden");
  }
}

if (!customElements.get("select-dropdown")) {
  customElements.define("select-dropdown", SelectDropdown);
}
