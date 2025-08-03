import {render, screen} from '@testing-library/react';
import {BrowserRouter} from 'react-router-dom';
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

const renderWithRouter = (component: React.ReactElement) => {
    return render(<BrowserRouter>{component}</BrowserRouter>);
};

describe('App', () => {
    it('renders StrikePad welcome page', () => {
        renderWithRouter(<App/>);
        const titleElement = screen.getByText(/StrikePad/i);
    expect(titleElement).toBeInTheDocument();
  });

    it('renders login and signup buttons', () => {
        renderWithRouter(<App/>);
        const loginButton = screen.getByText(/Login/i);
        const signupButton = screen.getByText(/Sign Up/i);
        expect(loginButton).toBeInTheDocument();
        expect(signupButton).toBeInTheDocument();
  });

    it('renders welcome message', () => {
        renderWithRouter(<App/>);
        const welcomeMessage = screen.getByText(/Welcome to StrikePad Application/i);
        expect(welcomeMessage).toBeInTheDocument();
  });
});