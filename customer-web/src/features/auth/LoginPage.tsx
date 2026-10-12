import React, { useState } from 'react';
import { Link, Navigate, useNavigate, useSearchParams } from 'react-router-dom';
import { ArrowRight, Lock, Phone } from 'lucide-react';
import { toast } from 'sonner';
import { useAuthStore } from '../../store/useAuthStore';
import { authService } from '../../services/apiServices';
import { getErrorMessage, isPhone } from '../../lib/utils';
import { AuthShell, authInput } from './AuthShell';

const STAFF = ['SUPER_ADMIN', 'STORE_MANAGER', 'WAREHOUSE_STAFF', 'DRIVER'];

// Only follow same-site paths after login, never a full URL from the query.
const safeNext = (next: string | null) => (next && next.startsWith('/') && !next.startsWith('//') ? next : '/');

export const LoginPage: React.FC = () => {
  const navigate = useNavigate();
  const [params] = useSearchParams();
  const next = safeNext(params.get('next'));
  const { isAuthenticated, setAuth } = useAuthStore();

  const [phoneNumber, setPhoneNumber] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [forgot, setForgot] = useState(false);

  if (isAuthenticated) return <Navigate to={next} replace />;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!isPhone(phoneNumber)) return toast.error('Nomor HP tidak valid (contoh: 081234567890).');
    setLoading(true);
    try {
      if (forgot) {
        const message = await authService.forgotPassword(phoneNumber.trim());
        toast.success(message, { duration: 8000 });
        setForgot(false);
        return;
      }
      const res = await authService.login({ phone_number: phoneNumber.trim(), password });
      if (STAFF.includes(res.user.role)) {
        toast.error('Ini akun staf toko. Silakan masuk lewat dashboard admin.');
        return;
      }
      setAuth(res.user, res.access_token, res.refresh_token);
      toast.success(`Selamat datang, ${res.user.full_name}!`);
      navigate(next, { replace: true });
    } catch (err) {
      toast.error(getErrorMessage(err, 'Gagal masuk. Periksa nomor HP dan kata sandi.'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <AuthShell title={forgot ? 'Lupa kata sandi' : 'Masuk'} subtitle={forgot ? 'Masukkan nomor HP akun Anda' : 'Belanja kebutuhan dapur dengan mudah'}>
      <form className="space-y-4" onSubmit={handleSubmit}>
        <label className="block">
          <span className="mb-1 block text-xs font-semibold uppercase text-slate-400">Nomor HP</span>
          <div className="relative">
            <Phone className="absolute left-3.5 top-3 h-5 w-5 text-slate-500" />
            <input type="tel" required autoComplete="username" value={phoneNumber} onChange={(e) => setPhoneNumber(e.target.value)} placeholder="08xxxxxxxxxx" className={authInput} />
          </div>
        </label>

        {!forgot && (
          <label className="block">
            <span className="mb-1 block text-xs font-semibold uppercase text-slate-400">Kata sandi</span>
            <div className="relative">
              <Lock className="absolute left-3.5 top-3 h-5 w-5 text-slate-500" />
              <input type="password" required autoComplete="current-password" value={password} onChange={(e) => setPassword(e.target.value)} className={authInput} />
            </div>
          </label>
        )}

        <button
          type="submit"
          disabled={loading}
          className="flex w-full items-center justify-center gap-2 rounded-xl bg-primary-600 py-3 text-sm font-bold text-white shadow-lg shadow-primary-600/30 hover:bg-primary-500 disabled:opacity-50"
        >
          {loading ? 'Memproses...' : forgot ? 'Kirim Permintaan' : 'Masuk'}
          <ArrowRight className="h-4 w-4" />
        </button>
      </form>

      <div className="space-y-2 text-center text-sm text-slate-400">
        <button onClick={() => setForgot(!forgot)} className="font-semibold text-primary-400 hover:underline">
          {forgot ? 'Kembali ke halaman masuk' : 'Lupa kata sandi?'}
        </button>
        {!forgot && (
          <p>
            Belum punya akun?{' '}
            <Link to={`/register${next !== '/' ? `?next=${encodeURIComponent(next)}` : ''}`} className="font-bold text-primary-400 hover:underline">
              Daftar
            </Link>
          </p>
        )}
      </div>
    </AuthShell>
  );
};
