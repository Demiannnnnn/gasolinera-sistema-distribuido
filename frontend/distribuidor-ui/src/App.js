
import React, { useState } from 'react';
import Distribuidor from './pages/Distribuidor';
import Login from './pages/Login';
import './App.css';

function App() {
  const [user, setUser] = useState(null);

  const handleLogin = (username) => {
    setUser(username);
  };

  const getDistributorId = () => {
    switch (user) {
      case 'dist1':
        return 'DISTRIBUIDOR_1';
      case 'dist2':
        return 'DISTRIBUIDOR_2';
      case 'dist3':
        return 'DISTRIBUIDOR_3';
      default:
        return null;
    }
  };

  return (
    <div className="App">
        <main className="main-content">
            {user ? (
                <Distribuidor distributorId={getDistributorId()} />
            ) : (
                <Login onLogin={handleLogin} />
            )}
        </main>
    </div>
  );
}

export default App;
