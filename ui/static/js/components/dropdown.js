export class Dropdown {
    constructor({buttonClass, contentId, arrowClass = undefined, hideOnOffClick = false}) {
        // multiple buttons can toggle the content (open and close button)
        this.buttons = document.getElementsByClassName(buttonClass);
        this.content = document.getElementById(contentId);
        this.arrows = arrowClass ? document.getElementsByClassName(arrowClass) : undefined;
        this.hideOnOffClick = hideOnOffClick;

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
