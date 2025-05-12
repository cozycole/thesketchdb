import { CatalogFilter } from '../components/catalogFilter.js'
import { FilterContent } from '../components/filterMenu.js'

export function initViewCatalog() {
  const tabContainer = document.getElementById("tabContainer");
  const tabUnderline = document.getElementById("tabUnderline");
  const tabs = document.querySelectorAll(".tab");

  // Set active tab
  const path = window.location.pathname;
  const match = path.match(/^\/catalog\/(\w+)/);
  const currentTab = match ? match[1] : null;

  const activeTab = document.querySelector(".tab.active");
  if (activeTab) updateUnderline(activeTab);

  // On tab click
  tabs.forEach((tab) => {
    if (tab.dataset.tab === currentTab) styleTabs(tab, tabs);

    tab.addEventListener("click", () => {
      styleTabs(tab, tabs)
    });
  });

  customElements.define("catalog-filter", CatalogFilter);
  customElements.define("filter-content", FilterContent);
}

function styleTabs(activeTab, tabs) {
  tabs.forEach((t) => t.classList.remove("text-orange-600", "font-bold", "active",));
  activeTab.classList.add("text-orange-600", "font-bold", "active");
  activeTab.classList.remove("hover:text-slate-700")
  updateUnderline(activeTab);
}

function updateUnderline(tab) {
  const containerRect = tabContainer.getBoundingClientRect();
  const tabRect = tab.getBoundingClientRect();
  const left = tabRect.left - containerRect.left + tabContainer.scrollLeft;

  tabUnderline.style.left = `${left}px`;
  tabUnderline.style.width = `${tab.offsetWidth}px`;
}
