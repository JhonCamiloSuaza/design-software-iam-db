import React, { createContext, useCallback, useEffect, useMemo, useState } from 'react';
import { authApi } from '../services/api';

export const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [token, setToken] = useState(() => localStorage.getItem('iam_token') || '');
  const [profile, setProfile] = useState(null);
  const [summary, setSummary] = useState(null);

  const logout = useCallback(() => {
    localStorage.removeItem('iam_token');
    setToken('');
    setProfile(null);
    setSummary(null);
  }, []);

  const loadSession = useCallback(async (sessionToken) => {
    const [profileData, summaryData] = await Promise.all([
      authApi.profile(sessionToken),
      authApi.summary(sessionToken)
    ]);
    setProfile(profileData);
    setSummary(summaryData);
  }, []);

  const saveSession = useCallback(async (data) => {
    localStorage.setItem('iam_token', data.token);
    setToken(data.token);
    await loadSession(data.token);
    return data.message;
  }, [loadSession]);

  const login = useCallback(async (credentials) => saveSession(await authApi.login(credentials)), [saveSession]);
  const register = useCallback(async (user) => saveSession(await authApi.register(user)), [saveSession]);

  useEffect(() => {
    if (token && !profile) loadSession(token).catch(logout);
  }, [token, profile, loadSession, logout]);

  const canManageUsers = profile?.permissions?.some((permission) => permission.feature === 'Gestionar usuarios') || false;
  const value = useMemo(() => ({ token, profile, summary, canManageUsers, login, register, logout, loadSession }), [token, profile, summary, canManageUsers, login, register, logout, loadSession]);

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
