import { render } from '@testing-library/react';
import ChartComponent from './ChartComponent';

// Mock D3 to avoid DOM manipulation issues in tests
jest.mock('d3', () => ({
  select: jest.fn(() => ({
    selectAll: jest.fn(() => ({
      remove: jest.fn(),
    })),
    attr: jest.fn().mockReturnThis(),
    append: jest.fn().mockReturnThis(),
  })),
  scaleBand: jest.fn(() => ({
    domain: jest.fn().mockReturnThis(),
    range: jest.fn().mockReturnThis(),
    padding: jest.fn().mockReturnThis(),
    bandwidth: jest.fn(() => 50),
  })),
  scaleLinear: jest.fn(() => ({
    domain: jest.fn().mockReturnThis(),
    range: jest.fn().mockReturnThis(),
  })),
  max: jest.fn(() => 100),
  axisBottom: jest.fn(),
  axisLeft: jest.fn(),
}));

describe('ChartComponent', () => {
  it('renders without crashing', () => {
    render(<ChartComponent />);
  });

  it('renders an SVG element', () => {
    const { container } = render(<ChartComponent />);
    const svgElement = container.querySelector('svg');
    expect(svgElement).toBeInTheDocument();
  });

  it('has correct CSS classes', () => {
    const { container } = render(<ChartComponent />);
    const svgElement = container.querySelector('svg');
    expect(svgElement).toHaveClass('w-full', 'h-auto');
  });
});