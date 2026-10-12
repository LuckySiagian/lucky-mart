import React, { useCallback, useEffect, useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import { toast } from 'sonner';
import { ArrowLeft, Check, Clock, Loader2, MapPin, MessageCircle, Package, Store, Truck, Upload, XCircle } from 'lucide-react';
import { orderService, paymentService, shippingService } from '../../services/apiServices';
import { OrderStatusBadge, Panel, PaymentStatusBadge, inputClass } from '../../components/common/Badge';
import { ImageUpload } from '../../components/common/ImageUpload';
import { ErrorState } from '../../components/common/ErrorState';
import { formatCurrency, formatDate, getErrorMessage, orderStatusInfo, SHIPMENT_STATUS, whatsappLink } from '../../lib/utils';
import { useStoreInfo } from '../../hooks/useStoreInfo';
import type { Order, PaymentMethod, Shipment } from '../../types';

const STEP_INDEX: Record<string, number> = { CREATED: 0, PAID: 0, CONFIRMED: 0, PROCESSING_PACKING: 1, OUT_FOR_DELIVERY: 2, COMPLETED: 3 };

/** Four simple steps so customers see at a glance where their order is. */
const OrderSteps: React.FC<{ order: Order }> = ({ order }) => {
  const pickup = order.delivery_type === 'PICKUP';
  const steps = ['Pesanan Diterima', 'Disiapkan', pickup ? 'Siap Diambil' : 'Diantar', 'Selesai'];
  const current = STEP_INDEX[order.status] ?? -1;
  return (
    <ol className="grid grid-cols-4 gap-1">
      {steps.map((label, i) => {
        const done = i <= current;
        return (
          <li key={label} className="flex flex-col items-center gap-1.5 text-center">
            <div className="flex w-full items-center">
              <div className={`h-1 flex-1 rounded ${i === 0 ? 'invisible' : done ? 'bg-primary-500' : 'bg-slate-200 dark:bg-slate-700'}`} />
              <span
                className={`flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-full text-sm font-bold ${
                  done ? 'bg-primary-600 text-white' : 'bg-slate-200 text-slate-500 dark:bg-slate-700 dark:text-slate-400'
                } ${i === current && order.status !== 'COMPLETED' ? 'ring-4 ring-primary-200 dark:ring-primary-900' : ''}`}
              >
                {done ? <Check className="h-4 w-4" /> : i + 1}
              </span>
              <div className={`h-1 flex-1 rounded ${i === steps.length - 1 ? 'invisible' : i < current ? 'bg-primary-500' : 'bg-slate-200 dark:bg-slate-700'}`} />
            </div>
            <span className={`text-[11px] font-semibold sm:text-xs ${done ? 'text-slate-900 dark:text-white' : 'text-slate-400'}`}>{label}</span>
          </li>
        );
      })}
    </ol>
  );
};

export const OrderDetailPage: React.FC = () => {
  const { id = '' } = useParams();
  const store = useStoreInfo();
  const [order, setOrder] = useState<Order | null>(null);
  const [shipment, setShipment] = useState<Shipment | null>(null);
  const [methods, setMethods] = useState<PaymentMethod[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [proofUrl, setProofUrl] = useState('');
  const [busy, setBusy] = useState(false);
  const [cancelOpen, setCancelOpen] = useState(false);
  const [reason, setReason] = useState('');

  const load = useCallback(async () => {
    try {
      const o = await orderService.getOrderById(id);
      setOrder(o);
      shippingService.getShipmentByOrderId(id).then(setShipment).catch(() => setShipment(null));
    } catch (err) {
      setError(getErrorMessage(err, 'Pesanan tidak ditemukan.'));
    }
  }, [id]);

  useEffect(() => {
    load();
    paymentService.getMethods().then(setMethods).catch(() => undefined);
  }, [load]);

  if (error) {
    return (
      <div className="mx-auto max-w-3xl px-4 py-10">
        <ErrorState message={error} />
      </div>
    );
  }
  if (!order) {
    return (
      <div className="flex justify-center py-20">
        <Loader2 className="h-8 w-8 animate-spin text-primary-600" />
      </div>
    );
  }

  // The bill the customer still has to act on (newest first from the API).
  const payment = order.payments?.[0];
  const isTransfer = payment?.method_type === 'MANUAL_TRANSFER';
  const needsProof = isTransfer && (payment.status === 'PENDING' || payment.status === 'FAILED') && order.status !== 'CANCELLED';
  const transferInfo = methods.find((m) => m.code === payment?.method_code)?.description;

  const uploadProof = async () => {
    if (!payment) return;
    if (!proofUrl) return toast.error('Pilih foto bukti transfer terlebih dahulu.');
    setBusy(true);
    try {
      await paymentService.uploadProof({ payment_id: payment.id, proof_image_url: proofUrl });
      toast.success('Bukti transfer terkirim. Toko akan segera memverifikasi.');
      setProofUrl('');
      load();
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setBusy(false);
    }
  };

  const cancel = async () => {
    if (reason.trim().length < 3) return toast.error('Tulis alasan pembatalan.');
    setBusy(true);
    try {
      setOrder(await orderService.cancelOrder(order.id, reason.trim()));
      setCancelOpen(false);
      toast.success('Pesanan dibatalkan.');
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setBusy(false);
    }
  };

  return (
    <div className="mx-auto max-w-4xl space-y-5 px-4 py-8 sm:px-6 lg:px-8">
      <Link to="/orders" className="inline-flex items-center gap-1 text-sm font-semibold text-slate-500 hover:text-primary-600">
        <ArrowLeft className="h-4 w-4" /> Pesanan Saya
      </Link>

      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="font-mono text-xl font-black text-slate-900 dark:text-white">{order.order_number}</h1>
          <p className="text-xs text-slate-500">{formatDate(order.order_date)}</p>
        </div>
        <OrderStatusBadge status={order.status} deliveryType={order.delivery_type} />
      </div>

      {order.status !== 'CANCELLED' && order.status !== 'PENDING_PAYMENT' && (
        <div className="rounded-3xl border border-slate-200 bg-white p-5 shadow-sm dark:border-slate-800 dark:bg-slate-900">
          <OrderSteps order={order} />
          {order.delivery_type === 'PICKUP' && order.status === 'OUT_FOR_DELIVERY' && (
            <p className="mt-4 rounded-xl bg-sky-50 p-3 text-center text-sm font-semibold text-sky-800 dark:bg-sky-950/40 dark:text-sky-300">
              Pesanan sudah siap! Datang ke toko dan sebutkan nomor pesanan <span className="font-mono">{order.order_number}</span>.
            </p>
          )}
        </div>
      )}

      {payment && (needsProof || payment.status === 'WAITING_CONFIRMATION') && (
        <Panel title="Pembayaran transfer">
          <div className="text-sm">
            <div className="flex items-center justify-between">
              <span>Total yang harus ditransfer</span>
              <span className="text-xl font-black text-primary-600">{formatCurrency(payment.grand_total)}</span>
            </div>
            {transferInfo && <div className="mt-3 whitespace-pre-line rounded-xl bg-slate-50 p-3 font-semibold dark:bg-slate-800">{transferInfo}</div>}
            {payment.status === 'FAILED' && payment.notes && (
              <div className="mt-3 rounded-xl bg-rose-50 p-3 text-xs text-rose-700 dark:bg-rose-950/40 dark:text-rose-300">{payment.notes}</div>
            )}
            {payment.status === 'WAITING_CONFIRMATION' ? (
              <p className="mt-3 rounded-xl bg-sky-50 p-3 text-xs text-sky-800 dark:bg-sky-950/40 dark:text-sky-300">
                Bukti transfer sudah diterima dan sedang diverifikasi toko.
              </p>
            ) : (
              <div className="mt-4 space-y-2">
                <label className="text-xs font-semibold text-slate-500">Foto bukti transfer</label>
                <ImageUpload value={proofUrl} onChange={setProofUrl} label="Foto / pilih bukti transfer" />
                <button
                  onClick={uploadProof}
                  disabled={busy || !proofUrl}
                  className="flex w-full items-center justify-center gap-1.5 rounded-xl bg-primary-600 px-4 py-3 text-sm font-bold text-white hover:bg-primary-700 disabled:opacity-50 sm:w-auto"
                >
                  <Upload className="h-4 w-4" /> Kirim Bukti Transfer
                </button>
                <p className="text-[11px] text-slate-400">Foto struk ATM atau screenshot m-banking juga boleh.</p>
              </div>
            )}
          </div>
        </Panel>
      )}

      {order.delivery_type === 'PICKUP' ? (
        <Panel title={<><Store className="h-5 w-5 text-primary-600" /> Ambil sendiri di toko</>}>
          <p className="text-sm text-slate-600 dark:text-slate-400">
            Tidak ada ongkir. Kami kabari lewat notifikasi saat pesanan siap diambil.
          </p>
          {store && (store.address || store.open_hours || store.whatsapp) && (
            <div className="mt-3 space-y-1.5 rounded-xl bg-slate-50 p-3 text-sm dark:bg-slate-800/60">
              {store.address && (
                <div className="flex items-start gap-2">
                  <MapPin className="mt-0.5 h-4 w-4 flex-shrink-0 text-primary-600" />
                  <span>
                    {store.address}
                    {store.maps_url && (
                      <a href={store.maps_url} target="_blank" rel="noreferrer" className="ml-1 font-semibold text-primary-600 underline">Buka peta</a>
                    )}
                  </span>
                </div>
              )}
              {store.open_hours && <div className="flex items-center gap-2"><Clock className="h-4 w-4 text-primary-600" /> {store.open_hours}</div>}
              {whatsappLink(store.whatsapp) && (
                <a
                  href={whatsappLink(store.whatsapp, `Halo, saya mau tanya pesanan ${order.order_number}.`)!}
                  target="_blank"
                  rel="noreferrer"
                  className="mt-1 inline-flex items-center gap-1.5 rounded-lg bg-emerald-600 px-3 py-1.5 text-xs font-bold text-white hover:bg-emerald-700"
                >
                  <MessageCircle className="h-3.5 w-3.5" /> Chat toko
                </a>
              )}
            </div>
          )}
          {order.address?.notes && <div className="mt-1 text-xs italic text-slate-500">Catatan: {order.address.notes}</div>}
        </Panel>
      ) : (
        <div className="grid gap-5 md:grid-cols-2">
          <Panel title={<><MapPin className="h-5 w-5 text-primary-600" /> Alamat</>}>
            {order.address && (
              <div className="text-sm">
                <div className="font-bold">{order.address.recipient_name} · {order.address.phone}</div>
                <div className="text-slate-600 dark:text-slate-400">{order.address.full_address}, {order.address.city} {order.address.postal_code}</div>
                {order.address.notes && <div className="mt-1 text-xs italic text-slate-500">Catatan: {order.address.notes}</div>}
              </div>
            )}
          </Panel>
  
          <Panel title={<><Truck className="h-5 w-5 text-primary-600" /> Pengiriman</>}>
            {shipment ? (
              <div className="text-sm">
                <div className="flex items-center justify-between">
                  <span className="font-mono text-xs">{shipment.tracking_number}</span>
                  <span className="text-xs font-bold text-primary-600">{SHIPMENT_STATUS[shipment.status] || shipment.status}</span>
                </div>
                <ol className="mt-3 space-y-2 border-l-2 border-slate-200 pl-3 dark:border-slate-700">
                  {shipment.trackings.slice().reverse().map((t) => (
                    <li key={t.id} className="text-xs">
                      <div className="font-semibold">{SHIPMENT_STATUS[t.status] || t.status}</div>
                      <div className="text-slate-500">{[t.location !== '-' ? t.location : '', t.notes].filter(Boolean).join(' — ')}</div>
                      <div className="text-slate-400">{formatDate(t.created_at)}</div>
                    </li>
                  ))}
                </ol>
              </div>
            ) : (
              <p className="text-sm text-slate-500">Info pengantaran muncul saat pesanan mulai diantar.</p>
            )}
          </Panel>
        </div>
      )}

      <Panel title={<><Package className="h-5 w-5 text-primary-600" /> Barang</>}>
        <ul className="divide-y divide-slate-100 dark:divide-slate-800">
          {order.items.map((i) => (
            <li key={i.id} className="flex items-center justify-between gap-3 py-2.5 text-sm">
              <div>
                <Link to={`/products/${i.product_id}`} className="font-medium hover:text-primary-600">{i.product_name}</Link>
                <div className="text-xs text-slate-500">{i.quantity} {i.unit_name} × {formatCurrency(i.unit_price)}</div>
              </div>
              <span className="font-semibold">{formatCurrency(i.subtotal)}</span>
            </li>
          ))}
        </ul>
        <div className="mt-3 space-y-1 border-t border-slate-100 pt-3 text-sm dark:border-slate-800">
          <div className="flex justify-between text-slate-500"><span>Subtotal</span><span>{formatCurrency(order.total_amount)}</span></div>
          {order.delivery_type !== 'PICKUP' && <div className="flex justify-between text-slate-500"><span>Ongkir</span><span>{formatCurrency(order.shipping_fee)}</span></div>}
          <div className="flex justify-between text-base font-black"><span>Total</span><span className="text-primary-600">{formatCurrency(order.grand_total)}</span></div>
          {payment && (
            <div className="flex items-center justify-between pt-1 text-xs">
              <span className="text-slate-500">
                {payment.method_type === 'COD'
                  ? order.delivery_type === 'PICKUP'
                    ? 'Bayar tunai saat ambil di toko'
                    : 'Bayar tunai saat barang sampai'
                  : payment.method_name}
              </span>
              <PaymentStatusBadge status={payment.status} />
            </div>
          )}
          {order.payment_method === 'B2B_CREDIT_TEMPO' && order.due_date && (
            <div className="text-xs text-slate-500">Jatuh tempo kasbon: {formatDate(order.due_date)}</div>
          )}
        </div>
      </Panel>

      <Panel title="Riwayat">
        <ol className="space-y-3 border-l-2 border-slate-200 pl-4 dark:border-slate-700">
          {order.histories.slice().reverse().map((h) => (
            <li key={h.id} className="relative text-sm">
              <span className="absolute -left-[21px] top-1.5 h-2.5 w-2.5 rounded-full bg-primary-500" />
              <div className="font-semibold">{orderStatusInfo(h.status, order.delivery_type).label}</div>
              {h.notes && <div className="text-xs text-slate-500">{h.notes}</div>}
              <div className="text-[11px] text-slate-400">{formatDate(h.created_at)}</div>
            </li>
          ))}
        </ol>
      </Panel>

      {order.can_cancel && (
        <div className="rounded-3xl border border-rose-200 bg-rose-50/50 p-5 dark:border-rose-900 dark:bg-rose-950/20">
          {!cancelOpen ? (
            <button onClick={() => setCancelOpen(true)} className="flex items-center gap-2 text-sm font-bold text-rose-600">
              <XCircle className="h-4 w-4" /> Batalkan pesanan
            </button>
          ) : (
            <div className="space-y-2">
              <label className="text-sm font-semibold text-rose-700 dark:text-rose-300">Kenapa pesanan dibatalkan?</label>
              <input className={inputClass} value={reason} onChange={(e) => setReason(e.target.value)} placeholder="Misal: salah pilih barang" />
              <div className="flex gap-2">
                <button onClick={() => setCancelOpen(false)} className="rounded-xl px-4 py-2 text-sm font-semibold text-slate-600">Tidak jadi</button>
                <button onClick={cancel} disabled={busy} className="rounded-xl bg-rose-600 px-4 py-2 text-sm font-bold text-white disabled:opacity-50">Ya, batalkan</button>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
};
