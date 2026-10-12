import React, { useEffect, useMemo, useRef, useState } from 'react';
import { Link, Navigate, useNavigate } from 'react-router-dom';
import { toast } from 'sonner';
import { AlertTriangle, CheckCircle2, CreditCard, Loader2, MapPin, Package, Store, Truck } from 'lucide-react';
import { addressService, cartService, creditService, orderService, paymentService, shippingService, type DeliveryType } from '../../services/apiServices';
import { useCartStore } from '../../store/useCartStore';
import { useAuthStore } from '../../store/useAuthStore';
import { useStoreInfo } from '../../hooks/useStoreInfo';
import { Field, Panel, inputClass } from '../../components/common/Badge';
import { formatCurrency, getErrorMessage, isPhone } from '../../lib/utils';
import type { Address, CheckoutSummary, CreditLedger, PaymentMethod, ShippingMethod } from '../../types';

const NEW_ADDRESS = 'new';

export const CheckoutPage: React.FC = () => {
  const navigate = useNavigate();
  const user = useAuthStore((s) => s.user);
  const store = useStoreInfo();
  const { selectedIds, setSelected, refresh: refreshCart } = useCartStore();

  const [addresses, setAddresses] = useState<Address[]>([]);
  const [payments, setPayments] = useState<PaymentMethod[]>([]);
  const [shippings, setShippings] = useState<ShippingMethod[]>([]);
  const [credit, setCredit] = useState<CreditLedger | null>(null);
  const [summary, setSummary] = useState<CheckoutSummary | null>(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);

  const [deliveryType, setDeliveryType] = useState<DeliveryType>('DELIVERY');
  const [addressId, setAddressId] = useState<string>(NEW_ADDRESS);
  const [newAddr, setNewAddr] = useState({
    title: 'Rumah',
    recipient_name: user?.full_name || '',
    phone: user?.phone_number || '',
    full_address: '',
    city: '',
    province: '',
    postal_code: '',
  });
  const [saveAddress, setSaveAddress] = useState(true);
  const [shippingId, setShippingId] = useState('');
  const [paymentCode, setPaymentCode] = useState('');
  const [notes, setNotes] = useState('');
  const placedRef = useRef(false);

  const isApprovedB2B = user?.role === 'B2B_CUSTOMER' && user.b2b_profile?.verification_status === 'APPROVED';

  useEffect(() => {
    (async () => {
      const [addrRes, payRes, shipRes, creditRes] = await Promise.allSettled([
        addressService.getAddresses(),
        paymentService.getMethods(),
        shippingService.getMethods(),
        isApprovedB2B ? creditService.getMine() : Promise.reject(new Error('not b2b')),
      ]);
      if (addrRes.status === 'fulfilled') {
        setAddresses(addrRes.value);
        const def = addrRes.value.find((a) => a.is_default) || addrRes.value[0];
        if (def) setAddressId(def.id);
      }
      if (payRes.status === 'fulfilled') {
        const usable = payRes.value.filter((m) => m.type !== 'B2B_CREDIT' || isApprovedB2B);
        setPayments(usable);
        setPaymentCode(usable[0]?.code || '');
      }
      if (shipRes.status === 'fulfilled') {
        setShippings(shipRes.value);
        setShippingId(shipRes.value[0]?.id || '');
      }
      if (creditRes.status === 'fulfilled') setCredit(creditRes.value);
      setLoading(false);
    })();
  }, [isApprovedB2B]);

  const pickup = deliveryType === 'PICKUP';
  const selectedAddress = addresses.find((a) => a.id === addressId);
  const city = selectedAddress ? selectedAddress.city : newAddr.city;

  // Ask the server for the real totals whenever the selection, courier or city changes.
  useEffect(() => {
    if (selectedIds.length === 0) return;
    const t = setTimeout(() => {
      cartService
        .checkoutSummary({
          cart_item_ids: selectedIds,
          delivery_type: deliveryType,
          shipping_method_id: pickup ? undefined : shippingId || undefined,
          city: pickup ? undefined : city.trim() || undefined,
        })
        .then(setSummary)
        .catch((err) => toast.error(getErrorMessage(err)));
    }, 300);
    return () => clearTimeout(t);
  }, [selectedIds, shippingId, city, deliveryType, pickup]);

  const payment = payments.find((p) => p.code === paymentCode);
  const total = summary?.grand_total_amount ?? 0;
  const creditShort = payment?.type === 'B2B_CREDIT' && credit && credit.remaining_credit < total;

  const addressError = useMemo(() => {
    if (pickup || addressId !== NEW_ADDRESS) return null;
    if (newAddr.recipient_name.trim().length < 3) return 'Isi nama penerima.';
    if (!isPhone(newAddr.phone)) return 'Nomor HP penerima tidak valid (contoh: 081234567890).';
    if (!newAddr.full_address.trim()) return 'Isi alamat lengkap.';
    if (!newAddr.city.trim()) return 'Isi kota / kabupaten.';
    if (saveAddress && (!newAddr.province.trim() || !newAddr.postal_code.trim())) return 'Isi provinsi dan kode pos untuk menyimpan alamat.';
    return null;
  }, [pickup, addressId, newAddr, saveAddress]);

  // Nothing selected (e.g. page opened directly): go back to the cart. Not
  // after a successful order, which navigates to the order page itself.
  if (selectedIds.length === 0 && !placedRef.current) return <Navigate to="/cart" replace />;

  const placeOrder = async () => {
    if (addressError) return toast.error(addressError);
    if (!payment) return toast.error('Pilih metode pembayaran.');
    if (!summary?.can_proceed_to_order) return toast.error(summary?.validation_errors?.[0] || 'Pesanan belum bisa diproses.');

    setSubmitting(true);
    try {
      let useAddressId = addressId !== NEW_ADDRESS ? addressId : undefined;
      if (!pickup && !useAddressId && saveAddress) {
        const saved = await addressService.createAddress({ ...newAddr, is_default: addresses.length === 0 });
        useAddressId = saved.id;
      }
      const order = await orderService.createOrder({
        payment_method: payment.code,
        delivery_type: deliveryType,
        cart_item_ids: selectedIds,
        shipping_method_id: pickup ? undefined : shippingId || undefined,
        notes: notes.trim() || undefined,
        ...(pickup
          ? {}
          : useAddressId
          ? { address_id: useAddressId }
          : {
              recipient_name: newAddr.recipient_name.trim(),
              phone: newAddr.phone.trim(),
              full_address: newAddr.full_address.trim(),
              city: newAddr.city.trim(),
              postal_code: newAddr.postal_code.trim(),
            }),
      });
      placedRef.current = true;
      toast.success('Pesanan berhasil dibuat!');
      navigate(`/orders/${order.id}`, { replace: true });
      setSelected([]);
      refreshCart().catch(() => undefined);
    } catch (err) {
      toast.error(getErrorMessage(err, 'Gagal membuat pesanan.'));
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center py-20">
        <Loader2 className="h-8 w-8 animate-spin text-primary-600" />
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-6xl space-y-6 px-4 py-8 sm:px-6 lg:px-8">
      <h1 className="flex items-center gap-3 text-2xl font-black text-slate-900 dark:text-white">
        <CheckCircle2 className="h-7 w-7 text-primary-600" /> Checkout
      </h1>

      <div className="grid gap-6 lg:grid-cols-3">
        <div className="space-y-6 lg:col-span-2">
          <Panel title="Cara terima barang">
            <div className="grid gap-3 sm:grid-cols-2">
              {[
                { value: 'DELIVERY' as const, icon: Truck, title: 'Diantar ke rumah', sub: [
                    shippings[0] ? `Ongkir ${formatCurrency(Math.min(...shippings.map((s) => s.base_rate)))}` : 'Diantar kurir toko',
                    store?.free_shipping_min ? `gratis mulai ${formatCurrency(store.free_shipping_min)}` : '',
                  ].filter(Boolean).join(', ') },
                { value: 'PICKUP' as const, icon: Store, title: 'Ambil sendiri di toko', sub: store?.address ? `Gratis ongkir · ${store.address}` : 'Gratis ongkir' },
              ].map((o) => (
                <button
                  key={o.value}
                  type="button"
                  onClick={() => setDeliveryType(o.value)}
                  className={`flex items-center gap-3 rounded-2xl border-2 p-4 text-left transition-colors ${deliveryType === o.value ? 'border-primary-600 bg-primary-50/60 dark:bg-primary-950/30' : 'border-slate-200 hover:border-slate-300 dark:border-slate-700'}`}
                >
                  <o.icon className={`h-7 w-7 flex-shrink-0 ${deliveryType === o.value ? 'text-primary-600' : 'text-slate-400'}`} />
                  <div>
                    <div className="font-bold">{o.title}</div>
                    <div className="text-xs text-slate-500">{o.sub}</div>
                  </div>
                </button>
              ))}
            </div>
            {pickup && (
              <Field label="Catatan untuk toko (opsional)" className="mt-4">
                <input className={inputClass} value={notes} onChange={(e) => setNotes(e.target.value)} placeholder="Misal: saya ambil jam 5 sore" />
              </Field>
            )}
          </Panel>

          {!pickup && (
            <Panel title={<><MapPin className="h-5 w-5 text-primary-600" /> Alamat pengiriman</>}>
              {addresses.length > 0 && (
                <div className="mb-4 space-y-2">
                  {addresses.map((a) => (
                    <label key={a.id} className={`flex cursor-pointer gap-3 rounded-2xl border p-3 ${addressId === a.id ? 'border-primary-600 bg-primary-50/50 dark:bg-primary-950/30' : 'border-slate-200 dark:border-slate-700'}`}>
                      <input type="radio" name="address" className="mt-1" checked={addressId === a.id} onChange={() => setAddressId(a.id)} />
                      <div className="text-sm">
                        <div className="font-bold">{a.title} · {a.recipient_name} <span className="font-normal text-slate-500">({a.phone})</span></div>
                        <div className="text-slate-600 dark:text-slate-400">{a.full_address}, {a.city}, {a.province} {a.postal_code}</div>
                      </div>
                    </label>
                  ))}
                  <label className={`flex cursor-pointer items-center gap-3 rounded-2xl border p-3 text-sm font-semibold ${addressId === NEW_ADDRESS ? 'border-primary-600' : 'border-slate-200 dark:border-slate-700'}`}>
                    <input type="radio" name="address" checked={addressId === NEW_ADDRESS} onChange={() => setAddressId(NEW_ADDRESS)} />
                    Kirim ke alamat lain
                  </label>
                </div>
              )}
  
              {addressId === NEW_ADDRESS && (
                <div className="grid gap-3 sm:grid-cols-2">
                  <Field label="Nama penerima">
                    <input className={inputClass} value={newAddr.recipient_name} onChange={(e) => setNewAddr({ ...newAddr, recipient_name: e.target.value })} autoComplete="name" />
                  </Field>
                  <Field label="Nomor HP penerima">
                    <input type="tel" className={inputClass} value={newAddr.phone} onChange={(e) => setNewAddr({ ...newAddr, phone: e.target.value })} placeholder="08xxxxxxxxxx" autoComplete="tel" />
                  </Field>
                  <Field label="Alamat lengkap" className="sm:col-span-2">
                    <textarea rows={2} className={inputClass} value={newAddr.full_address} onChange={(e) => setNewAddr({ ...newAddr, full_address: e.target.value })} placeholder="Nama jalan, nomor rumah, RT/RW, kelurahan, patokan" autoComplete="street-address" />
                  </Field>
                  <Field label="Kota / kabupaten">
                    <input className={inputClass} value={newAddr.city} onChange={(e) => setNewAddr({ ...newAddr, city: e.target.value })} autoComplete="address-level2" />
                  </Field>
                  <Field label="Provinsi">
                    <input className={inputClass} value={newAddr.province} onChange={(e) => setNewAddr({ ...newAddr, province: e.target.value })} autoComplete="address-level1" />
                  </Field>
                  <Field label="Kode pos">
                    <input className={inputClass} value={newAddr.postal_code} onChange={(e) => setNewAddr({ ...newAddr, postal_code: e.target.value })} autoComplete="postal-code" />
                  </Field>
                  <Field label="Simpan sebagai">
                    <input className={inputClass} value={newAddr.title} onChange={(e) => setNewAddr({ ...newAddr, title: e.target.value })} placeholder="Rumah / Toko" />
                  </Field>
                  <label className="flex items-center gap-2 text-sm sm:col-span-2">
                    <input type="checkbox" checked={saveAddress} onChange={(e) => setSaveAddress(e.target.checked)} />
                    Simpan alamat ini untuk pesanan berikutnya
                  </label>
                </div>
              )}
  
              <Field label="Catatan untuk kurir (opsional)" className="mt-4">
                <input className={inputClass} value={notes} onChange={(e) => setNotes(e.target.value)} placeholder="Misal: rumah pagar hijau, titip ke tetangga" />
              </Field>
            </Panel>
          )}

          {!pickup && shippings.length > 1 && (
            <Panel title={<><Truck className="h-5 w-5 text-primary-600" /> Pengiriman</>}>
              <div className="space-y-2">
                {shippings.map((s) => (
                  <label key={s.id} className={`flex cursor-pointer items-center gap-3 rounded-2xl border p-3 text-sm ${shippingId === s.id ? 'border-primary-600' : 'border-slate-200 dark:border-slate-700'}`}>
                    <input type="radio" checked={shippingId === s.id} onChange={() => setShippingId(s.id)} />
                    <div className="flex-1">
                      <div className="font-bold">{s.name}</div>
                      <div className="text-xs text-slate-500">{s.courier_name} · {s.estimated_days}</div>
                    </div>
                    <div className="font-semibold">{formatCurrency(s.base_rate)}</div>
                  </label>
                ))}
              </div>
            </Panel>
          )}

          <Panel title={<><CreditCard className="h-5 w-5 text-primary-600" /> Pembayaran</>}>
            <div className="space-y-2">
              {payments.map((m) => (
                <label key={m.id} className={`flex cursor-pointer gap-3 rounded-2xl border p-3 ${paymentCode === m.code ? 'border-primary-600 bg-primary-50/50 dark:bg-primary-950/30' : 'border-slate-200 dark:border-slate-700'}`}>
                  <input type="radio" name="payment" className="mt-1" checked={paymentCode === m.code} onChange={() => setPaymentCode(m.code)} />
                  <div className="text-sm">
                    <div className="font-bold">{m.type === 'COD' ? (pickup ? 'Bayar tunai saat ambil di toko' : 'Bayar tunai saat barang sampai') : m.name}</div>
                    {m.type !== 'COD' && m.description && <div className="whitespace-pre-line text-xs text-slate-500">{m.description}</div>}
                    {m.type === 'B2B_CREDIT' && credit && (
                      <div className="mt-1 text-xs text-slate-600 dark:text-slate-400">Sisa limit kasbon: <b>{formatCurrency(credit.remaining_credit)}</b></div>
                    )}
                  </div>
                </label>
              ))}
            </div>
            {payment?.type === 'MANUAL_TRANSFER' && (
              <p className="mt-3 rounded-xl bg-sky-50 p-3 text-xs text-sky-800 dark:bg-sky-950/40 dark:text-sky-300">
                Setelah pesanan dibuat, transfer sesuai total lalu unggah bukti transfer di halaman pesanan. Pesanan diproses setelah pembayaran diverifikasi toko.
              </p>
            )}
            {creditShort && (
              <p className="mt-3 rounded-xl bg-rose-50 p-3 text-xs text-rose-700 dark:bg-rose-950/40 dark:text-rose-300">
                Sisa limit kasbon tidak cukup untuk pesanan ini. Pilih metode lain atau lunasi kasbon terlebih dahulu.
              </p>
            )}
          </Panel>
        </div>

        <div className="h-fit space-y-4 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900 lg:sticky lg:top-28">
          <h2 className="text-base font-extrabold">Ringkasan</h2>
          {!summary ? (
            <Loader2 className="mx-auto h-6 w-6 animate-spin text-primary-600" />
          ) : (
            <>
              <ul className="max-h-60 space-y-2 overflow-y-auto">
                {summary.items.map((i) => (
                  <li key={i.id} className="flex items-center gap-3 text-sm">
                    <div className="flex h-10 w-10 flex-shrink-0 items-center justify-center overflow-hidden rounded-xl bg-slate-100 dark:bg-slate-800">
                      {i.image_url ? <img src={i.image_url} alt="" className="h-full w-full object-cover" /> : <Package className="h-4 w-4 text-slate-400" />}
                    </div>
                    <div className="min-w-0 flex-1">
                      <div className="truncate font-medium">{i.product_name}</div>
                      <div className="text-xs text-slate-500">{i.quantity} {i.unit_name} × {formatCurrency(i.unit_price)}</div>
                    </div>
                    <div className="font-semibold">{formatCurrency(i.subtotal)}</div>
                  </li>
                ))}
              </ul>
              <div className="space-y-1.5 border-t border-slate-100 pt-3 text-sm dark:border-slate-800">
                <div className="flex justify-between text-slate-600 dark:text-slate-400"><span>Subtotal ({summary.total_items} barang)</span><span>{formatCurrency(summary.subtotal_amount)}</span></div>
                <div className="flex justify-between text-slate-600 dark:text-slate-400"><span>Ongkir</span><span>{pickup ? 'Gratis (ambil di toko)' : summary.estimated_shipping_fee === 0 ? 'Gratis' : formatCurrency(summary.estimated_shipping_fee)}</span></div>
                {!pickup && summary.free_shipping_min > 0 && (
                  <div className={`rounded-lg px-2 py-1 text-xs font-semibold ${summary.subtotal_amount >= summary.free_shipping_min ? 'bg-primary-50 text-primary-700 dark:bg-primary-950/40 dark:text-primary-300' : 'bg-amber-50 text-amber-800 dark:bg-amber-950/40 dark:text-amber-300'}`}>
                    {summary.subtotal_amount >= summary.free_shipping_min
                      ? `Gratis ongkir untuk belanja di atas ${formatCurrency(summary.free_shipping_min)}`
                      : `Tambah ${formatCurrency(summary.free_shipping_min - summary.subtotal_amount)} lagi untuk gratis ongkir`}
                  </div>
                )}
                <div className="flex justify-between pt-2 text-lg font-black"><span>Total</span><span className="text-primary-600">{formatCurrency(total)}</span></div>
              </div>
              {summary.validation_errors && summary.validation_errors.length > 0 && (
                <div className="space-y-1 rounded-xl bg-rose-50 p-3 text-xs text-rose-700 dark:bg-rose-950/40 dark:text-rose-300">
                  {summary.validation_errors.map((e) => (
                    <div key={e} className="flex gap-1"><AlertTriangle className="h-3.5 w-3.5 flex-shrink-0" /> {e}</div>
                  ))}
                  <Link to="/cart" className="font-bold underline">Ubah keranjang</Link>
                </div>
              )}
            </>
          )}
          <button
            onClick={placeOrder}
            disabled={submitting || !summary?.can_proceed_to_order || !payment || !!creditShort}
            className="flex w-full items-center justify-center gap-2 rounded-2xl bg-primary-600 py-3.5 text-sm font-bold text-white shadow-lg shadow-primary-600/30 hover:bg-primary-700 disabled:opacity-50"
          >
            {submitting && <Loader2 className="h-4 w-4 animate-spin" />}
            {submitting ? 'Memproses...' : 'Buat Pesanan'}
          </button>
        </div>
      </div>
    </div>
  );
};
