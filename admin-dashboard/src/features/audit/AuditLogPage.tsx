import React, { useState, useEffect } from 'react';
import { PageHeader } from '../../components/common/PageHeader';
import { ShieldCheck, Lock, Activity, Search, CheckCircle2, AlertCircle, Info } from 'lucide-react';
import { notificationService } from '../../services/apiServices';

interface AuditEvent {
  id: string;
  time: string;
  user: string;
  role: string;
  eventName: string;
  module: string;
  ipAddress: string;
  status: 'SUCCESS' | 'WARNING' | 'FAILED';
  details: string;
}

const mockEvents: AuditEvent[] = [
  {
    id: 'AUD-9001',
    time: '2026-07-27 16:22:10',
    user: 'Super Admin Sembako',
    role: 'SUPER_ADMIN',
    eventName: 'CREATE_NEW_PRODUCT',
    module: 'KATALOG_PRODUK',
    ipAddress: '127.0.0.1 (Localhost)',
    status: 'SUCCESS',
    details: 'Menambahkan produk baru "Beras Rojo Lele Super 5kg" ke database.',
  },
  {
    id: 'AUD-9002',
    time: '2026-07-27 15:45:00',
    user: 'Finance & Operations Admin',
    role: 'STORE_MANAGER',
    eventName: 'VERIFY_PAYMENT_PROOF',
    module: 'KEUANGAN_PAYMENT',
    ipAddress: '192.168.1.105',
    status: 'SUCCESS',
    details: 'Menyetujui bukti transfer manual Rp 145.000 untuk Order #8821.',
  },
  {
    id: 'AUD-9003',
    time: '2026-07-27 14:10:30',
    user: 'Super Admin Sembako',
    role: 'SUPER_ADMIN',
    eventName: 'ADJUST_STOCK_INVENTORY',
    module: 'INVENTARIS_GUDANG',
    ipAddress: '127.0.0.1 (Localhost)',
    status: 'SUCCESS',
    details: 'Penyesuaian stok Telur Ayam (-2kg) karena handling damage.',
  },
  {
    id: 'AUD-9004',
    time: '2026-07-27 11:30:15',
    user: 'Super Admin Sembako',
    role: 'SUPER_ADMIN',
    eventName: 'UPDATE_B2B_CREDIT_LIMIT',
    module: 'B2B_CREDIT_LEDGER',
    ipAddress: '127.0.0.1 (Localhost)',
    status: 'SUCCESS',
    details: 'Mengubah limit kredit Toko Kelontong Berkah Ibu menjadi Rp 50.000.000.',
  },
  {
    id: 'AUD-9005',
    time: '2026-07-27 09:15:00',
    user: 'System Auto-Monitor',
    role: 'SYSTEM',
    eventName: 'LOW_STOCK_THRESHOLD_ALERT',
    module: 'STOK_MONITOR',
    ipAddress: '127.0.0.1 (Internal Worker)',
    status: 'WARNING',
    details: 'Deteksi stok Minyak Bimoli 2L mencapai ambang kritis (12 Unit).',
  },
];

