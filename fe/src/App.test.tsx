import {render, screen, waitFor} from '@testing-library/react';
import {MemoryRouter, Routes, Route} from 'react-router-dom';
import {useAuth} from './contexts/AuthContext';

// Mock the AuthContext
const mockUseAuth = jest.fn();

jest.mock('./contexts/AuthContext', () => ({
    AuthProvider: ({children}: { children: React.ReactNode }) => <div>{children}</div>,
    useAuth: () => mockUseAuth(),
}));

// Mock components for testing
const MockLandingPage = () => <div data-testid="landing-page">Landing Page</div>;
const MockLoginForm = () => <div data-testid="login-form">Login Form</div>;
const MockSignupForm = () => <div data-testid="signup-form">Signup Form</div>;
const MockDashboard = () => <div data-testid="dashboard">Dashboard</div>;

// Mock the actual components
jest.mock('./pages/LandingPage', () => () => <div data-testid="landing-page">Landing Page</div>);
jest.mock('./components/LoginForm', () => () => <div data-testid="login-form">Login Form</div>);
jest.mock('./components/SignupForm', () => () => <div data-testid="signup-form">Signup Form</div>);
jest.mock('./components/Dashboard', () => () => <div data-testid="dashboard">Dashboard</div>);

// Create route components similar to App.tsx
const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({children}) => {
    const {isAuthenticated, isLoading} = useAuth();

    if (isLoading) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"
                     data-testid="loading-spinner"></div>
            </div>
        );
    }

    return isAuthenticated ? <>{children}</> : <MockLoginForm/>;
};

const PublicRoute: React.FC<{ children: React.ReactNode }> = ({children}) => {
    const {isAuthenticated, isLoading} = useAuth();

    if (isLoading) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"
                     data-testid="loading-spinner"></div>
            </div>
        );
    }

    return isAuthenticated ? <MockDashboard/> : <>{children}</>;
};

// Test App component
const TestApp = () => {
    return (
        <Routes>
            <Route path="/" element={<MockLandingPage/>}/>
            <Route
                path="/login"
                element={
                    <PublicRoute>
                        <MockLoginForm/>
                    </PublicRoute>
                }
            />
            <Route
                path="/signup"
                element={
                    <PublicRoute>
                        <MockSignupForm/>
                    </PublicRoute>
                }
            />
            <Route
                path="/dashboard"
                element={
                    <ProtectedRoute>
                        <MockDashboard/>
                    </ProtectedRoute>
                }
            />
        </Routes>
    );
};

const renderWithRoute = (initialRoute: string) => {
    return render(
        <MemoryRouter initialEntries={[initialRoute]}>
            <TestApp/>
        </MemoryRouter>
    );
};

