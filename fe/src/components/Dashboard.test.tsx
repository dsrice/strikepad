import {render, screen, fireEvent} from '@testing-library/react';
import Dashboard from './Dashboard';

// Mock the AuthContext
const mockLogout = jest.fn();
const mockUser = {
    id: 1,
    email: 'test@example.com',
    display_name: 'Test User',
    email_verified: true,
    createdAt: '2023-01-01T00:00:00Z',
    updatedAt: '2023-01-01T00:00:00Z',
};

const mockUseAuth = jest.fn();

jest.mock('../contexts/AuthContext', () => ({
    useAuth: () => mockUseAuth(),
}));

describe('Dashboard', () => {
    beforeEach(() => {
        jest.clearAllMocks();
        mockUseAuth.mockReturnValue({
            user: mockUser,
            logout: mockLogout,
        });
    });

    it('renders dashboard with user information', () => {
        render(<Dashboard/>);

        expect(screen.getByText('StrikePad Dashboard')).toBeInTheDocument();
        expect(screen.getByText('ようこそ、Test Userさん')).toBeInTheDocument();
        expect(screen.getByText('ユーザー情報')).toBeInTheDocument();
        expect(screen.getByText('1')).toBeInTheDocument();
        expect(screen.getByText('test@example.com')).toBeInTheDocument();
        expect(screen.getByText('Test User')).toBeInTheDocument();
        expect(screen.getByText('認証済み')).toBeInTheDocument();
        expect(screen.getByText('D3.js Chart Demo')).toBeInTheDocument();
    });

    it('displays unverified email status when user email is not verified', () => {
        const unverifiedUser = {
            ...mockUser,
            email_verified: false,
        };

        mockUseAuth.mockReturnValue({
            user: unverifiedUser,
            logout: mockLogout,
        });

        render(<Dashboard/>);

        expect(screen.getByText('未認証')).toBeInTheDocument();
        expect(screen.getByText('未認証')).toHaveClass('text-red-600');
    });

    it('applies correct styling for verified email status', () => {
        render(<Dashboard/>);

        const verifiedStatus = screen.getByText('認証済み');
        expect(verifiedStatus).toHaveClass('text-green-600');
    });

    it('calls logout function when logout button is clicked', () => {
        render(<Dashboard/>);

        const logoutButton = screen.getByRole('button', {name: /ログアウト/i});
        fireEvent.click(logoutButton);

        expect(mockLogout).toHaveBeenCalledTimes(1);
    });

    it('renders properly when user displayName is null or undefined', () => {
        const userWithoutDisplayName = {
            ...mockUser,
            display_name: undefined,
        };

        mockUseAuth.mockReturnValue({
            user: userWithoutDisplayName,
            logout: mockLogout,
        });

        render(<Dashboard/>);

        expect(screen.getByText('ようこそ、さん')).toBeInTheDocument();
    });

    it('renders properly when user data is null', () => {
        mockUseAuth.mockReturnValue({
            user: null,
            logout: mockLogout,
        });

        render(<Dashboard/>);

        expect(screen.getByText('ようこそ、さん')).toBeInTheDocument();
        expect(screen.queryByText('1')).not.toBeInTheDocument();
    });

    it('renders all user information fields', () => {
        render(<Dashboard/>);

        expect(screen.getByText('ID:')).toBeInTheDocument();
        expect(screen.getByText('メール:')).toBeInTheDocument();
        expect(screen.getByText('表示名:')).toBeInTheDocument();
        expect(screen.getByText('メール認証:')).toBeInTheDocument();
    });

    it('has correct layout structure', () => {
        render(<Dashboard/>);

        // Header section
        const header = screen.getByRole('banner');
        expect(header).toHaveClass('bg-blue-600');

        // Main content
        const main = screen.getByRole('main');
        expect(main).toBeInTheDocument();

        // Logout button
        const logoutButton = screen.getByRole('button', {name: /ログアウト/i});
        expect(logoutButton).toHaveClass('bg-blue-500');
    });

    it('displays user information in correct format', () => {
        render(<Dashboard/>);

        // Check that ID is displayed as text content
        expect(screen.getByText((_, element) => {
            return element?.textContent === 'ID: 1';
        })).toBeInTheDocument();

        // Check email format
        expect(screen.getByText((_, element) => {
            return element?.textContent === 'メール: test@example.com';
        })).toBeInTheDocument();

        // Check display name format
        expect(screen.getByText((_, element) => {
            return element?.textContent === '表示名: Test User';
        })).toBeInTheDocument();
    });

    it('has hover effect on logout button', () => {
        render(<Dashboard/>);

        const logoutButton = screen.getByRole('button', {name: /ログアウト/i});
        expect(logoutButton).toHaveClass('hover:bg-blue-400');
    });

    it('uses grid layout for content sections', () => {
        render(<Dashboard/>);

        const gridContainer = screen.getByText('ユーザー情報').closest('div')?.parentElement;
        expect(gridContainer).toHaveClass('grid');
        expect(gridContainer).toHaveClass('grid-cols-1');
        expect(gridContainer).toHaveClass('lg:grid-cols-2');
    });
});