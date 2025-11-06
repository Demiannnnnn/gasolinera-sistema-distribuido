import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './App.css';
import DistribuidorMonitor from './components/DistribuidorMonitor';
import PricesSidebar from './components/PricesSidebar';
import Login from './pages/Login';

function App() {
  const [user, setUser] = useState(localStorage.getItem('empresa_user') || null);

  const [distribuidores, setDistribuidores] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [paused, setPaused] = useState(false);
  const [query, setQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');

  const distribuidoresConfig = [
      { id: 'DISTRIBUIDOR_1', name: 'Punto de Venta 1', url: process.env.REACT_APP_DIST1_URL || 'http://localhost:8080', color: '#4CAF50' },
      { id: 'DISTRIBUIDOR_2', name: 'Punto de Venta 2', url: process.env.REACT_APP_DIST2_URL || 'http://localhost:8083', color: '#2196F3' },
      { id: 'DISTRIBUIDOR_3', name: 'Punto de Venta 3', url: process.env.REACT_APP_DIST3_URL || 'http://localhost:8084', color: '#9C27B0' }
  ];

  const checkDistribuidoresStatus = async () => {
      const distribuidoresActivos = [];

      for (const config of distribuidoresConfig) {
          try {
              const response = await axios.get(`${config.url}/status`, {
                  timeout: 3000
              });

              if (response.status === 200) {
                  distribuidoresActivos.push({
                      ...config,
                      status: 'online',
                      data: response.data,
                      lastUpdate: new Date()
                  });
              }
          } catch (error) {
              distribuidoresActivos.push({
                  ...config,
                  status: 'offline',
                  data: null,
                  lastUpdate: new Date(),
                  error: error.message
              });
          }
      }
      return distribuidoresActivos;
  };

  const fetchDistribuidorData = async (distribuidor) => {
      try {
          const [ventasResponse, statsResponse] = await Promise.all([
              axios.get(`${distribuidor.url}/ventas/recientes`, { timeout: 5000 }),
              axios.get(`${distribuidor.url}/stats`, { timeout: 5000 })
          ]);

          return {
              ...distribuidor,
              status: 'online',
              ventas: ventasResponse.data || [],
              stats: statsResponse.data || {},
              lastUpdate: new Date()
          };
      } catch (error) {
          console.warn(`Error fetching data from ${distribuidor.name}:`, error.message);
          return {
              ...distribuidor,
              status: 'offline',
              ventas: [],
              stats: {},
              lastUpdate: new Date(),
              error: error.message
          };
      }
  };

  const updateDistribuidoresData = async () => {
      try {
          setLoading(true);
          const distribuidoresActivos = await checkDistribuidoresStatus();
          const distribuidoresConDatos = await Promise.all(
              distribuidoresActivos.map(async (distribuidor) => {
                  if (distribuidor.status === 'online') {
                      return await fetchDistribuidorData(distribuidor);
                  }
                  return distribuidor;
              })
          );
          setDistribuidores(distribuidoresConDatos);
          setError(null);
      } catch (error) {
          console.error('Error updating distribuidores data:', error);
          setError('Error al conectar con los distribuidores');
      } finally {
          setLoading(false);
      }
  };

  const handleUpdatePrices = async (distribuidorId, prices) => {
      try {
          const response = await axios.put(`http://localhost:8081/distribuidores/${distribuidorId}/prices`, { precios: prices });
          if (response.status === 200) {
              alert('Precios actualizados correctamente');
              updateDistribuidoresData();
          }
      } catch (error) {
          console.error('Error updating prices:', error);
          alert('Error al actualizar los precios');
      }
  };

  useEffect(() => {
      updateDistribuidoresData();
      if (paused) return;
      const interval = setInterval(updateDistribuidoresData, 1000);
      return () => clearInterval(interval);
  }, [paused]);

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
      <main className="App-main-layout">
        <DistribuidorMonitor
            distribuidores={distribuidores}
            loading={loading}
            error={error}
            paused={paused}
            setPaused={setPaused}
            query={query}
            setQuery={setQuery}
            statusFilter={statusFilter}
            setStatusFilter={setStatusFilter}
            handleUpdatePrices={handleUpdatePrices}
            getDistribuidoresOnline={() => distribuidores.filter(d => d.status === 'online').length}
            getTotalVentas={() => distribuidores.reduce((total, d) => total + (d.stats?.total_ventas || 0), 0)}
            getTotalIngresos={() => distribuidores.reduce((total, d) => total + (d.stats?.total_ingresos || 0), 0)}
        />
        <PricesSidebar distribuidores={distribuidores} />
      </main>
    </div>
  );
}

export default App;
