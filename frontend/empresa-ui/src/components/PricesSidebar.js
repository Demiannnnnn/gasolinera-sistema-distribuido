import React from 'react';

const PricesSidebar = ({ distribuidores }) => {
    return (
        <div className="prices-sidebar">
            <h3>Precios Actuales</h3>
            {distribuidores.map(distribuidor => (
                <div key={distribuidor.id} className="distribuidor-prices">
                    <h4>{distribuidor.name}</h4>
                    {distribuidor.status === 'online' && distribuidor.stats?.precios_base ? (
                        <ul>
                            {Object.entries(distribuidor.stats.precios_base).map(([tipo, precio]) => (
                                <li key={tipo}>
                                    <span>{tipo}:</span>
                                    <span>${precio.toFixed(2)}</span>
                                </li>
                            ))}
                        </ul>
                    ) : (
                        <p>No disponible</p>
                    )}
                </div>
            ))}
        </div>
    );
};

export default PricesSidebar;