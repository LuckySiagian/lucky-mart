import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuthStore } from '../../store/useAuthStore';
import { authService } from '../../services/apiServices';
import { Store, Lock, Phone, ArrowRight } from 'lucide-react';
import { toast } from 'sonner';

export const LoginPage: React.FC = () => {
  const navigate = useNavigate();
  const setAuth = useAuthStore((state) => state.setAuth);

  const [phoneNumber, setPhoneNumber] = useState('081234567890');
  const [password, setPassword] = useState('SuperAdmin123!');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      const res = await authService.login({
        phone_number: phoneNumber,
        password,
      });

      const { user, access_token, refresh_token } = res.data;
      setAuth(user, access_token, refresh_token);
      toast.success(`Selamat datang kembali, ${user.full_name}!`);
      navigate('/products');
    } catch (err: any) {
      toast.error(err.response?.data?.message || 'Gagal masuk. Periksa nomor telepon & password Anda.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-slate-950 p-4 font-sans">
      <div className="w-full max-w-md space-y-8 rounded-3xl bg-slate-900 p-8 shadow-2xl border border-slate-800">
        <div className="text-center">
          <div className="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl bg-primary-600 text-white shadow-xl shadow-primary-600/30">
            <Store className="h-8 w-8" />
          </div>
          <h2 className="mt-4 text-2xl font-extrabold text-white">LuckyMart</h2>
          <p className="mt-1 text-sm text-slate-400">Masuk Akun Pelanggan & B2B Mitra Warung</p>
        </div>

        <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
          <div className="space-y-4">
            <div>
              <label className="block text-xs font-semibold uppercase text-slate-400 mb-1">Nomor Telepon</label>
              <div className="relative">
                <Phone className="absolute left-3.5 top-3 h-5 w-5 text-slate-500" />
                <input
                  type="text"
                  required
                  value={phoneNumber}
                  onChange={(e) => setPhoneNumber(e.target.value)}
                  placeholder="081234567890"
                  className="w-full rounded-xl bg-slate-800/80 border border-slate-700 pl-11 pr-4 py-2.5 text-sm text-white placeholder-slate-500 focus:border-primary-500 focus:outline-none"
                />
              </div>
            </div>

            <div>
              <label className="block text-xs font-semibold uppercase text-slate-400 mb-1">Kata Sandi</label>
              <div className="relative">
                <Lock className="absolute left-3.5 top-3 h-5 w-5 text-slate-500" />
                <input
                  type="password"
                  required
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="••••••••"
                  className="w-full rounded-xl bg-slate-800/80 border border-slate-700 pl-11 pr-4 py-2.5 text-sm text-white placeholder-slate-500 focus:border-primary-500 focus:outline-none"
                />
              </div>
            </div>
          </div>

          <button
            type="submit"
            disabled={loading}
            className="flex w-full justify-center items-center space-x-2 rounded-xl bg-primary-600 py-3 text-sm font-bold text-white hover:bg-primary-500 transition-all shadow-lg shadow-primary-600/30 disabled:opacity-50"
          >
            <span>{loading ? 'Memproses...' : 'Masuk Sekarang'}</span>
            <ArrowRight className="h-4 w-4" />
          </button>
        </form>

        <div className="text-center text-xs text-slate-400">
          Belum punya akun?{' '}
          <Link to="/register" className="text-primary-400 font-bold hover:underline">
            Daftar Akun Baru
          </Link>
        </div>
      </div>
    </div>
  );
};
