import React, { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { AlertTriangle, ArrowDownToLine, ArrowUpFromLine, Bell, PencilRuler } from 'lucide-react';
import { PageHeader } from '../../components/common/PageHeader';
import { StatusBadge } from '../../components/common/StatusBadge';
import { LoadingSkeleton } from '../../components/common/LoadingSkeleton';
import { EmptyState } from '../../components/common/EmptyState';
import { ErrorState } from '../../components/common/ErrorState';
import { Field, Modal, btnPrimary, btnSecondary, inputClass } from '../../components/common/Modal';
import { Card, FilterTabs, Pagination, SearchInput, Table } from '../../components/common/Pagination';
import { inventoryService } from '../../services/apiServices';
import { useList } from '../../hooks/useList';
import { formatNumber, getErrorMessage } from '../../lib/utils';
import type { InventoryItem } from '../../types';

type Mode = 'STOCK_IN' | 'STOCK_OUT' | 'ADJUSTMENT' | 'MIN' | null;

const MODE_TEXT: Record<Exclude<Mode, null>, { title: string; hint: string }> = {
  STOCK_IN: { title: 'Barang masuk', hint: 'Jumlah barang yang datang dari supplier / gudang.' },
  STOCK_OUT: { title: 'Barang keluar', hint: 'Barang rusak, kedaluwarsa, atau dipakai sendiri (bukan penjualan).' },
  ADJUSTMENT: { title: 'Koreksi stok (stok opname)', hint: 'Isi jumlah stok yang benar menurut hitungan fisik.' },
  MIN: { title: 'Batas peringatan stok', hint: 'Peringatan muncul jika stok sama atau di bawah angka ini.' },
};

export const InventoryPage: React.FC = () => {
  const [view, setView] = useState<'all' | 'low'>('all');
  const list = useList((p) => inventoryService.getInventory(p));
  const [lowStock, setLowStock] = useState<InventoryItem[]>([]);
  const [lowLoading, setLowLoading] = useState(false);
  const [target, setTarget] = useState<InventoryItem | null>(null);
  const [mode, setMode] = useState<Mode>(null);
  const [qty, setQty] = useState('');
  const [notes, setNotes] = useState('');
  const [reference, setReference] = useState('');
  const [busy, setBusy] = useState(false);

  const loadLow = async () => {
    setLowLoading(true);
    try {
      setLowStock(await inventoryService.getLowStock());
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setLowLoading(false);
    }
  };

  useEffect(() => {
    loadLow();
  }, []);

  const open = (m: Mode, item: InventoryItem) => {
    setTarget(item);
    setMode(m);
    setQty(m === 'ADJUSTMENT' ? String(item.current_stock) : m === 'MIN' ? String(item.min_stock_alert) : '');
    setNotes('');
    setReference('');
  };

  const submit = async () => {
    if (!target || !mode) return;
    const n = Number(qty);
    if (!Number.isInteger(n) || n < 0) {
      toast.error('Masukkan angka bulat yang valid.');
      return;
    }
    setBusy(true);
    try {
      if (mode === 'MIN') {
        await inventoryService.updateMinStock(target.id, n);
        toast.success('Batas peringatan diperbarui.');
      } else {
        const delta = mode === 'ADJUSTMENT' ? n - target.current_stock : n;
        if (delta === 0) {
          toast.error(mode === 'ADJUSTMENT' ? 'Stok sudah sama dengan angka tersebut.' : 'Jumlah harus lebih dari 0.');
          setBusy(false);
          return;
        }
        if (notes.trim().length < 3) {
          toast.error('Tulis keterangan (minimal 3 huruf).');
          setBusy(false);
          return;
        }
        await inventoryService.adjustStock({
          product_unit_id: target.product_unit_id,
          adjustment_type: mode,
          quantity_changed: delta,
          reference_number: reference.trim() || undefined,
          notes: notes.trim(),
        });
        toast.success(`Stok ${target.product_name} diperbarui.`);
      }
      setMode(null);
      list.reload();
      loadLow();
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setBusy(false);
    }
  };

  const rows = view === 'low' ? lowStock : list.items;
  const loading = view === 'low' ? lowLoading : list.loading;

  return (
    <div className="space-y-5">
      <PageHeader title="Stok Barang" subtitle="Stok yang tampil di toko online. Penjualan dan pembatalan mengubah stok otomatis." />

      {lowStock.length > 0 && view === 'all' && (
        <button
          onClick={() => setView('low')}
          className="flex w-full items-center gap-3 rounded-2xl border border-amber-200 bg-amber-50 p-4 text-left text-sm text-amber-900 dark:border-amber-900 dark:bg-amber-950/30 dark:text-amber-200"
        >
          <AlertTriangle className="h-5 w-5 flex-shrink-0" />
          <span><b>{lowStock.length} barang</b> stoknya menipis. Klik untuk melihat.</span>
        </button>
      )}

      <div className="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
        <FilterTabs
          value={view}
          onChange={(v) => setView(v as 'all' | 'low')}
          options={[
            { value: 'all', label: 'Semua barang' },
            { value: 'low', label: `Stok menipis (${lowStock.length})` },
          ]}
        />
        {view === 'all' && <SearchInput value={list.searchInput} onChange={list.setSearchInput} placeholder="Cari nama / SKU..." />}
      </div>

      {loading ? (
        <LoadingSkeleton count={6} />
      ) : view === 'all' && list.error ? (
        <ErrorState message={list.error} onRetry={list.reload} />
      ) : rows.length === 0 ? (
        <EmptyState title={view === 'low' ? 'Tidak ada stok menipis' : 'Belum ada barang'} description={view === 'low' ? 'Semua stok aman.' : 'Tambahkan produk di menu Produk.'} />
      ) : (
        <Card className="p-2 sm:p-4">
          <Table
            head={
              <tr>
                <th className="px-3 py-3">Barang</th>
                <th className="px-3 py-3">Satuan</th>
                <th className="px-3 py-3 text-right">Stok</th>
                <th className="px-3 py-3 text-right">Batas</th>
                <th className="px-3 py-3">Status</th>
                <th className="px-3 py-3" />
              </tr>
            }
          >
            {rows.map((item) => (
              <tr key={item.id} className="hover:bg-slate-50 dark:hover:bg-slate-800/40">
                <td className="px-3 py-3">
                  <div className="font-semibold">{item.product_name}</div>
                  <div className="text-xs text-slate-500">{item.sku}{item.category_name ? ` · ${item.category_name}` : ''}</div>
                </td>
                <td className="px-3 py-3 text-xs">{item.unit_name}</td>
                <td className={`px-3 py-3 text-right text-base font-bold ${item.is_low_stock ? 'text-rose-600' : ''}`}>{formatNumber(item.current_stock)}</td>
                <td className="px-3 py-3 text-right text-xs text-slate-500">{formatNumber(item.min_stock_alert)}</td>
                <td className="px-3 py-3">
                  {item.current_stock === 0 ? (
                    <StatusBadge status="CANCELLED" label="Habis" />
                  ) : item.is_low_stock ? (
                    <StatusBadge status="PENDING" label="Menipis" />
                  ) : (
                    <StatusBadge status="APPROVED" label="Aman" />
                  )}
                </td>
                <td className="px-3 py-3">
                  <div className="flex justify-end gap-1">
                    <button title="Barang masuk" onClick={() => open('STOCK_IN', item)} className="rounded-lg p-2 text-emerald-600 hover:bg-emerald-50 dark:hover:bg-emerald-950/40">
                      <ArrowDownToLine className="h-4 w-4" />
                    </button>
                    <button title="Barang keluar" onClick={() => open('STOCK_OUT', item)} className="rounded-lg p-2 text-rose-600 hover:bg-rose-50 dark:hover:bg-rose-950/30">
                      <ArrowUpFromLine className="h-4 w-4" />
                    </button>
                    <button title="Koreksi stok" onClick={() => open('ADJUSTMENT', item)} className="rounded-lg p-2 text-slate-600 hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-800">
                      <PencilRuler className="h-4 w-4" />
                    </button>
                    <button title="Batas peringatan" onClick={() => open('MIN', item)} className="rounded-lg p-2 text-slate-600 hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-800">
                      <Bell className="h-4 w-4" />
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </Table>
          {view === 'all' && <Pagination meta={list.meta} onPage={list.setPage} />}
        </Card>
      )}

      <Modal
        open={!!mode && !!target}
        onClose={() => setMode(null)}
        size="sm"
        title={mode ? MODE_TEXT[mode].title : ''}
        subtitle={target ? `${target.product_name} (${target.unit_name}) · stok sekarang ${formatNumber(target.current_stock)}` : undefined}
        footer={
          <>
            <button className={btnSecondary} onClick={() => setMode(null)}>Batal</button>
            <button className={btnPrimary} disabled={busy} onClick={submit}>Simpan</button>
          </>
        }
      >
        {mode && (
          <div className="grid gap-3">
            <Field
              label={mode === 'ADJUSTMENT' ? 'Stok sebenarnya' : mode === 'MIN' ? 'Batas minimum' : 'Jumlah'}
              hint={MODE_TEXT[mode].hint}
            >
              <input type="number" min={0} className={inputClass} value={qty} onChange={(e) => setQty(e.target.value)} autoFocus />
            </Field>
            {mode !== 'MIN' && (
              <>
                <Field label="Keterangan">
                  <input
                    className={inputClass}
                    value={notes}
                    onChange={(e) => setNotes(e.target.value)}
                    placeholder={mode === 'STOCK_IN' ? 'Misal: kiriman dari distributor' : mode === 'STOCK_OUT' ? 'Misal: 2 karung bocor' : 'Misal: stok opname akhir bulan'}
                  />
                </Field>
                <Field label="No. referensi / nota (opsional)">
                  <input className={inputClass} value={reference} onChange={(e) => setReference(e.target.value)} />
                </Field>
              </>
            )}
          </div>
        )}
      </Modal>
    </div>
  );
};
