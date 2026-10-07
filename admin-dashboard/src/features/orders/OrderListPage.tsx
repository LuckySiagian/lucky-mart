import React, { useState } from 'react';
import { PageHeader } from '../../components/common/PageHeader';
import { formatCurrency } from '../../lib/utils';
import { ShoppingCart, Eye, Ban, Search, CheckCircle2, Truck, Clock, X, FileText, Plus, Trash2 } from 'lucide-react';
import { toast } from 'sonner';
import { orderService } from '../../services/apiServices';

interface OrderItem {
  name: string;
  unit: string;
  qty: number;
  price: number;
}

interface OrderRecord {
  id: string;
  orderNumber: string;
  customerName: string;
  customerType: 'B2C_CUSTOMER' | 'B2B_MITRA';
  orderDate: string;
  paymentMethod: string;
  grandTotal: number;
  status: 'PENDING_PAYMENT' | 'PAID' | 'PROCESSING' | 'SHIPPED' | 'COMPLETED' | 'CANCELLED';
  items: OrderItem[];
  shippingAddress: string;
}

export const OrderListPage: React.FC = () => {
  const [orders, setOrders] = useState<OrderRecord[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState('');

  const fetchOrders = async () => {
    setLoading(true);
    try {
      const res = await orderService.getOrders(1, 100);
      const items = Array.isArray(res.data)
        ? res.data
        : (res.data?.items || (Array.isArray(res) ? res : []));
      
      const mapped = items.map((o: any) => ({
        id: o.id,
        orderNumber: o.order_number || `ORD-${o.id.slice(0, 8)}`,
        customerName: o.user?.full_name || o.customer_name || 'Pelanggan LuckyMart',
        customerType: (o.user?.role === 'B2B_CUSTOMER' ? 'B2B_MITRA' : 'B2C_CUSTOMER') as any,
        orderDate: o.created_at ? new Date(o.created_at).toLocaleString('id-ID') : 'Terbaru',
        paymentMethod: o.payment_method || 'Transfer Bank',
        grandTotal: o.total_amount || 0,
        status: o.status || 'PENDING_PAYMENT',
        items: o.order_items?.map((item: any) => ({
          name: item.product_name || 'Produk Sembako',
          unit: item.unit_name || 'Pcs',
          qty: item.quantity || 1,
          price: item.price || 0,
        })) || [],
        shippingAddress: o.address?.address_line || 'Alamat Penerima',
      }));
      setOrders(mapped);
    } catch (e) {
      setOrders([]);
    } finally {
      setLoading(false);
    }
  };

  React.useEffect(() => {
    fetchOrders();
  }, []);

  // Detail Modal State
  const [selectedOrder, setSelectedOrder] = useState<OrderRecord | null>(null);

  // Manual Create Order Modal State
  const [isAddOpen, setIsAddOpen] = useState(false);
  const [custName, setCustName] = useState('');
  const [custType, setCustType] = useState<'B2C_CUSTOMER' | 'B2B_MITRA'>('B2C_CUSTOMER');
  const [payMethod, setPayMethod] = useState('Tunai / Cash Toko');
  const [itemSelect, setItemSelect] = useState('Beras Pandan Wangi Super 10kg');
  const [itemQty, setItemQty] = useState('1');
  const [itemPrice, setItemPrice] = useState('145000');
  const [address, setAddress] = useState('Alamat Toko Utama');

  const handleUpdateStatus = async (id: string, newStatus: OrderRecord['status']) => {
    try {
      await orderService.updateOrderStatus(id, newStatus);
      setOrders((prev) =>
        prev.map((o) => (o.id === id ? { ...o, status: newStatus } : o))
      );
      toast.success(`Status Order "${id}" diperbarui menjadi ${newStatus}!`);
      if (selectedOrder && selectedOrder.id === id) {
        setSelectedOrder((prev) => (prev ? { ...prev, status: newStatus } : null));
      }
    } catch (e: any) {
      toast.error('Gagal memperbarui status order ke server.');
    }
  };

  const handleDeleteOrder = (id: string, num: string) => {
    if (window.confirm(`Hapus permanen data pesanan "${num}" dari database?`)) {
      setOrders((prev) => prev.filter((o) => o.id !== id));
      toast.success(`Pesanan "${num}" berhasil dihapus!`);
      if (selectedOrder?.id === id) setSelectedOrder(null);
    }
  };

  const handleCreateManualOrder = (e: React.FormEvent) => {
    e.preventDefault();
    if (!custName.trim() || !itemPrice) {
      toast.error('Nama pelanggan dan total harga wajib diisi.');
      return;
    }

    const qty = parseInt(itemQty, 10) || 1;
    const price = parseFloat(itemPrice) || 50000;
    const orderNum = `ORD-${new Date().toISOString().slice(0, 10).replace(/-/g, '')}-${Math.floor(1000 + Math.random() * 9000)}`;

    const newOrder: OrderRecord = {
      id: `ORD-${Date.now()}`,
      orderNumber: orderNum,
      customerName: custName.trim(),
      customerType: custType,
      orderDate: new Date().toISOString().replace('T', ' ').slice(0, 16),
      paymentMethod: payMethod,
      grandTotal: qty * price,
      status: 'PAID',
      items: [{ name: itemSelect, unit: 'Unit', qty: qty, price: price }],
      shippingAddress: address.trim() || 'Diambil di Gudang',
    };

    setOrders((prev) => [newOrder, ...prev]);
    toast.success(`Pesanan manual "${orderNum}" berhasil dibuat!`);
    setIsAddOpen(false);
    setCustName('');
  };

  const filtered = orders.filter(
    (o) =>
      o.orderNumber.toLowerCase().includes(search.toLowerCase()) ||
      o.customerName.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className="space-y-6">
      <PageHeader
        title="Manajemen Pesanan & Fulfillment (Orders)"
        subtitle="Pantau siklus pemrosesan pesanan dari Pembeli B2C & Mitra Warung B2B dari pembayaran hingga pengiriman barang."
      >
        <button
          onClick={() => setIsAddOpen(true)}
          className="inline-flex items-center space-x-2 rounded-xl bg-brand-600 px-4 py-2.5 text-sm font-semibold text-white hover:bg-brand-700 transition-colors shadow-lg shadow-brand-600/20"
        >
          <Plus className="h-4 w-4" />
          <span>Buat Pesanan Manual</span>
        </button>
      </PageHeader>

      {/* Filter & Search Bar */}
      <div className="flex flex-col sm:flex-row items-center justify-between gap-4 bg-white dark:bg-slate-900 p-4 rounded-2xl border border-slate-200 dark:border-slate-800">
        <div className="relative w-full sm:w-80">
          <Search className="absolute left-3.5 top-3 h-4 w-4 text-slate-400" />
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Cari no. order atau nama pembeli..."
            className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 pl-10 pr-4 py-2 text-sm text-slate-900 dark:text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-brand-500"
          />
        </div>
        <div className="text-xs text-slate-500 font-semibold">
          Total Transaksi: <span className="text-brand-600 font-bold">{orders.length} Pesanan</span>
        </div>
      </div>

      {/* Orders Table */}
      <div className="overflow-x-auto rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 shadow-sm">
        <table className="w-full text-left text-sm text-slate-600 dark:text-slate-400">
          <thead className="bg-slate-50 dark:bg-slate-800/50 text-xs uppercase text-slate-500 font-bold border-b border-slate-200 dark:border-slate-800">
            <tr>
              <th className="px-6 py-4">Nomor Order & Waktu</th>
              <th className="px-6 py-4">Pelanggan & Tipe</th>
              <th className="px-6 py-4">Metode Bayar</th>
              <th className="px-6 py-4">Total Belanja</th>
              <th className="px-6 py-4">Status Pesanan</th>
              <th className="px-6 py-4 text-right">Aksi Kelola (CRUD)</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-200 dark:divide-slate-800">
            {filtered.map((o) => (
              <tr key={o.id} className="hover:bg-slate-50/80 dark:hover:bg-slate-800/30">
                <td className="px-6 py-4">
                  <div className="font-mono font-bold text-slate-900 dark:text-white">{o.orderNumber}</div>
                  <div className="text-xs text-slate-400">{o.orderDate}</div>
                </td>
                <td className="px-6 py-4">
                  <div className="font-bold text-slate-900 dark:text-white">{o.customerName}</div>
                  <span className={`inline-block text-[10px] font-bold px-2 py-0.5 rounded-full ${o.customerType === 'B2B_MITRA' ? 'bg-emerald-100 text-emerald-800' : 'bg-slate-100 text-slate-700'}`}>
                    {o.customerType === 'B2B_MITRA' ? 'B2B MITRA WARUNG' : 'B2C ECERAN'}
                  </span>
                </td>
                <td className="px-6 py-4 text-xs font-semibold text-slate-700 dark:text-slate-300">{o.paymentMethod}</td>
                <td className="px-6 py-4 font-black text-slate-900 dark:text-white">{formatCurrency(o.grandTotal)}</td>
                <td className="px-6 py-4">
                  {o.status === 'PENDING_PAYMENT' && (
                    <span className="inline-flex items-center rounded-full bg-amber-100 text-amber-800 px-2.5 py-0.5 text-xs font-bold">
                      MENUNGGU BAYAR
                    </span>
                  )}
                  {o.status === 'PAID' && (
                    <span className="inline-flex items-center rounded-full bg-blue-100 text-blue-800 px-2.5 py-0.5 text-xs font-bold">
                      LUNAS (PAID)
                    </span>
                  )}
                  {o.status === 'PROCESSING' && (
                    <span className="inline-flex items-center rounded-full bg-purple-100 text-purple-800 px-2.5 py-0.5 text-xs font-bold">
                      DIPROSES GUDANG
                    </span>
                  )}
                  {o.status === 'SHIPPED' && (
                    <span className="inline-flex items-center rounded-full bg-indigo-100 text-indigo-800 px-2.5 py-0.5 text-xs font-bold">
                      DALAM PENGIRIMAN
                    </span>
                  )}
                  {o.status === 'COMPLETED' && (
                    <span className="inline-flex items-center rounded-full bg-emerald-100 text-emerald-800 px-2.5 py-0.5 text-xs font-bold">
                      SELESAI (COMPLETED)
                    </span>
                  )}
                  {o.status === 'CANCELLED' && (
                    <span className="inline-flex items-center rounded-full bg-rose-100 text-rose-800 px-2.5 py-0.5 text-xs font-bold">
                      DIBATALKAN
                    </span>
                  )}
                </td>
                <td className="px-6 py-4 text-right space-x-2">
                  <button
                    onClick={() => setSelectedOrder(o)}
                    className="inline-flex items-center space-x-1 rounded-lg bg-slate-100 dark:bg-slate-800 px-2.5 py-1 text-xs font-bold text-slate-700 dark:text-slate-300 hover:bg-brand-600 hover:text-white transition-colors"
                  >
                    <Eye className="h-3.5 w-3.5" />
                    <span>Detail</span>
                  </button>
                  <button
                    onClick={() => handleDeleteOrder(o.id, o.orderNumber)}
                    className="inline-flex items-center space-x-1 rounded-lg bg-rose-100 text-rose-700 dark:bg-rose-950/50 dark:text-rose-400 px-2.5 py-1 text-xs font-bold hover:bg-rose-200"
                  >
                    <Trash2 className="h-3.5 w-3.5" />
                    <span>Hapus</span>
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Modal Detail Order */}
      {selectedOrder && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 backdrop-blur-sm p-4">
          <div className="relative w-full max-w-lg rounded-3xl bg-white dark:bg-slate-900 p-6 shadow-2xl border border-slate-200 dark:border-slate-800 space-y-6">
            <div className="flex items-center justify-between border-b border-slate-200 dark:border-slate-800 pb-4">
              <div className="flex items-center space-x-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-brand-100 text-brand-600">
                  <FileText className="h-5 w-5" />
                </div>
                <div>
                  <h3 className="text-lg font-bold text-slate-900 dark:text-white">Rincian Invoice Order</h3>
                  <p className="text-xs text-slate-500">{selectedOrder.orderNumber}</p>
                </div>
              </div>
              <button onClick={() => setSelectedOrder(null)} className="rounded-xl p-2 text-slate-400 hover:bg-slate-100">
                <X className="h-5 w-5" />
              </button>
            </div>

            <div className="space-y-4">
              <div className="rounded-2xl bg-slate-100 dark:bg-slate-800 p-4 space-y-2 text-xs">
                <div className="flex justify-between">
                  <span className="text-slate-400">Pembeli:</span>
                  <span className="font-bold text-slate-900 dark:text-white">{selectedOrder.customerName}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-slate-400">Metode Bayar:</span>
                  <span className="font-bold text-slate-900 dark:text-white">{selectedOrder.paymentMethod}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-slate-400">Alamat Pengiriman:</span>
                  <span className="font-bold text-slate-900 dark:text-white max-w-xs text-right">{selectedOrder.shippingAddress}</span>
                </div>
              </div>

              {/* Items Table */}
              <div>
                <h4 className="text-xs font-bold text-slate-500 mb-2">Item Sembako Dipesan:</h4>
                <div className="space-y-2">
                  {selectedOrder.items.map((item, idx) => (
                    <div key={idx} className="flex justify-between items-center bg-slate-50 dark:bg-slate-800/50 p-3 rounded-xl text-xs">
                      <div>
                        <div className="font-bold text-slate-900 dark:text-white">{item.name}</div>
                        <div className="text-slate-400">{item.qty} {item.unit} x {formatCurrency(item.price)}</div>
                      </div>
                      <div className="font-bold text-brand-600 dark:text-brand-400">
                        {formatCurrency(item.qty * item.price)}
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              <div className="flex justify-between items-center pt-3 border-t border-slate-200 dark:border-slate-800">
                <span className="text-sm font-bold text-slate-900 dark:text-white">Total Tagihan:</span>
                <span className="text-lg font-black text-brand-600 dark:text-brand-400">{formatCurrency(selectedOrder.grandTotal)}</span>
              </div>
            </div>

            {/* Workflow Action Buttons */}
            <div className="space-y-2 pt-2 border-t border-slate-200 dark:border-slate-800">
              <label className="block text-xs font-bold text-slate-500 mb-1">Perbarui Status Fulfillment:</label>
              <div className="grid grid-cols-2 gap-2">
                <button
                  onClick={() => handleUpdateStatus(selectedOrder.id, 'PROCESSING')}
                  className="rounded-xl bg-purple-600 px-3 py-2 text-xs font-bold text-white hover:bg-purple-700"
                >
                  Set diproses Gudang
                </button>
                <button
                  onClick={() => handleUpdateStatus(selectedOrder.id, 'SHIPPED')}
                  className="rounded-xl bg-indigo-600 px-3 py-2 text-xs font-bold text-white hover:bg-indigo-700"
                >
                  Set Kirim Resi
                </button>
                <button
                  onClick={() => handleUpdateStatus(selectedOrder.id, 'COMPLETED')}
                  className="col-span-2 rounded-xl bg-emerald-600 px-3 py-2 text-xs font-bold text-white hover:bg-emerald-700"
                >
                  Tandai Selesai (Completed)
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Modal Buat Pesanan Manual Admin */}
      {isAddOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 backdrop-blur-sm p-4">
          <div className="relative w-full max-w-md rounded-3xl bg-white dark:bg-slate-900 p-6 shadow-2xl border border-slate-200 dark:border-slate-800 space-y-6">
            <div className="flex items-center justify-between border-b border-slate-200 dark:border-slate-800 pb-4">
              <div className="flex items-center space-x-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-brand-100 text-brand-600">
                  <Plus className="h-5 w-5" />
                </div>
                <div>
                  <h3 className="text-lg font-bold text-slate-900 dark:text-white">Buat Pesanan Manual Admin</h3>
                  <p className="text-xs text-slate-500">Pencatatan order pembelian offline atau pesanan via telpon.</p>
                </div>
              </div>
              <button onClick={() => setIsAddOpen(false)} className="rounded-xl p-2 text-slate-400 hover:bg-slate-100">
                <X className="h-5 w-5" />
              </button>
            </div>

            <form onSubmit={handleCreateManualOrder} className="space-y-4">
              <div>
                <label className="block text-xs font-bold text-slate-500 mb-1">Nama Pembeli *</label>
                <input
                  type="text"
                  required
                  placeholder="Contoh: Pembeli Walk-in Toko"
                  value={custName}
                  onChange={(e) => setCustName(e.target.value)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white font-bold"
                />
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs font-bold text-slate-500 mb-1">Tipe Pelanggan</label>
                  <select
                    value={custType}
                    onChange={(e) => setCustType(e.target.value as any)}
                    className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-xs text-slate-900 dark:text-white font-bold"
                  >
                    <option value="B2C_CUSTOMER">B2C Eceran</option>
                    <option value="B2B_MITRA">B2B Mitra Warung</option>
                  </select>
                </div>

                <div>
                  <label className="block text-xs font-bold text-slate-500 mb-1">Metode Bayar</label>
                  <select
                    value={payMethod}
                    onChange={(e) => setPayMethod(e.target.value)}
                    className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-xs text-slate-900 dark:text-white font-bold"
                  >
                    <option value="Tunai / Cash Toko">Tunai / Cash Toko</option>
                    <option value="B2B Kasbon Tempo 30 Hari">B2B Kasbon Tempo</option>
                  </select>
                </div>
              </div>

              <div>
                <label className="block text-xs font-bold text-slate-500 mb-1">Pilih Produk Sembako</label>
                <select
                  value={itemSelect}
                  onChange={(e) => setItemSelect(e.target.value)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-xs text-slate-900 dark:text-white font-bold"
                >
                  <option value="Beras Pandan Wangi Super 10kg">Beras Pandan Wangi Super 10kg</option>
                  <option value="Minyak Goreng Bimoli Refill 2L">Minyak Goreng Bimoli Refill 2L</option>
                  <option value="Telur Ayam Negeri Super (1kg)">Telur Ayam Negeri Super (1kg)</option>
                  <option value="Indomie Goreng Spesial (Karton / 40 Pcs)">Indomie Goreng Spesial (Karton)</option>
                </select>
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs font-bold text-slate-500 mb-1">Jumlah (Qty)</label>
                  <input
                    type="number"
                    min="1"
                    value={itemQty}
                    onChange={(e) => setItemQty(e.target.value)}
                    className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white font-bold"
                  />
                </div>

                <div>
                  <label className="block text-xs font-bold text-slate-500 mb-1">Harga Satuan (Rp)</label>
                  <input
                    type="number"
                    value={itemPrice}
                    onChange={(e) => setItemPrice(e.target.value)}
                    className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white font-bold"
                  />
                </div>
              </div>

              <div className="flex items-center justify-end space-x-3 pt-4 border-t border-slate-200 dark:border-slate-800">
                <button
                  type="button"
                  onClick={() => setIsAddOpen(false)}
                  className="rounded-xl px-4 py-2.5 text-xs font-bold text-slate-600 hover:bg-slate-100"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  className="rounded-xl bg-brand-600 px-5 py-2.5 text-xs font-bold text-white hover:bg-brand-700 shadow-lg shadow-brand-600/20"
                >
                  Simpan & Buat Pesanan
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
