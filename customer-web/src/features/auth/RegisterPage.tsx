import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { authService } from '../../services/apiServices';
import { Store, User, Phone, Mail, Lock, ArrowRight } from 'lucide-react';
import { toast } from 'sonner';

export const RegisterPage: React.FC = () => {
  const navigate = useNavigate();

  const [fullName, setFullName] = useState('');
  const [phoneNumber, setPhoneNumber] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [role, setRole] = useState<'B2C_CUSTOMER' | 'B2B_CUSTOMER'>('B2C_CUSTOMER');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      await authService.register({
        full_name: fullName,
        phone_number: phoneNumber,
        email,
        password,
        role,
      });

      toast.success('Pendaftaran akun berhasil! Silakan masuk.');
      navigate('/login');
    } catch (err: any) {
      toast.error(err.response?.data?.message || 'Gagal mendaftar. Silakan periksa data Anda.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-slate-950 p-4 font-sans">
      <div className="w-full max-w-md space-y-6 rounded-3xl bg-slate-900 p-8 shadow-2xl border border-slate-800">
        <div className="text-center">
          <div className="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl bg-primary-600 text-white shadow-xl shadow-primary-600/30">
            <Store className="h-8 w-8" />
          </div>
          <h2 className="mt-4 text-2xl font-extrabold text-white">Daftar LuckyMart</h2>
          <p className="mt-1 text-sm text-slate-400">Pilih jenis akun rumah tangga atau mitra toko</p>
        </div>

        <form className="space-y-4" onSubmit={handleSubmit}>
          {/* Account Role Selector */}
          <div className="grid grid-cols-2 gap-2 p-1 bg-slate-800 rounded-xl">
            <button
              type="button"
              onClick={() => setRole('B2C_CUSTOMER')}
              className={`py-2 text-xs font-bold rounded-lg transition-colors ${
                role === 'B2C_CUSTOMER' ? 'bg-primary-600 text-white' : 'text-slate-400'
              }`}
            >
              Rumah Tangga (Eceran)
            </button>
            <button
              type="button"
              onClick={() => setRole('B2B_CUSTOMER')}
              className={`py-2 text-xs font-bold rounded-lg transition-colors ${
                role === 'B2B_CUSTOMER' ? 'bg-primary-600 text-white' : 'text-slate-400'
              }`}
            >
              Mitra Warung (B2B Grosir)
            </button>
          </div>

          <div>
            <label className="block text-xs font-semibold uppercase text-slate-400 mb-1">Nama Lengkap</label>
            <div className="relative">
              <User className="absolute left-3.5 top-3 h-5 w-5 text-slate-500" />
              <input
                type="text"
                required
                value={fullName}
                onChange={(e) => setFullName(e.target.value)}
                placeholder="Budi Santoso"
                className="w-full rounded-xl bg-slate-800/80 border border-slate-700 pl-11 pr-4 py-2.5 text-sm text-white placeholder-slate-500 focus:border-primary-500 focus:outline-none"
              />
            </div>
          </div>

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
            <label className="block text-xs font-semibold uppercase text-slate-400 mb-1">Email</label>
            <div className="relative">
              <Mail className="absolute left-3.5 top-3 h-5 w-5 text-slate-500" />
              <input
                type="email"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="pelanggan@luckymart.com"
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

          <button
            type="submit"
            disabled={loading}
            className="flex w-full justify-center items-center space-x-2 rounded-xl bg-primary-600 py-3 text-sm font-bold text-white hover:bg-primary-500 transition-all shadow-lg shadow-primary-600/30 disabled:opacity-50"
          >
            <span>{loading ? 'Mendaftarkan...' : 'Daftar Sekarang'}</span>
            <ArrowRight className="h-4 w-4" />
          </button>
        </form>

        <div className="text-center text-xs text-slate-400">
          Sudah punya akun?{' '}
          <Link to="/login" className="text-primary-400 font-bold hover:underline">
            Masuk Akun
          </Link>
        </div>
      </div>
    </div>
  );
};
