import { UploadImagePreview } from '../components/uploadImagePreview.js'
import { FormSearchDropdown } from '../components/formSearchDropdown.js'

export function initUpdateShow() {
  customElements.define("img-preview", UploadImagePreview);

  document.body.addEventListener("htmx:configRequest", function (evt) {
    // this adds the value of the triggering element to the query parameter of the 
    // url request
    evt.detail.parameters["query"] = evt.detail.elt.value;
  });

  // remove elements on DELETE since delete endpoints return flash messages
  // we always display those, but need to also remove the html of the deleted either
  // season or episode
  document.body.addEventListener("htmx:afterRequest", (e) => {
    // Check that it was a DELETE and was successful
    //console.log(e.detail.requestConfig.method);
    //console.log(e.detail.xhr.status);
    //console.log(e);
    if (e.detail.requestConfig.verb === "delete" && e.detail.xhr.status === 200) {
      const triggeringElement = e.detail.elt;
      const type = triggeringElement.dataset.type;
      if (type === "episode") {
        const seasonId = triggeringElement.dataset.sid;
        const episodeId = triggeringElement.dataset.eid;
        const episodeEl = document.getElementById(`s${seasonId}e${episodeId}`);
        if (episodeEl) episodeEl.remove();
      }
    }
  });
}
