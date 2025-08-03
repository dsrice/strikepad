import axios, { AxiosResponse } from 'axios';
import { LoginRequest, SignupRequest, UserInfo, SignupResponse, ErrorResponse } from '../types/auth';

// Create axios instance with default config
const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/api',
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

// Add response interceptor for error handling
api.interceptors.response.use(
  (response) => {
    console.log(`✅ API Response: ${response.status} ${response.config.url}`);
    return response;
  },
  (error) => {
    console.error('❌ API Response Error:', error.response?.data || error.message);
    return Promise.reject(error);
  },
);

// Auth API functions
export const authAPI = {
  // Login user
  login: async (credentials: LoginRequest): Promise<UserInfo> => {
    try {
      const response: AxiosResponse<UserInfo> = await api.post('/auth/login', credentials);
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
      const response: AxiosResponse<SignupResponse> = await api.post('/auth/signup', userData);
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

  // Get current user profile (if we have token-based auth in the future)
  getProfile: async (): Promise<UserInfo> => {
    try {
      const response: AxiosResponse<UserInfo> = await api.get('/auth/profile');
      return response.data;
    } catch (error: any) {
      if (error.response?.data) {
        throw new Error(error.response.data.message || 'Failed to fetch profile');
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