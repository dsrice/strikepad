// Authentication related types

export interface LoginRequest {
  email: string;
  password: string;
}

export interface SignupRequest {
  email: string;
  password: string;
  display_name: string;
}

export interface GoogleSignupRequest {
  access_token: string;
}

export interface GoogleLoginRequest {
  access_token: string;
}

export interface UserInfo {
  id: number;
  email: string;
  display_name: string;
  email_verified: boolean;
  access_token?: string;
  refresh_token?: string;
  expires_at?: string;
}

export interface SignupResponse {
  id: number;
  email: string;
  display_name: string;
  email_verified: boolean;
  created_at: string;
  access_token: string;
  refresh_token: string;
  expires_at: string;
}

export interface ValidationError {
  field: string;
  tag: string;
  value: string;
  message: string;
}

export interface ErrorResponse {
  code: string;
  message: string;
  description: string;
  details?: ValidationError[];
}

export interface AuthContextType {
  user: UserInfo | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  signup: (email: string, password: string, display_name: string) => Promise<void>;
  googleSignup: (accessToken: string) => Promise<void>;
  googleLogin: (accessToken: string) => Promise<void>;
  logout: () => Promise<void>;
  error: string | null;
}