import React from 'react';

const DistribuidorCard = ({ distribuidor }) => {
    const {
        id,
        name,
        status,
        ventas = [],
        stats = {},
        lastUpdate,
        error,
        color
    } = distribuidor;

    // Función para formatear tiempo relativo
    const getTimeAgo = (date) => {
        if (!date) return 'Nunca';
        const now = new Date();
        const diff = Math.floor((now - new Date(date)) / 1000);

        if (diff < 60) return `${diff}s`;
        if (diff < 3600) return `${Math.floor(diff / 60)}m`;
        return `${Math.floor(diff / 3600)}h`;
    };

    // Función para obtener la clase CSS del combustible
    const getFuelClass = (tipo) => {
        const fuel = tipo ? tipo.toLowerCase().replace(/[^a-z0-9]/g, '') : '';
        return `fuel-${fuel}`;
    };

    // Función para formatear el tipo de combustible
    const formatFuelType = (tipo) => {
        const fuelMap = {
            '93': 'Gasolina 93',
            '95': 'Gasolina 95',
            '97': 'Gasolina 97',
            'DIESEL': 'Diésel',
            'KEROSENE': 'Keroseno',
            'PETROLEO': 'Petróleo'
        };
        return fuelMap[tipo] || tipo;
    };

    return (
        <div className="distribuidor-card">
            {/* Header del distribuidor */}
            <div className="distribuidor-header">
                <h2 className="distribuidor-title" style={{ color: color }}>
                    {name}
                </h2>
                <span className={`status-indicator status-${status}`}>
                    {status === 'online' ? '🟢 Online' : status === 'offline' ? '🔴 Offline' : '🟡 Cargando'}
                </span>
            </div>

            {/* Estadísticas del distribuidor */}
            {status === 'online' ? (
                <>
                    <div className="stats-grid">
                        <div className="stat-item">
                            <h3 className="stat-value">{stats.total_ventas || 0}</h3>
                            <p className="stat-label">Ventas Hoy</p>
                        </div>
                        <div className="stat-item">
                            <h3 className="stat-value">${(stats.total_ingresos || 0).toFixed(2)}</h3>
                            <p className="stat-label">Ingresos</p>
                        </div>
                        <div className="stat-item">
                            <h3 className="stat-value">{stats.surtidores_activos || 0}</h3>
                            <p className="stat-label">Surtidores</p>
                        </div>
                        <div className="stat-item">
                            <h3 className="stat-value">{(stats.factor_utilidad || 1).toFixed(3)}</h3>
                            <p className="stat-label">Factor Util.</p>
                        </div>
                    </div>

                    {/* Lista de transacciones recientes */}
                    <div className="transactions-section">
                        <h3 className="section-title">
                            🔄 Transacciones Recientes
                            {ventas.length > 0 && (
                                <span style={{ fontSize: '0.8rem', fontWeight: 'normal' }}>
                                    ({ventas.length})
                                </span>
                            )}
                        </h3>

                        <div className="transactions-list">
                            {ventas.length > 0 ? (
                                ventas.slice(0, 10).map((venta, index) => (
                                    <div
                                        key={venta.id || index}
                                        className={`transaction-item ${getFuelClass(venta.tipo_combustible)}`}
                                    >
                                        <div className="transaction-header">
                                            <span className="transaction-fuel">
                                                {formatFuelType(venta.tipo_combustible)}
                                            </span>
                                            <span className="transaction-time">
                                                {venta.fecha_venta ? getTimeAgo(venta.fecha_venta) : 'Ahora'}
                                            </span>
                                        </div>

                                        <div className="transaction-details">
                                            <div className="transaction-detail">
                                                <div className="transaction-detail-label">Litros</div>
                                                <div className="transaction-detail-value">
                                                    {(venta.litros_vendidos || 0).toFixed(2)}L
                                                </div>
                                            </div>
                                            <div className="transaction-detail">
                                                <div className="transaction-detail-label">Precio/L</div>
                                                <div className="transaction-detail-value">
                                                    ${(venta.precio_por_litro || 0).toFixed(2)}
                                                </div>
                                            </div>
                                            <div className="transaction-detail">
                                                <div className="transaction-detail-label">Total</div>
                                                <div className="transaction-detail-value">
                                                    ${(venta.total_venta || 0).toFixed(2)}
                                                </div>
                                            </div>
                                        </div>

                                        {venta.hilo_id && (
                                            <div style={{
                                                marginTop: '8px',
                                                fontSize: '0.8rem',
                                                color: '#666',
                                                textAlign: 'center'
                                            }}>
                                                Surtidor: {venta.hilo_id} • Factura: {venta.numero_factura || 'N/A'}
                                            </div>
                                        )}
                                    </div>
                                ))
                            ) : (
                                <div className="no-data">
                                    <p>No hay transacciones recientes</p>
                                </div>
                            )}
                        </div>
                    </div>
                </>
            ) : (
                /* Estado offline */
                <div className="no-data">
                    <h3>🔌 Distribuidor Desconectado</h3>
                    <p>No se puede conectar al servicio</p>
                    {error && (
                        <div style={{
                            fontSize: '0.8rem',
                            color: '#f44336',
                            marginTop: '10px',
                            padding: '10px',
                            background: '#ffebee',
                            borderRadius: '5px'
                        }}>
                            Error: {error}
                        </div>
                    )}
                </div>
            )}

            {/* Footer con información de última actualización */}
            <div style={{
                marginTop: '15px',
                paddingTop: '15px',
                borderTop: '1px solid #eee',
                fontSize: '0.8rem',
                color: '#666',
                textAlign: 'center'
            }}>
                <div>ID: {id}</div>
                <div>Última actualización: {lastUpdate ? getTimeAgo(lastUpdate) : 'Nunca'}</div>
            </div>
        </div>
    );
};

export default DistribuidorCard;