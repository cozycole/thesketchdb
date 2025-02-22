import { UploadImagePreview } from '../components/uploadImagePreview.js'
import { FormSearchDropdown } from '../components/formSearchDropdown.js'

export function initUpdateVideo() {
  new UploadImagePreview('vidThumbPrev');
  new UploadImagePreview('actorThumbnailPreview');
  new UploadImagePreview('actorProfilePreview');
  new FormSearchDropdown('creatorHtmxDropdown');
  new FormSearchDropdown('personHtmxDropdown');
  new FormSearchDropdown('characterHtmxDropdown');


  let addCastButton = document.getElementById('addCastButton');
  let formViewer = document.getElementById('addCastFormViewer');

  // display addCastForm
  addCastButton.addEventListener('click', (e) => {
    formViewer.classList.toggle('hidden');
    formViewer.classList.toggle('flex');

  });

  // hide addCastForm
  formViewer.addEventListener('click', (e) => {
    if (e.target === e.currentTarget) {
      formViewer.classList.toggle('hidden');
      formViewer.classList.toggle('flex');
    }
  })

  document.body.addEventListener("htmx:configRequest", function (evt) {
    // this adds the value of the triggering element to the query parameter of the 
    // url request
    if (evt.detail.path.includes('search')) {
      evt.detail.parameters["query"] = evt.detail.elt.value;
    }
  });

  document.body.addEventListener("htmx:afterSwap", function (evt) {
    //console.log(evt.detail.target);
    if (evt.detail.target.id === 'castForm') {
      new UploadImagePreview('actorThumbnailPreview');
      new UploadImagePreview('actorProfilePreview');
    }


    if (evt.detail.target.id === 'castTable') {
      formViewer.click();

      const template = document.getElementById('actorFormInputs');
      const form = document.getElementById('addCastForm');

      if (template && form) {
        const clonedContent = template.content.cloneNode(true);
        form.innerHTML = ""; 
        form.appendChild(clonedContent);
        new UploadImagePreview('actorThumbnailPreview');
        new UploadImagePreview('actorProfilePreview');
      }
    }
  });
}
