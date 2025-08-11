import {render, screen, fireEvent} from '@testing-library/react';
import {BrowserRouter} from 'react-router-dom';
import {Button, LinkButton} from './Button';

const renderWithRouter = (component: React.ReactElement) => {
    return render(<BrowserRouter>{component}</BrowserRouter>);
};

describe('Button', () => {
    it('renders button with default props', () => {
        render(<Button>Click me</Button>);

        const button = screen.getByRole('button', {name: 'Click me'});
        expect(button).toBeInTheDocument();
        expect(button).toHaveAttribute('type', 'button');
        expect(button).not.toBeDisabled();
    });

    it('renders button with custom type', () => {
        render(<Button type="submit">Submit</Button>);

        const button = screen.getByRole('button', {name: 'Submit'});
        expect(button).toHaveAttribute('type', 'submit');
    });

    it('renders disabled button', () => {
        render(<Button disabled>Disabled</Button>);

        const button = screen.getByRole('button', {name: 'Disabled'});
        expect(button).toBeDisabled();
    });

    it('calls onClick when clicked', () => {
        const mockOnClick = jest.fn();
        render(<Button onClick={mockOnClick}>Click me</Button>);

        const button = screen.getByRole('button', {name: 'Click me'});
        fireEvent.click(button);

        expect(mockOnClick).toHaveBeenCalledTimes(1);
    });

    it('does not call onClick when disabled', () => {
        const mockOnClick = jest.fn();
        render(<Button onClick={mockOnClick} disabled>Disabled</Button>);

        const button = screen.getByRole('button', {name: 'Disabled'});
        fireEvent.click(button);

        expect(mockOnClick).not.toHaveBeenCalled();
    });

    it('applies primary variant styles by default', () => {
        render(<Button>Primary</Button>);

        const button = screen.getByRole('button', {name: 'Primary'});
        expect(button).toHaveClass('bg-blue-600');
        expect(button).toHaveClass('text-white');
        expect(button).toHaveClass('hover:bg-blue-700');
    });

    it('applies secondary variant styles', () => {
        render(<Button variant="secondary">Secondary</Button>);

        const button = screen.getByRole('button', {name: 'Secondary'});
        expect(button).toHaveClass('bg-gray-100');
        expect(button).toHaveClass('text-gray-900');
        expect(button).toHaveClass('hover:bg-gray-200');
    });

    it('applies outline variant styles', () => {
        render(<Button variant="outline">Outline</Button>);

        const button = screen.getByRole('button', {name: 'Outline'});
        expect(button).toHaveClass('border');
        expect(button).toHaveClass('border-gray-300');
        expect(button).toHaveClass('bg-white');
        expect(button).toHaveClass('text-gray-700');
    });

    it('applies medium size styles by default', () => {
        render(<Button>Medium</Button>);

        const button = screen.getByRole('button', {name: 'Medium'});
        expect(button).toHaveClass('px-4');
        expect(button).toHaveClass('py-2');
        expect(button).toHaveClass('text-base');
    });

    it('applies small size styles', () => {
        render(<Button size="sm">Small</Button>);

        const button = screen.getByRole('button', {name: 'Small'});
        expect(button).toHaveClass('px-3');
        expect(button).toHaveClass('py-2');
        expect(button).toHaveClass('text-sm');
    });

    it('applies large size styles', () => {
        render(<Button size="lg">Large</Button>);

        const button = screen.getByRole('button', {name: 'Large'});
        expect(button).toHaveClass('px-6');
        expect(button).toHaveClass('py-3');
        expect(button).toHaveClass('text-lg');
    });

    it('applies custom className', () => {
        render(<Button className="custom-class">Custom</Button>);

        const button = screen.getByRole('button', {name: 'Custom'});
        expect(button).toHaveClass('custom-class');
    });

    it('applies base classes to all buttons', () => {
        render(<Button>Base</Button>);

        const button = screen.getByRole('button', {name: 'Base'});
        expect(button).toHaveClass('inline-flex');
        expect(button).toHaveClass('items-center');
        expect(button).toHaveClass('justify-center');
        expect(button).toHaveClass('rounded-lg');
        expect(button).toHaveClass('font-semibold');
        expect(button).toHaveClass('transition-colors');
        expect(button).toHaveClass('focus:outline-none');
        expect(button).toHaveClass('focus:ring-2');
        expect(button).toHaveClass('focus:ring-offset-2');
        expect(button).toHaveClass('disabled:opacity-50');
        expect(button).toHaveClass('disabled:cursor-not-allowed');
    });

    it('applies focus ring color for primary variant', () => {
        render(<Button variant="primary">Primary</Button>);

        const button = screen.getByRole('button', {name: 'Primary'});
        expect(button).toHaveClass('focus:ring-blue-500');
    });

    it('applies focus ring color for secondary variant', () => {
        render(<Button variant="secondary">Secondary</Button>);

        const button = screen.getByRole('button', {name: 'Secondary'});
        expect(button).toHaveClass('focus:ring-gray-500');
    });

    it('applies focus ring color for outline variant', () => {
        render(<Button variant="outline">Outline</Button>);

        const button = screen.getByRole('button', {name: 'Outline'});
        expect(button).toHaveClass('focus:ring-gray-500');
    });
});

