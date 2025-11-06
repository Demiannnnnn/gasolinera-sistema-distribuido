import React, { useState } from 'react';

const PriceUpdateForm = ({ distribuidor, onUpdatePrices }) => {
    const [prices, setPrices] = useState({
        '93': 1.0,
        '95': 1.1,
        '97': 1.2,
        'DIESEL': 0.9,
        'KEROSENE': 0.8,
    });

    const handleChange = (e) => {
        const { name, value } = e.target;
        setPrices(prevPrices => ({
            ...prevPrices,
            [name]: parseFloat(value)
        }));
    };

    const handleSubmit = (e) => {
        e.preventDefault();
        onUpdatePrices(distribuidor.id, prices);
    };

    return (
        <form onSubmit={handleSubmit} className="price-update-form">
            <h4>Actualizar Precios Base</h4>
            <div className="price-inputs">
                {Object.keys(prices).map(fuelType => (
                    <div key={fuelType} className="price-input">
                        <label>{fuelType}</label>
                        <input
                            type="number"
                            name={fuelType}
                            value={prices[fuelType]}
                            onChange={handleChange}
                            step="0.01"
                        />
                    </div>
                ))}
            </div>
            <button type="submit" className="btn btn-primary">Actualizar</button>
        </form>
    );
};

export default PriceUpdateForm;