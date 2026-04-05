import { type ReactNode, createContext, useCallback, useMemo, useState } from "react";
import { getToken, removeToken, saveToken } from "@/utils/tokenStorage";

export interface AuthContextValue {
  token: string | null;
  isAuthenticated: boolean;
  login: (token: string) => void;
  logout: () => void;
}

export const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(() => getToken());

  const login = useCallback((newToken: string) => {
    saveToken(newToken);
    setToken(newToken);
  }, []);

  const logout = useCallback(() => {
    removeToken();
    setToken(null);
  }, []);

  const value = useMemo<AuthContextValue>(
    () => ({ token, isAuthenticated: token !== null, login, logout }),
    [token, login, logout],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
