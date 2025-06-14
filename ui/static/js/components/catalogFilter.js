export class CatalogFilter extends HTMLElement {
  constructor() {
    super();

    const template = document
      .getElementById("filterTemplate")
      .content.cloneNode(true);
    if (!this.children.length) {
      this.appendChild(template);
    }

    this.filterProfileTemplate = document.getElementById("filterProfile");

    this.input = this.querySelector(".filterSearch");

    this.input.setAttribute("hx-get", this.dataset.url);
    this.input.setAttribute("placeholder", this.dataset.placeholder);

    this.dropdown = this.querySelector(".dropdown");
    this.filtersDiv = this.querySelector(".filters");

    // close dropdown on click outside and escape
    document.body.addEventListener("click", (e) => {
      setTimeout(() => {
        if (
          !(this.dropdown.contains(e.target) || this.input.contains(e.target))
        ) {
          this.dropdown.innerHTML = "";
        }
      }, 100);
    });

    document.addEventListener("keydown", (e) => {
      if (e.key === "Escape") {
        this.dropdown.innerHTML = "";
        this.input.blur();
      }
    });

    // Once search results are inserted into the dropdown,
    // add the dropdownItemEvent to each child
    this.dropdown.addEventListener("htmx:afterSwap", (e) => {
      let dropdownItems = e.detail.target.children;

      // clear an empty response such that there isnt a gray
      // little line below on empty query
      if (!dropdownItems.length) {
        e.detail.target.innerHTML = "";
      }

      for (let child of dropdownItems) {
        this.dropdownItemEvent(child);
      }
    });

    // only add the query param for dropdown search
    this.input.addEventListener("htmx:configRequest", (e) => {
      if (e.detail.target.classList.contains("dropdown")) {
        if (!e.detail.elt.value) {
          this.dropdown.innerHTML = "";
          e.preventDefault();
          return;
        }
        e.detail.parameters["query"] = e.detail.elt.value;
      }
    });

    this.filterType = this.dataset.type;
    const selectedData = this.dataset.selected;
    if (selectedData) {
      try {
        this.selectedFilters = JSON.parse(selectedData);
      } catch (error) {
        console.error("Error parsing selected persons:", error);
      }
    }
  }

  connectedCallback() {
    for (let filter of this.selectedFilters) {
      if (!this.filterIsDisplayed(filter.id)) {
        this.displayFilter(filter.id);
      }
    }
  }

  dropdownItemEvent(ele) {
    ele.addEventListener("click", (e) => {
      setTimeout(() => {
        let currentURL = new URL(window.location.href);
        let id = ele.dataset.id;
        if (currentURL.searchParams.has(this.dataset.type, id)) {
          this.dropdown.innerHTML = "";
          this.input.value = "";
          return;
        }

        let img = ele.querySelector("img");
        let text = ele.querySelector("p").textContent;

        this.addFilter(id, text, img?.src);
        this.displayFilter(id);
        this.dropdown.innerHTML = "";
        this.input.value = "";
      }, 100);
    });
  }

  addFilter(id, text, imgSrc) {
    let filter = this.selectedFilters.filter((e) => e.id === id);
    // guard against adding the same filter twice
    if (!filter.length) {
      let newFilter = {};
      newFilter.id = id;
      newFilter.name = text;
      newFilter.image = imgSrc;

      this.selectedFilters.push(newFilter);
      return newFilter;
    }
  }

  displayFilter(id) {
    let filter = this.selectedFilters.filter((e) => e.id === id);
    if (!filter.length) {
      throw Error(`No filter with id ${id}`);
    }

    // guard against adding the same filter twice
    filter = filter[0];
    if (this.filterIsDisplayed(filter.id)) {
      return;
    }

    let filterElement = this.filterProfileTemplate.content.cloneNode(true);
    let img = filterElement.querySelector("img");
    if (this.dataset.displayImg === "false") {
      img.remove();
    } else {
      img.src = filter.image;
    }

    let title = filterElement.querySelector("h4");
    title.textContent = filter.name;

    let deleteButton = filterElement.querySelector("button");
    let insertedFilter = filterElement.firstElementChild;
    insertedFilter.dataset.id = id;

    this.filtersDiv.appendChild(insertedFilter);
    deleteButton.addEventListener("click", (e) => {
      setTimeout(() => {
        this.removeFilter(filter);
      }, 100);
    });

    console.log(this.filtersDiv.children);

    filter.applied = true;
  }

  removeFilter(filter) {
    let filterElement = null;
    for (let e of this.filtersDiv.children) {
      if (e.dataset.id == filter.id) {
        filterElement = e;
      }
    }

    if (!filterElement) {
      throw Error(`Cannot hide filter with id ${filter.id}`);
    }

    this.filtersDiv.removeChild(filterElement);
    this.selectedFilters = this.selectedFilters.filter(
      (item) => item.id !== filter.id,
    );
  }

  getFilterIds() {
    let ids = [];
    for (let jsonFilter of this.selectedFilters) {
      ids.push(jsonFilter.id);
    }
    return ids;
  }

  filterIsDisplayed(filterId) {
    for (let filter of this.filtersDiv.children) {
      if (filter.dataset.id == filterId) {
        return true;
      }
    }
    return false;
  }
}

if (!customElements.get("catalog-filter")) {
  customElements.define("catalog-filter", CatalogFilter);
}
