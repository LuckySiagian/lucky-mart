import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { AlertTriangle, BadgeCheck, Boxes, ClipboardCheck, CreditCard, PackageOpen, Store, Truck, Users, Wallet } from 'lucide-react';
import { Bar, BarChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts';
import { PageHeader } from '../../components/common/PageHeader';
import { StatusBadge } from '../../components/common/StatusBadge';
import { LoadingSkeleton } from '../../components/common/LoadingSkeleton';
import { ErrorState } from '../../components/common/ErrorState';
import { Card } from '../../components/common/Pagination';
import { dashboardService } from '../../services/apiServices';
import { useAuthStore, MANAGER_ROLES, hasRole } from '../../store/useAuthStore';
import { formatCurrency, formatDate, formatNumber, getErrorMessage } from '../../lib/utils';
import type { DashboardStats } from '../../types';

const dayLabel = (iso: string) => new Intl.DateTimeFormat('id-ID', { weekday: 'short', day: 'numeric' }).format(new Date(iso + 'T00:00:00'));

export const OverviewPage: React.FC = () => {
  const isManager = hasRole(useAuthStore((s) => s.user), MANAGER_ROLES);
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [error, setError] = useState<string | null>(null);

  const load = async () => {
    setError(null);
    try {
      setStats(await dashboardService.getStats());
    } catch (err) {
      setError(getErrorMessage(err));
    }
  };

  useEffect(() => {
    load();
    const t = setInterval(load, 60000);
    return () => clearInterval(t);
  }, []);

  if (error) return <ErrorState message={error} onRetry={load} />;
  if (!stats) return <LoadingSkeleton count={6} />;

  const p = stats.pending;
  const todo = [
    { show: isManager, count: p.payments_to_verify, label: 'Bukti transfer perlu dicek', to: '/orders?tab=bayar', icon: CreditCard },
    { show: true, count: p.orders_to_pack, label: 'Pesanan baru perlu diproses', to: '/orders?tab=proses', icon: ClipboardCheck },
    { show: true, count: p.orders_packing, label: 'Pesanan sedang disiapkan', to: '/orders?tab=proses', icon: PackageOpen },
    { show: true, count: p.orders_delivering, label: 'Pesanan diantar / menunggu diambil', to: '/orders?tab=jalan', icon: Truck },
    { show: true, count: p.low_stock_items, label: 'Barang stoknya menipis', to: '/inventory', icon: AlertTriangle },
    { show: isManager, count: p.b2b_to_verify, label: 'Pendaftaran mitra toko menunggu', to: '/users', icon: BadgeCheck },
    { show: true, count: p.awaiting_payment, label: 'Pesanan menunggu transfer pelanggan', to: '/orders?tab=bayar', icon: Wallet },
  ].filter((t) => t.show);

  const kpis = [
    { label: 'Omzet hari ini', value: formatCurrency(stats.today_revenue), sub: `${stats.today_orders} pesanan`, icon: Wallet, tone: 'emerald' },
    { label: 'Omzet bulan ini', value: formatCurrency(stats.month_revenue), sub: `${stats.month_orders} pesanan`, icon: Store, tone: 'sky' },
    { label: 'Produk aktif', value: formatNumber(stats.active_products), sub: `${p.low_stock_items} stok menipis`, icon: Boxes, tone: 'amber' },
    isManager
      ? { label: 'Kasbon belum lunas', value: formatCurrency(stats.kasbon_outstanding), sub: `${stats.total_customers} pelanggan terdaftar`, icon: CreditCard, tone: 'rose' }
      : { label: 'Pelanggan', value: formatNumber(stats.total_customers), sub: 'terdaftar', icon: Users, tone: 'rose' },
  ];
  const toneClass: Record<string, string> = {
    emerald: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-950/60 dark:text-emerald-400',
    sky: 'bg-sky-100 text-sky-700 dark:bg-sky-950/60 dark:text-sky-400',
    amber: 'bg-amber-100 text-amber-700 dark:bg-amber-950/60 dark:text-amber-400',
    rose: 'bg-rose-100 text-rose-700 dark:bg-rose-950/60 dark:text-rose-400',
  };

  return (
    <div className="space-y-6">
      <PageHeader title="Ringkasan Toko" subtitle="Diperbarui otomatis setiap menit." />

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {kpis.map((k) => (
          <Card key={k.label} className="p-5">
            <div className="flex items-center justify-between">
              <span className="text-xs font-semibold uppercase text-slate-500">{k.label}</span>
              <div className={`flex h-9 w-9 items-center justify-center rounded-xl ${toneClass[k.tone]}`}>
                <k.icon className="h-5 w-5" />
              </div>
            </div>
            <div className="mt-3 text-2xl font-black text-slate-900 dark:text-white">{k.value}</div>
            <div className="mt-1 text-xs text-slate-500">{k.sub}</div>
          </Card>
        ))}
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        <Card className="p-5 lg:col-span-1">
          <h3 className="mb-3 text-sm font-bold">Perlu ditindaklanjuti</h3>
          <ul className="space-y-1">
            {todo.map((t) => (
              <li key={t.label}>
                <Link
                  to={t.to}
                  className={`flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm hover:bg-slate-50 dark:hover:bg-slate-800 ${t.count === 0 ? 'opacity-50' : ''}`}
                >
                  <t.icon className="h-4 w-4 text-slate-400" />
                  <span className="flex-1">{t.label}</span>
                  <span className={`min-w-[2rem] rounded-full px-2 py-0.5 text-center text-xs font-bold ${t.count > 0 ? 'bg-brand-600 text-white' : 'bg-slate-100 text-slate-500 dark:bg-slate-800'}`}>
                    {t.count}
                  </span>
                </Link>
              </li>
            ))}
          </ul>
        </Card>

        <Card className="p-5 lg:col-span-2">
          <h3 className="mb-4 text-sm font-bold">Penjualan 7 hari terakhir</h3>
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={stats.sales_last_7_days.map((d) => ({ ...d, label: dayLabel(d.date) }))}>
                <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#e2e8f0" />
                <XAxis dataKey="label" tick={{ fontSize: 12 }} />
                <YAxis tick={{ fontSize: 11 }} tickFormatter={(v) => (v >= 1000000 ? `${v / 1000000}jt` : v >= 1000 ? `${v / 1000}rb` : v)} width={48} />
                <Tooltip formatter={(v: number) => formatCurrency(v)} labelFormatter={(l) => l} />
                <Bar dataKey="revenue" name="Omzet" fill="#16a34a" radius={[6, 6, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </Card>
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        <Card className="p-5">
          <h3 className="mb-3 text-sm font-bold">Produk terlaris (30 hari)</h3>
          {stats.top_products.length === 0 ? (
            <p className="py-6 text-center text-sm text-slate-500">Belum ada penjualan.</p>
          ) : (
            <ul className="divide-y divide-slate-100 dark:divide-slate-800">
              {stats.top_products.map((tp, i) => (
                <li key={tp.product_name} className="flex items-center gap-3 py-2.5 text-sm">
                  <span className="flex h-6 w-6 items-center justify-center rounded-full bg-slate-100 text-xs font-bold dark:bg-slate-800">{i + 1}</span>
                  <span className="flex-1 font-medium">{tp.product_name}</span>
                  <span className="text-xs text-slate-500">{formatNumber(tp.quantity)} terjual</span>
                  <span className="w-28 text-right font-semibold">{formatCurrency(tp.revenue)}</span>
                </li>
              ))}
            </ul>
          )}
        </Card>

        <Card className="p-5">
          <div className="mb-3 flex items-center justify-between">
            <h3 className="text-sm font-bold">Pesanan terbaru</h3>
            <Link to="/orders" className="text-xs font-semibold text-brand-600">Lihat semua</Link>
          </div>
          {stats.recent_orders.length === 0 ? (
            <p className="py-6 text-center text-sm text-slate-500">Belum ada pesanan.</p>
          ) : (
            <ul className="divide-y divide-slate-100 dark:divide-slate-800">
              {stats.recent_orders.map((o) => (
                <li key={o.id} className="flex items-center justify-between gap-3 py-2.5 text-sm">
                  <div>
                    <div className="font-mono text-xs font-semibold">{o.order_number}</div>
                    <div className="text-xs text-slate-500">{o.customer_name} · {formatDate(o.order_date)}</div>
                  </div>
                  <div className="text-right">
                    <div className="font-semibold">{formatCurrency(o.grand_total)}</div>
                    <StatusBadge status={o.status} />
                  </div>
                </li>
              ))}
            </ul>
          )}
        </Card>
      </div>
    </div>
  );
};
