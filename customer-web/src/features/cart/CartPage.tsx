import React, { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { cartService } from '../../services/apiServices';
import { useCartStore, CART_MAX_AGE_MS } from '../../store/useCartStore';
import { formatCurrency } from '../../lib/utils';
import { Trash2, Plus, Minus, ShoppingBag, ArrowRight, ShieldCheck, Clock, CheckSquare, Square } from 'lucide-react';
import { toast } from 'sonner';

export const CartPage: React.FC = () => {
  const navigate = useNavigate();
  const { items, totalAmount, totalItems, setCart, addItem, updateQuantity, removeItem, removeMultipleItems, clearCart, cleanExpiredItems, setSelectedForCheckout } = useCartStore();
  const [loading, setLoading] = useState(true);
  const [selectedIds, setSelectedIds] = useState<string[]>([]);

  const fetchCart = async () => {
    setLoading(true);
    try {
      const res = await cartService.getCart();
      if (res.data) {
        setCart(res.data.items || [], res.data.total_items || 0, res.data.total_amount || 0);
      }
    } catch (e) {
      // Mock local fallback
    } finally {
      const removed = cleanExpiredItems();
      if (removed > 0) {
        toast.info(`${removed} item yang melebihi batas 30 hari telah dibersihkan secara otomatis.`);
      }
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCart();
  }, []);

  // Keep only valid selections when items change (defaults to unselected on open)
  useEffect(() => {
    if (items.length > 0) {
      setSelectedIds((prev) => prev.filter((id) => items.some((item) => item.id === id)));
    } else {
      setSelectedIds([]);
    }
  }, [items]);

  const getRemainingDays = (addedAt?: number) => {
    if (!addedAt) return 30;
    const elapsedMs = Date.now() - addedAt;
    const remainingMs = CART_MAX_AGE_MS - elapsedMs;
    const remainingDays = Math.ceil(remainingMs / (1000 * 60 * 60 * 24));
    return Math.max(1, Math.min(30, remainingDays));
  };

  const handleToggleSelect = (id: string) => {
    setSelectedIds((prev) =>
      prev.includes(id) ? prev.filter((i) => i !== id) : [...prev, id]
    );
  };

  const handleToggleSelectAll = () => {
    if (selectedIds.length === items.length) {
      setSelectedIds([]);
    } else {
      setSelectedIds(items.map((i) => i.id));
    }
  };

  const handleUpdateQuantity = async (id: string, newQty: number) => {
    if (newQty <= 0) {
      await handleRemoveItem(id);
      return;
    }
    updateQuantity(id, newQty);
    try {
      await cartService.updateCartItem(id, { quantity: newQty });
    } catch (e) {
      // Guest or local fallback handled silently
    }
  };

  const handleRemoveItem = async (id: string) => {
    try {
      await cartService.removeCartItem(id);
      removeItem(id);
      setSelectedIds((prev) => prev.filter((i) => i !== id));
      toast.success('Item dihapus dari keranjang');
    } catch (e) {
      toast.error('Gagal menghapus item');
    }
  };

  const handleDeleteSelected = async () => {
    if (selectedIds.length === 0) {
      toast.error('Pilih minimal 1 produk untuk dihapus');
      return;
    }

    if (!window.confirm(`Hapus ${selectedIds.length} produk yang dipilih dari keranjang?`)) {
      return;
    }

    try {
      await Promise.all(selectedIds.map((id) => cartService.removeCartItem(id)));
      removeMultipleItems(selectedIds);
      setSelectedIds([]);
      toast.success(`${selectedIds.length} produk berhasil dihapus dari keranjang.`);
    } catch (e) {
      toast.error('Gagal menghapus beberapa produk yang dipilih.');
    }
  };

  const handleClearCart = async () => {
    try {
      await cartService.clearCart();
      clearCart();
      setSelectedIds([]);
      toast.success('Keranjang dibersihkan');
    } catch (e) {
      toast.error('Gagal mengosongkan keranjang');
    }
  };

  const handleCheckoutSelected = () => {
    if (selectedIds.length === 0) {
      toast.error('Pilih minimal 1 produk untuk di-checkout!');
      return;
    }
    setSelectedForCheckout(selectedIds);
    navigate('/checkout');
  };

  // Selected items statistics
  const selectedItems = items.filter((i) => selectedIds.includes(i.id));
  const selectedTotalQuantity = selectedItems.reduce((sum, i) => sum + i.quantity, 0);
  const selectedTotalAmount = selectedItems.reduce((sum, i) => sum + i.subtotal, 0);
  const isAllSelected = items.length > 0 && selectedIds.length === items.length;

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-8">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <h1 className="text-2xl font-black text-slate-900 dark:text-white flex items-center space-x-3">
          <ShoppingBag className="h-7 w-7 text-primary-600" />
          <span>Keranjang Belanja Sembako Saya ({totalItems} Item)</span>
        </h1>

        <div className="inline-flex items-center space-x-2 rounded-2xl bg-amber-50 dark:bg-amber-950/40 px-4 py-2 border border-amber-200 dark:border-amber-800 text-amber-800 dark:text-amber-300 text-xs font-semibold">
          <Clock className="h-4 w-4 flex-shrink-0 text-amber-600" />
          <span>Batas simpan produk di keranjang: <b>30 Hari</b></span>
        </div>
      </div>

      {items.length === 0 ? (
        <div className="rounded-3xl border border-dashed border-slate-300 dark:border-slate-800 bg-white dark:bg-slate-900 p-12 text-center">
          <ShoppingBag className="h-16 w-16 text-slate-300 dark:text-slate-600 mx-auto mb-4" />
          <h3 className="text-lg font-bold text-slate-900 dark:text-white mb-1">Keranjang Belanja Masih Kosong</h3>
          <p className="text-sm text-slate-500 mb-6">Jelajahi katalog sembako segar dan pilih barang kebutuhan Anda.</p>
          <Link to="/products" className="rounded-2xl bg-primary-600 px-6 py-3 text-sm font-bold text-white shadow-lg shadow-primary-600/20">
            Mulai Belanja Sekarang
          </Link>
        </div>
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Cart Items List */}
          <div className="lg:col-span-2 space-y-4">
            {/* Header Control Bar */}
            <div className="flex flex-wrap justify-between items-center bg-white dark:bg-slate-900 p-4 rounded-2xl border border-slate-200 dark:border-slate-800 gap-3">
              <label className="flex items-center space-x-3 cursor-pointer select-none">
                <input
                  type="checkbox"
                  checked={isAllSelected}
                  onChange={handleToggleSelectAll}
                  className="h-5 w-5 rounded-md border-slate-300 text-primary-600 focus:ring-primary-500 cursor-pointer accent-primary-600"
                />
                <span className="text-xs font-bold text-slate-700 dark:text-slate-300">
                  Pilih Semua ({selectedIds.length}/{items.length} Produk)
                </span>
              </label>

              <div className="flex items-center space-x-3">
                {selectedIds.length > 0 && (
                  <button
                    onClick={handleDeleteSelected}
                    className="flex items-center space-x-1 text-xs font-bold text-rose-600 hover:text-rose-700 bg-rose-50 dark:bg-rose-950/40 px-3 py-1.5 rounded-xl border border-rose-200 dark:border-rose-800 transition-colors"
                  >
                    <Trash2 className="h-3.5 w-3.5" />
                    <span>Hapus Terpilih ({selectedIds.length})</span>
                  </button>
                )}

                <button onClick={handleClearCart} className="text-xs font-bold text-slate-400 hover:text-rose-500 hover:underline">
                  Kosongkan Semua
                </button>
              </div>
            </div>

            {/* Cart Items */}
            {items.map((item) => {
              const remainingDays = getRemainingDays(item.added_at);
              const isSelected = selectedIds.includes(item.id);

              return (
                <div
                  key={item.id}
                  className={`flex flex-col sm:flex-row items-start sm:items-center justify-between p-4 rounded-3xl border transition-all gap-4 ${
                    isSelected
                      ? 'border-primary-500/50 bg-primary-50/20 dark:bg-primary-950/20 shadow-sm'
                      : 'border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900'
                  }`}
                >
                  <div className="flex items-center space-x-3 w-full sm:w-auto">
                    {/* Item Checkbox */}
                    <input
                      type="checkbox"
                      checked={isSelected}
                      onChange={() => handleToggleSelect(item.id)}
                      className="h-5 w-5 rounded-md border-slate-300 text-primary-600 focus:ring-primary-500 cursor-pointer accent-primary-600 flex-shrink-0"
                    />

                    <div className="flex items-center space-x-4">
                      <div className="h-16 w-16 rounded-2xl bg-slate-100 dark:bg-slate-800 flex items-center justify-center font-bold text-slate-400 flex-shrink-0 text-xs">
                        SEMBAKO
                      </div>
                      <div>
                        <h4 className="font-bold text-sm text-slate-900 dark:text-white">{item.product_name}</h4>
                        <div className="flex items-center space-x-2 mt-0.5">
                          <span className="text-xs text-slate-400">Satuan: {item.unit_name}</span>
                          <span className="inline-flex items-center space-x-1 text-[10px] font-bold text-amber-700 dark:text-amber-300 bg-amber-100 dark:bg-amber-900/40 px-2 py-0.5 rounded-full">
                            <Clock className="h-3 w-3" />
                            <span>Tersisa {remainingDays} hari</span>
                          </span>
                        </div>
                        <div className="text-sm font-black text-primary-600 mt-1">{formatCurrency(item.unit_price)}</div>
                      </div>
                    </div>
                  </div>

                  <div className="flex items-center space-x-4 w-full sm:w-auto justify-between sm:justify-end border-t sm:border-t-0 border-slate-100 dark:border-slate-800 pt-3 sm:pt-0">
                    {/* Quantity Stepper Controls */}
                    <div className="flex items-center space-x-1.5 bg-slate-100 dark:bg-slate-800 p-1.5 rounded-2xl border border-slate-200 dark:border-slate-700">
                      <button
                        onClick={() => handleUpdateQuantity(item.id, item.quantity - 1)}
                        className="p-1 rounded-xl text-slate-600 hover:bg-white dark:text-slate-300 dark:hover:bg-slate-700 transition-colors shadow-sm"
                        title="Kurangi 1"
                      >
                        <Minus className="h-3.5 w-3.5" />
                      </button>

                      <span className="text-xs font-black px-2 text-slate-900 dark:text-white min-w-[28px] text-center select-none">
                        {item.quantity}
                      </span>

                      <button
                        onClick={() => handleUpdateQuantity(item.id, item.quantity + 1)}
                        className="p-1 rounded-xl text-slate-600 hover:bg-white dark:text-slate-300 dark:hover:bg-slate-700 transition-colors shadow-sm"
                        title="Tambah 1"
                      >
                        <Plus className="h-3.5 w-3.5" />
                      </button>
                    </div>

                    <div className="text-right">
                      <span className="text-xs text-slate-400">Subtotal ({item.quantity} {item.unit_name})</span>
                      <div className="font-bold text-slate-900 dark:text-white">{formatCurrency(item.subtotal)}</div>
                    </div>

                    <button
                      onClick={() => handleRemoveItem(item.id)}
                      className="p-2 text-slate-400 hover:text-rose-500 rounded-xl hover:bg-rose-50 dark:hover:bg-rose-950/30 transition-colors"
                      title="Hapus Item Ini"
                    >
                      <Trash2 className="h-4 w-4" />
                    </button>
                  </div>
                </div>
              );
            })}
          </div>

          {/* Sticky Checkout Summary Box */}
          <div className="h-fit sticky top-24 rounded-3xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 p-6 space-y-6 shadow-sm">
            <h3 className="text-base font-extrabold text-slate-900 dark:text-white border-b border-slate-100 dark:border-slate-800 pb-3">
              Ringkasan Belanja
            </h3>

            <div className="space-y-2 text-sm">
              <div className="flex justify-between text-slate-600 dark:text-slate-400">
                <span>Produk Terpilih</span>
                <span className="font-semibold text-slate-900 dark:text-white">{selectedIds.length} dari {items.length} Item</span>
              </div>
              <div className="flex justify-between text-slate-600 dark:text-slate-400">
                <span>Total Kuantitas</span>
                <span className="font-semibold text-slate-900 dark:text-white">{selectedTotalQuantity} Item</span>
              </div>
              <div className="flex justify-between text-slate-600 dark:text-slate-400">
                <span>Subtotal Terpilih</span>
                <span className="font-bold text-slate-900 dark:text-white">{formatCurrency(selectedTotalAmount)}</span>
              </div>
            </div>

            <div className="pt-3 border-t border-slate-100 dark:border-slate-800 flex justify-between items-center">
              <span className="font-bold text-slate-900 dark:text-white">Total Checkout</span>
              <span className="text-xl font-black text-primary-600">{formatCurrency(selectedTotalAmount)}</span>
            </div>

            <button
              onClick={handleCheckoutSelected}
              disabled={selectedIds.length === 0}
              className={`flex w-full items-center justify-center space-x-2 rounded-2xl py-3.5 text-sm font-bold text-white transition-all shadow-lg ${
                selectedIds.length > 0
                  ? 'bg-primary-600 hover:bg-primary-700 shadow-primary-600/30 cursor-pointer'
                  : 'bg-slate-300 dark:bg-slate-800 cursor-not-allowed opacity-60'
              }`}
            >
              <span>Lanjut ke Checkout ({selectedIds.length})</span>
              <ArrowRight className="h-4 w-4" />
            </button>
          </div>
        </div>
      )}
    </div>
  );
};
