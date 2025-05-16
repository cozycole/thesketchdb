import { FavoriteButton } from '../components/favoriteButton.js'
import { YoutubeEmbed } from '../components/youtubeEmbed.js'

export function initViewVideo() {
  customElements.define("favorite-button", FavoriteButton);
  customElements.define("youtube-embed", YoutubeEmbed);
}
