import React from 'react';
import { toast } from 'sonner';
import { Bell, CheckCheck, Trash2 } from 'lucide-react';
import { PageHeader } from '../../components/common/PageHeader';
import { LoadingSkeleton } from '../../components/common/LoadingSkeleton';
import { EmptyState } from '../../components/common/EmptyState';
import { ErrorState } from '../../components/common/ErrorState';
import { btnSecondary } from '../../components/common/Modal';
import { Card, Pagination } from '../../components/common/Pagination';
import { notificationService } from '../../services/apiServices';
import { useList } from '../../hooks/useList';
import { formatDate, getErrorMessage } from '../../lib/utils';

export const NotificationPage: React.FC = () => {
  const list = useList((p) => notificationService.getNotifications(p));

  const markAll = async () => {
    try {
      await notificationService.markAllRead();
      list.setItems((items) => items.map((n) => ({ ...n, is_read: true })));
      window.dispatchEvent(new Event('notifications:changed'));
    } catch (err) {
      toast.error(getErrorMessage(err));
    }
  };

  const markOne = async (id: string) => {
    try {
      await notificationService.markRead(id);
      list.setItems((items) => items.map((n) => (n.id === id ? { ...n, is_read: true } : n)));
      window.dispatchEvent(new Event('notifications:changed'));
    } catch {
      /* not important */
    }
  };

  const remove = async (id: string) => {
    try {
      await notificationService.remove(id);
      list.setItems((items) => items.filter((n) => n.id !== id));
      window.dispatchEvent(new Event('notifications:changed'));
    } catch (err) {
      toast.error(getErrorMessage(err));
    }
  };

  return (
    <div className="space-y-5">
      <PageHeader title="Notifikasi" subtitle="Pesanan baru, bukti transfer masuk, pendaftaran mitra, dan permintaan reset kata sandi.">
        <button className={btnSecondary} onClick={markAll}>
          <CheckCheck className="h-4 w-4" /> Tandai semua dibaca
        </button>
      </PageHeader>

      {list.loading ? (
        <LoadingSkeleton count={5} />
      ) : list.error ? (
        <ErrorState message={list.error} onRetry={list.reload} />
      ) : list.items.length === 0 ? (
        <EmptyState title="Belum ada notifikasi" description="Pemberitahuan toko akan muncul di sini." />
      ) : (
        <Card className="divide-y divide-slate-100 dark:divide-slate-800">
          {list.items.map((n) => (
            <div
              key={n.id}
              onClick={() => !n.is_read && markOne(n.id)}
              className={`flex cursor-pointer items-start gap-3 p-4 ${n.is_read ? '' : 'bg-brand-50/60 dark:bg-brand-950/20'}`}
            >
              <div className={`mt-0.5 flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-xl ${n.is_read ? 'bg-slate-100 text-slate-400 dark:bg-slate-800' : 'bg-brand-100 text-brand-700 dark:bg-brand-900/50 dark:text-brand-300'}`}>
                <Bell className="h-4 w-4" />
              </div>
              <div className="min-w-0 flex-1">
                <div className="flex items-center gap-2">
                  <span className={`text-sm ${n.is_read ? 'font-medium' : 'font-bold'} text-slate-900 dark:text-white`}>{n.title}</span>
                  {!n.is_read && <span className="h-2 w-2 rounded-full bg-brand-500" />}
                </div>
                <p className="mt-0.5 text-sm text-slate-600 dark:text-slate-400">{n.message}</p>
                <span className="text-xs text-slate-400">{formatDate(n.created_at)}</span>
              </div>
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  remove(n.id);
                }}
                className="rounded-lg p-2 text-slate-400 hover:bg-rose-50 hover:text-rose-600 dark:hover:bg-rose-950/30"
                title="Hapus"
              >
                <Trash2 className="h-4 w-4" />
              </button>
            </div>
          ))}
        </Card>
      )}
      <Pagination meta={list.meta} onPage={list.setPage} />
    </div>
  );
};
