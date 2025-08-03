import axios from 'axios';
import {authAPI, healthAPI} from './api';

// Mock axios
jest.mock('axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

// Mock the axios instance
const mockAxiosInstance = {
    post: jest.fn(),
    get: jest.fn(),
    interceptors: {
        request: {use: jest.fn()},
        response: {use: jest.fn()},
    },
};

mockedAxios.create.mockReturnValue(mockAxiosInstance as any);

describe('authAPI', () => {
    beforeEach(() => {
        jest.clearAllMocks();
    });

    describe('login', () => {
        it('successfully logs in user', async () => {
            const mockResponse = {
                data: {
                    id: 1,
                    email: 'test@example.com',
                    username: 'testuser',
                    createdAt: '2023-01-01T00:00:00Z',
                    updatedAt: '2023-01-01T00:00:00Z',
                },
            };

            mockAxiosInstance.post.mockResolvedValue(mockResponse);

            const credentials = {
                email: 'test@example.com',
                password: 'password123',
            };

            const result = await authAPI.login(credentials);

            expect(mockAxiosInstance.post).toHaveBeenCalledWith('/auth/login', credentials);
            expect(result).toEqual(mockResponse.data);
        });

        it('handles login error with error response', async () => {
            const errorResponse = {
                response: {
                    data: {
                        message: 'Invalid credentials',
                    },
                },
            };

            mockAxiosInstance.post.mockRejectedValue(errorResponse);

            const credentials = {
                email: 'test@example.com',
                password: 'wrongpassword',
            };

            await expect(authAPI.login(credentials)).rejects.toThrow('Invalid credentials');
        });

        it('handles network error', async () => {
            mockAxiosInstance.post.mockRejectedValue(new Error('Network Error'));

            const credentials = {
                email: 'test@example.com',
                password: 'password123',
            };

            await expect(authAPI.login(credentials)).rejects.toThrow('Network error occurred');
        });
    });

    describe('signup', () => {
        it('successfully registers user', async () => {
            const mockResponse = {
                data: {
                    message: 'User created successfully',
                    user: {
                        id: 1,
                        email: 'test@example.com',
                        username: 'testuser',
                    },
                },
            };

            mockAxiosInstance.post.mockResolvedValue(mockResponse);

            const userData = {
                email: 'test@example.com',
                username: 'testuser',
                password: 'Password123!',
            };

            const result = await authAPI.signup(userData);

            expect(mockAxiosInstance.post).toHaveBeenCalledWith('/auth/signup', userData);
            expect(result).toEqual(mockResponse.data);
        });

        it('handles signup validation errors', async () => {
            const errorResponse = {
                response: {
                    data: {
                        message: 'Validation failed',
                        details: [
                            {field: 'email', message: 'Email is required'},
                            {field: 'password', message: 'Password must contain uppercase letter'},
                        ],
                    },
                },
            };

            mockAxiosInstance.post.mockRejectedValue(errorResponse);

            const userData = {
                email: '',
                username: 'testuser',
                password: 'weak',
            };

            await expect(authAPI.signup(userData)).rejects.toThrow(
                'Email is required, Password must contain uppercase letter'
            );
        });

        it('handles signup error without details', async () => {
            const errorResponse = {
                response: {
                    data: {
                        message: 'Email already exists',
                    },
                },
            };

            mockAxiosInstance.post.mockRejectedValue(errorResponse);

            const userData = {
                email: 'existing@example.com',
                username: 'testuser',
                password: 'Password123!',
            };

            await expect(authAPI.signup(userData)).rejects.toThrow('Email already exists');
        });
    });

    describe('getProfile', () => {
        it('successfully gets user profile', async () => {
            const mockResponse = {
                data: {
                    id: 1,
                    email: 'test@example.com',
                    username: 'testuser',
                    createdAt: '2023-01-01T00:00:00Z',
                    updatedAt: '2023-01-01T00:00:00Z',
                },
            };

            mockAxiosInstance.get.mockResolvedValue(mockResponse);

            const result = await authAPI.getProfile();

            expect(mockAxiosInstance.get).toHaveBeenCalledWith('/auth/profile');
            expect(result).toEqual(mockResponse.data);
        });

        it('handles profile fetch error', async () => {
            const errorResponse = {
                response: {
                    data: {
                        message: 'Unauthorized',
                    },
                },
            };

            mockAxiosInstance.get.mockRejectedValue(errorResponse);

            await expect(authAPI.getProfile()).rejects.toThrow('Unauthorized');
        });
    });
});

describe('healthAPI', () => {
    beforeEach(() => {
        jest.clearAllMocks();
    });

    describe('check', () => {
        it('successfully performs health check', async () => {
            const mockResponse = {
                data: {
                    status: 'OK',
                    message: 'Server is healthy',
                },
            };

            mockAxiosInstance.get.mockResolvedValue(mockResponse);

            const result = await healthAPI.check();

            expect(mockAxiosInstance.get).toHaveBeenCalledWith('/health');
            expect(result).toEqual(mockResponse.data);
        });

        it('handles health check failure', async () => {
            mockAxiosInstance.get.mockRejectedValue(new Error('Server Error'));

            await expect(healthAPI.check()).rejects.toThrow('Health check failed');
        });
    });
});

describe('axios instance configuration', () => {
    it('creates axios instance with correct config', () => {
        expect(mockedAxios.create).toHaveBeenCalledWith({
            baseURL: 'http://localhost:8080/api',
            headers: {
                'Content-Type': 'application/json',
            },
            timeout: 10000,
        });
    });

    it('sets up request and response interceptors', () => {
        expect(mockAxiosInstance.interceptors.request.use).toHaveBeenCalled();
        expect(mockAxiosInstance.interceptors.response.use).toHaveBeenCalled();
    });
});