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
        const storedUser = localStorage.getItem('user');
        if (storedUser) {
          setUser(JSON.parse(storedUser));
        }
      } catch (error) {
        console.error('Failed to load user from storage:', error);
        localStorage.removeItem('user');
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
      
      const userInfo = await authAPI.login({ email, password });
      
      setUser(userInfo);
      localStorage.setItem('user', JSON.stringify(userInfo));
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
      
      const response = await authAPI.signup({ email, password, displayName });
      
      // Convert signup response to UserInfo format
      const userInfo: UserInfo = {
        id: response.id,
        email: response.email,
        displayName: response.displayName,
        emailVerified: response.emailVerified,
      };
      
      setUser(userInfo);
      localStorage.setItem('user', JSON.stringify(userInfo));
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

      // Convert signup response to UserInfo format
      const userInfo: UserInfo = {
        id: response.id,
        email: response.email,
        displayName: response.displayName,
        emailVerified: response.emailVerified,
      };

      setUser(userInfo);
      localStorage.setItem('user', JSON.stringify(userInfo));
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

      const userInfo = await authAPI.googleLogin({access_token: accessToken});

      setUser(userInfo);
      localStorage.setItem('user', JSON.stringify(userInfo));
    } catch (error: any) {
      setError(error.message || 'Google login failed');
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const logout = (): void => {
    setUser(null);
    setError(null);
    localStorage.removeItem('user');
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