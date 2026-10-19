import { useNavigate } from 'react-router-dom';

function ProjectList({ projects, onRemove }) {
  const navigate = useNavigate();

  const goToProject = (name) => {
    navigate(`/projects/${encodeURIComponent(name)}`);
  };

  const handleRemoveClick = (projectName) => onRemove(projectName);

  return (
    <div className="project-list">
      {projects.map((project) => (
        <div
          key={project}
          className="project-card"
        >
          <span className="project-card__title">{project}</span>
          <div className="project-card__actions">
            <button
              type="button"
              className="btn btn--secondary"
              onClick={() => goToProject(project)}
            >
              Open
            </button>
            <button
              type="button"
              className="btn btn--danger"
              onClick={() => handleRemoveClick(project)}
            >
              Remove
            </button>
          </div>
        </div>
      ))}
    </div>
  );
}

export default ProjectList;
