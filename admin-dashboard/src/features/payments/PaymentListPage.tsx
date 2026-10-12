import React, { useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { toast } from 'sonner';
import { CheckCircle2, Eye, ImageOff, RotateCcw, XCircle } from 'lucide-react';
import { PageHeader } from '../../components/common/PageHeader';
import { StatusBadge } from '../../components/common/StatusBadge';
import { LoadingSkeleton } from '../../components/common/LoadingSkeleton';
import { EmptyState } from '../../components/common/EmptyState';
import { ErrorState } from '../../components/common/ErrorState';
import { Field, Modal, btnDanger, btnPrimary, btnSecondary, inputClass } from '../../components/common/Modal';
import { Card, FilterTabs, Pagination, SearchInput, Table } from '../../components/common/Pagination';
import { paymentService } from '../../services/apiServices';
import { useList } from '../../hooks/useList';
import { formatCurrency, formatDate, getErrorMessage, statusLabel } from '../../lib/utils';
import type { Payment } from '../../types';

const TABS = [
  { value: 'WAITING_CONFIRMATION', label: 'Perlu Verifikasi' },
  { value: 'PENDING', label: 'Belum Dibayar' },
  { value: 'PAID', label: 'Lunas' },
  { value: 'FAILED', label: 'Ditolak' },
  { value: '', label: 'Semua' },
];

type Action = 'approve' | 'reject' | 'refund' | null;

export const PaymentListPage: React.FC = () => {
  const [searchParams] = useSearchParams();
  const list = useList((p) => paymentService.getPayments(p), { status: searchParams.get('status') ?? 'WAITING_CONFIRMATION' });
  const [selected, setSelected] = useState<Payment | null>(null);
  const [action, setAction] = useState<Action>(null);
  const [notes, setNotes] = useState('');
  const [refundAmount, setRefundAmount] = useState('');
  const [busy, setBusy] = useState(false);

  const open = async (p: Payment) => {
    try {
      setSelected(await paymentService.getPaymentById(p.id));
    } catch (err) {
      toast.error(getErrorMessage(err));
    }
  };

  const startAction = (a: Action) => {
    setAction(a);
    setNotes('');
    setRefundAmount(a === 'refund' && selected ? String(selected.grand_total) : '');
  };

  const submit = async () => {
    if (!selected || !action) return;
    if ((action === 'reject' || action === 'refund') && notes.trim().length < 3) {
      toast.error('Tulis alasannya (minimal 3 huruf).');
      return;
    }
    setBusy(true);
    try {
      let updated: Payment;
      if (action === 'approve') updated = await paymentService.approvePayment(selected.id, notes.trim());
      else if (action === 'reject') updated = await paymentService.rejectPayment(selected.id, notes.trim());
      else updated = await paymentService.refundPayment(selected.id, Number(refundAmount), notes.trim());
      setSelected(updated);
      setAction(null);
      toast.success(`Pembayaran ${updated.payment_number}: ${statusLabel(updated.status)}`);
      list.reload();
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setBusy(false);
    }
  };

  const canVerify = selected && (selected.status === 'WAITING_CONFIRMATION' || (selected.status === 'PENDING' && selected.method_type === 'MANUAL_TRANSFER'));
  const canRefund = selected && selected.status === 'PAID' && selected.method_type !== 'B2B_CREDIT';

  return (
    <div className="space-y-5">
      <PageHeader title="Pembayaran" subtitle="Verifikasi bukti transfer pelanggan dan kelola pengembalian dana." />

      <div className="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
        <FilterTabs value={(list.params.status as string) ?? ''} onChange={(v) => list.setFilter('status', v)} options={TABS} />
        <SearchInput value={list.searchInput} onChange={list.setSearchInput} placeholder="Cari nomor pembayaran..." />
      </div>

      {list.loading ? (
        <LoadingSkeleton count={4} />
      ) : list.error ? (
        <ErrorState message={list.error} onRetry={list.reload} />
      ) : list.items.length === 0 ? (
        <EmptyState
          title={list.params.status === 'WAITING_CONFIRMATION' ? 'Tidak ada yang perlu diverifikasi' : 'Belum ada data pembayaran'}
          description="Bukti transfer yang diunggah pelanggan akan muncul di tab Perlu Verifikasi."
        />
      ) : (
        <Card className="p-2 sm:p-4">
          <Table
            head={
              <tr>
                <th className="px-3 py-3">Pembayaran</th>
                <th className="px-3 py-3">Pelanggan</th>
                <th className="px-3 py-3">Metode</th>
                <th className="px-3 py-3 text-right">Jumlah</th>
                <th className="px-3 py-3">Status</th>
                <th className="px-3 py-3" />
              </tr>
            }
          >
            {list.items.map((p) => (
              <tr key={p.id} className="hover:bg-slate-50 dark:hover:bg-slate-800/40">
                <td className="px-3 py-3">
                  <div className="font-mono text-xs font-semibold">{p.payment_number}</div>
                  <div className="text-xs text-slate-500">Pesanan {p.order_number} · {formatDate(p.created_at)}</div>
                </td>
                <td className="px-3 py-3">
                  <div className="font-medium">{p.customer_name || '-'}</div>
                  <div className="text-xs text-slate-500">{p.customer_phone}</div>
                </td>
                <td className="px-3 py-3 text-xs">{p.method_name}</td>
                <td className="px-3 py-3 text-right font-semibold">{formatCurrency(p.grand_total)}</td>
                <td className="px-3 py-3"><StatusBadge status={p.status} /></td>
                <td className="px-3 py-3 text-right">
                  <button onClick={() => open(p)} className="inline-flex items-center gap-1 rounded-lg px-2.5 py-1.5 text-xs font-semibold text-brand-700 hover:bg-brand-50 dark:text-brand-400 dark:hover:bg-brand-950/40">
                    <Eye className="h-4 w-4" /> {p.status === 'WAITING_CONFIRMATION' ? 'Verifikasi' : 'Detail'}
                  </button>
                </td>
              </tr>
            ))}
          </Table>
          <Pagination meta={list.meta} onPage={list.setPage} />
        </Card>
      )}

      <Modal
        open={!!selected && !action}
        onClose={() => setSelected(null)}
        size="lg"
        title={selected ? `Pembayaran ${selected.payment_number}` : ''}
        subtitle={selected ? `Pesanan ${selected.order_number} · ${selected.customer_name}` : undefined}
        footer={
          selected && (
            <>
              {canRefund && (
                <button className={btnSecondary} onClick={() => startAction('refund')}>
                  <RotateCcw className="h-4 w-4" /> Refund
                </button>
              )}
              {canVerify && selected.status === 'WAITING_CONFIRMATION' && (
                <button className={btnDanger} onClick={() => startAction('reject')}>
                  <XCircle className="h-4 w-4" /> Tolak
                </button>
              )}
              {canVerify && (
                <button className={btnPrimary} onClick={() => startAction('approve')}>
                  <CheckCircle2 className="h-4 w-4" /> Setujui Pembayaran
                </button>
              )}
            </>
          )
        }
      >
        {selected && (
          <div className="space-y-4 text-sm">
            <div className="grid grid-cols-2 gap-3">
              <div className="rounded-xl bg-slate-50 dark:bg-slate-800/50 p-3">
                <div className="text-xs text-slate-500">Jumlah</div>
                <div className="text-lg font-bold">{formatCurrency(selected.grand_total)}</div>
              </div>
              <div className="rounded-xl bg-slate-50 dark:bg-slate-800/50 p-3">
                <div className="text-xs text-slate-500">Status</div>
                <div className="mt-1"><StatusBadge status={selected.status} /></div>
              </div>
            </div>

            {selected.method_type === 'MANUAL_TRANSFER' && (
              <div>
                <div className="mb-2 text-xs font-semibold text-slate-500">Bukti transfer</div>
                {selected.proof_image_url ? (
                  <a href={selected.proof_image_url} target="_blank" rel="noreferrer" className="block">
                    <img
                      src={selected.proof_image_url}
                      alt="Bukti transfer"
                      className="max-h-80 w-full rounded-xl border border-slate-200 dark:border-slate-800 object-contain bg-slate-50 dark:bg-slate-800"
                    />
                    <span className="mt-1 block text-xs text-brand-600">Buka gambar ukuran penuh</span>
                  </a>
                ) : (
                  <div className="flex items-center gap-2 rounded-xl border border-dashed border-slate-300 dark:border-slate-700 p-4 text-xs text-slate-500">
                    <ImageOff className="h-4 w-4" /> Pelanggan belum mengunggah bukti. Setujui hanya jika dana sudah terlihat di rekening toko.
                  </div>
                )}
              </div>
            )}

            {selected.notes && (
              <div className="rounded-xl bg-amber-50 dark:bg-amber-950/30 p-3 text-xs text-amber-800 dark:text-amber-300">Catatan: {selected.notes}</div>
            )}

            <div>
              <div className="mb-2 text-xs font-semibold text-slate-500">Riwayat</div>
              <ul className="space-y-2">
                {selected.logs.map((l) => (
                  <li key={l.id} className="rounded-lg border border-slate-100 dark:border-slate-800 p-2.5 text-xs">
                    <div className="flex items-center justify-between gap-2">
                      <span className="font-semibold">{statusLabel(l.old_status)} → {statusLabel(l.new_status)}</span>
                      <span className="text-slate-400">{formatDate(l.created_at)}</span>
                    </div>
                    {l.notes && <div className="mt-0.5 text-slate-500">{l.notes}</div>}
                  </li>
                ))}
              </ul>
            </div>
          </div>
        )}
      </Modal>

      <Modal
        open={!!selected && !!action}
        onClose={() => setAction(null)}
        size="sm"
        title={action === 'approve' ? 'Setujui pembayaran?' : action === 'reject' ? 'Tolak bukti transfer?' : 'Refund pembayaran'}
        subtitle={
          action === 'approve'
            ? 'Pastikan dana sudah masuk ke rekening toko. Pesanan akan langsung siap diproses.'
            : action === 'reject'
              ? 'Pelanggan akan diberi tahu dan bisa mengunggah ulang bukti transfer.'
              : 'Catat pengembalian dana yang sudah Anda transfer ke pelanggan.'
        }
        footer={
          <>
            <button className={btnSecondary} onClick={() => setAction(null)}>Kembali</button>
            <button className={action === 'approve' ? btnPrimary : btnDanger} disabled={busy} onClick={submit}>
              {action === 'approve' ? 'Setujui' : action === 'reject' ? 'Tolak' : 'Simpan Refund'}
            </button>
          </>
        }
      >
        <div className="space-y-3">
          {action === 'refund' && (
            <Field label="Jumlah refund (Rp)">
              <input type="number" min={1} className={inputClass} value={refundAmount} onChange={(e) => setRefundAmount(e.target.value)} />
            </Field>
          )}
          <Field label={action === 'approve' ? 'Catatan (opsional)' : 'Alasan'}>
            <textarea
              className={inputClass}
              rows={3}
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              placeholder={action === 'reject' ? 'Misal: nominal tidak sesuai / bukti tidak terbaca' : ''}
            />
          </Field>
        </div>
      </Modal>
    </div>
  );
};
