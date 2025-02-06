// class that implements search auto complete dropdown and fills in the input box on click
export class SearchDropdown {
    constructor({buttonClass, contentId, hideOnOffClick = false}) {
        // multiple buttons can toggle the content (open and close button)
        this.buttons = document.getElementsByClassName(buttonClass);
        this.content = document.getElementById(contentId);
        this.arrows = arrowClass ? document.getElementsByClassName(arrowClass) : undefined;
        this.hideOnOffClick = hideOnOffClick;
        console.log(this.buttons);
        this.init();
    }

    init() {
        for (let button of this.buttons) {
            button.addEventListener('click', (event) => this.toggleDropdown(event));
        }

        if (this.hideOnOffClick) {
            document.addEventListener('click', (event) => {
                if (!this.content.contains(event.target) && !this.content.classList.contains('hidden')) {
                    this.toggleDropdown(event);
                }
            });
        }
    }

    toggleDropdown(event) {
        this.content.classList.toggle('hidden');
        if (this.arrows) {
            this.toggleArrow();
        }
        event.stopPropagation();
    }

    toggleArrow() {
        for (let arrow of this.arrows) {
            arrow.classList.toggle('hidden');
        }
    }
}
