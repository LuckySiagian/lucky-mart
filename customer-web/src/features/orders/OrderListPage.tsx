import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { ChevronRight, ShoppingBag } from 'lucide-react';
import { orderService } from '../../services/apiServices';
import { OrderStatusBadge } from '../../components/common/Badge';
import { EmptyState } from '../../components/common/EmptyState';
import { ErrorState } from '../../components/common/ErrorState';
import { formatCurrency, formatDate, getErrorMessage } from '../../lib/utils';
import type { Order } from '../../types';

const TABS = [
  { value: '', label: 'Semua' },
  { value: 'PENDING_PAYMENT', label: 'Belum Bayar' },
  { value: 'CONFIRMED', label: 'Dikonfirmasi' },
  { value: 'PROCESSING_PACKING', label: 'Dikemas' },
  { value: 'OUT_FOR_DELIVERY', label: 'Diantar' },
  { value: 'COMPLETED', label: 'Selesai' },
  { value: 'CANCELLED', label: 'Dibatalkan' },
];

export const OrderListPage: React.FC = () => {
  const [status, setStatus] = useState('');
  const [orders, setOrders] = useState<Order[]>([]);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => setPage(1), [status]);

  useEffect(() => {
    let alive = true;
    setLoading(true);
    setError(null);
    orderService
      .getOrders({ page, limit: 10, status })
      .then((res) => {
        if (!alive) return;
        setOrders((prev) => (page === 1 ? res.items : [...prev, ...res.items]));
        setTotalPages(res.meta.total_pages);
      })
      .catch((err) => alive && setError(getErrorMessage(err)))
      .finally(() => alive && setLoading(false));
    return () => {
      alive = false;
    };
  }, [page, status]);

  return (
    <div className="mx-auto max-w-4xl space-y-5 px-4 py-8 sm:px-6 lg:px-8">
      <h1 className="flex items-center gap-3 text-2xl font-black text-slate-900 dark:text-white">
        <ShoppingBag className="h-7 w-7 text-primary-600" /> Pesanan Saya
      </h1>

      <div className="flex gap-2 overflow-x-auto pb-1">
        {TABS.map((t) => (
          <button
            key={t.value}
            onClick={() => setStatus(t.value)}
            className={`flex-shrink-0 rounded-xl px-4 py-2 text-xs font-bold ${
              status === t.value ? 'bg-primary-600 text-white' : 'bg-white text-slate-600 border border-slate-200 dark:bg-slate-900 dark:border-slate-800 dark:text-slate-300'
            }`}
          >
            {t.label}
          </button>
        ))}
      </div>

      {error ? (
        <ErrorState message={error} onRetry={() => setPage(1)} />
      ) : loading && orders.length === 0 ? (
        <div className="space-y-3">
          {[0, 1, 2].map((i) => (
            <div key={i} className="h-28 animate-pulse rounded-3xl bg-slate-200 dark:bg-slate-800" />
          ))}
        </div>
      ) : orders.length === 0 ? (
        <EmptyState title="Belum ada pesanan" description="Pesanan yang Anda buat akan muncul di sini." actionText="Mulai Belanja" />
      ) : (
        <div className="space-y-3">
          {orders.map((o) => (
            <Link
              key={o.id}
              to={`/orders/${o.id}`}
              className="block rounded-3xl border border-slate-200 bg-white p-5 shadow-sm transition-colors hover:border-primary-300 dark:border-slate-800 dark:bg-slate-900"
            >
              <div className="flex flex-wrap items-center justify-between gap-2">
                <div>
                  <span className="font-mono text-sm font-bold">{o.order_number}</span>
                  <span className="ml-2 text-xs text-slate-400">{formatDate(o.order_date)}</span>
                </div>
                <OrderStatusBadge status={o.status} deliveryType={o.delivery_type} />
              </div>
              <div className="mt-3 flex items-center justify-between gap-3">
                <div className="min-w-0 text-sm text-slate-600 dark:text-slate-400">
                  <div className="truncate">{o.items.map((i) => i.product_name).join(', ')}</div>
                  <div className="text-xs text-slate-400">{o.total_items} barang</div>
                </div>
                <div className="flex items-center gap-2">
                  <span className="text-base font-black text-primary-600">{formatCurrency(o.grand_total)}</span>
                  <ChevronRight className="h-4 w-4 text-slate-400" />
                </div>
              </div>
              {o.status === 'PENDING_PAYMENT' && (
                <div className="mt-3 rounded-xl bg-amber-50 px-3 py-2 text-xs font-semibold text-amber-800 dark:bg-amber-950/40 dark:text-amber-300">
                  Selesaikan pembayaran dan unggah bukti transfer →
                </div>
              )}
            </Link>
          ))}
          {page < totalPages && (
            <div className="text-center">
              <button
                disabled={loading}
                onClick={() => setPage((p) => p + 1)}
                className="rounded-2xl border border-slate-200 bg-white px-6 py-2.5 text-sm font-bold disabled:opacity-50 dark:border-slate-700 dark:bg-slate-900"
              >
                {loading ? 'Memuat...' : 'Muat lebih banyak'}
              </button>
            </div>
          )}
        </div>
      )}
    </div>
  );
};
