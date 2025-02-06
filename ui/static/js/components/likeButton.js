// button must contain two elements with id 'likeButton' and 'unlikeButton' to toggle
export class LikeButton {
  constructor(buttonId) {
    this.likeButton = document.getElementById(buttonId);
    if (!this.likeButton) {
      throw Error(`Like button with id ${buttonId}`);
    }

    this.likeIcon = this.likeButton.querySelector('#likeIcon');
    this.unlikeIcon = this.likeButton.querySelector('#unlikeIcon');
    this.vidId = likeButton.dataset.id;
    if (!(this.vidId && this.likeIcon && this.unlikeIcon)) {
      throw Error(`Like button with id ${buttonId} error`);
    }

    this.likeButton.addEventListener('click', () =>{
      this.like();
    });
  }

  async like() {
    const isLiked = window.getComputedStyle(this.likeIcon).display === 'none';
    const method = isLiked ? 'DELETE' : 'POST';
    const response = await fetch(`/video/like/${this.vidId}`, {
        method: method,
        credentials: 'include',
        redirect: 'manual'
    });

    if (response.status == 0) {
      this.showPopup(this.likeButton, 'Sign in to like videos');
      return;
    }

    if (!response.ok) {
      this.showPopup(this.likeButton, 'Error liking video...');
      return;
    }

    this.likeIcon.style.display = isLiked ? 'inline' : 'none';
    this.unlikeIcon.style.display = isLiked ? 'none' : 'inline';
  }

  showPopup(element, message) {
    const popup = document.createElement('div');
    popup.textContent = message;
    popup.className = 'absolute -top-6 left-1/2 -translate-x-1/2 whitespace-nowrap transform bg-black text-white text-xs px-2 py-1 rounded opacity-0 transition-opacity duration-300';
    element.appendChild(popup);

    setTimeout(() => popup.classList.add('opacity-100'), 10);
    setTimeout(() => {
      popup.classList.remove('opacity-100');
      setTimeout(() => popup.remove(), 300); 
    }, 2000);
  }
}
