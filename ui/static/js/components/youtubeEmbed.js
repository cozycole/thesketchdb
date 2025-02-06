export class YoutubeEmbed {
  constructor(embedDivId, toggleButtonClass) {
    this.embedDiv = document.getElementById(embedDivId);
    this.toggleButtons = document.getElementsByClassName(toggleButtonClass);
    this.watchNowIframe = document.querySelector(`#${embedDivId} iframe`);

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
