// We assume that we will have the following html
// <img>
// <label><input></label>
// So label.prevElementSibling is where we will be placing the image
export class UploadImagePreview {
    constructor(labelId) {
        this.label = document.getElementById(labelId);
        if (!this.label) {
            throw Error(`No label found with id ${labelId}`);
        }

        this.input = this.label.querySelector('input[type=file]');
        this.previewImage(this.input);
        this.input.addEventListener('change', (e) => this.previewImage());
    }

    previewImage() {
        const file = this.input.files[0];

        if (file) {
            const reader = new FileReader();

            reader.onload = (e) => {
                let prevPreview = this.label.previousElementSibling;
                if (prevPreview && prevPreview.nodeName == "IMG") {
                    prevPreview.remove();
                }

                const img = document.createElement('img');
                img.src = e.target.result;
                img.style.maxWidth = '300px'; 
                img.style.maxWidth = '200px'; 

                this.label.before(img);
            };

            reader.readAsDataURL(file);
        } 
    }
}
