import { render, screen } from '@testing-library/react';
import DistribuidorCard from './DistribuidorCard';

const base = {
  id: 'DISTRIBUIDOR_1',
  name: 'Distribuidor 1',
  color: '#4CAF50',
};

const venta = (id, extra = {}) => ({
  id,
  tipo_combustible: '95',
  litros_vendidos: 10,
  precio_por_litro: 1.5,
  total_venta: 15,
  hilo_id: 'surtidor-1',
  numero_factura: `F-DISTRIBUIDOR_1-${id}`,
  ...extra,
});

describe('DistribuidorCard', () => {
  test('online: muestra estadísticas formateadas', () => {
    render(
      <DistribuidorCard
        distribuidor={{
          ...base,
          status: 'online',
          stats: { total_ventas: 42, ventas_hoy: 7, total_ingresos: 1234.5, surtidores_activos: 5, factor_utilidad: 1.2 },
          ventas: [],
        }}
      />
    );
    expect(screen.getByText('🟢 Online')).toBeInTheDocument();
    expect(screen.getByText('7')).toBeInTheDocument(); // ventas de hoy, no el total histórico
    expect(screen.getByText('$1234.50')).toBeInTheDocument();
    expect(screen.getByText('5')).toBeInTheDocument();
    expect(screen.getByText('1.200')).toBeInTheDocument();
    expect(screen.getByText('No hay transacciones recientes')).toBeInTheDocument();
  });

  test('online sin stats usa valores por defecto', () => {
    render(<DistribuidorCard distribuidor={{ ...base, status: 'online' }} />);
    expect(screen.getByText('$0.00')).toBeInTheDocument();
    expect(screen.getByText('1.000')).toBeInTheDocument();
  });

  test('traduce los tipos de combustible y muestra detalles de la venta', () => {
    const tipos = ['93', '95', '97', 'DIESEL', 'KEROSENE', 'PETROLEO'];
    render(
      <DistribuidorCard
        distribuidor={{ ...base, status: 'online', ventas: tipos.map((t, i) => venta(i + 1, { tipo_combustible: t })) }}
      />
    );
    ['Gasolina 93', 'Gasolina 95', 'Gasolina 97', 'Diésel', 'Keroseno', 'Petróleo'].forEach((label) =>
      expect(screen.getByText(label)).toBeInTheDocument()
    );
    expect(screen.getByText(/Factura: F-DISTRIBUIDOR_1-1$/)).toBeInTheDocument();
    expect(screen.getAllByText('10.00L')).toHaveLength(6);
  });

  test('tipo de combustible desconocido se muestra tal cual y con clase CSS saneada', () => {
    const { container } = render(
      <DistribuidorCard distribuidor={{ ...base, status: 'online', ventas: [venta(1, { tipo_combustible: 'GAS-X' })] }} />
    );
    expect(screen.getByText('GAS-X')).toBeInTheDocument();
    expect(container.querySelector('.transaction-item.fuel-gasx')).not.toBeNull();
  });

  test('muestra como máximo 10 transacciones pero indica el total', () => {
    const ventas = Array.from({ length: 15 }, (_, i) => venta(i + 1));
    const { container } = render(<DistribuidorCard distribuidor={{ ...base, status: 'online', ventas }} />);
    expect(container.querySelectorAll('.transaction-item')).toHaveLength(10);
    expect(screen.getByText('(15)')).toBeInTheDocument();
  });

  test('tiempo relativo en segundos, minutos y horas', () => {
    const ahora = Date.now();
    const ventas = [
      venta(1, { fecha_venta: new Date(ahora - 30 * 1000).toISOString() }),
      venta(2, { fecha_venta: new Date(ahora - 5 * 60 * 1000).toISOString() }),
      venta(3, { fecha_venta: new Date(ahora - 3 * 3600 * 1000).toISOString() }),
      venta(4, { fecha_venta: undefined }),
    ];
    render(<DistribuidorCard distribuidor={{ ...base, status: 'online', ventas }} />);
    expect(screen.getByText(/^(30|31)s$/)).toBeInTheDocument();
    expect(screen.getByText('5m')).toBeInTheDocument();
    expect(screen.getByText('3h')).toBeInTheDocument();
    expect(screen.getByText('Ahora')).toBeInTheDocument();
    expect(screen.getByText('Última actualización: Nunca')).toBeInTheDocument();
  });

  test('offline: muestra desconexión y el error', () => {
    render(<DistribuidorCard distribuidor={{ ...base, status: 'offline', error: 'timeout of 3000ms exceeded' }} />);
    expect(screen.getByText('🔴 Offline')).toBeInTheDocument();
    expect(screen.getByText(/Distribuidor Desconectado/)).toBeInTheDocument();
    expect(screen.getByText('Error: timeout of 3000ms exceeded')).toBeInTheDocument();
    expect(screen.queryByText('Ingresos')).toBeNull();
  });

  test('estado desconocido se muestra como cargando', () => {
    render(<DistribuidorCard distribuidor={{ ...base, status: 'loading' }} />);
    expect(screen.getByText('🟡 Cargando')).toBeInTheDocument();
    expect(screen.getByText('ID: DISTRIBUIDOR_1')).toBeInTheDocument();
  });
});
