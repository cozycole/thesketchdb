export function initViewSeason() {
  document.body.addEventListener("htmx:configRequest", function (evt) {
    if (evt.detail.path.includes("season")) {
      evt.detail.path = evt.detail.path + `/${evt.detail.elt.value}?format=sub`;
      evt.detail.parameters = {};
    }
  });
}
