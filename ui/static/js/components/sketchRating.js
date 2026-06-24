export class SketchRating extends HTMLElement {
  constructor() {
    super();
    this.rater = document.querySelector("#sketchRating");
    this.raterTitle = document.querySelector("#ratingTitle");
    this.raterSubtitle = document.querySelector("#ratingSubtitle");
    this.rateToggle = this.querySelector(".toggleRatingWidget");
    this.rateModal = document.querySelector("#ratingModal");
    this.ratingDisplay = document.querySelector("#ratingDisplay");
    this.rateFormInput = this.querySelector('input[name="rating"]');

    this.rateToggle.addEventListener("click", (e) => {
      this.rateModal.classList.toggle("hidden");
      this.rateModal.classList.toggle("flex");
    });

    document.addEventListener("click", (e) => {
      if (this.isOpen()) {
        if (e.target === this.rateModal || e.target.id === "ratingWrapper") {
          this.toggleModal();
        }
      }
    });

    document.addEventListener("scroll", () => {
      if (this.isOpen()) {
        this.toggleModal();
      }
    });

    window.addEventListener("resize", () => {
      if (this.isOpen()) {
        this.toggleModal();
      }
    });

    this.rater.addEventListener("change", (e) => {
      let rating = e.target.value;
      this.rateFormInput.value = rating;
    });

    document.addEventListener("change", (e) => {
      if (e.target === this.rater) {
        this.ratingDisplay.textContent = e.target.value;
      }
    });
  }

  isOpen() {
    return this.rateModal.classList.contains("flex");
  }

  toggleModal() {
    this.rateModal.classList.toggle("hidden");
    this.rateModal.classList.toggle("flex");
  }
}

if (!customElements.get("rating-widget")) {
  customElements.define("rating-widget", SketchRating);
}
