import React, { useState } from 'react';
import { toast } from 'sonner';
import { BadgeCheck, KeyRound, Pencil, Plus, Trash2 } from 'lucide-react';
import { PageHeader } from '../../components/common/PageHeader';
import { StatusBadge } from '../../components/common/StatusBadge';
import { LoadingSkeleton } from '../../components/common/LoadingSkeleton';
import { EmptyState } from '../../components/common/EmptyState';
import { ErrorState } from '../../components/common/ErrorState';
import { Field, Modal, btnDanger, btnPrimary, btnSecondary, inputClass } from '../../components/common/Modal';
import { Card, FilterTabs, Pagination, SearchInput, Table } from '../../components/common/Pagination';
import { userService } from '../../services/apiServices';
import { useList } from '../../hooks/useList';
import { useAuthStore } from '../../store/useAuthStore';
import { formatCurrency, formatDate, getErrorMessage, ROLE_LABELS } from '../../lib/utils';
import type { Role, User } from '../../types';

const TABS = [
  { value: '', label: 'Semua' },
  { value: 'B2C_CUSTOMER', label: 'Pelanggan' },
  { value: 'B2B_CUSTOMER', label: 'Mitra Toko' },
  { value: 'WAREHOUSE_STAFF', label: 'Staf Gudang' },
  { value: 'DRIVER', label: 'Kurir' },
  { value: 'STORE_MANAGER', label: 'Manajer' },
  { value: 'SUPER_ADMIN', label: 'Super Admin' },
];

type Mode = 'create' | 'edit' | 'password' | 'b2b' | 'delete' | null;

const PASSWORD_HINT = 'Minimal 8 karakter, ada huruf besar, huruf kecil, angka, dan simbol.';

