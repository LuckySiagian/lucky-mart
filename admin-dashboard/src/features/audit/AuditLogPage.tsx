import React from 'react';
import { PageHeader } from '../../components/common/PageHeader';
import { LoadingSkeleton } from '../../components/common/LoadingSkeleton';
import { EmptyState } from '../../components/common/EmptyState';
import { ErrorState } from '../../components/common/ErrorState';
import { Card, Pagination, Table } from '../../components/common/Pagination';
import { eventService } from '../../services/apiServices';
import { useList } from '../../hooks/useList';
import { formatDate } from '../../lib/utils';

const EVENT_LABELS: Record<string, string> = {
  ORDER_CREATED: 'Pesanan dibuat',
  ORDER_CANCELLED: 'Pesanan dibatalkan',
  ORDER_STATUS_CHANGED: 'Status pesanan diubah',
  ORDER_COMPLETED: 'Pesanan selesai',
  PAYMENT_PROOF_UPLOADED: 'Bukti transfer diunggah',
  PAYMENT_PAID: 'Pembayaran disetujui',
  PAYMENT_REJECTED: 'Bukti transfer ditolak',
  SHIPMENT_OUT_FOR_DELIVERY: 'Pesanan diantar',
  CREDIT_SETTLED: 'Pelunasan kasbon',
  B2B_REGISTERED: 'Pendaftaran mitra toko',
  B2B_VERIFIED: 'Verifikasi mitra toko',
  PASSWORD_RESET_REQUESTED: 'Permintaan reset kata sandi',
};

function describe(payload: string): string {
  try {
    const p = JSON.parse(payload) as Record<string, unknown>;
    const parts = [p.order_number, p.message || p.staff_message].filter(Boolean);
    return parts.join(' — ') || '-';
  } catch {
    return '-';
  }
}

export const AuditLogPage: React.FC = () => {
  const list = useList((p) => eventService.getEvents(p), { limit: 30 });

  return (
    <div className="space-y-5">
      <PageHeader title="Log Aktivitas" subtitle="Catatan kejadian penting di sistem: pesanan, pembayaran, pengiriman, dan akun." />

      {list.loading ? (
        <LoadingSkeleton count={6} />
      ) : list.error ? (
        <ErrorState message={list.error} onRetry={list.reload} />
      ) : list.items.length === 0 ? (
        <EmptyState title="Belum ada aktivitas" />
      ) : (
        <Card className="p-2 sm:p-4">
          <Table
            head={
              <tr>
                <th className="px-3 py-3">Waktu</th>
                <th className="px-3 py-3">Kejadian</th>
                <th className="px-3 py-3">Keterangan</th>
              </tr>
            }
          >
            {list.items.map((e) => (
              <tr key={e.id}>
                <td className="whitespace-nowrap px-3 py-3 text-xs text-slate-500">{formatDate(e.created_at)}</td>
                <td className="px-3 py-3">
                  <div className="font-semibold">{EVENT_LABELS[e.event_name] || e.event_name}</div>
                  <div className="text-[11px] text-slate-400">{e.aggregate}</div>
                </td>
                <td className="px-3 py-3 text-xs text-slate-600 dark:text-slate-400">{describe(e.payload)}</td>
              </tr>
            ))}
          </Table>
          <Pagination meta={list.meta} onPage={list.setPage} />
        </Card>
      )}
    </div>
  );
};
