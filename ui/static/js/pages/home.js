import { VideoCarousel } from '../components/videoCarousel.js'
import Glide from '@glidejs/glide'

export function initHome() {
  new Glide('.glide', {
    type: 'carousel',
    startAt: 0,
    perView: 1,
    focusAt: "center",
    autoplay: 5000,
    gap: 0,
    //autoplay: 5000
  }).mount();
  const carouselSections = document.getElementsByClassName('carouselSection');
  for (let section of carouselSections) {
    new VideoCarousel(
      section.querySelector('.carousel'),
      section.querySelector('.carouselPrevBtn'),
      section.querySelector('.carouselNextBtn'),
    );
  }
}
