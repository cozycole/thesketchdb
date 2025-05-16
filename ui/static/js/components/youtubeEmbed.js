export class YoutubeEmbed extends HTMLElement {
  constructor() {
    super();

    this.embedDiv = this.querySelector("#watchNow");
    this.toggleButtons = document.getElementsByClassName("toggleVideo");
    this.watchNowIframe = this.querySelector("iframe");

    if (!(this.embedDiv && this.toggleButtons && this.watchNowIframe)) {
      throw Error('embed error');
    }

    for (let button of this.toggleButtons) {
      button.addEventListener('click', () => {
        this.embedDiv.classList.toggle("hidden");
        this.embedDiv.classList.toggle("flex");
        this.watchNowIframe.contentWindow.postMessage('{"event":"command","func":"stopVideo","args":""}', '*');
      })
    }
  }
}
