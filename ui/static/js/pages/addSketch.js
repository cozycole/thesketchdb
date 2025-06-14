import { UploadImagePreview } from "../components/uploadImagePreview.js";
import { FormSearchDropdown } from "../components/formSearchDropdown.js";

export function initAddSketch() {
  let imgPreview = new UploadImagePreview("vidThumbPrev");

  customElements.define("form-search", FormSearchDropdown);

  document.body.addEventListener("htmx:configRequest", function (evt) {
    // this adds the value of the triggering element to the query parameter of the
    // url request
    evt.detail.parameters["query"] = evt.detail.elt.value;
  });
}
