export function initViewSeason() {
  document.body.addEventListener("htmx:configRequest", function (evt) {
    //if (evt.detail.path.includes("season")) {
    //  console.log(evt.detail.elt.value);
    //  evt.detail.path = `/${evt.detail.elt.value}?format=sub`;
    //  evt.detail.parameters = {};
    //}
  });
}
