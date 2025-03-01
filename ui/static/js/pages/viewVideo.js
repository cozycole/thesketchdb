import { Dropdown } from '../components/dropdown.js'
import { LikeButton } from '../components/likeButton.js'
import { YoutubeEmbed } from '../components/youtubeEmbed.js'

export function initViewVideo() {
    new Dropdown({
        buttonClass: 'castButton',
        contentId: 'castGallery',
        arrowClass: 'galleryArrow',
    });
    new Dropdown({
        buttonClass: 'tagButton',
        contentId: 'tags',
        arrowClass: 'galleryArrow',
    });
    new LikeButton('likeButton');
    new YoutubeEmbed('watchNow', 'toggleVideo');
}
