import React, { useState, useEffect } from 'react';

// Nota: actualmente no existe un endpoint de autenticación en el backend del distribuidor.
// Mantengo un fallback local con los usuarios conocidos (dist1/dist2/dist3) y contraseña '1234'.
// Si añades un endpoint /auth en el backend puedo cambiar esto para validar contra el servidor.

const ALLOWED_USERS = ['dist1', 'dist2', 'dist3'];
const DEFAULT_PASSWORD = '1234';

const Login = ({ onLogin }) => {
    const [username, setUsername] = useState(localStorage.getItem('dist_user') || '');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [remember, setRemember] = useState(!!localStorage.getItem('dist_user'));

    useEffect(() => {
        // prefills username when recordado
        if (remember && username) {
            localStorage.setItem('dist_user', username);
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

        // FUTURO: aquí podríamos llamar a un endpoint de autenticación en el backend.
        // Por ahora mantenemos el comportamiento existente (fallback local).
        if (ALLOWED_USERS.includes(username) && password === DEFAULT_PASSWORD) {
            if (remember) localStorage.setItem('dist_user', username);
            else localStorage.removeItem('dist_user');
            onLogin(username);
        } else {
            setError('Usuario o contraseña incorrectos');
        }
    };

    return (
        <div className="login-container">
            <form className="login-form" onSubmit={handleSubmit}>
                <h2>Acceso Distribuidor</h2>
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
                    <small>Usuarios de prueba: dist1 / dist2 / dist3 (contraseña: 1234)</small>
                </div>

                {error && <p className="error-message">{error}</p>}
            </form>
        </div>
    );
};

export default Login;
