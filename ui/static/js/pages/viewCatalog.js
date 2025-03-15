import { CatalogFilter } from '../components/catalogFilter.js'

export function initViewCatalog() {
  document.querySelector("#sortDropdown").addEventListener("change", (e) => {
    let currentURL = new URL(window.location.href);
    currentURL.searchParams.set("sort", e.target.value);
    let newURL = currentURL.toString();
    
    let resultsDiv = document.getElementById("results");
    resultsDiv.setAttribute("hx-get", newURL);

    htmx.process(resultsDiv);
    htmx.trigger(resultsDiv, "filter-change");
  })
  customElements.define("catalog-filter", CatalogFilter);
}
