import React, { useState } from 'react';
import { PageHeader } from '../../components/common/PageHeader';
import { formatCurrency } from '../../lib/utils';
import { CreditCard, CheckCircle, XCircle, Image, Search, ShieldCheck, X, Plus, Trash2 } from 'lucide-react';
import { toast } from 'sonner';

interface PaymentRecord {
  id: string;
  paymentNumber: string;
  orderNumber: string;
  customerName: string;
  method: string;
  amount: number;
  date: string;
  status: 'PENDING_VERIFICATION' | 'VERIFIED' | 'REJECTED';
  proofUrl: string;
  notes?: string;
}

export const PaymentListPage: React.FC = () => {
  const [payments, setPayments] = useState<PaymentRecord[]>([]);
  const [search, setSearch] = useState('');

  // Preview Modal State
  const [selectedPayment, setSelectedPayment] = useState<PaymentRecord | null>(null);

  // Manual Add Payment State
  const [isAddOpen, setIsAddOpen] = useState(false);
  const [payOrderNum, setPayOrderNum] = useState('');
  const [payCustName, setPayCustName] = useState('');
  const [payMethod, setPayMethod] = useState('Transfer Bank Mandiri (Manual)');
  const [payAmount, setPayAmount] = useState('');

  const handleApprove = (id: string) => {
    setPayments((prev) =>
      prev.map((p) => (p.id === id ? { ...p, status: 'VERIFIED' } : p))
    );
    toast.success('Pembayaran berhasil diverifikasi & dikonfirmasi!');
    setSelectedPayment(null);
  };

  const handleReject = (id: string) => {
    setPayments((prev) =>
      prev.map((p) => (p.id === id ? { ...p, status: 'REJECTED' } : p))
    );
    toast.error('Pembayaran ditolak. Pelanggan akan diminta mengunggah ulang bukti transfer.');
    setSelectedPayment(null);
  };

  const handleDeletePayment = (id: string, num: string) => {
    if (window.confirm(`Hapus permanen data pembayaran "${num}"?`)) {
      setPayments((prev) => prev.filter((p) => p.id !== id));
      toast.success(`Data pembayaran "${num}" berhasil dihapus!`);
      if (selectedPayment?.id === id) setSelectedPayment(null);
    }
  };

  const handleCreatePayment = (e: React.FormEvent) => {
    e.preventDefault();
    if (!payCustName.trim() || !payAmount) {
      toast.error('Nama pelanggan dan nominal pembayaran wajib diisi.');
      return;
    }

    const payNum = `PAY-${new Date().toISOString().slice(0, 10).replace(/-/g, '')}-${Math.floor(1000 + Math.random() * 9000)}`;
    const newPay: PaymentRecord = {
      id: `PAY-${Date.now()}`,
      paymentNumber: payNum,
      orderNumber: payOrderNum.trim() || `ORD-${Math.floor(1000 + Math.random() * 9000)}`,
      customerName: payCustName.trim(),
      method: payMethod,
      amount: parseFloat(payAmount),
      date: new Date().toISOString().replace('T', ' ').slice(0, 16),
      status: 'VERIFIED',
      proofUrl: 'https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=500&auto=format&fit=crop',
      notes: 'Pencatatan Pembayaran Manual oleh Admin',
    };

    setPayments((prev) => [newPay, ...prev]);
    toast.success(`Pembayaran manual "${payNum}" berhasil dicatat!`);
    setIsAddOpen(false);
    setPayCustName('');
    setPayAmount('');
  };

  const filtered = payments.filter(
    (p) =>
      p.paymentNumber.toLowerCase().includes(search.toLowerCase()) ||
      p.orderNumber.toLowerCase().includes(search.toLowerCase()) ||
      p.customerName.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className="space-y-6">
      <PageHeader
        title="Manajemen Pembayaran & Keuangan"
        subtitle="Verifikasi bukti transfer manual, transaksi QRIS, dan pencatatan B2B Credit Tempo."
      >
        <button
          onClick={() => setIsAddOpen(true)}
          className="inline-flex items-center space-x-2 rounded-xl bg-brand-600 px-4 py-2.5 text-sm font-semibold text-white hover:bg-brand-700 transition-colors shadow-lg shadow-brand-600/20"
        >
          <Plus className="h-4 w-4" />
          <span>Catat Pembayaran Manual</span>
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
            placeholder="Cari no. payment, order, atau nama..."
            className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 pl-10 pr-4 py-2 text-sm text-slate-900 dark:text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-brand-500"
          />
        </div>
        <div className="text-xs text-slate-500 font-semibold">
          Menunggu Verifikasi: <span className="text-amber-600 font-bold">{payments.filter((p) => p.status === 'PENDING_VERIFICATION').length} Transaksi</span>
        </div>
      </div>

      {/* Payment Table */}
      <div className="overflow-x-auto rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 shadow-sm">
        <table className="w-full text-left text-sm text-slate-600 dark:text-slate-400">
          <thead className="bg-slate-50 dark:bg-slate-800/50 text-xs uppercase text-slate-500 font-bold border-b border-slate-200 dark:border-slate-800">
            <tr>
              <th className="px-6 py-4">No. Payment & Order</th>
              <th className="px-6 py-4">Nama Pelanggan</th>
              <th className="px-6 py-4">Metode Bayar</th>
              <th className="px-6 py-4">Nominal</th>
              <th className="px-6 py-4">Bukti Foto</th>
              <th className="px-6 py-4">Status</th>
              <th className="px-6 py-4 text-right">Aksi Kelola (CRUD)</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-200 dark:divide-slate-800">
            {filtered.map((p) => (
              <tr key={p.id} className="hover:bg-slate-50/80 dark:hover:bg-slate-800/30">
                <td className="px-6 py-4">
                  <div className="font-mono font-bold text-slate-900 dark:text-white">{p.paymentNumber}</div>
                  <div className="text-xs font-mono text-brand-600 dark:text-brand-400">{p.orderNumber}</div>
                </td>
                <td className="px-6 py-4 font-bold text-slate-800 dark:text-slate-200">{p.customerName}</td>
                <td className="px-6 py-4 text-xs font-semibold text-slate-600 dark:text-slate-400">{p.method}</td>
                <td className="px-6 py-4 font-black text-slate-900 dark:text-white">{formatCurrency(p.amount)}</td>
                <td className="px-6 py-4">
                  <button
                    onClick={() => setSelectedPayment(p)}
                    className="inline-flex items-center space-x-1.5 text-xs text-brand-600 font-bold hover:underline"
                  >
                    <Image className="h-4 w-4" />
                    <span>Lihat Bukti</span>
                  </button>
                </td>
                <td className="px-6 py-4">
                  {p.status === 'PENDING_VERIFICATION' && (
                    <span className="inline-flex items-center rounded-full bg-amber-100 dark:bg-amber-950/50 px-2.5 py-0.5 text-xs font-bold text-amber-700 dark:text-amber-400">
                      PERLU VERIFIKASI
                    </span>
                  )}
                  {p.status === 'VERIFIED' && (
                    <span className="inline-flex items-center rounded-full bg-emerald-100 dark:bg-emerald-950/50 px-2.5 py-0.5 text-xs font-bold text-emerald-700 dark:text-emerald-400">
                      DIVERIFIKASI (PAID)
                    </span>
                  )}
                  {p.status === 'REJECTED' && (
                    <span className="inline-flex items-center rounded-full bg-rose-100 dark:bg-rose-950/50 px-2.5 py-0.5 text-xs font-bold text-rose-700 dark:text-rose-400">
                      DITOLAK
                    </span>
                  )}
                </td>
                <td className="px-6 py-4 text-right space-x-2">
                  {p.status === 'PENDING_VERIFICATION' ? (
                    <div className="flex items-center justify-end space-x-2">
                      <button
                        onClick={() => handleApprove(p.id)}
                        className="inline-flex items-center space-x-1 rounded-xl bg-emerald-600 px-3 py-1.5 text-xs font-bold text-white hover:bg-emerald-700 shadow-md shadow-emerald-600/20"
                      >
                        <CheckCircle className="h-3.5 w-3.5" />
                        <span>Setujui</span>
                      </button>
                      <button
                        onClick={() => handleReject(p.id)}
                        className="inline-flex items-center space-x-1 rounded-xl bg-rose-100 text-rose-700 dark:bg-rose-950/50 dark:text-rose-400 px-3 py-1.5 text-xs font-bold hover:bg-rose-200"
                      >
                        <XCircle className="h-3.5 w-3.5" />
                        <span>Tolak</span>
                      </button>
                    </div>
                  ) : (
                    <button
                      onClick={() => handleDeletePayment(p.id, p.paymentNumber)}
                      className="inline-flex items-center space-x-1 rounded-lg bg-rose-100 text-rose-700 dark:bg-rose-950/50 dark:text-rose-400 px-2.5 py-1 text-xs font-bold hover:bg-rose-200"
                    >
                      <Trash2 className="h-3.5 w-3.5" />
                      <span>Hapus</span>
                    </button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Modal Preview Bukti Transfer */}
      {selectedPayment && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 backdrop-blur-sm p-4">
          <div className="relative w-full max-w-lg rounded-3xl bg-white dark:bg-slate-900 p-6 shadow-2xl border border-slate-200 dark:border-slate-800 space-y-6">
            <div className="flex items-center justify-between border-b border-slate-200 dark:border-slate-800 pb-4">
              <div className="flex items-center space-x-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-emerald-100 text-emerald-600">
                  <ShieldCheck className="h-5 w-5" />
                </div>
                <div>
                  <h3 className="text-lg font-bold text-slate-900 dark:text-white">Detail & Bukti Pembayaran</h3>
                  <p className="text-xs text-slate-500">{selectedPayment.paymentNumber}</p>
                </div>
              </div>
              <button onClick={() => setSelectedPayment(null)} className="rounded-xl p-2 text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800">
                <X className="h-5 w-5" />
              </button>
            </div>

            <div className="space-y-4">
              <div className="rounded-2xl bg-slate-100 dark:bg-slate-800 p-4 grid grid-cols-2 gap-4 text-xs">
                <div>
                  <span className="text-slate-400">Pelanggan:</span>
                  <div className="font-bold text-slate-900 dark:text-white mt-0.5">{selectedPayment.customerName}</div>
                </div>
                <div>
                  <span className="text-slate-400">Nominal Transfer:</span>
                  <div className="font-black text-brand-600 dark:text-brand-400 text-sm mt-0.5">{formatCurrency(selectedPayment.amount)}</div>
                </div>
              </div>

              <div>
                <label className="block text-xs font-bold text-slate-500 mb-2">Foto Bukti Transfer Pemesan:</label>
                <div className="rounded-2xl border border-slate-200 dark:border-slate-800 overflow-hidden max-h-64 flex items-center justify-center bg-slate-900/5 dark:bg-slate-800/50 p-2">
                  <img src={selectedPayment.proofUrl} alt="Bukti Transfer" className="max-h-60 w-auto rounded-xl object-contain shadow-md" />
                </div>
              </div>
            </div>

            <div className="flex items-center justify-end space-x-3 pt-4 border-t border-slate-200 dark:border-slate-800">
              {selectedPayment.status === 'PENDING_VERIFICATION' ? (
                <>
                  <button
                    type="button"
                    onClick={() => handleReject(selectedPayment.id)}
                    className="rounded-xl bg-rose-100 text-rose-700 px-4 py-2.5 text-xs font-bold hover:bg-rose-200"
                  >
                    Tolak Pembayaran
                  </button>
                  <button
                    type="button"
                    onClick={() => handleApprove(selectedPayment.id)}
                    className="rounded-xl bg-emerald-600 px-5 py-2.5 text-xs font-bold text-white hover:bg-emerald-700 shadow-lg shadow-emerald-600/20"
                  >
                    Verifikasi & Setujui
                  </button>
                </>
              ) : (
                <button
                  type="button"
                  onClick={() => setSelectedPayment(null)}
                  className="rounded-xl bg-slate-100 dark:bg-slate-800 px-4 py-2.5 text-xs font-bold text-slate-700 dark:text-slate-200"
                >
                  Tutup
                </button>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Modal Catat Pembayaran Manual */}
      {isAddOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 backdrop-blur-sm p-4">
          <div className="relative w-full max-w-md rounded-3xl bg-white dark:bg-slate-900 p-6 shadow-2xl border border-slate-200 dark:border-slate-800 space-y-6">
            <div className="flex items-center justify-between border-b border-slate-200 dark:border-slate-800 pb-4">
              <div className="flex items-center space-x-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-brand-100 text-brand-600">
                  <Plus className="h-5 w-5" />
                </div>
                <div>
                  <h3 className="text-lg font-bold text-slate-900 dark:text-white">Catat Pembayaran Manual</h3>
                  <p className="text-xs text-slate-500">Pencatatan pembayaran tunai / transfer langsung oleh admin.</p>
                </div>
              </div>
              <button onClick={() => setIsAddOpen(false)} className="rounded-xl p-2 text-slate-400 hover:bg-slate-100">
                <X className="h-5 w-5" />
              </button>
            </div>

            <form onSubmit={handleCreatePayment} className="space-y-4">
              <div>
                <label className="block text-xs font-bold text-slate-500 mb-1">Nama Pelanggan / Mitra *</label>
                <input
                  type="text"
                  required
                  placeholder="Contoh: Budi Santoso"
                  value={payCustName}
                  onChange={(e) => setPayCustName(e.target.value)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white font-bold"
                />
              </div>

              <div>
                <label className="block text-xs font-bold text-slate-500 mb-1">Nomor Order Terkait (Opsional)</label>
                <input
                  type="text"
                  placeholder="ORD-20260727-xxxx"
                  value={payOrderNum}
                  onChange={(e) => setPayOrderNum(e.target.value)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white font-mono"
                />
              </div>

              <div>
                <label className="block text-xs font-bold text-slate-500 mb-1">Metode Pembayaran *</label>
                <select
                  value={payMethod}
                  onChange={(e) => setPayMethod(e.target.value)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white font-bold"
                >
                  <option value="Transfer Bank Mandiri (Manual)">Transfer Bank Mandiri</option>
                  <option value="Tunai / Cash Kasir Toko">Tunai / Cash Kasir Toko</option>
                  <option value="QRIS Instant Payment">QRIS Instant Payment</option>
                </select>
              </div>

              <div>
                <label className="block text-xs font-bold text-slate-500 mb-1">Nominal Diterima (Rp) *</label>
                <input
                  type="number"
                  required
                  min="1"
                  placeholder="145000"
                  value={payAmount}
                  onChange={(e) => setPayAmount(e.target.value)}
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
                  Simpan Pembayaran
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
