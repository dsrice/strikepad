import {render, screen, fireEvent, waitFor, act} from '@testing-library/react';
import GoogleSignButton from './GoogleSignButton';

// Mock window.google
const mockInitTokenClient = jest.fn();
const mockRequestAccessToken = jest.fn();

const mockGoogleAPI = {
    accounts: {
        oauth2: {
            initTokenClient: mockInitTokenClient,
        },
    },
};

// Mock import.meta.env is already handled in setupTests.ts
// We'll just modify the mock as needed for specific tests
const mockEnv = {
    VITE_GOOGLE_CLIENT_ID: 'test-client-id'
};

describe('GoogleSignButton', () => {
    const mockOnSuccess = jest.fn();
    const mockOnError = jest.fn();

    beforeEach(() => {
        jest.clearAllMocks();
        jest.clearAllTimers();
        jest.useFakeTimers();

        // Reset window.google
        (window as any).google = mockGoogleAPI;

        mockInitTokenClient.mockReturnValue({
            requestAccessToken: mockRequestAccessToken,
        });
    });

    afterEach(() => {
        jest.runOnlyPendingTimers();
        jest.useRealTimers();
        delete (window as any).google;
    });

    it('renders Google sign button when API is available', async () => {
        render(
            <GoogleSignButton
                onSuccess={mockOnSuccess}
                onError={mockOnError}
            />
        );

        // Fast-forward timers to complete initialization
        act(() => {
            jest.advanceTimersByTime(1000);
        });

        await waitFor(() => {
            expect(screen.getByText('Googleでサインイン')).toBeInTheDocument();
        });

        expect(mockInitTokenClient).toHaveBeenCalledWith({
            client_id: 'test-client-id',
            scope: 'email profile',
            callback: expect.any(Function),
        });
    });

    it('renders custom button text', async () => {
        render(
            <GoogleSignButton
                onSuccess={mockOnSuccess}
                onError={mockOnError}
                buttonText="カスタムテキスト"
            />
        );

        act(() => {
            jest.advanceTimersByTime(1000);
        });

        await waitFor(() => {
            expect(screen.getByText('カスタムテキスト')).toBeInTheDocument();
        });
    });

    it('does not render when Google API is not available', async () => {
        delete (window as any).google;

        render(
            <GoogleSignButton
                onSuccess={mockOnSuccess}
                onError={mockOnError}
            />
        );

        act(() => {
            jest.advanceTimersByTime(6000); // More than max retry time
        });

        await waitFor(() => {
            expect(screen.queryByText('Googleでサインイン')).not.toBeInTheDocument();
        });
    });

    it('does not render when client ID is not configured', async () => {
        // Mock empty client ID for this test
        const originalImport = (global as any).import;
        (global as any).import = {
            meta: {
                env: {
                    VITE_GOOGLE_CLIENT_ID: ''
                }
            }
        };

        render(
            <GoogleSignButton
                onSuccess={mockOnSuccess}
                onError={mockOnError}
            />
        );

        act(() => {
            jest.advanceTimersByTime(1000);
        });

        // Should not render the button when client ID is missing
        expect(screen.queryByText('Googleでサインイン')).not.toBeInTheDocument();

        // Restore original mock
        (global as any).import = originalImport;
    });

    it('calls requestAccessToken when button is clicked', async () => {
        render(
            <GoogleSignButton
                onSuccess={mockOnSuccess}
                onError={mockOnError}
            />
        );

        act(() => {
            jest.advanceTimersByTime(1000);
        });

        await waitFor(() => {
            expect(screen.getByText('Googleでサインイン')).toBeInTheDocument();
        });

        const button = screen.getByText('Googleでサインイン');
        fireEvent.click(button);

        expect(mockRequestAccessToken).toHaveBeenCalledTimes(1);
    });

    it('calls onSuccess when OAuth succeeds', async () => {
        const mockCallback = jest.fn();
        mockInitTokenClient.mockImplementation((config) => {
            mockCallback.mockImplementation(config.callback);
            return {
                requestAccessToken: mockRequestAccessToken,
            };
        });

        render(
            <GoogleSignButton
                onSuccess={mockOnSuccess}
                onError={mockOnError}
            />
        );

        act(() => {
            jest.advanceTimersByTime(1000);
        });

        await waitFor(() => {
            expect(mockInitTokenClient).toHaveBeenCalled();
        });

        // Simulate successful OAuth response
        const successResponse = {
            access_token: 'test-access-token',
        };

        act(() => {
            mockCallback(successResponse);
        });

        expect(mockOnSuccess).toHaveBeenCalledWith('test-access-token');
    });

    it('calls onError when OAuth fails with error', async () => {
        const mockCallback = jest.fn();
        mockInitTokenClient.mockImplementation((config) => {
            mockCallback.mockImplementation(config.callback);
            return {
                requestAccessToken: mockRequestAccessToken,
            };
        });

        render(
            <GoogleSignButton
                onSuccess={mockOnSuccess}
                onError={mockOnError}
            />
        );

        act(() => {
            jest.advanceTimersByTime(1000);
        });

        await waitFor(() => {
            expect(mockInitTokenClient).toHaveBeenCalled();
        });

        // Simulate OAuth error response
        const errorResponse = {
            error: 'access_denied',
        };

        act(() => {
            mockCallback(errorResponse);
        });

        expect(mockOnError).toHaveBeenCalledWith('Google認証に失敗しました');
    });

    it('calls onError when no access token is received', async () => {
        const mockCallback = jest.fn();
        mockInitTokenClient.mockImplementation((config) => {
            mockCallback.mockImplementation(config.callback);
            return {
                requestAccessToken: mockRequestAccessToken,
            };
        });

        render(
            <GoogleSignButton
                onSuccess={mockOnSuccess}
                onError={mockOnError}
            />
        );

        act(() => {
            jest.advanceTimersByTime(1000);
        });

        await waitFor(() => {
            expect(mockInitTokenClient).toHaveBeenCalled();
        });

        // Simulate response without access token
        const emptyResponse = {};

        act(() => {
            mockCallback(emptyResponse);
        });

        expect(mockOnError).toHaveBeenCalledWith('アクセストークンが取得できませんでした');
    });

    it('handles initialization errors gracefully', async () => {
        mockInitTokenClient.mockImplementation(() => {
            throw new Error('Initialization failed');
        });

        render(
            <GoogleSignButton
                onSuccess={mockOnSuccess}
                onError={mockOnError}
            />
        );

        act(() => {
            jest.advanceTimersByTime(1000);
        });

        expect(mockOnError).toHaveBeenCalledWith('Google認証の初期化に失敗しました');
    });

    it('handles request access token errors', async () => {
        mockRequestAccessToken.mockImplementation(() => {
            throw new Error('Request failed');
        });

        render(
            <GoogleSignButton
                onSuccess={mockOnSuccess}
                onError={mockOnError}
            />
        );

        act(() => {
            jest.advanceTimersByTime(1000);
        });

        await waitFor(() => {
            expect(screen.getByText('Googleでサインイン')).toBeInTheDocument();
        });

        const button = screen.getByText('Googleでサインイン');
        fireEvent.click(button);

        expect(mockOnError).toHaveBeenCalledWith('Google認証に失敗しました');
    });

    it('shows error when client is not ready', async () => {
        mockInitTokenClient.mockReturnValue(null);

        render(
            <GoogleSignButton
                onSuccess={mockOnSuccess}
                onError={mockOnError}
            />
        );

        act(() => {
            jest.advanceTimersByTime(1000);
        });

        await waitFor(() => {
            expect(screen.getByText('Googleでサインイン')).toBeInTheDocument();
        });

        const button = screen.getByText('Googleでサインイン');
        fireEvent.click(button);

        expect(mockOnError).toHaveBeenCalledWith('Google認証の準備ができていません');
    });

    it('disables button when disabled prop is true', async () => {
        render(
            <GoogleSignButton
                onSuccess={mockOnSuccess}
                onError={mockOnError}
                disabled={true}
            />
        );

        act(() => {
            jest.advanceTimersByTime(1000);
        });

        await waitFor(() => {
            const button = screen.getByText('Googleでサインイン').closest('button');
            expect(button).toBeDisabled();
        });
    });

    it('does not call requestAccessToken when disabled', async () => {
        render(
            <GoogleSignButton
                onSuccess={mockOnSuccess}
                onError={mockOnError}
                disabled={true}
            />
        );

        act(() => {
            jest.advanceTimersByTime(1000);
        });

        await waitFor(() => {
            expect(screen.getByText('Googleでサインイン')).toBeInTheDocument();
        });

        const button = screen.getByText('Googleでサインイン');
        fireEvent.click(button);

        expect(mockRequestAccessToken).not.toHaveBeenCalled();
    });

    it('applies hover styles correctly', async () => {
        render(
            <GoogleSignButton
                onSuccess={mockOnSuccess}
                onError={mockOnError}
            />
        );

        act(() => {
            jest.advanceTimersByTime(1000);
        });

        await waitFor(() => {
            expect(screen.getByText('Googleでサインイン')).toBeInTheDocument();
        });

        const button = screen.getByText('Googleでサインイン').closest('button') as HTMLButtonElement;

        // Test hover effect
        fireEvent.mouseEnter(button);
        expect(button.style.backgroundColor).toBe('rgb(249, 250, 251)');

        fireEvent.mouseLeave(button);
        expect(button.style.backgroundColor).toBe('rgb(255, 255, 255)');
    });

    it('retries initialization when Google API loads later', async () => {
        // Start without Google API
        delete (window as any).google;

        render(
            <GoogleSignButton
                onSuccess={mockOnSuccess}
                onError={mockOnError}
            />
        );

        // Advance time for first retry attempts
        act(() => {
            jest.advanceTimersByTime(1000);
        });

        // API still not available
        expect(screen.queryByText('Googleでサインイン')).not.toBeInTheDocument();

        // Now make Google API available
        (window as any).google = mockGoogleAPI;

        // Continue advancing time to trigger retry
        act(() => {
            jest.advanceTimersByTime(1000);
        });

        await waitFor(() => {
            expect(screen.getByText('Googleでサインイン')).toBeInTheDocument();
        });

        expect(mockInitTokenClient).toHaveBeenCalled();
    });
});