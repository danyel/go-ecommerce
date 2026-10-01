import {createContext} from 'react';

export interface LoginResponse {
    token: string;
    user_id: string;
    provider: string;
}

export interface AuthContextValue {
    token: string | null;
    isAuthenticated: boolean;
    loginWithSSO: (idToken: string) => Promise<void>;
    logout: () => void;
}

export const AuthContext = createContext<AuthContextValue | undefined>(undefined);
