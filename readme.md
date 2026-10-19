# Todo List

Todo List contains a Go HTTP API and the TaskFlow React frontend in one repository. Users can register, sign in, manage projects, create and edit tasks, mark tasks complete, and deactivate their accounts.

## Repository layout

```text
cmd/server/          Go server entry point
cmd/cassandra/       Cassandra schema initialization
internal/handlers/   HTTP endpoints
internal/services/   Application logic
internal/repository/ Cassandra and in-memory storage implementations
internal/middleware/ Basic authentication and CORS
internal/models/     User and task models
deploy/Dockerfile    Backend container build
frontend/            React application, npm manifest, and lockfile
docker-compose.yml   Backend and Cassandra services
```

The backend remains a Go module at the repository root. The frontend has its own npm dependencies and build inside `frontend/`.

## Prerequisites

- Go 1.22.2 or newer.
- Node.js 20 or newer and npm. The frontend's `.nvmrc` selects Node 22 when using nvm.
- Docker with Compose when using the containerized Cassandra backend.

On Windows, run the following commands inside WSL when the checkout is stored on the WSL filesystem.

## Local development

From the repository root, start the backend in one terminal:

```bash
go mod download
STORAGE_TYPE=inmem SERVER_PORT=7071 go run ./cmd/server
```

This uses temporary in-memory storage, which is cleared when the server stops. Stop the server with Ctrl+C.

Start React in a second terminal:

```bash
cd frontend
npm ci
npm start
```

Open `http://localhost:3000`. The browser calls the API at `http://localhost:7071` directly; the Go server handles CORS. Register through the Register tab, then sign in.

To use another API address, copy the frontend environment example and edit its value:

```bash
cp frontend/.env.example frontend/.env.local
```

Set `REACT_APP_API_URL` in that file and restart the React development server. Create React App also embeds this value during `npm run build`.

## Cassandra and Docker

The existing Compose setup builds the backend and runs Cassandra. For a new database, initialize its schema before starting the app:

```bash
docker compose up -d cassandra
docker compose run --rm cassandra-init
docker compose up --build -d app
```

The API is available at `http://localhost:7071`. Set `SERVER_PORT` in a root `.env` file to change the exposed API port, and set the frontend API URL to match.

```bash
docker compose logs -f app
docker compose down
```

`docker compose down` retains the database volume. Run the React development server separately as described above.

For a locally running Cassandra instance initialized with `cmd/cassandra/init.cql`, start Go directly with:

```bash
CASSANDRA_HOSTS=127.0.0.1:9042 CASSANDRA_KEYSPACE=todolist go run ./cmd/server
```

The backend defaults to `STORAGE_TYPE=cassandra`, `SERVER_PORT=7071`, and the Cassandra values shown above.

## Build and verification

From the repository root:

```bash
go test ./...
go vet ./...
go build -o bin/todolist-server ./cmd/server
npm --prefix frontend ci
npm --prefix frontend run build
CI=true npm --prefix frontend test -- --watchAll=false --passWithNoTests
```

The frontend production output is `frontend/build/`. Configure its static host to serve `index.html` for client routes such as `/projects`. The Go backend and React frontend retain separate build and deployment artifacts.

No automated test cases existed at the time of consolidation; the frontend test command allows that baseline. Before importing the frontend, `go test ./...` and `go vet ./...` reported an existing recursive formatting call in `internal/models/task.go`. `go test -vet=off ./...` and the backend build pass. In-memory project listing also has an existing issue: a new user's projects are not listed until their first task exists. These issues are separate from the repository move.

## Repository consolidation

The frontend was imported from `xyluo2028/todolist-fe-react` at commit `812ba755ebb0a31584d125e75c9244ceb45b3dd1` into `frontend/` using Git subtree without squashing. Its original commits and authors are retained in this repository's history. The backend starting commit was `7d2b1d059143663cbf50489029a7848a113ef539`.

Make future frontend changes in `frontend/` here. The original frontend repository can be deprecated after the consolidation is reviewed and merged.

