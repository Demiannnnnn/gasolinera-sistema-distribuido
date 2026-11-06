import { render, screen } from '@testing-library/react';
import axios from 'axios';
import App from './App';

jest.mock('axios', () => ({ get: jest.fn() }));

test('muestra el encabezado del monitor y el estado de carga inicial', () => {
  axios.get.mockReturnValue(new Promise(() => {}));
  render(<App />);
  expect(screen.getByText(/Monitor de Distribuidores/i)).toBeInTheDocument();
  expect(screen.getByText(/Conectando con distribuidores/i)).toBeInTheDocument();
});
