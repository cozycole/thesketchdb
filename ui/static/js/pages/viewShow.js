
export function initViewShow() {
  document.body.addEventListener("htmx:configRequest", function (evt) {
    // this adds the value of the triggering element to the query parameter of the 
    // url request
    if (evt.detail.path.includes('season')) {
      evt.detail.path = evt.detail.path + `/${evt.detail.elt.value}`
      evt.detail.parameters = {};
    }
  });
}
