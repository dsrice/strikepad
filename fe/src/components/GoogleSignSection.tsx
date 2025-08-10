import React, {useEffect, useState} from 'react';
import GoogleSignButton from './GoogleSignButton';

interface GoogleSignSectionProps {
    onSuccess: (accessToken: string) => void;
    onError: (error: string) => void;
    buttonText?: string;
    disabled?: boolean;
}

const GoogleSignSection: React.FC<GoogleSignSectionProps> = ({
                                                                 onSuccess,
                                                                 onError,
                                                                 buttonText = 'Googleでサインイン',
                                                                 disabled = false,
                                                             }) => {
    const [isGoogleAvailable, setIsGoogleAvailable] = useState(false);

    useEffect(() => {
        // Google APIの利用可能性をチェック
        const checkGoogleAvailability = () => {
            const hasGoogleScript = document.querySelector('script[src*="accounts.google.com"]');
            const hasClientId = !!(import.meta as any).env?.VITE_GOOGLE_CLIENT_ID;

            if (hasGoogleScript && hasClientId) {
                setIsGoogleAvailable(true);
            }
        };

        // 少し遅らせてチェック
        setTimeout(checkGoogleAvailability, 100);
    }, []);

    if (!isGoogleAvailable) {
        return null; // Google認証が利用できない場合は何も表示しない
    }

    return (
        <div className="space-y-4">
            <GoogleSignButton
                onSuccess={onSuccess}
                onError={onError}
                buttonText={buttonText}
                disabled={disabled}
            />

            <div className="relative">
                <div className="absolute inset-0 flex items-center">
                    <div className="w-full border-t border-gray-300"/>
                </div>
                <div className="relative flex justify-center text-sm">
                    <span className="px-2 bg-gray-50 text-gray-500">または</span>
                </div>
            </div>
        </div>
    );
};

export default GoogleSignSection;