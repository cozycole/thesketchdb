import { UploadImagePreview } from '../components/uploadImagePreview.js'
import { FormSearchDropdown } from '../components/formSearchDropdown.js'

export function initUpdateVideo() {
  new UploadImagePreview('vidThumbPrev');
  new UploadImagePreview('actorThumbnailPreview');
  new UploadImagePreview('actorProfilePreview');

  const creatorHtmxDropdown = document.getElementById('creatorHtmxDropdown');
  const personHtmxDropdown = document.getElementById('personHtmxDropdown');
  const characterHtmxDropdown = document.getElementById('characterHtmxDropdown');
  new FormSearchDropdown(creatorHtmxDropdown);
  new FormSearchDropdown(personHtmxDropdown);
  new FormSearchDropdown(characterHtmxDropdown);

  const tagHtmxDropdowns = document.getElementsByClassName('tagHtmxDropdown');
  for (let drop of tagHtmxDropdowns) {
    new FormSearchDropdown(drop);
  }

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

  const addTagButton = document.getElementById('addTagButton');
  addTagButton.addEventListener("click", (e) => {
    const template = document.getElementById('tagInput');
    const node = template.content.firstElementChild.cloneNode(true);
    const newInput = document.querySelector('#tagTable tbody').appendChild(node);
    new FormSearchDropdown(newInput.querySelector('input[type="search"]').parentElement);
    htmx.process(document.getElementById('tagTable'));
  });

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
