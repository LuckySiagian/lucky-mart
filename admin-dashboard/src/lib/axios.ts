import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios';
import { useAuthStore } from '../store/useAuthStore';

const baseURL = import.meta.env.VITE_API_URL || '/api/v1';

export const api = axios.create({
  baseURL,
  headers: { 'Content-Type': 'application/json' },
  timeout: 20000,
});

api.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const accessToken = useAuthStore.getState().accessToken;
  if (accessToken && config.headers) {
    config.headers.Authorization = `Bearer ${accessToken}`;
  }
  return config;
});

// One refresh at a time: requests that fail while a refresh is running wait
// for it instead of each spending (and revoking) the refresh token.
let refreshing: Promise<string> | null = null;

function refreshAccessToken(): Promise<string> {
  if (!refreshing) {
    const refreshToken = useAuthStore.getState().refreshToken;
    refreshing = axios
      .post(`${baseURL}/auth/refresh`, { refresh_token: refreshToken })
      .then((res) => {
        const { access_token, refresh_token, user } = res.data.data;
        useAuthStore.getState().setTokens(access_token, refresh_token, user);
        return access_token as string;
      })
      .finally(() => {
        refreshing = null;
      });
  }
  return refreshing;
}

function forceLogout() {
  useAuthStore.getState().logout();
  if (window.location.pathname !== '/login') {
    window.location.href = '/login';
  }
}

api.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const original = error.config as InternalAxiosRequestConfig & { _retry?: boolean };
    const isAuthCall = original?.url?.includes('/auth/login') || original?.url?.includes('/auth/refresh');

    if (error.response?.status === 401 && original && !original._retry && !isAuthCall) {
      original._retry = true;
      if (!useAuthStore.getState().refreshToken) {
        forceLogout();
        return Promise.reject(error);
      }
      try {
        const token = await refreshAccessToken();
        original.headers.Authorization = `Bearer ${token}`;
        return api(original);
      } catch (refreshErr) {
        forceLogout();
        return Promise.reject(refreshErr);
      }
    }

    return Promise.reject(error);
  }
);
