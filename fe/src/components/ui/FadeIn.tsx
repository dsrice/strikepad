import React from 'react';
import {clsx} from 'clsx';

interface FadeInProps {
    className?: string;
    children: React.ReactNode;
}

export function FadeIn({className, children}: FadeInProps) {
    return (
        <div className={clsx('animate-in fade-in duration-700', className)}>
            {children}
        </div>
    );
}

interface FadeInStaggerProps {
    className?: string;
    children: React.ReactNode;
    faster?: boolean;
}

export function FadeInStagger({className, children, faster = false}: FadeInStaggerProps) {
    return (
        <div className={clsx(
            'animate-in fade-in duration-700',
            faster ? 'stagger-children-300' : 'stagger-children-500',
            className
        )}>
            {children}
        </div>
    );
}