import React, { useState, useEffect } from 'react';

// Implementación basada en el login del distribuidor.
// Proyecto educativo: única cuenta de prueba para la empresa.
const ALLOWED_USERS = ['empleado'];
const DEFAULT_PASSWORD = '1234';

const Login = ({ onLogin }) => {
  const [username, setUsername] = useState(localStorage.getItem('empresa_user') || '');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [remember, setRemember] = useState(!!localStorage.getItem('empresa_user'));

  useEffect(() => {
    if (remember && username) {
      localStorage.setItem('empresa_user', username);
    }
  }, [remember, username]);

  const handleSubmit = (e) => {
    if (e && e.preventDefault) e.preventDefault();
    setError('');

    if (!username) {
      setError('Por favor ingrese su usuario');
      return;
    }
    if (!password) {
      setError('Por favor ingrese su contraseña');
      return;
    }

    // Comportamiento educativo: validación local usando usuarios/clave conocidos
    if (ALLOWED_USERS.includes(username) && password === DEFAULT_PASSWORD) {
      if (remember) localStorage.setItem('empresa_user', username);
      else localStorage.removeItem('empresa_user');
      onLogin(username);
    } else {
      setError('Usuario o contraseña incorrectos');
    }
  };

  return (
    <div className="login-container">
      <form className="login-form" onSubmit={handleSubmit}>
  <h2>Acceso Empresa</h2>
        <label htmlFor="username">Usuario</label>
        <input
          id="username"
          aria-label="usuario"
          type="text"
          placeholder="ej. dist1"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />

        <label htmlFor="password">Contraseña</label>
        <input
          id="password"
          aria-label="contraseña"
          type="password"
          placeholder="Contraseña"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />

        <div className="login-row">
          <label style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            <input type="checkbox" checked={remember} onChange={(e) => setRemember(e.target.checked)} />
            Recordarme
          </label>
          <button className="login-submit" type="submit" disabled={!username || !password}>Iniciar sesión</button>
        </div>

        <div className="login-help">
          <small>Usuario de prueba: empleado (contraseña: 1234)</small>
        </div>

        {error && <p className="error-message">{error}</p>}
      </form>
    </div>
  );
};

export default Login;
