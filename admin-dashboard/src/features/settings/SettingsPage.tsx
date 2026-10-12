import React, { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { CreditCard, Lock, Store, Truck, User as UserIcon } from 'lucide-react';
import { PageHeader } from '../../components/common/PageHeader';
import { StatusBadge } from '../../components/common/StatusBadge';
import { Field, btnPrimary, inputClass } from '../../components/common/Modal';
import { Card, FilterTabs } from '../../components/common/Pagination';
import { authService, paymentService, shippingService, storeService } from '../../services/apiServices';
import { useAuthStore, MANAGER_ROLES, hasRole } from '../../store/useAuthStore';
import { formatCurrency, getErrorMessage, ROLE_LABELS } from '../../lib/utils';
import type { PaymentMethod, ShippingMethod, StoreInfo } from '../../types';

type Tab = 'store' | 'profile' | 'password' | 'payment' | 'shipping';

const STORE_FIELDS: { key: Exclude<keyof StoreInfo, 'free_shipping_min'>; label: string; placeholder: string; hint?: string; multiline?: boolean }[] = [
  { key: 'name', label: 'Nama toko', placeholder: 'Lucky Mart' },
  { key: 'tagline', label: 'Slogan singkat', placeholder: 'Toko sembako keluarga, grosir & eceran' },
  { key: 'address', label: 'Alamat toko', placeholder: 'Jl. ..., Kel. ..., Kec. ...', hint: 'Ditampilkan untuk pelanggan yang ambil sendiri di toko.', multiline: true },
  { key: 'maps_url', label: 'Link Google Maps (opsional)', placeholder: 'https://maps.app.goo.gl/...' },
  { key: 'whatsapp', label: 'Nomor WhatsApp toko', placeholder: '081234567890', hint: 'Muncul sebagai tombol "Chat WhatsApp toko" untuk pelanggan.' },
  { key: 'open_hours', label: 'Jam buka', placeholder: 'Setiap hari 07.00 – 21.00' },
  { key: 'announcement', label: 'Pengumuman di bagian atas toko online', placeholder: 'Minyak goreng promo minggu ini!' },
  { key: 'receipt_note', label: 'Catatan di bawah nota', placeholder: 'Terima kasih sudah berbelanja!' },
];

export const SettingsPage: React.FC = () => {
  const { user, setUser } = useAuthStore();
  const isManager = hasRole(user, MANAGER_ROLES);
  const [tab, setTab] = useState<Tab>(isManager ? 'store' : 'profile');
  const [store, setStore] = useState<StoreInfo | null>(null);
  const [busy, setBusy] = useState(false);

  const [profile, setProfile] = useState({ full_name: user?.full_name || '', email: user?.email || '' });
  const [pw, setPw] = useState({ old_password: '', new_password: '', confirm: '' });
  const [methods, setMethods] = useState<PaymentMethod[]>([]);
  const [shipping, setShipping] = useState<ShippingMethod[]>([]);

  useEffect(() => {
    if (tab === 'payment' && isManager) paymentService.getMethods().then(setMethods).catch((e) => toast.error(getErrorMessage(e)));
    if (tab === 'store' && isManager) storeService.get().then(setStore).catch((e) => toast.error(getErrorMessage(e)));
    if (tab === 'shipping') shippingService.getMethods().then(setShipping).catch((e) => toast.error(getErrorMessage(e)));
  }, [tab, isManager]);

  const saveProfile = async () => {
    setBusy(true);
    try {
      const updated = await authService.updateProfile({ full_name: profile.full_name.trim(), email: profile.email.trim() });
      setUser(updated);
      toast.success('Profil disimpan.');
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setBusy(false);
    }
  };

  const savePassword = async () => {
    if (pw.new_password !== pw.confirm) return toast.error('Konfirmasi kata sandi tidak sama.');
    setBusy(true);
    try {
      await authService.changePassword({ old_password: pw.old_password, new_password: pw.new_password });
      setPw({ old_password: '', new_password: '', confirm: '' });
      toast.success('Kata sandi diganti.');
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setBusy(false);
    }
  };

  const saveMethod = async (m: PaymentMethod) => {
    setBusy(true);
    try {
      const updated = await paymentService.updateMethod(m.id, { name: m.name, description: m.description || '', is_active: m.is_active });
      setMethods((ms) => ms.map((x) => (x.id === updated.id ? updated : x)));
      toast.success(`${updated.name} disimpan.`);
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setBusy(false);
    }
  };

  const saveStore = async () => {
    if (!store) return;
    setBusy(true);
    try {
      setStore(await storeService.update(store));
      toast.success('Info toko disimpan.');
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setBusy(false);
    }
  };

  const tabs = [
    ...(isManager ? [{ value: 'store', label: 'Info toko' }] : []),
    { value: 'profile', label: 'Profil saya' },
    { value: 'password', label: 'Kata sandi' },
    ...(isManager ? [{ value: 'payment', label: 'Metode pembayaran' }] : []),
    { value: 'shipping', label: 'Pengiriman' },
  ];

  return (
    <div className="space-y-5">
      <PageHeader title="Pengaturan" subtitle="Akun Anda dan pengaturan toko." />
      <FilterTabs value={tab} onChange={(v) => setTab(v as Tab)} options={tabs} />

      {tab === 'store' && isManager && (
        <Card className="max-w-2xl p-5">
          <div className="mb-4 flex items-center gap-2 text-sm font-semibold">
            <Store className="h-4 w-4" /> Info toko yang dilihat pelanggan dan tercetak di nota
          </div>
          {!store ? (
            <div className="text-sm text-slate-500">Memuat...</div>
          ) : (
            <div className="grid gap-3">
              {STORE_FIELDS.map((f) => (
                <Field key={f.key} label={f.label} hint={f.hint}>
                  {f.multiline ? (
                    <textarea className={inputClass} rows={2} value={store[f.key]} placeholder={f.placeholder} onChange={(e) => setStore({ ...store, [f.key]: e.target.value })} />
                  ) : (
                    <input className={inputClass} value={store[f.key]} placeholder={f.placeholder} onChange={(e) => setStore({ ...store, [f.key]: e.target.value })} />
                  )}
                </Field>
              ))}
              <Field label="Gratis ongkir untuk belanja minimal (Rp)" hint="Isi 0 kalau tidak ada gratis ongkir. Tidak termasuk ongkir, dihitung dari total belanja barang.">
                <input
                  type="number"
                  min={0}
                  step={10000}
                  className={inputClass}
                  value={store.free_shipping_min}
                  onChange={(e) => setStore({ ...store, free_shipping_min: Math.max(0, Number(e.target.value) || 0) })}
                />
              </Field>
              <div>
                <button className={btnPrimary} disabled={busy} onClick={saveStore}>Simpan Info Toko</button>
              </div>
            </div>
          )}
        </Card>
      )}

      {tab === 'profile' && (
        <Card className="max-w-xl p-5">
          <div className="mb-4 flex items-center gap-3">
            <UserIcon className="h-5 w-5 text-slate-400" />
            <div>
              <div className="font-semibold">{user?.phone_number}</div>
              <div className="text-xs text-slate-500">{user && ROLE_LABELS[user.role]} · nomor HP dipakai untuk login dan tidak bisa diubah</div>
            </div>
          </div>
          <div className="grid gap-3">
            <Field label="Nama lengkap">
              <input className={inputClass} value={profile.full_name} onChange={(e) => setProfile({ ...profile, full_name: e.target.value })} />
            </Field>
            <Field label="Email">
              <input type="email" className={inputClass} value={profile.email} onChange={(e) => setProfile({ ...profile, email: e.target.value })} />
            </Field>
            <div>
              <button className={btnPrimary} disabled={busy} onClick={saveProfile}>Simpan</button>
            </div>
          </div>
        </Card>
      )}

      {tab === 'password' && (
        <Card className="max-w-xl p-5">
          <div className="mb-4 flex items-center gap-2 text-sm font-semibold"><Lock className="h-4 w-4" /> Ganti kata sandi</div>
          <div className="grid gap-3">
            <Field label="Kata sandi lama">
              <input type="password" className={inputClass} value={pw.old_password} onChange={(e) => setPw({ ...pw, old_password: e.target.value })} autoComplete="current-password" />
            </Field>
            <Field label="Kata sandi baru" hint="Minimal 8 karakter, ada huruf besar, huruf kecil, angka, dan simbol.">
              <input type="password" className={inputClass} value={pw.new_password} onChange={(e) => setPw({ ...pw, new_password: e.target.value })} autoComplete="new-password" />
            </Field>
            <Field label="Ulangi kata sandi baru">
              <input type="password" className={inputClass} value={pw.confirm} onChange={(e) => setPw({ ...pw, confirm: e.target.value })} autoComplete="new-password" />
            </Field>
            <div>
              <button className={btnPrimary} disabled={busy || !pw.old_password || !pw.new_password} onClick={savePassword}>Ganti Kata Sandi</button>
            </div>
          </div>
        </Card>
      )}

      {tab === 'payment' && isManager && (
        <div className="space-y-3">
          <p className="text-sm text-slate-500">
            Keterangan metode <b>Transfer Bank</b> ditampilkan ke pelanggan saat checkout — tulis nomor rekening dan nama pemilik di sana.
          </p>
          {methods.map((m) => (
            <Card key={m.id} className="p-4">
              <div className="mb-3 flex items-center gap-2">
                <CreditCard className="h-4 w-4 text-slate-400" />
                <span className="font-mono text-xs text-slate-500">{m.code}</span>
                <StatusBadge status={m.is_active ? 'APPROVED' : 'CANCELLED'} label={m.is_active ? 'Aktif' : 'Nonaktif'} />
              </div>
              <div className="grid gap-3 sm:grid-cols-2">
                <Field label="Nama yang dilihat pelanggan">
                  <input className={inputClass} value={m.name} onChange={(e) => setMethods((ms) => ms.map((x) => (x.id === m.id ? { ...x, name: e.target.value } : x)))} />
                </Field>
                <Field label="Keterangan / instruksi">
                  <textarea
                    className={inputClass}
                    rows={2}
                    value={m.description || ''}
                    placeholder={m.type === 'MANUAL_TRANSFER' ? 'BCA 1234567890 a.n. Toko Lucky Mart' : ''}
                    onChange={(e) => setMethods((ms) => ms.map((x) => (x.id === m.id ? { ...x, description: e.target.value } : x)))}
                  />
                </Field>
              </div>
              <div className="mt-3 flex items-center justify-between">
                <label className="flex items-center gap-2 text-sm">
                  <input type="checkbox" checked={m.is_active} onChange={(e) => setMethods((ms) => ms.map((x) => (x.id === m.id ? { ...x, is_active: e.target.checked } : x)))} />
                  Tersedia untuk pelanggan
                </label>
                <button className={btnPrimary} disabled={busy} onClick={() => saveMethod(m)}>Simpan</button>
              </div>
            </Card>
          ))}
        </div>
      )}

      {tab === 'shipping' && (
        <div className="space-y-3">
          {shipping.map((s) => (
            <Card key={s.id} className="flex items-center gap-3 p-4">
              <Truck className="h-5 w-5 text-slate-400" />
              <div className="flex-1">
                <div className="font-semibold">{s.name}</div>
                <div className="text-xs text-slate-500">{s.courier_name} · {s.estimated_days}{s.description ? ` · ${s.description}` : ''}</div>
              </div>
              <div className="font-semibold">{formatCurrency(s.base_rate)}</div>
            </Card>
          ))}
          <p className="text-xs text-slate-500">Ongkir saat ini memakai tarif dasar metode pengiriman. Pengaturan ongkir per wilayah / gratis ongkir dijadwalkan di tahap berikutnya.</p>
        </div>
      )}
    </div>
  );
};
