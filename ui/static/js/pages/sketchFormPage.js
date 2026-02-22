import htmx from "htmx.org";
import Sortable from "sortablejs";

export function initSketchFormPage() {
  // Drag and drop sortable cast table
  htmx.onLoad(function (content) {
    var sortables = content.querySelectorAll(".sortable");

    for (var i = 0; i < sortables.length; i++) {
      var sortable = sortables[i];
      let onEnd = function (evt) {
        this.option("disabled", true);
      };
      if (sortable.id === "quoteRows") {
        onEnd = function (evt) {};
      }

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
        onEnd: onEnd,
      });

      // Re-enable sorting on the `htmx:afterSwap` event
      sortable.addEventListener("htmx:afterSwap", function () {
        sortableInstance.option("disabled", false);
      });
    }
  });

  document.body.addEventListener("htmx:afterSwap", function (evt) {
    // Process formViewer to enable closing on off click
    let formModal = document.body.querySelector("#formModal");
    if (formModal && evt.target.id === "formViewer") {
      formModal.addEventListener("click", (e) => {
        let menu = formModal.querySelector("div");
        let dropDown = document.querySelector(".dropdown");
        if (!(menu.contains(e.target) || dropDown.contains(e.target))) {
          formModal.classList.remove("flex");
          formModal.classList.add("hidden");
        }
      });
    }

    if (evt.target.tagName === "TBODY") {
      document.querySelector("#noTagRow")?.remove();
    }
  });
}
