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

async function toggleLike(videoId) {
  const likeButton = document.getElementById('likeButton');
  const unlikeButton = document.getElementById('unlikeButton');
  const buttonDiv = document.getElementById('buttonDiv');
  const isLiked = window.getComputedStyle(likeButton).display === 'none';

  try {
    const method = isLiked ? 'DELETE' : 'POST';
    const response = await fetch(`/video/like/${videoId}`, {
        method: method,
        credentials: 'include',
        redirect: 'manual'
    });

    if (response.status == 0) {
      throw new Error("Sign in to like videos");
    }

    if (!response.ok) {
      throw new Error(`Network error ${response.status}: ${response.body}`);
    }

    likeButton.style.display = isLiked ? 'inline' : 'none';
    unlikeButton.style.display = isLiked ? 'none' : 'inline';
  } catch (error) {
      showPopup(buttonDiv, "Sign in to like videos")
  }
}

function showPopup(element, message) {
  const popup = document.createElement('div');
  popup.textContent = message;
  popup.className = 'absolute -top-6 left-1/2 -translate-x-1/2 whitespace-nowrap transform bg-black text-white text-xs px-2 py-1 rounded opacity-0 transition-opacity duration-300';
  element.appendChild(popup);
  console.log("appended!")

  setTimeout(() => popup.classList.add('opacity-100'), 10);
  setTimeout(() => {
    popup.classList.remove('opacity-100');
    setTimeout(() => popup.remove(), 300); 
  }, 2000);
}
