// button must contain two elements with id 'likeButton' and 'unlikeButton' to toggle
export class FavoriteButton extends HTMLElement {
  constructor() {
    super();

    this._onClick = () => this.like();

    this.favorited = this.dataset.favorited.trim() === "true" ? true : false;
    this.sketchId = this.dataset.id;

    this.favButton = this.querySelector("#favBtn");
    this.icon = this.querySelector("#favIcon");
    this.favText = this.querySelector("#favText");

    if (!(this.favButton && this.icon && this.favText)) {
      throw Error(`Favorite button error`);
    }
  }

  connectedCallback() {
    this.favButton.addEventListener("click", this._onClick);
  }

  disconnectedCallback() {
    this.favButton.removeEventListener("click", this._onClick);
  }

  async like() {
    console.log("button clicked, favorite status: ", this.favorited);
    const method = this.favorited ? "DELETE" : "POST";
    const response = await fetch(`/sketch/like/${this.sketchId}`, {
      method: method,
      credentials: "include",
      redirect: "manual",
    });

    if (response.status == 0) {
      this.showPopup(this.favButton, "Sign in to favorite sketches");
      return;
    }

    if (!response.ok) {
      this.showPopup(this.favButton, "Error favoriting sketch...");
      return;
    }

    this.toggleButtonState();
  }

  toggleButtonState() {
    this.favorited = !this.favorited;

    this.icon.innerHTML = this.favorited
      ? `<path fill="currentColor" d="M3.172 5.172a4 4 0 015.656 0L12 8.343l3.172-3.171a4 4 0 115.656 5.656L12 21.343 3.172 10.828a4 4 0 010-5.656z"/>`
      : `<path stroke-linecap="round" stroke-linejoin="round" stroke="currentColor" fill="none"
          d="M3.172 5.172a4 4 0 015.656 0L12 8.343l3.172-3.171a4 4 0 115.656 5.656L12 21.343 3.172 10.828a4 4 0 010-5.656z" />`;
    this.favText.textContent = this.favorited ? "Favorited" : "Favorite";
  }

  showPopup(element, message) {
    const popup = document.createElement("div");
    popup.textContent = message;
    popup.className =
      "absolute -top-6 left-1/2 -translate-x-1/2 whitespace-nowrap transform bg-slate-900 text-white text-xs px-2 py-1 rounded opacity-0 transition-opacity duration-300";
    element.appendChild(popup);

    setTimeout(() => popup.classList.add("opacity-100"), 10);
    setTimeout(() => {
      popup.classList.remove("opacity-100");
      setTimeout(() => popup.remove(), 300);
    }, 2000);
  }
}
