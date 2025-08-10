import {render, screen, fireEvent, waitFor} from '@testing-library/react';
import {BrowserRouter} from 'react-router-dom';
import {AuthProvider, useAuth} from './AuthContext';
import {authAPI} from '../services/api';

// Mock the API
jest.mock('../services/api', () => ({
    authAPI: {
        login: jest.fn(),
        signup: jest.fn(),
        getProfile: jest.fn(),
    },
}));

const mockedAuthAPI = authAPI as jest.Mocked<typeof authAPI>;

// Test component to interact with AuthContext
const TestComponent = () => {
    const {isAuthenticated, isLoading, error, login, logout, user} = useAuth();

    return (
        <div>
            <div data-testid="auth-status">
                {isAuthenticated ? 'authenticated' : 'not-authenticated'}
            </div>
            <div data-testid="loading-status">
                {isLoading ? 'loading' : 'not-loading'}
            </div>
            <div data-testid="error-status">
                {error || 'no-error'}
            </div>
            <div data-testid="user-info">
                {user ? `${user.email}-${user.displayName}` : 'no-user'}
            </div>
            <button onClick={async () => {
                try {
                    await login('test@example.com', 'password');
                } catch (error) {
                    // Error is already handled in context
                }
            }}>
                Login
            </button>
            <button onClick={logout}>
                Logout
            </button>
        </div>
    );
};

const renderWithAuthProvider = () => {
    return render(
        <BrowserRouter>
            <AuthProvider>
                <TestComponent/>
            </AuthProvider>
        </BrowserRouter>
    );
};

// Mock localStorage
const localStorageMock = (() => {
    let store: Record<string, string> = {};

    return {
        getItem: jest.fn((key: string) => store[key] || null),
        setItem: jest.fn((key: string, value: string) => {
            store[key] = value;
        }),
        removeItem: jest.fn((key: string) => {
            delete store[key];
        }),
        clear: jest.fn(() => {
            store = {};
        }),
    };
})();

Object.defineProperty(global, 'localStorage', {
    value: localStorageMock,
});

describe('AuthContext', () => {
    beforeEach(() => {
        jest.clearAllMocks();
        localStorageMock.clear();
    });

    it('provides initial auth state', async () => {
        renderWithAuthProvider();

        // Wait for initialization to complete
        await waitFor(() => {
            expect(screen.getByTestId('loading-status')).toHaveTextContent('not-loading');
        });

        expect(screen.getByTestId('auth-status')).toHaveTextContent('not-authenticated');
        expect(screen.getByTestId('error-status')).toHaveTextContent('no-error');
        expect(screen.getByTestId('user-info')).toHaveTextContent('no-user');
    });

    it('handles successful login', async () => {
        const mockUser = {
            id: 1,
            email: 'test@example.com',
            displayName: 'testuser',
            emailVerified: false,
        };

        mockedAuthAPI.login.mockResolvedValue(mockUser);

        renderWithAuthProvider();

        const loginButton = screen.getByText('Login');
        fireEvent.click(loginButton);

        // Should show loading state
        await waitFor(() => {
            expect(screen.getByTestId('loading-status')).toHaveTextContent('loading');
        });

        // Should complete login
        await waitFor(() => {
            expect(screen.getByTestId('auth-status')).toHaveTextContent('authenticated');
            expect(screen.getByTestId('loading-status')).toHaveTextContent('not-loading');
            expect(screen.getByTestId('user-info')).toHaveTextContent('test@example.com-testuser');
            expect(screen.getByTestId('error-status')).toHaveTextContent('no-error');
        });

        expect(mockedAuthAPI.login).toHaveBeenCalledWith({
            email: 'test@example.com',
            password: 'password',
        });
    });

    it('handles login failure', async () => {
        mockedAuthAPI.login.mockRejectedValue(new Error('Invalid credentials'));

        renderWithAuthProvider();

        // Wait for initialization to complete first
        await waitFor(() => {
            expect(screen.getByTestId('loading-status')).toHaveTextContent('not-loading');
        });

        const loginButton = screen.getByText('Login');
        fireEvent.click(loginButton);

        await waitFor(() => {
            expect(screen.getByTestId('auth-status')).toHaveTextContent('not-authenticated');
            expect(screen.getByTestId('loading-status')).toHaveTextContent('not-loading');
            expect(screen.getByTestId('error-status')).toHaveTextContent('Invalid credentials');
            expect(screen.getByTestId('user-info')).toHaveTextContent('no-user');
        });
    });

    it('handles logout', async () => {
        const mockUser = {
            id: 1,
            email: 'test@example.com',
            displayName: 'testuser',
            emailVerified: false,
        };

        mockedAuthAPI.login.mockResolvedValue(mockUser);

        renderWithAuthProvider();

        // First login
        const loginButton = screen.getByText('Login');
        fireEvent.click(loginButton);

        await waitFor(() => {
            expect(screen.getByTestId('auth-status')).toHaveTextContent('authenticated');
        });

        // Then logout
        const logoutButton = screen.getByText('Logout');
        fireEvent.click(logoutButton);

        await waitFor(() => {
            expect(screen.getByTestId('auth-status')).toHaveTextContent('not-authenticated');
            expect(screen.getByTestId('user-info')).toHaveTextContent('no-user');
            expect(screen.getByTestId('error-status')).toHaveTextContent('no-error');
        });
    });

    it('persists user session in localStorage', async () => {
        const mockUser = {
            id: 1,
            email: 'test@example.com',
            displayName: 'testuser',
            emailVerified: false,
        };

        mockedAuthAPI.login.mockResolvedValue(mockUser);

        renderWithAuthProvider();

        const loginButton = screen.getByText('Login');
        fireEvent.click(loginButton);

        await waitFor(() => {
            expect(screen.getByTestId('auth-status')).toHaveTextContent('authenticated');
        });

        // Check if user is stored in localStorage
        const storedUser = localStorage.getItem('user');
        expect(storedUser).toBe(JSON.stringify(mockUser));
    });

    it('restores session from localStorage on initialization', () => {
        const mockUser = {
            id: 1,
            email: 'test@example.com',
            displayName: 'testuser',
            emailVerified: false,
        };

        // Pre-populate localStorage
        localStorage.setItem('user', JSON.stringify(mockUser));

        renderWithAuthProvider();

        expect(screen.getByTestId('auth-status')).toHaveTextContent('authenticated');
        expect(screen.getByTestId('user-info')).toHaveTextContent('test@example.com-testuser');
    });

    it('clears localStorage on logout', async () => {
        const mockUser = {
            id: 1,
            email: 'test@example.com',
            displayName: 'testuser',
            emailVerified: false,
        };

        localStorage.setItem('user', JSON.stringify(mockUser));

        renderWithAuthProvider();

        // Should start authenticated due to localStorage
        expect(screen.getByTestId('auth-status')).toHaveTextContent('authenticated');

        // Logout
        const logoutButton = screen.getByText('Logout');
        fireEvent.click(logoutButton);

        await waitFor(() => {
            expect(screen.getByTestId('auth-status')).toHaveTextContent('not-authenticated');
        });

        // Check localStorage is cleared
        expect(localStorage.getItem('user')).toBeNull();
    });
});