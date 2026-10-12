import React from 'react';
import { PageHeader } from '../../components/common/PageHeader';
import { StatusBadge } from '../../components/common/StatusBadge';
import { LoadingSkeleton } from '../../components/common/LoadingSkeleton';
import { EmptyState } from '../../components/common/EmptyState';
import { ErrorState } from '../../components/common/ErrorState';
import { Card, FilterTabs, Pagination, SearchInput, Table } from '../../components/common/Pagination';
import { inventoryService } from '../../services/apiServices';
import { useList } from '../../hooks/useList';
import { formatDate, formatNumber } from '../../lib/utils';

const TABS = [
  { value: '', label: 'Semua' },
  { value: 'STOCK_IN', label: 'Barang Masuk' },
  { value: 'STOCK_OUT', label: 'Barang Keluar' },
  { value: 'ADJUSTMENT', label: 'Koreksi' },
  { value: 'INITIAL_STOCK', label: 'Stok Awal' },
];

export const StockMovementPage: React.FC = () => {
  const list = useList((p) => inventoryService.getMovements(p), { limit: 30 });

  return (
    <div className="space-y-5">
      <PageHeader title="Mutasi Stok" subtitle="Setiap perubahan stok tercatat otomatis: penjualan, pembatalan, barang masuk, dan koreksi." />

      <div className="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
        <FilterTabs value={(list.params.type as string) || ''} onChange={(v) => list.setFilter('type', v)} options={TABS} />
        <SearchInput value={list.searchInput} onChange={list.setSearchInput} placeholder="Cari nomor referensi / pesanan..." />
      </div>

      {list.loading ? (
        <LoadingSkeleton count={6} />
      ) : list.error ? (
        <ErrorState message={list.error} onRetry={list.reload} />
      ) : list.items.length === 0 ? (
        <EmptyState title="Belum ada mutasi stok" />
      ) : (
        <Card className="p-2 sm:p-4">
          <Table
            head={
              <tr>
                <th className="px-3 py-3">Waktu</th>
                <th className="px-3 py-3">Produk</th>
                <th className="px-3 py-3">Jenis</th>
                <th className="px-3 py-3 text-right">Sebelum</th>
                <th className="px-3 py-3 text-right">Perubahan</th>
                <th className="px-3 py-3 text-right">Sesudah</th>
                <th className="px-3 py-3">Referensi</th>
              </tr>
            }
          >
            {list.items.map((m) => (
              <tr key={m.id}>
                <td className="whitespace-nowrap px-3 py-3 text-xs text-slate-500">{formatDate(m.created_at)}</td>
                <td className="px-3 py-3">
                  <div className="font-medium">{m.product_name || '-'}</div>
                  <div className="text-xs text-slate-500">{m.unit_name}</div>
                </td>
                <td className="px-3 py-3"><StatusBadge status={m.movement_type} /></td>
                <td className="px-3 py-3 text-right">{formatNumber(m.quantity_before)}</td>
                <td className={`px-3 py-3 text-right font-semibold ${m.quantity_changed < 0 ? 'text-rose-600' : 'text-emerald-600'}`}>
                  {m.quantity_changed > 0 ? '+' : ''}{formatNumber(m.quantity_changed)}
                </td>
                <td className="px-3 py-3 text-right font-semibold">{formatNumber(m.quantity_after)}</td>
                <td className="px-3 py-3">
                  <div className="font-mono text-xs">{m.reference_number}</div>
                  {m.notes && <div className="text-xs text-slate-500">{m.notes}</div>}
                </td>
              </tr>
            ))}
          </Table>
          <Pagination meta={list.meta} onPage={list.setPage} />
        </Card>
      )}
    </div>
  );
};
