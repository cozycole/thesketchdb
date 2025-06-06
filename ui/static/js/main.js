import { initHome } from './pages/home.js'
import { initViewVideo } from './pages/viewVideo.js';
import { initAddVideo } from './pages/addVideo.js';
import { initUpdateVideo } from './pages/updateVideo.js';
import { initAddTag } from './pages/addTag.js';
import { initAddCategory } from './pages/addCategory.js';
import { initSearch } from './pages/search.js';
import { initBrowse } from './pages/browse.js';
import { initUpdateShow } from './pages/updateShow.js';
import { initViewShow } from './pages/viewShow.js';
import { initViewSeason } from './pages/viewSeason.js';
import { initViewCatalog } from './pages/viewCatalog.js';
import { CollapsibleContent } from './components/collapseContent.js';

import './components/flashMessage.js';

(function() {
  const firstDiv = document.querySelector('main');
  const pageType = firstDiv ? firstDiv.dataset.page : 'No page attribute found!';
  console.log(pageType)

  const dropdownMenuButtons = document.querySelectorAll(".dropdownBtn");
  const dropdownMenus = document.querySelectorAll(".dropdownMenu");

  document.addEventListener("DOMContentLoaded", () => {
    dropdownMenuButtons.forEach(button => {
        button.addEventListener("click", function () {
            let dropdown = this.nextElementSibling;
            let isOpen = dropdown.classList.contains("opacity-100");

            // Close all dropdowns
            document.querySelectorAll(".dropdownMenu").forEach(menu => {
                menu.classList.remove("opacity-100", "scale-100", "pointer-events-auto");
                menu.classList.add("opacity-0", "scale-95", "pointer-events-none");
            });

            // Toggle the clicked dropdown
            if (!isOpen) {
                dropdown.classList.remove("opacity-0", "scale-95", "pointer-events-none");
                dropdown.classList.add("opacity-100", "scale-100", "pointer-events-auto");
            }
        });
    });

    // Close dropdown when clicking outside
    document.addEventListener("click", function (event) {
        let clickInside = false;
        dropdownMenus.forEach(e => {
          if (e.contains(event.target)) {
            clickInside = true;
          };
        })
        dropdownMenuButtons.forEach(e => {
          if (e.contains(event.target)) {
            clickInside = true;
          };
        })
        if (!clickInside) {
            document.querySelectorAll(".dropdownMenu").forEach(menu => {
                menu.classList.remove("opacity-100", "scale-100", "pointer-events-auto");
                menu.classList.add("opacity-0", "scale-95", "pointer-events-none");
            });
        }
    });

    const menuBtn = document.getElementById('mobileBtn');
    const mobileMenu = document.getElementById('mobileMenu');
    const closeMenu = document.getElementById('closeMenuBtn');
    const searchBtn = document.getElementById('mobileSearchBtn');
    const mobileSearch = document.getElementById('mobileSearch');
    const mobileBackground = document.getElementById('mobileBackground');
    const mobileInput = mobileSearch.querySelector('input');

    menuBtn.addEventListener('click', () => {
      mobileMenu.classList.remove('translate-x-full');
      mobileBackground.classList.remove('hidden');
    });

    closeMenu.addEventListener('click', () => {
      mobileMenu.classList.add('translate-x-full');
      mobileBackground.classList.add('hidden');
    });

    mobileBackground.addEventListener('click', () => {
      mobileMenu.classList.add('translate-x-full');
      mobileBackground.classList.add('hidden');
    });

    searchBtn.addEventListener('click', () => {
      mobileSearch.classList.toggle('-translate-y-full');
      mobileSearch.classList.toggle('translate-y-0');
      setTimeout(() => {
        if (mobileSearch.classList.contains('translate-y-0')) {
          // focus to end of the text
          let length = mobileInput.value.length;
          mobileInput.focus();
          mobileInput.setSelectionRange(length, length);
        } else {
          mobileSearch.querySelector('input').blur();
        }
      }, 300);
    });

    let clearMobileSearch = document.getElementById('clearMobileSearch');
    clearMobileSearch.addEventListener('click', () => {
      mobileInput.value = '';
      clearMobileSearch.classList.toggle('hidden', mobileInput.value.trim() === '');
    });

    mobileInput.addEventListener('input', () => {
      clearMobileSearch.classList.toggle('hidden', mobileInput.value.trim() === '');
    });

    let desktopSearch = document.getElementById('desktopSearch');
    let desktopInput = desktopSearch.querySelector('input');
    let clearDesktopSearch = document.getElementById('clearDesktopSearch');
    desktopInput.addEventListener('input', () => {
      clearDesktopSearch.classList.toggle('hidden', desktopInput.value.trim() === '');
    });

    clearDesktopSearch.addEventListener('click', () => {
      let input = desktopSearch.querySelector('input');
      input.value = '';
      clearDesktopSearch.classList.toggle('hidden', desktopInput.value.trim() === '');
    });
  });

  switch (pageType) {
    case 'home':
      initHome();
      break;
    case 'browse':
      initBrowse();
      break;
    case 'catalog':
      initViewCatalog();
      break;
    case 'view-video':
      initViewVideo();
      break;
    case 'search':
      initSearch();
      break;
    case 'add-video':
      initAddVideo();
      break;
    case 'update-video':
      initUpdateVideo();
      break;
    case 'add-tag':
      initAddTag();
      break;
    case 'add-category':
      initAddCategory();
      break;
    case 'view-show':
      initViewShow();
      break;
    case 'update-show':
      initUpdateShow();
      break;
    case 'view-season':
      initViewSeason();
      break;
  }

}());
