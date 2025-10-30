import React, { useState, useEffect } from 'react';
import axios from 'axios';
import DistribuidorCard from './DistribuidorCard';

const DistribuidorMonitor = () => {
    const [distribuidores, setDistribuidores] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    // Lista de distribuidores posibles (dinámico)
    const distribuidoresConfig = [
        { id: 'DISTRIBUIDOR_1', name: 'Distribuidor 1', url: process.env.REACT_APP_DIST1_URL || 'http://localhost:8080', color: '#4CAF50' },
        { id: 'DISTRIBUIDOR_2', name: 'Distribuidor 2', url: process.env.REACT_APP_DIST2_URL || 'http://localhost:8083', color: '#2196F3' },
        { id: 'DISTRIBUIDOR_3', name: 'Distribuidor 3', url: process.env.REACT_APP_DIST3_URL || 'http://localhost:8084', color: '#9C27B0' }
    ];

    // Verificar qué distribuidores están activos
    const checkDistribuidoresStatus = async () => {
        const distribuidoresActivos = [];

        for (const config of distribuidoresConfig) {
            try {
                // Intentar hacer ping al distribuidor
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
                // Si falla, aún lo agregamos pero como offline
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

    // Obtener datos de un distribuidor específico
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

    // Actualizar datos de todos los distribuidores
    const updateDistribuidoresData = async () => {
        try {
            setLoading(true);

            // Primero verificar cuáles están activos
            const distribuidoresActivos = await checkDistribuidoresStatus();

            // Luego obtener datos detallados solo de los que están online
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

    // Efecto para cargar datos iniciales y configurar polling
    useEffect(() => {
        updateDistribuidoresData();

        // Actualizar cada 3 segundos
        const interval = setInterval(updateDistribuidoresData, 3000);

        return () => clearInterval(interval);
    }, []);

    // Funciones de utilidad
    const getDistribuidoresOnline = () => distribuidores.filter(d => d.status === 'online').length;
    const getTotalVentas = () => distribuidores.reduce((total, d) => total + (d.stats?.total_ventas || 0), 0);
    const getTotalIngresos = () => distribuidores.reduce((total, d) => total + (d.stats?.total_ingresos || 0), 0);

    if (loading && distribuidores.length === 0) {
        return (
            <div className="monitor-container">
                <div className="no-data">
                    <div className="loading-spinner"></div>
                    <p>Conectando con distribuidores...</p>
                </div>
            </div>
        );
    }

    return (
        <div>
            {/* Panel de estadísticas globales */}
            <div className="global-stats">
                <div className="stats-grid">
                    <div className="stat-item">
                        <h3 className="stat-value">{getDistribuidoresOnline()}/{distribuidores.length}</h3>
                        <p className="stat-label">Distribuidores Online</p>
                    </div>
                    <div className="stat-item">
                        <h3 className="stat-value">{getTotalVentas()}</h3>
                        <p className="stat-label">Total Ventas</p>
                    </div>
                    <div className="stat-item">
                        <h3 className="stat-value">${getTotalIngresos().toFixed(2)}</h3>
                        <p className="stat-label">Ingresos Totales</p>
                    </div>
                    <div className="stat-item">
                        <h3 className="stat-value">
                            {distribuidores.filter(d => d.status === 'online').length > 0 ? 'Activo' : 'Inactivo'}
                        </h3>
                        <p className="stat-label">Estado Sistema</p>
                    </div>
                </div>
            </div>

            {error && (
                <div className="error-message">
                    <strong>Error:</strong> {error}
                </div>
            )}

            {/* Grid de distribuidores */}
            <div className="monitor-container">
                {distribuidores.map(distribuidor => (
                    <DistribuidorCard
                        key={distribuidor.id}
                        distribuidor={distribuidor}
                    />
                ))}
            </div>

            {distribuidores.length === 0 && !loading && (
                <div className="no-data">
                    <h3>No se encontraron distribuidores</h3>
                    <p>Verifica que los servicios estén corriendo</p>
                </div>
            )}

            {/* Información de última actualización */}
            <div style={{ textAlign: 'center', marginTop: '20px', color: '#666', fontSize: '0.9rem' }}>
                Última actualización: {new Date().toLocaleTimeString()}
                {loading && <span> (Actualizando...)</span>}
            </div>
        </div>
    );
};

export default DistribuidorMonitor;