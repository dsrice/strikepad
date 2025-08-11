import {render, screen, fireEvent, waitFor} from '@testing-library/react';
import {BrowserRouter} from 'react-router-dom';
import LoginForm from './LoginForm';

// Mock the AuthContext
const mockLogin = jest.fn();
const mockNavigate = jest.fn();
let mockIsLoading = false;
let mockError = null;

jest.mock('../contexts/AuthContext', () => ({
    useAuth: () => ({
        isAuthenticated: false,
        isLoading: mockIsLoading,
        error: mockError,
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
        mockIsLoading = false;
        mockError = null;
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

        // Check that login was not called due to validation
        expect(mockLogin).not.toHaveBeenCalled();
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
        let resolvePromise: (value?: any) => void;
        mockLogin.mockImplementation(() => new Promise(resolve => {
            resolvePromise = resolve;
        }));
        
        renderWithRouter(<LoginForm/>);

        const emailInput = screen.getByLabelText(/メールアドレス/i);
        const passwordInput = screen.getByLabelText(/パスワード/i);
        const submitButton = screen.getByRole('button', {name: /ログイン/i});

        fireEvent.change(emailInput, {target: {value: 'test@example.com'}});
        fireEvent.change(passwordInput, {target: {value: 'password123'}});
        fireEvent.click(submitButton);

        // Check loading state immediately after form submission
        await waitFor(() => {
            expect(screen.getByText(/ログイン中.../i)).toBeInTheDocument();
        });
        expect(submitButton).toBeDisabled();

        // Resolve the promise to complete the test
        resolvePromise!({});
    });

    it('handles login error', async () => {
        mockError = 'ログインに失敗しました';
        mockLogin.mockRejectedValue(new Error('ログインに失敗しました'));

        const {rerender} = renderWithRouter(<LoginForm/>);
        rerender(<LoginForm/>);

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
        let resolvePromise: (value?: any) => void;
        mockLogin.mockImplementation(() => new Promise(resolve => {
            resolvePromise = resolve;
        }));

        renderWithRouter(<LoginForm/>);

        const emailInput = screen.getByLabelText(/メールアドレス/i);
        const passwordInput = screen.getByLabelText(/パスワード/i);
        const submitButton = screen.getByRole('button', {name: /ログイン/i});

        fireEvent.change(emailInput, {target: {value: 'test@example.com'}});
        fireEvent.change(passwordInput, {target: {value: 'password123'}});
        fireEvent.click(submitButton);

        // Check that inputs are disabled during loading
        await waitFor(() => {
            expect(emailInput).toBeDisabled();
            expect(passwordInput).toBeDisabled();
        });

        // Resolve the promise to complete the test
        resolvePromise!({});
    });
});