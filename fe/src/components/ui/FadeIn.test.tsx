import {render, screen} from '@testing-library/react';
import {FadeIn, FadeInStagger} from './FadeIn';

describe('FadeIn', () => {
    it('renders children correctly', () => {
        render(
            <FadeIn>
                <div>Fade in content</div>
            </FadeIn>
        );

        expect(screen.getByText('Fade in content')).toBeInTheDocument();
    });

    it('applies default animation classes', () => {
        render(
            <FadeIn>
                <div data-testid="content">Test content</div>
            </FadeIn>
        );

        const container = screen.getByTestId('content').parentElement;
        expect(container).toHaveClass('animate-in');
        expect(container).toHaveClass('fade-in');
        expect(container).toHaveClass('duration-700');
    });

    it('applies custom className along with default classes', () => {
        render(
            <FadeIn className="custom-fade">
                <div data-testid="content">Test content</div>
            </FadeIn>
        );

        const container = screen.getByTestId('content').parentElement;
        expect(container).toHaveClass('animate-in');
        expect(container).toHaveClass('fade-in');
        expect(container).toHaveClass('duration-700');
        expect(container).toHaveClass('custom-fade');
    });

    it('applies multiple custom classes', () => {
        render(
            <FadeIn className="delay-100 slide-in-from-bottom-4">
                <div data-testid="content">Test content</div>
            </FadeIn>
        );

        const container = screen.getByTestId('content').parentElement;
        expect(container).toHaveClass('animate-in');
        expect(container).toHaveClass('fade-in');
        expect(container).toHaveClass('duration-700');
        expect(container).toHaveClass('delay-100');
        expect(container).toHaveClass('slide-in-from-bottom-4');
    });

    it('renders without custom className', () => {
        render(
            <FadeIn>
                <div data-testid="content">Test content</div>
            </FadeIn>
        );

        const container = screen.getByTestId('content').parentElement;
        expect(container).toHaveClass('animate-in', 'fade-in', 'duration-700');
    });

    it('renders complex children correctly', () => {
        render(
            <FadeIn>
                <div>
                    <h1>Title</h1>
                    <p>Description</p>
                    <button>Action</button>
                </div>
            </FadeIn>
        );

        expect(screen.getByText('Title')).toBeInTheDocument();
        expect(screen.getByText('Description')).toBeInTheDocument();
        expect(screen.getByRole('button', {name: 'Action'})).toBeInTheDocument();
    });

    it('has correct HTML structure', () => {
        const {container} = render(
            <FadeIn>
                <div>Content</div>
            </FadeIn>
        );

        const fadeInDiv = container.firstChild;
        expect(fadeInDiv).toBeInstanceOf(HTMLDivElement);
        expect(fadeInDiv?.textContent).toBe('Content');
    });

    it('preserves event handlers on children', () => {
        const handleClick = jest.fn();

        render(
            <FadeIn>
                <button onClick={handleClick}>Click me</button>
            </FadeIn>
        );

        const button = screen.getByRole('button', {name: 'Click me'});
        button.click();

        expect(handleClick).toHaveBeenCalledTimes(1);
    });
});

