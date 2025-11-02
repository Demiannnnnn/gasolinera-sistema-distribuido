import React, { useState, useEffect } from 'react';
import axios from 'axios';

const Surtidor = ({ distribuidorUrl }) => {
    const [surtidorId, setSurtidorId] = useState('surtidor-1');
    const [ventas, setVentas] = useState([]);
    const [error, setError] = useState(null);

    const surtidores = ['surtidor-1', 'surtidor-2', 'surtidor-3', 'surtidor-4', 'surtidor-5'];

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

        fetchVentas();
    // Polling más frecuente para el surtidor: cada 1 segundo
    const interval = setInterval(fetchVentas, 1000);
        return () => clearInterval(interval);
    }, [surtidorId, distribuidorUrl]);

    // resumen rápido del surtidor seleccionado
    const resumen = {
        totalVentas: ventas.length,
        ingresos: ventas.reduce((s, v) => s + (v.total_venta || 0), 0),
        litros: ventas.reduce((s, v) => s + (v.litros_vendidos || 0), 0),
        ultimaVenta: ventas.length > 0 ? ventas[0].fecha_venta : null
    };

    // datos para sparkline: tomar las últimas 12 ventas y representar total_venta
    const sparkValues = ventas.slice(0, 12).map(v => Number(v.total_venta || 0)).reverse();

    const renderSparkline = (values = []) => {
        const width = 220;
        const height = 48;
        const padding = 4;
        const n = values.length || 1;
        const max = Math.max(...values, 1);

        const points = values.map((v, i) => {
            const x = padding + (i * (width - padding * 2)) / (n - 1 || 1);
            const y = height - padding - (v / max) * (height - padding * 2);
            return `${x},${y}`;
        }).join(' ');

        return (
            <svg className="sparkline" viewBox={`0 0 ${width} ${height}`} preserveAspectRatio="none" aria-hidden>
                <polyline fill="none" stroke="#0ea5a4" strokeWidth="2" points={points} strokeLinecap="round" strokeLinejoin="round" />
            </svg>
        );
    };

    return (
        <div className="surtidor-container">
            <h2>Monitor de Surtidor</h2>

            <div className="selectors">
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

            <div className="surtidor-summary">
                <div className="summary-card">
                    <div className="summary-title">Ventas</div>
                    <div className="summary-value">{resumen.totalVentas}</div>
                </div>
                <div className="summary-card">
                    <div className="summary-title">Ingresos</div>
                    <div className="summary-value">${resumen.ingresos.toFixed(2)}</div>
                </div>
                <div className="summary-card">
                    <div className="summary-title">Litros</div>
                    <div className="summary-value">{resumen.litros.toFixed(2)}L</div>
                </div>
                <div className="summary-card">
                    <div className="summary-title">Última</div>
                    <div className="summary-value">{resumen.ultimaVenta ? new Date(resumen.ultimaVenta).toLocaleTimeString() : 'N/A'}</div>
                </div>
            </div>

            <div className="ventas-list">
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <h3 style={{ margin: 0 }}>Últimas Ventas de {surtidorId}</h3>
                    <div style={{ display: 'flex', gap: 12, alignItems: 'center' }}>
                        <div style={{ fontSize: 12, color: '#6b7280' }}>Últimas {sparkValues.length} ventas</div>
                        {renderSparkline(sparkValues)}
                    </div>
                </div>

                {ventas.length > 0 ? (
                    <div className="ventas-table-wrapper">
                        <table className="ventas-table">
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
                                        <td>{(venta.litros_vendidos || 0).toFixed(2)}</td>
                                        <td>${(venta.precio_por_litro || 0).toFixed(2)}</td>
                                        <td>${(venta.total_venta || 0).toFixed(2)}</td>
                                        <td>{venta.fecha_venta ? new Date(venta.fecha_venta).toLocaleTimeString() : 'Ahora'}</td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                ) : (
                    <p>No hay ventas para mostrar.</p>
                )}
            </div>
        </div>
    );
};

export default Surtidor;
