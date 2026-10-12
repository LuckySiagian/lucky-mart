import React, { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Bell, Heart, LogOut, Moon, Package, Search, ShoppingBag, Sun, User } from 'lucide-react';
import { toast } from 'sonner';
import { useAuthStore } from '../../store/useAuthStore';
import { useCartStore } from '../../store/useCartStore';
import { useWishlistStore } from '../../store/useWishlistStore';
import { useThemeStore } from '../../store/useThemeStore';
import { authService, notificationService } from '../../services/apiServices';
import { useStoreInfo } from '../../hooks/useStoreInfo';

const CountBadge: React.FC<{ n: number; tone?: string }> = ({ n, tone = 'bg-primary-600' }) =>
  n > 0 ? (
    <span className={`absolute -right-1 -top-1 flex h-5 min-w-[1.25rem] items-center justify-center rounded-full ${tone} px-1 text-[10px] font-bold text-white shadow-md ring-2 ring-white dark:ring-slate-900`}>
      {n > 99 ? '99+' : n}
    </span>
  ) : null;

export const Navbar: React.FC = () => {
  const navigate = useNavigate();
  const store = useStoreInfo();
  const { user, isAuthenticated, logout } = useAuthStore();
  const { totalItems, refresh: refreshCart, reset: resetCart } = useCartStore();
  const wishlistCount = useWishlistStore((s) => s.wishlistIds.length);
  const { theme, toggleTheme } = useThemeStore();
  const [searchQuery, setSearchQuery] = useState('');
  const [menuOpen, setMenuOpen] = useState(false);
  const [unread, setUnread] = useState(0);

  // Keep the cart badge and notification count in sync with the server.
  useEffect(() => {
    if (!isAuthenticated) {
      resetCart();
      setUnread(0);
      return;
    }
    refreshCart().catch(() => undefined);
    const loadUnread = () =>
      notificationService
        .unreadCount()
        .then((r) => setUnread(r.unread))
        .catch(() => undefined);
    loadUnread();
    const timer = setInterval(loadUnread, 60000);
    window.addEventListener('notifications:changed', loadUnread);
    return () => {
      clearInterval(timer);
      window.removeEventListener('notifications:changed', loadUnread);
    };
  }, [isAuthenticated, refreshCart, resetCart]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    navigate(searchQuery.trim() ? `/products?search=${encodeURIComponent(searchQuery.trim())}` : '/products');
  };

  const handleLogout = async () => {
    setMenuOpen(false);
    try {
      await authService.logout();
    } catch {
      // Local logout is enough when offline.
    }
    logout();
    resetCart();
    toast.success('Berhasil keluar.');
    navigate('/');
  };

  const iconBtn = 'relative rounded-xl p-2 text-slate-600 transition-colors hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-800';

  return (
    <header className="sticky top-0 z-50 border-b border-slate-200/80 bg-white/90 shadow-sm backdrop-blur-md dark:border-slate-800/80 dark:bg-slate-900/90">
      <div className="hidden bg-primary-700 px-4 py-1.5 text-center text-xs font-medium text-white sm:block">
        {store?.announcement || 'Pesan dari rumah — diantar atau ambil sendiri di toko. Mitra warung dapat harga grosir.'}
      </div>

      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="flex h-16 items-center justify-between gap-3 sm:h-20">
          <Link to="/" className="flex flex-shrink-0 items-center gap-3">
            <img src="/logo.png" alt="Lucky Mart" className="h-10 w-auto object-contain sm:h-12" />
            <div className="hidden sm:block">
              <span className="block text-lg font-black leading-none text-slate-900 dark:text-white">
                LUCKY<span className="text-primary-600">MART</span>
              </span>
              <span className="block text-[10px] font-semibold uppercase tracking-wider text-slate-500">Grosir & Eceran</span>
            </div>
          </Link>

          <form onSubmit={handleSearch} className="hidden max-w-xl flex-1 md:block">
            <div className="relative">
              <input
                type="search"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="Cari beras, minyak goreng, gula, telur..."
                className="w-full rounded-2xl border-0 bg-slate-100 py-2.5 pl-11 pr-20 text-sm text-slate-900 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-primary-500 dark:bg-slate-800 dark:text-white"
              />
              <Search className="absolute left-4 top-3 h-4 w-4 text-slate-400" />
              <button type="submit" className="absolute right-1.5 top-1.5 rounded-xl bg-primary-600 px-4 py-1 text-xs font-bold text-white hover:bg-primary-700">
                Cari
              </button>
            </div>
          </form>

          <div className="flex items-center gap-1 sm:gap-2">
            <button onClick={() => navigate('/products')} className={`${iconBtn} md:hidden`} aria-label="Cari">
              <Search className="h-5 w-5" />
            </button>
            <button onClick={toggleTheme} className={`${iconBtn} hidden sm:block`} aria-label="Mode terang / gelap">
              {theme === 'light' ? <Moon className="h-5 w-5" /> : <Sun className="h-5 w-5 text-amber-400" />}
            </button>
            {isAuthenticated && (
              <>
                <Link to="/orders" className={`${iconBtn} flex items-center gap-1.5`} aria-label="Pesanan Saya">
                  <Package className="h-5 w-5" />
                  <span className="hidden text-sm font-semibold sm:inline">Pesanan</span>
                </Link>
                <Link to="/notifications" className={iconBtn} aria-label="Notifikasi">
                  <Bell className="h-5 w-5" />
                  <CountBadge n={unread} tone="bg-rose-500" />
                </Link>
              </>
            )}
            <Link to="/cart" className={iconBtn} aria-label="Keranjang">
              <ShoppingBag className="h-5 w-5" />
              <CountBadge n={totalItems} />
            </Link>

            {isAuthenticated ? (
              <div className="relative">
                <button onClick={() => setMenuOpen(!menuOpen)} className="flex items-center gap-2 rounded-xl p-1.5 hover:bg-slate-100 dark:hover:bg-slate-800">
                  <div className="flex h-8 w-8 items-center justify-center rounded-full bg-primary-600 text-sm font-bold text-white">
                    {user?.full_name?.charAt(0).toUpperCase() || 'U'}
                  </div>
                  <span className="hidden text-sm font-bold text-slate-800 dark:text-slate-200 lg:block">{user?.full_name}</span>
                </button>

                {menuOpen && (
                  <>
                    <div className="fixed inset-0 z-0" onClick={() => setMenuOpen(false)} />
                    <div className="absolute right-0 z-10 mt-2 w-52 rounded-2xl bg-white p-2 shadow-2xl ring-1 ring-slate-200 dark:bg-slate-900 dark:ring-slate-800">
                      <Link to="/profile" onClick={() => setMenuOpen(false)} className="flex items-center gap-2 rounded-xl px-3 py-2 text-sm text-slate-700 hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-800">
                        <User className="h-4 w-4" /> Akun Saya
                      </Link>
                      <Link to="/orders" onClick={() => setMenuOpen(false)} className="flex items-center gap-2 rounded-xl px-3 py-2 text-sm text-slate-700 hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-800">
                        <Package className="h-4 w-4" /> Pesanan Saya
                      </Link>
                      <Link to="/wishlist" onClick={() => setMenuOpen(false)} className="flex items-center gap-2 rounded-xl px-3 py-2 text-sm text-slate-700 hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-800">
                        <Heart className="h-4 w-4" /> Favorit {wishlistCount > 0 && <span className="text-xs text-slate-400">({wishlistCount})</span>}
                      </Link>
                      <button onClick={toggleTheme} className="flex w-full items-center gap-2 rounded-xl px-3 py-2 text-sm text-slate-700 hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-800 sm:hidden">
                        {theme === 'light' ? <Moon className="h-4 w-4" /> : <Sun className="h-4 w-4" />} Mode {theme === 'light' ? 'gelap' : 'terang'}
                      </button>
                      <div className="my-1 border-t border-slate-200 dark:border-slate-800" />
                      <button onClick={handleLogout} className="flex w-full items-center gap-2 rounded-xl px-3 py-2 text-sm text-rose-600 hover:bg-rose-50 dark:hover:bg-rose-950/30">
                        <LogOut className="h-4 w-4" /> Keluar
                      </button>
                    </div>
                  </>
                )}
              </div>
            ) : (
              <Link to="/login" className="rounded-xl bg-primary-600 px-4 py-2 text-sm font-semibold text-white shadow-md shadow-primary-600/20 hover:bg-primary-700">
                Masuk
              </Link>
            )}
          </div>
        </div>
      </div>
    </header>
  );
};
