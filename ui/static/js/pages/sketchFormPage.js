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

  const tableBodies = document.querySelectorAll(".quoteTable tbody");
  for (const tb of tableBodies) {
    tb.addEventListener("click", function (evt) {
      const trashButton = evt.target.closest(".trashButton");
      if (trashButton) {
        const tbody = trashButton.closest("tbody");

        const row = evt.target.closest("tr");
        // this needs to be done to avoid an error
        setTimeout(() => {
          row.remove();
          if (tbody.innerHTML.trim() === "") {
            const template = document.getElementById("noQuoteRowTemplate");
            const emptyRow = template.content.cloneNode(true);

            tbody.appendChild(emptyRow);
          }
        }, 0);
      }
    });
  }

  // get Add Quote buttons to add quote row to respective quote table
  let addQuoteButtons = document.querySelectorAll(".addQuoteButton");
  addQuoteBtnListeners(addQuoteButtons);

  document.body.addEventListener("htmx:configRequest", function (evt) {
    // this adds the value of the triggering element to the query parameter of the
    // url request
    if (evt.detail.path.includes("search")) {
      evt.detail.parameters["query"] = evt.detail.elt.value;
    }
  });

  document.body.addEventListener("htmx:afterSwap", function (evt) {
    // Process formViewer to enable closing on off click
    console.log(evt.target);
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
    // add new quote btn functionality to newly swapped form
    if (
      evt.target.classList.contains("quoteForm") ||
      evt.target.id === "moments"
    ) {
      let btns = evt.target.querySelectorAll(".addQuoteButton");
      addQuoteBtnListeners(btns);
    }

    if (evt.target.tagName === "TBODY") {
      document.querySelector("#noTagRow")?.remove();
    }
  });
}

function addQuoteBtnListeners(qtBtns) {
  for (const btn of qtBtns) {
    btn.addEventListener("click", function (evt) {
      // go up to the enclosing form
      const form = btn.closest("div");
      const tableBody = form.querySelector(".quoteTable tbody");

      const template = document.getElementById("quoteRowTemplate");
      const newRow = template.content.cloneNode(true);

      const noQuoteRow = tableBody.querySelector("#noQuoteRow");
      if (noQuoteRow) noQuoteRow.parentElement.remove();

      // append the new row
      tableBody.appendChild(newRow);
      htmx.process(tableBody);
    });
  }
}
