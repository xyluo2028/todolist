import { useState } from 'react';
import { api } from '../api';
import { useNavigate } from 'react-router-dom';

function Register() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (event) => {
    event.preventDefault();
    if (!username || !password) {
      alert('Please provide both a username and password.');
      return;
    }

    try {
      await api.register(username, password);
      alert('Registration successful! Please log in.');
      navigate('/auth?mode=login');
    } catch (error) {
      alert('Registration failed');
    }
  };

  return (
    <form className="form-body" onSubmit={handleSubmit}>
      <div className="form-field">
        <label className="form-label" htmlFor="register-username">
          Username
        </label>
        <input
          id="register-username"
          className="form-input"
          type="text"
          placeholder="Choose a username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          autoComplete="username"
        />
      </div>
      <div className="form-field">
        <label className="form-label" htmlFor="register-password">
          Password
        </label>
        <input
          id="register-password"
          className="form-input"
          type="password"
          placeholder="Create a strong password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          autoComplete="new-password"
        />
      </div>
      <div className="form-actions">
        <button type="submit" className="btn btn--primary">
          Create account
        </button>
        <p className="muted-note">We never share your credentials with third parties.</p>
      </div>
    </form>
  );
}

export default Register;
