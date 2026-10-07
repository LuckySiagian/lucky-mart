import React, { useEffect, useState, useRef } from 'react';
import { useSearchParams } from 'react-router-dom';
import { productService } from '../../services/apiServices';
import { ProductCard } from '../../components/product/ProductCard';
import { LoadingSkeleton } from '../../components/common/LoadingSkeleton';
import { EmptyState } from '../../components/common/EmptyState';
import { Search, Filter } from 'lucide-react';

export const ProductListPage: React.FC = () => {
  const [searchParams] = useSearchParams();
  const searchArg = searchParams.get('search') || '';
  const categoryArg = searchParams.get('category_id') || searchParams.get('category') || '';

  const [products, setProducts] = useState<any[]>([]);
  const [categories, setCategories] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState(searchArg);
  const [selectedCategory, setSelectedCategory] = useState(categoryArg);

  const selectedCategoryRef = useRef(selectedCategory);
  selectedCategoryRef.current = selectedCategory;

  const loadData = async (isInitial = false) => {
    if (isInitial) setLoading(true);

    try {
      const prodRes = await productService.getProducts(1, 100, selectedCategoryRef.current || undefined);
      const fetchedProds = Array.isArray(prodRes.data)
        ? prodRes.data
        : (prodRes.data?.items || (Array.isArray(prodRes) ? prodRes : []));
      setProducts(fetchedProds);
    } catch (e) {
      if (isInitial) setProducts([]);
    }

    try {
      const catRes = await productService.getCategories();
      const fetchedCats = Array.isArray(catRes.data)
        ? catRes.data
        : (Array.isArray(catRes) ? catRes : []);
      setCategories(fetchedCats);
    } catch (e) {
      if (isInitial) setCategories([]);
    }

    if (isInitial) setLoading(false);
  };

  useEffect(() => {
    loadData(true);

    // Setup 3-second realtime polling interval
    const interval = setInterval(() => {
      loadData(false);
    }, 3000);

    return () => clearInterval(interval);
  }, [selectedCategory]);

  const filtered = Array.isArray(products)
    ? products.filter((p) =>
        p.name?.toLowerCase().includes(searchQuery.toLowerCase())
      )
    : [];

  return (
    <div className="space-y-8 pb-16">
      {/* Search & Filter Header */}
      <div className="rounded-3xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 p-6 shadow-sm space-y-4">
        <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
          <div className="relative w-full sm:w-96">
            <Search className="absolute left-3.5 top-3.5 h-4 w-4 text-slate-400" />
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="Cari produk sembako..."
              className="w-full rounded-2xl bg-slate-100 dark:bg-slate-800 pl-10 pr-4 py-2.5 text-sm text-slate-900 dark:text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-primary-500"
            />
          </div>

          <div className="flex items-center space-x-2 w-full sm:w-auto">
            <Filter className="h-4 w-4 text-slate-400" />
            <select
              value={selectedCategory}
              onChange={(e) => setSelectedCategory(e.target.value)}
              className="rounded-2xl bg-slate-100 dark:bg-slate-800 px-4 py-2.5 text-xs font-bold text-slate-700 dark:text-slate-200 focus:outline-none"
            >
              <option value="">Semua Kategori Sembako</option>
              {categories.map((c) => (
                <option key={c.id} value={c.id}>
                  {c.name}
                </option>
              ))}
            </select>
          </div>
        </div>

        {/* Category Pills Filter */}
        <div className="flex items-center space-x-2 overflow-x-auto pt-2 scrollbar-none">
          <button
            onClick={() => setSelectedCategory('')}
            className={`flex-shrink-0 px-4 py-2 rounded-xl text-xs font-bold transition-all ${
              !selectedCategory
                ? 'bg-primary-600 text-white shadow-md shadow-primary-600/20'
                : 'bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 hover:bg-slate-200'
            }`}
          >
            Semua
          </button>
          {categories.map((c) => (
            <button
              key={c.id}
              onClick={() => setSelectedCategory(c.id)}
              className={`flex-shrink-0 px-4 py-2 rounded-xl text-xs font-bold transition-all ${
                selectedCategory === c.id
                  ? 'bg-primary-600 text-white shadow-md shadow-primary-600/20'
                  : 'bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 hover:bg-slate-200'
              }`}
            >
              {c.name}
            </button>
          ))}
        </div>
      </div>

      {/* Product List Grid */}
      {loading ? (
        <LoadingSkeleton count={8} />
      ) : filtered.length === 0 ? (
        <EmptyState title="Produk Tidak Ditemukan" description="Belum ada produk sembako yang sesuai dengan kriteria pencarian Anda." />
      ) : (
        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-4 sm:gap-6">
          {filtered.map((p) => (
            <ProductCard key={p.id} product={p} />
          ))}
        </div>
      )}
    </div>
  );
};
