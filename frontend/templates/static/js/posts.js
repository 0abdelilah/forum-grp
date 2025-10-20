function openPopup() {
    const overlay = document.getElementById("overlay");
    const popup = document.getElementById("popup");
    overlay.style.display = "flex";
    popup.style.display = "block";
}

function closePopup() {
    const overlay = document.getElementById("overlay");
    const popup = document.getElementById("popup");
    overlay.style.display = "none";
    popup.style.display = "none";
}

function createPost() {
    fetch('/api/create_post', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ title: 'title', content: 'content' })
    })
    .then(res => res.json())
    .then(data => {
        if (data.success) {
            console.log("Success:", data);
        } else {
            console.log("Error:", data);
        }
    })
    .catch(err => console.error(err));
}
