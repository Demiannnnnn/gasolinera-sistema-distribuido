
import React, { useState, useEffect } from 'react';
import axios from 'axios';

const Surtidor = () => {
    const [surtidorId, setSurtidorId] = useState('surtidor-1');
    const [ventas, setVentas] = useState([]);
    const [error, setError] = useState(null);
    const [distribuidorUrl, setDistribuidorUrl] = useState('http://localhost:8080'); // Default to distribuidor1

    const surtidores = ['surtidor-1', 'surtidor-2', 'surtidor-3', 'surtidor-4', 'surtidor-5'];
    const distribuidores = [
        { name: 'Distribuidor 1', url: 'http://localhost:8080' },
        { name: 'Distribuidor 2', url: 'http://localhost:8083' },
        { name: 'Distribuidor 3', url: 'http://localhost:8084' },
    ];

    useEffect(() => {
        const fetchVentas = async () => {
            try {
                const response = await axios.get(`${distribuidorUrl}/ventas/recientes?limit=100`);
                const surtidorVentas = response.data.filter(venta => venta.hilo_id === surtidorId);
                setVentas(surtidorVentas);
                setError(null);
            } catch (err) {
                setError('Error al conectar con el distribuidor. Asegúrate de que esté en línea.');
                setVentas([]);
            }
        };

        const interval = setInterval(fetchVentas, 2000);
        return () => clearInterval(interval);
    }, [surtidorId, distribuidorUrl]);

    return (
        <div className="surtidor-container">
            <h2>Monitor de Surtidor</h2>

            <div className="selectors">
                <div className="selector">
                    <label htmlFor="distribuidor-select">Distribuidor:</label>
                    <select id="distribuidor-select" value={distribuidorUrl} onChange={e => setDistribuidorUrl(e.target.value)}>
                        {distribuidores.map(d => (
                            <option key={d.url} value={d.url}>{d.name}</option>
                        ))}
                    </select>
                </div>
                <div className="selector">
                    <label htmlFor="surtidor-select">Surtidor:</label>
                    <select id="surtidor-select" value={surtidorId} onChange={e => setSurtidorId(e.target.value)}>
                        {surtidores.map(id => (
                            <option key={id} value={id}>{id}</option>
                        ))}
                    </select>
                </div>
            </div>

            {error && <p className="error-message">{error}</p>}

            <div className="ventas-list">
                <h3>Últimas Ventas de {surtidorId}</h3>
                {ventas.length > 0 ? (
                    <table>
                        <thead>
                            <tr>
                                <th>ID Venta</th>
                                <th>Combustible</th>
                                <th>Litros</th>
                                <th>Precio/Litro</th>
                                <th>Total</th>
                                <th>Fecha</th>
                            </tr>
                        </thead>
                        <tbody>
                            {ventas.map(venta => (
                                <tr key={venta.id}>
                                    <td>{venta.id}</td>
                                    <td>{venta.tipo_combustible}</td>
                                    <td>{venta.litros_vendidos.toFixed(2)}</td>
                                    <td>${venta.precio_por_litro.toFixed(2)}</td>
                                    <td>${venta.total_venta.toFixed(2)}</td>
                                    <td>{new Date(venta.fecha_venta).toLocaleTimeString()}</td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                ) : (
                    <p>No hay ventas para mostrar.</p>
                )}
            </div>
        </div>
    );
};

export default Surtidor;
