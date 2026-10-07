import React, { useEffect, useState } from 'react';
import { orderService, paymentService } from '../../services/apiServices';
import { formatCurrency, formatDate } from '../../lib/utils';
import { ShoppingBag, Eye, Ban, Upload, CheckCircle2 } from 'lucide-react';
import { toast } from 'sonner';

export const OrderListPage: React.FC = () => {
  const [orders, setOrders] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchOrders = async () => {
    setLoading(true);
    try {
      const res = await orderService.getOrders(1, 20);
      const items = Array.isArray(res.data)
        ? res.data
        : (res.data?.items || (Array.isArray(res) ? res : []));
      setOrders(Array.isArray(items) ? items : []);
    } catch (e) {
      setOrders([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchOrders();
  }, []);

  const handleUploadProof = async (paymentId: string) => {
    const url = window.prompt('Masukkan URL / Link Foto Bukti Transfer:');
    if (url) {
      try {
        await paymentService.uploadProof({
          payment_id: paymentId,
          proof_image_url: url,
        });
        toast.success('Bukti transfer berhasil diunggah. Menunggu konfirmasi admin finance.');
        fetchOrders();
      } catch (err: any) {
        toast.error('Gagal mengunggah bukti transfer.');
      }
    }
  };

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-6">
      <h1 className="text-2xl font-black text-slate-900 dark:text-white flex items-center space-x-3">
        <ShoppingBag className="h-7 w-7 text-primary-600" />
        <span>Pesanan Saya</span>
      </h1>

      <div className="space-y-4">
        {Array.isArray(orders) && orders.map((o) => (
          <div key={o.id} className="rounded-3xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 p-6 space-y-4 shadow-sm">
            <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center border-b border-slate-100 dark:border-slate-800 pb-3 gap-2">
              <div>
                <span className="font-mono font-bold text-sm text-slate-900 dark:text-white">{o.order_number || `ORD-${o.id?.slice(0, 8)}`}</span>
                <span className="text-xs text-slate-400 ml-3">{formatDate(o.order_date || o.created_at)}</span>
              </div>
              <span className="inline-flex items-center rounded-full bg-amber-100 text-amber-800 px-3 py-0.5 text-xs font-bold">
                {o.status}
              </span>
            </div>

            <div className="flex justify-between items-center text-sm">
              <div>
                <span className="text-xs text-slate-400">Total Pembayaran</span>
                <div className="text-lg font-black text-primary-600">{formatCurrency(o.grand_total || 0)}</div>
              </div>

              <div className="flex space-x-2">
                <button
                  onClick={() => handleUploadProof(o.id)}
                  className="inline-flex items-center space-x-1.5 rounded-xl bg-primary-600 px-4 py-2 text-xs font-bold text-white hover:bg-primary-700 transition-colors shadow-md shadow-primary-600/20"
                >
                  <Upload className="h-3.5 w-3.5" />
                  <span>Upload Bukti Bayar</span>
                </button>
              </div>
            </div>
          </div>
        ))}

        {(!Array.isArray(orders) || orders.length === 0) && !loading && (
          <div className="text-center py-12 text-slate-400 text-sm">Belum ada pesanan aktif.</div>
        )}
      </div>
    </div>
  );
};
