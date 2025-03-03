import { initHome } from './pages/home.js'
import { initViewVideo } from './pages/viewVideo.js';
import { initAddVideo } from './pages/addVideo.js';
import { initUpdateVideo } from './pages/updateVideo.js';
import { initAddTag } from './pages/addTag.js';
import { initAddCategory } from './pages/addCategory.js';
import { initSearch } from './pages/search.js';
import { initBrowse } from './pages/browse.js';

(function() {
  const firstDiv = document.querySelector('main > div');
  const pageType = firstDiv ? firstDiv.dataset.page : 'No page attribute found!';

  console.log(`On ${pageType} page`);

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

    menuBtn.addEventListener('click', () => {
        mobileMenu.classList.remove('translate-x-full');
    });

    closeMenu.addEventListener('click', () => {
        mobileMenu.classList.add('translate-x-full');
    }
    );

    searchBtn.addEventListener('click', () => {
      mobileSearch.classList.toggle('hidden');
    });
  });

  switch (pageType) {
    case 'home':
      initHome();
      break;
    case 'browse':
      initBrowse();
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
  }

}());