export const AuditLogPage: React.FC = () => {
  const [events, setEvents] = useState<AuditEvent[]>(mockEvents);
  const [search, setSearch] = useState('');

  useEffect(() => {
    const fetchLogs = async () => {
      try {
        const res = await notificationService.getEventLogs(1, 50);
        const list = Array.isArray(res.data) ? res.data : (res.data?.items || []);
        if (list.length > 0) {
          const mapped = list.map((ev: any) => ({
            id: ev.id,
            time: ev.created_at ? new Date(ev.created_at).toLocaleString('id-ID') : 'Baru saja',
            user: ev.user?.full_name || 'System Operator',
            role: ev.user?.role || 'SYSTEM',
            eventName: ev.event_name || 'SYSTEM_ACTION',
            module: ev.module_name || 'CORE_SYSTEM',
            ipAddress: ev.ip_address || '127.0.0.1',
            status: ev.status || 'SUCCESS',
            details: ev.description || ev.details || 'Aktivitas tercatat di log audit system.',
          }));
          setEvents(mapped);
        }
      } catch (e) {
        // Fallback to mock
      }
    };
    fetchLogs();
  }, []);

  const filtered = events.filter(
    (e) =>
      e.user.toLowerCase().includes(search.toLowerCase()) ||
      e.eventName.toLowerCase().includes(search.toLowerCase()) ||
      e.details.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className="space-y-6">
      <PageHeader
        title="Audit Events Log (Jejak Keamanan Sistem)"
        subtitle="Catatan log immutable (Audit Trail) seluruh aktivitas sensitif perubah data, transaksi keuangan, dan perubahan stok oleh admin."
      >
        <div className="flex items-center space-x-2 text-xs font-semibold text-slate-500 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 px-3 py-2 rounded-xl">
          <ShieldCheck className="h-4 w-4 text-emerald-600" />
          <span>Security Audit Trail Active</span>
        </div>
      </PageHeader>

      {/* Filter & Search Bar */}
      <div className="flex flex-col sm:flex-row items-center justify-between gap-4 bg-white dark:bg-slate-900 p-4 rounded-2xl border border-slate-200 dark:border-slate-800">
        <div className="relative w-full sm:w-80">
          <Search className="absolute left-3.5 top-3 h-4 w-4 text-slate-400" />
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Cari event, user, atau aktivitas..."
            className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 pl-10 pr-4 py-2 text-sm text-slate-900 dark:text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-brand-500"
          />
        </div>
        <div className="text-xs text-slate-500 font-semibold">
          Total Log Audit: <span className="text-brand-600 font-bold">{events.length} Events</span>
        </div>
      </div>

      {/* Events Log Table */}
      <div className="overflow-x-auto rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 shadow-sm">
        <table className="w-full text-left text-sm text-slate-600 dark:text-slate-400">
          <thead className="bg-slate-50 dark:bg-slate-800/50 text-xs uppercase text-slate-500 font-bold border-b border-slate-200 dark:border-slate-800">
            <tr>
              <th className="px-6 py-4">Waktu Event</th>
              <th className="px-6 py-4">Pengguna / Operator</th>
              <th className="px-6 py-4">Nama Event Aktivitas</th>
              <th className="px-6 py-4">Modul System</th>
              <th className="px-6 py-4">Status & IP Address</th>
              <th className="px-6 py-4">Rincian Perubahan Data</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-200 dark:divide-slate-800">
            {filtered.map((e) => (
              <tr key={e.id} className="hover:bg-slate-50/80 dark:hover:bg-slate-800/30">
                <td className="px-6 py-4 font-mono text-xs text-slate-500">{e.time}</td>
                <td className="px-6 py-4">
                  <div className="font-bold text-slate-900 dark:text-white">{e.user}</div>
                  <div className="text-[11px] text-brand-600 font-mono font-semibold">{e.role}</div>
                </td>
                <td className="px-6 py-4 font-mono text-xs font-bold text-slate-800 dark:text-slate-200">
                  {e.eventName}
                </td>
                <td className="px-6 py-4 text-xs font-semibold text-slate-500">{e.module}</td>
                <td className="px-6 py-4">
                  <div className="space-y-1">
                    {e.status === 'SUCCESS' && (
                      <span className="inline-flex items-center space-x-1 rounded-full bg-emerald-100 dark:bg-emerald-950/50 px-2.5 py-0.5 text-xs font-bold text-emerald-700 dark:text-emerald-400">
                        <CheckCircle2 className="h-3.5 w-3.5" />
                        <span>SUCCESS</span>
                      </span>
                    )}
                    {e.status === 'WARNING' && (
                      <span className="inline-flex items-center space-x-1 rounded-full bg-amber-100 dark:bg-amber-950/50 px-2.5 py-0.5 text-xs font-bold text-amber-700 dark:text-amber-400">
                        <AlertCircle className="h-3.5 w-3.5" />
                        <span>WARNING</span>
                      </span>
                    )}
                    <div className="text-[11px] text-slate-400 font-mono">{e.ipAddress}</div>
                  </div>
                </td>
                <td className="px-6 py-4 text-xs text-slate-600 dark:text-slate-400 max-w-sm">
                  {e.details}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};
