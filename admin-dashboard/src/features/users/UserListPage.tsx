import React, { useState, useEffect } from 'react';
import { PageHeader } from '../../components/common/PageHeader';
import { Users, ShieldCheck, UserCheck, Edit, Plus, X, Search, Lock, UserPlus, Trash2 } from 'lucide-react';
import { toast } from 'sonner';
import { userService } from '../../services/apiServices';

interface UserAccount {
  id: string;
  name: string;
  phone: string;
  email: string;
  role: 'SUPER_ADMIN' | 'STORE_MANAGER' | 'B2B_CUSTOMER' | 'B2C_CUSTOMER';
  status: 'ACTIVE' | 'SUSPENDED';
}

const mockUsers: UserAccount[] = [
  {
    id: 'USR-001',
    name: 'Super Admin LuckyMart (Pemilik Utama)',
    phone: '081234567890',
    email: 'superadmin@luckymart.com',
    role: 'SUPER_ADMIN',
    status: 'ACTIVE',
  },
];

export const UserListPage: React.FC = () => {
  const [users, setUsers] = useState<UserAccount[]>(mockUsers);
  const [search, setSearch] = useState('');
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const fetchUsers = async () => {
      setLoading(true);
      try {
        const res = await userService.getUsers(1, 50);
        const list = Array.isArray(res.data) ? res.data : (res.data?.items || []);
        if (list.length > 0) {
          const mapped = list.map((u: any) => ({
            id: u.id,
            name: u.full_name || u.name || 'Pengguna LuckyMart',
            phone: u.phone_number || u.phone || '-',
            email: u.email || '-',
            role: u.role || 'B2C_CUSTOMER',
            status: u.is_active ? 'ACTIVE' : 'SUSPENDED',
          }));
          setUsers(mapped);
        }
      } catch (e) {
        // Keep mock user fallback
      } finally {
        setLoading(false);
      }
    };
    fetchUsers();
  }, []);

  // Edit Role Modal State
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState<UserAccount | null>(null);
  const [newRole, setNewRole] = useState<'B2B_CUSTOMER' | 'B2C_CUSTOMER'>('B2C_CUSTOMER');

  // Add Customer Modal State
  const [isAddOpen, setIsAddOpen] = useState(false);
  const [addName, setAddName] = useState('');
  const [addPhone, setAddPhone] = useState('');
  const [addEmail, setAddEmail] = useState('');
  const [addRole, setAddRole] = useState<'B2B_CUSTOMER' | 'B2C_CUSTOMER'>('B2C_CUSTOMER');

  const handleOpenEdit = (u: UserAccount) => {
    if (u.role === 'SUPER_ADMIN' || u.role === 'STORE_MANAGER') {
      toast.info('Role SuperAdmin dan Admin dibatasi 1 akun utama sesuai proteksi keamanan.');
      return;
    }
    setSelectedUser(u);
    setNewRole(u.role as any);
    setIsModalOpen(true);
  };

  const handleDeleteUser = (id: string, name: string) => {
    if (window.confirm(`Hapus permanen akun pengguna "${name}"?`)) {
      setUsers((prev) => prev.filter((u) => u.id !== id));
      toast.success(`Akun pengguna "${name}" berhasil dihapus!`);
    }
  };

  const handleSaveRole = (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedUser) return;

    setUsers((prev) =>
      prev.map((u) => (u.id === selectedUser.id ? { ...u, role: newRole } : u))
    );

    toast.success(`Role pengguna "${selectedUser.name}" berhasil diubah menjadi ${newRole}!`);
    setIsModalOpen(false);
  };

  const handleAddCustomer = (e: React.FormEvent) => {
    e.preventDefault();
    if (!addName.trim() || !addPhone.trim()) {
      toast.error('Nama dan No. HP wajib diisi.');
      return;
    }

    const newUser: UserAccount = {
      id: `USR-${Date.now()}`,
      name: addName.trim(),
      phone: addPhone.trim(),
      email: addEmail.trim() || `${addName.toLowerCase().replace(/\s+/g, '')}@gmail.com`,
      role: addRole,
      status: 'ACTIVE',
    };

    setUsers((prev) => [...prev, newUser]);
    toast.success(`Pelanggan baru "${addName}" berhasil didaftarkan!`);
    setIsAddOpen(false);
    setAddName('');
    setAddPhone('');
    setAddEmail('');
  };

  const filtered = users.filter(
    (u) =>
      u.name.toLowerCase().includes(search.toLowerCase()) ||
      u.email.toLowerCase().includes(search.toLowerCase()) ||
      u.phone.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className="space-y-6">
      <PageHeader
        title="Manajemen Pengguna & Hak Akses (RBAC)"
        subtitle="Sistem Hak Akses Proteksi: 1 SuperAdmin, 1 Admin Operasional, serta Manajemen Akun Pelanggan B2B Mitra Warung & B2C Eceran."
      >
        <button
          onClick={() => setIsAddOpen(true)}
          className="inline-flex items-center space-x-2 rounded-xl bg-brand-600 px-4 py-2.5 text-sm font-semibold text-white hover:bg-brand-700 transition-colors shadow-lg shadow-brand-600/20"
        >
          <UserPlus className="h-4 w-4" />
          <span>Tambah Pelanggan Baru</span>
        </button>
      </PageHeader>

      {/* Role Protection Info Box */}
      <div className="rounded-2xl border border-blue-200 dark:border-blue-950 bg-blue-50/50 dark:bg-blue-950/20 p-4 flex items-center justify-between">
        <div className="flex items-center space-x-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-blue-100 dark:bg-blue-900/50 text-blue-600">
            <Lock className="h-5 w-5" />
          </div>
          <div>
            <h4 className="text-xs font-bold text-blue-900 dark:text-blue-300">Proteksi Role Sistem Sembako</h4>
            <p className="text-xs text-blue-700 dark:text-blue-400">
              Sesuai aturan keamanan: Role <b>SUPER_ADMIN</b> dan <b>STORE_MANAGER (Admin)</b> dikunci masing-masing 1 akun utama. Penambahan pengguna baru diperuntukkan untuk <b>B2B Mitra Warung</b> & <b>B2C Customer</b>.
            </p>
          </div>
        </div>
      </div>

      {/* Filter & Search Bar */}
      <div className="flex flex-col sm:flex-row items-center justify-between gap-4 bg-white dark:bg-slate-900 p-4 rounded-2xl border border-slate-200 dark:border-slate-800">
        <div className="relative w-full sm:w-80">
          <Search className="absolute left-3.5 top-3 h-4 w-4 text-slate-400" />
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Cari nama, email, atau HP..."
            className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 pl-10 pr-4 py-2 text-sm text-slate-900 dark:text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-brand-500"
          />
        </div>
        <div className="text-xs text-slate-500 font-semibold">
          Total Pengguna: <span className="text-brand-600 font-bold">{users.length} Akun</span>
        </div>
      </div>

      {/* Users Table */}
      <div className="overflow-x-auto rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 shadow-sm">
        <table className="w-full text-left text-sm text-slate-600 dark:text-slate-400">
          <thead className="bg-slate-50 dark:bg-slate-800/50 text-xs uppercase text-slate-500 font-bold border-b border-slate-200 dark:border-slate-800">
            <tr>
              <th className="px-6 py-4">Nama Pengguna</th>
              <th className="px-6 py-4">Nomor Kontak</th>
              <th className="px-6 py-4">Email</th>
              <th className="px-6 py-4">Role Akses</th>
              <th className="px-6 py-4 text-right">Aksi Kelola (CRUD)</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-200 dark:divide-slate-800">
            {filtered.map((u) => (
              <tr key={u.id} className="hover:bg-slate-50/80 dark:hover:bg-slate-800/30">
                <td className="px-6 py-4">
                  <div className="font-bold text-slate-900 dark:text-white">{u.name}</div>
                  <div className="text-[11px] text-slate-400">ID: {u.id}</div>
                </td>
                <td className="px-6 py-4 font-mono text-xs font-semibold">{u.phone}</td>
                <td className="px-6 py-4 text-xs">{u.email}</td>
                <td className="px-6 py-4">
                  {u.role === 'SUPER_ADMIN' && (
                    <span className="inline-flex items-center space-x-1 rounded-full bg-purple-100 dark:bg-purple-950/60 px-3 py-1 text-xs font-bold text-purple-700 dark:text-purple-300">
                      <ShieldCheck className="h-3.5 w-3.5" />
                      <span>SUPER ADMIN (UTAMA)</span>
                    </span>
                  )}
                  {u.role === 'STORE_MANAGER' && (
                    <span className="inline-flex items-center space-x-1 rounded-full bg-blue-100 dark:bg-blue-950/60 px-3 py-1 text-xs font-bold text-blue-700 dark:text-blue-300">
                      <UserCheck className="h-3.5 w-3.5" />
                      <span>ADMIN OPERASIONAL</span>
                    </span>
                  )}
                  {u.role === 'B2B_CUSTOMER' && (
                    <span className="inline-flex items-center space-x-1 rounded-full bg-emerald-100 dark:bg-emerald-950/60 px-3 py-1 text-xs font-bold text-emerald-700 dark:text-emerald-300">
                      <span>B2B MITRA WARUNG</span>
                    </span>
                  )}
                  {u.role === 'B2C_CUSTOMER' && (
                    <span className="inline-flex items-center space-x-1 rounded-full bg-slate-100 dark:bg-slate-800 px-3 py-1 text-xs font-bold text-slate-600 dark:text-slate-400">
                      <span>B2C CUSTOMER ECERAN</span>
                    </span>
                  )}
                </td>
                <td className="px-6 py-4 text-right space-x-2">
                  {u.role === 'SUPER_ADMIN' || u.role === 'STORE_MANAGER' ? (
                    <span className="text-xs text-slate-400 italic font-semibold">Terkunci (Proteksi)</span>
                  ) : (
                    <>
                      <button
                        onClick={() => handleOpenEdit(u)}
                        className="inline-flex items-center space-x-1 text-xs font-bold text-brand-600 hover:text-brand-700 bg-brand-50 dark:bg-brand-950/40 px-2.5 py-1 rounded-lg"
                      >
                        <Edit className="h-3.5 w-3.5" />
                        <span>Edit Role</span>
                      </button>
                      <button
                        onClick={() => handleDeleteUser(u.id, u.name)}
                        className="inline-flex items-center space-x-1 rounded-lg bg-rose-100 text-rose-700 dark:bg-rose-950/50 dark:text-rose-400 px-2.5 py-1 text-xs font-bold hover:bg-rose-200"
                      >
                        <Trash2 className="h-3.5 w-3.5" />
                        <span>Hapus</span>
                      </button>
                    </>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Modal Edit Role */}
      {isModalOpen && selectedUser && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 backdrop-blur-sm p-4">
          <div className="relative w-full max-w-md rounded-3xl bg-white dark:bg-slate-900 p-6 shadow-2xl border border-slate-200 dark:border-slate-800 space-y-6">
            <div className="flex items-center justify-between border-b border-slate-200 dark:border-slate-800 pb-4">
              <div className="flex items-center space-x-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-brand-100 text-brand-600">
                  <Edit className="h-5 w-5" />
                </div>
                <div>
                  <h3 className="text-lg font-bold text-slate-900 dark:text-white">Ubah Role Pengguna</h3>
                  <p className="text-xs text-slate-500">{selectedUser.name}</p>
                </div>
              </div>
              <button onClick={() => setIsModalOpen(false)} className="rounded-xl p-2 text-slate-400 hover:bg-slate-100">
                <X className="h-5 w-5" />
              </button>
            </div>

            <form onSubmit={handleSaveRole} className="space-y-4">
              <div>
                <label className="block text-xs font-bold text-slate-500 mb-1">Pilih Role Akses Baru *</label>
                <select
                  value={newRole}
                  onChange={(e) => setNewRole(e.target.value as any)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white font-bold"
                >
                  <option value="B2B_CUSTOMER">B2B MITRA WARUNG (Kasbon Tempo 30 Hari & Harga Grosir)</option>
                  <option value="B2C_CUSTOMER">B2C CUSTOMER ECERAN (Pembayaran Langsung)</option>
                </select>
              </div>

              <div className="flex items-center justify-end space-x-3 pt-4 border-t border-slate-200 dark:border-slate-800">
                <button
                  type="button"
                  onClick={() => setIsModalOpen(false)}
                  className="rounded-xl px-4 py-2.5 text-xs font-bold text-slate-600 hover:bg-slate-100"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  className="rounded-xl bg-brand-600 px-5 py-2.5 text-xs font-bold text-white hover:bg-brand-700 shadow-lg shadow-brand-600/20"
                >
                  Simpan Perubahan Role
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Modal Tambah Pelanggan Baru */}
      {isAddOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 backdrop-blur-sm p-4">
          <div className="relative w-full max-w-md rounded-3xl bg-white dark:bg-slate-900 p-6 shadow-2xl border border-slate-200 dark:border-slate-800 space-y-6">
            <div className="flex items-center justify-between border-b border-slate-200 dark:border-slate-800 pb-4">
              <div className="flex items-center space-x-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-brand-100 text-brand-600">
                  <UserPlus className="h-5 w-5" />
                </div>
                <div>
                  <h3 className="text-lg font-bold text-slate-900 dark:text-white">Tambah Pelanggan Baru</h3>
                  <p className="text-xs text-slate-500">Daftarkan pelanggan B2B Mitra Warung atau B2C Eceran.</p>
                </div>
              </div>
              <button onClick={() => setIsAddOpen(false)} className="rounded-xl p-2 text-slate-400 hover:bg-slate-100">
                <X className="h-5 w-5" />
              </button>
            </div>

            <form onSubmit={handleAddCustomer} className="space-y-4">
              <div>
                <label className="block text-xs font-bold text-slate-500 mb-1">Nama Lengkap / Nama Toko *</label>
                <input
                  type="text"
                  required
                  placeholder="Contoh: Toko Kelontong Rejeki"
                  value={addName}
                  onChange={(e) => setAddName(e.target.value)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white font-bold"
                />
              </div>

              <div>
                <label className="block text-xs font-bold text-slate-500 mb-1">Nomor Handphone / WhatsApp *</label>
                <input
                  type="text"
                  required
                  placeholder="0812-xxxx-xxxx"
                  value={addPhone}
                  onChange={(e) => setAddPhone(e.target.value)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white font-mono"
                />
              </div>

              <div>
                <label className="block text-xs font-bold text-slate-500 mb-1">Email (Opsional)</label>
                <input
                  type="email"
                  placeholder="pelanggan@gmail.com"
                  value={addEmail}
                  onChange={(e) => setAddEmail(e.target.value)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white"
                />
              </div>

              <div>
                <label className="block text-xs font-bold text-slate-500 mb-1">Role Akses Pelanggan *</label>
                <select
                  value={addRole}
                  onChange={(e) => setAddRole(e.target.value as any)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white font-bold"
                >
                  <option value="B2B_CUSTOMER">B2B MITRA WARUNG (Harga Grosir & Kasbon Tempo 30 Hari)</option>
                  <option value="B2C_CUSTOMER">B2C CUSTOMER ECERAN (Pembayaran Tunai / Transfer / QRIS)</option>
                </select>
              </div>

              <div className="flex items-center justify-end space-x-3 pt-4 border-t border-slate-200 dark:border-slate-800">
                <button
                  type="button"
                  onClick={() => setIsAddOpen(false)}
                  className="rounded-xl px-4 py-2.5 text-xs font-bold text-slate-600 hover:bg-slate-100"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  className="rounded-xl bg-brand-600 px-5 py-2.5 text-xs font-bold text-white hover:bg-brand-700 shadow-lg shadow-brand-600/20"
                >
                  Daftarkan Pelanggan
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
