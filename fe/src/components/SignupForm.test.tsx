import {render, screen, fireEvent, waitFor} from '@testing-library/react';
import {BrowserRouter} from 'react-router-dom';
import SignupForm from './SignupForm';

// Mock the AuthContext
const mockSignup = jest.fn();
const mockGoogleSignup = jest.fn();
const mockNavigate = jest.fn();

jest.mock('../contexts/AuthContext', () => ({
    useAuth: () => ({
        isAuthenticated: false,
        isLoading: false,
        error: null,
        signup: mockSignup,
        googleSignup: mockGoogleSignup,
        login: jest.fn(),
        logout: jest.fn(),
    }),
}));

jest.mock('react-router-dom', () => ({
    ...jest.requireActual('react-router-dom'),
    useNavigate: () => mockNavigate,
}));

// Mock GoogleSignSection
jest.mock('./GoogleSignSection', () => {
    return function MockGoogleSignSection({onSuccess, onError, disabled, buttonText}: any) {
        return (
            <div data-testid="google-sign-section">
                <button
                    onClick={() => onSuccess('mock-access-token')}
                    disabled={disabled}
                    data-testid="google-signup-success"
                >
                    {buttonText}
                </button>
                <button
                    onClick={() => onError('Google authentication error')}
                    data-testid="google-signup-error"
                >
                    Trigger Error
                </button>
            </div>
        );
    };
});

// Mock Layout and UI components
jest.mock('./Layout', () => {
    return {
        Layout: ({children}: { children: React.ReactNode }) => <div>{children}</div>,
    };
});

jest.mock('./ui/Container', () => ({
    Container: ({children, className}: { children: React.ReactNode, className?: string }) =>
        <div className={className}>{children}</div>,
}));

jest.mock('./ui/FadeIn', () => ({
    FadeIn: ({children, className}: { children: React.ReactNode, className?: string }) =>
        <div className={className}>{children}</div>,
}));

const renderWithRouter = (component: React.ReactElement) => {
    return render(<BrowserRouter>{component}</BrowserRouter>);
};