## API authentication

All API endpoints except registration require HTTP Basic Authentication.

## Endpoints

### Register a New User
- **URL:** `/register`
- **Method:** POST
- **Headers:** `Content-Type: application/json`
- **Body Example:**
  ```json
  {
      "username": "test",
      "password": "test123"
  }
  ```
- **cURL Example:**
  ```bash
  curl -X POST \
    http://localhost:7071/register \
    -H 'Content-Type: application/json' \
    -d '{"username":"test","password":"test123"}'
  ```

### Welcome
- **URL:** `/welcome`
- **Method:** GET
- **Authentication:** Basic (e.g., `test:test123`)
- **cURL Example:**
  ```bash
  curl -X GET -u test:test123 http://localhost:7071/welcome
  ```

### Write (Add or Update) a Task
- **URL:** `/writeTask`
- **Method:** POST
- **Query Parameter:** `pjt` _(project name, required)_
- **Authentication:** Basic
- **Body:** JSON representation of a task (see [Task](./internal/models/task.go)). `title` is required (max 100 characters). `content` is optional (max 1000 characters). `id` is optional; if omitted, the server will generate a new one. The `updatedTime` field is set server-side. Project must exist, if not method return error
- **Body Example:**
  ```json
  {
      "title": "Grocery run",
      "content": "Buy milk, eggs, and bread",
      "priority": 2,
      "due": "2025-05-09T15:04:05Z",
      "completed": false
  }
  ```
- **cURL Example:**
```bash
curl -X POST -u test:test123 "http://localhost:7071/writeTask?pjt=home" \
  -H 'Content-Type: application/json' \
  -d '{"title":"Grocery run","content":"Buy milk, eggs, and bread","priority":2,"due":"2025-05-09T15:04:05Z","completed":false}'
```

### Get All Tasks for a Project
- **URL:** `/printTasks`
- **Method:** GET
- **Query Parameter:** `pjt` _(project name, required)_
- **Authentication:** Basic
- **cURL Example:**
  ```bash
  curl -X GET -u test:test123 "http://localhost:7071/printTasks?pjt=home"
  ```

### Get All Projects
- **URL:** `/printProjects`
- **Method:** GET
- **Authentication:** Basic
- **cURL Example:**
  ```bash
  curl -X GET -u test:test123 http://localhost:7071/printProjects
  ```

### Complete a Task
- **URL:** `/completeTask`
- **Method:** GET
- **Query Parameters:** 
  - `pjt` _(project name, required)_
  - `key` _(task ID, required)_
- **Authentication:** Basic
- **cURL Example:**
  ```bash
  curl -X GET -u test:test123 "http://localhost:7071/completeTask?pjt=home&key=task_xxx"
  ```

### Remove a Task
- **URL:** `/removeTask`
- **Method:** DELETE
- **Query Parameters:** 
  - `pjt` _(project name, required)_
  - `key` _(task ID, required)_
- **Authentication:** Basic
- **cURL Example:**
  ```bash
  curl -X DELETE -u test:test123 "http://localhost:7071/removeTask?pjt=home&key=task_xxx"
  ```

### Remove a Project
- **URL:** `/removeProject`
- **Method:** DELETE
- **Query Parameter:** `pjt` _(project name, required)_
- **Authentication:** Basic
- **cURL Example:**
  ```bash
  curl -X DELETE -u test:test123 "http://localhost:7071/removeProject?pjt=home"
  ```

### Create a Project
- **URL:** `/createProject`
- **Method:** POST
- **Query Parameter:** `pjt` _(project name, required)_
- **Authentication:** Basic
- **cURL Example:**
  ```bash
  curl -X POST -u test:test123 "http://localhost:7071/createProject?pjt=home"
  ```

### Deactivate (Delete) a User
- **URL:** `/deactivate`
- **Method:** DELETE
- **Authentication:** Basic (the user to delete)
- **cURL Example:**
  ```bash
  curl -X DELETE -u test:test123 http://localhost:7071/deactivate
  ```
