import React from 'react';

const Dashboard = () => {
    return (
        <div className="page-content">
            <h1>Dashboard General</h1>
            <p>Resumen consolidado de toda la red de distribuidores.</p>
            
            <div className="stats-grid">
                <div className="stat-card">
                    <h3>Distribuidores Online</h3>
                    <p>Cargando...</p>
                </div>
                <div className="stat-card">
                    <h3>Ventas Totales (unidades)</h3>
                    <p>Cargando...</p>
                </div>
                <div className="stat-card">
                    <h3>Ingresos Totales</h3>
                    <p>Cargando...</p>
                </div>
            </div>
            <p style={{marginTop: '20px', color: 'gray'}}><i>Componente de diagnóstico: La lógica de datos está desactivada temporalmente.</i></p>
        </div>
    );
};

export default Dashboard;