export const UserListPage: React.FC = () => {
  const me = useAuthStore((s) => s.user);
  const list = useList((p) => userService.getUsers(p));
  const [mode, setMode] = useState<Mode>(null);
  const [target, setTarget] = useState<User | null>(null);
  const [busy, setBusy] = useState(false);
  const [form, setForm] = useState({ full_name: '', phone_number: '', email: '', password: '', role: 'WAREHOUSE_STAFF' as Role, is_active: true });
  const [b2b, setB2b] = useState({ status: 'APPROVED' as 'APPROVED' | 'REJECTED' | 'PENDING', credit_limit: '0' });

  const assignableRoles: Role[] =
    me?.role === 'SUPER_ADMIN'
      ? ['SUPER_ADMIN', 'STORE_MANAGER', 'WAREHOUSE_STAFF', 'DRIVER', 'B2B_CUSTOMER', 'B2C_CUSTOMER']
      : ['WAREHOUSE_STAFF', 'DRIVER', 'B2B_CUSTOMER', 'B2C_CUSTOMER'];
  const canManage = (u: User) => me?.role === 'SUPER_ADMIN' || (u.role !== 'SUPER_ADMIN' && u.role !== 'STORE_MANAGER');

  const open = (m: Mode, u: User | null) => {
    setMode(m);
    setTarget(u);
    if (m === 'create') setForm({ full_name: '', phone_number: '', email: '', password: '', role: 'WAREHOUSE_STAFF', is_active: true });
    if (m === 'edit' && u) setForm({ full_name: u.full_name, phone_number: u.phone_number, email: u.email || '', password: '', role: u.role, is_active: u.is_active });
    if (m === 'password') setForm((f) => ({ ...f, password: '' }));
    if (m === 'b2b' && u) setB2b({ status: 'APPROVED', credit_limit: String(u.b2b_profile?.approved_credit_limit || 0) });
  };
  const close = () => {
    setMode(null);
    setTarget(null);
  };

  const run = async (fn: () => Promise<unknown>, success: string) => {
    setBusy(true);
    try {
      await fn();
      toast.success(success);
      close();
      list.reload();
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setBusy(false);
    }
  };

  const submit = () => {
    if (mode === 'create') {
      return run(
        () => userService.createUser({ full_name: form.full_name.trim(), phone_number: form.phone_number.trim(), email: form.email.trim() || undefined, password: form.password, role: form.role }),
        'Akun berhasil dibuat.'
      );
    }
    if (mode === 'edit' && target) {
      return run(
        () => userService.updateUser(target.id, { full_name: form.full_name.trim(), email: form.email.trim(), role: form.role, is_active: form.is_active }),
        'Data pengguna diperbarui.'
      );
    }
    if (mode === 'password' && target) {
      return run(() => userService.resetPassword(target.id, form.password), 'Kata sandi direset. Sampaikan kata sandi baru ke pengguna secara langsung.');
    }
    if (mode === 'b2b' && target) {
      return run(
        () => userService.verifyB2B(target.id, b2b.status, b2b.status === 'APPROVED' ? Number(b2b.credit_limit) || 0 : undefined),
        b2b.status === 'APPROVED' ? 'Mitra toko disetujui.' : 'Status mitra diperbarui.'
      );
    }
    if (mode === 'delete' && target) {
      return run(() => userService.deleteUser(target.id), 'Pengguna dihapus.');
    }
  };

  return (
    <div className="space-y-5">
      <PageHeader title="Pengguna" subtitle="Akun staf toko, pelanggan, dan verifikasi mitra toko (B2B).">
        <button className={btnPrimary} onClick={() => open('create', null)}>
          <Plus className="h-4 w-4" /> Tambah Akun
        </button>
      </PageHeader>

      <div className="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
        <FilterTabs value={(list.params.role as string) || ''} onChange={(v) => list.setFilter('role', v)} options={TABS} />
        <SearchInput value={list.searchInput} onChange={list.setSearchInput} placeholder="Cari nama / nomor HP..." />
      </div>

      {list.loading ? (
        <LoadingSkeleton count={5} />
      ) : list.error ? (
        <ErrorState message={list.error} onRetry={list.reload} />
      ) : list.items.length === 0 ? (
        <EmptyState title="Tidak ada pengguna" />
      ) : (
        <Card className="p-2 sm:p-4">
          <Table
            head={
              <tr>
                <th className="px-3 py-3">Nama</th>
                <th className="px-3 py-3">Peran</th>
                <th className="px-3 py-3">Status</th>
                <th className="px-3 py-3">Login terakhir</th>
                <th className="px-3 py-3" />
              </tr>
            }
          >
            {list.items.map((u) => (
              <tr key={u.id} className="hover:bg-slate-50 dark:hover:bg-slate-800/40">
                <td className="px-3 py-3">
                  <div className="font-semibold">{u.full_name} {u.id === me?.id && <span className="text-xs font-normal text-slate-400">(Anda)</span>}</div>
                  <div className="text-xs text-slate-500">{u.phone_number}{u.email ? ` · ${u.email}` : ''}</div>
                  {u.b2b_profile && <div className="text-xs text-slate-500">Toko: {u.b2b_profile.store_name}</div>}
                </td>
                <td className="px-3 py-3">
                  <div className="text-xs font-semibold">{ROLE_LABELS[u.role]}</div>
                  {u.b2b_profile && <div className="mt-1"><StatusBadge status={u.b2b_profile.verification_status} /></div>}
                </td>
                <td className="px-3 py-3">
                  {u.is_active ? <StatusBadge status="APPROVED" label="Aktif" /> : <StatusBadge status="CANCELLED" label="Nonaktif" />}
                </td>
                <td className="px-3 py-3 text-xs text-slate-500">{formatDate(u.last_login_at)}</td>
                <td className="px-3 py-3">
                  {canManage(u) && (
                    <div className="flex justify-end gap-1">
                      {u.role === 'B2B_CUSTOMER' && u.b2b_profile && (
                        <button title="Verifikasi mitra toko" onClick={() => open('b2b', u)} className="rounded-lg p-2 text-emerald-600 hover:bg-emerald-50 dark:hover:bg-emerald-950/40">
                          <BadgeCheck className="h-4 w-4" />
                        </button>
                      )}
                      <button title="Ubah" onClick={() => open('edit', u)} className="rounded-lg p-2 text-slate-600 hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-800">
                        <Pencil className="h-4 w-4" />
                      </button>
                      <button title="Reset kata sandi" onClick={() => open('password', u)} className="rounded-lg p-2 text-slate-600 hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-800">
                        <KeyRound className="h-4 w-4" />
                      </button>
                      {u.id !== me?.id && (
                        <button title="Hapus" onClick={() => open('delete', u)} className="rounded-lg p-2 text-rose-600 hover:bg-rose-50 dark:hover:bg-rose-950/30">
                          <Trash2 className="h-4 w-4" />
                        </button>
                      )}
                    </div>
                  )}
                </td>
              </tr>
            ))}
          </Table>
          <Pagination meta={list.meta} onPage={list.setPage} />
        </Card>
      )}

      <Modal
        open={mode === 'create' || mode === 'edit'}
        onClose={close}
        title={mode === 'create' ? 'Tambah akun' : `Ubah ${target?.full_name}`}
        footer={
          <>
            <button className={btnSecondary} onClick={close}>Batal</button>
            <button className={btnPrimary} disabled={busy} onClick={submit}>Simpan</button>
          </>
        }
      >
        <div className="grid gap-3">
          <Field label="Nama lengkap">
            <input className={inputClass} value={form.full_name} onChange={(e) => setForm({ ...form, full_name: e.target.value })} />
          </Field>
          {mode === 'create' && (
            <Field label="Nomor HP (untuk login)" hint="Format: 08xxxxxxxxxx">
              <input className={inputClass} value={form.phone_number} onChange={(e) => setForm({ ...form, phone_number: e.target.value })} />
            </Field>
          )}
          <Field label="Email (opsional)">
            <input type="email" className={inputClass} value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })} />
          </Field>
          <Field label="Peran">
            <select className={inputClass} value={form.role} disabled={target?.id === me?.id} onChange={(e) => setForm({ ...form, role: e.target.value as Role })}>
              {assignableRoles.map((r) => (
                <option key={r} value={r}>{ROLE_LABELS[r]}</option>
              ))}
            </select>
          </Field>
          {mode === 'create' && (
            <Field label="Kata sandi awal" hint={PASSWORD_HINT}>
              <input type="text" className={inputClass} value={form.password} onChange={(e) => setForm({ ...form, password: e.target.value })} />
            </Field>
          )}
          {mode === 'edit' && target?.id !== me?.id && (
            <label className="flex items-center gap-2 text-sm">
              <input type="checkbox" checked={form.is_active} onChange={(e) => setForm({ ...form, is_active: e.target.checked })} />
              Akun aktif (bisa login)
            </label>
          )}
        </div>
      </Modal>

      <Modal
        open={mode === 'password'}
        onClose={close}
        size="sm"
        title="Reset kata sandi"
        subtitle={`Untuk ${target?.full_name} (${target?.phone_number}). Pastikan identitasnya sudah diverifikasi.`}
        footer={
          <>
            <button className={btnSecondary} onClick={close}>Batal</button>
            <button className={btnPrimary} disabled={busy} onClick={submit}>Reset</button>
          </>
        }
      >
        <Field label="Kata sandi baru" hint={PASSWORD_HINT}>
          <input type="text" className={inputClass} value={form.password} onChange={(e) => setForm({ ...form, password: e.target.value })} />
        </Field>
      </Modal>

      <Modal
        open={mode === 'b2b'}
        onClose={close}
        size="sm"
        title="Verifikasi mitra toko"
        subtitle={target?.b2b_profile ? `${target.b2b_profile.store_name} — ${target.b2b_profile.store_address}` : undefined}
        footer={
          <>
            <button className={btnSecondary} onClick={close}>Batal</button>
            <button className={btnPrimary} disabled={busy} onClick={submit}>Simpan</button>
          </>
        }
      >
        <div className="grid gap-3">
          <Field label="Keputusan">
            <select className={inputClass} value={b2b.status} onChange={(e) => setB2b({ ...b2b, status: e.target.value as typeof b2b.status })}>
              <option value="APPROVED">Setujui — dapat harga grosir</option>
              <option value="REJECTED">Tolak</option>
              <option value="PENDING">Tunda</option>
            </select>
          </Field>
          {b2b.status === 'APPROVED' && (
            <Field label="Limit kasbon (Rp)" hint={`Isi 0 jika tidak diberi kasbon. Saat ini: ${formatCurrency(target?.b2b_profile?.approved_credit_limit || 0)}`}>
              <input type="number" min={0} className={inputClass} value={b2b.credit_limit} onChange={(e) => setB2b({ ...b2b, credit_limit: e.target.value })} />
            </Field>
          )}
        </div>
      </Modal>

      <Modal
        open={mode === 'delete'}
        onClose={close}
        size="sm"
        title={`Hapus ${target?.full_name}?`}
        subtitle="Sebaiknya nonaktifkan akun saja agar riwayat pesanan tetap rapi. Hapus hanya untuk akun yang salah dibuat."
        footer={
          <>
            <button className={btnSecondary} onClick={close}>Batal</button>
            <button className={btnDanger} disabled={busy} onClick={submit}>Hapus</button>
          </>
        }
      >
        <p className="text-sm text-slate-600 dark:text-slate-400">Tindakan ini tidak bisa dibatalkan dari dashboard.</p>
      </Modal>
    </div>
  );
};
