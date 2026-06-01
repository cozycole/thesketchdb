export class QuoteList extends HTMLElement {
  connectedCallback() {
    this.buttons = this.querySelectorAll("[data-sort]");
    this.container = this.querySelector("[data-quotes]");

    this.buttons.forEach((button) => {
      button.addEventListener("click", () => {
        this.setActive(button);
        this.sort(button.dataset.sort);
      });
    });
  }

  setActive(activeButton) {
    this.buttons.forEach((button) => {
      button.classList.remove(
        "bg-orange-500",
        "border-orange-500",
        "text-white",
        "shadow-sm",
      );

      button.classList.add(
        "text-slate-600",
        "border-slate-500",
        "hover:text-slate-900",
      );
    });

    activeButton.classList.remove(
      "text-slate-600",
      "border-slate-500",
      "hover:text-slate-900",
    );

    activeButton.classList.add(
      "bg-orange-500",
      "border-orange-500",
      "text-white",
      "shadow-sm",
    );
  }

  sort(mode) {
    const sorted = [...this.container.children].sort((a, b) => {
      if (mode === "top") {
        return Number(b.dataset.likeCount) - Number(a.dataset.likeCount);
      }

      return Number(a.dataset.timestamp) - Number(b.dataset.timestamp);
    });

    this.container.replaceChildren(...sorted);
  }
}

if (!customElements.get("quote-list")) {
  customElements.define("quote-list", QuoteList);
}
