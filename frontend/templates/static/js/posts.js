function createPost() {
    // Placeholder for post creation functionality
    console.log("Create Post button clicked! Add modal or redirect logic here.");
}

function renderPostToContainer(post, container) {
    const parentDiv = document.createElement('div');
    parentDiv.className = 'post_comments';
    parentDiv.dataset.id = post.id;  // lowercase

    parentDiv.innerHTML = `
        <div class="post">
            <div class="post-header">
                <img src="./static/icon.png">
                <h3>${post.author}</h3>
            </div>
                    
            <div class="post-content">
                <h3>${post.title}</h3>
                <p>${post.content}</p>
            </div>

            <div class="post-actions">
                <button onclick="likePost(${post.id})">
                    <span>👍</span> Like (${post.reactions.likes.length})
                </button>
                <button onclick="dislikePost(${post.id})">
                    <span>👎</span> Dislike (${post.reactions.dislikes.length})
                </button>
                <button onclick="ShowComments(${post.id})">
                    <span>💬</span> Comments (${post.comments})
                </button>
                <button>
                    <span>#</span> ${post.category}
                </button>
            </div>
        </div>

        <!-- comments popup -->
        <div id="overlay" onclick="closePopup()"></div>
        <div id="popup">
            <div class="comments">
                <textarea id="cmntarea_${post.id}" placeholder="Write a comment..."></textarea>
                <button onclick="CreateComment(${post.id})">Submit</button>
            </div>
        </div>
    `;

    container.appendChild(parentDiv);
}

function loadPosts() {
    fetch('/api/posts')
        .then(response => response.json())
        .then(data => {
            if (!data.success) {
                console.error('Failed to load posts');
                return;
            }

            const parent = document.getElementById('posts'); // parent div in HTML
            const container = document.createElement('div'); // create container
            container.id = 'posts-container';

            data.posts.forEach(post => {
                renderPostToContainer(post, container); // fill container
            });

            parent.appendChild(container); // append to parent

        })
        .catch(err => console.error(err));
}