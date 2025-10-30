
import React from 'react';
import { NavLink } from 'react-router-dom';
import './Sidebar.css';

const Sidebar = () => {
  return (
    <div className="sidebar">
      <div className="sidebar-header">
        <h2>Empresa Central</h2>
      </div>
      <nav className="sidebar-nav">
        <ul>
          <li>
            <NavLink to="/" end>Resumen</NavLink>
          </li>
          <li>
            <NavLink to="/precios">Precios</NavLink>
          </li>
          <li>
            <NavLink to="/reportes">Reportes</NavLink>
          </li>
        </ul>
      </nav>
    </div>
  );
};

export default Sidebar;
