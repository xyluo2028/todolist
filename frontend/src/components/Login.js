import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../api';

const parseStoredAuth = () => {
  try {
    const stored = sessionStorage.getItem('auth');
    return stored ? JSON.parse(stored) : null;
  } catch (err) {
    console.warn('Failed reading stored auth payload, ignoring it.', err);
    return null;
  }
};

function Login() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [auth, setAuth] = useState(() => parseStoredAuth());
  const navigate = useNavigate();

  const handleSubmit = async (event) => {
    event.preventDefault();
    if (!username || !password) {
      alert('Please supply both a username and password.');
      return;
    }

    try {
      await api.getProjects({ username, password });
      const payload = { username, password };
      sessionStorage.setItem('auth', JSON.stringify(payload));
      setAuth(payload);
      navigate('/projects');
    } catch (error) {
      alert('Login failed');
    }
  };

  const handleLogout = () => {
    sessionStorage.removeItem('auth');
    setAuth(null);
    setUsername('');
    setPassword('');
  };

  if (auth) {
    return (
      <div className="form-body">
        <p className="section-subtitle">You are already signed in.</p>
        <div className="form-actions">
          <button
            type="button"
            className="btn btn--secondary"
            onClick={() => navigate('/projects')}
          >
            Go to projects
          </button>
          <button type="button" className="btn btn--danger" onClick={handleLogout}>
            Sign out
          </button>
        </div>
      </div>
    );
  }

  return (
    <form className="form-body" onSubmit={handleSubmit}>
      <div className="form-field">
        <label className="form-label" htmlFor="login-username">
          Username
        </label>
        <input
          id="login-username"
          className="form-input"
          type="text"
          placeholder="you@example.com"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          autoComplete="username"
        />
      </div>
      <div className="form-field">
        <label className="form-label" htmlFor="login-password">
          Password
        </label>
        <input
          id="login-password"
          className="form-input"
          type="password"
          placeholder="Your password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          autoComplete="current-password"
        />
      </div>
      <div className="form-actions">
        <button type="submit" className="btn btn--primary">
          Sign in
        </button>
        <p className="muted-note">Forgot your login? Contact an administrator.</p>
      </div>
    </form>
  );
}

export default Login;
