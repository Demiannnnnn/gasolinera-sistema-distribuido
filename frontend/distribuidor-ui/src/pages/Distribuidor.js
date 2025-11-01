import React, { useState, useEffect } from 'react';
import axios from 'axios';
import Surtidor from '../components/Surtidor';

const Distribuidor = ({ distributorId }) => {
    const [distribuidor, setDistribuidor] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    const distribuidoresConfig = [
        { id: 'DISTRIBUIDOR_1', name: 'Distribuidor 1', url: process.env.REACT_APP_DIST1_URL || 'http://localhost:8080', color: '#4CAF50' },
        { id: 'DISTRIBUIDOR_2', name: 'Distribuidor 2', url: process.env.REACT_APP_DIST2_URL || 'http://localhost:8083', color: '#2196F3' },
        { id: 'DISTRIBUIDOR_3', name: 'Distribuidor 3', url: process.env.REACT_APP_DIST3_URL || 'http://localhost:8084', color: '#9C27B0' }
    ];

    const fetchDistribuidorData = async () => {
        const config = distribuidoresConfig.find(d => d.id === distributorId);
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
            setDistribuidor({
                ...config,
                status: 'offline',
                ventas: [],
                stats: {},
                lastUpdate: new Date(),
                error: error.message
            });
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchDistribuidorData();
        const interval = setInterval(fetchDistribuidorData, 3000);
        return () => clearInterval(interval);
    }, [distributorId]);

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
        <div>
            <h1>{distribuidor.name}</h1>
            <Surtidor distribuidorUrl={distribuidor.url} />
        </div>
    );
};

export default Distribuidor;
