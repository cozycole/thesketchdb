import { VideoCarousel } from '../components/videoCarousel.js'

export function initBrowse() {
  const browseSections = document.getElementById('browseSections');
  for (let section of browseSections.children) {
    new VideoCarousel(
      section.querySelector('.carousel'),
      section.querySelector('.carouselPrevBtn'),
      section.querySelector('.carouselNextBtn'),
    );
  }
}
