// label must contain input of type file
export class UploadImagePreview {
    constructor(labelClass) {
        this.labelClass = labelClass;
        this.init();
    }

    init() {
        this.labels = document.getElementsByClassName(this.labelClass);
        if (!this.labels.length) {
            throw Error(`No label found with id ${labelClass}`)
        }

        this.addListeners();
    }

    refresh() {
        // there is functionality to add labeled inputs to the form,
        // so when you add one, you need to refresh this to include the
        // newly added one
        this.init();
    }

    addListeners() {
        for (let label of this.labels) {
            this.label = label;
            this.input = this.label.querySelector('input[type=file]');

            if (this.input.getAttribute('uploadPreview') !== 'true') {
                this.input.addEventListener('change', (e) => this.previewImage(e));
                this.input.setAttribute('uploadPreview', 'true');
            }
        }
    }

    previewImage(event) {
        const file = event.target.files[0];

        if (file) {
            const reader = new FileReader();

            reader.onload = (e) => {
                let prevPreview = this.label.nextElementSibling;
                if (prevPreview && prevPreview.nodeName == "IMG") {
                    prevPreview.remove();
                }

                const img = document.createElement('img');
                img.src = e.target.result;
                img.style.maxWidth = '300px'; 

                this.label.after(img);
            };

            reader.readAsDataURL(file);
        } 
    }
}
