import { useEffect, useMemo, useState } from 'react';

const toDateTimeLocalValue = (isoString) => {
  if (!isoString) {
    return '';
  }
  const parsed = new Date(isoString);
  if (Number.isNaN(parsed.getTime())) {
    return '';
  }
  const tzOffsetMs = parsed.getTimezoneOffset() * 60 * 1000;
  const local = new Date(parsed.getTime() - tzOffsetMs);
  return local.toISOString().slice(0, 16);
};

function TaskItem({ task, onComplete, onRemove, onUpdate }) {
  const isCompleted = Boolean(task.completed);
  const dueDate = useMemo(() => (task.due ? new Date(task.due) : null), [task.due]);
  const dueLabel = dueDate && !Number.isNaN(dueDate.getTime())
    ? dueDate.toLocaleString()
    : 'No due date';
  const rawTitle = typeof task.title === 'string' ? task.title.trim() : '';
  const rawContent = typeof task.content === 'string' ? task.content.trim() : '';
  const displayTitle = rawTitle || rawContent || 'Untitled task';
  const showDescription = Boolean(rawContent && rawContent !== displayTitle);
  const [isEditing, setIsEditing] = useState(false);
  const [editTitle, setEditTitle] = useState(task.title ?? '');
  const [editContent, setEditContent] = useState(task.content ?? '');
  const [editDue, setEditDue] = useState(toDateTimeLocalValue(task.due));
  const [editPriority, setEditPriority] = useState(
    typeof task.priority === 'number' ? String(task.priority) : '1'
  );

  const resetForm = () => {
    setEditTitle(task.title ?? '');
    setEditContent(task.content ?? '');
    setEditDue(toDateTimeLocalValue(task.due));
    setEditPriority(typeof task.priority === 'number' ? String(task.priority) : '1');
  };

  useEffect(() => {
    if (!isEditing) {
      setEditTitle(task.title ?? '');
      setEditContent(task.content ?? '');
      setEditDue(toDateTimeLocalValue(task.due));
      setEditPriority(typeof task.priority === 'number' ? String(task.priority) : '1');
    }
  }, [task, isEditing]);

  const startEditing = () => {
    resetForm();
    setIsEditing(true);
  };

  const handleCancel = () => {
    resetForm();
    setIsEditing(false);
  };

  const handleSave = async (event) => {
    event.preventDefault();
    const trimmedTitle = editTitle.trim();
    const trimmedContent = editContent.trim();

    if (!trimmedTitle) {
      alert('Please provide a task title.');
      return;
    }
    if (trimmedTitle.length > 100) {
      alert('Task titles must be 100 characters or fewer.');
      return;
    }
    if (!editDue) {
      alert('Please select a due date.');
      return;
    }
    if (trimmedContent.length > 1000) {
      alert('Task details must be 1000 characters or fewer.');
      return;
    }

    const payload = {
      id: task.id,
      title: trimmedTitle,
      content: trimmedContent,
      due: new Date(editDue).toISOString(),
      priority: parseInt(editPriority, 10),
      completed: Boolean(task.completed),
    };

    try {
      await onUpdate(payload);
      setIsEditing(false);
    } catch (error) {
      console.error('Failed to update task:', error);
      alert('Unable to update this task right now. Please try again later.');
    }
  };

  const priorityOptions = useMemo(() => {
    const base = [1, 2, 3, 4, 5];
    if (typeof task.priority === 'number' && !base.includes(task.priority)) {
      base.push(task.priority);
      base.sort((a, b) => a - b);
    }
    return base;
  }, [task.priority]);

  return (
    <div
      className={`task-card${isCompleted ? ' task-card--completed' : ''}${
        isEditing ? ' task-card--editing' : ''
      }`}
    >
      <div className="task-card__info">
        {isEditing ? (
          <form className="task-card__edit-fields" onSubmit={handleSave}>
            <input
              className="form-input"
              value={editTitle}
              maxLength={100}
              onChange={(e) => setEditTitle(e.target.value)}
              placeholder="Task title"
            />
            <textarea
              className="form-input form-input--textarea"
              value={editContent}
              maxLength={1000}
              rows={3}
              onChange={(e) => setEditContent(e.target.value)}
              placeholder="Task details (optional)"
            />
            <div className="task-card__edit-row">
              <input
                className="form-input"
                type="datetime-local"
                value={editDue}
                onChange={(e) => setEditDue(e.target.value)}
              />
              <select
                className="form-select"
                value={editPriority}
                onChange={(e) => setEditPriority(e.target.value)}
              >
                {priorityOptions.map((num) => (
                  <option key={num} value={String(num)}>
                    Priority {num}
                  </option>
                ))}
              </select>
            </div>
            <div className="task-card__actions task-card__actions--editing">
              <button type="submit" className="btn btn--primary">
                Save
              </button>
              <button type="button" className="btn btn--ghost" onClick={handleCancel}>
                Cancel
              </button>
            </div>
          </form>
        ) : (
          <>
            <p className={`task-card__title${isCompleted ? ' task-card__title--done' : ''}`}>
              {displayTitle}
            </p>
            {showDescription && (
              <p className="task-card__description">
                {rawContent}
              </p>
            )}
            <div className="task-card__meta">
              <span className="badge badge--due">Due · {dueLabel}</span>
              {typeof task.priority === 'number' && (
                <span className="badge badge--priority">Priority {task.priority}</span>
              )}
              {isCompleted && <span className="badge badge--done">Completed</span>}
            </div>
          </>
        )}
      </div>
      {!isEditing && (
        <div className="task-card__actions">
          {!isCompleted && (
            <button
              type="button"
              className="btn btn--secondary"
              onClick={() => onComplete(task.id)}
            >
              Mark done
            </button>
          )}
          <button
            type="button"
            className="btn btn--ghost"
            onClick={startEditing}
          >
            Edit
          </button>
          <button
            type="button"
            className="btn btn--danger"
            onClick={() => onRemove(task.id)}
          >
            Delete
          </button>
        </div>
      )}
    </div>
  );
}

export default TaskItem;
