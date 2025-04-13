class FlashMessage extends HTMLElement {
  connectedCallback() {
    const duration = parseInt(this.getAttribute('duration'), 10) || 3000;

    setTimeout(() => {
      this.classList.add('fade-out');
      setTimeout(() => {
        this.remove();
      }, 500); // match CSS transition
    }, duration);
  }
}

if (!customElements.get('flash-message')) {
  customElements.define('flash-message', FlashMessage);
}
