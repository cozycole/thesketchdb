class ShowMore extends HTMLElement {
  connectedCallback() {
    if (this.initialized) return;
    this.initialized = true;

    this.maxHeight = Number(this.getAttribute("max-height")) || 240;
    this.expanded = false;

    this.classList.add("block");

    const content = document.createElement("div");
    content.className = "show-more-content";

    while (this.firstChild) {
      content.appendChild(this.firstChild);
    }

    this.appendChild(content);

    requestAnimationFrame(() => {
      if (content.scrollHeight <= this.maxHeight + 1) {
        return;
      }

      content.classList.add("relative", "overflow-hidden");
      content.style.maxHeight = `${this.maxHeight}px`;

      const fade = document.createElement("div");
      fade.className =
        "show-more-fade pointer-events-none absolute inset-x-0 bottom-0 h-16 bg-gradient-to-t from-white to-transparent";

      const button = document.createElement("button");
      button.type = "button";
      button.className =
        "show-more-button mt-3 block mx-auto text-sm font-bold text-slate-700 hover:text-slate-950";
      button.textContent = "Show more";

      content.appendChild(fade);
      this.appendChild(button);

      this.content = content;
      this.fade = fade;
      this.button = button;

      button.addEventListener("click", () => this.toggle());
    });
  }

  toggle() {
    this.expanded = !this.expanded;

    this.content.style.maxHeight = this.expanded
      ? `${this.content.scrollHeight}px`
      : `${this.maxHeight}px`;

    this.fade.hidden = this.expanded;
    this.button.textContent = this.expanded ? "Show less" : "Show more";
  }
}

if (!customElements.get("show-more")) {
  customElements.define("show-more", ShowMore);
}
