import {render, screen, fireEvent, waitFor} from '@testing-library/react';
import {BrowserRouter} from 'react-router-dom';
import LoginForm from './LoginForm';

// Mock the AuthContext
const mockLogin = jest.fn();
const mockNavigate = jest.fn();

jest.mock('../contexts/AuthContext', () => ({
    useAuth: () => ({
        isAuthenticated: false,
        isLoading: false,
        error: null,
        login: mockLogin,
        logout: jest.fn(),
    }),
}));

jest.mock('react-router-dom', () => ({
    ...jest.requireActual('react-router-dom'),
    useNavigate: () => mockNavigate,
}));

const renderWithRouter = (component: React.ReactElement) => {
    return render(<BrowserRouter>{component}</BrowserRouter>);
};

describe('LoginForm', () => {
    beforeEach(() => {
        jest.clearAllMocks();
    });

    it('renders login form elements', () => {
        renderWithRouter(<LoginForm/>);

        expect(screen.getByText(/StrikePadにログイン/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/メールアドレス/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/パスワード/i)).toBeInTheDocument();
        expect(screen.getByRole('button', {name: /ログイン/i})).toBeInTheDocument();
        expect(screen.getByText(/新規登録/i)).toBeInTheDocument();
    });

    it('shows validation errors for empty fields', async () => {
        renderWithRouter(<LoginForm/>);

        const submitButton = screen.getByRole('button', {name: /ログイン/i});
        fireEvent.click(submitButton);

        await waitFor(() => {
            expect(screen.getByText(/有効なメールアドレスを入力してください/i)).toBeInTheDocument();
            expect(screen.getByText(/パスワードを入力してください/i)).toBeInTheDocument();
        });
    });

    it('shows validation error for invalid email', async () => {
        renderWithRouter(<LoginForm/>);

        const emailInput = screen.getByLabelText(/メールアドレス/i);
        const submitButton = screen.getByRole('button', {name: /ログイン/i});

        fireEvent.change(emailInput, {target: {value: 'invalid-email'}});
        fireEvent.click(submitButton);

        await waitFor(() => {
            expect(screen.getByText(/有効なメールアドレスを入力してください/i)).toBeInTheDocument();
        });
    });

    it('calls login function with correct credentials', async () => {
        mockLogin.mockResolvedValue({});
        renderWithRouter(<LoginForm/>);

        const emailInput = screen.getByLabelText(/メールアドレス/i);
        const passwordInput = screen.getByLabelText(/パスワード/i);
        const submitButton = screen.getByRole('button', {name: /ログイン/i});

        fireEvent.change(emailInput, {target: {value: 'test@example.com'}});
        fireEvent.change(passwordInput, {target: {value: 'password123'}});
        fireEvent.click(submitButton);

        await waitFor(() => {
            expect(mockLogin).toHaveBeenCalledWith('test@example.com', 'password123');
            expect(mockNavigate).toHaveBeenCalledWith('/dashboard');
        });
    });

    it('shows loading state during login', async () => {
        mockLogin.mockImplementation(() => new Promise(resolve => setTimeout(resolve, 100)));
        renderWithRouter(<LoginForm/>);

        const emailInput = screen.getByLabelText(/メールアドレス/i);
        const passwordInput = screen.getByLabelText(/パスワード/i);
        const submitButton = screen.getByRole('button', {name: /ログイン/i});

        fireEvent.change(emailInput, {target: {value: 'test@example.com'}});
        fireEvent.change(passwordInput, {target: {value: 'password123'}});
        fireEvent.click(submitButton);

        expect(screen.getByText(/ログイン中.../i)).toBeInTheDocument();
        expect(submitButton).toBeDisabled();
    });

    it('handles login error', async () => {
        mockLogin.mockRejectedValue(new Error('ログインに失敗しました'));
        renderWithRouter(<LoginForm/>);

        const emailInput = screen.getByLabelText(/メールアドレス/i);
        const passwordInput = screen.getByLabelText(/パスワード/i);
        const submitButton = screen.getByRole('button', {name: /ログイン/i});

        fireEvent.change(emailInput, {target: {value: 'test@example.com'}});
        fireEvent.change(passwordInput, {target: {value: 'wrongpassword'}});
        fireEvent.click(submitButton);

        await waitFor(() => {
            expect(screen.getByText(/ログインに失敗しました/i)).toBeInTheDocument();
        });
    });

    it('navigates to signup when signup button is clicked', () => {
        renderWithRouter(<LoginForm/>);

        const signupButton = screen.getByText(/新規登録/i);
        fireEvent.click(signupButton);

        expect(mockNavigate).toHaveBeenCalledWith('/signup');
    });

    it('disables form inputs during loading', async () => {
        mockLogin.mockImplementation(() => new Promise(resolve => setTimeout(resolve, 100)));
        renderWithRouter(<LoginForm/>);

        const emailInput = screen.getByLabelText(/メールアドレス/i);
        const passwordInput = screen.getByLabelText(/パスワード/i);
        const submitButton = screen.getByRole('button', {name: /ログイン/i});

        fireEvent.change(emailInput, {target: {value: 'test@example.com'}});
        fireEvent.change(passwordInput, {target: {value: 'password123'}});
        fireEvent.click(submitButton);

        expect(emailInput).toBeDisabled();
        expect(passwordInput).toBeDisabled();
    });
});