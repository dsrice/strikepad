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
                {user ? `${user.email}-${user.username}` : 'no-user'}
            </div>
            <button onClick={() => login('test@example.com', 'password')}>
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

describe('AuthContext', () => {
    beforeEach(() => {
        jest.clearAllMocks();
        localStorage.clear();
    });

    it('provides initial auth state', () => {
        renderWithAuthProvider();

        expect(screen.getByTestId('auth-status')).toHaveTextContent('not-authenticated');
        expect(screen.getByTestId('loading-status')).toHaveTextContent('not-loading');
        expect(screen.getByTestId('error-status')).toHaveTextContent('no-error');
        expect(screen.getByTestId('user-info')).toHaveTextContent('no-user');
    });

    it('handles successful login', async () => {
        const mockUser = {
            id: 1,
            email: 'test@example.com',
            username: 'testuser',
            createdAt: '2023-01-01T00:00:00Z',
            updatedAt: '2023-01-01T00:00:00Z',
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
            username: 'testuser',
            createdAt: '2023-01-01T00:00:00Z',
            updatedAt: '2023-01-01T00:00:00Z',
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
            username: 'testuser',
            createdAt: '2023-01-01T00:00:00Z',
            updatedAt: '2023-01-01T00:00:00Z',
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
            username: 'testuser',
            createdAt: '2023-01-01T00:00:00Z',
            updatedAt: '2023-01-01T00:00:00Z',
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
            username: 'testuser',
            createdAt: '2023-01-01T00:00:00Z',
            updatedAt: '2023-01-01T00:00:00Z',
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