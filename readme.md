# Todo List API

This API implements a simple Todo List using HTTP endpoints. All endpoints (except registration) require Basic Authentication.

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
- **Body:** JSON representation of a task (see [Task](./internal/models/task.go)); `id` is optional. If omitted, the server will generate a new one. The `updatedTime` field is set server-side.
- **Body Example:**
  ```json
  {
      "content": "Buy groceries",
      "priority": 2,
      "due": "2025-05-09T15:04:05Z",
      "completed": false
  }
  ```
- **cURL Example:**
  ```bash
  curl -X POST -u test:test123 "http://localhost:7071/writeTask?pjt=home" \
    -H 'Content-Type: application/json' \
    -d '{"content":"Buy groceries","priority":2,"due":"2025-05-09T15:04:05Z","completed":false}'
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
- **Method:** GET
- **Query Parameters:** 
  - `pjt` _(project name, required)_
  - `key` _(task ID, required)_
- **Authentication:** Basic
- **cURL Example:**
  ```bash
  curl -X GET -u test:test123 "http://localhost:7071/removeTask?pjt=home&key=task_xxx"
  ```

### Remove a Project
- **URL:** `/removeProject`
- **Method:** GET
- **Query Parameter:** `pjt` _(project name, required)_
- **Authentication:** Basic
- **cURL Example:**
  ```bash
  curl -X GET -u test:test123 "http://localhost:7071/removeProject?pjt=home"
  ```

### Deactivate (Delete) a User
- **URL:** `/deactivate`
- **Method:** DELETE
- **Authentication:** Basic (the user to delete)
- **cURL Example:**
  ```bash
  curl -X DELETE -u test:test123 http://localhost:7071/deactivate
  ```

## Running the Server

Start the server using `go run`:

```bash
nohup go run cmd/server/main.go > logs/todolist.log 2>&1 &
```

To run the server without a Cassandra database, using an in-memory store, you can set the `STORAGE_TYPE` environment variable to `inmem`. You can also specify a different port using `SERVER_PORT`. For example:

```bash
STORAGE_TYPE=inmem SERVER_PORT={your_port} nohup go run cmd/server/main.go > logs/todolist.log 2>&1 &
```

To stop the server completely, terminate the process:

```bash
pkill -TERM -f "go run cmd/server/main.go"
```

Replace `<PID>` with the actual process ID of your server.

### Using Docker

Alternatively, you can run the server using Docker Compose, which will also set up and initialize the Cassandra database. Ensure you have a `docker-compose.yml` file in the project root (as provided in the context).

1.  **Build and run the services (app and Cassandra):**
    ```bash
    docker-compose up --build -d
    ```
    This command will build the `todolist-api` image (if it doesn't exist or if changes are detected) and start all services defined in `docker-compose.yml` in detached mode. The `-d` flag runs containers in the background. The `--build` flag forces a rebuild of the image.

    The application will be available at `http://localhost:7071`.

2.  **To view logs:**
    ```bash
    docker-compose logs -f app
    ```
    Replace `app` with `cassandra` or `cassandra-init` to view logs for those services.

3.  **To stop and remove the containers, networks, and volumes created by `up`:**
    ```bash
    docker-compose down
    ```
    If you want to remove volumes defined in the `volumes` section of `docker-compose.yml` (like `cassandra_data`), use:
    ```bash
    docker-compose down -v
    ```

---
