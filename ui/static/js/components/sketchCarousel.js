export class SketchCarousel {
  constructor(carousel, prevBtn, nextBtn) {
    this.carousel = carousel;
    this.prevBtn = prevBtn;
    this.nextBtn = nextBtn;
    this.scrollAmount = 600;

    this.nextBtn.addEventListener("click", () => {
      this.carousel.scrollBy({ left: this.scrollAmount, behavior: "smooth" });
    });

    this.prevBtn.addEventListener("click", () => {
      this.carousel.scrollBy({ left: -this.scrollAmount, behavior: "smooth" });
    });
  }
}
