import '@testing-library/jest-dom';

// Add TextEncoder/TextDecoder polyfill for Node.js
if (typeof TextEncoder === 'undefined') {
    const {TextEncoder, TextDecoder} = require('util');
    global.TextEncoder = TextEncoder;
    global.TextDecoder = TextDecoder;
}

// Mock import.meta.env for tests
Object.defineProperty(globalThis, 'import', {
    value: {
        meta: {
            env: {
                VITE_API_URL: 'http://localhost:8080/api',
            },
        },
    },
    writable: true,
});

// Mock console methods to reduce noise in tests
global.console = {
    ...console,
    log: jest.fn(),
    error: jest.fn(),
    warn: jest.fn(),
    info: jest.fn(),
};

// Mock localStorage
const localStorageMock = {
    getItem: jest.fn(),
    setItem: jest.fn(),
    removeItem: jest.fn(),
    clear: jest.fn(),
};

Object.defineProperty(window, 'localStorage', {
    value: localStorageMock,
    writable: true,
});

// Reset localStorage mock before each test
beforeEach(() => {
    localStorageMock.getItem.mockClear();
    localStorageMock.setItem.mockClear();
    localStorageMock.removeItem.mockClear();
    localStorageMock.clear.mockClear();
});