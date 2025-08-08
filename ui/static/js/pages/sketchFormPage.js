import htmx from "htmx.org";
import Sortable from "sortablejs";

export function initSketchFormPage() {
  // Drag and drop sortable cast table
  htmx.onLoad(function (content) {
    var sortables = content.querySelectorAll(".sortable");
    for (var i = 0; i < sortables.length; i++) {
      var sortable = sortables[i];
      var sortableInstance = new Sortable(sortable, {
        animation: 150,
        ghostClass: "bg-slate-300",
        dragClass: "bg-white",
        handle: ".drag-icon",

        // Make the `.htmx-indicator` unsortable
        filter: ".htmx-indicator",
        onMove: function (evt) {
          return evt.related.className.indexOf("htmx-indicator") === -1;
        },

        // Disable sorting on the `end` event
        onEnd: function (evt) {
          this.option("disabled", true);
        },
      });

      // Re-enable sorting on the `htmx:afterSwap` event
      sortable.addEventListener("htmx:afterSwap", function () {
        sortableInstance.option("disabled", false);
      });
    }
  });

  document.body.addEventListener("htmx:configRequest", function (evt) {
    // this adds the value of the triggering element to the query parameter of the
    // url request
    if (evt.detail.path.includes("search")) {
      evt.detail.parameters["query"] = evt.detail.elt.value;
    }
  });

  document.body.addEventListener("htmx:afterSwap", function (evt) {
    // Process formViewer to enable closing on off click
    let formViewer = document.body.querySelector("#formViewer");
    if (formViewer && evt.target.id === "castFormViewer") {
      formViewer.addEventListener("click", (e) => {
        let menu = formViewer.querySelector("div");
        let dropDown = document.querySelector(".dropdown");
        if (!(menu.contains(e.target) || dropDown.contains(e.target))) {
          formViewer.classList.remove("flex");
          formViewer.classList.add("hidden");
        }
      });
    }
    // Hide modal if there's been a swap into the castTable
    if (formViewer && evt.target.id === "castTable") {
      formViewer.classList.remove("flex");
      formViewer.classList.add("hidden");
    }

    if (evt.target.tagName === "TBODY") {
      document.querySelector("#noTagRow")?.remove();
    }
  });
}
