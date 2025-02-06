import { initViewVideo } from './pages/viewVideo.js';
import { initSearch } from './pages/search.js';
import { initAddVideo } from './pages/addVideo.js';

const firstDiv = document.querySelector('main > div')
const pageType = firstDiv ? firstDiv.dataset.page : 'No page attribute found!';

console.log(`On ${pageType} page`);

switch (pageType) {
  case 'view-video':
    initViewVideo();
    break;
  case 'search':
    initSearch();
    break;
  case 'addVideo':
    initAddVideo();
    break;
}
