# Todolist React Frontend

This project provides the React client used to interact with the [Todo List backend](https://github.com/xyluo2028/todolist). It offers registration, authentication, project management, and task tracking features backed by the Go API.

## Prerequisites
- Node.js 18+
- npm 8+
- A running instance of the backend API (defaults to `http://localhost:7071`).

## Getting Started
1. **Install dependencies**
   ```bash
   npm install
   ```
2. **Run the development server**
   ```bash
   npm start
   ```
   The app is served at `http://localhost:3000` and proxies requests directly to the backend URL configured in `src/api/index.js`.

3. **Build for production**
   ```bash
   npm run build
   ```

## Configuring the API URL
The frontend currently targets `http://localhost:7071`. If your backend runs elsewhere, update the `baseURL` in `src/api/index.js` to match your environment.

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

## Related Repositories
- Backend service: https://github.com/xyluo2028/todolist

## Contributing
Pull requests are welcome. Please run the existing lint/build scripts before submitting changes.
