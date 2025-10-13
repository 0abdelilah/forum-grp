RUN: - cd backend
     - go run .

Person 1: Handles user authentication (login, registration).

Person 2: Manages post creation and display.

Person 3: Implements commenting on posts.

Person 4: Adds liking functionality and post filtering.

Person 5: Sets up the database, writes tests, and manages Docker/deployment.


/api/posts:

[
  {
    "id": 1,
    "title": "Getting Started with Bug Bounty",
    "author": "Gadr",
    "created_at": "2025-10-13T18:00:00Z",
    "content": "Learn how to find your first vulnerability and report it responsibly.",
    "comments": [
      {
        "author": "Alice",
        "content": "Great guide for beginners!",
        "created_at": "2025-10-13T19:10:00Z"
      },
      {
        "author": "Bob",
        "content": "Would love a follow-up on responsible disclosure.",
        "created_at": "2025-10-13T19:45:00Z"
      }
    ]
  },
  {
    "id": 2,
    "title": "Top 5 Tools for Web Pentesting",
    "author": "Gadr",
    "created_at": "2025-10-12T16:30:00Z",
    "content": "A quick overview of the most effective tools for modern web app pentesting.",
    "comments": [
      {
        "author": "Charlie",
        "content": "Burp Suite and Nuclei are must-haves!",
        "created_at": "2025-10-12T17:05:00Z"
      },
      {
        "author": "Dana",
        "content": "You should include some open-source tools too.",
        "created_at": "2025-10-12T18:20:00Z"
      }
    ]
  }
]