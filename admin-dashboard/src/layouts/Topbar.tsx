import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Menu, Sun, Moon, Bell, LogOut, User as UserIcon } from 'lucide-react';
import { toast } from 'sonner';
import { useAuthStore } from '../store/useAuthStore';
import { useThemeStore } from '../store/useThemeStore';
import { authService, notificationService } from '../services/apiServices';
import { ROLE_LABELS } from '../lib/utils';

export const Topbar: React.FC<{ onMenuClick: () => void }> = ({ onMenuClick }) => {
  const { user, logout } = useAuthStore();
  const { theme, toggleTheme } = useThemeStore();
  const navigate = useNavigate();
  const [showProfileMenu, setShowProfileMenu] = useState(false);
  const [unread, setUnread] = useState(0);

  useEffect(() => {
    let alive = true;
    const refresh = () =>
      notificationService
        .unreadCount()
        .then((r) => alive && setUnread(r.unread))
        .catch(() => undefined);
    refresh();
    const timer = setInterval(refresh, 30000);
    window.addEventListener('notifications:changed', refresh);
    return () => {
      alive = false;
      clearInterval(timer);
      window.removeEventListener('notifications:changed', refresh);
    };
  }, []);

  const handleLogout = async () => {
    try {
      await authService.logout();
    } catch {
      // Logging out locally is enough if the server is unreachable.
    } finally {
      logout();
      toast.success('Berhasil keluar.');
      navigate('/login', { replace: true });
    }
  };

  return (
    <header className="sticky top-0 z-30 flex h-16 w-full items-center justify-between border-b border-slate-200 bg-white/80 px-4 backdrop-blur-md dark:border-slate-800 dark:bg-slate-900/80 sm:px-6">
      <button onClick={onMenuClick} className="rounded-lg p-2 text-slate-600 hover:bg-slate-100 dark:text-slate-400 dark:hover:bg-slate-800 lg:hidden" aria-label="Menu">
        <Menu className="h-6 w-6" />
      </button>
      <div className="hidden lg:block" />

      <div className="flex items-center gap-2">
        <button
          onClick={toggleTheme}
          className="rounded-lg p-2.5 text-slate-600 transition-colors hover:bg-slate-100 dark:text-slate-400 dark:hover:bg-slate-800"
          title="Mode terang / gelap"
        >
          {theme === 'light' ? <Moon className="h-5 w-5" /> : <Sun className="h-5 w-5 text-amber-400" />}
        </button>

        <button
          onClick={() => navigate('/notifications')}
          className="relative rounded-lg p-2.5 text-slate-600 transition-colors hover:bg-slate-100 dark:text-slate-400 dark:hover:bg-slate-800"
          title="Notifikasi"
        >
          <Bell className="h-5 w-5" />
          {unread > 0 && (
            <span className="absolute right-1 top-1 flex h-4 min-w-[1rem] items-center justify-center rounded-full bg-rose-500 px-1 text-[10px] font-bold text-white">
              {unread > 99 ? '99+' : unread}
            </span>
          )}
        </button>

        <div className="relative">
          <button
            onClick={() => setShowProfileMenu(!showProfileMenu)}
            className="flex items-center gap-3 rounded-lg p-1.5 transition-colors hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <div className="flex h-9 w-9 items-center justify-center rounded-full bg-brand-600 font-bold text-white shadow-md">
              {user?.full_name?.charAt(0).toUpperCase() || 'A'}
            </div>
            <div className="hidden text-left md:block">
              <div className="text-sm font-semibold leading-tight text-slate-900 dark:text-white">{user?.full_name}</div>
              <div className="text-xs text-slate-500 dark:text-slate-400">{user ? ROLE_LABELS[user.role] : ''}</div>
            </div>
          </button>

          {showProfileMenu && (
            <>
              <div className="fixed inset-0 z-0" onClick={() => setShowProfileMenu(false)} />
              <div className="absolute right-0 z-10 mt-2 w-48 rounded-xl bg-white p-1 shadow-xl ring-1 ring-slate-200 dark:bg-slate-900 dark:ring-slate-800">
                <button
                  onClick={() => {
                    setShowProfileMenu(false);
                    navigate('/settings');
                  }}
                  className="flex w-full items-center gap-2 rounded-lg px-3 py-2 text-sm text-slate-700 hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-800"
                >
                  <UserIcon className="h-4 w-4" />
                  <span>Profil & Pengaturan</span>
                </button>
                <div className="my-1 border-t border-slate-200 dark:border-slate-800" />
                <button
                  onClick={handleLogout}
                  className="flex w-full items-center gap-2 rounded-lg px-3 py-2 text-sm text-rose-600 hover:bg-rose-50 dark:hover:bg-rose-950/30"
                >
                  <LogOut className="h-4 w-4" />
                  <span>Keluar</span>
                </button>
              </div>
            </>
          )}
        </div>
      </div>
    </header>
  );
};
