import React, {useEffect, useRef} from 'react';

declare global {
    interface Window {
        google: {
            accounts: {
                id: {
                    initialize: (config: any) => void;
                    renderButton: (element: HTMLElement | null, options: any) => void;
                    prompt: () => void;
                };
                oauth2: {
                    initTokenClient: (config: any) => any;
                    revoke: (accessToken: string, callback: () => void) => void;
                };
            };
        };
    }
}

interface GoogleSignButtonProps {
    onSuccess: (accessToken: string) => void;
    onError: (error: string) => void;
    buttonText?: string;
    disabled?: boolean;
}

const GoogleSignButton: React.FC<GoogleSignButtonProps> = ({
                                                               onSuccess,
                                                               onError,
                                                               buttonText = 'Googleでサインイン',
                                                               disabled = false,
                                                           }) => {
    const buttonRef = useRef<HTMLDivElement>(null);
    const clientRef = useRef<any>(null);
    const [isAvailable, setIsAvailable] = React.useState(false);

    useEffect(() => {
        let retryCount = 0;
        const maxRetries = 50; // 5秒間リトライ (100ms * 50)

        const initializeGoogleSignIn = () => {
            if (window.google && window.google.accounts && window.google.accounts.oauth2) {
                try {
                    const googleClientId = (import.meta as any).env?.VITE_GOOGLE_CLIENT_ID || '';

                    if (!googleClientId) {
                        console.warn('Google Client ID not configured');
                        return; // エラーメッセージを表示しない、単に初期化をスキップ
                    }

                    clientRef.current = window.google.accounts.oauth2.initTokenClient({
                        client_id: googleClientId,
                        scope: 'email profile',
                        callback: (response: any) => {
                            if (response.error) {
                                console.error('Google OAuth error:', response.error);
                                onError('Google認証に失敗しました');
                                return;
                            }

                            if (response.access_token) {
                                onSuccess(response.access_token);
                            } else {
                                onError('アクセストークンが取得できませんでした');
                            }
                        },
                    });

                    setIsAvailable(true);
                } catch (error) {
                    console.error('Failed to initialize Google OAuth:', error);
                    onError('Google認証の初期化に失敗しました');
                }
            } else {
                // Google API not loaded yet, retry after a short delay
                retryCount++;
                if (retryCount < maxRetries) {
                    setTimeout(initializeGoogleSignIn, 100);
                } else {
                    console.warn('Google Identity Services API could not be loaded');
                    // 最大リトライ数に達してもエラーメッセージを表示しない
                }
            }
        };

        // 初期化を少し遅らせる
        const timeoutId = setTimeout(initializeGoogleSignIn, 500);

        return () => {
            clearTimeout(timeoutId);
        };
    }, [onSuccess, onError]);

    const handleClick = () => {
        if (disabled) {
            return;
        }

        if (clientRef.current) {
            try {
                clientRef.current.requestAccessToken();
            } catch (error) {
                console.error('Failed to request access token:', error);
                onError('Google認証に失敗しました');
            }
        } else {
            onError('Google認証の準備ができていません');
        }
    };

    // Google認証が利用可能でない場合はnullを返す（非表示）
    if (!isAvailable) {
        return null;
    }

    return (
        <div ref={buttonRef}>
            <button
                type="button"
                onClick={handleClick}
                disabled={disabled}
                className="w-full flex justify-center items-center px-3 py-3 border border-black rounded focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                style={{
                    fontFamily: 'Roboto, sans-serif',
                    fontSize: '14px',
                    fontWeight: 500,
                    color: '#1f1f1f',
                    minHeight: '48px',
                    backgroundColor: '#ffffff',
                    borderColor: '#1f1f1f',
                }}
                onMouseEnter={(e) => {
                    e.currentTarget.style.backgroundColor = '#f9fafb';
                }}
                onMouseLeave={(e) => {
                    e.currentTarget.style.backgroundColor = '#ffffff';
                }}
            >
                <div className="flex items-center justify-center">
                    {/* Google "G" Logo - Official SVG */}
                    <div className="w-5 h-5 mr-3 flex-shrink-0"
                         style={{background: 'white', borderRadius: '2px', padding: '1px'}}>
                        <svg width="18" height="18" viewBox="0 0 18 18">
                            <path
                                fill="#4285F4"
                                d="M16.51 8H8.98v3h4.3c-.18 1-.74 1.48-1.6 2.04v2.01h2.6a7.8 7.8 0 0 0 2.38-5.88c0-.57-.05-.66-.15-1.18z"
                            />
                            <path
                                fill="#34A853"
                                d="M8.98 17c2.16 0 3.97-.72 5.3-1.94l-2.6-2.04a4.8 4.8 0 0 1-7.18-2.51H1.83v2.07A8 8 0 0 0 8.98 17z"
                            />
                            <path
                                fill="#FBBC05"
                                d="M4.5 10.52a4.8 4.8 0 0 1 0-3.04V5.41H1.83a8 8 0 0 0 0 7.18l2.67-2.07z"
                            />
                            <path
                                fill="#EA4335"
                                d="M8.98 4.18c1.17 0 2.23.4 3.06 1.2l2.3-2.3A8 8 0 0 0 1.83 5.41L4.5 7.49a4.77 4.77 0 0 1 4.48-3.3z"
                            />
                        </svg>
                    </div>
                    <span>{buttonText}</span>
                </div>
            </button>
        </div>
    );
};

export default GoogleSignButton;