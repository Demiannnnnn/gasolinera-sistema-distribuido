import React from 'react';
import './App.css';
import DistribuidorMonitor from './components/DistribuidorMonitor';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <h1>Panel de Monitoreo Empresarial</h1>
        <p>Visor centralizado de estado y métricas.</p>
      </header>
      <main className="App-main">
        <DistribuidorMonitor />
      </main>
    </div>
  );
}

export default App;
