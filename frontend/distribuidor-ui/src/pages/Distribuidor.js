import React, { useState, useEffect } from 'react';
import axios from 'axios';
import Surtidor from '../components/Surtidor';

const Distribuidor = ({ distributorId }) => {
    const [distribuidor, setDistribuidor] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [paused, setPaused] = useState(false); // Pausa de la UI
    const [simulationPaused, setSimulationPaused] = useState(false); // Pausa de la simulación

    const distribuidoresConfig = [
        { id: 'DISTRIBUIDOR_1', name: 'Distribuidor 1', url: process.env.REACT_APP_DIST1_URL || 'http://localhost:8080', color: '#4CAF50' },
        { id: 'DISTRIBUIDOR_2', name: 'Distribuidor 2', url: process.env.REACT_APP_DIST2_URL || 'http://localhost:8083', color: '#2196F3' },
        { id: 'DISTRIBUIDOR_3', name: 'Distribuidor 3', url: process.env.REACT_APP_DIST3_URL || 'http://localhost:8084', color: '#9C27B0' }
    ];

    const getConfig = () => distribuidoresConfig.find(d => d.id === distributorId);

    const fetchDistribuidorData = async () => {
        const config = getConfig();
        if (!config) {
            setError('Distribuidor no encontrado');
            setLoading(false);
            return;
        }

        try {
            const [ventasResponse, statsResponse] = await Promise.all([
                axios.get(`${config.url}/ventas/recientes`, { timeout: 5000 }),
                axios.get(`${config.url}/stats`, { timeout: 5000 })
            ]);

            setDistribuidor({
                ...config,
                status: 'online',
                ventas: ventasResponse.data || [],
                stats: statsResponse.data || {},
                lastUpdate: new Date()
            });
            setError(null);
        } catch (error) {
            console.warn(`Error fetching data from ${config.name}:`, error.message);
            setDistribuidor(prev => ({ ...prev, ...config, status: 'offline', lastUpdate: new Date(), error: error.message }));
        } finally {
            setLoading(false);
        }
    };

    // Nueva función para controlar la simulación
    const toggleSimulation = async () => {
        const config = getConfig();
        const action = simulationPaused ? 'resume' : 'pause';
        try {
            await axios.post(`${config.url}/surtidores/control`, { action });
            setSimulationPaused(!simulationPaused);
        } catch (error) {
            console.error(`Error al intentar ${action} la simulación:`, error);
            alert(`No se pudo ${action} la simulación.`);
        }
    };

    useEffect(() => {
        fetchDistribuidorData();
        if (paused) return;
        const interval = setInterval(fetchDistribuidorData, 1000);
        return () => clearInterval(interval);
    }, [distributorId, paused]);

    if (loading) {
        return (
            <div className="monitor-container">
                <div className="no-data">
                    <div className="loading-spinner"></div>
                    <p>Conectando con el distribuidor...</p>
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="error-message">
                <strong>Error:</strong> {error}
            </div>
        );
    }

    if (!distribuidor) {
        return null;
    }

    return (
        <div className="distribuidor-page">
            <header className="dist-header">
                <div>
                    <h1 className="dist-title">{distribuidor.name}</h1>
                    <div className="dist-sub">Estado: <span className={`status-badge status-${distribuidor.status}`}>{distribuidor.status}</span></div>
                    <div className="dist-sub">Última actualización: {distribuidor.lastUpdate ? new Date(distribuidor.lastUpdate).toLocaleTimeString() : 'Nunca'}</div>
                </div>

                <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
                    <button className="btn" onClick={() => fetchDistribuidorData()}>Actualizar</button>
                    {/* Botón para pausar la UI */}
                    <button className={`btn ${paused ? '' : 'btn-primary'}`} onClick={() => setPaused(!paused)}>{paused ? 'Reanudar UI' : 'Pausar UI'}</button>
                    {/* NUEVO Botón para pausar la simulación */}
                    <button className={`btn ${simulationPaused ? '' : 'btn-danger'}`} onClick={toggleSimulation}>{simulationPaused ? 'Reanudar Simulación' : 'Pausar Simulación'}</button>
                </div>
            </header>

            {/* Estadísticas globales */}
            <section className="stats-grid">
                <div className="stat-card">
                    <h3>Ventas hoy</h3>
                    <p className="stat-value">{distribuidor.stats?.total_ventas || 0}</p>
                </div>
                <div className="stat-card">
                    <h3>Ingresos</h3>
                    <p className="stat-value">${(distribuidor.stats?.total_ingresos || 0).toFixed(2)}</p>
                </div>
                <div className="stat-card">
                    <h3>Surtidores activos</h3>
                    <p className="stat-value">{distribuidor.stats?.surtidores_activos || 0}</p>
                </div>
                <div className="stat-card">
                    <h3>Factor util.</h3>
                    <p className="stat-value">{(distribuidor.stats?.factor_utilidad || 1).toFixed(3)}</p>
                </div>
            </section>

            <section className="stats-grid">
                {distribuidor.stats?.precios_base && Object.entries(distribuidor.stats.precios_base).map(([tipo, precio]) => (
                    <div className="stat-card" key={tipo}>
                        <h3>Precio {tipo}</h3>
                        <p className="stat-value">${precio.toFixed(2)}</p>
                    </div>
                ))}
            </section>

            <main className="dist-main">
                <aside className="surtidor-aside">
                    <h3>Surtidores</h3>
                    {/* render Surtidor component which provides selector + table */}
                    <Surtidor distribuidorUrl={distribuidor.url} />
                </aside>

                <section className="dist-info">
                    <h3>Actividad reciente</h3>
                    <div className="info-placeholder">
                        {/* space to later add charts or transaction stream */}
                        <p>Resumen y gráficos rápidos serán agregados aquí.</p>
                    </div>
                </section>
            </main>
        </div>
    );
};

export default Distribuidor;