describe('FadeInStagger', () => {
    it('renders children correctly', () => {
        render(
            <FadeInStagger>
                <div>Stagger content</div>
            </FadeInStagger>
        );

        expect(screen.getByText('Stagger content')).toBeInTheDocument();
    });

    it('applies default animation classes with normal stagger timing', () => {
        render(
            <FadeInStagger>
                <div data-testid="content">Test content</div>
            </FadeInStagger>
        );

        const container = screen.getByTestId('content').parentElement;
        expect(container).toHaveClass('animate-in');
        expect(container).toHaveClass('fade-in');
        expect(container).toHaveClass('duration-700');
        expect(container).toHaveClass('stagger-children-500');
        expect(container).not.toHaveClass('stagger-children-300');
    });

    it('applies faster stagger timing when faster prop is true', () => {
        render(
            <FadeInStagger faster={true}>
                <div data-testid="content">Test content</div>
            </FadeInStagger>
        );

        const container = screen.getByTestId('content').parentElement;
        expect(container).toHaveClass('animate-in');
        expect(container).toHaveClass('fade-in');
        expect(container).toHaveClass('duration-700');
        expect(container).toHaveClass('stagger-children-300');
        expect(container).not.toHaveClass('stagger-children-500');
    });

    it('applies normal stagger timing when faster prop is false', () => {
        render(
            <FadeInStagger faster={false}>
                <div data-testid="content">Test content</div>
            </FadeInStagger>
        );

        const container = screen.getByTestId('content').parentElement;
        expect(container).toHaveClass('stagger-children-500');
        expect(container).not.toHaveClass('stagger-children-300');
    });

    it('applies custom className along with default classes', () => {
        render(
            <FadeInStagger className="custom-stagger">
                <div data-testid="content">Test content</div>
            </FadeInStagger>
        );

        const container = screen.getByTestId('content').parentElement;
        expect(container).toHaveClass('animate-in');
        expect(container).toHaveClass('fade-in');
        expect(container).toHaveClass('duration-700');
        expect(container).toHaveClass('stagger-children-500');
        expect(container).toHaveClass('custom-stagger');
    });

    it('combines faster prop with custom className', () => {
        render(
            <FadeInStagger faster={true} className="custom-fast">
                <div data-testid="content">Test content</div>
            </FadeInStagger>
        );

        const container = screen.getByTestId('content').parentElement;
        expect(container).toHaveClass('animate-in');
        expect(container).toHaveClass('fade-in');
        expect(container).toHaveClass('duration-700');
        expect(container).toHaveClass('stagger-children-300');
        expect(container).toHaveClass('custom-fast');
    });

    it('renders multiple children correctly', () => {
        render(
            <FadeInStagger>
                <div>Child 1</div>
                <div>Child 2</div>
                <div>Child 3</div>
            </FadeInStagger>
        );

        expect(screen.getByText('Child 1')).toBeInTheDocument();
        expect(screen.getByText('Child 2')).toBeInTheDocument();
        expect(screen.getByText('Child 3')).toBeInTheDocument();
    });

    it('renders without custom className', () => {
        render(
            <FadeInStagger>
                <div data-testid="content">Test content</div>
            </FadeInStagger>
        );

        const container = screen.getByTestId('content').parentElement;
        expect(container).toHaveClass('animate-in', 'fade-in', 'duration-700', 'stagger-children-500');
    });

    it('has correct HTML structure', () => {
        const {container} = render(
            <FadeInStagger>
                <div>Content</div>
            </FadeInStagger>
        );

        const staggerDiv = container.firstChild;
        expect(staggerDiv).toBeInstanceOf(HTMLDivElement);
        expect(staggerDiv?.textContent).toBe('Content');
    });

    it('preserves event handlers on children', () => {
        const handleClick1 = jest.fn();
        const handleClick2 = jest.fn();

        render(
            <FadeInStagger>
                <button onClick={handleClick1}>Button 1</button>
                <button onClick={handleClick2}>Button 2</button>
            </FadeInStagger>
        );

        const button1 = screen.getByRole('button', {name: 'Button 1'});
        const button2 = screen.getByRole('button', {name: 'Button 2'});

        button1.click();
        button2.click();

        expect(handleClick1).toHaveBeenCalledTimes(1);
        expect(handleClick2).toHaveBeenCalledTimes(1);
    });

    it('handles React fragments as children', () => {
        render(
            <FadeInStagger>
                <>
                    <div>Fragment child 1</div>
                    <div>Fragment child 2</div>
                </>
            </FadeInStagger>
        );

        expect(screen.getByText('Fragment child 1')).toBeInTheDocument();
        expect(screen.getByText('Fragment child 2')).toBeInTheDocument();
    });

    it('can be nested inside FadeIn', () => {
        render(
            <FadeIn className="outer-fade">
                <FadeInStagger className="inner-stagger">
                    <div data-testid="nested-content">Nested content</div>
                </FadeInStagger>
            </FadeIn>
        );

        const nestedContent = screen.getByTestId('nested-content');
        const innerStagger = nestedContent.parentElement;
        const outerFade = innerStagger?.parentElement;

        expect(innerStagger).toHaveClass('inner-stagger', 'stagger-children-500');
        expect(outerFade).toHaveClass('outer-fade', 'fade-in');
    });
});