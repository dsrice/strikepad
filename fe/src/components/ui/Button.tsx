import React from 'react';
import {Link} from 'react-router-dom';
import {clsx} from 'clsx';

interface BaseButtonProps {
    variant?: 'primary' | 'secondary' | 'outline';
    size?: 'sm' | 'md' | 'lg';
    className?: string;
    children: React.ReactNode;
    disabled?: boolean;
}

interface ButtonProps extends BaseButtonProps {
    onClick?: () => void;
    type?: 'button' | 'submit' | 'reset';
}

interface LinkButtonProps extends BaseButtonProps {
    href: string;
}

const buttonClasses = {
    base: 'inline-flex items-center justify-center rounded-lg font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed',
    variant: {
        primary: 'bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500',
        secondary: 'bg-gray-100 text-gray-900 hover:bg-gray-200 focus:ring-gray-500',
        outline: 'border border-gray-300 bg-white text-gray-700 hover:bg-gray-50 focus:ring-gray-500',
    },
    size: {
        sm: 'px-3 py-2 text-sm',
        md: 'px-4 py-2 text-base',
        lg: 'px-6 py-3 text-lg',
    },
};

export function Button({
                           variant = 'primary',
                           size = 'md',
                           className,
                           children,
                           onClick,
                           type = 'button',
                           disabled = false,
                       }: ButtonProps) {
    return (
        <button
            type={type}
            onClick={onClick}
            disabled={disabled}
            className={clsx(
                buttonClasses.base,
                buttonClasses.variant[variant],
                buttonClasses.size[size],
                className
            )}
        >
            {children}
        </button>
    );
}

export function LinkButton({
                               href,
                               variant = 'primary',
                               size = 'md',
                               className,
                               children,
                           }: LinkButtonProps) {
    return (
        <Link
            to={href}
            className={clsx(
                buttonClasses.base,
                buttonClasses.variant[variant],
                buttonClasses.size[size],
                className
            )}
        >
            {children}
        </Link>
    );
}