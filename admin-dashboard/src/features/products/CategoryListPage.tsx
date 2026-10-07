import React, { useState, useEffect } from 'react';
import { PageHeader } from '../../components/common/PageHeader';
import { LoadingSkeleton } from '../../components/common/LoadingSkeleton';
import { EmptyState } from '../../components/common/EmptyState';
import { ErrorState } from '../../components/common/ErrorState';
import { productService } from '../../services/apiServices';
import { Plus, Search, Layers, CheckCircle2, X, Loader2, Edit3, Trash2 } from 'lucide-react';
import { toast } from 'sonner';

export const CategoryListPage: React.FC = () => {
  const [categories, setCategories] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [search, setSearch] = useState('');

  // Modal Create/Edit state
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingCategory, setEditingCategory] = useState<any | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');

  const fetchCategories = async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await productService.getCategories();
      const items = Array.isArray(res.data)
        ? res.data
        : (Array.isArray(res) ? res : []);
      setCategories(items);
    } catch (err: any) {
      setError(err.message || 'Gagal memuat kategori dari server database.');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCategories();
  }, []);

  const handleOpenAdd = () => {
    setEditingCategory(null);
    setName('');
    setDescription('');
    setIsModalOpen(true);
  };

  const handleOpenEdit = (c: any) => {
    setEditingCategory(c);
    setName(c.name || '');
    setDescription(c.description || '');
    setIsModalOpen(true);
  };

  const handleDeleteCategory = async (id: string, catName: string) => {
    if (window.confirm(`Hapus kategori "${catName}" dari database?`)) {
      try {
        await productService.deleteCategory(id);
        toast.success(`Kategori "${catName}" berhasil dihapus!`);
        fetchCategories();
      } catch (err: any) {
        toast.error(`Gagal menghapus kategori: ${err.message || 'Error'}`);
      }
    }
  };

  const handleSaveCategory = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) {
      toast.error('Nama kategori tidak boleh kosong.');
      return;
    }
    setSubmitting(true);
    try {
      if (editingCategory) {
        await productService.updateCategory(editingCategory.id, { name: name.trim(), description: description.trim() });
        toast.success(`Kategori "${name}" berhasil diperbarui!`);
      } else {
        await productService.createCategory({ name: name.trim(), description: description.trim() });
        toast.success(`Kategori "${name}" berhasil ditambahkan & tayang di Customer Web!`);
      }

      setIsModalOpen(false);
      fetchCategories();
    } catch (err: any) {
      toast.error(err.message || 'Gagal menyimpan kategori.');
    } finally {
      setSubmitting(false);
    }
  };

  const filteredCategories = categories.filter((c) =>
    c.name?.toLowerCase().includes(search.toLowerCase()) ||
    c.slug?.toLowerCase().includes(search.toLowerCase())
  );

  const [selectedIds, setSelectedIds] = useState<string[]>([]);

  const handleSelectAll = () => {
    if (selectedIds.length === filteredCategories.length) {
      setSelectedIds([]);
    } else {
      setSelectedIds(filteredCategories.map((c) => c.id));
    }
  };

  const handleSelectRow = (id: string) => {
    setSelectedIds((prev) =>
      prev.includes(id) ? prev.filter((item) => item !== id) : [...prev, id]
    );
  };

  const handleBulkDelete = async () => {
    if (selectedIds.length === 0) return;
    if (window.confirm(`Apakah Anda yakin ingin menghapus ${selectedIds.length} kategori terpilih sekaligus?`)) {
      const idsToDelete = [...selectedIds];
      setLoading(true);
      
      for (const id of idsToDelete) {
        try {
          await productService.deleteCategory(id);
        } catch (e) {}
      }

      toast.success(`Berhasil menghapus ${idsToDelete.length} kategori terpilih!`);
      setSelectedIds([]);
      fetchCategories();
    }
  };

  return (
    <div className="space-y-6">
      <PageHeader
        title="Kategori Produk Sembako"
        subtitle="Kelola kelompok dan klasifikasi barang jualan di usaha sembako Anda."
      >
        <div className="flex items-center space-x-3">
          {selectedIds.length > 0 && (
            <button
              onClick={handleBulkDelete}
              className="inline-flex items-center space-x-1.5 rounded-xl bg-rose-600 px-4 py-2.5 text-sm font-semibold text-white hover:bg-rose-700 transition-colors shadow-lg shadow-rose-600/20 animate-in fade-in duration-150"
            >
              <Trash2 className="h-4 w-4" />
              <span>Hapus ({selectedIds.length}) Kategori</span>
            </button>
          )}
          <button
            onClick={handleOpenAdd}
            className="inline-flex items-center space-x-2 rounded-xl bg-brand-600 px-4 py-2.5 text-sm font-semibold text-white hover:bg-brand-700 transition-colors shadow-lg shadow-brand-600/20"
          >
            <Plus className="h-4 w-4" />
            <span>Tambah Kategori</span>
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
            placeholder="Cari nama kategori atau slug..."
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
            Total Kategori: <span className="text-brand-600 font-bold">{filteredCategories.length}</span>
          </div>
        </div>
      </div>

      {loading ? (
        <LoadingSkeleton count={4} />
      ) : error ? (
        <ErrorState message={error} onRetry={fetchCategories} />
      ) : filteredCategories.length === 0 ? (
        <EmptyState title="Tidak Ada Kategori Ditemukan" description="Belum ada data kategori yang cocok dengan pencarian Anda." />
      ) : (
        <div className="overflow-x-auto rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 shadow-sm">
          <table className="w-full text-left text-sm text-slate-600 dark:text-slate-400">
            <thead className="bg-slate-50 dark:bg-slate-800/50 text-xs uppercase text-slate-500 font-bold border-b border-slate-200 dark:border-slate-800">
              <tr>
                <th className="px-6 py-4">Nama Kategori</th>
                <th className="px-6 py-4">Slug URL</th>
                <th className="px-6 py-4">Deskripsi</th>
                <th className="px-6 py-4">Status</th>
                <th className="px-4 py-4 w-12 text-center">
                  <input
                    type="checkbox"
                    checked={filteredCategories.length > 0 && selectedIds.length === filteredCategories.length}
                    onChange={handleSelectAll}
                    className="h-4 w-4 rounded border-slate-300 text-brand-600 focus:ring-brand-500 cursor-pointer"
                    title="Pilih Semua"
                  />
                </th>
                <th className="px-6 py-4 text-right">Aksi (CRUD)</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-200 dark:divide-slate-800">
              {filteredCategories.map((c) => (
                <tr
                  key={c.id}
                  className={`hover:bg-slate-50/80 dark:hover:bg-slate-800/30 transition-colors ${
                    selectedIds.includes(c.id) ? 'bg-brand-50/50 dark:bg-brand-950/20' : ''
                  }`}
                >
                  <td className="px-6 py-4">
                    <div className="flex items-center space-x-3">
                      <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-brand-50 dark:bg-brand-950/40 text-brand-600 font-bold">
                        <Layers className="h-5 w-5" />
                      </div>
                      <div className="font-bold text-slate-900 dark:text-white">{c.name}</div>
                    </div>
                  </td>
                  <td className="px-6 py-4 font-mono text-xs text-brand-600 dark:text-brand-400 font-semibold">
                    /{c.slug || c.name?.toLowerCase().replace(/\s+/g, '-')}
                  </td>
                  <td className="px-6 py-4 text-xs text-slate-500 max-w-md truncate">
                    {c.description || 'Kategori utama pasokan barang sembako.'}
                  </td>
                  <td className="px-6 py-4">
                    <span className="inline-flex items-center space-x-1.5 rounded-full bg-emerald-100 dark:bg-emerald-950/50 px-2.5 py-1 text-xs font-semibold text-emerald-700 dark:text-emerald-400">
                      <CheckCircle2 className="h-3.5 w-3.5" />
                      <span>Aktif</span>
                    </span>
                  </td>
                  <td className="px-4 py-4 text-center">
                    <input
                      type="checkbox"
                      checked={selectedIds.includes(c.id)}
                      onChange={() => handleSelectRow(c.id)}
                      className="h-4 w-4 rounded border-slate-300 text-brand-600 focus:ring-brand-500 cursor-pointer"
                    />
                  </td>
                  <td className="px-6 py-4 text-right space-x-2">
                    <button
                      onClick={() => handleOpenEdit(c)}
                      className="inline-flex items-center space-x-1 rounded-lg bg-slate-100 dark:bg-slate-800 px-2.5 py-1 text-xs font-bold text-brand-600 hover:bg-brand-600 hover:text-white transition-colors"
                    >
                      <Edit3 className="h-3.5 w-3.5" />
                      <span>Edit</span>
                    </button>
                    <button
                      onClick={() => handleDeleteCategory(c.id, c.name)}
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
      )}

      {/* Modal Create/Edit Kategori */}
      {isModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 backdrop-blur-sm p-4">
          <div className="relative w-full max-w-md rounded-3xl bg-white dark:bg-slate-900 p-6 shadow-2xl border border-slate-200 dark:border-slate-800 space-y-6">
            <div className="flex items-center justify-between border-b border-slate-200 dark:border-slate-800 pb-4">
              <div className="flex items-center space-x-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-brand-100 text-brand-600 dark:bg-brand-950/60">
                  <Layers className="h-5 w-5" />
                </div>
                <div>
                  <h3 className="text-lg font-bold text-slate-900 dark:text-white">
                    {editingCategory ? 'Edit Data Kategori' : 'Tambah Kategori Baru'}
                  </h3>
                  <p className="text-xs text-slate-500">Perubahan akan langsung ter-update di Customer Web.</p>
                </div>
              </div>
              <button onClick={() => setIsModalOpen(false)} className="rounded-xl p-2 text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800">
                <X className="h-5 w-5" />
              </button>
            </div>

            <form onSubmit={handleSaveCategory} className="space-y-4">
              <div>
                <label className="block text-xs font-bold text-slate-700 dark:text-slate-300 mb-1">
                  Nama Kategori *
                </label>
                <input
                  type="text"
                  required
                  placeholder="Contoh: Minuman & Bumbu Dapur"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-brand-500"
                />
              </div>

              <div>
                <label className="block text-xs font-bold text-slate-700 dark:text-slate-300 mb-1">
                  Deskripsi Singkat
                </label>
                <textarea
                  rows={3}
                  placeholder="Deskripsikan jenis barang sembako dalam kategori ini..."
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-brand-500 resize-none"
                />
              </div>

              <div className="flex items-center justify-end space-x-3 pt-4 border-t border-slate-200 dark:border-slate-800">
                <button
                  type="button"
                  onClick={() => setIsModalOpen(false)}
                  className="rounded-xl px-4 py-2.5 text-xs font-bold text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  disabled={submitting}
                  className="inline-flex items-center space-x-2 rounded-xl bg-brand-600 px-5 py-2.5 text-xs font-bold text-white hover:bg-brand-700 disabled:opacity-50"
                >
                  {submitting && <Loader2 className="h-4 w-4 animate-spin" />}
                  <span>{editingCategory ? 'Simpan Perubahan' : 'Simpan Kategori'}</span>
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
