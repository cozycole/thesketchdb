export class FormDropdown extends HTMLElement {
  constructor() {
    super();

    this.idInput = this.querySelector('input[type="hidden"]');
    this.selector = this.querySelector(".selector");
    this.dropdown = this.querySelector(".dropdown");

    if (!this.selector) {
      log.error("selector not defined!");
    }

    this.addEventListener("insertDropdownItem", (e) => {
      let dropDownItems = this.querySelectorAll("li.result");

      for (let el of dropDownItems) {
        el.addEventListener("click", (e) => {
          const target = e.currentTarget;
          const text = target.outerText;
          const id = target.dataset.id;
          const img = target.querySelector("img");
          setTimeout(() => this.insertDropdownItem(id, text, img), 100);
        });
      }
    });

    document.body.addEventListener("click", (e) => {
      if (
        !(
          this.dropdown?.contains(e.target) || this.selector?.contains(e.target)
        )
      ) {
        this.dropdown.innerHTML = "";
      }
    });

    document.addEventListener("keydown", (e) => {
      if (e.key === "Escape") {
        this.removeDropdown();
      }
    });
  }

  insertDropdownItem(id, text, img) {
    this.selector.querySelector(".text").innerHTML = text;
    this.idInput.value = id;

    if (img) {
      this.selector.querySelector("img").src = img.src;
    }

    this.removeDropdown();
  }

  removeDropdown() {
    this.dropdown.innerHTML = "";
  }
}

if (!customElements.get("form-dropdown")) {
  customElements.define("form-dropdown", FormDropdown);
}
