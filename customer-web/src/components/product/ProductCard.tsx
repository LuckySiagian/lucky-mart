import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { motion } from 'framer-motion';
import { Heart, Package, ShoppingBag, Zap } from 'lucide-react';
import { formatCurrency } from '../../lib/utils';
import { useWishlistStore } from '../../store/useWishlistStore';
import { useAuthStore } from '../../store/useAuthStore';
import { useAddToCart } from '../../hooks/useAddToCart';
import type { Product } from '../../types';

export const ProductCard: React.FC<{ product: Product }> = ({ product }) => {
  const { wishlistIds, toggleWishlist } = useWishlistStore();
  const isB2B = useAuthStore((s) => s.user?.role === 'B2B_CUSTOMER' && s.user?.b2b_profile?.verification_status === 'APPROVED');
  const { add, adding } = useAddToCart();
  const [unitIndex, setUnitIndex] = useState(0);

  const unit = product.units[unitIndex] || product.units[0];
  const isLiked = wishlistIds.includes(product.id);
  const outOfStock = !unit || unit.stock_quantity <= 0;
  const hasWholesale = unit && unit.wholesale_price < unit.retail_price;

  return (
    <motion.div
      whileHover={{ y: -4 }}
      className="group relative flex w-full flex-col overflow-hidden rounded-3xl border border-slate-200 bg-white p-3 shadow-sm transition-all duration-300 hover:shadow-xl dark:border-slate-800 dark:bg-slate-900 sm:p-4"
    >
      <button
        onClick={() => toggleWishlist(product.id)}
        aria-label={isLiked ? 'Hapus dari favorit' : 'Tambah ke favorit'}
        className="absolute right-5 top-5 z-10 rounded-full bg-white/80 p-2 text-slate-400 shadow-sm backdrop-blur-md transition-colors hover:text-rose-500 dark:bg-slate-800/80 sm:right-6 sm:top-6"
      >
        <Heart className={`h-4 w-4 ${isLiked ? 'fill-rose-500 text-rose-500' : ''}`} />
      </button>

      <Link to={`/products/${product.id}`} className="relative mb-3 flex h-36 w-full items-center justify-center overflow-hidden rounded-2xl bg-slate-100 dark:bg-slate-800 sm:h-44">
        {product.image_url ? (
          <img src={product.image_url} alt={product.name} loading="lazy" className="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105" />
        ) : (
          <Package className="h-14 w-14 text-slate-300 dark:text-slate-600" />
        )}
        {outOfStock ? (
          <span className="absolute left-2 top-2 rounded-md bg-slate-800 px-2 py-0.5 text-[10px] font-bold uppercase text-white">Stok Habis</span>
        ) : unit.stock_quantity < 10 ? (
          <span className="absolute left-2 top-2 rounded-md bg-rose-500 px-2 py-0.5 text-[10px] font-bold uppercase text-white">Sisa {unit.stock_quantity}</span>
        ) : null}
      </Link>

      <div className="flex flex-1 flex-col">
        <span className="mb-1 text-[11px] font-semibold text-slate-400">{product.category?.name}</span>
        <Link to={`/products/${product.id}`} className="mb-2 line-clamp-2 text-sm font-bold text-slate-900 transition-colors hover:text-primary-600 dark:text-white">
          {product.name}
        </Link>

        {unit && (
          <div className="mt-auto space-y-1 border-t border-slate-100 pt-2 dark:border-slate-800">
            {product.units.length > 1 ? (
              <select
                value={unitIndex}
                onChange={(e) => setUnitIndex(Number(e.target.value))}
                className="rounded-lg border-0 bg-slate-100 px-2 py-0.5 text-xs font-bold text-slate-800 focus:ring-1 focus:ring-primary-500 dark:bg-slate-800 dark:text-slate-200"
                aria-label="Pilih satuan"
              >
                {product.units.map((u, idx) => (
                  <option key={u.id} value={idx}>{u.unit_name}</option>
                ))}
              </select>
            ) : (
              <span className="inline-flex rounded-md bg-slate-100 px-2 py-0.5 text-[10px] font-bold text-slate-600 dark:bg-slate-800 dark:text-slate-300">{unit.unit_name}</span>
            )}

            <div className="text-base font-black text-slate-900 dark:text-white">
              {formatCurrency(unit.retail_price)}
              <span className="ml-1 text-[10px] font-normal text-slate-400">/ {unit.unit_name}</span>
            </div>

            {hasWholesale && (
              <div className={`text-[11px] font-semibold ${isB2B ? 'text-primary-600 dark:text-primary-400' : 'text-slate-500'}`}>
                Grosir {formatCurrency(unit.wholesale_price)} <span className="font-normal">(min {unit.min_wholesale_qty}{isB2B ? '' : ', mitra toko'})</span>
              </div>
            )}
          </div>
        )}
      </div>

      <div className="mt-3 flex items-center gap-2">
        <button
          onClick={() => unit && add(product, unit)}
          disabled={outOfStock || adding}
          title="Tambah ke keranjang"
          aria-label="Tambah ke keranjang"
          className="flex flex-shrink-0 items-center justify-center rounded-2xl border border-primary-600/30 bg-primary-50 p-2.5 text-primary-600 shadow-sm transition-all hover:bg-primary-100 active:scale-95 disabled:opacity-40 dark:bg-primary-950/40 dark:text-primary-400"
        >
          <ShoppingBag className="h-4 w-4" />
        </button>
        <button
          onClick={() => unit && add(product, unit, 1, true)}
          disabled={outOfStock || adding}
          className="flex flex-1 items-center justify-center gap-1.5 rounded-2xl bg-primary-600 px-3 py-2.5 text-xs font-extrabold text-white shadow-md shadow-primary-600/20 transition-all hover:bg-primary-700 active:scale-95 disabled:opacity-40"
        >
          <Zap className="h-3.5 w-3.5 fill-amber-300 text-amber-300" />
          <span>{outOfStock ? 'Habis' : 'Beli'}</span>
        </button>
      </div>
    </motion.div>
  );
};
