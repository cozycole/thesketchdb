import { UploadImagePreview } from '../components/uploadImagePreview.js'
import { FormSearchDropdown } from '../components/formSearchDropdown.js'

export function initUpdateVideo() {
  let imgPreview = new UploadImagePreview('imagePreview');
  let searchDropdown = new FormSearchDropdown('htmxDropdown');
  //new ActorForm(
  //  'personInputs', 
  //  'addPersonButton', 
  //  'personInputTemplate',
  //  imgPreview);

  document.body.addEventListener("htmx:configRequest", function (evt) {
    // this adds the value of the triggering element to the query parameter of the 
    // url request
    evt.detail.parameters["query"] = evt.detail.elt.value;
  });

}
