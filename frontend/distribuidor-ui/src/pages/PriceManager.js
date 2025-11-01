
import React, { useState } from 'react';

const fuelTypes = [
    { id: '93', name: 'Gasolina 93' },
    { id: '95', name: 'Gasolina 95' },
    { id: '97', name: 'Gasolina 97' },
    { id: 'DIESEL', name: 'Diesel' },
    { id: 'KEROSENE', name: 'Kerosene' },
];

const PriceManager = () => {
    const [prices, setPrices] = useState(
        fuelTypes.reduce((acc, fuel) => {
            acc[fuel.id] = (1.5 + Math.random() * 0.5).toFixed(3); // Initial random price
            return acc;
        }, {})
    );
    const [message, setMessage] = useState('');

    const handlePriceChange = (fuelId, value) => {
        setPrices(prevPrices => ({
            ...prevPrices,
            [fuelId]: value
        }));
    };

    const handleSavePrices = () => {
        // Mock save operation
        console.log('Saving prices:', prices);
        setMessage('Precios actualizados con éxito (simulación). Los distribuidores recibirán la actualización en el próximo ciclo.');
        setTimeout(() => setMessage(''), 5000); // Clear message after 5 seconds
    };

    return (
        <div className="page-content">
            <h1>Gestión de Precios Base</h1>
            <p>Aquí puede definir los precios base para cada tipo de combustible. Estos precios serán enviados a los distribuidores.</p>
            
            <div className="price-list">
                {fuelTypes.map(fuel => (
                    <div className="price-item" key={fuel.id}>
                        <label htmlFor={`price-${fuel.id}`}>{fuel.name}</label>
                        <input 
                            type="number" 
                            id={`price-${fuel.id}`} 
                            value={prices[fuel.id]} 
                            onChange={e => handlePriceChange(fuel.id, e.target.value)}
                            step="0.01"
                        />
                    </div>
                ))}
            </div>

            <button onClick={handleSavePrices} className="save-button">Guardar y Notificar a Distribuidores</button>

            {message && <div className="success-message">{message}</div>}
        </div>
    );
};

export default PriceManager;
