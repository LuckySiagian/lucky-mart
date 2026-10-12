import React, { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { CreditCard, Lock, MapPin, Star, Trash2, User as UserIcon } from 'lucide-react';
import { addressService, authService, creditService } from '../../services/apiServices';
import { useAuthStore } from '../../store/useAuthStore';
import { Badge, Field, Panel, btnPrimary, btnSecondary, inputClass } from '../../components/common/Badge';
import { formatCurrency, getErrorMessage, isPhone, PASSWORD_HINT } from '../../lib/utils';
import type { Address, CreditLedger } from '../../types';

type Tab = 'profile' | 'address' | 'password' | 'kasbon';

const emptyAddress = { title: 'Rumah', recipient_name: '', phone: '', full_address: '', city: '', province: '', postal_code: '', is_default: false };

export const ProfilePage: React.FC = () => {
  const { user, setUser } = useAuthStore();
  const isB2B = user?.role === 'B2B_CUSTOMER';
  const [tab, setTab] = useState<Tab>('profile');
  const [busy, setBusy] = useState(false);

  const [profile, setProfile] = useState({ full_name: user?.full_name || '', email: user?.email || '' });
  const [pw, setPw] = useState({ old_password: '', new_password: '', confirm: '' });
  const [addresses, setAddresses] = useState<Address[]>([]);
  const [addrForm, setAddrForm] = useState<typeof emptyAddress | null>(null);
  const [editingAddr, setEditingAddr] = useState<string | null>(null);
  const [credit, setCredit] = useState<CreditLedger | null>(null);
  const [creditMissing, setCreditMissing] = useState(false);

  // Refresh the profile so B2B verification changes show without re-login.
  useEffect(() => {
    authService.me().then(setUser).catch(() => undefined);
  }, [setUser]);

  const loadAddresses = () => addressService.getAddresses().then(setAddresses).catch((e) => toast.error(getErrorMessage(e)));

  useEffect(() => {
    if (tab === 'address') loadAddresses();
    if (tab === 'kasbon') {
      creditService
        .getMine()
        .then((c) => {
          setCredit(c);
          setCreditMissing(false);
        })
        .catch(() => setCreditMissing(true));
    }
  }, [tab]);

  const wrap = async (fn: () => Promise<unknown>, ok: string) => {
    setBusy(true);
    try {
      await fn();
      toast.success(ok);
      return true;
    } catch (err) {
      toast.error(getErrorMessage(err));
      return false;
    } finally {
      setBusy(false);
    }
  };

  const saveProfile = () =>
    wrap(async () => setUser(await authService.updateProfile({ full_name: profile.full_name.trim(), email: profile.email.trim() })), 'Profil disimpan.');

  const savePassword = async () => {
    if (pw.new_password !== pw.confirm) return toast.error('Konfirmasi kata sandi tidak sama.');
    if (await wrap(() => authService.changePassword({ old_password: pw.old_password, new_password: pw.new_password }), 'Kata sandi diganti.')) {
      setPw({ old_password: '', new_password: '', confirm: '' });
    }
  };

  const saveAddress = async () => {
    if (!addrForm) return;
    if (!isPhone(addrForm.phone)) return toast.error('Nomor HP tidak valid (contoh: 081234567890).');
    const ok = await wrap(
      () => (editingAddr ? addressService.updateAddress(editingAddr, addrForm) : addressService.createAddress(addrForm)),
      'Alamat disimpan.'
    );
    if (ok) {
      setAddrForm(null);
      setEditingAddr(null);
      loadAddresses();
    }
  };

  const tabs: { value: Tab; label: string; icon: React.ElementType }[] = [
    { value: 'profile', label: 'Profil', icon: UserIcon },
    { value: 'address', label: 'Alamat', icon: MapPin },
    { value: 'password', label: 'Kata Sandi', icon: Lock },
    ...(isB2B ? [{ value: 'kasbon' as Tab, label: 'Kasbon', icon: CreditCard }] : []),
  ];

  const b2b = user?.b2b_profile;

  return (
    <div className="mx-auto max-w-3xl space-y-5 px-4 py-8 sm:px-6 lg:px-8">
      <h1 className="text-2xl font-black text-slate-900 dark:text-white">Akun Saya</h1>

      {b2b && (
        <div className="rounded-3xl border border-slate-200 bg-white p-4 text-sm dark:border-slate-800 dark:bg-slate-900">
          <div className="flex items-center justify-between gap-2">
            <span className="font-bold">Mitra toko: {b2b.store_name}</span>
            <Badge tone={b2b.verification_status === 'APPROVED' ? 'green' : b2b.verification_status === 'REJECTED' ? 'red' : 'amber'}>
              {b2b.verification_status === 'APPROVED' ? 'Terverifikasi' : b2b.verification_status === 'REJECTED' ? 'Ditolak' : 'Menunggu verifikasi'}
            </Badge>
          </div>
          {b2b.verification_status === 'PENDING' && <p className="mt-1 text-xs text-slate-500">Toko akan memverifikasi data Anda. Setelah disetujui, Anda mendapat harga grosir.</p>}
        </div>
      )}

      <div className="flex gap-2 overflow-x-auto">
        {tabs.map((t) => (
          <button
            key={t.value}
            onClick={() => setTab(t.value)}
            className={`flex flex-shrink-0 items-center gap-1.5 rounded-xl px-4 py-2 text-xs font-bold ${
              tab === t.value ? 'bg-primary-600 text-white' : 'border border-slate-200 bg-white text-slate-600 dark:border-slate-800 dark:bg-slate-900 dark:text-slate-300'
            }`}
          >
            <t.icon className="h-4 w-4" /> {t.label}
          </button>
        ))}
      </div>

      {tab === 'profile' && (
        <Panel>
          <div className="grid gap-3">
            <Field label="Nomor HP (untuk masuk)">
              <input className={`${inputClass} opacity-60`} value={user?.phone_number || ''} disabled />
            </Field>
            <Field label="Nama lengkap">
              <input className={inputClass} value={profile.full_name} onChange={(e) => setProfile({ ...profile, full_name: e.target.value })} />
            </Field>
            <Field label="Email (opsional)">
              <input type="email" className={inputClass} value={profile.email} onChange={(e) => setProfile({ ...profile, email: e.target.value })} />
            </Field>
            <div><button className={btnPrimary} disabled={busy} onClick={saveProfile}>Simpan</button></div>
          </div>
        </Panel>
      )}

      {tab === 'password' && (
        <Panel>
          <div className="grid gap-3">
            <Field label="Kata sandi lama">
              <input type="password" className={inputClass} value={pw.old_password} onChange={(e) => setPw({ ...pw, old_password: e.target.value })} autoComplete="current-password" />
            </Field>
            <Field label="Kata sandi baru" hint={PASSWORD_HINT}>
              <input type="password" className={inputClass} value={pw.new_password} onChange={(e) => setPw({ ...pw, new_password: e.target.value })} autoComplete="new-password" />
            </Field>
            <Field label="Ulangi kata sandi baru">
              <input type="password" className={inputClass} value={pw.confirm} onChange={(e) => setPw({ ...pw, confirm: e.target.value })} autoComplete="new-password" />
            </Field>
            <div><button className={btnPrimary} disabled={busy || !pw.old_password || !pw.new_password} onClick={savePassword}>Ganti Kata Sandi</button></div>
          </div>
        </Panel>
      )}

      {tab === 'address' && (
        <div className="space-y-3">
          {addresses.map((a) => (
            <Panel key={a.id}>
              <div className="flex items-start justify-between gap-3 text-sm">
                <div>
                  <div className="flex items-center gap-2 font-bold">
                    {a.title} {a.is_default && <Badge tone="green">Utama</Badge>}
                  </div>
                  <div>{a.recipient_name} · {a.phone}</div>
                  <div className="text-slate-500">{a.full_address}, {a.city}, {a.province} {a.postal_code}</div>
                </div>
                <div className="flex flex-shrink-0 gap-1">
                  {!a.is_default && (
                    <button title="Jadikan utama" onClick={() => wrap(() => addressService.setDefault(a.id), 'Alamat utama diubah.').then(loadAddresses)} className="rounded-lg p-2 text-slate-500 hover:bg-slate-100 dark:hover:bg-slate-800">
                      <Star className="h-4 w-4" />
                    </button>
                  )}
                  <button
                    onClick={() => {
                      setEditingAddr(a.id);
                      setAddrForm({ title: a.title, recipient_name: a.recipient_name, phone: a.phone, full_address: a.full_address, city: a.city, province: a.province, postal_code: a.postal_code, is_default: a.is_default });
                    }}
                    className="rounded-lg px-2 text-xs font-bold text-primary-600"
                  >
                    Ubah
                  </button>
                  <button title="Hapus" onClick={() => window.confirm('Hapus alamat ini?') && wrap(() => addressService.deleteAddress(a.id), 'Alamat dihapus.').then(loadAddresses)} className="rounded-lg p-2 text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-950/30">
                    <Trash2 className="h-4 w-4" />
                  </button>
                </div>
              </div>
            </Panel>
          ))}

          {addrForm ? (
            <Panel title={editingAddr ? 'Ubah alamat' : 'Alamat baru'}>
              <div className="grid gap-3 sm:grid-cols-2">
                <Field label="Label"><input className={inputClass} value={addrForm.title} onChange={(e) => setAddrForm({ ...addrForm, title: e.target.value })} placeholder="Rumah / Toko" /></Field>
                <Field label="Nama penerima"><input className={inputClass} value={addrForm.recipient_name} onChange={(e) => setAddrForm({ ...addrForm, recipient_name: e.target.value })} /></Field>
                <Field label="Nomor HP"><input type="tel" className={inputClass} value={addrForm.phone} onChange={(e) => setAddrForm({ ...addrForm, phone: e.target.value })} placeholder="08xxxxxxxxxx" /></Field>
                <Field label="Kota / kabupaten"><input className={inputClass} value={addrForm.city} onChange={(e) => setAddrForm({ ...addrForm, city: e.target.value })} /></Field>
                <Field label="Alamat lengkap" className="sm:col-span-2"><textarea rows={2} className={inputClass} value={addrForm.full_address} onChange={(e) => setAddrForm({ ...addrForm, full_address: e.target.value })} /></Field>
                <Field label="Provinsi"><input className={inputClass} value={addrForm.province} onChange={(e) => setAddrForm({ ...addrForm, province: e.target.value })} /></Field>
                <Field label="Kode pos"><input className={inputClass} value={addrForm.postal_code} onChange={(e) => setAddrForm({ ...addrForm, postal_code: e.target.value })} /></Field>
                <label className="flex items-center gap-2 text-sm sm:col-span-2">
                  <input type="checkbox" checked={addrForm.is_default} onChange={(e) => setAddrForm({ ...addrForm, is_default: e.target.checked })} /> Jadikan alamat utama
                </label>
              </div>
              <div className="mt-4 flex gap-2">
                <button className={btnSecondary} onClick={() => { setAddrForm(null); setEditingAddr(null); }}>Batal</button>
                <button className={btnPrimary} disabled={busy} onClick={saveAddress}>Simpan</button>
              </div>
            </Panel>
          ) : (
            <button className={btnPrimary} onClick={() => { setEditingAddr(null); setAddrForm({ ...emptyAddress, recipient_name: user?.full_name || '', phone: user?.phone_number || '', is_default: addresses.length === 0 }); }}>
              + Tambah Alamat
            </button>
          )}
        </div>
      )}

      {tab === 'kasbon' && (
        <Panel title="Kasbon">
          {creditMissing ? (
            <p className="text-sm text-slate-500">Anda belum mendapat fasilitas kasbon. Hubungi toko setelah akun mitra Anda terverifikasi.</p>
          ) : !credit ? (
            <p className="text-sm text-slate-500">Memuat...</p>
          ) : (
            <div className="grid gap-3 sm:grid-cols-3">
              <div className="rounded-2xl bg-slate-50 p-4 dark:bg-slate-800"><div className="text-xs text-slate-500">Limit</div><div className="text-lg font-black">{formatCurrency(credit.credit_limit)}</div></div>
              <div className="rounded-2xl bg-rose-50 p-4 dark:bg-rose-950/30"><div className="text-xs text-slate-500">Hutang berjalan</div><div className="text-lg font-black text-rose-600">{formatCurrency(credit.outstanding_balance)}</div></div>
              <div className="rounded-2xl bg-emerald-50 p-4 dark:bg-emerald-950/30"><div className="text-xs text-slate-500">Sisa limit</div><div className="text-lg font-black text-emerald-600">{formatCurrency(credit.remaining_credit)}</div></div>
              <p className="text-xs text-slate-500 sm:col-span-3">Pelunasan kasbon dilakukan langsung ke toko dan akan dicatat oleh admin.</p>
            </div>
          )}
        </Panel>
      )}
    </div>
  );
};
