import React, { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Edit3, ImageOff, Plus, Trash2, X } from 'lucide-react';
import { PageHeader } from '../../components/common/PageHeader';
import { StatusBadge } from '../../components/common/StatusBadge';
import { LoadingSkeleton } from '../../components/common/LoadingSkeleton';
import { EmptyState } from '../../components/common/EmptyState';
import { ErrorState } from '../../components/common/ErrorState';
import { Field, Modal, btnDanger, btnPrimary, btnSecondary, inputClass } from '../../components/common/Modal';
import { Card, Pagination, SearchInput, Table } from '../../components/common/Pagination';
import { productService } from '../../services/apiServices';
import { useList } from '../../hooks/useList';
import { useAuthStore, MANAGER_ROLES, hasRole } from '../../store/useAuthStore';
import { formatCurrency, formatNumber, getErrorMessage } from '../../lib/utils';
import type { Category, Product } from '../../types';

interface UnitForm {
  id?: string;
  unit_name: string;
  retail_price: string;
  wholesale_price: string;
  min_wholesale_qty: string;
  stock: string;
}

const emptyUnit = (): UnitForm => ({ unit_name: 'Pcs', retail_price: '', wholesale_price: '', min_wholesale_qty: '5', stock: '0' });

export const ProductListPage: React.FC = () => {
  const canEdit = hasRole(useAuthStore((s) => s.user), MANAGER_ROLES);
  const list = useList((p) => productService.getProducts(p));
  const [categories, setCategories] = useState<Category[]>([]);
  const [editing, setEditing] = useState<Product | null>(null);
  const [formOpen, setFormOpen] = useState(false);
  const [deleting, setDeleting] = useState<Product | null>(null);
  const [busy, setBusy] = useState(false);
  const [form, setForm] = useState({ name: '', sku: '', category_id: '', description: '', image_url: '', is_active: true });
  const [units, setUnits] = useState<UnitForm[]>([emptyUnit()]);

  useEffect(() => {
    productService.getCategories().then(setCategories).catch(() => toast.error('Gagal memuat kategori'));
  }, []);

  const openCreate = () => {
    setEditing(null);
    setForm({ name: '', sku: '', category_id: categories[0]?.id || '', description: '', image_url: '', is_active: true });
    setUnits([emptyUnit()]);
    setFormOpen(true);
  };

  const openEdit = (p: Product) => {
    setEditing(p);
    setForm({ name: p.name, sku: p.sku, category_id: p.category_id, description: p.description || '', image_url: p.image_url || '', is_active: p.is_active });
    setUnits(
      p.units.map((u) => ({
        id: u.id,
        unit_name: u.unit_name,
        retail_price: String(u.retail_price),
        wholesale_price: String(u.wholesale_price),
        min_wholesale_qty: String(u.min_wholesale_qty),
        stock: String(u.stock_quantity),
      }))
    );
    setFormOpen(true);
  };

  const updateUnit = (i: number, patch: Partial<UnitForm>) => setUnits((us) => us.map((u, idx) => (idx === i ? { ...u, ...patch } : u)));

  const submit = async () => {
    if (form.name.trim().length < 2) return toast.error('Nama produk wajib diisi.');
    if (!form.category_id) return toast.error('Pilih kategori. Buat kategori dulu jika belum ada.');
    for (const u of units) {
      if (!u.unit_name.trim() || !(Number(u.retail_price) > 0)) return toast.error('Setiap satuan wajib punya nama dan harga eceran.');
    }

    setBusy(true);
    try {
      if (editing) {
        await productService.updateProduct(editing.id, {
          name: form.name.trim(),
          sku: form.sku.trim() || undefined,
          category_id: form.category_id,
          description: form.description.trim(),
          image_url: form.image_url.trim(),
          is_active: form.is_active,
          units: units.map((u) => ({
            id: u.id,
            unit_name: u.unit_name.trim(),
            retail_price: Number(u.retail_price),
            wholesale_price: Number(u.wholesale_price) || Number(u.retail_price),
            min_wholesale_qty: Number(u.min_wholesale_qty) || 1,
            stock_quantity: Number(u.stock) || 0,
          })),
        });
        toast.success('Produk diperbarui.');
      } else {
        await productService.createProduct({
          name: form.name.trim(),
          sku: form.sku.trim() || undefined,
          category_id: form.category_id,
          description: form.description.trim(),
          image_url: form.image_url.trim() || undefined,
          units: units.map((u) => ({
            unit_name: u.unit_name.trim(),
            retail_price: Number(u.retail_price),
            wholesale_price: Number(u.wholesale_price) || undefined,
            min_wholesale_qty: Number(u.min_wholesale_qty) || undefined,
            initial_stock: Number(u.stock) || 0,
          })),
        });
        toast.success('Produk ditambahkan dan langsung tampil di toko.');
      }
      setFormOpen(false);
      list.reload();
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setBusy(false);
    }
  };

  const confirmDelete = async () => {
    if (!deleting) return;
    setBusy(true);
    try {
      await productService.deleteProduct(deleting.id);
      toast.success('Produk dihapus dari katalog.');
      setDeleting(null);
      list.reload();
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setBusy(false);
    }
  };

  return (
    <div className="space-y-5">
      <PageHeader title="Produk" subtitle="Katalog barang, satuan, harga eceran dan grosir.">
        {canEdit && (
          <button className={btnPrimary} onClick={openCreate}>
            <Plus className="h-4 w-4" /> Tambah Produk
          </button>
        )}
      </PageHeader>

      <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
        <SearchInput value={list.searchInput} onChange={list.setSearchInput} placeholder="Cari nama / SKU..." />
        <select
          className={`${inputClass} sm:w-60`}
          value={(list.params.category_id as string) || ''}
          onChange={(e) => list.setFilter('category_id', e.target.value)}
        >
          <option value="">Semua kategori</option>
          {categories.map((c) => (
            <option key={c.id} value={c.id}>{c.name}</option>
          ))}
        </select>
      </div>

      {list.loading ? (
        <LoadingSkeleton count={6} />
      ) : list.error ? (
        <ErrorState message={list.error} onRetry={list.reload} />
      ) : list.items.length === 0 ? (
        <EmptyState title="Belum ada produk" actionText={canEdit ? 'Tambah Produk' : undefined} onAction={canEdit ? openCreate : undefined} />
      ) : (
        <Card className="p-2 sm:p-4">
          <Table
            head={
              <tr>
                <th className="px-3 py-3">Produk</th>
                <th className="px-3 py-3">Kategori</th>
                <th className="px-3 py-3">Harga</th>
                <th className="px-3 py-3 text-right">Stok</th>
                <th className="px-3 py-3">Status</th>
                <th className="px-3 py-3" />
              </tr>
            }
          >
            {list.items.map((p) => (
              <tr key={p.id} className="hover:bg-slate-50 dark:hover:bg-slate-800/40">
                <td className="px-3 py-3">
                  <div className="flex items-center gap-3">
                    {p.image_url ? (
                      <img src={p.image_url} alt="" className="h-11 w-11 flex-shrink-0 rounded-xl object-cover bg-slate-100" />
                    ) : (
                      <div className="flex h-11 w-11 flex-shrink-0 items-center justify-center rounded-xl bg-slate-100 dark:bg-slate-800 text-slate-400"><ImageOff className="h-4 w-4" /></div>
                    )}
                    <div>
                      <div className="font-semibold">{p.name}</div>
                      <div className="font-mono text-[11px] text-slate-500">{p.sku}</div>
                    </div>
                  </div>
                </td>
                <td className="px-3 py-3 text-xs">{p.category?.name}</td>
                <td className="px-3 py-3 text-xs">
                  {p.units.map((u) => (
                    <div key={u.id}>
                      <b>{formatCurrency(u.retail_price)}</b> / {u.unit_name}
                      {u.wholesale_price < u.retail_price && <span className="text-slate-500"> · grosir {formatCurrency(u.wholesale_price)} (≥{u.min_wholesale_qty})</span>}
                    </div>
                  ))}
                </td>
                <td className="px-3 py-3 text-right text-xs">
                  {p.units.map((u) => (
                    <div key={u.id} className={u.stock_quantity === 0 ? 'font-semibold text-rose-600' : ''}>{formatNumber(u.stock_quantity)} {u.unit_name}</div>
                  ))}
                </td>
                <td className="px-3 py-3">{p.is_active ? <StatusBadge status="APPROVED" label="Tampil" /> : <StatusBadge status="CANCELLED" label="Disembunyikan" />}</td>
                <td className="px-3 py-3">
                  {canEdit && (
                    <div className="flex justify-end gap-1">
                      <button title="Ubah" onClick={() => openEdit(p)} className="rounded-lg p-2 text-slate-600 hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-800">
                        <Edit3 className="h-4 w-4" />
                      </button>
                      <button title="Hapus" onClick={() => setDeleting(p)} className="rounded-lg p-2 text-rose-600 hover:bg-rose-50 dark:hover:bg-rose-950/30">
                        <Trash2 className="h-4 w-4" />
                      </button>
                    </div>
                  )}
                </td>
              </tr>
            ))}
          </Table>
          <Pagination meta={list.meta} onPage={list.setPage} />
        </Card>
      )}

      <Modal
        open={formOpen}
        onClose={() => setFormOpen(false)}
        size="lg"
        title={editing ? `Ubah ${editing.name}` : 'Tambah produk'}
        footer={
          <>
            <button className={btnSecondary} onClick={() => setFormOpen(false)}>Batal</button>
            <button className={btnPrimary} disabled={busy} onClick={submit}>{busy ? 'Menyimpan...' : 'Simpan'}</button>
          </>
        }
      >
        <div className="grid gap-3 sm:grid-cols-2">
          <div className="sm:col-span-2">
            <Field label="Nama produk">
              <input className={inputClass} value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} placeholder="Misal: Beras Pandan Wangi 5kg" />
            </Field>
          </div>
          <Field label="Kategori">
            <select className={inputClass} value={form.category_id} onChange={(e) => setForm({ ...form, category_id: e.target.value })}>
              <option value="">— Pilih —</option>
              {categories.map((c) => (
                <option key={c.id} value={c.id}>{c.name}</option>
              ))}
            </select>
          </Field>
          <Field label="SKU / kode barang" hint={editing ? undefined : 'Kosongkan untuk dibuat otomatis'}>
            <input className={inputClass} value={form.sku} onChange={(e) => setForm({ ...form, sku: e.target.value.toUpperCase() })} />
          </Field>
          <div className="sm:col-span-2">
            <Field label="Link foto produk" hint="Sementara berupa link gambar. Fitur unggah foto langsung menyusul.">
              <input className={inputClass} value={form.image_url} onChange={(e) => setForm({ ...form, image_url: e.target.value })} placeholder="https://..." />
            </Field>
          </div>
          <div className="sm:col-span-2">
            <Field label="Deskripsi">
              <textarea className={inputClass} rows={2} value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} />
            </Field>
          </div>
        </div>

        <div className="mt-5">
          <div className="mb-2 flex items-center justify-between">
            <span className="text-sm font-bold">Satuan & harga</span>
            {!editing && (
              <button className="text-xs font-semibold text-brand-600" onClick={() => setUnits((u) => [...u, { ...emptyUnit(), unit_name: '' }])}>
                + Tambah satuan
              </button>
            )}
          </div>
          <div className="space-y-3">
            {units.map((u, i) => (
              <div key={u.id || i} className="rounded-xl border border-slate-200 dark:border-slate-800 p-3">
                <div className="grid grid-cols-2 gap-2 sm:grid-cols-5">
                  <Field label="Satuan">
                    <input className={inputClass} value={u.unit_name} onChange={(e) => updateUnit(i, { unit_name: e.target.value })} placeholder="Pcs / Dus / Kg" />
                  </Field>
                  <Field label="Harga eceran">
                    <input type="number" min={0} className={inputClass} value={u.retail_price} onChange={(e) => updateUnit(i, { retail_price: e.target.value })} />
                  </Field>
                  <Field label="Harga grosir">
                    <input type="number" min={0} className={inputClass} value={u.wholesale_price} onChange={(e) => updateUnit(i, { wholesale_price: e.target.value })} placeholder="= eceran" />
                  </Field>
                  <Field label="Min. grosir">
                    <input type="number" min={1} className={inputClass} value={u.min_wholesale_qty} onChange={(e) => updateUnit(i, { min_wholesale_qty: e.target.value })} />
                  </Field>
                  <Field label={editing ? 'Stok' : 'Stok awal'}>
                    <input type="number" min={0} className={inputClass} value={u.stock} onChange={(e) => updateUnit(i, { stock: e.target.value })} />
                  </Field>
                </div>
                {!editing && units.length > 1 && (
                  <button className="mt-2 inline-flex items-center gap-1 text-xs text-rose-600" onClick={() => setUnits((us) => us.filter((_, idx) => idx !== i))}>
                    <X className="h-3 w-3" /> Hapus satuan
                  </button>
                )}
              </div>
            ))}
          </div>
          <p className="mt-2 text-[11px] text-slate-500">
            Harga grosir berlaku untuk mitra toko (B2B) yang membeli minimal sejumlah "Min. grosir".
            {editing && ' Perubahan stok di sini dicatat sebagai koreksi di Mutasi Stok.'}
          </p>
        </div>

        {editing && (
          <label className="mt-4 flex items-center gap-2 text-sm">
            <input type="checkbox" checked={form.is_active} onChange={(e) => setForm({ ...form, is_active: e.target.checked })} />
            Tampilkan di toko online
          </label>
        )}
      </Modal>

      <Modal
        open={!!deleting}
        onClose={() => setDeleting(null)}
        size="sm"
        title={`Hapus ${deleting?.name}?`}
        subtitle="Produk hilang dari toko dan keranjang pelanggan. Riwayat pesanan lama tetap tersimpan. Untuk menyembunyikan sementara, pakai Ubah → matikan 'Tampilkan di toko online'."
        footer={
          <>
            <button className={btnSecondary} onClick={() => setDeleting(null)}>Batal</button>
            <button className={btnDanger} disabled={busy} onClick={confirmDelete}>Hapus</button>
          </>
        }
      >
        <p className="text-sm text-slate-600 dark:text-slate-400">Tindakan ini tidak bisa dibatalkan dari dashboard.</p>
      </Modal>
    </div>
  );
};
