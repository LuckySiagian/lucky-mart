import React, { useState } from 'react';
import { PageHeader } from '../../components/common/PageHeader';
import { formatCurrency } from '../../lib/utils';
import { Store, Edit, CreditCard, ShieldCheck, DollarSign, Plus, X, Search, CheckCircle2, Clock, Trash2 } from 'lucide-react';
import { toast } from 'sonner';

interface B2BLedger {
  id: string;
  mitraName: string;
  ownerName: string;
  phone: string;
  creditLimit: number;
  usedCredit: number;
  remainingLimit: number;
  dueDate: string;
  status: 'TEMPO_LANCAR' | 'JATUH_TEMPO' | 'OVERDUE';
}

export const CreditLedgerPage: React.FC = () => {
  const [ledgers, setLedgers] = useState<B2BLedger[]>([]);
  const [search, setSearch] = useState('');

  // Modal Limit Adjustment State
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedMitra, setSelectedMitra] = useState<B2BLedger | null>(null);
  const [newLimit, setNewLimit] = useState<string>('50000000');

  // Modal Register B2B State
  const [isAddOpen, setIsAddOpen] = useState(false);
  const [mitraName, setMitraName] = useState('');
  const [ownerName, setOwnerName] = useState('');
  const [phone, setPhone] = useState('');
  const [initLimit, setInitLimit] = useState('50000000');

  const handleOpenModal = (mitra: B2BLedger) => {
    setSelectedMitra(mitra);
    setNewLimit(mitra.creditLimit.toString());
    setIsModalOpen(true);
  };

  const handleSaveLimit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedMitra) return;
    const limit = parseFloat(newLimit);
    if (isNaN(limit) || limit < 0) {
      toast.error('Batas limit kredit tidak valid.');
      return;
    }

    setLedgers((prev) =>
      prev.map((l) => {
        if (l.id === selectedMitra.id) {
          const remaining = Math.max(0, limit - l.usedCredit);
          return {
            ...l,
            creditLimit: limit,
            remainingLimit: remaining,
          };
        }
        return l;
      })
    );

    toast.success(`Plafon limit kredit untuk "${selectedMitra.mitraName}" diperbarui menjadi ${formatCurrency(limit)}!`);
    setIsModalOpen(false);
  };

  const handleRepay = (mitra: B2BLedger) => {
    setLedgers((prev) =>
      prev.map((l) => {
        if (l.id === mitra.id) {
          return {
            ...l,
            usedCredit: 0,
            remainingLimit: l.creditLimit,
            status: 'TEMPO_LANCAR',
          };
        }
        return l;
      })
    );
    toast.success(`Pelunasan Kasbon sebesar ${formatCurrency(mitra.usedCredit)} oleh "${mitra.mitraName}" berhasil dicatat!`);
  };

  const handleDeleteMitra = (id: string, name: string) => {
    if (window.confirm(`Hapus permanen kemitraan B2B "${name}"?`)) {
      setLedgers((prev) => prev.filter((l) => l.id !== id));
      toast.success(`Kemitraan B2B "${name}" berhasil dihapus!`);
    }
  };

  const handleRegisterB2B = (e: React.FormEvent) => {
    e.preventDefault();
    if (!mitraName.trim() || !phone.trim()) {
      toast.error('Nama toko mitra dan nomor HP wajib diisi.');
      return;
    }

    const limit = parseFloat(initLimit) || 50000000;
    const newMitra: B2BLedger = {
      id: `LED-${Date.now()}`,
      mitraName: mitraName.trim(),
      ownerName: ownerName.trim() || 'Pemilik Toko',
      phone: phone.trim(),
      creditLimit: limit,
      usedCredit: 0,
      remainingLimit: limit,
      dueDate: '2026-08-26 (Sisa 30 Hari)',
      status: 'TEMPO_LANCAR',
    };

    setLedgers((prev) => [...prev, newMitra]);
    toast.success(`Mitra B2B Baru "${mitraName}" berhasil terdaftar dengan limit ${formatCurrency(limit)}!`);
    setIsAddOpen(false);
    setMitraName('');
    setOwnerName('');
    setPhone('');
  };

  const filtered = ledgers.filter(
    (l) =>
      l.mitraName.toLowerCase().includes(search.toLowerCase()) ||
      l.ownerName.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className="space-y-6">
      <PageHeader
        title="B2B Credit Ledger (Mitra Kasbon Tempo 30 Hari)"
        subtitle="Fasilitas pembiayaan kasbon tempo 30 hari tanpa bunga khusus Toko Kelontong Mitra Warung."
      >
        <button
          onClick={() => setIsAddOpen(true)}
          className="inline-flex items-center space-x-2 rounded-xl bg-brand-600 px-4 py-2.5 text-sm font-semibold text-white hover:bg-brand-700 transition-colors shadow-lg shadow-brand-600/20"
        >
          <Plus className="h-4 w-4" />
          <span>Registrasi Mitra B2B Baru</span>
        </button>
      </PageHeader>

      {/* Filter & Search Bar */}
      <div className="flex flex-col sm:flex-row items-center justify-between gap-4 bg-white dark:bg-slate-900 p-4 rounded-2xl border border-slate-200 dark:border-slate-800">
        <div className="relative w-full sm:w-80">
          <Search className="absolute left-3.5 top-3 h-4 w-4 text-slate-400" />
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Cari nama toko mitra atau pemilik..."
            className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 pl-10 pr-4 py-2 text-sm text-slate-900 dark:text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-brand-500"
          />
        </div>
        <div className="text-xs text-slate-500 font-semibold">
          Total Piutang Aktif B2B:{' '}
          <span className="text-rose-600 font-bold">
            {formatCurrency(ledgers.reduce((acc, l) => acc + l.usedCredit, 0))}
          </span>
        </div>
      </div>

      {/* B2B Table */}
      <div className="overflow-x-auto rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 shadow-sm">
        <table className="w-full text-left text-sm text-slate-600 dark:text-slate-400">
          <thead className="bg-slate-50 dark:bg-slate-800/50 text-xs uppercase text-slate-500 font-bold border-b border-slate-200 dark:border-slate-800">
            <tr>
              <th className="px-6 py-4">Mitra Toko Kelontong B2B</th>
              <th className="px-6 py-4">Plafon Limit Kredit</th>
              <th className="px-6 py-4">Kredit Terpakai</th>
              <th className="px-6 py-4">Sisa Limit Tersedia</th>
              <th className="px-6 py-4">Status & Jatuh Tempo</th>
              <th className="px-6 py-4 text-right">Aksi Kelola (CRUD)</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-200 dark:divide-slate-800">
            {filtered.map((l) => (
              <tr key={l.id} className="hover:bg-slate-50/80 dark:hover:bg-slate-800/30">
                <td className="px-6 py-4">
                  <div className="flex items-center space-x-3">
                    <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-brand-100 text-brand-700 dark:bg-brand-950/60 font-bold">
                      <Store className="h-5 w-5" />
                    </div>
                    <div>
                      <div className="font-bold text-slate-900 dark:text-white">{l.mitraName}</div>
                      <div className="text-xs text-slate-400">Pemilik: {l.ownerName} ({l.phone})</div>
                    </div>
                  </div>
                </td>
                <td className="px-6 py-4 font-bold text-slate-900 dark:text-white">{formatCurrency(l.creditLimit)}</td>
                <td className="px-6 py-4 font-black text-rose-600">{formatCurrency(l.usedCredit)}</td>
                <td className="px-6 py-4 font-black text-emerald-600">{formatCurrency(l.remainingLimit)}</td>
                <td className="px-6 py-4">
                  <div className="space-y-1">
                    {l.status === 'TEMPO_LANCAR' && (
                      <span className="inline-flex items-center space-x-1 rounded-full bg-emerald-100 dark:bg-emerald-950/50 px-2.5 py-0.5 text-xs font-bold text-emerald-700 dark:text-emerald-400">
                        <CheckCircle2 className="h-3.5 w-3.5" />
                        <span>TEMPO LANCAR</span>
                      </span>
                    )}
                    {l.status === 'JATUH_TEMPO' && (
                      <span className="inline-flex items-center space-x-1 rounded-full bg-amber-100 dark:bg-amber-950/50 px-2.5 py-0.5 text-xs font-bold text-amber-700 dark:text-amber-400">
                        <Clock className="h-3.5 w-3.5" />
                        <span>SEGERA JATUH TEMPO</span>
                      </span>
                    )}
                    <div className="text-[11px] text-slate-400 font-medium">{l.dueDate}</div>
                  </div>
                </td>
                <td className="px-6 py-4 text-right space-x-2">
                  <button
                    onClick={() => handleOpenModal(l)}
                    className="inline-flex items-center space-x-1 rounded-lg bg-slate-100 dark:bg-slate-800 px-2.5 py-1 text-xs font-bold text-brand-600 hover:bg-brand-600 hover:text-white transition-colors"
                  >
                    <Edit className="h-3.5 w-3.5" />
                    <span>Set Limit</span>
                  </button>
                  {l.usedCredit > 0 && (
                    <button
                      onClick={() => handleRepay(l)}
                      className="inline-flex items-center space-x-1 rounded-lg bg-emerald-600 px-2.5 py-1 text-xs font-bold text-white hover:bg-emerald-700"
                    >
                      <DollarSign className="h-3.5 w-3.5" />
                      <span>Pelunasan</span>
                    </button>
                  )}
                  <button
                    onClick={() => handleDeleteMitra(l.id, l.mitraName)}
                    className="inline-flex items-center space-x-1 rounded-lg bg-rose-100 text-rose-700 dark:bg-rose-950/50 dark:text-rose-400 px-2.5 py-1 text-xs font-bold hover:bg-rose-200"
                  >
                    <Trash2 className="h-3.5 w-3.5" />
                    <span>Hapus</span>
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Modal Adjust Limit */}
      {isModalOpen && selectedMitra && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 backdrop-blur-sm p-4">
          <div className="relative w-full max-w-md rounded-3xl bg-white dark:bg-slate-900 p-6 shadow-2xl border border-slate-200 dark:border-slate-800 space-y-6">
            <div className="flex items-center justify-between border-b border-slate-200 dark:border-slate-800 pb-4">
              <div className="flex items-center space-x-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-brand-100 text-brand-600">
                  <CreditCard className="h-5 w-5" />
                </div>
                <div>
                  <h3 className="text-lg font-bold text-slate-900 dark:text-white">Ubah Limit Kasbon B2B</h3>
                  <p className="text-xs text-slate-500">{selectedMitra.mitraName}</p>
                </div>
              </div>
              <button onClick={() => setIsModalOpen(false)} className="rounded-xl p-2 text-slate-400 hover:bg-slate-100">
                <X className="h-5 w-5" />
              </button>
            </div>

            <form onSubmit={handleSaveLimit} className="space-y-4">
              <div>
                <label className="block text-xs font-bold text-slate-500 mb-1">Limit Plafon Baru (Rp) *</label>
                <input
                  type="number"
                  required
                  min="0"
                  value={newLimit}
                  onChange={(e) => setNewLimit(e.target.value)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white font-bold"
                />
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
                  Simpan Limit Baru
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Modal Registrasi Mitra B2B Baru */}
      {isAddOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 backdrop-blur-sm p-4">
          <div className="relative w-full max-w-md rounded-3xl bg-white dark:bg-slate-900 p-6 shadow-2xl border border-slate-200 dark:border-slate-800 space-y-6">
            <div className="flex items-center justify-between border-b border-slate-200 dark:border-slate-800 pb-4">
              <div className="flex items-center space-x-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-brand-100 text-brand-600">
                  <Plus className="h-5 w-5" />
                </div>
                <div>
                  <h3 className="text-lg font-bold text-slate-900 dark:text-white">Registrasi Mitra B2B Baru</h3>
                  <p className="text-xs text-slate-500">Pendaftaran Toko Kelontong Mitra Kasbon Tempo 30 Hari.</p>
                </div>
              </div>
              <button onClick={() => setIsAddOpen(false)} className="rounded-xl p-2 text-slate-400 hover:bg-slate-100">
                <X className="h-5 w-5" />
              </button>
            </div>

            <form onSubmit={handleRegisterB2B} className="space-y-4">
              <div>
                <label className="block text-xs font-bold text-slate-500 mb-1">Nama Toko Kelontong / Usaha *</label>
                <input
                  type="text"
                  required
                  placeholder="Contoh: Toko Sembako Barokah"
                  value={mitraName}
                  onChange={(e) => setMitraName(e.target.value)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white font-bold"
                />
              </div>

              <div>
                <label className="block text-xs font-bold text-slate-500 mb-1">Nama Pemilik Toko *</label>
                <input
                  type="text"
                  required
                  placeholder="Contoh: Hj. Aminah"
                  value={ownerName}
                  onChange={(e) => setOwnerName(e.target.value)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white"
                />
              </div>

              <div>
                <label className="block text-xs font-bold text-slate-500 mb-1">Nomor Handphone / WhatsApp *</label>
                <input
                  type="text"
                  required
                  placeholder="0812-xxxx-xxxx"
                  value={phone}
                  onChange={(e) => setPhone(e.target.value)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white font-mono"
                />
              </div>

              <div>
                <label className="block text-xs font-bold text-slate-500 mb-1">Plafon Limit Kasbon Awal (Rp) *</label>
                <input
                  type="number"
                  required
                  min="1000000"
                  value={initLimit}
                  onChange={(e) => setInitLimit(e.target.value)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white font-bold"
                />
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
                  Daftarkan Mitra B2B
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
