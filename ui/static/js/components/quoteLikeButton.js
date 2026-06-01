import { showToast } from "../utils/toast";

export class QuoteLikeButton extends HTMLElement {
  constructor() {
    super();

    this._onClick = () => this.like();

    this.liked = "liked" in this.dataset;
    this.quoteId = this.dataset.id;

    this.likeButton = this.querySelector("button");
    this.icon = this.querySelector("svg");
    this.likeCount = this.querySelector(".count");

    if (!(this.likeButton && this.icon)) {
      throw Error(`Like button error`);
    }
  }

  connectedCallback() {
    this.likeButton.addEventListener("click", this._onClick);
  }

  disconnectedCallback() {
    this.likeButton.removeEventListener("click", this._onClick);
  }

  async like() {
    const method = this.liked ? "DELETE" : "POST";
    const response = await fetch(
      `/api/v1/quotes/like?quoteId=${this.quoteId}`,
      {
        method: method,
        credentials: "include",
        redirect: "manual",
      },
    );

    if (response.status == 401) {
      showToast("Sign in to like quotes");
      return;
    }

    if (!response.ok) {
      showToast("Error liking quote");
      return;
    }

    this.toggleButtonState();
  }

  toggleButtonState() {
    this.liked = !this.liked;
    this.icon.classList.toggle("fill-none", !this.liked);
    this.icon.classList.toggle("stroke-current", !this.liked);
    this.icon.classList.toggle("stroke-[50]", !this.liked);
    this.likeButton.classList.toggle("text-slate-800", !this.liked);

    this.icon.classList.toggle("fill-current", this.liked);
    this.likeButton.classList.toggle("text-orange-500", this.liked);

    let likeDiff = this.liked ? 1 : -1;
    let current = Number(this.likeCount.textContent);
    this.likeCount.textContent = current + likeDiff;
  }
}

if (!customElements.get("quote-like-button")) {
  customElements.define("quote-like-button", QuoteLikeButton);
}
