import { useState, useEffect, useMemo, useCallback } from 'react';
import ProjectList from '../components/ProjectList';
import { api } from '../api';

function Projects() {
  const auth = useMemo(() => JSON.parse(sessionStorage.getItem('auth')), []);
  const [projects, setProjects] = useState([]);
  const [newProjectName, setNewProjectName] = useState(''); // State for the new project name

  // useCallback for fetchProjects to call it on demand
  const fetchProjects = useCallback(async () => {
    if (!auth) return;
    try {
      const res = await api.getProjects(auth);
      const payload = res.data;
      let list = [];

      if (Array.isArray(payload)) {
        list = payload;
      } else if (typeof payload === 'string') {
        const trimmed = payload.trim();
        if (trimmed && trimmed.toLowerCase() !== 'no projects found!') {
          list = trimmed
            .split('\n')
            .map((s) => s.trim())
            .filter(Boolean)
            .filter((item) => item.toLowerCase() !== 'no projects found!');
        }
      }

      setProjects(list);
    } catch (err) {
      console.error('Failed loading projects:', err);
      setProjects([]);
    }
  }, [auth]);

  const handleRemoveProject = async (projectName) => {
    if (!auth || !auth.username || !auth.password) {
      alert('You must be signed in to remove a project.');
      return;
    }
    if (!window.confirm(`Are you sure you want to delete project "${projectName}"?`)) {
      return;
    }
    try {
      await api.removeProject(projectName, auth);
      fetchProjects(); // Refresh the list after deletion
    } catch (err) {
      console.error('Failed to remove project:', err);
      alert('Failed to remove project.');
    }
  };

  useEffect(() => {
    fetchProjects();
  }, [fetchProjects]); // Call fetchProjects when it (or its dependencies like auth) changes

  const handleCreateProject = async (e) => {
    e.preventDefault(); // Prevent default form submission
    if (!auth || !auth.username || !auth.password) {
      console.warn('🔒 Missing credentials, aborting createProject')
      alert('You must be signed in to create a project.')
      return
    }
    if (!newProjectName.trim()) {
      alert('Project name cannot be empty.');
      return;
    }
    try {
      await api.createProject(newProjectName.trim(), auth);
      setNewProjectName(''); // Clear the input field
      fetchProjects(); // Refresh the projects list
      // alert('Project created successfully!'); // Optional success feedback
    } catch (err) {
      console.error('Failed to create project:', err);
      const errorMsg = err.response?.data?.message || err.message || 'Failed to create project.';
      alert(`Error creating project: ${errorMsg}`); // User-friendly error
    }
  };

  return (
    <section className="section-card">
      <div className="section-header">
        <h2 className="section-title">Projects</h2>
        <p className="section-subtitle">
          Group work by initiative, team, or area of your life so every task has a home.
        </p>
      </div>

      <form onSubmit={handleCreateProject} className="project-form">
        <input
          type="text"
          value={newProjectName}
          onChange={(e) => setNewProjectName(e.target.value)}
          placeholder="Name your next project"
        />
        <button type="submit" className="btn btn--primary">
          Create project
        </button>
      </form>

      <div className="refresh-row">
        <span className="section-subtitle">
          {projects.length > 0
            ? `You are tracking ${projects.length} ${projects.length === 1 ? 'project' : 'projects'}.`
            : 'No projects yet — let’s make your first one!'}
        </span>
        <button type="button" className="btn btn--ghost" onClick={fetchProjects}>
          Refresh list
        </button>
      </div>

      {projects.length > 0 ? (
        <ProjectList projects={projects} onRemove={handleRemoveProject} />
      ) : (
        <div className="empty-state">
          Add a project to begin grouping related tasks and tracking deadlines together.
        </div>
      )}
    </section>
  );
}

export default Projects;
