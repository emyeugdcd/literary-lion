<div align="center">
  <h1> Literary Lions (Lion's Library)</h1>
  <p>
    <strong>A modernized, fully-featured discussion forum connecting book enthusiasts globally.</strong>
  </p>
  <p>
    <a href="https://lions.tancodes.com/"><strong>View Live Demo »</strong></a>
    ·
    <a href="https://dev.to/tan1193/from-localhost3000-to-the-world-deploying-your-dockerized-website-with-cloudflared-traefik-1g12">Read Deployment Blog Post</a>
  </p>
  <p>
    <img alt="Go" src="https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white" />
    <img alt="SQLite" src="https://img.shields.io/badge/sqlite-%2307405e.svg?style=for-the-badge&logo=sqlite&logoColor=white" />
    <img alt="Docker" src="https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white" />
  </p>
</div>

<br />

## Overview

**Literary Lions** is a digital haven built from the ground up to modernize traditional book forums. It connects readers through a performant, custom-built web application completely devoid of bulky frontend frameworks. Relying strictly on raw Go, HTML, CSS, and SQLite, the application stands as a testament to efficient, foundational back-end and web development.

### Key Features

- **Forum Mechanics:** Create deep-dive posts with customizable titles, content, categories, and #hashtags. 
- **Nested Discussions:** Seamless comment system mapped contextually to individual posts.
- **Engagement Tools:** Global like system with backend rate limiting and duplication prevention.
- **Custom Authentication:** From-scratch session handling natively tracking UUID-based secure cookies alongside salted password hashes.
- **Filtering & Search:** Easily parse through categories, genres, and themes to discover relevant discourse.

## Built With

| Layer | Technology |
| :--- | :--- |
| **Backend API & Server** | Go (Golang) |
| **Frontend Rendering** | HTML5, CSS3, Go Templates |
| **Persistence** | SQLite (`go-sqlite3`) |
| **DevOps & Hosting** | Docker |

## Architecture highlights

This monolithic service implements proper separation of concerns (MVC inspired):
- `/handlers`: Pure HTTP controllers managing requests and formulating protocol responses.
- `/services`: The business logic layer governing application rules.
- `/models`: Strong typing dictating exactly how internal services reflect the SQLite schema.
- `/static & /templates`: The visual rendering engine.

*An internally accessible `lions.drawio` visualizes our relational SQLite footprint comprehensively.*

## Getting Started

### Credentials for demonstration
- **Email:** `test@example.com`
- **Password:** `password`

### Prerequisites
- Go (v1.21+) or Docker
- GCC (required natively for `go-sqlite3` bindings)

### Installation

**Using Docker (Recommended):**
```bash
git clone https://gitea.kood.tech/anhle/literary-lions.git
cd literaryLions

docker build -t literary-lions .
docker run -p 8080:8080 literary-lions
```
*Visit `http://localhost:8080` to view the service.*

**Running Locally:**
```bash
go run main.go
```

## The Team
- **Anh Le**  - Backend, Architecture, & Containerization
- **Tan Hoang**  - Deployment & Infrastructure
- **Nooa Niinikangas**  - Frontend & Interactivity
