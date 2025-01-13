function toggleWatchNow() {
    const watchNow = document.getElementById('watchNow');
    watchNow.classList.toggle("hidden");
    watchNow.classList.toggle("flex");
    

    const watchNowIframe = document.getElementById('watchNowIframe');
    watchNowIframe.contentWindow.postMessage('{"event":"command","func":"stopVideo","args":""}', '*')
}

function toggleCastDropDown() {
  const typeDropToggleUp = document.getElementById('typeDropUp');
  const typeDropToggleDown = document.getElementById('typeDropDown');
  const castGallery = document.getElementById('castGallery');
  typeDropToggleUp.classList.toggle('hidden');
  typeDropToggleDown.classList.toggle('hidden');
  castGallery.classList.toggle('hidden');

}
