export class CollapsibleContent extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {
    this.toggleBtn = this.querySelector(".toggleBtn");
    this.arrow = this.querySelector(".arrow");
    this.content = this.querySelector(".content");
    this.content.style.maxHeight = this.content.scrollHeight + "px";


    this.toggleBtn.addEventListener("click", () => this.toggleContent());
    window.addEventListener("resize", () => this.updateHeight());
  }

  toggleContent() {
    const isOpen = this.content.style.maxHeight !== "0px";
    this.content.style.maxHeight = isOpen ?  "0px" : this.content.scrollHeight + "px";
    this.arrow.classList.toggle("rotate-180", isOpen);
  }

  updateHeight() {
    const isOpen = this.content.style.maxHeight !== "0px";
    if (isOpen) {
      this.content.style.maxHeight = this.content.scrollHeight + "px";
    }

  }
}
