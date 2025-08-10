import '@testing-library/jest-dom';
import {TextEncoder as NodeTextEncoder, TextDecoder as NodeTextDecoder} from 'util';

// Add TextEncoder/TextDecoder polyfill for Node.js
if (typeof TextEncoder === 'undefined') {
    // Assign Node's TextEncoder/TextDecoder to global if not present
    // @ts-expect-error: Assign Node's TextEncoder to global in test env
    global.TextEncoder = NodeTextEncoder as unknown as typeof TextEncoder;
    // @ts-expect-error: Assign Node's TextDecoder to global in test env
    global.TextDecoder = NodeTextDecoder as unknown as typeof TextDecoder;
}

// Mock import.meta.env for tests
Object.defineProperty(globalThis, 'import', {
    value: {
        meta: {
            env: {
                VITE_API_URL: 'http://localhost:8080/api',
                VITE_GOOGLE_CLIENT_ID: 'test-client-id',
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