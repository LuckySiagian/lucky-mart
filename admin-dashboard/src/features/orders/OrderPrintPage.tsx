import React, { useEffect, useState } from 'react';
import { useParams, useSearchParams } from 'react-router-dom';
import { Loader2, Printer } from 'lucide-react';
import { orderService, storeService } from '../../services/apiServices';
import { DELIVERY_LABELS, formatCurrency, formatDate, getErrorMessage, paymentMethodLabel, statusLabel } from '../../lib/utils';
import type { Order, StoreInfo } from '../../types';

/**
 * Printable receipt (?mode=nota) or packing list (?mode=packing) for one
 * order. Narrow on purpose so it fits a 58/80 mm thermal printer, and still
 * prints fine on A4.
 */
export const OrderPrintPage: React.FC = () => {
  const { id = '' } = useParams();
  const [params] = useSearchParams();
  const packing = params.get('mode') === 'packing';
  const [order, setOrder] = useState<Order | null>(null);
  const [store, setStore] = useState<StoreInfo | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    Promise.all([orderService.getOrderById(id), storeService.get().catch(() => null)])
      .then(([o, s]) => {
        setOrder(o);
        setStore(s);
      })
      .catch((err) => setError(getErrorMessage(err, 'Pesanan tidak ditemukan.')));
  }, [id]);

  // Open the print dialog once everything is on screen.
  useEffect(() => {
    if (!order) return;
    document.title = `${packing ? 'Daftar barang' : 'Nota'} ${order.order_number}`;
    const t = setTimeout(() => window.print(), 400);
    return () => clearTimeout(t);
  }, [order, packing]);

  if (error) return <div className="p-6 text-sm text-rose-600">{error}</div>;
  if (!order) {
    return (
      <div className="flex justify-center p-10">
        <Loader2 className="h-6 w-6 animate-spin text-slate-400" />
      </div>
    );
  }

  const pickup = order.delivery_type === 'PICKUP';
  const payment = order.payments?.[0];

  return (
    <div className="min-h-screen bg-slate-100 py-6 print:bg-white print:py-0">
      <style>{`@media print { @page { margin: 4mm; } .no-print { display: none !important; } }`}</style>

      <div className="no-print mx-auto mb-4 flex max-w-[80mm] justify-center gap-2">
        <button onClick={() => window.print()} className="inline-flex items-center gap-2 rounded-lg bg-slate-900 px-4 py-2 text-sm font-semibold text-white">
          <Printer className="h-4 w-4" /> Cetak
        </button>
      </div>

      <div className="mx-auto max-w-[80mm] bg-white p-4 font-mono text-[12px] leading-snug text-black shadow print:max-w-none print:p-0 print:shadow-none">
        <div className="text-center">
          <div className="text-[15px] font-bold">{store?.name || 'Lucky Mart'}</div>
          {store?.address && <div>{store.address}</div>}
          {store?.whatsapp && <div>WA {store.whatsapp}</div>}
        </div>

        <Divider />
        <div className="text-center font-bold">{packing ? 'DAFTAR BARANG' : 'NOTA PENJUALAN'}</div>
        <Divider />

        <Row left="No" right={order.order_number} />
        <Row left="Tanggal" right={formatDate(order.order_date)} />
        <Row left="Pelanggan" right={order.customer?.full_name || order.address?.recipient_name || '-'} />
        <Row left="HP" right={order.customer?.phone_number || order.address?.phone || '-'} />
        <Row left="Terima" right={DELIVERY_LABELS[order.delivery_type] || 'Diantar'} />
        {!pickup && order.address && <div className="mt-1">Alamat: {order.address.full_address}, {order.address.city}</div>}
        {order.address?.notes && <div className="mt-1">Catatan: {order.address.notes}</div>}

        <Divider />
        {order.items.map((it) => (
          <div key={it.id} className="mb-1.5">
            {packing ? (
              <div className="flex items-start gap-2">
                <span className="mt-0.5 inline-block h-3.5 w-3.5 flex-shrink-0 border border-black" />
                <span>
                  <b>{it.quantity} {it.unit_name}</b> {it.product_name}
                </span>
              </div>
            ) : (
              <>
                <div>{it.product_name}</div>
                <Row left={`  ${it.quantity} ${it.unit_name} x ${formatCurrency(it.unit_price)}`} right={formatCurrency(it.subtotal)} />
              </>
            )}
          </div>
        ))}
        <Divider />

        {packing ? (
          <Row left="Jumlah barang" right={`${order.total_items}`} />
        ) : (
          <>
            <Row left="Subtotal" right={formatCurrency(order.total_amount)} />
            {order.shipping_fee > 0 && <Row left="Ongkir" right={formatCurrency(order.shipping_fee)} />}
            <Row left={<b>TOTAL</b>} right={<b>{formatCurrency(order.grand_total)}</b>} />
            <Row left="Bayar" right={paymentMethodLabel(order.payment_method)} />
            {payment && <Row left="Status" right={statusLabel(payment.status)} />}
            {order.payment_method === 'B2B_CREDIT_TEMPO' && order.due_date && <Row left="Jatuh tempo" right={formatDate(order.due_date)} />}
          </>
        )}

        <Divider />
        <div className="text-center">{packing ? 'Disiapkan oleh: ____________' : store?.receipt_note || 'Terima kasih sudah berbelanja!'}</div>
      </div>
    </div>
  );
};

const Divider = () => <div className="my-2 border-t border-dashed border-black" />;

const Row: React.FC<{ left: React.ReactNode; right: React.ReactNode }> = ({ left, right }) => (
  <div className="flex justify-between gap-2">
    <span className="whitespace-pre">{left}</span>
    <span className="text-right">{right}</span>
  </div>
);
