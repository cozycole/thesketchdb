import { UploadImagePreview } from './uploadImagePreview.js'

export class ActorForm {
  constructor(divId, addButtonId, templateId, uploadPreview) {
    this.divId = divId;
    this.addButton = document.getElementById(addButtonId);
    this.templateId = templateId;
    this.uploadPreview = uploadPreview;

    this.addButton.addEventListener('click', (e) => this.addPersonDiv());

    let dropDownDivs = document.getElementsByClassName('htmxDropdown');
    for (let div of dropDownDivs) {
      // this event gets triggered after every htmx request for this given search
      div.addEventListener('insertDropdownItem', (e) => {
        let dropDownItems = div.querySelectorAll('li.result');

        for (let el of dropDownItems) {
          el.addEventListener('click', (e) => {
            this.insertDropdownItem(e);
          })
        }
      })
    }

    // remove dropdown if user clicks outside of dropdown
    document.addEventListener("click", (e) => {
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
    
  addPersonDiv() {
    let personInputDivs = document.querySelectorAll(`#${this.divId} > div`);

    const lastInputDivs = Array.from(personInputDivs).sort((a,b) => {
      a = a.getAttribute("id")
      b = b.getAttribute("id")
      return a-b
    })

    const lastInputDiv = lastInputDivs.pop();
    const lastInputDivId = lastInputDiv.getAttribute("id");
    const newInputNum = Number(lastInputDivId.match(/\d+/)[0]) + 1;

    const template = document.getElementById(this.templateId).content;
    const newInput = document.importNode(template, true);

    let regex = /\[\d+\]/;
    let formInputs = newInput.querySelectorAll("input");
    formInputs.forEach((e) => {
      e.name.replace(regex, `[${newInputNum}]`);
    });

    lastInputDiv.parentNode.insertBefore(newInput, this.addButton);

    // you need to register the newly created inputs with htmx 
    // for the dropdown search to work
    htmx.process(document.getElementById(this.divId));

    this.uploadPreview.refresh();
  }

  insertDropdownItem(e) {
      let text = e.target.innerText;
      let id = e.target.dataset.id;

      let dropDownList = e.target.parentNode;
      // dropdown list is contained in div
      let searchInput = dropDownList.parentNode.previousElementSibling;
      searchInput.value = text;

      let idInput = searchInput.previousElementSibling;
      idInput.value = id;

      dropDownList.remove();
  }
}
