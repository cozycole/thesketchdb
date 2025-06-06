export class CollapsibleContent extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {
    const defaultState = this.getAttribute("default") || "open";
    this.isOpen = defaultState === "open";

    this.init();
    window.addEventListener("resize", () => this.updateHeight());
  }

  init() {
    if (!document.contains(this.toggleBtn)) {
      this.toggleBtn = this.querySelector(".toggleBtn");
      this.toggleBtn.addEventListener("click", () => this.toggleContent());
    }

    if (!document.contains(this.arrow)) {
      this.arrow = this.querySelector(".arrow");

    }

    if (!document.contains(this.content)) {
      this.content = this.querySelector(".content");
      this.content.style.maxHeight = this.content.scrollHeight + "px";
    }
    setTimeout(() => {
      this.content.style.maxHeight = this.isOpen ? this.content.scrollHeight + "px" : "0px";
    }, 0);

    this.arrow.classList.toggle("rotate-180", !this.isOpen);

    if (this.dataset.status === "closed") {
      this.toggleContent();
    }

    this.observeContentMutations();
  }

  toggleContent() {
    this.isOpen = this.content.style.maxHeight !== "0px";
    this.content.style.maxHeight = this.isOpen ?  "0px" : this.content.scrollHeight + "px";
    this.arrow.classList.toggle("rotate-180", this.isOpen);
    this.isOpen = this.content.style.maxHeight !== "0px";
  }

  updateHeight() {
    if (this.isOpen) {
      this.content.style.maxHeight = this.content.scrollHeight + "px";
    }
  }
  
  // if the content is mutated in any way we want to update the height to include
  // any new content AND reset all class attributes for the custom element
  observeContentMutations() {
    const observer = new MutationObserver((mutations) => {
      for (const mutation of mutations) {
        for (const node of mutation.addedNodes) {
          if (node.tagName === "IMG") {
            node.addEventListener("load", () => this.updateHeight());
          } else if (node.querySelectorAll) {
            // In case images are nested in fragments or other elements
            node.querySelectorAll("img").forEach(img =>
              img.addEventListener("load", () => this.updateHeight())
            );
          }
        }
      }

      // Run initial height update in case non-image content was added
      this.updateHeight();
    });

      observer.observe(this, {
      childList: true,
      subtree: true,
    });

    this.mutationObserver = observer;
  }
}

if (!customElements.get("collapse-content")) {
  customElements.define("collapse-content", CollapsibleContent);
}
