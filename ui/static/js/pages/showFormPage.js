export function initShowFormPage() {
  document.body.addEventListener("htmx:afterSwap", function (evt) {
    // Process formViewer to enable closing on off click
    let formViewer = document.body.querySelector("#formModal");
    console.log(evt.target.id);
    if (formViewer && evt.target.id === "episodeFormViewer") {
      formViewer.addEventListener("click", (e) => {
        let menu = formViewer.querySelector("div");
        let dropDown = document.querySelector(".dropdown");
        if (!(menu.contains(e.target) || dropDown?.contains(e.target))) {
          formViewer.classList.remove("flex");
          formViewer.classList.add("hidden");
        }
      });
    }
    // Hide modal if there's been a swap into the castTable
    if (formViewer && evt.target.id?.includes("EpisodeTable")) {
      formViewer.classList.remove("flex");
      formViewer.classList.add("hidden");
    }
  });
}
