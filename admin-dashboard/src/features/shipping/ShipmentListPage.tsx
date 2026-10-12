import React, { useState } from 'react';
import { toast } from 'sonner';
import { MapPin, Phone, Truck } from 'lucide-react';
import { PageHeader } from '../../components/common/PageHeader';
import { StatusBadge } from '../../components/common/StatusBadge';
import { LoadingSkeleton } from '../../components/common/LoadingSkeleton';
import { EmptyState } from '../../components/common/EmptyState';
import { ErrorState } from '../../components/common/ErrorState';
import { Field, Modal, btnPrimary, btnSecondary, inputClass } from '../../components/common/Modal';
import { Card, FilterTabs, Pagination, SearchInput } from '../../components/common/Pagination';
import { shippingService } from '../../services/apiServices';
import { useList } from '../../hooks/useList';
import { formatCurrency, formatDate, getErrorMessage, statusLabel } from '../../lib/utils';
import type { Shipment } from '../../types';

const TABS = [
  { value: '', label: 'Semua' },
  { value: 'READY_TO_SHIP', label: 'Siap Kirim' },
  { value: 'OUT_FOR_DELIVERY', label: 'Diantar' },
  { value: 'FAILED', label: 'Gagal' },
  { value: 'DELIVERED', label: 'Terkirim' },
];

const STATUS_OPTIONS = ['READY_TO_SHIP', 'PICKED_UP', 'OUT_FOR_DELIVERY', 'DELIVERED', 'FAILED', 'RETURNED'];

export const ShipmentListPage: React.FC = () => {
  const list = useList((p) => shippingService.getShipments(p));
  const [selected, setSelected] = useState<Shipment | null>(null);
  const [status, setStatus] = useState('');
  const [location, setLocation] = useState('');
  const [notes, setNotes] = useState('');
  const [busy, setBusy] = useState(false);

  const open = (s: Shipment) => {
    setSelected(s);
    setStatus(s.status === 'READY_TO_SHIP' ? 'OUT_FOR_DELIVERY' : s.status === 'OUT_FOR_DELIVERY' ? 'DELIVERED' : s.status);
    setLocation('');
    setNotes('');
  };

  const finished = selected?.status === 'DELIVERED' || selected?.status === 'RETURNED';

  const submit = async () => {
    if (!selected) return;
    setBusy(true);
    try {
      const updated = await shippingService.updateStatus(selected.id, { status, location: location.trim(), notes: notes.trim() });
      toast.success(`Resi ${updated.tracking_number}: ${statusLabel(updated.status)}`);
      setSelected(null);
      list.reload();
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setBusy(false);
    }
  };

  return (
    <div className="space-y-5">
      <PageHeader
        title="Pengiriman"
        subtitle="Resi dibuat otomatis saat pesanan dikirim. Tandai Terkirim agar pesanan selesai dan pembayaran COD tercatat lunas."
      />

      <div className="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
        <FilterTabs value={(list.params.status as string) || ''} onChange={(v) => list.setFilter('status', v)} options={TABS} />
        <SearchInput value={list.searchInput} onChange={list.setSearchInput} placeholder="Cari nomor resi..." />
      </div>

      {list.loading ? (
        <LoadingSkeleton count={4} />
      ) : list.error ? (
        <ErrorState message={list.error} onRetry={list.reload} />
      ) : list.items.length === 0 ? (
        <EmptyState title="Belum ada pengiriman" description='Ubah status pesanan menjadi "Kirim ke Pelanggan" di menu Pesanan untuk membuat resi.' />
      ) : (
        <div className="grid gap-3 md:grid-cols-2">
          {list.items.map((s) => (
            <Card key={s.id} className="p-4">
              <div className="flex items-start justify-between gap-3">
                <div>
                  <div className="font-mono text-xs font-semibold">{s.tracking_number}</div>
                  <div className="text-xs text-slate-500">Pesanan {s.order_number} · {s.courier_name}</div>
                </div>
                <StatusBadge status={s.status} />
              </div>
              <div className="mt-3 space-y-1 text-sm">
                <div className="font-semibold">{s.recipient_name || s.customer_name}</div>
                {s.recipient_phone && <div className="flex items-center gap-1.5 text-xs text-slate-500"><Phone className="h-3 w-3" />{s.recipient_phone}</div>}
                {s.address && <div className="flex items-start gap-1.5 text-xs text-slate-500"><MapPin className="h-3 w-3 mt-0.5" />{s.address}</div>}
              </div>
              <div className="mt-3 flex items-center justify-between border-t border-slate-100 dark:border-slate-800 pt-3">
                <span className="text-xs text-slate-500">Ongkir {formatCurrency(s.shipping_fee)}</span>
                <button onClick={() => open(s)} className="inline-flex items-center gap-1 rounded-lg px-2.5 py-1.5 text-xs font-semibold text-brand-700 hover:bg-brand-50 dark:text-brand-400 dark:hover:bg-brand-950/40">
                  <Truck className="h-4 w-4" /> {s.status === 'DELIVERED' || s.status === 'RETURNED' ? 'Lihat' : 'Perbarui'}
                </button>
              </div>
            </Card>
          ))}
        </div>
      )}
      {!list.loading && <Pagination meta={list.meta} onPage={list.setPage} />}

      <Modal
        open={!!selected}
        onClose={() => setSelected(null)}
        size="md"
        title={selected ? `Resi ${selected.tracking_number}` : ''}
        subtitle={selected ? `Pesanan ${selected.order_number}` : undefined}
        footer={
          !finished && (
            <>
              <button className={btnSecondary} onClick={() => setSelected(null)}>Tutup</button>
              <button className={btnPrimary} disabled={busy || status === selected?.status} onClick={submit}>Simpan Status</button>
            </>
          )
        }
      >
        {selected && (
          <div className="space-y-4 text-sm">
            {!finished && (
              <div className="grid gap-3">
                <Field label="Status baru">
                  <select className={inputClass} value={status} onChange={(e) => setStatus(e.target.value)}>
                    {STATUS_OPTIONS.map((s) => (
                      <option key={s} value={s}>{statusLabel(s)}</option>
                    ))}
                  </select>
                </Field>
                <Field label="Lokasi (opsional)">
                  <input className={inputClass} value={location} onChange={(e) => setLocation(e.target.value)} placeholder="Misal: Jl. Gatot Subroto" />
                </Field>
                <Field label="Catatan (opsional)">
                  <input className={inputClass} value={notes} onChange={(e) => setNotes(e.target.value)} placeholder="Misal: diterima oleh Ibu Sari" />
                </Field>
              </div>
            )}
            <div>
              <div className="mb-2 text-xs font-semibold text-slate-500">Lacak</div>
              <ol className="space-y-3 border-l-2 border-slate-200 dark:border-slate-800 pl-4">
                {selected.trackings.map((t) => (
                  <li key={t.id} className="relative">
                    <span className="absolute -left-[21px] top-1 h-2.5 w-2.5 rounded-full bg-brand-500" />
                    <div className="flex flex-wrap items-center gap-2">
                      <StatusBadge status={t.status} />
                      <span className="text-xs text-slate-400">{formatDate(t.created_at)}</span>
                    </div>
                    <p className="mt-1 text-xs text-slate-600 dark:text-slate-400">
                      {t.location && t.location !== '-' ? `${t.location} — ` : ''}{t.notes}
                    </p>
                  </li>
                ))}
              </ol>
            </div>
          </div>
        )}
      </Modal>
    </div>
  );
};
