import React, { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Edit3, Layers, Plus, Trash2 } from 'lucide-react';
import { PageHeader } from '../../components/common/PageHeader';
import { LoadingSkeleton } from '../../components/common/LoadingSkeleton';
import { EmptyState } from '../../components/common/EmptyState';
import { ErrorState } from '../../components/common/ErrorState';
import { Field, Modal, btnDanger, btnPrimary, btnSecondary, inputClass } from '../../components/common/Modal';
import { Card } from '../../components/common/Pagination';
import { productService } from '../../services/apiServices';
import { useAuthStore, MANAGER_ROLES, hasRole } from '../../store/useAuthStore';
import { getErrorMessage } from '../../lib/utils';
import type { Category } from '../../types';

export const CategoryListPage: React.FC = () => {
  const canEdit = hasRole(useAuthStore((s) => s.user), MANAGER_ROLES);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [editing, setEditing] = useState<Category | null>(null);
  const [formOpen, setFormOpen] = useState(false);
  const [deleting, setDeleting] = useState<Category | null>(null);
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [busy, setBusy] = useState(false);

  const load = async () => {
    setLoading(true);
    setError(null);
    try {
      setCategories(await productService.getCategories());
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const open = (c: Category | null) => {
    setEditing(c);
    setName(c?.name || '');
    setDescription(c?.description || '');
    setFormOpen(true);
  };

  const submit = async () => {
    if (name.trim().length < 2) return toast.error('Nama kategori minimal 2 huruf.');
    setBusy(true);
    try {
      if (editing) await productService.updateCategory(editing.id, { name: name.trim(), description: description.trim() });
      else await productService.createCategory({ name: name.trim(), description: description.trim() });
      toast.success(editing ? 'Kategori diperbarui.' : 'Kategori ditambahkan.');
      setFormOpen(false);
      load();
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setBusy(false);
    }
  };

  const remove = async () => {
    if (!deleting) return;
    setBusy(true);
    try {
      await productService.deleteCategory(deleting.id);
      toast.success('Kategori dihapus.');
      setDeleting(null);
      load();
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setBusy(false);
    }
  };

  return (
    <div className="space-y-5">
      <PageHeader title="Kategori" subtitle="Pengelompokan produk di toko online.">
        {canEdit && (
          <button className={btnPrimary} onClick={() => open(null)}>
            <Plus className="h-4 w-4" /> Tambah Kategori
          </button>
        )}
      </PageHeader>

      {loading ? (
        <LoadingSkeleton count={4} />
      ) : error ? (
        <ErrorState message={error} onRetry={load} />
      ) : categories.length === 0 ? (
        <EmptyState title="Belum ada kategori" actionText={canEdit ? 'Tambah Kategori' : undefined} onAction={canEdit ? () => open(null) : undefined} />
      ) : (
        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {categories.map((c) => (
            <Card key={c.id} className="flex items-start gap-3 p-4">
              <div className="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-xl bg-brand-50 text-brand-700 dark:bg-brand-950/40 dark:text-brand-400">
                <Layers className="h-5 w-5" />
              </div>
              <div className="min-w-0 flex-1">
                <div className="font-semibold">{c.name}</div>
                <p className="text-xs text-slate-500">{c.description || 'Tanpa deskripsi'}</p>
              </div>
              {canEdit && (
                <div className="flex gap-1">
                  <button title="Ubah" onClick={() => open(c)} className="rounded-lg p-2 text-slate-600 hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-800">
                    <Edit3 className="h-4 w-4" />
                  </button>
                  <button title="Hapus" onClick={() => setDeleting(c)} className="rounded-lg p-2 text-rose-600 hover:bg-rose-50 dark:hover:bg-rose-950/30">
                    <Trash2 className="h-4 w-4" />
                  </button>
                </div>
              )}
            </Card>
          ))}
        </div>
      )}

      <Modal
        open={formOpen}
        onClose={() => setFormOpen(false)}
        size="sm"
        title={editing ? 'Ubah kategori' : 'Tambah kategori'}
        footer={
          <>
            <button className={btnSecondary} onClick={() => setFormOpen(false)}>Batal</button>
            <button className={btnPrimary} disabled={busy} onClick={submit}>Simpan</button>
          </>
        }
      >
        <div className="grid gap-3">
          <Field label="Nama">
            <input className={inputClass} value={name} onChange={(e) => setName(e.target.value)} autoFocus />
          </Field>
          <Field label="Deskripsi (opsional)">
            <textarea className={inputClass} rows={2} value={description} onChange={(e) => setDescription(e.target.value)} />
          </Field>
        </div>
      </Modal>

      <Modal
        open={!!deleting}
        onClose={() => setDeleting(null)}
        size="sm"
        title={`Hapus kategori ${deleting?.name}?`}
        subtitle="Kategori yang masih dipakai produk tidak bisa dihapus."
        footer={
          <>
            <button className={btnSecondary} onClick={() => setDeleting(null)}>Batal</button>
            <button className={btnDanger} disabled={busy} onClick={remove}>Hapus</button>
          </>
        }
      >
        <p className="text-sm text-slate-600 dark:text-slate-400">Tindakan ini tidak bisa dibatalkan dari dashboard.</p>
      </Modal>
    </div>
  );
};
