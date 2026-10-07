import React, { useState } from 'react';
import { PageHeader } from '../../components/common/PageHeader';
import { Bell, CheckCheck, ShieldAlert, ShoppingBag, Boxes, Send, Plus, X, Trash2 } from 'lucide-react';
import { toast } from 'sonner';

interface SystemNotification {
  id: string;
  title: string;
  message: string;
  type: 'ORDER' | 'STOCK' | 'PAYMENT' | 'BROADCAST';
  time: string;
  isRead: boolean;
  channel: string;
}

export const NotificationPage: React.FC = () => {
  const [notifications, setNotifications] = useState<SystemNotification[]>([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [broadcastTitle, setBroadcastTitle] = useState('');
  const [broadcastMessage, setBroadcastMessage] = useState('');
  const [targetAudience, setTargetAudience] = useState('ALL');

  const handleMarkAllRead = () => {
    setNotifications((prev) => prev.map((n) => ({ ...n, isRead: true })));
    toast.success('Semua notifikasi telah ditandai dibaca.');
  };

  const handleDeleteNotif = (id: string) => {
    setNotifications((prev) => prev.filter((n) => n.id !== id));
    toast.success('Notifikasi berhasil dihapus.');
  };

  const handleSendBroadcast = (e: React.FormEvent) => {
    e.preventDefault();
    if (!broadcastTitle.trim() || !broadcastMessage.trim()) {
      toast.error('Judul dan pesan pengumuman wajib diisi.');
      return;
    }

    const newBroadcast: SystemNotification = {
      id: `NOT-${Date.now()}`,
      title: `[SIARAN] ${broadcastTitle.trim()}`,
      message: broadcastMessage.trim(),
      type: 'BROADCAST',
      time: 'Baru Saja',
      isRead: false,
      channel: `BROADCAST (${targetAudience})`,
    };

    setNotifications((prev) => [newBroadcast, ...prev]);
    toast.success(`Pengumuman siaran berhasil dikirim ke target ${targetAudience}!`);
    setIsModalOpen(false);
  };

  return (
    <div className="space-y-6">
      <PageHeader
        title="Notifikasi & Alert Sistem"
        subtitle="Pusat kontrol alert transaksi real-time, stok kritis, serta fasilitas siaran pengumuman ke pembeli/mitra warung."
      >
        <div className="flex items-center space-x-3">
          <button
            onClick={() => setIsModalOpen(true)}
            className="inline-flex items-center space-x-2 rounded-xl bg-brand-600 px-4 py-2.5 text-xs font-semibold text-white hover:bg-brand-700 transition-colors shadow-lg shadow-brand-600/20"
          >
            <Send className="h-4 w-4" />
            <span>Kirim Siaran Pengumuman</span>
          </button>

          <button
            onClick={handleMarkAllRead}
            className="inline-flex items-center space-x-2 rounded-xl bg-slate-100 dark:bg-slate-800 px-4 py-2.5 text-xs font-semibold text-slate-700 dark:text-slate-300 hover:bg-slate-200"
          >
            <CheckCheck className="h-4 w-4" />
            <span>Tandai Semua Dibaca</span>
          </button>
        </div>
      </PageHeader>

      {/* Notification List */}
      <div className="space-y-3">
        {notifications.map((n) => (
          <div
            key={n.id}
            className={`flex items-start space-x-4 rounded-2xl border p-4 shadow-sm transition-all ${
              !n.isRead
                ? 'bg-brand-50/40 dark:bg-brand-950/20 border-brand-200 dark:border-brand-900/50'
                : 'bg-white dark:bg-slate-900 border-slate-200 dark:border-slate-800'
            }`}
          >
            <div className="flex h-10 w-10 items-center justify-center rounded-xl flex-shrink-0">
              {n.type === 'ORDER' && <div className="h-10 w-10 rounded-xl bg-blue-100 text-blue-600 flex items-center justify-center"><ShoppingBag className="h-5 w-5" /></div>}
              {n.type === 'STOCK' && <div className="h-10 w-10 rounded-xl bg-amber-100 text-amber-600 flex items-center justify-center"><Boxes className="h-5 w-5" /></div>}
              {n.type === 'PAYMENT' && <div className="h-10 w-10 rounded-xl bg-emerald-100 text-emerald-600 flex items-center justify-center"><Bell className="h-5 w-5" /></div>}
              {n.type === 'BROADCAST' && <div className="h-10 w-10 rounded-xl bg-purple-100 text-purple-600 flex items-center justify-center"><Send className="h-5 w-5" /></div>}
            </div>

            <div className="flex-1">
              <div className="flex items-center justify-between">
                <h4 className="text-sm font-bold text-slate-900 dark:text-white">{n.title}</h4>
                <div className="flex items-center space-x-2">
                  <span className="text-[11px] text-slate-400 font-mono">{n.time}</span>
                  <button
                    onClick={() => handleDeleteNotif(n.id)}
                    className="p-1 rounded-lg text-slate-400 hover:text-rose-600 hover:bg-rose-50 dark:hover:bg-rose-950/40 transition-colors"
                  >
                    <Trash2 className="h-3.5 w-3.5" />
                  </button>
                </div>
              </div>
              <p className="text-xs text-slate-600 dark:text-slate-400 mt-1">{n.message}</p>
              <div className="mt-2 flex items-center space-x-2 text-[10px] font-bold text-slate-400">
                <span className="uppercase px-2 py-0.5 rounded-md bg-slate-100 dark:bg-slate-800">{n.channel}</span>
                {!n.isRead && <span className="text-brand-600 font-black">• BARU</span>}
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Modal Kirim Siaran Pengumuman */}
      {isModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 backdrop-blur-sm p-4">
          <div className="relative w-full max-w-md rounded-3xl bg-white dark:bg-slate-900 p-6 shadow-2xl border border-slate-200 dark:border-slate-800 space-y-6">
            <div className="flex items-center justify-between border-b border-slate-200 dark:border-slate-800 pb-4">
              <div className="flex items-center space-x-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-brand-100 text-brand-600">
                  <Send className="h-5 w-5" />
                </div>
                <div>
                  <h3 className="text-lg font-bold text-slate-900 dark:text-white">Kirim Siaran Pengumuman</h3>
                  <p className="text-xs text-slate-500">Notifikasi push langsung ke Customer & Mitra Warung B2B.</p>
                </div>
              </div>
              <button onClick={() => setIsModalOpen(false)} className="rounded-xl p-2 text-slate-400 hover:bg-slate-100">
                <X className="h-5 w-5" />
              </button>
            </div>

            <form onSubmit={handleSendBroadcast} className="space-y-4">
              <div>
                <label className="block text-xs font-bold text-slate-500 mb-1">Target Penerima *</label>
                <select
                  value={targetAudience}
                  onChange={(e) => setTargetAudience(e.target.value)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white font-bold"
                >
                  <option value="ALL">Semua Pengguna (B2C & B2B Mitra)</option>
                  <option value="B2B_PARTNERS">Khusus Mitra Toko Kelontong B2B</option>
                  <option value="B2C_CUSTOMERS">Khusus Pembeli Eceran B2C</option>
                </select>
              </div>

              <div>
                <label className="block text-xs font-bold text-slate-500 mb-1">Judul Pengumuman *</label>
                <input
                  type="text"
                  required
                  placeholder="Contoh: Promo Bebas Ongkir Pasokan Beras & Minyak"
                  value={broadcastTitle}
                  onChange={(e) => setBroadcastTitle(e.target.value)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white font-bold"
                />
              </div>

              <div>
                <label className="block text-xs font-bold text-slate-500 mb-1">Isi Pesan Pengumuman *</label>
                <textarea
                  rows={3}
                  required
                  placeholder="Tuliskan isi pengumuman atau promo terkini..."
                  value={broadcastMessage}
                  onChange={(e) => setBroadcastMessage(e.target.value)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white resize-none"
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
                  Kirim Sekarang
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
