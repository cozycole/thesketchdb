// deletes img after input if present, then inserts an img directly after input
function previewImage(event) {
    const input = event.target
    const file = event.target.files[0];

    if (file) {
        const reader = new FileReader();

        reader.onload = function(e) {
            let prevPreview = input.nextElementSibling;
            if (prevPreview && prevPreview.nodeName == "IMG") {
                prevPreview.remove()
            }

            const img = document.createElement('img');
            img.src = e.target.result;
            img.style.maxWidth = '300px'; 

            let inputParent =  input.parentNode;
            inputParent.insertBefore(img, input.nextSibling)
        };

        reader.readAsDataURL(file);
    }
}

document.getElementById("addPersonButton").addEventListener("click", () => {
    const personInputDivs = document.querySelectorAll("#personInputs > div")
    const lastInputDivs = Array.from(personInputDivs).sort((a,b) => {
        a = a.getAttribute("id")
        b = b.getAttribute("id")
        return a-b
    })
    console.log(lastInputDivs)


    const lastInputDiv = lastInputDivs.pop()
    // We are just trying increment all the name attributes of each input
    // input[0] was the last, so now we need to name them all input[1]
    const lastInputDivId = lastInputDiv.getAttribute("id")
    let newInputNum = Number(lastInputDivId.match(/\d+/)[0]) + 1

    const template = document.getElementById("personInputTemplate").content;
    const newInput = document.importNode(template, true);
    const addPersonButton = document.getElementById("addPersonButton")

    newInput.firstElementChild.id = `people[${newInputNum}]`
    let formInputs = newInput.querySelectorAll("input")
    formInputs.forEach((e) => {
        switch (e.name) {
            case "peopleId":
                e.name = `peopleId[${newInputNum}]`
                break
            case "peopleText":
                e.name = `peopleText[${newInputNum}]`
                break
            case "characterId":
                e.name = `characterId[${newInputNum}]`
                break
            case "characterText":
                e.name = `characterText[${newInputNum}]`
                break
            case "characterThumbnail":
                e.name = `characterThumbnail[${newInputNum}]`
                break
        }
    })

    lastInputDiv.parentNode.insertBefore(newInput, addPersonButton)

    htmx.process(document.getElementById("personInputs"))
})

function insertDropdownItem(e) {
    text = e.target.outerText
    id = e.target.dataset.id

    dropDownList = e.target.parentNode
    // dropdown list is contained in div
    searchInput = dropDownList.parentNode.previousElementSibling
    searchInput.value = text

    idInput = searchInput.previousElementSibling
    idInput.value = id

    dropDownList.remove()
}

// remove dropdown if user clicks outside of dropdown
document.addEventListener("click", (e) => {
    const dropdown = document.getElementById('dropdown')
    if (!dropdown) {
        return
    }
    const input = dropdown.parentNode.previousElementSibling

    const isClickInside = input.contains(e.target) || dropdown.contains(e.target)
    if (!isClickInside) {
        dropdown.remove()
    }
})
