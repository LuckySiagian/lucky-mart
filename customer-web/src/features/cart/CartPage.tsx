import React, { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { toast } from 'sonner';
import { AlertTriangle, ArrowRight, CheckSquare, Minus, Package, Plus, ShoppingBag, Square, Trash2 } from 'lucide-react';
import { cartService } from '../../services/apiServices';
import { useCartStore } from '../../store/useCartStore';
import { EmptyState } from '../../components/common/EmptyState';
import { ErrorState } from '../../components/common/ErrorState';
import { formatCurrency, getErrorMessage } from '../../lib/utils';

export const CartPage: React.FC = () => {
  const navigate = useNavigate();
  const { items, loaded, selectedIds, setSelected, apply, refresh } = useCartStore();
  const [error, setError] = useState<string | null>(null);
  const [busyId, setBusyId] = useState<string | null>(null);

  useEffect(() => {
    refresh()
      .then(() => {
        // Pre-select everything that can be ordered when nothing is chosen yet.
        const s = useCartStore.getState();
        if (s.selectedIds.length === 0) {
          s.setSelected(s.items.filter((i) => i.is_product_active && i.is_stock_sufficient).map((i) => i.id));
        }
      })
      .catch((err) => setError(getErrorMessage(err, 'Gagal memuat keranjang.')));
  }, [refresh]);

  const run = async (id: string, fn: () => ReturnType<typeof cartService.getCart>) => {
    setBusyId(id);
    try {
      apply(await fn());
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setBusyId(null);
    }
  };

  const toggle = (id: string) => setSelected(selectedIds.includes(id) ? selectedIds.filter((x) => x !== id) : [...selectedIds, id]);
  const selectable = items.filter((i) => i.is_product_active);
  const allSelected = selectable.length > 0 && selectable.every((i) => selectedIds.includes(i.id));

  const selectedItems = items.filter((i) => selectedIds.includes(i.id));
  const selectedTotal = selectedItems.reduce((s, i) => s + i.subtotal, 0);
  const selectedQty = selectedItems.reduce((s, i) => s + i.quantity, 0);
  const blocked = selectedItems.some((i) => !i.is_product_active || !i.is_stock_sufficient);

  const removeSelected = async () => {
    if (selectedIds.length === 0 || !window.confirm(`Hapus ${selectedIds.length} barang dari keranjang?`)) return;
    try {
      let cart = null;
      for (const id of selectedIds) cart = await cartService.removeCartItem(id);
      if (cart) apply(cart);
      setSelected([]);
      toast.success('Barang dihapus dari keranjang.');
    } catch (err) {
      toast.error(getErrorMessage(err));
      refresh().catch(() => undefined);
    }
  };

  if (error) return <div className="mx-auto max-w-3xl px-4 py-10"><ErrorState message={error} onRetry={() => window.location.reload()} /></div>;
  if (!loaded) return <div className="mx-auto max-w-5xl animate-pulse px-4 py-10"><div className="h-64 rounded-3xl bg-slate-200 dark:bg-slate-800" /></div>;
  if (items.length === 0) {
    return (
      <div className="mx-auto max-w-3xl px-4 py-10">
        <EmptyState title="Keranjang masih kosong" description="Yuk pilih kebutuhan dapur Anda." actionText="Mulai Belanja" />
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-6xl space-y-6 px-4 py-8 pb-32 sm:px-6 lg:px-8 lg:pb-8">
      <h1 className="flex items-center gap-3 text-2xl font-black text-slate-900 dark:text-white">
        <ShoppingBag className="h-7 w-7 text-primary-600" /> Keranjang
      </h1>

      <div className="grid gap-6 lg:grid-cols-3">
        <div className="space-y-3 lg:col-span-2">
          <div className="flex items-center justify-between rounded-2xl border border-slate-200 bg-white px-4 py-3 dark:border-slate-800 dark:bg-slate-900">
            <button className="flex items-center gap-2 text-sm font-semibold" onClick={() => setSelected(allSelected ? [] : selectable.map((i) => i.id))}>
              {allSelected ? <CheckSquare className="h-5 w-5 text-primary-600" /> : <Square className="h-5 w-5 text-slate-400" />}
              Pilih semua ({selectable.length})
            </button>
            {selectedIds.length > 0 && (
              <button onClick={removeSelected} className="flex items-center gap-1 text-xs font-bold text-rose-600">
                <Trash2 className="h-4 w-4" /> Hapus
              </button>
            )}
          </div>

          {items.map((item) => {
            const selected = selectedIds.includes(item.id);
            const problem = !item.is_product_active ? 'Produk sudah tidak tersedia' : !item.is_stock_sufficient ? `Stok tinggal ${item.available_stock}` : null;
            return (
              <div key={item.id} className={`flex gap-3 rounded-3xl border bg-white p-4 dark:bg-slate-900 ${selected ? 'border-primary-300 dark:border-primary-800' : 'border-slate-200 dark:border-slate-800'}`}>
                <button onClick={() => toggle(item.id)} disabled={!item.is_product_active} className="self-center disabled:opacity-30" aria-label="Pilih">
                  {selected ? <CheckSquare className="h-5 w-5 text-primary-600" /> : <Square className="h-5 w-5 text-slate-400" />}
                </button>
                <Link to={`/products/${item.product_id}`} className="flex h-20 w-20 flex-shrink-0 items-center justify-center overflow-hidden rounded-2xl bg-slate-100 dark:bg-slate-800">
                  {item.image_url ? <img src={item.image_url} alt="" className="h-full w-full object-cover" /> : <Package className="h-8 w-8 text-slate-300" />}
                </Link>
                <div className="min-w-0 flex-1">
                  <Link to={`/products/${item.product_id}`} className="line-clamp-2 text-sm font-bold text-slate-900 dark:text-white">{item.product_name}</Link>
                  <div className="text-xs text-slate-500">{formatCurrency(item.unit_price)} / {item.unit_name}</div>
                  {problem && (
                    <div className="mt-1 flex items-center gap-1 text-xs font-semibold text-rose-600"><AlertTriangle className="h-3.5 w-3.5" /> {problem}</div>
                  )}
                  <div className="mt-2 flex items-center justify-between gap-2">
                    <div className="flex items-center rounded-xl border border-slate-200 dark:border-slate-700">
                      <button
                        className="p-2 disabled:opacity-30"
                        disabled={busyId === item.id}
                        onClick={() =>
                          item.quantity <= 1
                            ? run(item.id, () => cartService.removeCartItem(item.id))
                            : run(item.id, () => cartService.updateCartItem(item.id, item.quantity - 1))
                        }
                        aria-label="Kurangi"
                      >
                        {item.quantity <= 1 ? <Trash2 className="h-3.5 w-3.5 text-rose-500" /> : <Minus className="h-3.5 w-3.5" />}
                      </button>
                      <span className="w-8 text-center text-sm font-bold">{item.quantity}</span>
                      <button
                        className="p-2 disabled:opacity-30"
                        disabled={busyId === item.id || item.quantity >= item.available_stock}
                        onClick={() => run(item.id, () => cartService.updateCartItem(item.id, item.quantity + 1))}
                        aria-label="Tambah"
                      >
                        <Plus className="h-3.5 w-3.5" />
                      </button>
                    </div>
                    <div className="text-sm font-black text-primary-600">{formatCurrency(item.subtotal)}</div>
                  </div>
                </div>
              </div>
            );
          })}
        </div>

        <div className="fixed inset-x-0 bottom-0 z-40 border-t border-slate-200 bg-white p-4 dark:border-slate-800 dark:bg-slate-900 lg:sticky lg:top-28 lg:h-fit lg:rounded-3xl lg:border lg:p-6 lg:shadow-sm">
          <div className="flex items-center justify-between lg:block lg:space-y-4">
            <div>
              <div className="text-xs text-slate-500">{selectedQty} barang dipilih</div>
              <div className="text-xl font-black text-slate-900 dark:text-white">{formatCurrency(selectedTotal)}</div>
              <div className="hidden text-[11px] text-slate-400 lg:block">Ongkir dihitung di halaman checkout</div>
            </div>
            <button
              disabled={selectedIds.length === 0 || blocked}
              onClick={() => navigate('/checkout')}
              className="flex items-center justify-center gap-2 rounded-2xl bg-primary-600 px-6 py-3 text-sm font-bold text-white shadow-lg shadow-primary-600/30 hover:bg-primary-700 disabled:opacity-40 lg:w-full"
            >
              Checkout <ArrowRight className="h-4 w-4" />
            </button>
          </div>
          {blocked && <p className="mt-2 text-xs text-rose-600">Ada barang terpilih yang stoknya kurang atau tidak tersedia.</p>}
        </div>
      </div>
    </div>
  );
};
