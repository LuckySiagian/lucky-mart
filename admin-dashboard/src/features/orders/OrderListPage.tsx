import React, { useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { toast } from 'sonner';
import { ArrowRight, CheckCheck, ClipboardList, ExternalLink, Loader2, MapPin, MessageCircle, Package, Phone, Printer, Store, Truck, XCircle } from 'lucide-react';
import { PageHeader } from '../../components/common/PageHeader';
import { StatusBadge } from '../../components/common/StatusBadge';
import { LoadingSkeleton } from '../../components/common/LoadingSkeleton';
import { EmptyState } from '../../components/common/EmptyState';
import { ErrorState } from '../../components/common/ErrorState';
import { Modal, btnDanger, btnPrimary, btnSecondary, inputClass } from '../../components/common/Modal';
import { Card, FilterTabs, Pagination, SearchInput } from '../../components/common/Pagination';
import { orderService, paymentService } from '../../services/apiServices';
import { useList } from '../../hooks/useList';
import { hasRole, MANAGER_ROLES, useAuthStore } from '../../store/useAuthStore';
import {
  DELIVERY_LABELS,
  formatCurrency,
  formatDate,
  getErrorMessage,
  orderStatusLabel,
  paymentMethodLabel,
  whatsappLink,
} from '../../lib/utils';
import type { Order, OrderStatus } from '../../types';

// Tabs follow the work to do, not every internal status.
const TABS = [
  { key: 'proses', value: 'CREATED,CONFIRMED,PAID,PROCESSING_PACKING', label: 'Perlu Diproses' },
  { key: 'bayar', value: 'PENDING_PAYMENT', label: 'Menunggu Transfer' },
  { key: 'jalan', value: 'OUT_FOR_DELIVERY', label: 'Diantar / Siap Diambil' },
  { key: 'selesai', value: 'COMPLETED', label: 'Selesai' },
  { key: 'batal', value: 'CANCELLED', label: 'Batal' },
  { key: 'semua', value: '', label: 'Semua' },
];

/** Button text for moving an order to `next`, worded for the store owner. */
function actionLabel(next: OrderStatus, o: Order): string {
  const pickup = o.delivery_type === 'PICKUP';
  switch (next) {
    case 'CONFIRMED':
      return 'Terima Pesanan';
    case 'PROCESSING_PACKING':
      return 'Siapkan Barang';
    case 'OUT_FOR_DELIVERY':
      return pickup ? 'Siap Diambil' : 'Antar Sekarang';
    case 'COMPLETED':
      return pickup ? 'Sudah Diambil' : 'Sudah Sampai';
    default:
      return orderStatusLabel(next, o.delivery_type);
  }
}

export const OrderListPage: React.FC = () => {
  const user = useAuthStore((s) => s.user);
  const isManager = hasRole(user, MANAGER_ROLES);
  const isDriver = user?.role === 'DRIVER';
  const [searchParams] = useSearchParams();
  const initialTab = TABS.find((t) => t.key === searchParams.get('tab')) || TABS[0];

  const list = useList((p) => orderService.getOrders(p), { status: initialTab.value });
  const [selected, setSelected] = useState<Order | null>(null);
  const [detailLoading, setDetailLoading] = useState(false);
  const [busyId, setBusyId] = useState<string | null>(null);
  const [cancelOpen, setCancelOpen] = useState(false);
  const [cancelReason, setCancelReason] = useState('');
  const [rejectOpen, setRejectOpen] = useState(false);
  const [rejectReason, setRejectReason] = useState('');

  // Steps this user may take, cancel excluded (it has its own button).
  const nextSteps = (o: Order) =>
    o.next_statuses.filter((s) => s !== 'CANCELLED' && (!isDriver || s === 'COMPLETED'));

  const openDetail = async (id: string) => {
    setDetailLoading(true);
    try {
      setSelected(await orderService.getOrderById(id));
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setDetailLoading(false);
    }
  };

  const changeStatus = async (o: Order, status: OrderStatus) => {
    setBusyId(o.id);
    try {
      const updated = await orderService.updateOrderStatus(o.id, status);
      if (selected?.id === o.id) setSelected(updated);
      toast.success(`${updated.order_number}: ${orderStatusLabel(updated.status, updated.delivery_type)}`);
      list.reload();
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setBusyId(null);
    }
  };

  const cancelOrder = async () => {
    if (!selected) return;
    if (cancelReason.trim().length < 3) return toast.error('Tulis alasan pembatalan (minimal 3 huruf).');
    setBusyId(selected.id);
    try {
      const updated = await orderService.cancelOrder(selected.id, cancelReason.trim());
      setSelected(updated);
      setCancelOpen(false);
      setCancelReason('');
      toast.success('Pesanan dibatalkan dan stok dikembalikan.');
      list.reload();
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setBusyId(null);
    }
  };

  const proofPayment = selected?.payments?.find((p) => p.status === 'WAITING_CONFIRMATION');

  const approvePayment = async () => {
    if (!selected || !proofPayment) return;
    setBusyId(selected.id);
    try {
      await paymentService.approvePayment(proofPayment.id, 'Transfer sudah masuk');
      setSelected(await orderService.getOrderById(selected.id));
      toast.success('Pembayaran diterima. Pesanan siap diproses.');
      list.reload();
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setBusyId(null);
    }
  };

  const rejectPayment = async () => {
    if (!selected || !proofPayment) return;
    if (rejectReason.trim().length < 3) return toast.error('Tulis alasan penolakan agar pelanggan tahu.');
    setBusyId(selected.id);
    try {
      await paymentService.rejectPayment(proofPayment.id, rejectReason.trim());
      setSelected(await orderService.getOrderById(selected.id));
      setRejectOpen(false);
      setRejectReason('');
      toast.success('Bukti transfer ditolak. Pelanggan diminta mengirim ulang.');
      list.reload();
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setBusyId(null);
    }
  };

  const waLink = (o: Order) =>
    whatsappLink(
      o.customer?.phone_number || o.address?.phone,
      `Halo ${o.customer?.full_name || ''}, ini dari Lucky Mart soal pesanan ${o.order_number}.`
    );

  const activeTab = TABS.find((t) => t.value === ((list.params.status as string) || '')) || TABS[TABS.length - 1];

  return (
    <div className="space-y-5">
      <PageHeader title="Pesanan" subtitle="Klik tombol hijau untuk lanjut ke langkah berikutnya." />

      <div className="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
        <FilterTabs value={activeTab.value} onChange={(v) => list.setFilter('status', v)} options={TABS} />
        <SearchInput value={list.searchInput} onChange={list.setSearchInput} placeholder="Cari nama, nomor HP, atau nomor pesanan..." />
      </div>

      {list.loading ? (
        <LoadingSkeleton count={5} />
      ) : list.error ? (
        <ErrorState message={list.error} onRetry={list.reload} />
      ) : list.items.length === 0 ? (
        <EmptyState title="Tidak ada pesanan" description={activeTab.key === 'proses' ? 'Semua pesanan sudah diproses. Mantap!' : 'Belum ada pesanan di sini.'} />
      ) : (
        <div className="space-y-3">
          {list.items.map((o) => {
            const next = nextSteps(o)[0];
            const pickup = o.delivery_type === 'PICKUP';
            return (
              <Card key={o.id} className="p-4">
                <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                  <button onClick={() => openDetail(o.id)} className="min-w-0 flex-1 text-left">
                    <div className="flex flex-wrap items-center gap-2">
                      <span className="font-bold text-slate-900 dark:text-white">{o.customer?.full_name || o.address?.recipient_name || '-'}</span>
                      <StatusBadge status={o.status} label={orderStatusLabel(o.status, o.delivery_type)} />
                      <span className="inline-flex items-center gap-1 rounded-full bg-slate-100 px-2 py-0.5 text-[11px] font-semibold text-slate-600 dark:bg-slate-800 dark:text-slate-300">
                        {pickup ? <Store className="h-3 w-3" /> : <Truck className="h-3 w-3" />}
                        {DELIVERY_LABELS[o.delivery_type] || 'Diantar'}
                      </span>
                    </div>
                    <div className="mt-1 text-xs text-slate-500">
                      {o.order_number} · {formatDate(o.order_date)} · {o.total_items} barang · {paymentMethodLabel(o.payment_method)}
                    </div>
                    {!pickup && o.address && (
                      <div className="mt-0.5 truncate text-xs text-slate-500">
                        <MapPin className="mr-1 inline h-3 w-3" />
                        {o.address.full_address}
                      </div>
                    )}
                  </button>
                  <div className="flex items-center justify-between gap-3 sm:justify-end">
                    <div className="text-base font-bold text-slate-900 dark:text-white">{formatCurrency(o.grand_total)}</div>
                    {next ? (
                      <button className={btnPrimary} disabled={busyId === o.id} onClick={() => changeStatus(o, next)}>
                        {busyId === o.id ? <Loader2 className="h-4 w-4 animate-spin" /> : <ArrowRight className="h-4 w-4" />}
                        {actionLabel(next, o)}
                      </button>
                    ) : (
                      <button className={btnSecondary} onClick={() => openDetail(o.id)}>
                        {o.status === 'PENDING_PAYMENT' ? 'Cek Pembayaran' : 'Detail'}
                      </button>
                    )}
                  </div>
                </div>
              </Card>
            );
          })}
          <Pagination meta={list.meta} onPage={list.setPage} />
        </div>
      )}

      {detailLoading && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/40">
          <Loader2 className="h-8 w-8 animate-spin text-white" />
        </div>
      )}

      <Modal
        open={!!selected}
        onClose={() => setSelected(null)}
        size="lg"
        title={selected ? selected.customer?.full_name || selected.order_number : ''}
        subtitle={selected ? `${selected.order_number} · ${formatDate(selected.order_date)}` : undefined}
        footer={
          selected && (
            <>
              {selected.can_cancel && !isDriver && (
                <button className={btnDanger} disabled={busyId === selected.id} onClick={() => setCancelOpen(true)}>
                  <XCircle className="h-4 w-4" /> Batalkan
                </button>
              )}
              {nextSteps(selected)
                .slice()
                .reverse()
                .map((s, i, arr) => (
                  <button
                    key={s}
                    className={i === arr.length - 1 ? btnPrimary : btnSecondary}
                    disabled={busyId === selected.id}
                    onClick={() => changeStatus(selected, s)}
                  >
                    {s === 'COMPLETED' ? <CheckCheck className="h-4 w-4" /> : <ArrowRight className="h-4 w-4" />}
                    {actionLabel(s, selected)}
                  </button>
                ))}
            </>
          )
        }
      >
        {selected && (
          <div className="space-y-5 text-sm">
            <div className="flex flex-wrap items-center gap-2">
              <StatusBadge status={selected.status} label={orderStatusLabel(selected.status, selected.delivery_type)} />
              <span className="text-xs text-slate-500">
                {DELIVERY_LABELS[selected.delivery_type] || 'Diantar'} · {paymentMethodLabel(selected.payment_method)}
              </span>
              <span className="ml-auto flex gap-1.5">
                <a href={`/orders/${selected.id}/print?mode=packing`} target="_blank" rel="noreferrer" className={btnSecondary}>
                  <ClipboardList className="h-4 w-4" /> Daftar Barang
                </a>
                <a href={`/orders/${selected.id}/print`} target="_blank" rel="noreferrer" className={btnSecondary}>
                  <Printer className="h-4 w-4" /> Cetak Nota
                </a>
              </span>
            </div>

            {selected.status === 'PENDING_PAYMENT' && (
              <div className="rounded-xl border border-amber-200 bg-amber-50 p-3 dark:border-amber-900 dark:bg-amber-950/30">
                {proofPayment ? (
                  <div className="space-y-3">
                    <div className="font-semibold text-amber-900 dark:text-amber-200">
                      Pelanggan sudah kirim bukti transfer {formatCurrency(proofPayment.grand_total)}. Cek rekening toko, lalu:
                    </div>
                    {proofPayment.proof_image_url && (
                      <a href={proofPayment.proof_image_url} target="_blank" rel="noreferrer" className="block w-fit" title="Buka ukuran penuh">
                        <img
                          src={proofPayment.proof_image_url}
                          alt="Bukti transfer"
                          className="max-h-64 rounded-lg border border-amber-200 bg-white object-contain dark:border-amber-900"
                        />
                        <span className="mt-1 inline-flex items-center gap-1 text-xs font-semibold text-brand-700 underline dark:text-brand-400">
                          <ExternalLink className="h-3.5 w-3.5" /> Buka ukuran penuh
                        </span>
                      </a>
                    )}
                    {isManager ? (
                      <div className="flex flex-wrap gap-2">
                        <button className={btnPrimary} disabled={busyId === selected.id} onClick={approvePayment}>
                          <CheckCheck className="h-4 w-4" /> Uang Sudah Masuk
                        </button>
                        <button className={btnSecondary} disabled={busyId === selected.id} onClick={() => setRejectOpen(true)}>
                          Belum Masuk / Tolak
                        </button>
                      </div>
                    ) : (
                      <div className="text-xs text-amber-800 dark:text-amber-300">Minta pemilik atau manajer toko mengecek pembayaran ini.</div>
                    )}
                  </div>
                ) : (
                  <div className="text-amber-900 dark:text-amber-200">Menunggu pelanggan transfer dan mengirim bukti. Pesanan diproses setelah uang masuk.</div>
                )}
              </div>
            )}

            <div className="grid gap-3 sm:grid-cols-2">
              <div className="rounded-xl bg-slate-50 p-3 dark:bg-slate-800/50">
                <div className="mb-1 text-xs font-semibold text-slate-500">Pelanggan</div>
                <div className="font-semibold">{selected.customer?.full_name}</div>
                <div className="flex items-center gap-1 text-xs text-slate-500">
                  <Phone className="h-3 w-3" /> {selected.customer?.phone_number}
                </div>
                {waLink(selected) && (
                  <a href={waLink(selected)!} target="_blank" rel="noreferrer" className="mt-2 inline-flex items-center gap-1.5 rounded-lg bg-emerald-600 px-3 py-1.5 text-xs font-bold text-white hover:bg-emerald-700">
                    <MessageCircle className="h-3.5 w-3.5" /> Chat WhatsApp
                  </a>
                )}
              </div>
              <div className="rounded-xl bg-slate-50 p-3 dark:bg-slate-800/50">
                {selected.delivery_type === 'PICKUP' ? (
                  <>
                    <div className="mb-1 text-xs font-semibold text-slate-500">Cara terima</div>
                    <div className="flex items-center gap-1.5 font-semibold"><Store className="h-4 w-4" /> Diambil sendiri di toko</div>
                  </>
                ) : (
                  selected.address && (
                    <>
                      <div className="mb-1 text-xs font-semibold text-slate-500">Antar ke</div>
                      <div className="font-semibold">{selected.address.recipient_name} · {selected.address.phone}</div>
                      <div className="flex items-start gap-1 text-xs text-slate-500">
                        <MapPin className="mt-0.5 h-3 w-3 flex-shrink-0" /> {selected.address.full_address}, {selected.address.city} {selected.address.postal_code}
                      </div>
                    </>
                  )
                )}
                {selected.address?.notes && <div className="mt-1 text-xs italic text-slate-500">Catatan: {selected.address.notes}</div>}
              </div>
            </div>

            <div>
              <div className="mb-2 text-xs font-semibold text-slate-500">Barang yang disiapkan</div>
              <div className="divide-y divide-slate-100 rounded-xl border border-slate-200 dark:divide-slate-800 dark:border-slate-800">
                {selected.items.map((it) => (
                  <div key={it.id} className="flex items-center justify-between gap-3 p-3">
                    <div className="flex items-center gap-3">
                      <Package className="h-4 w-4 text-slate-400" />
                      <div>
                        <div className="font-medium">{it.product_name}</div>
                        <div className="text-xs text-slate-500">
                          <b className="text-slate-700 dark:text-slate-300">{it.quantity} {it.unit_name}</b> × {formatCurrency(it.unit_price)}
                        </div>
                      </div>
                    </div>
                    <div className="font-semibold">{formatCurrency(it.subtotal)}</div>
                  </div>
                ))}
              </div>
              <div className="mt-3 space-y-1 text-right text-sm">
                <div className="text-slate-500">Belanja {formatCurrency(selected.total_amount)}</div>
                {selected.shipping_fee > 0 && <div className="text-slate-500">Ongkir {formatCurrency(selected.shipping_fee)}</div>}
                <div className="text-base font-bold">Total {formatCurrency(selected.grand_total)}</div>
              </div>
            </div>

            <details className="rounded-xl border border-slate-200 p-3 dark:border-slate-800">
              <summary className="cursor-pointer text-xs font-semibold text-slate-500">Riwayat pesanan</summary>
              <ol className="mt-3 space-y-3 border-l-2 border-slate-200 pl-4 dark:border-slate-800">
                {selected.histories.map((h) => (
                  <li key={h.id} className="relative">
                    <span className="absolute -left-[21px] top-1 h-2.5 w-2.5 rounded-full bg-brand-500" />
                    <div className="flex flex-wrap items-center gap-2">
                      <StatusBadge status={h.status} label={orderStatusLabel(h.status, selected.delivery_type)} />
                      <span className="text-xs text-slate-400">{formatDate(h.created_at)}</span>
                    </div>
                    {h.notes && <p className="mt-1 text-xs text-slate-600 dark:text-slate-400">{h.notes}</p>}
                  </li>
                ))}
              </ol>
            </details>
          </div>
        )}
      </Modal>

      <Modal
        open={cancelOpen}
        onClose={() => setCancelOpen(false)}
        size="sm"
        title="Batalkan pesanan?"
        subtitle="Stok barang dan kasbon dikembalikan otomatis. Kalau pelanggan sudah transfer, kembalikan uangnya lewat menu Pembayaran."
        footer={
          <>
            <button className={btnSecondary} onClick={() => setCancelOpen(false)}>Kembali</button>
            <button className={btnDanger} disabled={!!busyId} onClick={cancelOrder}>Ya, Batalkan</button>
          </>
        }
      >
        <textarea
          className={inputClass}
          rows={3}
          placeholder="Alasan pembatalan, misal: stok habis / pelanggan minta batal"
          value={cancelReason}
          onChange={(e) => setCancelReason(e.target.value)}
        />
      </Modal>

      <Modal
        open={rejectOpen}
        onClose={() => setRejectOpen(false)}
        size="sm"
        title="Tolak bukti transfer?"
        subtitle="Pelanggan akan diminta mengirim bukti transfer yang benar."
        footer={
          <>
            <button className={btnSecondary} onClick={() => setRejectOpen(false)}>Kembali</button>
            <button className={btnDanger} disabled={!!busyId} onClick={rejectPayment}>Tolak</button>
          </>
        }
      >
        <textarea
          className={inputClass}
          rows={3}
          placeholder="Alasan, misal: uang belum masuk / nominal kurang"
          value={rejectReason}
          onChange={(e) => setRejectReason(e.target.value)}
        />
      </Modal>
    </div>
  );
};
