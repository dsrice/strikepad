import React from 'react';
import {render, screen, fireEvent} from '@testing-library/react';
import {BrowserRouter} from 'react-router-dom';
import {Layout} from './Layout';

// Mock the LinkButton component to simplify testing
jest.mock('./ui/Button', () => ({
    LinkButton: ({children, href, variant, size, ...props}: any) => (
        <a href={href} data-variant={variant} data-size={size} {...props}>
            {children}
        </a>
    ),
}));

// Mock the Container component
jest.mock('./ui/Container', () => ({
    Container: ({children, className}: any) => (
        <div className={className} data-testid="container">
            {children}
        </div>
    ),
}));

const renderWithRouter = (ui: React.ReactElement) => {
    return render(
        <BrowserRouter>
            {ui}
        </BrowserRouter>
    );
};

describe('Layout', () => {
    const TestContent = () => <div data-testid="test-content">Test Content</div>;

    describe('Basic Rendering', () => {
        it('renders children content', () => {
            renderWithRouter(
                <Layout>
                    <TestContent/>
                </Layout>
            );
            expect(screen.getByTestId('test-content')).toBeInTheDocument();
        });

        it('renders logo in header', () => {
            renderWithRouter(
                <Layout>
                    <TestContent/>
                </Layout>
            );
            const logos = screen.getAllByText('StrikePad');
            expect(logos.length).toBeGreaterThan(0);
        });

        it('renders navigation buttons', () => {
            renderWithRouter(
                <Layout>
                    <TestContent/>
                </Layout>
            );
            expect(screen.getByRole('link', {name: /ログイン/i})).toBeInTheDocument();
            expect(screen.getByRole('link', {name: /サインアップ/i})).toBeInTheDocument();
        });

        it('renders footer content', () => {
            renderWithRouter(
                <Layout>
                    <TestContent/>
                </Layout>
            );
            expect(screen.getByText(/© \d{4} StrikePad. All rights reserved./)).toBeInTheDocument();
            expect(screen.getByText('Made with ❤️ for better productivity')).toBeInTheDocument();
            expect(screen.getByText('Version 1.0.0')).toBeInTheDocument();
        });
    });

    describe('Header Navigation', () => {
        it('renders menu toggle button', () => {
            renderWithRouter(
                <Layout>
                    <TestContent/>
                </Layout>
            );
            const menuButton = screen.getByRole('button', {name: /toggle navigation/i});
            expect(menuButton).toBeInTheDocument();
        });

        it('opens mobile menu when toggle button is clicked', () => {
            renderWithRouter(
                <Layout>
                    <TestContent/>
                </Layout>
            );

            // Initially, navigation menu should not be visible
            expect(screen.queryByText('機能')).not.toBeInTheDocument();
            expect(screen.queryByText('料金')).not.toBeInTheDocument();

            // Click the menu toggle button
            const menuButton = screen.getByRole('button', {name: /toggle navigation/i});
            fireEvent.click(menuButton);

            // Now navigation menu should be visible
            expect(screen.getByText('機能')).toBeInTheDocument();
            expect(screen.getByText('料金')).toBeInTheDocument();
        });

        it('closes mobile menu when close button is clicked', () => {
            renderWithRouter(
                <Layout>
                    <TestContent/>
                </Layout>
            );

            // Open the menu
            const menuButton = screen.getByRole('button', {name: /toggle navigation/i});
            fireEvent.click(menuButton);

            // Verify menu is open
            expect(screen.getByText('機能')).toBeInTheDocument();

            // Click the close button
            const closeButton = screen.getByRole('button', {name: /close navigation/i});
            fireEvent.click(closeButton);

            // Menu should be closed
            expect(screen.queryByText('機能')).not.toBeInTheDocument();
        });

        it('closes mobile menu when navigation link is clicked', () => {
            renderWithRouter(
                <Layout>
                    <TestContent/>
                </Layout>
            );

            // Open the menu
            const menuButton = screen.getByRole('button', {name: /toggle navigation/i});
            fireEvent.click(menuButton);

            // Verify menu is open
            expect(screen.getByText('機能')).toBeInTheDocument();

            // Click a navigation link
            const featuresLink = screen.getByRole('link', {name: '機能'});
            fireEvent.click(featuresLink);

            // Menu should be closed
            expect(screen.queryByText('機能')).not.toBeInTheDocument();
        });

        it('shows correct icon based on menu state', () => {
            renderWithRouter(
                <Layout>
                    <TestContent/>
                </Layout>
            );

            const menuButton = screen.getByRole('button', {name: /toggle navigation/i});

            // Initially should show menu icon (not X icon)
            expect(menuButton.querySelector('svg')).toBeInTheDocument();

            // Click to open menu
            fireEvent.click(menuButton);

            // Should show X icon when menu is open
            const closeButton = screen.getByRole('button', {name: /close navigation/i});
            expect(closeButton.querySelector('svg')).toBeInTheDocument();
        });
    });

    describe('Footer Navigation', () => {
        it('renders all footer links', () => {
            renderWithRouter(
                <Layout>
                    <TestContent/>
                </Layout>
            );

            expect(screen.getByRole('link', {name: /プライバシーポリシー/i})).toBeInTheDocument();
            expect(screen.getByRole('link', {name: /利用規約/i})).toBeInTheDocument();
            expect(screen.getByRole('link', {name: /お問い合わせ/i})).toBeInTheDocument();
            expect(screen.getByRole('link', {name: /ヘルプ/i})).toBeInTheDocument();
        });

        it('renders footer logo as link to home', () => {
            renderWithRouter(
                <Layout>
                    <TestContent/>
                </Layout>
            );

            const footerLogo = screen.getAllByRole('link', {name: /home/i});
            expect(footerLogo.length).toBeGreaterThan(0);
        });

        it('displays current year in copyright', () => {
            renderWithRouter(
                <Layout>
                    <TestContent/>
                </Layout>
            );

            const currentYear = new Date().getFullYear();
            expect(screen.getByText(new RegExp(`© ${currentYear} StrikePad`))).toBeInTheDocument();
        });
    });

    describe('Accessibility', () => {
        it('has proper aria-labels for navigation elements', () => {
            renderWithRouter(
                <Layout>
                    <TestContent/>
                </Layout>
            );

            expect(screen.getByRole('button', {name: /toggle navigation/i})).toHaveAttribute('aria-label', 'Toggle navigation');
            expect(screen.getAllByRole('link', {name: /home/i})[0]).toHaveAttribute('aria-label', 'Home');
        });

        it('sets aria-expanded correctly on menu button', () => {
            renderWithRouter(
                <Layout>
                    <TestContent/>
                </Layout>
            );

            const menuButton = screen.getByRole('button', {name: /toggle navigation/i});

            // Initially closed
            expect(menuButton).toHaveAttribute('aria-expanded', 'false');

            // Open menu
            fireEvent.click(menuButton);
            expect(menuButton).toHaveAttribute('aria-expanded', 'true');
        });

        it('provides proper svg accessibility attributes', () => {
            renderWithRouter(
                <Layout>
                    <TestContent/>
                </Layout>
            );

            const svgElements = document.querySelectorAll('svg');
            svgElements.forEach(svg => {
                expect(svg).toHaveAttribute('aria-hidden', 'true');
            });
        });
    });

    describe('Layout Structure', () => {
        it('has proper layout structure with header, main, and footer', () => {
            renderWithRouter(
                <Layout>
                    <TestContent/>
                </Layout>
            );

            expect(screen.getByRole('banner')).toBeInTheDocument(); // header
            expect(screen.getByRole('main')).toBeInTheDocument(); // main
            expect(screen.getByRole('contentinfo')).toBeInTheDocument(); // footer
        });

        it('applies correct CSS classes for styling', () => {
            renderWithRouter(
                <Layout>
                    <TestContent/>
                </Layout>
            );

            const mainElement = screen.getByRole('main');
            expect(mainElement).toHaveClass('w-full', 'flex-1');
        });
    });

    describe('Mobile Menu Behavior', () => {
        it('menu toggle state is managed correctly', () => {
            renderWithRouter(
                <Layout>
                    <TestContent/>
                </Layout>
            );

            const menuButton = screen.getByRole('button', {name: /toggle navigation/i});

            // Click to open
            fireEvent.click(menuButton);
            expect(screen.getByText('機能')).toBeInTheDocument();

            // Click again to close
            fireEvent.click(menuButton);
            expect(screen.queryByText('機能')).not.toBeInTheDocument();
        });

        it('clicking pricing link closes menu', () => {
            renderWithRouter(
                <Layout>
                    <TestContent/>
                </Layout>
            );

            // Open menu
            const menuButton = screen.getByRole('button', {name: /toggle navigation/i});
            fireEvent.click(menuButton);

            // Click pricing link
            const pricingLink = screen.getByRole('link', {name: '料金'});
            fireEvent.click(pricingLink);

            // Menu should close
            expect(screen.queryByText('機能')).not.toBeInTheDocument();
        });
    });
});