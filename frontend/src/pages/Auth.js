import React, { useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import Login from '../components/Login';
import Register from '../components/Register';

export default function Auth() {
  const [searchParams, setSearchParams] = useSearchParams();
  const [mode, setMode] = useState('login');

  useEffect(() => {
    const paramMode = searchParams.get('mode') === 'register' ? 'register' : 'login';
    if (paramMode !== mode) {
      setMode(paramMode);
    }
  }, [mode, searchParams]);

  const switchMode = (nextMode) => {
    if (nextMode === mode) {
      return;
    }
    setMode(nextMode);
    setSearchParams({ mode: nextMode }, { replace: true });
  };

  return (
    <div className="form-card">
      <div className="form-header">
        <h2>{mode === 'login' ? 'Welcome back' : 'Create your account'}</h2>
        <p>
          {mode === 'login'
            ? 'Enter your credentials to access your projects and tasks.'
            : 'Fill in the details below to start tracking everything that matters.'}
        </p>
      </div>

      <div className="form-tabs">
        <button
          type="button"
          className={`form-tab${mode === 'login' ? ' form-tab--active' : ''}`}
          onClick={() => switchMode('login')}
        >
          Login
        </button>
        <button
          type="button"
          className={`form-tab${mode === 'register' ? ' form-tab--active' : ''}`}
          onClick={() => switchMode('register')}
        >
          Register
        </button>
      </div>

      <div className="form-body">
        {mode === 'login' ? <Login /> : <Register />}
      </div>
    </div>
  );
}
