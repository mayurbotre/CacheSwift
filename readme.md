# CacheSwift

## Description
LRU (Least Recently Used) Cache is a type of caching mechanism that maintains a limited size cache, evicting the least recently used items when the cache reaches its maximum capacity. It operates on the principle that items accessed recently are more likely to be accessed again in the near future. LRU Cache implementations typically involve efficient data structures like doubly linked lists and hash maps to ensure quick access, insertion, and deletion operations, making it suitable for scenarios where quick access to frequently used data is crucial for optimizing performance.

## Technologies Used
- Go (Golang)
- React.js
- Other technologies used (e.g., Axios, Gin, etc.)

## Directory Structure

project-root/
│
├── client/
├── server/
├── README.md

markdown


## Prerequisites
- Go installed
- Node.js and npm installed
- IDE (e.g., VS Code) or preferred text editor

## Getting Started

### Backend (Go / Golang)

1. **Clone the repository**
   ```bash
   git clone https://github.com/mayurbotre/CacheSwift.git
   cd CacheSwift/server

    Install dependencies

    bash

go mod tidy

Run the backend server

bash

    go run main.go

    This will start the backend server at http://localhost:8080.

    API Endpoints
        GET /cache/:key: Retrieve a cached item by key.
        GET /cache: Retrieve all cached items.
        POST /cache: Add a new item to the cache.
        DELETE /cache/:key: Delete a cached item by key.

Frontend (React.js)

    Navigate to frontend directory

    bash

cd CacheSwift/client

Install dependencies

bash

npm install

Start the development server

bash

    npm start

    This will start the React development server at http://localhost:3000.

    Accessing the Application
    Open your web browser and go to http://localhost:3000 to access the frontend.
