// Make it so enter button clicks button
document.getElementById('searchInput')
    .addEventListener('keyup', e => {
        e.preventDefault()
        if (e.key === 'Enter') {
            document.getElementById("searchButton").click()
        }
    })

// API call to search videos on button click
document.getElementById('searchButton').addEventListener('click', function(event) {
    const query = document.getElementById('searchInput').value.trim();

    fetch('/search/videos?search=' + encodeURIComponent(query))
    .then(response => response.json())
    .then(data => {
        displaySearchResults(data);
    })
    .catch(error => {
        console.error('Error:', error);
    });
});

function displaySearchResults(results) {
    // Clear previous search results
    document.getElementById('searchResults').innerHTML = '';
    // console.log(results)
    
    // Display new search results
    results.forEach(result => {
        const resultElement = document.createElement('div');
        resultElement.className = 'videoResult'


        let videoThumbnail = document.createElement('img');
        videoThumbnail.src = result.thumbnail_name;
        videoThumbnail.className = 'videoThumbnail'

        videoThumbnail = wrapElementInLink(videoThumbnail, result.video_url)

        let videoTitle = document.createElement('h2');
        videoTitle.textContent = result.title;


        videoTitle = wrapElementInLink(videoTitle, result.video_url)
        
        const creatorName = document.createElement('h3');
        creatorName.textContent = result.creator_name;
        
        const creatorImage = document.createElement('img');
        creatorImage.src = result.creator_img_path;
        
        const uploadDate = document.createElement('p');
        uploadDate.textContent = result.creation_date.split('T')[0]

        
        const vidMetaData = document.createElement('div');
        vidMetaData.className = 'vidMetaData';
        
        // Data below the thumbnail that has the video's
        // title, channel name, and upload date which is next
        // to the creator's profile image
        vidMetaData.appendChild(videoTitle);
        vidMetaData.appendChild(creatorName);
        vidMetaData.appendChild(uploadDate);

        const videoDetails = document.createElement('div');
        videoDetails.className = 'videoDetails'
        videoDetails.append(creatorImage);
        videoDetails.append(vidMetaData);
        
        resultElement.appendChild(videoThumbnail);
        resultElement.append(videoDetails);

        document.getElementById('searchResults').appendChild(resultElement);
    });
}

function wrapElementInLink(element, link) {
    const linkElement = document.createElement('a');
    linkElement.href = link;
    linkElement.target = '_blank';
    
    linkElement.appendChild(element);

    return linkElement;
}