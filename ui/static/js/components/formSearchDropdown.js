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
        el.addEventListener("click", (e) => this.insertDropdownItem(e));
      }
    });

    document.body.addEventListener("click", (e) => {
      if (
        !(
          this.dropdown.contains(e.target) ||
          this.searchInput.contains(e.target)
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

  insertDropdownItem(e) {
    const text = e.currentTarget.outerText;
    const id = e.currentTarget.dataset.id;

    this.searchInput.value = text;
    this.idInput.value = id;

    this.removeDropdown();
  }

  removeDropdown() {
    this.dropdown.innerHTML = "";
    this.searchInput.blur();
  }
}
