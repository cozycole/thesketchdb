import { LikeButton } from '../components/likeButton.js'
import { YoutubeEmbed } from '../components/youtubeEmbed.js'

export function initViewVideo() {
    new LikeButton('likeButton');
    new YoutubeEmbed('watchNow', 'toggleVideo');
}
