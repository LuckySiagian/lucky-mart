import React, { useEffect, useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import { ArrowLeft, Heart, Minus, Package, Plus, ShoppingBag, Zap } from 'lucide-react';
import { productService } from '../../services/apiServices';
import { useAddToCart } from '../../hooks/useAddToCart';
import { useWishlistStore } from '../../store/useWishlistStore';
import { useAuthStore } from '../../store/useAuthStore';
import { ErrorState } from '../../components/common/ErrorState';
import { formatCurrency, getErrorMessage } from '../../lib/utils';
import type { Product } from '../../types';

export const ProductDetailPage: React.FC = () => {
  const { id = '' } = useParams();
  const [product, setProduct] = useState<Product | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [unitIndex, setUnitIndex] = useState(0);
  const [qty, setQty] = useState(1);
  const [imageIndex, setImageIndex] = useState(0);
  const { add, adding } = useAddToCart();
  const { wishlistIds, toggleWishlist } = useWishlistStore();
  const user = useAuthStore((s) => s.user);
  const isB2B = user?.role === 'B2B_CUSTOMER' && user.b2b_profile?.verification_status === 'APPROVED';

  useEffect(() => {
    setProduct(null);
    setError(null);
    productService
      .getProductById(id)
      .then((p) => {
        setProduct(p);
        setUnitIndex(0);
        setQty(1);
      })
      .catch((err) => setError(getErrorMessage(err, 'Produk tidak ditemukan.')));
  }, [id]);

  if (error) {
    return (
      <div className="mx-auto max-w-3xl px-4 py-10">
        <ErrorState message={error} />
        <div className="text-center"><Link to="/products" className="text-sm font-bold text-primary-600">Kembali ke daftar produk</Link></div>
      </div>
    );
  }
  if (!product) {
    return <div className="mx-auto max-w-5xl animate-pulse px-4 py-10"><div className="h-96 rounded-3xl bg-slate-200 dark:bg-slate-800" /></div>;
  }

  const unit = product.units[unitIndex];
  const stock = unit?.stock_quantity ?? 0;
  const wholesaleApplies = isB2B && unit && qty >= unit.min_wholesale_qty && unit.wholesale_price < unit.retail_price;
  const price = unit ? (wholesaleApplies ? unit.wholesale_price : unit.retail_price) : 0;
  const images = product.images.length > 0 ? product.images.map((i) => i.image_url) : [];
  const liked = wishlistIds.includes(product.id);

  return (
    <div className="mx-auto max-w-6xl px-4 py-6 sm:px-6 lg:px-8">
      <Link to="/products" className="mb-4 inline-flex items-center gap-1 text-sm font-semibold text-slate-500 hover:text-primary-600">
        <ArrowLeft className="h-4 w-4" /> Produk
      </Link>

      <div className="grid gap-8 lg:grid-cols-2">
        <div>
          <div className="flex aspect-square items-center justify-center overflow-hidden rounded-3xl border border-slate-200 bg-slate-100 dark:border-slate-800 dark:bg-slate-800">
            {images.length > 0 ? (
              <img src={images[imageIndex]} alt={product.name} className="h-full w-full object-cover" />
            ) : (
              <Package className="h-24 w-24 text-slate-300" />
            )}
          </div>
          {images.length > 1 && (
            <div className="mt-3 flex gap-2">
              {images.map((src, i) => (
                <button key={src} onClick={() => setImageIndex(i)} className={`h-16 w-16 overflow-hidden rounded-xl border-2 ${i === imageIndex ? 'border-primary-600' : 'border-transparent'}`}>
                  <img src={src} alt="" className="h-full w-full object-cover" />
                </button>
              ))}
            </div>
          )}
        </div>

        <div className="space-y-5">
          <div>
            <Link to={`/products?category_id=${product.category_id}`} className="text-xs font-bold uppercase tracking-wide text-primary-600">
              {product.category?.name}
            </Link>
            <h1 className="mt-1 text-2xl font-black text-slate-900 dark:text-white sm:text-3xl">{product.name}</h1>
            <p className="mt-1 font-mono text-xs text-slate-400">SKU {product.sku}</p>
          </div>

          {unit && (
            <div className="rounded-3xl border border-slate-200 bg-white p-5 dark:border-slate-800 dark:bg-slate-900">
              <div className="text-3xl font-black text-slate-900 dark:text-white">
                {formatCurrency(price)}
                <span className="ml-1 text-sm font-normal text-slate-400">/ {unit.unit_name}</span>
              </div>
              {unit.wholesale_price < unit.retail_price && (
                <p className="mt-1 text-sm text-slate-500">
                  Harga grosir <b className="text-primary-600">{formatCurrency(unit.wholesale_price)}</b> untuk pembelian minimal {unit.min_wholesale_qty} {unit.unit_name}
                  {isB2B ? '.' : ' — khusus mitra toko yang sudah terdaftar.'}
                </p>
              )}

              {product.units.length > 1 && (
                <div className="mt-4">
                  <div className="mb-2 text-xs font-semibold text-slate-500">Pilih satuan</div>
                  <div className="flex flex-wrap gap-2">
                    {product.units.map((u, i) => (
                      <button
                        key={u.id}
                        onClick={() => {
                          setUnitIndex(i);
                          setQty(1);
                        }}
                        className={`rounded-xl border px-4 py-2 text-sm font-semibold ${
                          i === unitIndex ? 'border-primary-600 bg-primary-50 text-primary-700 dark:bg-primary-950/40' : 'border-slate-200 dark:border-slate-700'
                        }`}
                      >
                        {u.unit_name}
                      </button>
                    ))}
                  </div>
                </div>
              )}

              <div className="mt-4 flex items-center gap-4">
                <div className="flex items-center rounded-2xl border border-slate-200 dark:border-slate-700">
                  <button className="p-3 disabled:opacity-30" disabled={qty <= 1} onClick={() => setQty((q) => q - 1)} aria-label="Kurangi">
                    <Minus className="h-4 w-4" />
                  </button>
                  <input
                    type="number"
                    min={1}
                    max={stock}
                    value={qty}
                    onChange={(e) => setQty(Math.max(1, Math.min(stock || 1, Number(e.target.value) || 1)))}
                    className="w-14 border-0 bg-transparent text-center text-sm font-bold focus:ring-0"
                  />
                  <button className="p-3 disabled:opacity-30" disabled={qty >= stock} onClick={() => setQty((q) => q + 1)} aria-label="Tambah">
                    <Plus className="h-4 w-4" />
                  </button>
                </div>
                <span className={`text-sm ${stock === 0 ? 'font-bold text-rose-600' : 'text-slate-500'}`}>
                  {stock === 0 ? 'Stok habis' : `Stok ${stock}`}
                </span>
              </div>

              <div className="mt-5 flex gap-3">
                <button
                  onClick={() => toggleWishlist(product.id)}
                  className="rounded-2xl border border-slate-200 p-3 dark:border-slate-700"
                  aria-label="Favorit"
                >
                  <Heart className={`h-5 w-5 ${liked ? 'fill-rose-500 text-rose-500' : 'text-slate-400'}`} />
                </button>
                <button
                  disabled={stock === 0 || adding}
                  onClick={() => add(product, unit, qty)}
                  className="flex flex-1 items-center justify-center gap-2 rounded-2xl border border-primary-600 py-3 text-sm font-bold text-primary-700 hover:bg-primary-50 disabled:opacity-40 dark:text-primary-400 dark:hover:bg-primary-950/40"
                >
                  <ShoppingBag className="h-4 w-4" /> Keranjang
                </button>
                <button
                  disabled={stock === 0 || adding}
                  onClick={() => add(product, unit, qty, true)}
                  className="flex flex-1 items-center justify-center gap-2 rounded-2xl bg-primary-600 py-3 text-sm font-bold text-white shadow-lg shadow-primary-600/30 hover:bg-primary-700 disabled:opacity-40"
                >
                  <Zap className="h-4 w-4" /> Beli Sekarang
                </button>
              </div>
            </div>
          )}

          {product.description && (
            <div>
              <h2 className="mb-2 text-sm font-bold text-slate-900 dark:text-white">Deskripsi</h2>
              <p className="whitespace-pre-line text-sm leading-relaxed text-slate-600 dark:text-slate-400">{product.description}</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
