import {render, screen, fireEvent, waitFor, act} from '@testing-library/react';
import GoogleSignSection from './GoogleSignSection';

// Mock GoogleSignButton
jest.mock('./GoogleSignButton', () => {
    return function MockGoogleSignButton({onSuccess, onError, buttonText, disabled}: any) {
        return (
            <div data-testid="google-sign-button">
                <button
                    onClick={() => onSuccess('mock-access-token')}
                    disabled={disabled}
                    data-testid="google-button-success"
                >
                    {buttonText}
                </button>
                <button
                    onClick={() => onError('Mock error')}
                    data-testid="google-button-error"
                >
                    Trigger Error
                </button>
            </div>
        );
    };
});

// Mock import.meta.env is already handled in setupTests.ts
const mockEnv = {
    VITE_GOOGLE_CLIENT_ID: 'test-client-id'
};

// Mock document.querySelector
const mockQuerySelector = jest.fn();
Object.defineProperty(document, 'querySelector', {
    value: mockQuerySelector,
    writable: true
});

describe('GoogleSignSection', () => {
    const mockOnSuccess = jest.fn();
    const mockOnError = jest.fn();

    beforeEach(() => {
        jest.clearAllMocks();
        jest.useFakeTimers();

        // Mock Google script presence by default
        mockQuerySelector.mockReturnValue(document.createElement('script'));
    });

    afterEach(() => {
        jest.runOnlyPendingTimers();
        jest.useRealTimers();
    });

    it('renders Google sign section when Google is available', async () => {
        render(
            <GoogleSignSection
                onSuccess={mockOnSuccess}
                onError={mockOnError}
            />
        );

        act(() => {
            jest.advanceTimersByTime(200);
        });

        await waitFor(() => {
            expect(screen.getByTestId('google-sign-button')).toBeInTheDocument();
            expect(screen.getByText('または')).toBeInTheDocument();
        });
    });

    it('renders with custom button text', async () => {
        render(
            <GoogleSignSection
                onSuccess={mockOnSuccess}
                onError={mockOnError}
                buttonText="カスタムサインインテキスト"
            />
        );

        act(() => {
            jest.advanceTimersByTime(200);
        });

        await waitFor(() => {
            expect(screen.getByText('カスタムサインインテキスト')).toBeInTheDocument();
        });
    });

    it('does not render when Google script is not available', async () => {
        mockQuerySelector.mockReturnValue(null);

        render(
            <GoogleSignSection
                onSuccess={mockOnSuccess}
                onError={mockOnError}
            />
        );

        act(() => {
            jest.advanceTimersByTime(200);
        });

        await waitFor(() => {
            expect(screen.queryByTestId('google-sign-button')).not.toBeInTheDocument();
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
            <GoogleSignSection
                onSuccess={mockOnSuccess}
                onError={mockOnError}
            />
        );

        act(() => {
            jest.advanceTimersByTime(200);
        });

        await waitFor(() => {
            expect(screen.queryByTestId('google-sign-button')).not.toBeInTheDocument();
        });

        // Restore original mock
        (global as any).import = originalImport;
    });

    it('forwards onSuccess calls to GoogleSignButton', async () => {
        render(
            <GoogleSignSection
                onSuccess={mockOnSuccess}
                onError={mockOnError}
            />
        );

        act(() => {
            jest.advanceTimersByTime(200);
        });

        await waitFor(() => {
            expect(screen.getByTestId('google-button-success')).toBeInTheDocument();
        });

        const successButton = screen.getByTestId('google-button-success');
        fireEvent.click(successButton);

        expect(mockOnSuccess).toHaveBeenCalledWith('mock-access-token');
    });

    it('forwards onError calls to GoogleSignButton', async () => {
        render(
            <GoogleSignSection
                onSuccess={mockOnSuccess}
                onError={mockOnError}
            />
        );

        act(() => {
            jest.advanceTimersByTime(200);
        });

        await waitFor(() => {
            expect(screen.getByTestId('google-button-error')).toBeInTheDocument();
        });

        const errorButton = screen.getByTestId('google-button-error');
        fireEvent.click(errorButton);

        expect(mockOnError).toHaveBeenCalledWith('Mock error');
    });

    it('passes disabled prop to GoogleSignButton', async () => {
        render(
            <GoogleSignSection
                onSuccess={mockOnSuccess}
                onError={mockOnError}
                disabled={true}
            />
        );

        act(() => {
            jest.advanceTimersByTime(200);
        });

        await waitFor(() => {
            const button = screen.getByTestId('google-button-success');
            expect(button).toBeDisabled();
        });
    });

    it('renders divider with correct styling', async () => {
        render(
            <GoogleSignSection
                onSuccess={mockOnSuccess}
                onError={mockOnError}
            />
        );

        act(() => {
            jest.advanceTimersByTime(200);
        });

        await waitFor(() => {
            const divider = screen.getByText('または');
            expect(divider).toHaveClass('px-2', 'bg-gray-50', 'text-gray-500');
        });
    });

    it('has proper container structure', async () => {
        render(
            <GoogleSignSection
                onSuccess={mockOnSuccess}
                onError={mockOnError}
            />
        );

        act(() => {
            jest.advanceTimersByTime(200);
        });

        await waitFor(() => {
            const container = screen.getByTestId('google-sign-button').parentElement;
            expect(container).toHaveClass('space-y-4');
        });
    });

    it('checks for Google script with correct selector', async () => {
        render(
            <GoogleSignSection
                onSuccess={mockOnSuccess}
                onError={mockOnError}
            />
        );

        act(() => {
            jest.advanceTimersByTime(200);
        });

        await waitFor(() => {
            expect(mockQuerySelector).toHaveBeenCalledWith('script[src*="accounts.google.com"]');
        });
    });

    it('renders default button text when none provided', async () => {
        render(
            <GoogleSignSection
                onSuccess={mockOnSuccess}
                onError={mockOnError}
            />
        );

        act(() => {
            jest.advanceTimersByTime(200);
        });

        await waitFor(() => {
            expect(screen.getByText('Googleでサインイン')).toBeInTheDocument();
        });
    });

    it('delays availability check by 100ms', async () => {
        render(
            <GoogleSignSection
                onSuccess={mockOnSuccess}
                onError={mockOnError}
            />
        );

        // Should not be available immediately
        expect(screen.queryByTestId('google-sign-button')).not.toBeInTheDocument();

        // Fast forward just under 100ms
        act(() => {
            jest.advanceTimersByTime(99);
        });

        expect(screen.queryByTestId('google-sign-button')).not.toBeInTheDocument();

        // Fast forward past 100ms
        act(() => {
            jest.advanceTimersByTime(2);
        });

        await waitFor(() => {
            expect(screen.getByTestId('google-sign-button')).toBeInTheDocument();
        });
    });

    it('handles environment where both requirements are missing', async () => {
        mockQuerySelector.mockReturnValue(null);

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
            <GoogleSignSection
                onSuccess={mockOnSuccess}
                onError={mockOnError}
            />
        );

        act(() => {
            jest.advanceTimersByTime(200);
        });

        await waitFor(() => {
            expect(screen.queryByTestId('google-sign-button')).not.toBeInTheDocument();
            expect(screen.queryByText('または')).not.toBeInTheDocument();
        });

        // Restore original mock
        (global as any).import = originalImport;
    });

    it('returns null when Google is not available', () => {
        mockQuerySelector.mockReturnValue(null);

        // Mock empty client ID for this test
        const originalImport = (global as any).import;
        (global as any).import = {
            meta: {
                env: {
                    VITE_GOOGLE_CLIENT_ID: ''
                }
            }
        };

        const {container} = render(
            <GoogleSignSection
                onSuccess={mockOnSuccess}
                onError={mockOnError}
            />
        );

        act(() => {
            jest.advanceTimersByTime(200);
        });

        expect(container.firstChild).toBeNull();

        // Restore original mock
        (global as any).import = originalImport;
    });
});