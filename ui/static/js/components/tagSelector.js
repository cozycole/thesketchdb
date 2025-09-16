export class TagSelector extends HTMLElement {
  constructor() {
    super();

    let tags = this.querySelectorAll("li[data-id]");
    for (const tag of tags) {
      if (tag.dataset.selected) {
        this.styleSelected(tag);
      }
    }

    this.addEventListener("click", (evt) => {
      const li = evt.target.closest("li[data-id]");
      if (!li) {
        return;
      }

      const id = li.dataset.id;
      const form = this.querySelector("form");
      if (!form) {
        return;
      }

      const hiddenInput = form.querySelector(
        `input[type="hidden"][value="${id}"]`,
      );

      if (hiddenInput) {
        // Deselect
        hiddenInput.remove();
        li.removeAttribute("data-selected");
        this.styleUnselected(li);
      } else {
        // Select
        const input = document.createElement("input");
        input.type = "hidden";
        input.name = "tag_ids[]";
        input.value = id;
        form.appendChild(input);

        li.setAttribute("data-selected", "true");
        this.styleSelected(li);
      }
    });
  }

  styleSelected(e) {
    e.classList.remove("bg-slate-200", "hover:bg-slate-300", "text-slate-950");
    e.classList.add("bg-orange-500", "text-white");
  }
  styleUnselected(e) {
    e.classList.remove("bg-orange-500", "text-white");
    e.classList.add("bg-slate-200", "hover:bg-slate-300", "text-slate-950");
  }
}

if (!customElements.get("tag-selector")) {
  customElements.define("tag-selector", TagSelector);
}
