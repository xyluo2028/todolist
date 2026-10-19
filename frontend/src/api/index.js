import axios from 'axios';

const apiClient = axios.create({
  baseURL: process.env.REACT_APP_API_URL || 'http://localhost:7071',
});

const authConfig = (auth) => ({
  headers: {
    Authorization: `Basic ${btoa(`${auth.username}:${auth.password}`)}`,
  },
});

export const api = {
  register: (username, password) => apiClient.post('/register', { username, password }),
  getProjects: (auth) => apiClient.get('/printProjects', authConfig(auth)),
  getTasks: (project, auth) =>
    apiClient.get(`/printTasks?pjt=${encodeURIComponent(project)}`, authConfig(auth)),
  writeTask: (project, task, auth) =>
    apiClient.post(`/writeTask?pjt=${encodeURIComponent(project)}`, task, authConfig(auth)),
  deactivateAccount: (auth) => apiClient.delete('/deactivate', authConfig(auth)),
  completeTask: (project, key, auth) =>
    apiClient.get(
      `/completeTask?pjt=${encodeURIComponent(project)}&key=${encodeURIComponent(key)}`,
      authConfig(auth)
    ),
  removeTask: (project, key, auth) =>
    apiClient.delete(
      `/removeTask?pjt=${encodeURIComponent(project)}&key=${encodeURIComponent(key)}`,
      authConfig(auth)
    ),
  createProject: (project, auth) =>
    apiClient.post(`/createProject?pjt=${encodeURIComponent(project)}`, {}, authConfig(auth)),
  removeProject: (project, auth) =>
    apiClient.delete(`/removeProject?pjt=${encodeURIComponent(project)}`, authConfig(auth)),
};
