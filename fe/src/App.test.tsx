import {render, screen} from '@testing-library/react';
import App from './App';

// Mock the AuthContext since it makes API calls
jest.mock('./contexts/AuthContext', () => ({
    AuthProvider: ({children}: { children: React.ReactNode }) => <div>{children}</div>,
    useAuth: () => ({
        isAuthenticated: false,
        isLoading: false,
        error: null,
        login: jest.fn(),
        logout: jest.fn(),
    }),
}));

describe('App', () => {
    it('renders StrikePad welcome page', () => {
        render(<App/>);
        const titleElements = screen.getAllByText(/StrikePad/i);
        expect(titleElements.length).toBeGreaterThan(0);
    });

    it('renders login and signup buttons', () => {
        render(<App/>);
        const loginButtons = screen.getAllByText(/ログイン/i);
        const signupButtons = screen.getAllByText(/サインアップ/i);
        expect(loginButtons.length).toBeGreaterThan(0);
        expect(signupButtons.length).toBeGreaterThan(0);
    });

    it('renders welcome message', () => {
        render(<App/>);
        const welcomeMessage = screen.getByText(/Made with ❤️ for better productivity/i);
        expect(welcomeMessage).toBeInTheDocument();
    });
});