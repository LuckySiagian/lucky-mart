import React, { useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { Heart, Search } from 'lucide-react';
import { productService } from '../../services/apiServices';
import { ProductCard } from '../../components/product/ProductCard';
import { LoadingSkeleton } from '../../components/common/LoadingSkeleton';
import { EmptyState } from '../../components/common/EmptyState';
import { ErrorState } from '../../components/common/ErrorState';
import { useWishlistStore } from '../../store/useWishlistStore';
import { getErrorMessage } from '../../lib/utils';
import type { Category, Product } from '../../types';

const PAGE_SIZE = 24;

export const ProductListPage: React.FC<{ wishlistOnly?: boolean }> = ({ wishlistOnly = false }) => {
  const [searchParams, setSearchParams] = useSearchParams();
  const search = searchParams.get('search') || '';
  const categoryId = searchParams.get('category_id') || '';
  const wishlistIds = useWishlistStore((s) => s.wishlistIds);

  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchInput, setSearchInput] = useState(search);

  useEffect(() => {
    productService.getCategories().then(setCategories).catch(() => undefined);
  }, []);

  useEffect(() => setSearchInput(search), [search]);
  useEffect(() => setPage(1), [search, categoryId]);

  useEffect(() => {
    let alive = true;
    const load = async () => {
      setLoading(true);
      setError(null);
      try {
        if (wishlistOnly) {
          const results = await Promise.allSettled(wishlistIds.map((id) => productService.getProductById(id)));
          if (!alive) return;
          setProducts(results.flatMap((r) => (r.status === 'fulfilled' ? [r.value] : [])));
          setTotalPages(1);
        } else {
          const res = await productService.getProducts({ page, limit: PAGE_SIZE, search, category_id: categoryId });
          if (!alive) return;
          setProducts((prev) => (page === 1 ? res.items : [...prev, ...res.items]));
          setTotalPages(res.meta.total_pages);
        }
      } catch (err) {
        if (alive) setError(getErrorMessage(err, 'Gagal memuat produk.'));
      } finally {
        if (alive) setLoading(false);
      }
    };
    load();
    return () => {
      alive = false;
    };
  }, [wishlistOnly, wishlistIds, page, search, categoryId]);

  const setParam = (key: string, value: string) => {
    const next = new URLSearchParams(searchParams);
    if (value) next.set(key, value);
    else next.delete(key);
    setSearchParams(next);
  };

  return (
    <div className="mx-auto max-w-7xl space-y-6 px-4 py-8 pb-16 sm:px-6 lg:px-8">
      {wishlistOnly ? (
        <h1 className="flex items-center gap-2 text-2xl font-black text-slate-900 dark:text-white">
          <Heart className="h-6 w-6 fill-rose-500 text-rose-500" /> Produk Favorit
        </h1>
      ) : (
        <div className="space-y-4 rounded-3xl border border-slate-200 bg-white p-4 shadow-sm dark:border-slate-800 dark:bg-slate-900 sm:p-6">
          <form
            onSubmit={(e) => {
              e.preventDefault();
              setParam('search', searchInput.trim());
            }}
            className="relative"
          >
            <Search className="absolute left-3.5 top-3.5 h-4 w-4 text-slate-400" />
            <input
              type="search"
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
              onBlur={() => searchInput.trim() !== search && setParam('search', searchInput.trim())}
              placeholder="Cari beras, minyak, gula..."
              className="w-full rounded-2xl bg-slate-100 py-2.5 pl-10 pr-4 text-sm text-slate-900 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-primary-500 dark:bg-slate-800 dark:text-white"
            />
          </form>
          <div className="flex items-center gap-2 overflow-x-auto pb-1">
            {[{ id: '', name: 'Semua' }, ...categories].map((c) => (
              <button
                key={c.id || 'all'}
                onClick={() => setParam('category_id', c.id)}
                className={`flex-shrink-0 rounded-xl px-4 py-2 text-xs font-bold transition-all ${
                  categoryId === c.id
                    ? 'bg-primary-600 text-white shadow-md shadow-primary-600/20'
                    : 'bg-slate-100 text-slate-600 hover:bg-slate-200 dark:bg-slate-800 dark:text-slate-300'
                }`}
              >
                {c.name}
              </button>
            ))}
          </div>
        </div>
      )}

      {error ? (
        <ErrorState message={error} onRetry={() => setPage(1)} />
      ) : loading && products.length === 0 ? (
        <LoadingSkeleton count={8} />
      ) : products.length === 0 ? (
        wishlistOnly ? (
          <EmptyState title="Belum ada favorit" description="Ketuk ikon hati pada produk untuk menyimpannya di sini." actionText="Lihat Produk" />
        ) : (
          <EmptyState title="Produk tidak ditemukan" description="Coba kata kunci atau kategori lain." />
        )
      ) : (
        <>
          <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 sm:gap-6 lg:grid-cols-4 xl:grid-cols-6">
            {products.map((p) => (
              <ProductCard key={p.id} product={p} />
            ))}
          </div>
          {!wishlistOnly && page < totalPages && (
            <div className="text-center">
              <button
                disabled={loading}
                onClick={() => setPage((p) => p + 1)}
                className="rounded-2xl border border-slate-200 bg-white px-6 py-2.5 text-sm font-bold text-slate-700 hover:bg-slate-50 disabled:opacity-50 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-200"
              >
                {loading ? 'Memuat...' : 'Muat lebih banyak'}
              </button>
            </div>
          )}
        </>
      )}
    </div>
  );
};
