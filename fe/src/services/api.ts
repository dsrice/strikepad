import axios, { AxiosResponse } from 'axios';
import {
  LoginRequest,
  SignupRequest,
  GoogleSignupRequest,
  GoogleLoginRequest,
  UserInfo,
  SignupResponse,
    ErrorResponse,
    RefreshRequest,
    RefreshResponse,
    LoginResponse
} from '../types/auth';

// Flag to prevent infinite refresh loops
let isRefreshing = false;
let failedQueue: Array<{
    resolve: (value?: any) => void;
    reject: (reason?: any) => void;
}> = [];

const processQueue = (error: any, token: string | null = null) => {
    failedQueue.forEach(({resolve, reject}) => {
        if (error) {
            reject(error);
        } else {
            resolve(token);
        }
    });

    failedQueue = [];
};

// Create axios instance with default config
const api = axios.create({
  baseURL: (import.meta as any).env?.VITE_API_URL || 'http://localhost:8080/api',
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 10000, // 10 seconds timeout
});

// Add request interceptor for logging
api.interceptors.request.use(
  (config) => {
    console.log(`🚀 API Request: ${config.method?.toUpperCase()} ${config.url}`);
    return config;
  },
  (error) => {
    console.error('❌ API Request Error:', error);
    return Promise.reject(error);
  },
);

// Add response interceptor for error handling and automatic token refresh
api.interceptors.response.use(
  (response) => {
    console.log(`✅ API Response: ${response.status} ${response.config.url}`);
    return response;
  },
    async (error) => {
        const originalRequest = error.config;

        // Check if it's a 401 error and not a login/signup/refresh request
        if (error.response?.status === 401 &&
            !originalRequest._retry &&
            !originalRequest.url?.includes('/auth/login') &&
            !originalRequest.url?.includes('/auth/signup') &&
            !originalRequest.url?.includes('/auth/refresh') &&
            !originalRequest.url?.includes('/auth/google/login') &&
            !originalRequest.url?.includes('/auth/google/signup')) {

            if (isRefreshing) {
                // If already refreshing, queue this request
                return new Promise((resolve, reject) => {
                    failedQueue.push({resolve, reject});
                }).then(() => {
                    // Retry with new token
                    const token = localStorage.getItem('access_token');
                    if (token) {
                        originalRequest.headers.Authorization = `Bearer ${token}`;
                    }
                    return api(originalRequest);
                }).catch(err => {
                    return Promise.reject(err);
                });
            }

            originalRequest._retry = true;
            isRefreshing = true;

            try {
                const accessToken = localStorage.getItem('access_token');
                const refreshToken = localStorage.getItem('refresh_token');

                if (!accessToken || !refreshToken) {
                    throw new Error('No tokens available');
                }

                // Try to refresh the token
                const refreshResponse = await authAPI.refresh({
                    access_token: accessToken,
                    refresh_token: refreshToken
                });

                // Store new tokens
                localStorage.setItem('access_token', refreshResponse.access_token);
                localStorage.setItem('refresh_token', refreshResponse.refresh_token);

                // Update the authorization header for the failed request
                originalRequest.headers.Authorization = `Bearer ${refreshResponse.access_token}`;

                // Process queued requests
                processQueue(null, refreshResponse.access_token);

                // Retry the original request
                return api(originalRequest);

            } catch (refreshError) {
                console.error('Token refresh failed:', refreshError);

                // Clear tokens and redirect to login
                localStorage.removeItem('access_token');
                localStorage.removeItem('refresh_token');

                // Process queued requests with error
                processQueue(refreshError, null);

                // Redirect to login page
                if (typeof window !== 'undefined') {
                    window.location.href = '/login';
                }

                return Promise.reject(refreshError);
            } finally {
                isRefreshing = false;
            }
        }

    console.error('❌ API Response Error:', error.response?.data || error.message);
    return Promise.reject(error);
  },
);

// Helper function to get Authorization header with Bearer token
const getAuthHeaders = () => {
  const token = localStorage.getItem('access_token');
  return token ? {Authorization: `Bearer ${token}`} : {};
};

// Auth API functions
export const authAPI = {
  // Login user
    login: async (credentials: LoginRequest): Promise<LoginResponse> => {
    try {
        const response: AxiosResponse<LoginResponse> = await api.post('/auth/login', credentials);
      return response.data;
    } catch (error: any) {
      if (error.response?.data) {
        throw new Error(error.response.data.message || 'Login failed');
      }
      throw new Error('Network error occurred');
    }
  },

  // Register new user
  signup: async (userData: SignupRequest): Promise<SignupResponse> => {
    try {
      const requestData = {
        email: userData.email,
        password: userData.password,
        display_name: userData.display_name
      };
      const response: AxiosResponse<SignupResponse> = await api.post('/auth/signup', requestData);
      return response.data;
    } catch (error: any) {
      if (error.response?.data) {
        const errorData: ErrorResponse = error.response.data;
        if (errorData.details && errorData.details.length > 0) {
          // Format validation errors
          const validationMessages = errorData.details.map(detail => detail.message).join(', ');
          throw new Error(validationMessages);
        }
        throw new Error(errorData.message || 'Signup failed');
      }
      throw new Error('Network error occurred');
    }
  },

  // Google OAuth signup
  googleSignup: async (googleData: GoogleSignupRequest): Promise<SignupResponse> => {
    try {
      const response: AxiosResponse<SignupResponse> = await api.post('/auth/google/signup', googleData);
      return response.data;
    } catch (error: any) {
      if (error.response?.data) {
        const errorData: ErrorResponse = error.response.data;
        throw new Error(errorData.message || 'Google signup failed');
      }
      throw new Error('Network error occurred');
    }
  },

  // Google OAuth login
    googleLogin: async (googleData: GoogleLoginRequest): Promise<LoginResponse> => {
    try {
        const response: AxiosResponse<LoginResponse> = await api.post('/auth/google/login', googleData);
      return response.data;
    } catch (error: any) {
      if (error.response?.data) {
        throw new Error(error.response.data.message || 'Google login failed');
      }
      throw new Error('Network error occurred');
    }
  },

    // Get current user profile
  getProfile: async (): Promise<UserInfo> => {
    try {
        const response: AxiosResponse<UserInfo> = await api.get('/user/me', {
        headers: getAuthHeaders(),
      });
      return response.data;
    } catch (error: any) {
      if (error.response?.data) {
        throw new Error(error.response.data.message || 'Failed to fetch profile');
      }
      throw new Error('Network error occurred');
    }
  },

    // Refresh tokens
    refresh: async (refreshRequest: RefreshRequest): Promise<RefreshResponse> => {
        try {
            const response: AxiosResponse<RefreshResponse> = await api.post('/auth/refresh', refreshRequest);
            return response.data;
        } catch (error: any) {
            if (error.response?.data) {
                throw new Error(error.response.data.message || 'Token refresh failed');
            }
            throw new Error('Network error occurred');
        }
    },

  // Logout user
  logout: async (): Promise<{ message: string }> => {
    try {
      const response: AxiosResponse<{ message: string }> = await api.post('/auth/logout', {}, {
        headers: getAuthHeaders(),
      });
      return response.data;
    } catch (error: any) {
      if (error.response?.data) {
        throw new Error(error.response.data.message || 'Logout failed');
      }
      throw new Error('Network error occurred');
    }
  },
};

// Health check
export const healthAPI = {
  check: async (): Promise<{ status: string; message: string }> => {
    try {
      const response = await api.get('/health');
      return response.data;
    } catch {
      throw new Error('Health check failed');
    }
  },
};

export default api;