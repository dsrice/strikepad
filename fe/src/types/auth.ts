// Authentication related types

export interface LoginRequest {
  email: string;
  password: string;
}

export interface SignupRequest {
  email: string;
  password: string;
  displayName: string;
}

export interface UserInfo {
  id: number;
  email: string;
  displayName: string;
  emailVerified: boolean;
}

export interface SignupResponse {
  id: number;
  email: string;
  displayName: string;
  emailVerified: boolean;
  createdAt: string;
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
  signup: (email: string, password: string, displayName: string) => Promise<void>;
  logout: () => void;
  error: string | null;
}