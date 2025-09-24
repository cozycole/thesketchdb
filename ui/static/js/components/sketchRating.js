export class SketchRating extends HTMLElement {
  constructor() {
    super();
    this.rater = document.querySelector("#sketchRating");
    this.raterTitle = document.querySelector("#ratingTitle");
    this.raterSubtitle = document.querySelector("#ratingSubtitle");
    this.rateToggle = this.querySelector(".toggleRatingWidget");
    this.rateModal = document.querySelector("#ratingModal");
    this.rateFormInput = this.querySelector('input[name="rating"]');

    this.rateToggle.addEventListener("click", (e) => {
      this.rateModal.classList.toggle("hidden");
      this.rateModal.classList.toggle("flex");
    });

    this.rateModal.addEventListener("click", (e) => {
      if (e.target === this.rateModal) {
        this.toggleModal();
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

    // functionality for displaying rating info
    //this.rater.addEventListener("wa-hover", (e) => {
    //  let rating = e.detail.value;
    //  this.setRatingTitles(rating);
    //});
    //
    //this.rater.addEventListener("mouseleave", (e) => {
    //  this.setRatingTitles(this.rater.value);
    //});
  }

  isOpen() {
    return this.rateModal.classList.contains("flex");
  }

  toggleModal() {
    this.rateModal.classList.toggle("hidden");
    this.rateModal.classList.toggle("flex");
  }

  //setRatingTitles(rating) {
  //  let ratingTitles = ratingMap[rating];
  //  if (ratingTitles) {
  //    this.raterTitle.textContent = ratingTitles.title;
  //    this.raterSubtitle.textContent = ratingTitles.subtitle;
  //  } else {
  //    this.raterTitle.textContent = "Rate this sketch";
  //    this.raterSubtitle.textContent = "";
  //  }
  //}
}

const ratingMap = {
  1.0: {
    title: "Abysmal",
    subtitle: "Unwatchable. No comedic value.",
  },
  2.0: {
    title: "Terrible",
    subtitle: "Very poor execution. Little to no humor.",
  },
  3.0: {
    title: "Weak",
    subtitle: "Poorly done. A faint attempt at humor, but ineffective.",
  },
  4.0: {
    title: "Subpar",
    subtitle: "Below average. A few attempts land, most do not.",
  },
  5.0: {
    title: "Mediocre",
    subtitle: "Average at best. Occasionally amusing, mostly forgettable.",
  },
  6.0: {
    title: "Decent",
    subtitle: "Acceptable quality. Some laughs, though uneven.",
  },
  7.0: {
    title: "Good",
    subtitle: "Solid effort. Generally funny with minor shortcomings.",
  },
  8.0: {
    title: "Very Good",
    subtitle: "Consistently funny. Well-crafted and effective.",
  },
  9.0: {
    title: "Excellent",
    subtitle: "High quality. Strong writing, memorable delivery.",
  },
  10.0: {
    title: "Outstanding",
    subtitle: "Exceptional. Near-perfect sketch comedy.",
  },
};

if (!customElements.get("rating-widget")) {
  customElements.define("rating-widget", SketchRating);
}
