import { Dropdown } from '../components/dropdown.js'

export function initSearch() {
    new Dropdown({
        buttonClass: 'filterMenuButton', 
        contentId: 'filterMenu', 
        hideOnOffClick: true
    });
    new Dropdown({
        buttonClass: 'filterTypeButton', 
        contentId: 'typeMenu',
        arrowClass: 'filterTypeArrow'
    });

    document.body.addEventListener('htmx:configRequest', function (evt) {
        //console.log(evt.detail.headers);
        if (evt.detail.elt.classList.contains('htmxSearchPage')) {
            const searchType = document.querySelector('input[name="type"]:checked').value;
            evt.detail.parameters['type'] = searchType;
        }
    });
};
