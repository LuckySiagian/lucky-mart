import React, { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import { productService } from '../../services/apiServices';
import { ProductCard } from '../../components/product/ProductCard';
import { LoadingSkeleton } from '../../components/common/LoadingSkeleton';
import { ArrowRight, Flame, ShieldCheck, Truck, Sparkles, Store } from 'lucide-react';
import { Link } from 'react-router-dom';
import type { Category, Product } from '../../types';

export const HomePage: React.FC = () => {
  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    Promise.allSettled([productService.getProducts({ limit: 12 }), productService.getCategories()])
      .then(([prod, cats]) => {
        if (prod.status === 'fulfilled') setProducts(prod.value.items);
        if (cats.status === 'fulfilled') setCategories(cats.value);
      })
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="space-y-12 pb-16">
      {/* Hero Banner Slider Section */}
      <section className="relative overflow-hidden bg-gradient-to-r from-primary-900 via-primary-800 to-slate-900 py-16 sm:py-20 px-4 sm:px-6 lg:px-8 text-white">
        <div className="max-w-7xl mx-auto grid grid-cols-1 lg:grid-cols-2 gap-8 items-center">
          <motion.div initial={{ opacity: 0, x: -20 }} animate={{ opacity: 1, x: 0 }} className="space-y-6">
            <span className="inline-flex items-center space-x-2 rounded-full bg-primary-500/20 px-3.5 py-1 text-xs font-bold text-primary-300 border border-primary-500/30">
              <Sparkles className="h-4 w-4" />
              <span>Toko Sembako Keluarga</span>
            </span>
            <h1 className="text-3xl sm:text-5xl font-black tracking-tight leading-tight">
              Belanja Sembako Lengkap, Segar & Harga Grosir Direct!
            </h1>
            <p className="text-sm sm:text-base text-slate-300 max-w-lg leading-relaxed">
              Pasokan beras, minyak goreng, gula, dan telur segar langsung dari pemasok utama. Siap kirim untuk rumah tangga & mitra toko kelontong (Kasbon Tempo 30 Hari).
            </p>
            <div className="flex flex-wrap items-center gap-4 pt-2">
              <Link
                to="/products"
                className="rounded-2xl bg-primary-600 px-6 py-3.5 text-sm font-bold text-white hover:bg-primary-500 transition-colors shadow-lg shadow-primary-600/30 flex items-center space-x-2"
              >
                <span>Belanja Sekarang</span>
                <ArrowRight className="h-4 w-4" />
              </Link>
            </div>
          </motion.div>

          <div className="relative hidden lg:flex justify-center items-center">
            <div className="relative w-80 h-80 rounded-3xl overflow-hidden shadow-2xl border-4 border-white/10 group">
              <img
                src="/logo.png"
                alt="LuckyMart"
                className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
              />
              <div className="absolute inset-0 bg-gradient-to-t from-slate-900/90 via-slate-900/20 to-transparent flex flex-col justify-end p-6">
                <span className="text-xs font-bold text-primary-400 uppercase tracking-wider">Pasokan Utama</span>
                <span className="text-xl font-black text-white">LuckyMart</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Feature Highlights */}
      <section className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 grid grid-cols-1 sm:grid-cols-3 gap-6">
        <div className="flex items-center space-x-4 p-5 rounded-2xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 shadow-sm">
          <div className="p-3 rounded-xl bg-primary-100 dark:bg-primary-950/60 text-primary-600">
            <Truck className="h-6 w-6" />
          </div>
          <div>
            <h4 className="font-bold text-sm text-slate-900 dark:text-white">Pengiriman Same-Day</h4>
            <p className="text-xs text-slate-500">Armada toko siap kirim hari ini</p>
          </div>
        </div>
        <div className="flex items-center space-x-4 p-5 rounded-2xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 shadow-sm">
          <div className="p-3 rounded-xl bg-primary-100 dark:bg-primary-950/60 text-primary-600">
            <ShieldCheck className="h-6 w-6" />
          </div>
          <div>
            <h4 className="font-bold text-sm text-slate-900 dark:text-white">Jaminan Kualitas 100%</h4>
            <p className="text-xs text-slate-500">Barang sembako segar & asli</p>
          </div>
        </div>
        <div className="flex items-center space-x-4 p-5 rounded-2xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 shadow-sm">
          <div className="p-3 rounded-xl bg-primary-100 dark:bg-primary-950/60 text-primary-600">
            <Store className="h-6 w-6" />
          </div>
          <div>
            <h4 className="font-bold text-sm text-slate-900 dark:text-white">Fasilitas B2B Kasbon</h4>
            <p className="text-xs text-slate-500">Tempo 30 Hari Khusus Mitra Warung</p>
          </div>
        </div>
      </section>

      {/* Category Pills Bar */}
      <section className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-bold text-slate-900 dark:text-white flex items-center space-x-2">
            <span>Kategori Sembako Pilihan</span>
          </h2>
          <Link to="/products" className="text-xs font-bold text-primary-600 hover:underline flex items-center space-x-1">
            <span>Lihat Semua Katalog</span>
            <ArrowRight className="h-3 w-3" />
          </Link>
        </div>
        <div className="flex items-center space-x-3 overflow-x-auto pb-2 scrollbar-none">
          {categories.map((c) => (
            <Link
              key={c.id}
              to={`/products?category_id=${c.id}`}
              className="flex-shrink-0 rounded-2xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 px-4 py-3 text-xs font-bold text-slate-700 dark:text-slate-200 hover:border-primary-500 hover:text-primary-600 transition-all shadow-sm flex items-center space-x-2"
            >
              <Store className="h-4 w-4 text-primary-500" />
              <span>{c.name}</span>
            </Link>
          ))}
        </div>
      </section>

      {/* Popular Products Grid */}
      <section className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center space-x-2">
            <Flame className="h-5 w-5 text-rose-500 fill-rose-500" />
            <h2 className="text-xl font-bold text-slate-900 dark:text-white">Produk Pilihan</h2>
          </div>
        </div>

        {loading ? (
          <LoadingSkeleton count={6} />
        ) : (
          <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-4 sm:gap-6">
            {products.map((p) => (
              <ProductCard key={p.id} product={p} />
            ))}
          </div>
        )}
      </section>
    </div>
  );
};
