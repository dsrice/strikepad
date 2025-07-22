import { render, screen, fireEvent } from '@testing-library/react';
import App from './App';

describe('App', () => {
  it('renders StrikePad Frontend title', () => {
    render(<App />);
    const titleElement = screen.getByText(/StrikePad Frontend/i);
    expect(titleElement).toBeInTheDocument();
  });

  it('renders counter button with initial count', () => {
    render(<App />);
    const buttonElement = screen.getByText(/Count is 0/i);
    expect(buttonElement).toBeInTheDocument();
  });

  it('increments counter when button is clicked', () => {
    render(<App />);
    const buttonElement = screen.getByText(/Count is 0/i);
    
    fireEvent.click(buttonElement);
    expect(screen.getByText(/Count is 1/i)).toBeInTheDocument();
    
    fireEvent.click(buttonElement);
    expect(screen.getByText(/Count is 2/i)).toBeInTheDocument();
  });

  it('renders chart demo section', () => {
    render(<App />);
    const chartSection = screen.getByText(/D3.js Chart Demo/i);
    expect(chartSection).toBeInTheDocument();
  });

  it('has correct styling classes', () => {
    render(<App />);
    const mainElement = screen.getByRole('main');
    expect(mainElement).toHaveClass('max-w-6xl', 'mx-auto', 'p-6');
  });
});