describe('LinkButton', () => {
    it('renders link button with default props', () => {
        renderWithRouter(<LinkButton href="/test">Link</LinkButton>);

        const link = screen.getByRole('link', {name: 'Link'});
        expect(link).toBeInTheDocument();
        expect(link).toHaveAttribute('href', '/test');
    });

    it('applies primary variant styles by default', () => {
        renderWithRouter(<LinkButton href="/test">Primary Link</LinkButton>);

        const link = screen.getByRole('link', {name: 'Primary Link'});
        expect(link).toHaveClass('bg-blue-600');
        expect(link).toHaveClass('text-white');
        expect(link).toHaveClass('hover:bg-blue-700');
    });

    it('applies secondary variant styles', () => {
        renderWithRouter(<LinkButton href="/test" variant="secondary">Secondary Link</LinkButton>);

        const link = screen.getByRole('link', {name: 'Secondary Link'});
        expect(link).toHaveClass('bg-gray-100');
        expect(link).toHaveClass('text-gray-900');
        expect(link).toHaveClass('hover:bg-gray-200');
    });

    it('applies outline variant styles', () => {
        renderWithRouter(<LinkButton href="/test" variant="outline">Outline Link</LinkButton>);

        const link = screen.getByRole('link', {name: 'Outline Link'});
        expect(link).toHaveClass('border');
        expect(link).toHaveClass('border-gray-300');
        expect(link).toHaveClass('bg-white');
        expect(link).toHaveClass('text-gray-700');
    });

    it('applies medium size styles by default', () => {
        renderWithRouter(<LinkButton href="/test">Medium Link</LinkButton>);

        const link = screen.getByRole('link', {name: 'Medium Link'});
        expect(link).toHaveClass('px-4');
        expect(link).toHaveClass('py-2');
        expect(link).toHaveClass('text-base');
    });

    it('applies small size styles', () => {
        renderWithRouter(<LinkButton href="/test" size="sm">Small Link</LinkButton>);

        const link = screen.getByRole('link', {name: 'Small Link'});
        expect(link).toHaveClass('px-3');
        expect(link).toHaveClass('py-2');
        expect(link).toHaveClass('text-sm');
    });

    it('applies large size styles', () => {
        renderWithRouter(<LinkButton href="/test" size="lg">Large Link</LinkButton>);

        const link = screen.getByRole('link', {name: 'Large Link'});
        expect(link).toHaveClass('px-6');
        expect(link).toHaveClass('py-3');
        expect(link).toHaveClass('text-lg');
    });

    it('applies custom className', () => {
        renderWithRouter(<LinkButton href="/test" className="custom-link">Custom Link</LinkButton>);

        const link = screen.getByRole('link', {name: 'Custom Link'});
        expect(link).toHaveClass('custom-link');
    });

    it('applies base classes to all link buttons', () => {
        renderWithRouter(<LinkButton href="/test">Base Link</LinkButton>);

        const link = screen.getByRole('link', {name: 'Base Link'});
        expect(link).toHaveClass('inline-flex');
        expect(link).toHaveClass('items-center');
        expect(link).toHaveClass('justify-center');
        expect(link).toHaveClass('rounded-lg');
        expect(link).toHaveClass('font-semibold');
        expect(link).toHaveClass('transition-colors');
        expect(link).toHaveClass('focus:outline-none');
        expect(link).toHaveClass('focus:ring-2');
        expect(link).toHaveClass('focus:ring-offset-2');
    });

    it('navigates to correct href', () => {
        renderWithRouter(<LinkButton href="/dashboard">Go to Dashboard</LinkButton>);

        const link = screen.getByRole('link', {name: 'Go to Dashboard'});
        expect(link).toHaveAttribute('href', '/dashboard');
    });

    it('applies focus ring color for different variants', () => {
        renderWithRouter(
            <div>
                <LinkButton href="/test1" variant="primary">Primary</LinkButton>
                <LinkButton href="/test2" variant="secondary">Secondary</LinkButton>
                <LinkButton href="/test3" variant="outline">Outline</LinkButton>
            </div>
        );

        const primaryLink = screen.getByRole('link', {name: 'Primary'});
        const secondaryLink = screen.getByRole('link', {name: 'Secondary'});
        const outlineLink = screen.getByRole('link', {name: 'Outline'});

        expect(primaryLink).toHaveClass('focus:ring-blue-500');
        expect(secondaryLink).toHaveClass('focus:ring-gray-500');
        expect(outlineLink).toHaveClass('focus:ring-gray-500');
    });

    it('renders children content correctly', () => {
        renderWithRouter(
            <LinkButton href="/test">
                <span>Icon</span>
                Link Text
            </LinkButton>
        );

        const link = screen.getByRole('link');
        expect(link).toHaveTextContent('IconLink Text');
        expect(link.querySelector('span')).toHaveTextContent('Icon');
    });
});