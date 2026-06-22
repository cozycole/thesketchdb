class HorizontalFade extends HTMLElement {
  connectedCallback() {
    this.classList.add("relative");

    this.scrollEl = this.querySelector(".overflow-x-auto");
    if (!this.scrollEl) return;

    this.leftFade = document.createElement("div");
    this.rightFade = document.createElement("div");

    this.leftFade.className =
      "pointer-events-none absolute inset-y-0 left-0 w-14 bg-gradient-to-r from-white to-transparent opacity-0 transition-opacity";

    this.rightFade.className =
      "pointer-events-none absolute inset-y-0 right-0 w-14 bg-gradient-to-l from-white to-transparent opacity-0 transition-opacity";

    this.append(this.leftFade, this.rightFade);

    this.updateFades = this.updateFades.bind(this);

    this.scrollEl.addEventListener("scroll", this.updateFades);
    window.addEventListener("resize", this.updateFades);

    this.updateFades();
  }

  disconnectedCallback() {
    this.scrollEl?.removeEventListener("scroll", this.updateFades);
    window.removeEventListener("resize", this.updateFades);
  }

  updateFades() {
    const maxScroll = this.scrollEl.scrollWidth - this.scrollEl.clientWidth;

    this.leftFade.classList.toggle("opacity-100", this.scrollEl.scrollLeft > 1);
    this.leftFade.classList.toggle("opacity-0", this.scrollEl.scrollLeft <= 1);

    this.rightFade.classList.toggle(
      "opacity-100",
      this.scrollEl.scrollLeft < maxScroll - 1,
    );
    this.rightFade.classList.toggle(
      "opacity-0",
      this.scrollEl.scrollLeft >= maxScroll - 1,
    );
  }
}

if (!customElements.get("horizontal-fade")) {
  customElements.define("horizontal-fade", HorizontalFade);
}
