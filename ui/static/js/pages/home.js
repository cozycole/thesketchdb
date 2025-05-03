import { VideoCarousel } from '../components/videoCarousel.js'
import Glide from '@glidejs/glide'

const isTouchDevice = window.matchMedia('(pointer: coarse)').matches;

var FixBoundPeek = function (Glide, Components, Events) {
  return {
    /**
     * Fix peek 'after' with 'bound' option.
     *
     * @param  {Number} translate
     * @return {Number}
     */
    modify (translate) {
      var isBound = Components.Run.isBound
      // future method from 'master'
      if (typeof isBound !== 'function') {
        isBound = function () {
          return Glide.isType('slider') && Glide.settings.focusAt !== 'center' && Glide.settings.bound
        }
      }

      if (isBound() && Components.Run.isEnd()) {
        const peek = Components.Peek.value

        if (typeof peek === 'object' && peek.after) {
          return translate - peek.after
        }

        return translate - peek
      }

      return translate
    }
  }
}

export function initHome() {
  new Glide('.glide', {
    type: 'carousel',
    startAt: 0,
    perView: 1,
    focusAt: "center",
    autoplay: 5000,
    animationDuration: 800,
    gap: 0,
  }).mount();

  let elements = document.getElementsByClassName('glideCarousel');
  if (isTouchDevice) {
    for (carousel of elements) {
      let track = carousel.querySelector('.gt');
      let slides = carousel.querySelector('.gs');
      track.classList.toggle('glide__track');
      track.classList.toggle('flex');
      track.classList.toggle('overflow-x-auto');
      slides.classList.toggle('glide__slides');
      // hide nav buttons
      carousel.querySelector('.glide__arrows').classList.toggle('hidden');
    }
  } else {
    for (carousel of elements) {
      if (carousel.classList.contains('sketch')) {
        new Glide(carousel, {
          type: 'slider',
          startAt: 0,
          perView: 5,
          perSwipe: '|',
          focusAt: 0,
          bound: true,
          gap: 8,
          peek: {before:0, after: 50},
          breakpoints: {
            1100: {
              perView: 4,
              peek: {before:0, after: 50},
            },
            850: {
              perView: 3,
              peek: {before:0, after: 50},
            },
            575: {
              perView: 2,
              peek: {before:0, after: 50},
            },
            400: {
              perView: 1,
              peek: {before:0, after: 50},
            }
          },
          keyboard: false,
          rewind: false,
        }).mutate([FixBoundPeek])
          .mount();
      } else {
        // profile card carousel
        new Glide(carousel, {
          type: 'slider',
          startAt: 0,
          perView: 6,
          perSwipe: '|',
          focusAt: 0,
          bound: true,
          gap: 8,
          peek: {before:0, after: 50},
          breakpoints: {
            1100: {
              perView: 5,
              peek: {before:0, after: 50},
            },
            850: {
              perView: 4,
              peek: {before:0, after: 50},
            },
            575: {
              perView: 3,
              swipeThreshold: 40,
              peek: {before:0, after: 50},
            }
          },
          keyboard: false,
          rewind: false,
        }).mutate([FixBoundPeek])
          .mount();

      }
    }
  }

  //const carouselSections = document.getElementsByClassName('carouselSection');
  //for (let section of carouselSections) {
  //  new VideoCarousel(
  //    section.querySelector('.carousel'),
  //    section.querySelector('.carouselPrevBtn'),
  //    section.querySelector('.carouselNextBtn'),
  //  );
  //}
}
