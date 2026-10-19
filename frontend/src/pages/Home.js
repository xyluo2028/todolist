import { useEffect, useState } from 'react';
import { Link, useLocation } from 'react-router-dom';

const readStoredAuth = () => {
  try {
    const stored = sessionStorage.getItem('auth');
    return stored ? JSON.parse(stored) : null;
  } catch (err) {
    console.warn('Failed parsing stored auth credentials in Home.', err);
    return null;
  }
};

function Home() {
  const location = useLocation();
  const [auth, setAuth] = useState(() => readStoredAuth());

  useEffect(() => {
    setAuth(readStoredAuth());
  }, [location]);

  return (
    <section className="hero">
      <div className="hero__content">
        <span className="hero__eyebrow">Plan smarter</span>
        <h1 className="hero__title">
          Organize <span>every project</span> with clarity
        </h1>
        <p className="hero__subtitle">
          TaskFlow keeps your to-dos structured, deadlines visible, and priorities under
          control so you can focus on meaningful work.
        </p>
        {!auth && (
          <div className="hero__actions">
            <Link to="/auth?mode=register" className="btn btn--primary">
              Create an account
            </Link>
          </div>
        )}

        <div className="hero__stats">
          <div className="hero__stat">
            <strong>Unlimited projects</strong>
            <span>Keep personal and team work separated with ease.</span>
          </div>
          <div className="hero__stat">
            <strong>Deadline ready</strong>
            <span>Visualize due dates and react before they sneak up.</span>
          </div>
          <div className="hero__stat">
            <strong>Focus mode</strong>
            <span>Check off tasks and celebrate real progress.</span>
          </div>
        </div>
      </div>
    </section>
  );
}

export default Home;
