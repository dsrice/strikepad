import {render, screen} from '@testing-library/react';
import {Container} from './Container';

describe('Container', () => {
    it('renders children correctly', () => {
        render(
            <Container>
                <div>Test content</div>
            </Container>
        );

        expect(screen.getByText('Test content')).toBeInTheDocument();
    });

    it('applies default container classes', () => {
        render(
            <Container>
                <div data-testid="content">Test content</div>
            </Container>
        );

        const container = screen.getByTestId('content').parentElement;
        expect(container).toHaveClass('mx-auto');
        expect(container).toHaveClass('max-w-7xl');
        expect(container).toHaveClass('px-6');
        expect(container).toHaveClass('lg:px-8');
    });

    it('applies custom className along with default classes', () => {
        render(
            <Container className="custom-class">
                <div data-testid="content">Test content</div>
            </Container>
        );

        const container = screen.getByTestId('content').parentElement;
        expect(container).toHaveClass('mx-auto');
        expect(container).toHaveClass('max-w-7xl');
        expect(container).toHaveClass('px-6');
        expect(container).toHaveClass('lg:px-8');
        expect(container).toHaveClass('custom-class');
    });

    it('renders with multiple custom classes', () => {
        render(
            <Container className="py-8 bg-gray-100 custom-container">
                <div data-testid="content">Test content</div>
            </Container>
        );

        const container = screen.getByTestId('content').parentElement;
        expect(container).toHaveClass('py-8');
        expect(container).toHaveClass('bg-gray-100');
        expect(container).toHaveClass('custom-container');
        expect(container).toHaveClass('mx-auto');
        expect(container).toHaveClass('max-w-7xl');
    });

    it('renders without custom className', () => {
        render(
            <Container>
                <div data-testid="content">Test content</div>
            </Container>
        );

        const container = screen.getByTestId('content').parentElement;
        expect(container).toHaveClass('mx-auto', 'max-w-7xl', 'px-6', 'lg:px-8');
    });

    it('renders complex children correctly', () => {
        render(
            <Container>
                <header>Header content</header>
                <main>
                    <section>
                        <h1>Title</h1>
                        <p>Paragraph</p>
                    </section>
                </main>
                <footer>Footer content</footer>
            </Container>
        );

        expect(screen.getByText('Header content')).toBeInTheDocument();
        expect(screen.getByText('Title')).toBeInTheDocument();
        expect(screen.getByText('Paragraph')).toBeInTheDocument();
        expect(screen.getByText('Footer content')).toBeInTheDocument();

        expect(screen.getByRole('banner')).toBeInTheDocument();
        expect(screen.getByRole('main')).toBeInTheDocument();
        expect(screen.getByRole('contentinfo')).toBeInTheDocument();
    });

    it('has correct HTML structure', () => {
        const {container} = render(
            <Container>
                <div>Content</div>
            </Container>
        );

        const containerDiv = container.firstChild;
        expect(containerDiv).toBeInstanceOf(HTMLDivElement);
        expect(containerDiv?.textContent).toBe('Content');
    });

    it('handles empty children', () => {
        const {container} = render(<Container></Container>);

        const containerDiv = container.firstChild;
        expect(containerDiv).toBeInstanceOf(HTMLDivElement);
        expect(containerDiv?.textContent).toBe('');
    });

    it('handles null children gracefully', () => {
        render(
            <Container>
                {null}
                <div>Visible content</div>
                {undefined}
            </Container>
        );

        expect(screen.getByText('Visible content')).toBeInTheDocument();
    });

    it('handles React fragments as children', () => {
        render(
            <Container>
                <>
                    <div>Fragment content 1</div>
                    <div>Fragment content 2</div>
                </>
            </Container>
        );

        expect(screen.getByText('Fragment content 1')).toBeInTheDocument();
        expect(screen.getByText('Fragment content 2')).toBeInTheDocument();
    });

    it('preserves event handlers on children', () => {
        const handleClick = jest.fn();

        render(
            <Container>
                <button onClick={handleClick}>Click me</button>
            </Container>
        );

        const button = screen.getByRole('button', {name: 'Click me'});
        button.click();

        expect(handleClick).toHaveBeenCalledTimes(1);
    });

    it('allows nested containers', () => {
        render(
            <Container className="outer-container">
                <Container className="inner-container">
                    <div data-testid="nested-content">Nested content</div>
                </Container>
            </Container>
        );

        const nestedContent = screen.getByTestId('nested-content');
        const innerContainer = nestedContent.parentElement;
        const outerContainer = innerContainer?.parentElement;

        expect(innerContainer).toHaveClass('inner-container');
        expect(outerContainer).toHaveClass('outer-container');

        expect(innerContainer).toHaveClass('mx-auto', 'max-w-7xl', 'px-6', 'lg:px-8');
        expect(outerContainer).toHaveClass('mx-auto', 'max-w-7xl', 'px-6', 'lg:px-8');
    });
});