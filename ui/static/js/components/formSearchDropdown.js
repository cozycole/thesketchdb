// class that implements search auto complete dropdown
// and fills in the input box on click
export class FormSearchDropdown extends HTMLElement {
  constructor() {
    super();

    this.idInput = this.querySelector('input[type="hidden"]');
    this.searchInput = this.querySelector('input[type="search"]');
    this.dropdown = this.querySelector(".dropdown");

    this.addEventListener("insertDropdownItem", (e) => {
      let dropDownItems = this.querySelectorAll("li.result");

      for (let el of dropDownItems) {
        el.addEventListener("click", (e) => {
          const target = e.currentTarget;
          const text = target.outerText;
          const id = target.dataset.id;
          setTimeout(() => this.insertDropdownItem(id, text), 100);
        });
      }
    });

    document.body.addEventListener("htmx:configRequest", function (evt) {
      // this adds the value of the triggering element to the query parameter of the
      // url request
      evt.detail.parameters["query"] = evt.detail.elt.value;
    });

    document.body.addEventListener("click", (e) => {
      if (
        !(
          this.dropdown?.contains(e.target) ||
          this.searchInput?.contains(e.target)
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

  insertDropdownItem(id, text) {
    this.searchInput.value = text;
    this.idInput.value = id;

    this.removeDropdown();
  }

  removeDropdown() {
    this.dropdown.innerHTML = "";
    this.searchInput.blur();
  }
}

if (!customElements.get("form-search")) {
  customElements.define("form-search", FormSearchDropdown);
}
