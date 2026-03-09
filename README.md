# Literary Lions Forum (aka Lion'Library)

Welcome to **Literary Lions aka Lion's Library**, a digital haven for book lovers to roar out loud with their insights, howl discussions, and growl the passion for reading online. This project modernizes traditional book forums with a fully functional web forum using **Go**, **HTML/CSS**, and **SQLite** — all wrapped up in a cozy, dockerized package.

---
## Project Team

Tan Hoang 🤖  
Anh Le 🐳
Nooa Niinikangas 🐔  

## Getting Started
To run the project locally or in a Docker container, follow these steps:

### Live demo
If you dont want to run the project locally, you can visit the live demo at: https://lions.tancodes.com/

Currious about how can I publish the project? Check out my blog post at: [link](https://dev.to/tan1193/from-localhost3000-to-the-world-deploying-your-dockerized-website-with-cloudflared-traefik-1g12)

### Clone the repo
```
git clone https://gitea.kood.tech/anhle/literary-lions.git
cd literaryLions
```
- if run locally:

```
go run main.go
```

- or if using Docker:
```
docker build -t literary-lions .
docker run -p 8080:8080 literary-lions
```

Then visit http://localhost:8080

### Accessing the Forum
You can register an account and log in to explore the forum, create posts, comment on discussions, and interact with other book enthusiasts. If you want to test the project without registering, you can use the following credentials:

- **Email:** test@example.com
- **Password:** password

## Project Description

Literary Lions Forum is an online discussion platform where users can:

- **Create** posts and **comment** on posts
- **Search** for posts that spark your interests
- **Categorize** discussions by book genres, themes, and more
- **Like** posts and comments
- **Register and log in** with secure authentication
- **Handling data** using **SQLite**

---

## Key Features

### User Authentication
- Unique **username**, **email**, and **password** required
- Login via email and password only
- Passwords securely **encrypted**
- Sessions managed using **UUID-based cookies**

### Posts & Comments
- Create posts with a title, content, category and optional associated **hashtags**
- Comment on any post
- Posts can be **filtered** by category
  
### Book & Category Organization
- Posts are categorized into different discussion topics
- Categories also allow thematic discussions across books

### Likes 
- Users can like **posts** and **comments**
- Like counts are visible to all users
- Restrictions are in place to prevent multiple likes per user

### Dockerized Deployment
- Built-in `Dockerfile` to simplify deployment
- Easily run the entire app in a container

---

## Tech Stack

| Layer            | Technology              |
|---------------   |--------------------     |
| Backend          | Go (Golang)             |
| Frontend         | HTML + CSS              |
| Database         | SQLite (via go-sqlite3) |
| Containerization | Docker                  |

---

## Project Structure

/literaryLions  
├── /database # SQLite database files and schema  
├── /handlers # backend assets for web server  
├── /images # images for web server  
├── /models # structs for backends' services reflecting SQL schema  
├── /services # backend assets for web functionality  
├── /static # CSS and frontend assets  
├── /templates # HTML templates  
├── /utils # helper functions  
├── .dockerignore # Docker config  
├── database.db # Testing database from development phase  
├── Dockerfile # Docker config  
├── main.go # App entry point  
├── go.mod  
├── go.sum  
├── lions.drawio # SQL schema  
├── main.go # App entry point  
└── README.md    

## How to Run the Project

### Prerequisites

- Go (v1.21+)
- Docker (if using Docker)
- GCC (required for `go-sqlite3`)
  - On Linux: `sudo apt install build-essential`
  - On macOS: Xcode Command Line Tools
  - On Windows: Use MinGW or WSL
- For the lions.drawio file, it is also recommended to install draw.io extension for VS Code via this link for accessibility: https://marketplace.visualstudio.com/items?itemName=hediet.vscode-drawio



## Database schema overview
Refer to file lions.drawio for details (again, for the lions.drawio file, it is also recommended to install draw.io extension for VS Code via this link for accessibility: https://marketplace.visualstudio.com/items?itemName=hediet.vscode-drawio)

## Bonus feature 
- Session UUIDs
- Password hashing
- Docker cleanup and metadata
- Publicly accessible demo site
- Hashtags 

## Possible future improvements
- Admin/mod tools
- Improve user profile page where user can upload pictures and customize more
- Notifications 
