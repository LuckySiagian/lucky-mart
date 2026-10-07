import React, { useEffect, useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { cartService, orderService, addressService, paymentService } from '../../services/apiServices';
import { useCartStore } from '../../store/useCartStore';
import { formatCurrency } from '../../lib/utils';
import { MapPin, CreditCard, ShoppingBag, CheckCircle2, ArrowRight, Package, Truck, ShieldCheck } from 'lucide-react';
import { toast } from 'sonner';

export const CheckoutPage: React.FC = () => {
  const navigate = useNavigate();
  const { items: storeItems, selectedItemIdsForCheckout, removeMultipleItems, clearCart, setSelectedForCheckout } = useCartStore();
  const [summary, setSummary] = useState<any>(null);
  const [addresses, setAddresses] = useState<any[]>([]);
  const [paymentMethods, setPaymentMethods] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);

  const [selectedPaymentMethod, setSelectedPaymentMethod] = useState<string>('COD');
  const [recipientName, setRecipientName] = useState('Budi Santoso');
  const [phone, setPhone] = useState('081234567890');
  const [fullAddress, setFullAddress] = useState('Jl. Raya Cengkareng No. 45');
  const [city, setCity] = useState('Jakarta Barat');
  const [postalCode, setPostalCode] = useState('11730');

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      try {
        const [sumRes, addrRes, payRes] = await Promise.all([
          cartService.getCheckoutSummary(),
          addressService.getAddresses(),
          paymentService.getPaymentMethods(),
        ]);

        setSummary(sumRes.data || null);
        setAddresses(addrRes.data || []);
        setPaymentMethods(payRes.data || []);

        if (addrRes.data && addrRes.data.length > 0) {
          const def = addrRes.data.find((a: any) => a.is_default) || addrRes.data[0];
          setRecipientName(def.recipient_name);
          setPhone(def.phone);
          setFullAddress(def.full_address);
          setCity(def.city);
          setPostalCode(def.postal_code);
        }
      } catch (e) {
        // Fallback display
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  const handlePlaceOrder = async () => {
    setSubmitting(true);
    try {
      const orderPayload = {
        payment_method: selectedPaymentMethod,
        recipient_name: recipientName,
        phone: phone,
        full_address: fullAddress,
        city: city,
        postal_code: postalCode,
      };

      await orderService.createOrder(orderPayload);
      if (selectedItemIdsForCheckout && selectedItemIdsForCheckout.length > 0) {
        removeMultipleItems(selectedItemIdsForCheckout);
        setSelectedForCheckout([]);
      } else {
        clearCart();
      }
      toast.success('Pesanan berhasil dibuat! Stok dikunci secara otomatis.');
      navigate('/orders');
    } catch (err: any) {
      toast.error(err.response?.data?.message || 'Gagal membuat pesanan. Silakan periksa kembali.');
    } finally {
      setSubmitting(false);
    }
  };

  // Extract ordered items list from summary or Zustand store, filtered by user selection
  const rawItems = (summary?.items && summary.items.length > 0) ? summary.items : storeItems;
  const displayItems = (selectedItemIdsForCheckout && selectedItemIdsForCheckout.length > 0)
    ? rawItems.filter((i: any) => selectedItemIdsForCheckout.includes(i.id))
    : rawItems;
  const subtotalAmount = displayItems.reduce((acc: number, i: any) => acc + (i.subtotal || (i.unit_price * (i.quantity || 1))), 0);
  const shippingFee = summary?.estimated_shipping_fee || 15000;
  const grandTotalAmount = subtotalAmount + shippingFee;
  const totalItemQuantity = displayItems.reduce((acc: number, i: any) => acc + (i.quantity || 1), 0);

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-8">
      <h1 className="text-2xl font-black text-slate-900 dark:text-white flex items-center space-x-3">
        <CheckCircle2 className="h-7 w-7 text-primary-600" />
        <span>Konfirmasi Checkout Pesanan Sembako</span>
      </h1>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Left Form Column */}
        <div className="lg:col-span-2 space-y-6">
          {/* Section 1: Daftar Barang yang Dipesan */}
          <div className="rounded-3xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 p-6 space-y-4 shadow-sm">
            <div className="flex items-center justify-between border-b border-slate-100 dark:border-slate-800 pb-3">
              <h3 className="text-base font-bold text-slate-900 dark:text-white flex items-center space-x-2">
                <ShoppingBag className="h-5 w-5 text-primary-600" />
                <span>Rincian Barang yang Dipesan ({displayItems.length} Produk)</span>
              </h3>
              <Link to="/cart" className="text-xs font-bold text-primary-600 hover:underline">
                Ubah di Keranjang
              </Link>
            </div>

            {displayItems.length === 0 ? (
              <div className="text-center py-6 text-slate-400 text-sm">
                Belum ada produk terpilih untuk di-checkout.
              </div>
            ) : (
              <div className="divide-y divide-slate-100 dark:divide-slate-800 space-y-3">
                {displayItems.map((item: any, idx: number) => (
                  <div key={item.id || idx} className="pt-3 first:pt-0 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
                    <div className="flex items-center space-x-4">
                      <div className="h-14 w-14 rounded-2xl bg-slate-100 dark:bg-slate-800 flex items-center justify-center font-bold text-slate-400 flex-shrink-0 text-xs">
                        {item.image_url ? (
                          <img src={item.image_url} alt={item.product_name} className="h-14 w-14 object-cover rounded-2xl" />
                        ) : (
                          <Package className="h-6 w-6 text-slate-400" />
                        )}
                      </div>
                      <div>
                        <h4 className="font-bold text-sm text-slate-900 dark:text-white">{item.product_name}</h4>
                        <div className="flex items-center space-x-3 text-xs text-slate-500 mt-0.5">
                          <span>Satuan: <b className="text-slate-700 dark:text-slate-300">{item.unit_name || 'Pcs'}</b></span>
                          <span>•</span>
                          <span>Harga Satuan: <b className="text-slate-700 dark:text-slate-300">{formatCurrency(item.unit_price)}</b></span>
                        </div>
                      </div>
                    </div>

                    <div className="flex items-center justify-between sm:justify-end w-full sm:w-auto space-x-6 border-t sm:border-t-0 border-slate-100 dark:border-slate-800 pt-2 sm:pt-0">
                      <div className="inline-flex items-center px-3 py-1 rounded-xl bg-slate-100 dark:bg-slate-800 text-xs font-extrabold text-slate-800 dark:text-slate-200">
                        {item.quantity} {item.unit_name || 'Pcs'}
                      </div>
                      <div className="text-right">
                        <span className="text-[10px] text-slate-400 block">Subtotal Item</span>
                        <span className="font-black text-sm text-primary-600">{formatCurrency(item.subtotal || (item.unit_price * item.quantity))}</span>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Section 2: Alamat Pengiriman */}
          <div className="rounded-3xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 p-6 space-y-4 shadow-sm">
            <h3 className="text-base font-bold text-slate-900 dark:text-white flex items-center space-x-2">
              <MapPin className="h-5 w-5 text-primary-600" />
              <span>Alamat Pengiriman Sembako</span>
            </h3>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label className="block text-xs font-semibold text-slate-500 mb-1">Nama Penerima</label>
                <input
                  type="text"
                  value={recipientName}
                  onChange={(e) => setRecipientName(e.target.value)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 border-0 px-4 py-2.5 text-sm text-slate-900 dark:text-white font-medium focus:ring-2 focus:ring-primary-500"
                />
              </div>
              <div>
                <label className="block text-xs font-semibold text-slate-500 mb-1">Nomor Telepon</label>
                <input
                  type="text"
                  value={phone}
                  onChange={(e) => setPhone(e.target.value)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 border-0 px-4 py-2.5 text-sm text-slate-900 dark:text-white font-medium focus:ring-2 focus:ring-primary-500"
                />
              </div>
            </div>

            <div>
              <label className="block text-xs font-semibold text-slate-500 mb-1">Alamat Lengkap</label>
              <textarea
                value={fullAddress}
                onChange={(e) => setFullAddress(e.target.value)}
                rows={2}
                className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 border-0 px-4 py-2.5 text-sm text-slate-900 dark:text-white font-medium focus:ring-2 focus:ring-primary-500"
              />
            </div>
          </div>

          {/* Section 3: Metode Pembayaran */}
          <div className="rounded-3xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 p-6 space-y-4 shadow-sm">
            <h3 className="text-base font-bold text-slate-900 dark:text-white flex items-center space-x-2">
              <CreditCard className="h-5 w-5 text-primary-600" />
              <span>Pilih Metode Pembayaran</span>
            </h3>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
              {[
                { code: 'B2B_CREDIT_TEMPO', title: 'B2B Kasbon Tempo 30 Hari', desc: 'Khusus Mitra Toko Kelontong' },
                { code: 'COD', title: 'Bayar di Tempat (COD)', desc: 'Bayar tunai ke kurir saat barang sampai' },
              ].map((m) => (
                <div
                  key={m.code}
                  onClick={() => setSelectedPaymentMethod(m.code)}
                  className={`cursor-pointer p-4 rounded-2xl border transition-all ${
                    selectedPaymentMethod === m.code
                      ? 'border-primary-600 bg-primary-50/50 dark:bg-primary-950/40 shadow-sm'
                      : 'border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-800/40'
                  }`}
                >
                  <div className="font-bold text-sm text-slate-900 dark:text-white">{m.title}</div>
                  <div className="text-xs text-slate-500 mt-1">{m.desc}</div>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Right Sticky Checkout Summary */}
        <div className="h-fit sticky top-24 rounded-3xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 p-6 space-y-6 shadow-sm">
          <h3 className="text-base font-extrabold text-slate-900 dark:text-white border-b border-slate-100 dark:border-slate-800 pb-3">
            Rincian Pembayaran
          </h3>

          <div className="space-y-3 text-sm">
            <div className="flex justify-between text-slate-600 dark:text-slate-400">
              <span>Total Barang Dipesan</span>
              <span className="font-semibold text-slate-900 dark:text-white">{totalItemQuantity} Item</span>
            </div>
            <div className="flex justify-between text-slate-600 dark:text-slate-400">
              <span>Subtotal Produk</span>
              <span className="font-bold text-slate-900 dark:text-white">
                {formatCurrency(subtotalAmount)}
              </span>
            </div>
            <div className="flex justify-between text-slate-600 dark:text-slate-400">
              <span>Ongkos Kirim Estimasi</span>
              <span className="font-bold text-slate-900 dark:text-white">
                {formatCurrency(shippingFee)}
              </span>
            </div>
          </div>

          <div className="pt-3 border-t border-slate-100 dark:border-slate-800 flex justify-between items-center">
            <span className="font-bold text-slate-900 dark:text-white">Grand Total</span>
            <span className="text-xl font-black text-primary-600">
              {formatCurrency(grandTotalAmount)}
            </span>
          </div>

          <button
            onClick={handlePlaceOrder}
            disabled={submitting || displayItems.length === 0}
            className="flex w-full items-center justify-center space-x-2 rounded-2xl bg-primary-600 py-3.5 text-sm font-bold text-white hover:bg-primary-700 transition-colors shadow-lg shadow-primary-600/30 disabled:opacity-50"
          >
            <span>{submitting ? 'Membuat Pesanan (ACID Tx)...' : 'Buat Pesanan Sekarang'}</span>
            <ArrowRight className="h-4 w-4" />
          </button>
        </div>
      </div>
    </div>
  );
};
