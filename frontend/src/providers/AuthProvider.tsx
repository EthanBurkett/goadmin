import { createContext } from "react";
import type { ReactNode } from "react";
import { useAuth } from "@/hooks/useAuth";
import type { User } from "@/hooks/useAuth";

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  error: Error | null;
  isAuthenticated: boolean;
  login: (credentials: { username: string; password: string }) => Promise<User>;
  register: (credentials: {
    username: string;
    password: string;
  }) => Promise<User>;
  logout: () => Promise<void>;
  isLoggingIn: boolean;
  isRegistering: boolean;
  isLoggingOut: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const auth = useAuth();

  return <AuthContext.Provider value={auth}>{children}</AuthContext.Provider>;
}

export { AuthContext };
