export class CatalogFilter extends HTMLElement {
  constructor() {
    super();

    const template = document.getElementById("filterTemplate").content.cloneNode(true);
    if (!this.children.length) {
      this.appendChild(template);
    }

    this.filterProfileTemplate = document.getElementById("filterProfile");

    this.input = this.querySelector(".filterSearch");
    this.dropdown = this.querySelector(".dropdown");
    this.filtersDiv = this.querySelector(".filters"); 

    // close dropdown on events
    document.body.addEventListener("click", (e) => {
      if (!(this.dropdown.contains(e.target) || this.input.contains(e.target))) {
        this.dropdown.innerHTML = '';
      }
    });

    document.addEventListener("keydown", (e) => {
      if (e.key === "Escape") {
        this.dropdown.innerHTML = '';
        this.input.blur();
      }
    });

    this.dropdown.addEventListener("htmx:afterSwap", (e) => {
      let dropdownItems = e.detail.target.children;
      for (const child of dropdownItems) {
        this.dropdownItemEvent(child);
      }
    });

    this.input.addEventListener("htmx:configRequest", (e) => {
      e.detail.parameters["query"] = e.detail.elt.value;
    })
  }

  connectedCallback() {
    console.log("connected");
    const selectedData = this.dataset.selected;
    console.log(selectedData);
    if (selectedData) {
      try {
        this.selectedFilters = JSON.parse(selectedData);
        this.updateUI();
      } catch (error) {
        console.error("Error parsing selected persons:", error);
      }
    }

    this.filterType = this.dataset.type;
  }

  disconnectedCallback() {
    console.log("Disconnected")
  }

  updateUI() {
    const currentURL = new URL(window.location.href);
    //console.log(currentURL);
    //console.log(this.selectedFilters);
    for (let filter of this.selectedFilters) {
      let validFilter = currentURL.searchParams.has(this.dataset.type, filter.id);
      if (!filter.applied && validFilter) {
        this.displayFilter(filter.id);
        filter.applied = true;
      } else if (filter.applied && !validFilter) {
        this.hideFilter(filter);
        filter.applied = false;
      }
    }
  }

  dropdownItemEvent(ele) {
    ele.addEventListener("click", (e) => {
      let currentURL = new URL(window.location.href);
      let id = ele.dataset.id;
      if (currentURL.searchParams.has(this.dataset.type, id)) {
        this.dropdown.innerHTML = '';
        this.input.value = "";
        return
      }

      let img = ele.querySelector("img");
      let text = ele.querySelector("p").textContent;

      this.addFilter(id, text, img.src);
      this.displayFilter(id);
      this.dropdown.innerHTML = '';
      this.input.value = "";

      currentURL.searchParams.append(this.dataset.type, id);
      currentURL.searchParams.set("page", 1);
      let newURL = currentURL.toString();

      let resultsDiv = document.getElementById("results");
      resultsDiv.setAttribute("hx-get", newURL);

      htmx.process(resultsDiv);
      htmx.trigger(resultsDiv, "filter-change");

    })
  }

  addFilter(id, text, imgSrc) {
    let filter = this.selectedFilters.filter(e => e.id === id);
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
    let filter = this.selectedFilters.filter(e => e.id === id);
    if (!filter.length) {
      throw Error(`No filter with id ${id}`);
    }

    filter = filter[0];
    if (filter.applied) {
      return
    }

    let filterElement = this.filterProfileTemplate.content.cloneNode(true);
    let img = filterElement.querySelector('img');
    img.src = filter.image;

    let title = filterElement.querySelector('h4');
    title.textContent = filter.name;

    let deleteButton = filterElement.querySelector('button');
    let insertedFilter = filterElement.firstElementChild;
    insertedFilter.dataset.id = id;

    this.filtersDiv.appendChild(insertedFilter);
    deleteButton.addEventListener("click", (e) => {
      this.hideFilter(filter);
    });

    filter.applied = true;
  }

  hideFilter(filter) {
    let filterElement = null;
    for (let e of this.filtersDiv.children) {
      if (e.dataset.id == filter.id) {
        filterElement = e;
      }
    }

    if (!filterElement) {
      throw Error(`Cannot hide filter with id ${filter.id}`);
    }

    let currentURL = new URL(window.location.href);
    currentURL.searchParams.delete(this.dataset.type, filter.id);
    let newURL = currentURL.toString();

    let resultsDiv = document.getElementById("results");
    resultsDiv.setAttribute("hx-get", newURL);

    htmx.process(resultsDiv);
    htmx.trigger(resultsDiv, "filter-change");

    this.filtersDiv.removeChild(filterElement);
    filter.applied = false;
  }
}
