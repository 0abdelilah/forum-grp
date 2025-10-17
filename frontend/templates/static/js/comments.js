function ShowComments(postId) {
    document.getElementById('popup').style.display = 'block';
    document.getElementById('overlay').style.display = 'block';

    const url = '/api/comments?postId=' + postId
    fetch(url)
        .then(response => response.json())
        .then(data => console.log(data))
        .catch(error => console.error('Fetch error:', error))
}

function CreateComment(btn) {
    const postId = btn.closest('.post_comments').dataset.id;
    const content = btn.closest('.post_comments').querySelector('#cmntarea').value;

    fetch('/api/comment', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            postid: postId,
            content: content
        })
    })
        .then(response => response.json())
        .then(data => console.log(data))
        .catch(err => console.error(err));
}

function closePopup() {
    document.getElementById('popup').style.display = 'none';
    document.getElementById('overlay').style.display = 'none';
}