describe('App Routing Logic', () => {
    beforeEach(() => {
        jest.clearAllMocks();
    });

    describe('Landing Page Route', () => {
        it('renders landing page on root path', () => {
            mockUseAuth.mockReturnValue({
                isAuthenticated: false,
                isLoading: false,
                error: null,
                user: null,
                login: jest.fn(),
                logout: jest.fn(),
                signup: jest.fn(),
                googleSignup: jest.fn(),
                googleLogin: jest.fn(),
            });

            renderWithRoute('/');
            expect(screen.getByTestId('landing-page')).toBeInTheDocument();
        });
    });

    describe('PublicRoute - Not Authenticated', () => {
        beforeEach(() => {
            mockUseAuth.mockReturnValue({
                isAuthenticated: false,
                isLoading: false,
                error: null,
                user: null,
                login: jest.fn(),
                logout: jest.fn(),
                signup: jest.fn(),
                googleSignup: jest.fn(),
                googleLogin: jest.fn(),
            });
        });

        it('renders login form when accessing /login', () => {
            renderWithRoute('/login');
            expect(screen.getByTestId('login-form')).toBeInTheDocument();
        });

        it('renders signup form when accessing /signup', () => {
            renderWithRoute('/signup');
            expect(screen.getByTestId('signup-form')).toBeInTheDocument();
        });
    });

    describe('PublicRoute - Authenticated (should redirect)', () => {
        beforeEach(() => {
            mockUseAuth.mockReturnValue({
                isAuthenticated: true,
                isLoading: false,
                error: null,
                user: {id: 1, email: 'test@example.com', display_name: 'Test User', email_verified: true},
                login: jest.fn(),
                logout: jest.fn(),
                signup: jest.fn(),
                googleSignup: jest.fn(),
                googleLogin: jest.fn(),
            });
        });

        it('redirects to dashboard when authenticated user accesses /login', () => {
            renderWithRoute('/login');
            expect(screen.getByTestId('dashboard')).toBeInTheDocument();
        });

        it('redirects to dashboard when authenticated user accesses /signup', () => {
            renderWithRoute('/signup');
            expect(screen.getByTestId('dashboard')).toBeInTheDocument();
        });
    });

    describe('ProtectedRoute - Authenticated', () => {
        beforeEach(() => {
            mockUseAuth.mockReturnValue({
                isAuthenticated: true,
                isLoading: false,
                error: null,
                user: {id: 1, email: 'test@example.com', display_name: 'Test User', email_verified: true},
                login: jest.fn(),
                logout: jest.fn(),
                signup: jest.fn(),
                googleSignup: jest.fn(),
                googleLogin: jest.fn(),
            });
        });

        it('renders dashboard when authenticated user accesses /dashboard', () => {
            renderWithRoute('/dashboard');
            expect(screen.getByTestId('dashboard')).toBeInTheDocument();
        });
    });

    describe('ProtectedRoute - Not Authenticated (should redirect)', () => {
        beforeEach(() => {
            mockUseAuth.mockReturnValue({
                isAuthenticated: false,
                isLoading: false,
                error: null,
                user: null,
                login: jest.fn(),
                logout: jest.fn(),
                signup: jest.fn(),
                googleSignup: jest.fn(),
                googleLogin: jest.fn(),
            });
        });

        it('redirects to login when unauthenticated user accesses /dashboard', () => {
            renderWithRoute('/dashboard');
            expect(screen.getByTestId('login-form')).toBeInTheDocument();
        });
    });

    describe('Loading States', () => {
        it('shows loading spinner for ProtectedRoute while loading', () => {
            mockUseAuth.mockReturnValue({
                isAuthenticated: false,
                isLoading: true,
                error: null,
                user: null,
                login: jest.fn(),
                logout: jest.fn(),
                signup: jest.fn(),
                googleSignup: jest.fn(),
                googleLogin: jest.fn(),
            });

            renderWithRoute('/dashboard');
            expect(screen.getByTestId('loading-spinner')).toBeInTheDocument();
        });

        it('shows loading spinner for PublicRoute while loading', () => {
            mockUseAuth.mockReturnValue({
                isAuthenticated: false,
                isLoading: true,
                error: null,
                user: null,
                login: jest.fn(),
                logout: jest.fn(),
                signup: jest.fn(),
                googleSignup: jest.fn(),
                googleLogin: jest.fn(),
            });

            renderWithRoute('/login');
            expect(screen.getByTestId('loading-spinner')).toBeInTheDocument();
        });
    });

    describe('Route Component Logic', () => {
        it('handles authentication state changes in ProtectedRoute', () => {
            // Test not authenticated state
            mockUseAuth.mockReturnValue({
                isAuthenticated: false,
                isLoading: false,
                error: null,
                user: null,
                login: jest.fn(),
                logout: jest.fn(),
                signup: jest.fn(),
                googleSignup: jest.fn(),
                googleLogin: jest.fn(),
            });

            const {rerender} = renderWithRoute('/dashboard');
            expect(screen.getByTestId('login-form')).toBeInTheDocument();

            // Change to authenticated state
            mockUseAuth.mockReturnValue({
                isAuthenticated: true,
                isLoading: false,
                error: null,
                user: {id: 1, email: 'test@example.com', display_name: 'Test User', email_verified: true},
                login: jest.fn(),
                logout: jest.fn(),
                signup: jest.fn(),
                googleSignup: jest.fn(),
                googleLogin: jest.fn(),
            });

            rerender(
                <MemoryRouter initialEntries={['/dashboard']}>
                    <TestApp/>
                </MemoryRouter>
            );

            expect(screen.getByTestId('dashboard')).toBeInTheDocument();
        });

        it('handles authentication state changes in PublicRoute', () => {
            // Test not authenticated state
            mockUseAuth.mockReturnValue({
                isAuthenticated: false,
                isLoading: false,
                error: null,
                user: null,
                login: jest.fn(),
                logout: jest.fn(),
                signup: jest.fn(),
                googleSignup: jest.fn(),
                googleLogin: jest.fn(),
            });

            const {rerender} = renderWithRoute('/login');
            expect(screen.getByTestId('login-form')).toBeInTheDocument();

            // Change to authenticated state
            mockUseAuth.mockReturnValue({
                isAuthenticated: true,
                isLoading: false,
                error: null,
                user: {id: 1, email: 'test@example.com', display_name: 'Test User', email_verified: true},
                login: jest.fn(),
                logout: jest.fn(),
                signup: jest.fn(),
                googleSignup: jest.fn(),
                googleLogin: jest.fn(),
            });

            rerender(
                <MemoryRouter initialEntries={['/login']}>
                    <TestApp/>
                </MemoryRouter>
            );

            expect(screen.getByTestId('dashboard')).toBeInTheDocument();
        });
    });
});