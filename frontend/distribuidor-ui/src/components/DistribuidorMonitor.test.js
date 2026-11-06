import { act, render, screen, waitFor } from '@testing-library/react';
import axios from 'axios';
import DistribuidorMonitor from './DistribuidorMonitor';

jest.mock('axios', () => ({ get: jest.fn() }));

const respuestas = {
  'http://localhost:8080/status': { status: 'online' },
  'http://localhost:8080/ventas/recientes': [
    { id: 1, tipo_combustible: 'DIESEL', litros_vendidos: 20, precio_por_litro: 1.5, total_venta: 30 },
  ],
  'http://localhost:8080/stats': { total_ventas: 12, ventas_hoy: 12, total_ingresos: 300.5, surtidores_activos: 5, factor_utilidad: 1.2 },
  'http://localhost:8083/status': { status: 'online' },
  'http://localhost:8083/ventas/recientes': [],
  'http://localhost:8083/stats': { total_ventas: 8, ventas_hoy: 8, total_ingresos: 99.5, surtidores_activos: 5, factor_utilidad: 1.1 },
};

beforeEach(() => {
  jest.spyOn(console, 'warn').mockImplementation(() => {});
  axios.get.mockImplementation((url) =>
    url in respuestas
      ? Promise.resolve({ status: 200, data: respuestas[url] })
      : Promise.reject(new Error('Network Error'))
  );
});

afterEach(() => {
  jest.restoreAllMocks();
  axios.get.mockReset();
});

test('agrega las estadísticas de los distribuidores en línea y marca los caídos', async () => {
  render(<DistribuidorMonitor />);

  expect(await screen.findByText('2/3')).toBeInTheDocument();
  expect(screen.getByText('20')).toBeInTheDocument(); // 12 + 8 ventas
  expect(screen.getByText('$400.00')).toBeInTheDocument(); // 300.5 + 99.5
  expect(screen.getByText('Activo')).toBeInTheDocument();
  expect(screen.getAllByText('🟢 Online')).toHaveLength(2);
  expect(screen.getByText('🔴 Offline')).toBeInTheDocument();
  expect(screen.getByText('Error: Network Error')).toBeInTheDocument();
  expect(screen.getByText('Diésel')).toBeInTheDocument();
});

test('si /status responde pero /stats falla, el distribuidor queda offline', async () => {
  axios.get.mockImplementation((url) =>
    url === 'http://localhost:8080/status'
      ? Promise.resolve({ status: 200, data: {} })
      : Promise.reject(new Error('fallo'))
  );
  render(<DistribuidorMonitor />);

  expect(await screen.findByText('0/3')).toBeInTheDocument();
  expect(screen.getByText('Inactivo')).toBeInTheDocument();
  expect(screen.getAllByText('🔴 Offline')).toHaveLength(3);
});

test('consulta periódicamente y deja de hacerlo al desmontarse', async () => {
  jest.useFakeTimers();
  try {
    const { unmount } = render(<DistribuidorMonitor />);
    await screen.findByText('2/3');
    const llamadasIniciales = axios.get.mock.calls.length;

    await act(async () => { jest.advanceTimersByTime(3000); });
    await waitFor(() => expect(axios.get.mock.calls.length).toBeGreaterThan(llamadasIniciales));

    unmount();
    const trasDesmontar = axios.get.mock.calls.length;
    jest.advanceTimersByTime(10000);
    expect(axios.get.mock.calls.length).toBe(trasDesmontar);
  } finally {
    jest.useRealTimers();
  }
});
