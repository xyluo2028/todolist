# Todolist React Frontend

This directory contains the React client for the Go backend in the same repository. It offers registration, authentication, project management, and task tracking. See the [root README](../readme.md) for backend setup, Cassandra, and repository-wide verification.

## Prerequisites
- Node.js 20+ and npm (`.nvmrc` selects Node 22 for nvm users).
- A running instance of the backend API (defaults to `http://localhost:7071`).

## Getting Started
Run these commands from `frontend/`.

1. **Install the locked dependencies**
   ```bash
   npm ci
   ```
2. **Run the development server**
   ```bash
   npm start
   ```
   The app is served at `http://localhost:3000`. The browser calls the backend directly; the Go API handles CORS.

3. **Build for production**
   ```bash
   npm run build
   ```

## Configuring the API URL
The frontend defaults to `http://localhost:7071`. To use another API address:

```bash
cp .env.example .env.local
```

Edit `REACT_APP_API_URL` in `.env.local`, then restart `npm start`. This setting is embedded at build time for production, so set it before running `npm run build`.

Production files are written to `build/`. Configure your static host to serve `index.html` for client routes such as `/projects`.

## Authentication Flow
- Register a user via the **Register** tab under `/auth`.
- After successful registration you are redirected to the **Login** tab.
- Logging in stores credentials for the current browser session and unlocks protected routes such as `/projects` and project task views.

## Project Structure
```
src/
├── api/           # Axios client wrapper for backend API calls
├── components/    # Reusable UI components
├── pages/         # Route-level components
└── index.js       # App bootstrap
```

## Verification

```bash
npm run build
CI=true npm test -- --watchAll=false --passWithNoTests
```

There are currently no automated frontend test cases.

## Repository history

This client was imported from `xyluo2028/todolist-fe-react` with its complete Git history. Make future frontend changes in this directory.

## Contributing
Pull requests belong in `xyluo2028/todolist`. Run the build and verification commands before submitting changes.
