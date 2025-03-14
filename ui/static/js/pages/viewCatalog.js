import { CatalogFilter } from '../components/catalogFilter.js'

export function initViewCatalog() {
  customElements.define("catalog-filter", CatalogFilter);
}
