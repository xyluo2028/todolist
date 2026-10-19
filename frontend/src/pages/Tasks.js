import { useState, useEffect, useCallback, useMemo } from 'react';
import { useParams, Link } from 'react-router-dom';
import TaskList from '../components/TaskList';
import { api } from '../api';

function Tasks() {
  const { projectName } = useParams();
  const decodedProjectName = useMemo(
    () => (projectName ? decodeURIComponent(projectName) : ''),
    [projectName]
  );
  const auth = useMemo(() => JSON.parse(sessionStorage.getItem('auth')), []);
  const [tasks, setTasks] = useState([]);

  const loadTasks = useCallback(async (newTask, key, action) => {
    try {
      if (!decodedProjectName) {
        if (newTask || action) {
          console.warn("Cannot perform task operations without a project selected.");
        }
        setTasks([]);
        return;
      }

      if (action === 'complete') {
        await api.completeTask(decodedProjectName, key, auth);
      } else if (action === 'remove') {
        await api.removeTask(decodedProjectName, key, auth);
      } else if (newTask) {
        await api.writeTask(decodedProjectName, newTask, auth);
      }
      const res = await api.getTasks(decodedProjectName, auth);

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
    } catch (error) {
      console.error('Failed to process tasks:', error);
      setTasks([]);
    }
  }, [decodedProjectName, auth]);

  useEffect(() => {
    loadTasks();
  }, [loadTasks]);

  const subtitle = decodedProjectName
    ? 'Add, schedule, and complete tasks to keep this project on track.'
    : 'Pick a project to dive into its backlog and start moving items forward.';

  return (
    <section className="section-card">
      <div className="section-header">
        <h2 className="section-title">
          {decodedProjectName ? `Project · ${decodedProjectName}` : 'Tasks overview'}
        </h2>
        <p className="section-subtitle">{subtitle}</p>
      </div>

      {decodedProjectName ? (
        <TaskList tasks={tasks} project={decodedProjectName} refresh={loadTasks} />
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
