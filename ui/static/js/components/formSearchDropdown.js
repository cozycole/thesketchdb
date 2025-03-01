// class that implements search auto complete dropdown 
// and fills in the input box on click
export class FormSearchDropdown {
    constructor(divElement) {
        if (!(divElement instanceof HTMLElement)) {
            throw new Error("Expected an HTML element as the constructor argument");
        }

        this.div = divElement;

        this.div.addEventListener('insertDropdownItem', (e) => {
            let dropDownItems = this.div.querySelectorAll('li.result');

            for (let el of dropDownItems) {
                el.addEventListener('click', (e) => this.insertDropdownItem(e));
            }
        });

        // Remove dropdown if user clicks outside
        document.addEventListener('click', (e) => {
            const dropdown = document.getElementById('dropdown');
            if (!dropdown) {
                return;
            }
            const input = dropdown.parentNode.previousElementSibling;

            const isClickInside = input.contains(e.target) || dropdown.contains(e.target);
            if (!isClickInside) {
                dropdown.remove();
            }
        });
    }

    insertDropdownItem(e) {
        const text = e.target.outerText;
        const id = e.target.dataset.id;

        let dropDownList = e.target.parentNode;
        let searchInput = dropDownList.parentNode.previousElementSibling;
        searchInput.value = text;

        let idInput = searchInput.previousElementSibling;
        idInput.value = id;

        dropDownList.remove();
    }
}
