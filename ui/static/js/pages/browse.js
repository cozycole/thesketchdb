import { SketchCarousel } from "../components/sketchCarousel.js";

export function initBrowse() {
  const browseSections = document.getElementById("browseSections");
  for (let section of browseSections.children) {
    new SketchCarousel(
      section.querySelector(".carousel"),
      section.querySelector(".carouselPrevBtn"),
      section.querySelector(".carouselNextBtn"),
    );
  }
}
