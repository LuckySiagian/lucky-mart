import React, { useState, useEffect } from 'react';
import { PageHeader } from '../../components/common/PageHeader';
import { LoadingSkeleton } from '../../components/common/LoadingSkeleton';
import { EmptyState } from '../../components/common/EmptyState';
import { ErrorState } from '../../components/common/ErrorState';
import { productService } from '../../services/apiServices';
import { formatCurrency } from '../../lib/utils';
import { Plus, Search, Package, X, Loader2, Edit3, Trash2 } from 'lucide-react';
import { toast } from 'sonner';

export const ProductListPage: React.FC = () => {
  const [products, setProducts] = useState<any[]>([]);
  const [categories, setCategories] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [search, setSearch] = useState('');

  // Modal Create/Edit State
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingProduct, setEditingProduct] = useState<any | null>(null);
  const [submitting, setSubmitting] = useState(false);

  // Form Fields State
  const [formData, setFormData] = useState({
    name: '',
    sku: '',
    category_id: '',
    description: '',
    unit_name: 'Pcs',
    retail_price: '',
    wholesale_price: '',
    initial_stock: '100',
    image_url: '',
  });

  const fetchProducts = async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await productService.getProducts(1, 100);
      const items = Array.isArray(res.data)
        ? res.data
        : (res.data?.items || (Array.isArray(res) ? res : []));
      setProducts(items);
    } catch (err: any) {
      setError(err.message || 'Gagal memuat produk dari server database.');
    } finally {
      setLoading(false);
    }
  };

  const fetchCategories = async () => {
    try {
      const res = await productService.getCategories();
      const items = Array.isArray(res.data)
        ? res.data
        : (Array.isArray(res) ? res : []);

      setCategories(items);
      if (items.length > 0 && !formData.category_id) {
        setFormData((prev) => ({ ...prev, category_id: items[0].id }));
      }
    } catch (e) {}
  };

  useEffect(() => {
    fetchProducts();
    fetchCategories();
  }, []);

  const handleOpenAdd = () => {
    fetchCategories();
    setEditingProduct(null);
    const defaultSKU = `SKU-${Date.now().toString().slice(-6)}`;
    setFormData({
      name: '',
      sku: defaultSKU,
      category_id: categories.length > 0 ? categories[0].id : '',
      description: '',
      unit_name: 'Pcs',
      retail_price: '',
      wholesale_price: '',
      initial_stock: '100',
      image_url: '',
    });
    setIsModalOpen(true);
  };

  const handleOpenEdit = (p: any) => {
    setEditingProduct(p);
    const unit = p.product_units?.[0] || p.units?.[0];
    const image = p.images?.[0]?.image_url || '';
    setFormData({
      name: p.name || '',
      sku: p.sku || '',
      category_id: p.category_id || (categories.length > 0 ? categories[0].id : ''),
      description: p.description || '',
      unit_name: unit?.unit_name || 'Pcs',
      retail_price: unit?.retail_price ? unit.retail_price.toString() : '',
      wholesale_price: unit?.wholesale_price ? unit.wholesale_price.toString() : '',
      initial_stock: unit?.stock_quantity ? unit.stock_quantity.toString() : '100',
      image_url: image,
    });
    setIsModalOpen(true);
  };

  const handleDeleteProduct = async (id: string, name: string) => {
    if (window.confirm(`Apakah Anda yakin ingin menghapus produk "${name}" dari katalog?`)) {
      try {
        await productService.deleteProduct(id);
        toast.success(`Produk "${name}" berhasil dihapus!`);
        fetchProducts();
      } catch (e: any) {
        toast.error(`Gagal menghapus produk: ${e.message || 'Error'}`);
      }
      setSelectedIds((prev) => prev.filter((item) => item !== id));
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.name.trim() || !formData.sku.trim() || !formData.retail_price) {
      toast.error('Mohon isi nama produk, SKU, dan harga eceran.');
      return;
    }

    setSubmitting(true);
    try {
      if (editingProduct) {
        // UPDATE existing product
        const payload = {
          name: formData.name.trim(),
          sku: formData.sku.trim(),
          description: formData.description.trim(),
          image_url: formData.image_url.trim(),
          retail_price: parseFloat(formData.retail_price),
          wholesale_price: parseFloat(formData.wholesale_price || formData.retail_price),
          stock_quantity: parseInt(formData.initial_stock || '0', 10),
        };

        await productService.updateProduct(editingProduct.id, payload);
        toast.success(`Produk "${formData.name}" berhasil diperbarui!`);
      } else {
        // CREATE new product
        const payload = {
          name: formData.name.trim(),
          sku: formData.sku.trim(),
          category_id: formData.category_id || (categories[0]?.id || 'cat-1'),
          description: formData.description.trim(),
          image_url: formData.image_url.trim(),
          units: [
            {
              unit_id: '00000000-0000-0000-0000-000000000001',
              unit_name: formData.unit_name.trim() || 'Pcs',
              conversion_factor: 1,
              retail_price: parseFloat(formData.retail_price),
              wholesale_price: parseFloat(formData.wholesale_price || formData.retail_price),
              min_wholesale_qty: 5,
              initial_stock: parseInt(formData.initial_stock || '0', 10),
            },
          ],
        };

        await productService.createProduct(payload);
        toast.success(`Produk "${formData.name}" berhasil ditambahkan ke Database & Customer Web!`);
      }
      setIsModalOpen(false);
      await fetchProducts();
    } catch (err: any) {
      const msg = err.response?.data?.message || err.message || 'Gagal menyimpan produk ke database.';
      toast.error(`Gagal: ${msg}`);
    } finally {
      setSubmitting(false);
    }
  };

  const filteredProducts = Array.isArray(products)
    ? products.filter((p) =>
        p.name?.toLowerCase().includes(search.toLowerCase()) || p.sku?.toLowerCase().includes(search.toLowerCase())
      )
    : [];

  const [selectedIds, setSelectedIds] = useState<string[]>([]);

  const handleSelectAll = () => {
    if (selectedIds.length === filteredProducts.length) {
      setSelectedIds([]);
    } else {
      setSelectedIds(filteredProducts.map((p) => p.id));
    }
  };

  const handleSelectRow = (id: string) => {
    setSelectedIds((prev) =>
      prev.includes(id) ? prev.filter((item) => item !== id) : [...prev, id]
    );
  };

  const handleBulkDelete = async () => {
    if (selectedIds.length === 0) return;
    if (window.confirm(`Apakah Anda yakin ingin menghapus ${selectedIds.length} produk terpilih sekaligus dari katalog?`)) {
      const idsToDelete = [...selectedIds];
      setLoading(true);
      
      for (const id of idsToDelete) {
        try {
          await productService.deleteProduct(id);
        } catch (e) {}
      }

      toast.success(`Berhasil menghapus ${idsToDelete.length} produk terpilih!`);
      setSelectedIds([]);
      await fetchProducts();
    }
  };

  return (
    <div className="space-y-6">
      <PageHeader
        title="Katalog Produk Sembako"
        subtitle="Kelola seluruh barang, variasi satuan, harga eceran & grosir."
      >
        <div className="flex items-center space-x-3">
          {selectedIds.length > 0 && (
            <button
              onClick={handleBulkDelete}
              className="inline-flex items-center space-x-1.5 rounded-xl bg-rose-600 px-4 py-2.5 text-sm font-semibold text-white hover:bg-rose-700 transition-colors shadow-lg shadow-rose-600/20 animate-in fade-in duration-150"
            >
              <Trash2 className="h-4 w-4" />
              <span>Hapus ({selectedIds.length}) Produk</span>
            </button>
          )}
          <button
            onClick={handleOpenAdd}
            className="inline-flex items-center space-x-2 rounded-xl bg-brand-600 px-4 py-2.5 text-sm font-semibold text-white hover:bg-brand-700 transition-colors shadow-lg shadow-brand-600/20"
          >
            <Plus className="h-4 w-4" />
            <span>Tambah Produk</span>
          </button>
        </div>
      </PageHeader>

      {/* Filter & Search Bar */}
      <div className="flex flex-col sm:flex-row items-center justify-between gap-4 bg-white dark:bg-slate-900 p-4 rounded-2xl border border-slate-200 dark:border-slate-800">
        <div className="relative w-full sm:w-80">
          <Search className="absolute left-3.5 top-3 h-4 w-4 text-slate-400" />
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Cari nama produk atau SKU..."
            className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 pl-10 pr-4 py-2 text-sm text-slate-900 dark:text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-brand-500"
          />
        </div>
        <div className="flex items-center space-x-4 text-xs text-slate-500 font-semibold">
          {selectedIds.length > 0 && (
            <span className="text-rose-600 font-bold bg-rose-50 dark:bg-rose-950/40 px-2.5 py-1 rounded-lg">
              {selectedIds.length} Terpilih
            </span>
          )}
          <div>
            Total Produk: <span className="text-brand-600 font-bold">{filteredProducts.length}</span>
          </div>
        </div>
      </div>

      {loading ? (
        <LoadingSkeleton count={4} />
      ) : error ? (
        <ErrorState message={error} onRetry={fetchProducts} />
      ) : filteredProducts.length === 0 ? (
        <EmptyState title="Tidak ada produk ditemukan" description="Cobalah mengubah kriteria pencarian atau tambahkan produk baru." />
      ) : (
        <div className="overflow-x-auto rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 shadow-sm">
          <table className="w-full text-left text-sm text-slate-600 dark:text-slate-400">
            <thead className="bg-slate-50 dark:bg-slate-800/50 text-xs uppercase text-slate-500 font-bold border-b border-slate-200 dark:border-slate-800">
              <tr>
                <th className="px-6 py-4">Produk / SKU</th>
                <th className="px-6 py-4">Kategori</th>
                <th className="px-6 py-4">Satuan & Harga</th>
                <th className="px-6 py-4">Total Stok</th>
                <th className="px-6 py-4">Status</th>
                <th className="px-4 py-4 w-12 text-center">
                  <input
                    type="checkbox"
                    checked={filteredProducts.length > 0 && selectedIds.length === filteredProducts.length}
                    onChange={handleSelectAll}
                    className="h-4 w-4 rounded border-slate-300 text-brand-600 focus:ring-brand-500 cursor-pointer"
                    title="Pilih Semua"
                  />
                </th>
                <th className="px-6 py-4 text-right">Aksi (CRUD)</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-200 dark:divide-slate-800">
              {filteredProducts.map((p) => {
                const units = (p.units && p.units.length > 0) ? p.units : (p.product_units && p.product_units.length > 0) ? p.product_units : [];
                const primaryUnit = units[0];
                const totalStock = units.reduce((acc: number, u: any) => acc + (u.stock_quantity || 0), 0);

                return (
                  <tr
                    key={p.id}
                    className={`hover:bg-slate-50/80 dark:hover:bg-slate-800/30 transition-colors ${
                      selectedIds.includes(p.id) ? 'bg-brand-50/50 dark:bg-brand-950/20' : ''
                    }`}
                  >
                    <td className="px-6 py-4">
                      <div className="flex items-center space-x-3">
                        <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-slate-100 dark:bg-slate-800 text-slate-500 overflow-hidden">
                          {p.images && p.images.length > 0 ? (
                            <img src={p.images[0].image_url} alt={p.name} className="h-10 w-10 object-cover" />
                          ) : (
                            <Package className="h-5 w-5" />
                          )}
                        </div>
                        <div>
                          <div className="font-bold text-slate-900 dark:text-white">{p.name}</div>
                          <div className="text-xs text-slate-400">SKU: {p.sku}</div>
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4 font-medium text-slate-700 dark:text-slate-300">
                      {p.category?.name || 'Sembako Utama'}
                    </td>
                    <td className="px-6 py-4">
                      {primaryUnit ? (
                        <div className="space-y-1">
                          <div className="text-xs font-semibold text-slate-900 dark:text-white">
                            {primaryUnit.unit_name}: {formatCurrency(primaryUnit.retail_price)}
                          </div>
                          <div className="text-[11px] text-brand-600 font-medium">
                            Grosir: {formatCurrency(primaryUnit.wholesale_price)} (Min {primaryUnit.min_wholesale_qty || 5})
                          </div>
                        </div>
                      ) : (
                        <span className="text-xs text-slate-400">-</span>
                      )}
                    </td>
                    <td className="px-6 py-4">
                      <span className="font-bold text-slate-900 dark:text-white">
                        {totalStock}
                      </span>
                    </td>
                    <td className="px-6 py-4">
                      <span className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-semibold ${p.is_active !== false ? 'bg-emerald-100 text-emerald-800 dark:bg-emerald-950/60 dark:text-emerald-400' : 'bg-rose-100 text-rose-800'}`}>
                        {p.is_active !== false ? 'Aktif' : 'Non-Aktif'}
                      </span>
                    </td>
                    <td className="px-4 py-4 text-center">
                      <input
                        type="checkbox"
                        checked={selectedIds.includes(p.id)}
                        onChange={() => handleSelectRow(p.id)}
                        className="h-4 w-4 rounded border-slate-300 text-brand-600 focus:ring-brand-500 cursor-pointer"
                      />
                    </td>
                    <td className="px-6 py-4 text-right space-x-2">
                      <button
                        onClick={() => handleOpenEdit(p)}
                        className="inline-flex items-center space-x-1 rounded-lg bg-slate-100 dark:bg-slate-800 px-2.5 py-1 text-xs font-bold text-brand-600 hover:bg-brand-600 hover:text-white transition-colors"
                      >
                        <Edit3 className="h-3.5 w-3.5" />
                        <span>Edit</span>
                      </button>
                      <button
                        onClick={() => handleDeleteProduct(p.id, p.name)}
                        className="inline-flex items-center space-x-1 rounded-lg bg-rose-100 text-rose-700 dark:bg-rose-950/50 dark:text-rose-400 px-2.5 py-1 text-xs font-bold hover:bg-rose-200"
                      >
                        <Trash2 className="h-3.5 w-3.5" />
                        <span>Hapus</span>
                      </button>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}

      {/* Modal Form Create/Edit Produk */}
      {isModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 backdrop-blur-sm p-4 overflow-y-auto">
          <div className="relative w-full max-w-xl rounded-3xl bg-white dark:bg-slate-900 p-6 shadow-2xl border border-slate-200 dark:border-slate-800 space-y-6 animate-in fade-in zoom-in duration-200">
            <div className="flex items-center justify-between border-b border-slate-200 dark:border-slate-800 pb-4">
              <div className="flex items-center space-x-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-brand-100 text-brand-600 dark:bg-brand-950/60">
                  <Package className="h-5 w-5" />
                </div>
                <div>
                  <h3 className="text-lg font-bold text-slate-900 dark:text-white">
                    {editingProduct ? 'Edit Data Produk' : 'Tambah Produk Baru'}
                  </h3>
                  <p className="text-xs text-slate-500">Perubahan akan langsung ter-update di Customer Web.</p>
                </div>
              </div>
              <button onClick={() => setIsModalOpen(false)} className="rounded-xl p-2 text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800">
                <X className="h-5 w-5" />
              </button>
            </div>

            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-bold text-slate-700 dark:text-slate-300 mb-1">Nama Produk *</label>
                  <input
                    type="text"
                    required
                    placeholder="Contoh: Beras Rojo Lele 5kg"
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-brand-500"
                  />
                </div>

                <div>
                  <label className="block text-xs font-bold text-slate-700 dark:text-slate-300 mb-1">Kode SKU *</label>
                  <input
                    type="text"
                    required
                    placeholder="Contoh: BRS-001"
                    value={formData.sku}
                    onChange={(e) => setFormData({ ...formData, sku: e.target.value })}
                    className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-brand-500 font-mono"
                  />
                </div>
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-bold text-slate-700 dark:text-slate-300 mb-1">Kategori Sembako *</label>
                  <select
                    required
                    value={formData.category_id}
                    onChange={(e) => setFormData({ ...formData, category_id: e.target.value })}
                    className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-brand-500"
                  >
                    {categories.map((c) => (
                      <option key={c.id} value={c.id}>
                        {c.name}
                      </option>
                    ))}
                  </select>
                </div>

                <div>
                  <label className="block text-xs font-bold text-slate-700 dark:text-slate-300 mb-1">Nama Satuan & Kemasan *</label>
                  <input
                    type="text"
                    required
                    placeholder="Contoh: Karung 5kg / Pouch 2L"
                    value={formData.unit_name}
                    onChange={(e) => setFormData({ ...formData, unit_name: e.target.value })}
                    className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-brand-500"
                  />
                </div>
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                <div>
                  <label className="block text-xs font-bold text-slate-700 dark:text-slate-300 mb-1">Harga Eceran (Rp) *</label>
                  <input
                    type="text"
                    inputMode="numeric"
                    required
                    placeholder="75000"
                    value={formData.retail_price}
                    onChange={(e) => setFormData({ ...formData, retail_price: e.target.value.replace(/\D/g, '') })}
                    onKeyDown={(e) => {
                      if (!/[0-9]/.test(e.key) && !['Backspace', 'Delete', 'ArrowLeft', 'ArrowRight', 'Tab', 'Home', 'End'].includes(e.key) && !e.ctrlKey && !e.metaKey) {
                        e.preventDefault();
                      }
                    }}
                    className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-brand-500 font-mono"
                  />
                </div>

                <div>
                  <label className="block text-xs font-bold text-slate-700 dark:text-slate-300 mb-1">Harga Grosir (Rp)</label>
                  <input
                    type="text"
                    inputMode="numeric"
                    placeholder="70000"
                    value={formData.wholesale_price}
                    onChange={(e) => setFormData({ ...formData, wholesale_price: e.target.value.replace(/\D/g, '') })}
                    onKeyDown={(e) => {
                      if (!/[0-9]/.test(e.key) && !['Backspace', 'Delete', 'ArrowLeft', 'ArrowRight', 'Tab', 'Home', 'End'].includes(e.key) && !e.ctrlKey && !e.metaKey) {
                        e.preventDefault();
                      }
                    }}
                    className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-brand-500 font-mono"
                  />
                </div>

                <div>
                  <label className="block text-xs font-bold text-slate-700 dark:text-slate-300 mb-1">Stok Fisik</label>
                  <input
                    type="text"
                    inputMode="numeric"
                    placeholder="100"
                    value={formData.initial_stock}
                    onChange={(e) => setFormData({ ...formData, initial_stock: e.target.value.replace(/\D/g, '') })}
                    onKeyDown={(e) => {
                      if (!/[0-9]/.test(e.key) && !['Backspace', 'Delete', 'ArrowLeft', 'ArrowRight', 'Tab', 'Home', 'End'].includes(e.key) && !e.ctrlKey && !e.metaKey) {
                        e.preventDefault();
                      }
                    }}
                    className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-brand-500 font-mono"
                  />
                </div>
              </div>

              {/* Upload File Foto atau Input URL */}
              <div className="space-y-2">
                <label className="block text-xs font-bold text-slate-700 dark:text-slate-300 mb-1">
                  Foto Gambar Produk (File JPG, PNG, WEBP / URL)
                </label>
                
                <div className="flex flex-col sm:flex-row items-center gap-3">
                  {/* Preview Thumbnail */}
                  <div className="h-20 w-20 rounded-2xl bg-slate-100 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 overflow-hidden flex items-center justify-center flex-shrink-0 relative group">
                    {formData.image_url ? (
                      <img src={formData.image_url} alt="Preview" className="h-full w-full object-cover" />
                    ) : (
                      <Package className="h-8 w-8 text-slate-400" />
                    )}
                  </div>

                  <div className="flex-1 space-y-2 w-full">
                    {/* File Upload Input */}
                    <div className="relative">
                      <input
                        type="file"
                        id="product-image-upload"
                        accept="image/png, image/jpeg, image/jpg, image/webp"
                        onChange={(e) => {
                          const file = e.target.files?.[0];
                          if (file) {
                            if (file.size > 5 * 1024 * 1024) {
                              toast.error('Ukuran file foto maksimal 5MB.');
                              return;
                            }
                            const reader = new FileReader();
                            reader.onload = () => {
                              if (typeof reader.result === 'string') {
                                setFormData((prev) => ({ ...prev, image_url: reader.result as string }));
                                toast.success(`File foto "${file.name}" berhasil diunggah!`);
                              }
                            };
                            reader.readAsDataURL(file);
                          }
                        }}
                        className="hidden"
                      />
                      <label
                        htmlFor="product-image-upload"
                        className="inline-flex items-center justify-center space-x-2 w-full rounded-xl bg-brand-50 dark:bg-brand-950/40 border border-dashed border-brand-300 dark:border-brand-800 px-4 py-2.5 text-xs font-bold text-brand-600 hover:bg-brand-100 transition-colors cursor-pointer"
                      >
                        <Plus className="h-4 w-4" />
                        <span>Pilih & Upload File Foto (JPG / PNG / WEBP)</span>
                      </label>
                    </div>

                    {/* Fallback URL Input */}
                    <input
                      type="url"
                      placeholder="Atau tempel URL gambar web di sini..."
                      value={formData.image_url}
                      onChange={(e) => setFormData({ ...formData, image_url: e.target.value })}
                      className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-1.5 text-xs text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-brand-500 font-mono"
                    />
                  </div>
                </div>
              </div>

              <div>
                <label className="block text-xs font-bold text-slate-700 dark:text-slate-300 mb-1">Deskripsi Produk</label>
                <textarea
                  rows={2}
                  placeholder="Deskripsi singkat keunggulan produk sembako..."
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-brand-500 resize-none"
                />
              </div>

              <div className="flex items-center justify-end space-x-3 pt-4 border-t border-slate-200 dark:border-slate-800">
                <button
                  type="button"
                  onClick={() => setIsModalOpen(false)}
                  className="rounded-xl px-4 py-2.5 text-xs font-bold text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  disabled={submitting}
                  className="inline-flex items-center space-x-2 rounded-xl bg-brand-600 px-5 py-2.5 text-xs font-bold text-white hover:bg-brand-700 disabled:opacity-50 transition-colors shadow-lg shadow-brand-600/20"
                >
                  {submitting && <Loader2 className="h-4 w-4 animate-spin" />}
                  <span>{editingProduct ? 'Simpan Perubahan' : 'Simpan & Tayangkan Produk'}</span>
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
