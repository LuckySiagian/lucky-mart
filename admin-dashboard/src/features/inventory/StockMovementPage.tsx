import React, { useState } from 'react';
import { PageHeader } from '../../components/common/PageHeader';
import { FileSpreadsheet, ArrowUpRight, ArrowDownRight, RefreshCw, Search, Filter, Plus, Trash2, X } from 'lucide-react';
import { toast } from 'sonner';

interface MovementLog {
  id: string;
  time: string;
  ref: string;
  productName: string;
  unit: string;
  type: 'IN' | 'OUT' | 'ADJUSTMENT';
  qtyChange: number;
  stockBefore: number;
  stockAfter: number;
  operator: string;
  notes: string;
}

export const StockMovementPage: React.FC = () => {
  const [movements, setMovements] = useState<MovementLog[]>([]);
  const [search, setSearch] = useState('');
  const [filterType, setFilterType] = useState<string>('ALL');

  // Add Log Modal State
  const [isAddOpen, setIsAddOpen] = useState(false);
  const [productName, setProductName] = useState('Beras Pandan Wangi Super 10kg');
  const [moveType, setMoveType] = useState<'IN' | 'OUT' | 'ADJUSTMENT'>('IN');
  const [qty, setQty] = useState('10');
  const [notes, setNotes] = useState('');

  const handleDeleteLog = (id: string, ref: string) => {
    if (window.confirm(`Hapus catatan mutasi stok "${ref}" dari ledger audit?`)) {
      setMovements((prev) => prev.filter((m) => m.id !== id));
      toast.success(`Catatan mutasi stok "${ref}" berhasil dihapus!`);
    }
  };

  const handleCreateMovement = (e: React.FormEvent) => {
    e.preventDefault();
    const qtyVal = parseInt(qty, 10) || 1;
    const qtyChange = moveType === 'OUT' ? -qtyVal : qtyVal;
    const refCode = `MOV-${new Date().toISOString().slice(0, 10).replace(/-/g, '')}-${Math.floor(1000 + Math.random() * 9000)}`;

    const newMov: MovementLog = {
      id: `MOV-${Date.now()}`,
      time: new Date().toISOString().replace('T', ' ').slice(0, 19),
      ref: refCode,
      productName: productName,
      unit: 'Unit',
      type: moveType,
      qtyChange: qtyChange,
      stockBefore: 50,
      stockAfter: 50 + qtyChange,
      operator: 'Super Admin Sembako',
      notes: notes.trim() || 'Pencatatan Mutasi Stok Manual Admin',
    };

    setMovements((prev) => [newMov, ...prev]);
    toast.success(`Mutasi stok baru "${refCode}" berhasil dicatat!`);
    setIsAddOpen(false);
    setNotes('');
  };

  const filtered = movements.filter((m) => {
    const matchesSearch =
      m.productName.toLowerCase().includes(search.toLowerCase()) ||
      m.ref.toLowerCase().includes(search.toLowerCase()) ||
      m.notes.toLowerCase().includes(search.toLowerCase());
    const matchesFilter = filterType === 'ALL' || m.type === filterType;
    return matchesSearch && matchesFilter;
  });

  return (
    <div className="space-y-6">
      <PageHeader
        title="Mutasi & Log Pergerakan Stok"
        subtitle="Jejak audit lengkap (Audit Trail) seluruh transaksi barang masuk, barang keluar, dan penyesuaian stok."
      >
        <button
          onClick={() => setIsAddOpen(true)}
          className="inline-flex items-center space-x-2 rounded-xl bg-brand-600 px-4 py-2.5 text-sm font-semibold text-white hover:bg-brand-700 transition-colors shadow-lg shadow-brand-600/20"
        >
          <Plus className="h-4 w-4" />
          <span>Catat Mutasi Stok Manual</span>
        </button>
      </PageHeader>

      {/* Filter Bar */}
      <div className="flex flex-col sm:flex-row items-center justify-between gap-4 bg-white dark:bg-slate-900 p-4 rounded-2xl border border-slate-200 dark:border-slate-800">
        <div className="relative w-full sm:w-80">
          <Search className="absolute left-3.5 top-3 h-4 w-4 text-slate-400" />
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Cari produk, ref, atau catatan..."
            className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 pl-10 pr-4 py-2 text-sm text-slate-900 dark:text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-brand-500"
          />
        </div>

        <div className="flex items-center space-x-2 w-full sm:w-auto">
          <Filter className="h-4 w-4 text-slate-400" />
          <select
            value={filterType}
            onChange={(e) => setFilterType(e.target.value)}
            className="rounded-xl bg-slate-100 dark:bg-slate-800 px-3 py-2 text-xs font-bold text-slate-700 dark:text-slate-200 focus:outline-none"
          >
            <option value="ALL">Semua Jenis Mutasi</option>
            <option value="IN">Stok Masuk (IN)</option>
            <option value="OUT">Stok Keluar (OUT)</option>
            <option value="ADJUSTMENT">Penyesuaian (ADJUSTMENT)</option>
          </select>
        </div>
      </div>

      {/* Log Table */}
      <div className="overflow-x-auto rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 shadow-sm">
        <table className="w-full text-left text-sm text-slate-600 dark:text-slate-400">
          <thead className="bg-slate-50 dark:bg-slate-800/50 text-xs uppercase text-slate-500 font-bold border-b border-slate-200 dark:border-slate-800">
            <tr>
              <th className="px-6 py-4">Waktu Transaksi</th>
              <th className="px-6 py-4">No. Referensi</th>
              <th className="px-6 py-4">Produk Sembako</th>
              <th className="px-6 py-4">Tipe Mutasi</th>
              <th className="px-6 py-4">Perubahan Qty</th>
              <th className="px-6 py-4">Stok Sebelum & Sesudah</th>
              <th className="px-6 py-4">Operator & Catatan</th>
              <th className="px-6 py-4 text-right">Aksi (CRUD)</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-200 dark:divide-slate-800">
            {filtered.map((m) => (
              <tr key={m.id} className="hover:bg-slate-50/80 dark:hover:bg-slate-800/30">
                <td className="px-6 py-4 font-mono text-xs text-slate-500">{m.time}</td>
                <td className="px-6 py-4 font-mono text-xs font-bold text-brand-600 dark:text-brand-400">{m.ref}</td>
                <td className="px-6 py-4">
                  <div className="font-bold text-slate-900 dark:text-white">{m.productName}</div>
                  <div className="text-xs text-slate-400">Satuan: {m.unit}</div>
                </td>
                <td className="px-6 py-4">
                  {m.type === 'IN' && (
                    <span className="inline-flex items-center space-x-1 rounded-full bg-emerald-100 dark:bg-emerald-950/50 px-2.5 py-0.5 text-xs font-bold text-emerald-700 dark:text-emerald-400">
                      <ArrowDownRight className="h-3.5 w-3.5" />
                      <span>STOK MASUK</span>
                    </span>
                  )}
                  {m.type === 'OUT' && (
                    <span className="inline-flex items-center space-x-1 rounded-full bg-rose-100 dark:bg-rose-950/50 px-2.5 py-0.5 text-xs font-bold text-rose-700 dark:text-rose-400">
                      <ArrowUpRight className="h-3.5 w-3.5" />
                      <span>STOK KELUAR</span>
                    </span>
                  )}
                  {m.type === 'ADJUSTMENT' && (
                    <span className="inline-flex items-center space-x-1 rounded-full bg-amber-100 dark:bg-amber-950/50 px-2.5 py-0.5 text-xs font-bold text-amber-700 dark:text-amber-400">
                      <RefreshCw className="h-3.5 w-3.5" />
                      <span>ADJUSTMENT</span>
                    </span>
                  )}
                </td>
                <td className="px-6 py-4">
                  <span className={`font-black text-sm ${m.qtyChange > 0 ? 'text-emerald-600' : 'text-rose-600'}`}>
                    {m.qtyChange > 0 ? `+${m.qtyChange}` : m.qtyChange}
                  </span>
                </td>
                <td className="px-6 py-4 text-xs font-semibold">
                  <span className="text-slate-400">{m.stockBefore}</span> →{' '}
                  <span className="text-slate-900 dark:text-white font-bold">{m.stockAfter} Unit</span>
                </td>
                <td className="px-6 py-4 text-xs">
                  <div className="font-bold text-slate-700 dark:text-slate-300">{m.operator}</div>
                  <div className="text-slate-400 italic max-w-xs truncate">{m.notes}</div>
                </td>
                <td className="px-6 py-4 text-right">
                  <button
                    onClick={() => handleDeleteLog(m.id, m.ref)}
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

      {/* Modal Catat Mutasi Manual */}
      {isAddOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 backdrop-blur-sm p-4">
          <div className="relative w-full max-w-md rounded-3xl bg-white dark:bg-slate-900 p-6 shadow-2xl border border-slate-200 dark:border-slate-800 space-y-6">
            <div className="flex items-center justify-between border-b border-slate-200 dark:border-slate-800 pb-4">
              <div className="flex items-center space-x-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-brand-100 text-brand-600">
                  <Plus className="h-5 w-5" />
                </div>
                <div>
                  <h3 className="text-lg font-bold text-slate-900 dark:text-white">Catat Mutasi Stok Manual</h3>
                  <p className="text-xs text-slate-500">Pencatatan barang masuk/keluar ke audit ledger.</p>
                </div>
              </div>
              <button onClick={() => setIsAddOpen(false)} className="rounded-xl p-2 text-slate-400 hover:bg-slate-100">
                <X className="h-5 w-5" />
              </button>
            </div>

            <form onSubmit={handleCreateMovement} className="space-y-4">
              <div>
                <label className="block text-xs font-bold text-slate-500 mb-1">Pilih Produk Sembako *</label>
                <select
                  value={productName}
                  onChange={(e) => setProductName(e.target.value)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white font-bold"
                >
                  <option value="Beras Pandan Wangi Super 10kg">Beras Pandan Wangi Super 10kg</option>
                  <option value="Minyak Goreng Bimoli Refill 2L">Minyak Goreng Bimoli Refill 2L</option>
                  <option value="Telur Ayam Negeri Super (1kg)">Telur Ayam Negeri Super (1kg)</option>
                  <option value="Indomie Goreng Spesial (Karton / 40 Pcs)">Indomie Goreng Spesial (Karton)</option>
                </select>
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs font-bold text-slate-500 mb-1">Tipe Mutasi *</label>
                  <select
                    value={moveType}
                    onChange={(e) => setMoveType(e.target.value as any)}
                    className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-xs text-slate-900 dark:text-white font-bold"
                  >
                    <option value="IN">Stok Masuk (+ IN)</option>
                    <option value="OUT">Stok Keluar (- OUT)</option>
                    <option value="ADJUSTMENT">Adjustment (Koreksi)</option>
                  </select>
                </div>

                <div>
                  <label className="block text-xs font-bold text-slate-500 mb-1">Jumlah (Qty) *</label>
                  <input
                    type="number"
                    required
                    min="1"
                    value={qty}
                    onChange={(e) => setQty(e.target.value)}
                    className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white font-bold"
                  />
                </div>
              </div>

              <div>
                <label className="block text-xs font-bold text-slate-500 mb-1">Catatan / Alasan Mutasi</label>
                <input
                  type="text"
                  placeholder="Contoh: Penerimaan Pasokan Supplier Baru"
                  value={notes}
                  onChange={(e) => setNotes(e.target.value)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white"
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
                  Catat Mutasi Stok
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
