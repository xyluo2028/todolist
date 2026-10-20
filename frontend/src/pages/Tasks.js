import { useState, useEffect, useCallback, useMemo } from 'react';
import { useParams, Link } from 'react-router-dom';
import TaskList from '../components/TaskList';
import { api } from '../api';

function Tasks() {
  const { projectName = '' } = useParams();
  const auth = useMemo(() => JSON.parse(sessionStorage.getItem('auth')), []);
  const [tasks, setTasks] = useState([]);
  const [loadError, setLoadError] = useState('');

  const loadTasks = useCallback(async () => {
    try {
      if (!projectName) {
        setTasks([]);
        setLoadError('');
        return;
      }

      const res = await api.getTasks(projectName, auth);

      let data = res.data;

      let list = [];

      if (Array.isArray(data)) {
        list = data;
      } else if (typeof data === 'string') {
        const trimmed = data.trim();
        if (trimmed && trimmed.toLowerCase() !== 'no tasks found!') {
          try {
            const parsed = JSON.parse(trimmed);
            if (Array.isArray(parsed)) {
              list = parsed;
            }
          } catch (jsonErr) {
            console.warn('Unexpected string payload for tasks; falling back to empty list', jsonErr);
          }
        }
      }

      setTasks(list);
      setLoadError('');
    } catch (error) {
      console.error('Failed to load tasks:', error);
      setLoadError('Unable to refresh tasks. Please try again.');
    }
  }, [projectName, auth]);

  const refreshTasks = async (newTask, key, action) => {
    if (!projectName) {
      throw new Error('Select a project before changing tasks.');
    }

    if (action === 'complete') {
      await api.completeTask(projectName, key, auth);
    } else if (action === 'remove') {
      await api.removeTask(projectName, key, auth);
    } else if (newTask) {
      await api.writeTask(projectName, newTask, auth);
    }

    await loadTasks();
  };

  useEffect(() => {
    loadTasks();
  }, [loadTasks]);

  const subtitle = projectName
    ? 'Add, schedule, and complete tasks to keep this project on track.'
    : 'Pick a project to dive into its backlog and start moving items forward.';

  return (
    <section className="section-card">
      <div className="section-header">
        <h2 className="section-title">
          {projectName ? `Project · ${projectName}` : 'Tasks overview'}
        </h2>
        <p className="section-subtitle">{subtitle}</p>
      </div>

      {loadError && (
        <div className="form-actions" role="alert">
          <p className="section-subtitle">{loadError}</p>
          <button type="button" className="btn btn--secondary" onClick={loadTasks}>
            Retry
          </button>
        </div>
      )}

      {projectName ? (
        <TaskList tasks={tasks} project={projectName} refresh={refreshTasks} />
      ) : (
        <div className="empty-state">
          <p>Choose a project to see its tasks and start moving work forward.</p>
          <Link to="/projects" className="btn btn--primary">
            Browse projects
          </Link>
        </div>
      )}
    </section>
  );
}

export default Tasks;
