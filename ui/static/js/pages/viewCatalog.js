import { CatalogFilter } from '../components/catalogFilter.js'
import { FilterContent } from '../components/filterMenu.js'

export function initViewCatalog() {
  // Get all filters and trigger the results htmx element
  // with newly created url
  let applyButton = document.querySelector("#filterApply") 
  applyButton.addEventListener("click", (e) => {
    console.log("clicked!")
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
    console.log(newURL.toString());
    resultsDiv.setAttribute("hx-get", newURL.toString());

    htmx.process(resultsDiv);
    htmx.trigger(resultsDiv, "filter-change");
  });

  document.addEventListener("DOMContentLoaded", () => {
    const mobileFilterButton = document.querySelector("#mobileFilterMenu");
    const mobileFilters = document.querySelector("#mobileFilters");
    mobileFilterButton.addEventListener("click", function () {
      let isOpen = mobileFilters.classList.contains("opacity-100");
      if (isOpen) {
        mobileFilters.classList.remove("opacity-100", "scale-100", "pointer-events-auto");
        mobileFilters.classList.add("opacity-0", "scale-95", "pointer-events-none");
      } else {
        mobileFilters.classList.remove("opacity-0", "scale-95", "pointer-events-none");
        mobileFilters.classList.add("opacity-100", "scale-100", "pointer-events-auto");
      }
    });

    // close menu on outside click
    document.addEventListener("click", (e) => {
      let isOpen = mobileFilters.classList.contains("opacity-100");
      if (!isOpen) {
        return;
      }

      let clickMenu = mobileFilters.contains(e.target);
      let clickButton = mobileFilterButton.contains(e.target);
      let dropDowns = document.querySelectorAll(".dropdown");
      let clickDropDown = false;
      dropDowns.forEach((d) => {
        if (d.contains(e.target)) {
          clickDropDown = true;
        }
      })

      if (!(clickMenu || clickButton || clickDropDown)) {
        mobileFilters.classList.remove("opacity-100", "scale-100", "pointer-events-auto");
        mobileFilters.classList.add("opacity-0", "scale-95", "pointer-events-none");
      }
    });
  });

  customElements.define("catalog-filter", CatalogFilter);
  customElements.define("filter-content", FilterContent);
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
  "tag" : (f) => {
    return f.getFilterIds();
  }
}
