import { useState } from 'react';
import TaskItem from './TaskItem';

function TaskList({ tasks, project, refresh }) {
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [due, setDue] = useState('');
  const [priority, setPriority] = useState(1);

  const addTask = async () => {
    const trimmedTitle = title.trim();
    const trimmedContent = content.trim();

    if (!trimmedTitle) {
      alert('Please provide a task title.');
      return;
    }
    if (trimmedTitle.length > 100) {
      alert('Task titles must be 100 characters or fewer.');
      return;
    }
    if (!due) {
      alert('Please select a due date.');
      return;
    }
    if (trimmedContent.length > 1000) {
      alert('Task details must be 1000 characters or fewer.');
      return;
    }

    const payload = {
      title: trimmedTitle,
      due: new Date(due).toISOString(),
      priority: parseInt(priority, 10),
      completed: false,
    };

    if (trimmedContent) {
      payload.content = trimmedContent;
    }

    await refresh(payload);
    setTitle('');
    setContent('');
    setDue('');
    setPriority(1);
  };

  const onComplete = async (id) => {
    await refresh(null, id, 'complete');
  };

  const onRemove = async (id) => {
    await refresh(null, id, 'remove');
  };

  const onUpdate = async (updatedTask) => {
    await refresh(updatedTask, updatedTask.id, 'update');
  };

  return (
    <div className="task-area">
      <div className="task-header">
        <form
          className="task-inputs"
          onSubmit={(event) => {
            event.preventDefault();
            addTask();
          }}
        >
          <input
            className="form-input"
            value={title}
            maxLength={100}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Task title"
          />
          <textarea
            className="form-input form-input--textarea"
            value={content}
            maxLength={1000}
            rows={3}
            onChange={(e) => setContent(e.target.value)}
            placeholder={`Task details for ${project} (optional)`}
          />
          <input
            className="form-input"
            type="datetime-local"
            value={due}
            onChange={(e) => setDue(e.target.value)}
          />
          <select
            className="form-select"
            value={priority}
            onChange={(e) => setPriority(e.target.value)}
          >
            {[1, 2, 3, 4, 5].map((num) => (
              <option key={num} value={num}>
                Priority {num}
              </option>
            ))}
          </select>
          <button type="submit" className="btn btn--primary">
            Add task
          </button>
        </form>
      </div>

      {tasks.length > 0 ? (
        <div className="task-list">
          {tasks.map((task) => (
            <TaskItem
              key={task.id}
              task={task}
              onComplete={onComplete}
              onRemove={onRemove}
              onUpdate={onUpdate}
            />
          ))}
        </div>
      ) : (
        <div className="empty-state">
          No tasks yet — capture your first idea above to get the ball rolling.
        </div>
      )}
    </div>
  );
}

export default TaskList;
