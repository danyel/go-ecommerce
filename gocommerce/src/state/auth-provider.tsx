import {useMemo, useState, type ReactNode} from 'react';
import ApiClient from '../domain/common/api-client.tsx';
import {AuthContext, type AuthContextValue, type LoginResponse} from './auth-context.ts';

export function AuthProvider({children}: {children: ReactNode}) {
    const [token, setToken] = useState<string | null>(() => localStorage.getItem('app_session_token'));
    const value = useMemo<AuthContextValue>(() => ({
        token,
        isAuthenticated: Boolean(token),
        loginWithGoogle: async (idToken: string) => {
            const response = await ApiClient.POST<LoginResponse, {id_token: string}>('/api/auth/v1/google', {id_token: idToken});
            localStorage.setItem('app_session_token', response.token);
            setToken(response.token);
        },
        logout: () => {
            localStorage.removeItem('app_session_token');
            setToken(null);
        },
    }), [token]);
    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
