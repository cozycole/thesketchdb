// class that implements search auto complete dropdown 
// and fills in the input box on click
export class FormSearchDropdown {
    // div contains two inputs, one hidden above a visible one
    constructor(divId) {
        console.log(`FormSearchDropdown constructed with div ${divId}`);
        this.div = document.getElementById(divId);
        // this custom event gets triggered after every 
        // htmx request for this given search based on the server
        // response header value
        this.div.addEventListener('insertDropdownItem', (e) => {
          let dropDownItems = this.div.querySelectorAll('li.result');

          for (let el of dropDownItems) {
            el.addEventListener('click', (e) => {
              
              this.insertDropdownItem(e);
            })
          }
        })

        // remove dropdown if user clicks outside of dropdown
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
    // Assumed that there is a second hidden
    // input that's the previousElementSibiling
    // of the visible user input
    insertDropdownItem(e) {
      const text = e.target.outerText;
      const id = e.target.dataset.id;

      let dropDownList = e.target.parentNode;
      // dropdown list is contained in div
      let searchInput = dropDownList.parentNode.previousElementSibling;
      searchInput.value = text;

      let idInput = searchInput.previousElementSibling;
      idInput.value = id;

      dropDownList.remove();
    }
}
