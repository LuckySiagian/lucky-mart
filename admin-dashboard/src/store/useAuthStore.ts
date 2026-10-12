import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { Role, User } from '../types';

export type UserProfile = User;

export const STAFF_ROLES: Role[] = ['SUPER_ADMIN', 'STORE_MANAGER', 'WAREHOUSE_STAFF', 'DRIVER'];
export const MANAGER_ROLES: Role[] = ['SUPER_ADMIN', 'STORE_MANAGER'];

interface AuthState {
  user: UserProfile | null;
  accessToken: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
  setAuth: (user: UserProfile, accessToken: string, refreshToken: string) => void;
  setTokens: (accessToken: string, refreshToken: string, user?: UserProfile) => void;
  setUser: (user: UserProfile) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,
      setAuth: (user, accessToken, refreshToken) => set({ user, accessToken, refreshToken, isAuthenticated: true }),
      setTokens: (accessToken, refreshToken, user) =>
        set((state) => ({ accessToken, refreshToken, user: user ?? state.user })),
      setUser: (user) => set({ user }),
      logout: () => set({ user: null, accessToken: null, refreshToken: null, isAuthenticated: false }),
    }),
    { name: 'luckymart-admin-auth-storage' }
  )
);

export const hasRole = (user: UserProfile | null, roles: Role[]) => !!user && roles.includes(user.role);
