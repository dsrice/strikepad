import {render, screen, fireEvent, waitFor} from '@testing-library/react';
import {BrowserRouter} from 'react-router-dom';
import {AuthProvider, useAuth} from './AuthContext';
import {authAPI} from '../services/api';

// Mock the API
jest.mock('../services/api', () => ({
    authAPI: {
        login: jest.fn(),
        signup: jest.fn(),
        googleSignup: jest.fn(),
        googleLogin: jest.fn(),
        logout: jest.fn(),
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
                {user ? `${user.email}-${user.display_name}` : 'no-user'}
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
        const mockLoginResponse = {
            access_token: 'test-access-token',
            refresh_token: 'test-refresh-token',
            expires_at: '2024-12-31T23:59:59Z',
        };

        const mockUser = {
            id: 1,
            email: 'test@example.com',
            display_name: 'testuser',
            email_verified: false,
        };

        mockedAuthAPI.login.mockResolvedValue(mockLoginResponse);
        mockedAuthAPI.getProfile.mockResolvedValue(mockUser);

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
        expect(mockedAuthAPI.getProfile).toHaveBeenCalled();
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
        const mockLoginResponse = {
            access_token: 'test-access-token',
            refresh_token: 'test-refresh-token',
            expires_at: '2024-12-31T23:59:59Z',
        };

        const mockUser = {
            id: 1,
            email: 'test@example.com',
            display_name: 'testuser',
            email_verified: false,
        };

        mockedAuthAPI.login.mockResolvedValue(mockLoginResponse);
        mockedAuthAPI.getProfile.mockResolvedValue(mockUser);

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

    it('persists tokens in localStorage', async () => {
        const mockLoginResponse = {
            access_token: 'test-access-token',
            refresh_token: 'test-refresh-token',
            expires_at: '2024-12-31T23:59:59Z',
        };

        const mockUser = {
            id: 1,
            email: 'test@example.com',
            display_name: 'testuser',
            email_verified: false,
        };

        mockedAuthAPI.login.mockResolvedValue(mockLoginResponse);
        mockedAuthAPI.getProfile.mockResolvedValue(mockUser);

        renderWithAuthProvider();

        const loginButton = screen.getByText('Login');
        fireEvent.click(loginButton);

        await waitFor(() => {
            expect(screen.getByTestId('auth-status')).toHaveTextContent('authenticated');
        });

        // Check if tokens are stored in localStorage
        expect(localStorage.setItem).toHaveBeenCalledWith('access_token', 'test-access-token');
        expect(localStorage.setItem).toHaveBeenCalledWith('refresh_token', 'test-refresh-token');
    });

    it('restores session from localStorage on initialization', async () => {
        const mockUser = {
            id: 1,
            email: 'test@example.com',
            display_name: 'testuser',
            email_verified: false,
        };

        // Pre-populate localStorage with token
        localStorage.setItem('access_token', 'existing-token');
        mockedAuthAPI.getProfile.mockResolvedValue(mockUser);

        renderWithAuthProvider();

        await waitFor(() => {
            expect(screen.getByTestId('auth-status')).toHaveTextContent('authenticated');
            expect(screen.getByTestId('user-info')).toHaveTextContent('test@example.com-testuser');
        });

        expect(mockedAuthAPI.getProfile).toHaveBeenCalled();
    });

    it('clears localStorage on logout', async () => {
        const mockUser = {
            id: 1,
            email: 'test@example.com',
            display_name: 'testuser',
            email_verified: false,
        };

        localStorage.setItem('access_token', 'existing-token');
        localStorage.setItem('refresh_token', 'existing-refresh-token');
        mockedAuthAPI.getProfile.mockResolvedValue(mockUser);

        renderWithAuthProvider();

        // Should start authenticated due to localStorage
        await waitFor(() => {
            expect(screen.getByTestId('auth-status')).toHaveTextContent('authenticated');
        });

        // Logout
        const logoutButton = screen.getByText('Logout');
        fireEvent.click(logoutButton);

        await waitFor(() => {
            expect(screen.getByTestId('auth-status')).toHaveTextContent('not-authenticated');
        });

        // Check localStorage is cleared
        expect(localStorage.removeItem).toHaveBeenCalledWith('access_token');
        expect(localStorage.removeItem).toHaveBeenCalledWith('refresh_token');
    });

    it('handles successful signup', async () => {
        const mockSignupResponse = {
            id: 2,
            email: 'newuser@example.com',
            display_name: 'New User',
            email_verified: false,
            access_token: 'access_token_123',
            refresh_token: 'refresh_token_456',
            expires_at: '2024-12-31T23:59:59Z',
            created_at: '2023-01-01T00:00:00Z',
        };

        const mockUser = {
            id: 2,
            email: 'newuser@example.com',
            display_name: 'New User',
            email_verified: false,
        };

        mockedAuthAPI.signup.mockResolvedValue(mockSignupResponse);
        mockedAuthAPI.getProfile.mockResolvedValue(mockUser);

        const TestComponentWithSignup = () => {
            const {signup, isAuthenticated, user, isLoading, error} = useAuth();

            return (
                <div>
                    <div data-testid="signup-auth-status">
                        {isAuthenticated ? 'authenticated' : 'not-authenticated'}
                    </div>
                    <div data-testid="signup-loading-status">
                        {isLoading ? 'loading' : 'not-loading'}
                    </div>
                    <div data-testid="signup-error-status">
                        {error || 'no-error'}
                    </div>
                    <div data-testid="signup-user-info">
                        {user ? `${user.email}-${user.display_name}` : 'no-user'}
                    </div>
                    <button onClick={async () => {
                        try {
                            await signup('newuser@example.com', 'password123', 'New User');
                        } catch (error) {
                            // Error is already handled in context
                        }
                    }}>
                        Signup
                    </button>
                </div>
            );
        };

        render(
            <BrowserRouter>
                <AuthProvider>
                    <TestComponentWithSignup/>
                </AuthProvider>
            </BrowserRouter>
        );

        // Wait for initialization to complete first
        await waitFor(() => {
            expect(screen.getByTestId('signup-loading-status')).toHaveTextContent('not-loading');
        });

        const signupButton = screen.getByText('Signup');
        fireEvent.click(signupButton);

        // Should show loading state
        await waitFor(() => {
            expect(screen.getByTestId('signup-loading-status')).toHaveTextContent('loading');
        });

        // Should complete signup
        await waitFor(() => {
            expect(screen.getByTestId('signup-auth-status')).toHaveTextContent('authenticated');
            expect(screen.getByTestId('signup-loading-status')).toHaveTextContent('not-loading');
            expect(screen.getByTestId('signup-user-info')).toHaveTextContent('newuser@example.com-New User');
            expect(screen.getByTestId('signup-error-status')).toHaveTextContent('no-error');
        });

        expect(mockedAuthAPI.signup).toHaveBeenCalledWith({
            email: 'newuser@example.com',
            password: 'password123',
            display_name: 'New User',
        });
        expect(mockedAuthAPI.getProfile).toHaveBeenCalled();

        // Check tokens are stored
        expect(localStorage.setItem).toHaveBeenCalledWith('access_token', 'access_token_123');
        expect(localStorage.setItem).toHaveBeenCalledWith('refresh_token', 'refresh_token_456');
    });

    it('handles signup failure', async () => {
        mockedAuthAPI.signup.mockRejectedValue(new Error('Email already exists'));

        const TestComponentWithSignup = () => {
            const {signup, isAuthenticated, isLoading, error} = useAuth();

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
                    <button onClick={async () => {
                        try {
                            await signup('existing@example.com', 'password123', 'Test User');
                        } catch (error) {
                            // Error is already handled in context
                        }
                    }}>
                        Signup
                    </button>
                </div>
            );
        };

        render(
            <BrowserRouter>
                <AuthProvider>
                    <TestComponentWithSignup/>
                </AuthProvider>
            </BrowserRouter>
        );

        // Wait for initialization to complete first
        await waitFor(() => {
            expect(screen.getByTestId('loading-status')).toHaveTextContent('not-loading');
        });

        const signupButton = screen.getByText('Signup');
        fireEvent.click(signupButton);

        await waitFor(() => {
            expect(screen.getByTestId('auth-status')).toHaveTextContent('not-authenticated');
            expect(screen.getByTestId('loading-status')).toHaveTextContent('not-loading');
            expect(screen.getByTestId('error-status')).toHaveTextContent('Email already exists');
        });
    });

    it('handles Google signup', async () => {
        const mockGoogleSignupResponse = {
            id: 3,
            email: 'google@example.com',
            display_name: 'Google User',
            email_verified: true,
            access_token: 'google-access-token',
            refresh_token: 'google-refresh-token',
            expires_at: '2024-12-31T23:59:59Z',
            created_at: '2023-01-01T00:00:00Z',
        };

        const mockUser = {
            id: 3,
            email: 'google@example.com',
            display_name: 'Google User',
            email_verified: true,
        };

        mockedAuthAPI.googleSignup.mockResolvedValue(mockGoogleSignupResponse);
        mockedAuthAPI.getProfile.mockResolvedValue(mockUser);

        const TestComponentWithGoogleSignup = () => {
            const {googleSignup, isAuthenticated, user, isLoading, error} = useAuth();

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
                        {user ? `${user.email}-${user.display_name}` : 'no-user'}
                    </div>
                    <button onClick={async () => {
                        try {
                            await googleSignup('google-access-token');
                        } catch (error) {
                            // Error is already handled in context
                        }
                    }}>
                        Google Signup
                    </button>
                </div>
            );
        };

        render(
            <BrowserRouter>
                <AuthProvider>
                    <TestComponentWithGoogleSignup/>
                </AuthProvider>
            </BrowserRouter>
        );

        const googleSignupButton = screen.getByText('Google Signup');
        fireEvent.click(googleSignupButton);

        await waitFor(() => {
            expect(screen.getByTestId('auth-status')).toHaveTextContent('authenticated');
            expect(screen.getByTestId('loading-status')).toHaveTextContent('not-loading');
            expect(screen.getByTestId('user-info')).toHaveTextContent('google@example.com-Google User');
            expect(screen.getByTestId('error-status')).toHaveTextContent('no-error');
        });

        expect(mockedAuthAPI.googleSignup).toHaveBeenCalledWith({
            access_token: 'google-access-token',
        });
        expect(mockedAuthAPI.getProfile).toHaveBeenCalled();

        // Check tokens are stored
        expect(localStorage.setItem).toHaveBeenCalledWith('access_token', 'google-access-token');
        expect(localStorage.setItem).toHaveBeenCalledWith('refresh_token', 'google-refresh-token');
    });

    it('handles Google login', async () => {
        const mockGoogleLoginResponse = {
            access_token: 'google-access-token',
            refresh_token: 'google-refresh-token',
            expires_at: '2024-12-31T23:59:59Z',
        };

        const mockGoogleUser = {
            id: 4,
            email: 'googlelogin@example.com',
            display_name: 'Google Login User',
            email_verified: true,
        };

        mockedAuthAPI.googleLogin.mockResolvedValue(mockGoogleLoginResponse);
        mockedAuthAPI.getProfile.mockResolvedValue(mockGoogleUser);

        const TestComponentWithGoogleLogin = () => {
            const {googleLogin, isAuthenticated, user, isLoading, error} = useAuth();

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
                        {user ? `${user.email}-${user.display_name}` : 'no-user'}
                    </div>
                    <button onClick={async () => {
                        try {
                            await googleLogin('google-access-token');
                        } catch (error) {
                            // Error is already handled in context
                        }
                    }}>
                        Google Login
                    </button>
                </div>
            );
        };

        render(
            <BrowserRouter>
                <AuthProvider>
                    <TestComponentWithGoogleLogin/>
                </AuthProvider>
            </BrowserRouter>
        );

        const googleLoginButton = screen.getByText('Google Login');
        fireEvent.click(googleLoginButton);

        await waitFor(() => {
            expect(screen.getByTestId('auth-status')).toHaveTextContent('authenticated');
            expect(screen.getByTestId('loading-status')).toHaveTextContent('not-loading');
            expect(screen.getByTestId('user-info')).toHaveTextContent('googlelogin@example.com-Google Login User');
            expect(screen.getByTestId('error-status')).toHaveTextContent('no-error');
        });

        expect(mockedAuthAPI.googleLogin).toHaveBeenCalledWith({
            access_token: 'google-access-token',
        });
        expect(mockedAuthAPI.getProfile).toHaveBeenCalled();

        // Check tokens are stored
        expect(localStorage.setItem).toHaveBeenCalledWith('access_token', 'google-access-token');
        expect(localStorage.setItem).toHaveBeenCalledWith('refresh_token', 'google-refresh-token');
    });

    it('handles invalid access token on initialization', async () => {
        // Mock console.error to avoid noise in test output
        const consoleErrorSpy = jest.spyOn(console, 'error').mockImplementation(() => {
        });

        // Set invalid token data
        localStorage.setItem('access_token', 'invalid-token');
        mockedAuthAPI.getProfile.mockRejectedValue(new Error('Invalid token'));

        renderWithAuthProvider();

        // Should handle invalid token gracefully
        await waitFor(() => {
            expect(screen.getByTestId('auth-status')).toHaveTextContent('not-authenticated');
            expect(screen.getByTestId('user-info')).toHaveTextContent('no-user');
        });

        expect(mockedAuthAPI.getProfile).toHaveBeenCalled();

        consoleErrorSpy.mockRestore();
    });

    it('handles logout API failure gracefully', async () => {
        const consoleErrorSpy = jest.spyOn(console, 'error').mockImplementation(() => {
        });

        // First, log in a user
        const mockLoginResponse = {
            access_token: 'test-access-token',
            refresh_token: 'test-refresh-token',
            expires_at: '2024-12-31T23:59:59Z',
        };

        const mockUser = {
            id: 1,
            email: 'test@example.com',
            display_name: 'testuser',
            email_verified: false,
        };

        mockedAuthAPI.login.mockResolvedValue(mockLoginResponse);
        mockedAuthAPI.getProfile.mockResolvedValue(mockUser);
        mockedAuthAPI.logout.mockRejectedValue(new Error('Server error'));

        renderWithAuthProvider();

        // Login first
        const loginButton = screen.getByText('Login');
        fireEvent.click(loginButton);

        await waitFor(() => {
            expect(screen.getByTestId('auth-status')).toHaveTextContent('authenticated');
        });

        // Then logout - should work even if API fails
        const logoutButton = screen.getByText('Logout');
        fireEvent.click(logoutButton);

        await waitFor(() => {
            expect(screen.getByTestId('auth-status')).toHaveTextContent('not-authenticated');
            expect(screen.getByTestId('user-info')).toHaveTextContent('no-user');
            expect(screen.getByTestId('error-status')).toHaveTextContent('no-error');
        });

        // Should clear all localStorage items
        expect(localStorage.removeItem).toHaveBeenCalledWith('access_token');
        expect(localStorage.removeItem).toHaveBeenCalledWith('refresh_token');

        consoleErrorSpy.mockRestore();
    });

    it('throws error when useAuth is used outside AuthProvider', () => {
        const TestComponentOutsideProvider = () => {
            const auth = useAuth();
            return <div>{auth.user?.email}</div>;
        };

        // Suppress expected error message in console
        const consoleErrorSpy = jest.spyOn(console, 'error').mockImplementation(() => {
        });

        expect(() => {
            render(<TestComponentOutsideProvider/>);
        }).toThrow('useAuth must be used within an AuthProvider');

        consoleErrorSpy.mockRestore();
    });

    it('handles getProfile failure during initialization', async () => {
        // Set token in localStorage but make getProfile fail
        localStorage.setItem('access_token', 'invalid-token');
        mockedAuthAPI.getProfile.mockRejectedValue(new Error('Token expired'));

        renderWithAuthProvider();

        await waitFor(() => {
            expect(screen.getByTestId('auth-status')).toHaveTextContent('not-authenticated');
            expect(screen.getByTestId('user-info')).toHaveTextContent('no-user');
            expect(screen.getByTestId('loading-status')).toHaveTextContent('not-loading');
        });

        expect(mockedAuthAPI.getProfile).toHaveBeenCalled();
    });

    it('handles getProfile failure after successful login', async () => {
        const mockLoginResponse = {
            access_token: 'test-access-token',
            refresh_token: 'test-refresh-token',
            expires_at: '2024-12-31T23:59:59Z',
        };

        mockedAuthAPI.login.mockResolvedValue(mockLoginResponse);
        mockedAuthAPI.getProfile.mockRejectedValue(new Error('Profile fetch failed'));

        renderWithAuthProvider();

        const loginButton = screen.getByText('Login');
        fireEvent.click(loginButton);

        await waitFor(() => {
            expect(screen.getByTestId('auth-status')).toHaveTextContent('not-authenticated');
            expect(screen.getByTestId('error-status')).toHaveTextContent('Profile fetch failed');
            expect(screen.getByTestId('loading-status')).toHaveTextContent('not-loading');
        });

        expect(mockedAuthAPI.login).toHaveBeenCalled();
        expect(mockedAuthAPI.getProfile).toHaveBeenCalled();
    });

    it('handles getProfile failure after successful signup', async () => {
        const mockSignupResponse = {
            id: 2,
            email: 'newuser@example.com',
            display_name: 'New User',
            email_verified: false,
            access_token: 'access_token_123',
            refresh_token: 'refresh_token_456',
            expires_at: '2024-12-31T23:59:59Z',
            created_at: '2023-01-01T00:00:00Z',
        };

        mockedAuthAPI.signup.mockResolvedValue(mockSignupResponse);
        mockedAuthAPI.getProfile.mockRejectedValue(new Error('Profile fetch failed'));

        const TestComponentWithSignup = () => {
            const {signup, isAuthenticated, user, isLoading, error} = useAuth();

            return (
                <div>
                    <div data-testid="signup-auth-status">
                        {isAuthenticated ? 'authenticated' : 'not-authenticated'}
                    </div>
                    <div data-testid="signup-loading-status">
                        {isLoading ? 'loading' : 'not-loading'}
                    </div>
                    <div data-testid="signup-error-status">
                        {error || 'no-error'}
                    </div>
                    <div data-testid="signup-user-info">
                        {user ? `${user.email}-${user.display_name}` : 'no-user'}
                    </div>
                    <button onClick={async () => {
                        try {
                            await signup('newuser@example.com', 'password123', 'New User');
                        } catch (error) {
                            // Error is already handled in context
                        }
                    }}>
                        Signup
                    </button>
                </div>
            );
        };

        render(
            <BrowserRouter>
                <AuthProvider>
                    <TestComponentWithSignup/>
                </AuthProvider>
            </BrowserRouter>
        );

        const signupButton = screen.getByText('Signup');
        fireEvent.click(signupButton);

        await waitFor(() => {
            expect(screen.getByTestId('signup-auth-status')).toHaveTextContent('not-authenticated');
            expect(screen.getByTestId('signup-error-status')).toHaveTextContent('Profile fetch failed');
            expect(screen.getByTestId('signup-loading-status')).toHaveTextContent('not-loading');
        });

        expect(mockedAuthAPI.signup).toHaveBeenCalled();
        expect(mockedAuthAPI.getProfile).toHaveBeenCalled();
    });

    it('handles getProfile failure after Google signup', async () => {
        const mockGoogleSignupResponse = {
            id: 3,
            email: 'google@example.com',
            display_name: 'Google User',
            email_verified: true,
            access_token: 'google-access-token',
            refresh_token: 'google-refresh-token',
            expires_at: '2024-12-31T23:59:59Z',
            created_at: '2023-01-01T00:00:00Z',
        };

        mockedAuthAPI.googleSignup.mockResolvedValue(mockGoogleSignupResponse);
        mockedAuthAPI.getProfile.mockRejectedValue(new Error('Profile fetch failed'));

        const TestComponentWithGoogleSignup = () => {
            const {googleSignup, isAuthenticated, user, isLoading, error} = useAuth();

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
                        {user ? `${user.email}-${user.display_name}` : 'no-user'}
                    </div>
                    <button onClick={async () => {
                        try {
                            await googleSignup('google-access-token');
                        } catch (error) {
                            // Error is already handled in context
                        }
                    }}>
                        Google Signup
                    </button>
                </div>
            );
        };

        render(
            <BrowserRouter>
                <AuthProvider>
                    <TestComponentWithGoogleSignup/>
                </AuthProvider>
            </BrowserRouter>
        );

        const googleSignupButton = screen.getByText('Google Signup');
        fireEvent.click(googleSignupButton);

        await waitFor(() => {
            expect(screen.getByTestId('auth-status')).toHaveTextContent('not-authenticated');
            expect(screen.getByTestId('error-status')).toHaveTextContent('Profile fetch failed');
            expect(screen.getByTestId('loading-status')).toHaveTextContent('not-loading');
        });

        expect(mockedAuthAPI.googleSignup).toHaveBeenCalled();
        expect(mockedAuthAPI.getProfile).toHaveBeenCalled();
    });

    it('handles getProfile failure after Google login', async () => {
        const mockGoogleLoginResponse = {
            access_token: 'google-access-token',
            refresh_token: 'google-refresh-token',
            expires_at: '2024-12-31T23:59:59Z',
        };

        mockedAuthAPI.googleLogin.mockResolvedValue(mockGoogleLoginResponse);
        mockedAuthAPI.getProfile.mockRejectedValue(new Error('Profile fetch failed'));

        const TestComponentWithGoogleLogin = () => {
            const {googleLogin, isAuthenticated, user, isLoading, error} = useAuth();

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
                        {user ? `${user.email}-${user.display_name}` : 'no-user'}
                    </div>
                    <button onClick={async () => {
                        try {
                            await googleLogin('google-access-token');
                        } catch (error) {
                            // Error is already handled in context
                        }
                    }}>
                        Google Login
                    </button>
                </div>
            );
        };

        render(
            <BrowserRouter>
                <AuthProvider>
                    <TestComponentWithGoogleLogin/>
                </AuthProvider>
            </BrowserRouter>
        );

        const googleLoginButton = screen.getByText('Google Login');
        fireEvent.click(googleLoginButton);

        await waitFor(() => {
            expect(screen.getByTestId('auth-status')).toHaveTextContent('not-authenticated');
            expect(screen.getByTestId('error-status')).toHaveTextContent('Profile fetch failed');
            expect(screen.getByTestId('loading-status')).toHaveTextContent('not-loading');
        });

        expect(mockedAuthAPI.googleLogin).toHaveBeenCalled();
        expect(mockedAuthAPI.getProfile).toHaveBeenCalled();
    });
});