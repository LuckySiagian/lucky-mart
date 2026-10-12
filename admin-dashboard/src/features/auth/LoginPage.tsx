import React, { useState } from 'react';
import { Navigate, useNavigate } from 'react-router-dom';
import { Lock, Phone, ArrowRight } from 'lucide-react';
import { toast } from 'sonner';
import { useAuthStore, STAFF_ROLES } from '../../store/useAuthStore';
import { authService } from '../../services/apiServices';
import { getErrorMessage } from '../../lib/utils';

export const LoginPage: React.FC = () => {
  const [phoneNumber, setPhoneNumber] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const { isAuthenticated, user, setAuth } = useAuthStore();
  const navigate = useNavigate();

  if (isAuthenticated && user && STAFF_ROLES.includes(user.role)) {
    return <Navigate to={user.role === 'DRIVER' ? '/shipping' : '/dashboard'} replace />;
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      const { user, access_token, refresh_token } = await authService.login({ phone_number: phoneNumber.trim(), password });
      if (!STAFF_ROLES.includes(user.role)) {
        toast.error('Akun ini adalah akun pelanggan. Silakan masuk lewat toko online.');
        return;
      }
      setAuth(user, access_token, refresh_token);
      toast.success(`Selamat datang, ${user.full_name}!`);
      navigate(user.role === 'DRIVER' ? '/shipping' : '/dashboard', { replace: true });
    } catch (err) {
      toast.error(getErrorMessage(err, 'Gagal masuk. Periksa nomor HP dan kata sandi.'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-slate-950 p-4 font-sans">
      <div className="w-full max-w-md space-y-8 rounded-3xl border border-slate-800 bg-slate-900 p-8 shadow-2xl">
        <div className="text-center">
          <img src="/logo.png" alt="Lucky Mart" className="mx-auto h-16 w-auto object-contain" />
          <h2 className="mt-4 text-2xl font-extrabold text-white">Lucky Mart</h2>
          <p className="mt-1 text-sm text-slate-400">Dashboard pengelola & staf toko</p>
        </div>

        <form className="space-y-5" onSubmit={handleSubmit}>
          <label className="block">
            <span className="mb-1 block text-xs font-semibold uppercase text-slate-400">Nomor HP</span>
            <div className="relative">
              <Phone className="absolute left-3.5 top-3 h-5 w-5 text-slate-500" />
              <input
                type="tel"
                required
                autoComplete="username"
                value={phoneNumber}
                onChange={(e) => setPhoneNumber(e.target.value)}
                placeholder="08xxxxxxxxxx"
                className="w-full rounded-xl border border-slate-700 bg-slate-800/80 py-2.5 pl-11 pr-4 text-sm text-white placeholder-slate-500 focus:border-brand-500 focus:outline-none focus:ring-1 focus:ring-brand-500"
              />
            </div>
          </label>

          <label className="block">
            <span className="mb-1 block text-xs font-semibold uppercase text-slate-400">Kata sandi</span>
            <div className="relative">
              <Lock className="absolute left-3.5 top-3 h-5 w-5 text-slate-500" />
              <input
                type="password"
                required
                autoComplete="current-password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="w-full rounded-xl border border-slate-700 bg-slate-800/80 py-2.5 pl-11 pr-4 text-sm text-white placeholder-slate-500 focus:border-brand-500 focus:outline-none focus:ring-1 focus:ring-brand-500"
              />
            </div>
          </label>

          <button
            type="submit"
            disabled={loading}
            className="flex w-full items-center justify-center gap-2 rounded-xl bg-brand-600 py-3 text-sm font-bold text-white shadow-lg shadow-brand-600/30 transition-all hover:bg-brand-500 disabled:opacity-50"
          >
            <span>{loading ? 'Memeriksa...' : 'Masuk'}</span>
            <ArrowRight className="h-4 w-4" />
          </button>
        </form>

        <p className="text-center text-xs text-slate-500">Lupa kata sandi? Minta Super Admin mereset dari menu Pengguna.</p>
      </div>
    </div>
  );
};
