import { FormSearchDropdown } from "../components/formSearchDropdown.js";

export function initAddCategory() {
  const htmxDropdown = document.getElementById("categoryHtmxDropdown");
  let searchDropdown = new FormSearchDropdown(htmxDropdown);

  document.body.addEventListener("htmx:configRequest", function (evt) {
    // this adds the value of the triggering element to the query parameter of the
    // url request
    evt.detail.parameters["query"] = evt.detail.elt.value;
  });
}
