import { initHome } from './pages/home.js'
import { initViewVideo } from './pages/viewVideo.js';
import { initAddVideo } from './pages/addVideo.js';
import { initUpdateVideo } from './pages/updateVideo.js';
import { initAddTag } from './pages/addTag.js';
import { initAddCategory } from './pages/addCategory.js';
import { initSearch } from './pages/search.js';
import { initBrowse } from './pages/browse.js';

const firstDiv = document.querySelector('main > div');
const pageType = firstDiv ? firstDiv.dataset.page : 'No page attribute found!';

console.log(`On ${pageType} page`);

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
