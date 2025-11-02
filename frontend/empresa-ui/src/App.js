import React, { useState } from 'react';
import './App.css';
import DistribuidorMonitor from './components/DistribuidorMonitor';
import Login from './pages/Login';

function App() {
  const [user, setUser] = useState(localStorage.getItem('empresa_user') || null);

  const handleLogin = (username) => {
    setUser(username);
  };

  const handleLogout = () => {
    localStorage.removeItem('empresa_user');
    setUser(null);
  };

  if (!user) {
    return (
      <div className="App">
        <header className="App-header">
          <h1>Panel de Monitoreo Empresarial</h1>
          <p>Por favor, inicia sesión para continuar.</p>
        </header>
        <main className="App-main">
          <Login onLogin={handleLogin} />
        </main>
      </div>
    );
  }

  return (
    <div className="App">
      <header className="App-header">
        <h1>Panel de Monitoreo Empresarial</h1>
        <p>Bienvenido <strong>{user}</strong></p>
        <button onClick={handleLogout} style={{ position: 'absolute', right: 16, top: 16 }}>Cerrar sesión</button>
      </header>
      <main className="App-main">
        <DistribuidorMonitor />
      </main>
    </div>
  );
}

export default App;
