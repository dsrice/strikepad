import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { AuthContextType, UserInfo } from '../types/auth';
import { authAPI } from '../services/api';

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<UserInfo | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const isAuthenticated = user !== null;

  // Check if user is already logged in on app start
  useEffect(() => {
    const checkAuth = async () => {
      try {
          const token = localStorage.getItem('access_token');
          if (token) {
              // Fetch user info from API
              const userInfo = await authAPI.getProfile();
              setUser(userInfo);
        }
      } catch (error) {
          console.error('Failed to load user from API:', error);
          // Clear tokens if API call fails
          localStorage.removeItem('access_token');
          localStorage.removeItem('refresh_token');
      } finally {
        setIsLoading(false);
      }
    };

    checkAuth();
  }, []);

  const login = async (email: string, password: string): Promise<void> => {
    try {
      setIsLoading(true);
      setError(null);

      const loginResponse = await authAPI.login({email, password});

      // Store tokens separately
        if (loginResponse.access_token) {
            localStorage.setItem('access_token', loginResponse.access_token);
        }
        if (loginResponse.refresh_token) {
            localStorage.setItem('refresh_token', loginResponse.refresh_token);
        }

        // Fetch user info from API after login
        const userInfo = await authAPI.getProfile();
      setUser(userInfo);
    } catch (error: any) {
      setError(error.message || 'Login failed');
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const signup = async (email: string, password: string, displayName: string): Promise<void> => {
    try {
      setIsLoading(true);
      setError(null);

      const response = await authAPI.signup({email, password, display_name: displayName});
      
      // Store tokens separately
        if (response.access_token) {
            localStorage.setItem('access_token', response.access_token);
        }
        if (response.refresh_token) {
            localStorage.setItem('refresh_token', response.refresh_token);
        }

        // Fetch user info from API after signup
        const userInfo = await authAPI.getProfile();
      setUser(userInfo);
    } catch (error: any) {
      setError(error.message || 'Signup failed');
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const googleSignup = async (accessToken: string): Promise<void> => {
    try {
      setIsLoading(true);
      setError(null);

      const response = await authAPI.googleSignup({access_token: accessToken});

        // Store tokens if available
        if (response.access_token) {
            localStorage.setItem('access_token', response.access_token);
        }
        if (response.refresh_token) {
            localStorage.setItem('refresh_token', response.refresh_token);
        }

        // Fetch user info from API after signup
        const userInfo = await authAPI.getProfile();
      setUser(userInfo);
    } catch (error: any) {
      setError(error.message || 'Google signup failed');
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const googleLogin = async (accessToken: string): Promise<void> => {
    try {
      setIsLoading(true);
      setError(null);

        const loginResponse = await authAPI.googleLogin({access_token: accessToken});

        // Store tokens if available
        if (loginResponse.access_token) {
            localStorage.setItem('access_token', loginResponse.access_token);
        }
        if (loginResponse.refresh_token) {
            localStorage.setItem('refresh_token', loginResponse.refresh_token);
        }

        // Fetch user info from API after login
        const userInfo = await authAPI.getProfile();
      setUser(userInfo);
    } catch (error: any) {
      setError(error.message || 'Google login failed');
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const logout = async (): Promise<void> => {
    try {
      // Call logout API to invalidate session on server
      await authAPI.logout();
    } catch (error: any) {
      // Even if API call fails, we still want to clear local storage
      console.error('Logout API call failed:', error.message);
    } finally {
      // Clear local state and storage regardless of API call result
      setUser(null);
      setError(null);
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
    }
  };

  const value: AuthContextType = {
    user,
    isAuthenticated,
    isLoading,
    login,
    signup,
    googleSignup,
    googleLogin,
    logout,
    error,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};