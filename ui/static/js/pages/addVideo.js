import { UploadImagePreview } from '../components/uploadImagePreview.js'
import { ActorForm } from '../components/actorForm.js'

export function initAddVideo() {
  let imgPreview = new UploadImagePreview('imagePreview');
  new ActorForm(
    'personInputs', 
    'addPersonButton', 
    'personInputTemplate',
    imgPreview);

  document.body.addEventListener("htmx:configRequest", function (evt) {
    evt.detail.parameters["query"] = evt.detail.elt.value;
  });

  function insertDropdownItem(e) {
      text = e.target.outerText
      id = e.target.dataset.id

      dropDownList = e.target.parentNode
      // dropdown list is contained in div
      searchInput = dropDownList.parentNode.previousElementSibling
      searchInput.value = text

      idInput = searchInput.previousElementSibling
      idInput.value = id

      dropDownList.remove()
  }

  // remove dropdown if user clicks outside of dropdown
  document.addEventListener("click", (e) => {
      const dropdown = document.getElementById('dropdown')
      if (!dropdown) {
          return
      }
      const input = dropdown.parentNode.previousElementSibling

      const isClickInside = input.contains(e.target) || dropdown.contains(e.target)
      if (!isClickInside) {
          dropdown.remove()
      }
  })
}
