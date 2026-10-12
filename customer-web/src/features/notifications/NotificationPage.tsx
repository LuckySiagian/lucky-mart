import React, { useEffect, useState } from 'react';
import { Bell, CheckCheck } from 'lucide-react';
import { notificationService } from '../../services/apiServices';
import { EmptyState } from '../../components/common/EmptyState';
import { ErrorState } from '../../components/common/ErrorState';
import { formatDate, getErrorMessage } from '../../lib/utils';
import type { Notification } from '../../types';

export const NotificationPage: React.FC = () => {
  const [items, setItems] = useState<Notification[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const load = () => {
    setLoading(true);
    notificationService
      .getNotifications({ limit: 50 })
      .then((r) => setItems(r.items))
      .catch((e) => setError(getErrorMessage(e)))
      .finally(() => setLoading(false));
  };

  useEffect(load, []);

  const changed = () => window.dispatchEvent(new Event('notifications:changed'));

  const markAll = async () => {
    await notificationService.markAllRead().catch(() => undefined);
    setItems((ns) => ns.map((n) => ({ ...n, is_read: true })));
    changed();
  };

  const markOne = async (n: Notification) => {
    if (n.is_read) return;
    await notificationService.markRead(n.id).catch(() => undefined);
    setItems((ns) => ns.map((x) => (x.id === n.id ? { ...x, is_read: true } : x)));
    changed();
  };

  return (
    <div className="mx-auto max-w-3xl space-y-5 px-4 py-8 sm:px-6 lg:px-8">
      <div className="flex items-center justify-between">
        <h1 className="flex items-center gap-2 text-2xl font-black text-slate-900 dark:text-white">
          <Bell className="h-6 w-6 text-primary-600" /> Notifikasi
        </h1>
        {items.some((n) => !n.is_read) && (
          <button onClick={markAll} className="flex items-center gap-1 text-xs font-bold text-primary-600">
            <CheckCheck className="h-4 w-4" /> Tandai semua dibaca
          </button>
        )}
      </div>

      {error ? (
        <ErrorState message={error} onRetry={load} />
      ) : loading ? (
        <div className="h-40 animate-pulse rounded-3xl bg-slate-200 dark:bg-slate-800" />
      ) : items.length === 0 ? (
        <EmptyState title="Belum ada notifikasi" description="Kabar pesanan dan pembayaran Anda akan muncul di sini." />
      ) : (
        <div className="divide-y divide-slate-100 overflow-hidden rounded-3xl border border-slate-200 bg-white dark:divide-slate-800 dark:border-slate-800 dark:bg-slate-900">
          {items.map((n) => (
            <button key={n.id} onClick={() => markOne(n)} className={`block w-full p-4 text-left ${n.is_read ? '' : 'bg-primary-50/60 dark:bg-primary-950/20'}`}>
              <div className="flex items-center gap-2">
                <span className={`text-sm ${n.is_read ? 'font-medium' : 'font-bold'}`}>{n.title}</span>
                {!n.is_read && <span className="h-2 w-2 rounded-full bg-primary-500" />}
              </div>
              <p className="mt-0.5 text-sm text-slate-600 dark:text-slate-400">{n.message}</p>
              <span className="text-xs text-slate-400">{formatDate(n.created_at)}</span>
            </button>
          ))}
        </div>
      )}
    </div>
  );
};
