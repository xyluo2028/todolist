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
deploy/              Backend and Cassandra initialization container builds
frontend/            React application, npm manifest, and lockfile
docker-compose.yml   Frontend, backend, and Cassandra deployment
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

## Deploy with Docker Compose

From the repository root, build and start the whole application:

```bash
docker compose up --build -d --wait --wait-timeout 300
```

Open `http://localhost:3000`. Compose starts Cassandra, waits for it to become healthy, initializes the schema, starts the Go API, and then starts the React frontend served by Nginx. The first startup can take a few minutes.

The frontend sends requests to `/api` on its own origin, and Nginx forwards them to the Go service. Direct API access is available at `http://localhost:7071`; `GET /healthz` reports HTTP server health. No frontend environment file is needed for this deployment.

To change the host ports, copy the root `.env.example` to `.env` and edit `FRONTEND_PORT` and/or `SERVER_PORT`:

```bash
cp .env.example .env
```

Internal container ports stay fixed, so changing either host port does not require changing the frontend API URL. Go, Node.js, and npm do not need to be installed on the deployment host.

```bash
docker compose ps -a
docker compose logs -f app frontend
docker compose down
```

`docker compose down` retains the `cassandra_data` database volume. The schema initializer uses `CREATE ... IF NOT EXISTS`, adds the task `title` column when upgrading an older database, and preserves existing records. It is safe to run again. Cassandra, the API, and the frontend have health checks and restart policies. Re-run the startup command after code changes to rebuild and apply the updated containers.

## Run Go against local Cassandra

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

All API endpoints except registration and `/healthz` require HTTP Basic Authentication.

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