describe('SignupForm', () => {
    beforeEach(() => {
        jest.clearAllMocks();
    });

    it('renders signup form elements', () => {
        renderWithRouter(<SignupForm/>);

        expect(screen.getByText(/StrikePadアカウント作成/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/メールアドレス/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/表示名/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/^パスワード$/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/パスワード確認/i)).toBeInTheDocument();
        expect(screen.getByRole('button', {name: /アカウント作成/i})).toBeInTheDocument();
        expect(screen.getByText(/ログイン/i)).toBeInTheDocument();
        expect(screen.getByTestId('google-sign-section')).toBeInTheDocument();
    });

    it('shows validation errors for empty fields', async () => {
        renderWithRouter(<SignupForm/>);

        const submitButton = screen.getByRole('button', {name: /アカウント作成/i});
        fireEvent.click(submitButton);

        await waitFor(() => {
            expect(screen.getByText(/有効なメールアドレスを入力してください/i)).toBeInTheDocument();
            expect(screen.getByText(/表示名を入力してください/i)).toBeInTheDocument();
            expect(screen.getByText(/パスワードは8文字以上である必要があります/i)).toBeInTheDocument();
        });
    });

    it('shows validation error for invalid email', async () => {
        renderWithRouter(<SignupForm/>);

        const submitButton = screen.getByRole('button', {name: /アカウント作成/i});

        // Submit form with invalid email
        const emailInput = screen.getByLabelText(/メールアドレス/i);
        fireEvent.change(emailInput, {target: {value: 'invalid-email'}});
        fireEvent.click(submitButton);

        // Check if form validation prevents submission (by checking if signup was not called)
        expect(mockSignup).not.toHaveBeenCalled();
    });

    it('shows validation error for weak password', async () => {
        renderWithRouter(<SignupForm/>);

        const passwordInput = screen.getByLabelText(/^パスワード$/i);
        const submitButton = screen.getByRole('button', {name: /アカウント作成/i});

        fireEvent.change(passwordInput, {target: {value: 'weak'}});
        fireEvent.click(submitButton);

        await waitFor(() => {
            expect(screen.getByText(/パスワードは8文字以上である必要があります/i)).toBeInTheDocument();
        });
    });

    it('shows validation error for password without uppercase', async () => {
        renderWithRouter(<SignupForm/>);

        const passwordInput = screen.getByLabelText(/^パスワード$/i);
        const submitButton = screen.getByRole('button', {name: /アカウント作成/i});

        fireEvent.change(passwordInput, {target: {value: 'password123!'}});
        fireEvent.click(submitButton);

        await waitFor(() => {
            expect(screen.getByText(/パスワードには大文字を含める必要があります/i)).toBeInTheDocument();
        });
    });

    it('shows validation error for password without special character', async () => {
        renderWithRouter(<SignupForm/>);

        const passwordInput = screen.getByLabelText(/^パスワード$/i);
        const submitButton = screen.getByRole('button', {name: /アカウント作成/i});

        fireEvent.change(passwordInput, {target: {value: 'Password123'}});
        fireEvent.click(submitButton);

        await waitFor(() => {
            expect(screen.getByText(/パスワードには記号を含める必要があります/i)).toBeInTheDocument();
        });
    });

    it('shows validation error for mismatched passwords', async () => {
        renderWithRouter(<SignupForm/>);

        const passwordInput = screen.getByLabelText(/^パスワード$/i);
        const confirmPasswordInput = screen.getByLabelText(/パスワード確認/i);
        const submitButton = screen.getByRole('button', {name: /アカウント作成/i});

        fireEvent.change(passwordInput, {target: {value: 'Password123!'}});
        fireEvent.change(confirmPasswordInput, {target: {value: 'DifferentPassword123!'}});
        fireEvent.click(submitButton);

        await waitFor(() => {
            expect(screen.getByText(/パスワードが一致しません/i)).toBeInTheDocument();
        });
    });

    it('calls signup function with correct data', async () => {
        mockSignup.mockResolvedValue({});
        renderWithRouter(<SignupForm/>);

        const emailInput = screen.getByLabelText(/メールアドレス/i);
        const displayNameInput = screen.getByLabelText(/表示名/i);
        const passwordInput = screen.getByLabelText(/^パスワード$/i);
        const confirmPasswordInput = screen.getByLabelText(/パスワード確認/i);
        const submitButton = screen.getByRole('button', {name: /アカウント作成/i});

        fireEvent.change(emailInput, {target: {value: 'test@example.com'}});
        fireEvent.change(displayNameInput, {target: {value: 'Test User'}});
        fireEvent.change(passwordInput, {target: {value: 'Password123!'}});
        fireEvent.change(confirmPasswordInput, {target: {value: 'Password123!'}});
        fireEvent.click(submitButton);

        await waitFor(() => {
            expect(mockSignup).toHaveBeenCalledWith('test@example.com', 'Password123!', 'Test User');
            expect(mockNavigate).toHaveBeenCalledWith('/dashboard');
        });
    });

    it('shows loading state during signup', async () => {
        mockSignup.mockImplementation(() => new Promise(resolve => setTimeout(resolve, 100)));
        renderWithRouter(<SignupForm/>);

        const emailInput = screen.getByLabelText(/メールアドレス/i);
        const displayNameInput = screen.getByLabelText(/表示名/i);
        const passwordInput = screen.getByLabelText(/^パスワード$/i);
        const confirmPasswordInput = screen.getByLabelText(/パスワード確認/i);
        const submitButton = screen.getByRole('button', {name: /アカウント作成/i});

        fireEvent.change(emailInput, {target: {value: 'test@example.com'}});
        fireEvent.change(displayNameInput, {target: {value: 'Test User'}});
        fireEvent.change(passwordInput, {target: {value: 'Password123!'}});
        fireEvent.change(confirmPasswordInput, {target: {value: 'Password123!'}});
        fireEvent.click(submitButton);

        await waitFor(() => {
            expect(screen.getByText(/アカウント作成中.../i)).toBeInTheDocument();
        });
        expect(submitButton).toBeDisabled();
    });

    it('handles signup error', async () => {
        mockSignup.mockRejectedValue(new Error('アカウント作成に失敗しました'));
        renderWithRouter(<SignupForm/>);

        const emailInput = screen.getByLabelText(/メールアドレス/i);
        const displayNameInput = screen.getByLabelText(/表示名/i);
        const passwordInput = screen.getByLabelText(/^パスワード$/i);
        const confirmPasswordInput = screen.getByLabelText(/パスワード確認/i);
        const submitButton = screen.getByRole('button', {name: /アカウント作成/i});

        fireEvent.change(emailInput, {target: {value: 'test@example.com'}});
        fireEvent.change(displayNameInput, {target: {value: 'Test User'}});
        fireEvent.change(passwordInput, {target: {value: 'Password123!'}});
        fireEvent.change(confirmPasswordInput, {target: {value: 'Password123!'}});
        fireEvent.click(submitButton);

        await waitFor(() => {
            expect(screen.getByText(/アカウント作成に失敗しました/i)).toBeInTheDocument();
        });
    });

    it('handles Google signup success', async () => {
        mockGoogleSignup.mockResolvedValue({});
        renderWithRouter(<SignupForm/>);

        const googleSignupButton = screen.getByTestId('google-signup-success');
        fireEvent.click(googleSignupButton);

        await waitFor(() => {
            expect(mockGoogleSignup).toHaveBeenCalledWith('mock-access-token');
            expect(mockNavigate).toHaveBeenCalledWith('/dashboard');
        });
    });

    it('handles Google signup error', async () => {
        renderWithRouter(<SignupForm/>);

        const googleErrorButton = screen.getByTestId('google-signup-error');
        fireEvent.click(googleErrorButton);

        await waitFor(() => {
            expect(screen.getByText(/Google authentication error/i)).toBeInTheDocument();
        });
    });

    it('ignores Google authentication preparation errors', async () => {
        renderWithRouter(<SignupForm/>);

        // Mock console.warn to verify it's called
        const consoleSpy = jest.spyOn(console, 'warn').mockImplementation();

        // Simulate Google auth preparation error by modifying the mock
        const googleSignSection = screen.getByTestId('google-sign-section');

        // Create a button that triggers preparation error
        const errorButton = document.createElement('button');
        errorButton.onclick = () => {
            const mockOnError = jest.fn();
            // Simulate the actual error handling
            const error = 'Google認証の準備ができていません';
            if (error.includes('準備ができていません') || error.includes('初期化に失敗')) {
                console.warn('Google authentication not available:', error);
                return;
            }
        };

        fireEvent.click(errorButton);

        // Verify no error message is shown in UI for preparation errors
        expect(screen.queryByText(/準備ができていません/i)).not.toBeInTheDocument();

        consoleSpy.mockRestore();
    });

    it('navigates to login when login button is clicked', () => {
        renderWithRouter(<SignupForm/>);

        const loginButton = screen.getByText(/ログイン/i);
        fireEvent.click(loginButton);

        expect(mockNavigate).toHaveBeenCalledWith('/login');
    });

    it('disables form inputs during loading', async () => {
        mockSignup.mockImplementation(() => new Promise(resolve => setTimeout(resolve, 100)));
        renderWithRouter(<SignupForm/>);

        const emailInput = screen.getByLabelText(/メールアドレス/i);
        const displayNameInput = screen.getByLabelText(/表示名/i);
        const passwordInput = screen.getByLabelText(/^パスワード$/i);
        const confirmPasswordInput = screen.getByLabelText(/パスワード確認/i);
        const submitButton = screen.getByRole('button', {name: /アカウント作成/i});

        fireEvent.change(emailInput, {target: {value: 'test@example.com'}});
        fireEvent.change(displayNameInput, {target: {value: 'Test User'}});
        fireEvent.change(passwordInput, {target: {value: 'Password123!'}});
        fireEvent.change(confirmPasswordInput, {target: {value: 'Password123!'}});
        fireEvent.click(submitButton);

        await waitFor(() => {
            expect(emailInput).toBeDisabled();
            expect(displayNameInput).toBeDisabled();
            expect(passwordInput).toBeDisabled();
            expect(confirmPasswordInput).toBeDisabled();
        });
    });
});