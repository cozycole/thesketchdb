function previewImage(event) {
    const file = event.target.files[0];

    if (file) {
        const reader = new FileReader();

        reader.onload = function(e) {
            const imagePreview = document.getElementById('imagePreview');
            imagePreview.innerHTML = ''; // Clear previous image previews

            const img = document.createElement('img');
            img.src = e.target.result;
            img.style.maxWidth = '300px'; // Limit image width for display
            imagePreview.appendChild(img);
        };

        reader.readAsDataURL(file);
    }
}

async function autofillMetadata() {
    let url = document.getElementById('videoURL').value;
    console.log(url);

    try {
        var vidMetadata = await getYTVidMetadata(url);
        var imgBlob = await getYTThumbnail(vidMetadata.videoId);
    } catch (e) {
        console.log(e.message);
        return;
    }
    const imgUrl = URL.createObjectURL(imgBlob); 
    const file = new File([imgBlob], 'thumbnailImg', {type: 'image/jpeg'});

    const dataTransfer = new DataTransfer();
    dataTransfer.items.add(file)
    
    const imgInput = document.getElementById('imageInput');
    imgInput.files = dataTransfer.files;

    let vidTitleInput = document.getElementById('videoTitle');
    vidTitleInput.value = vidMetadata.title;
    
    let uploadDateInput = document.getElementById('uploadDate');
    uploadDateInput.value = vidMetadata.uploadDate;

    let creatorInput = document.getElementById('creatorInput');
    creatorInput.value = vidMetadata.channelTitle;

    console.log(`Setting preview to ${imgUrl}`) 
    let imgPreview = document.getElementById('imagePreview');
    imgPreview.innerHTML = ''; 

    const img = document.createElement('img');
    img.style.maxWidth = '300px'; 
    img.src = imgUrl;
    imgPreview.appendChild(img);
}

async function getYTVidMetadata(url) {
    const vidId = getVideoID(url);
    if (!vidId) {
        throw new Error('Vid ID not detected');
    }
    
    const response = await fetch(`/vid/metadata?vidId=${vidId}`);
    if (!response.ok) {
        throw new Error(`Response status: ${response.status}`);
    }
    
    const json = await response.json();
    console.log(json);
    return json;
}

async function getYTThumbnail(vidId) {
    const imgUrl = `/vid/thumbnail/?vidId=${vidId}`;
    console.log(`Getting thumbnail for ${imgUrl}`)

    const response = await fetch(imgUrl);
    if (!response.ok) {
        throw new Error('Unable to fetch thumbnail')
    }
    
    return response.blob()
}

function getVideoID(url) {
    let queryParamsString = url.split('?')[1]
    let queryParams = new URLSearchParams(queryParamsString)
    return queryParams.get('v')
}

document.getElementById("addPersonButton").addEventListener("click", () => {
    const personInputDivs = document.querySelectorAll("#personInputs input")
    const lastInputs = Array.from(personInputDivs).sort((a,b) => {
        a = a.getAttribute("name")
        b = b.getAttribute("name")
        return a-b
    })
    const template = document.getElementById("personInputTemplate").content;
    const newInput = document.importNode(template, true);

    const lastInput = lastInputs.pop()
    console.log(lastInput.getAttribute("name"))
    
    

    const lastInputName = lastInput.getAttribute("name")
    let lastInputNum = lastInputName.match(/\d+/)[0]

    newInput.name = `people[${Number(lastInputNum) + 1}]`

    const addPersonButton = document.querySelector("#personInputs button")
    lastInput.parentNode.insertBefore(newInput, addPersonButton)
})

function insertDropdownItem(e) {
    text = e.target.outerText
    dropDownList = e.target.parentNode
    // dropdown list is contained in div
    input = dropDownList.parentNode.previousElementSibling
    input.value = text

    dropDownList.remove()
}
