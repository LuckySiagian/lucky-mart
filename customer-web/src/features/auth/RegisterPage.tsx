import React, { useState } from 'react';
import { Link, Navigate, useNavigate, useSearchParams } from 'react-router-dom';
import { ArrowRight, Lock, Mail, MapPin, Phone, Store, User } from 'lucide-react';
import { toast } from 'sonner';
import { authService } from '../../services/apiServices';
import { useAuthStore } from '../../store/useAuthStore';
import { getErrorMessage, isPhone, PASSWORD_HINT } from '../../lib/utils';
import { AuthShell, authInput } from './AuthShell';

const safeNext = (next: string | null) => (next && next.startsWith('/') && !next.startsWith('//') ? next : '/');

export const RegisterPage: React.FC = () => {
  const navigate = useNavigate();
  const [params] = useSearchParams();
  const next = safeNext(params.get('next'));
  const { isAuthenticated, setAuth } = useAuthStore();

  const [form, setForm] = useState({ full_name: '', phone_number: '', email: '', password: '', store_name: '', store_address: '' });
  const [isB2B, setIsB2B] = useState(false);
  const [loading, setLoading] = useState(false);

  if (isAuthenticated) return <Navigate to={next} replace />;

  const set = (key: keyof typeof form) => (e: React.ChangeEvent<HTMLInputElement>) => setForm({ ...form, [key]: e.target.value });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (form.full_name.trim().length < 3) return toast.error('Nama minimal 3 huruf.');
    if (!isPhone(form.phone_number)) return toast.error('Nomor HP tidak valid (contoh: 081234567890).');
    if (isB2B && (form.store_name.trim().length < 3 || !form.store_address.trim())) return toast.error('Isi nama dan alamat toko.');

    setLoading(true);
    try {
      const res = await authService.register({
        full_name: form.full_name.trim(),
        phone_number: form.phone_number.trim(),
        email: form.email.trim() || undefined,
        password: form.password,
        is_b2b_registration: isB2B,
        store_name: isB2B ? form.store_name.trim() : undefined,
        store_address: isB2B ? form.store_address.trim() : undefined,
      });
      setAuth(res.user, res.access_token, res.refresh_token);
      toast.success(isB2B ? 'Akun dibuat! Toko akan memverifikasi data mitra Anda.' : 'Akun berhasil dibuat. Selamat berbelanja!');
      navigate(next, { replace: true });
    } catch (err) {
      toast.error(getErrorMessage(err, 'Gagal mendaftar. Periksa kembali data Anda.'));
    } finally {
      setLoading(false);
    }
  };

  const tab = (active: boolean) => `rounded-lg py-2 text-xs font-bold transition-colors ${active ? 'bg-primary-600 text-white' : 'text-slate-400'}`;

  return (
    <AuthShell title="Daftar" subtitle="Buat akun untuk mulai belanja">
      <form className="space-y-4" onSubmit={handleSubmit}>
        <div className="grid grid-cols-2 gap-2 rounded-xl bg-slate-800 p-1">
          <button type="button" onClick={() => setIsB2B(false)} className={tab(!isB2B)}>Rumah Tangga</button>
          <button type="button" onClick={() => setIsB2B(true)} className={tab(isB2B)}>Mitra Warung (Grosir)</button>
        </div>
        {isB2B && <p className="text-xs text-slate-400">Mitra warung mendapat harga grosir dan bisa mengajukan kasbon setelah diverifikasi toko.</p>}

        {[
          { key: 'full_name' as const, label: 'Nama lengkap', icon: User, type: 'text', auto: 'name', required: true },
          { key: 'phone_number' as const, label: 'Nomor HP', icon: Phone, type: 'tel', auto: 'tel', required: true, placeholder: '08xxxxxxxxxx' },
          { key: 'email' as const, label: 'Email (opsional)', icon: Mail, type: 'email', auto: 'email', required: false },
          ...(isB2B
            ? [
                { key: 'store_name' as const, label: 'Nama toko / warung', icon: Store, type: 'text', auto: 'organization', required: true },
                { key: 'store_address' as const, label: 'Alamat toko', icon: MapPin, type: 'text', auto: 'street-address', required: true },
              ]
            : []),
        ].map((f) => (
          <label key={f.key} className="block">
            <span className="mb-1 block text-xs font-semibold uppercase text-slate-400">{f.label}</span>
            <div className="relative">
              <f.icon className="absolute left-3.5 top-3 h-5 w-5 text-slate-500" />
              <input type={f.type} required={f.required} autoComplete={f.auto} value={form[f.key]} onChange={set(f.key)} placeholder={'placeholder' in f ? f.placeholder : undefined} className={authInput} />
            </div>
          </label>
        ))}

        <label className="block">
          <span className="mb-1 block text-xs font-semibold uppercase text-slate-400">Kata sandi</span>
          <div className="relative">
            <Lock className="absolute left-3.5 top-3 h-5 w-5 text-slate-500" />
            <input type="password" required autoComplete="new-password" value={form.password} onChange={set('password')} className={authInput} />
          </div>
          <span className="mt-1 block text-[11px] text-slate-500">{PASSWORD_HINT}</span>
        </label>

        <button
          type="submit"
          disabled={loading}
          className="flex w-full items-center justify-center gap-2 rounded-xl bg-primary-600 py-3 text-sm font-bold text-white shadow-lg shadow-primary-600/30 hover:bg-primary-500 disabled:opacity-50"
        >
          {loading ? 'Mendaftarkan...' : 'Daftar'}
          <ArrowRight className="h-4 w-4" />
        </button>
      </form>

      <p className="text-center text-sm text-slate-400">
        Sudah punya akun?{' '}
        <Link to={`/login${next !== '/' ? `?next=${encodeURIComponent(next)}` : ''}`} className="font-bold text-primary-400 hover:underline">
          Masuk
        </Link>
      </p>
    </AuthShell>
  );
};
