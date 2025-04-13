export class UploadImagePreview extends HTMLElement {
  constructor() {
    super();
    this.preview = this.querySelector('.imagePreview');
    if (!this.preview) {
      throw Error(`No image preview div`);
    }

    if (this.preview.firstElementChild) {
      this.originalImg = this.preview.firstElementChild;
    }

    this.input = this.querySelector('input[type=file]');
    this.input.addEventListener('change', (e) => this.previewImage());

    this.removeButton = this.querySelector('.remove');
    if (this.removeButton) {
      this.removeButton.addEventListener('click', () => {
        this.input.value = "";

        this.removePreviewImages();

        if (this.originalImg) {
          this.originalImg.style.display = "";
        }
      })
    }
  }

  connectedCallback() {
    this.input.value = "";
  }

  previewImage() {
    const file = this.input.files[0];

    if (file) {
      const reader = new FileReader();

      reader.onload = (e) => {
        this.removePreviewImages();

        const img = document.createElement('img');
        img.src = e.target.result;
        img.style.maxWidth = '300px'; 
        img.style.maxHeight = '300px'; 

        this.preview.appendChild(img);
      };

      reader.readAsDataURL(file);
    } 
  }

  removePreviewImages() {
    for (const child of this.preview.children) {
      if (!child.isEqualNode(this.originalImg)) {
        child.remove();
      }
      
      if (this.originalImg) {
        this.originalImg.style.display = "none";
      }
    }
  }
}
