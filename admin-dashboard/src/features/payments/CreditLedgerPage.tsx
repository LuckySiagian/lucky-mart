import React, { useState } from 'react';
import { toast } from 'sonner';
import { HandCoins, History, Plus, SlidersHorizontal } from 'lucide-react';
import { PageHeader } from '../../components/common/PageHeader';
import { StatusBadge } from '../../components/common/StatusBadge';
import { LoadingSkeleton } from '../../components/common/LoadingSkeleton';
import { EmptyState } from '../../components/common/EmptyState';
import { ErrorState } from '../../components/common/ErrorState';
import { Field, Modal, btnPrimary, btnSecondary, inputClass } from '../../components/common/Modal';
import { Card, Pagination, SearchInput, Table } from '../../components/common/Pagination';
import { creditService, userService } from '../../services/apiServices';
import { useList } from '../../hooks/useList';
import { formatCurrency, formatDate, getErrorMessage } from '../../lib/utils';
import type { CreditHistory, CreditLedger, User } from '../../types';

type Mode = 'limit' | 'settle' | 'history' | 'create' | null;

export const CreditLedgerPage: React.FC = () => {
  const list = useList((p) => creditService.getLedgers(p));
  const [selected, setSelected] = useState<CreditLedger | null>(null);
  const [mode, setMode] = useState<Mode>(null);
  const [amount, setAmount] = useState('');
  const [notes, setNotes] = useState('');
  const [busy, setBusy] = useState(false);
  const [history, setHistory] = useState<CreditHistory[]>([]);
  const [b2bUsers, setB2bUsers] = useState<User[]>([]);
  const [newUserId, setNewUserId] = useState('');

  const openMode = async (m: Mode, ledger: CreditLedger | null) => {
    setSelected(ledger);
    setMode(m);
    setNotes('');
    if (m === 'limit' && ledger) setAmount(String(ledger.credit_limit));
    else if (m === 'settle' && ledger) setAmount(String(ledger.outstanding_balance));
    else setAmount('');

    if (m === 'history' && ledger) {
      setHistory([]);
      try {
        setHistory((await creditService.getHistory(ledger.id, { limit: 100 })).items);
      } catch (err) {
        toast.error(getErrorMessage(err));
      }
    }
    if (m === 'create') {
      try {
        const res = await userService.getUsers({ role: 'B2B_CUSTOMER', limit: 100 });
        const withLedger = new Set(list.items.map((l) => l.user_id));
        setB2bUsers(res.items.filter((u) => !withLedger.has(u.id)));
        setNewUserId('');
      } catch (err) {
        toast.error(getErrorMessage(err));
      }
    }
  };

  const close = () => {
    setMode(null);
    setSelected(null);
  };

  const submit = async () => {
    const value = Number(amount);
    if (!Number.isFinite(value) || value < 0 || (mode === 'settle' && value <= 0)) {
      toast.error('Masukkan jumlah yang valid.');
      return;
    }
    setBusy(true);
    try {
      if (mode === 'limit' && selected) {
        await creditService.adjustLimit(selected.id, value);
        toast.success('Limit kasbon diperbarui.');
      } else if (mode === 'settle' && selected) {
        await creditService.settle(selected.id, value, notes.trim());
        toast.success(`Pelunasan ${formatCurrency(value)} dicatat.`);
      } else if (mode === 'create') {
        if (!newUserId) {
          toast.error('Pilih pelanggan mitra toko.');
          setBusy(false);
          return;
        }
        await creditService.createLedger(newUserId, value);
        toast.success('Akun kasbon dibuat.');
      }
      close();
      list.reload();
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setBusy(false);
    }
  };

  const totalOutstanding = list.items.reduce((s, l) => s + l.outstanding_balance, 0);

  return (
    <div className="space-y-5">
      <PageHeader title="Kasbon Mitra Toko" subtitle="Limit kredit, hutang berjalan, dan pelunasan pelanggan B2B.">
        <button className={btnPrimary} onClick={() => openMode('create', null)}>
          <Plus className="h-4 w-4" /> Buat Akun Kasbon
        </button>
      </PageHeader>

      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <SearchInput value={list.searchInput} onChange={list.setSearchInput} placeholder="Cari..." />
        <div className="text-sm text-slate-500">
          Total hutang berjalan (halaman ini): <b className="text-slate-900 dark:text-white">{formatCurrency(totalOutstanding)}</b>
        </div>
      </div>

      {list.loading ? (
        <LoadingSkeleton count={4} />
      ) : list.error ? (
        <ErrorState message={list.error} onRetry={list.reload} />
      ) : list.items.length === 0 ? (
        <EmptyState
          title="Belum ada akun kasbon"
          description="Setujui pendaftaran mitra toko di menu Pengguna dengan limit kredit, atau buat akun kasbon di sini."
        />
      ) : (
        <Card className="p-2 sm:p-4">
          <Table
            head={
              <tr>
                <th className="px-3 py-3">Mitra Toko</th>
                <th className="px-3 py-3 text-right">Limit</th>
                <th className="px-3 py-3 text-right">Hutang</th>
                <th className="px-3 py-3 text-right">Sisa Limit</th>
                <th className="px-3 py-3" />
              </tr>
            }
          >
            {list.items.map((l) => {
              const usage = l.credit_limit > 0 ? Math.min(100, (l.used_credit / l.credit_limit) * 100) : 0;
              return (
                <tr key={l.id} className="hover:bg-slate-50 dark:hover:bg-slate-800/40">
                  <td className="px-3 py-3">
                    <div className="font-semibold">{l.store_name || l.user_name}</div>
                    <div className="text-xs text-slate-500">{l.user_name} · {l.user_phone}</div>
                  </td>
                  <td className="px-3 py-3 text-right">{formatCurrency(l.credit_limit)}</td>
                  <td className="px-3 py-3 text-right">
                    <div className={`font-semibold ${l.outstanding_balance > 0 ? 'text-rose-600' : ''}`}>{formatCurrency(l.outstanding_balance)}</div>
                    <div className="ml-auto mt-1 h-1.5 w-24 rounded-full bg-slate-100 dark:bg-slate-800">
                      <div className={`h-1.5 rounded-full ${usage > 80 ? 'bg-rose-500' : 'bg-brand-500'}`} style={{ width: `${usage}%` }} />
                    </div>
                  </td>
                  <td className="px-3 py-3 text-right font-semibold text-emerald-600">{formatCurrency(l.remaining_credit)}</td>
                  <td className="px-3 py-3">
                    <div className="flex justify-end gap-1">
                      <button title="Catat pelunasan" disabled={l.outstanding_balance <= 0} onClick={() => openMode('settle', l)} className="rounded-lg p-2 text-emerald-600 hover:bg-emerald-50 disabled:opacity-30 dark:hover:bg-emerald-950/40">
                        <HandCoins className="h-4 w-4" />
                      </button>
                      <button title="Atur limit" onClick={() => openMode('limit', l)} className="rounded-lg p-2 text-slate-600 hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-800">
                        <SlidersHorizontal className="h-4 w-4" />
                      </button>
                      <button title="Riwayat" onClick={() => openMode('history', l)} className="rounded-lg p-2 text-slate-600 hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-800">
                        <History className="h-4 w-4" />
                      </button>
                    </div>
                  </td>
                </tr>
              );
            })}
          </Table>
          <Pagination meta={list.meta} onPage={list.setPage} />
        </Card>
      )}

      <Modal
        open={mode === 'limit' || mode === 'settle' || mode === 'create'}
        onClose={close}
        size="sm"
        title={mode === 'limit' ? 'Atur limit kasbon' : mode === 'settle' ? 'Catat pelunasan kasbon' : 'Buat akun kasbon'}
        subtitle={
          selected
            ? `${selected.store_name || selected.user_name} · hutang ${formatCurrency(selected.outstanding_balance)}`
            : 'Untuk pelanggan mitra toko (B2B) yang belum punya akun kasbon.'
        }
        footer={
          <>
            <button className={btnSecondary} onClick={close}>Batal</button>
            <button className={btnPrimary} disabled={busy} onClick={submit}>Simpan</button>
          </>
        }
      >
        <div className="space-y-3">
          {mode === 'create' && (
            <Field label="Mitra toko">
              <select className={inputClass} value={newUserId} onChange={(e) => setNewUserId(e.target.value)}>
                <option value="">— Pilih —</option>
                {b2bUsers.map((u) => (
                  <option key={u.id} value={u.id}>
                    {u.b2b_profile?.store_name || u.full_name} ({u.phone_number})
                  </option>
                ))}
              </select>
              {b2bUsers.length === 0 && <span className="mt-1 block text-[11px] text-slate-400">Semua mitra toko sudah punya akun kasbon.</span>}
            </Field>
          )}
          <Field label={mode === 'settle' ? 'Jumlah yang dibayar (Rp)' : 'Limit kredit (Rp)'}>
            <input type="number" min={0} className={inputClass} value={amount} onChange={(e) => setAmount(e.target.value)} />
          </Field>
          {mode === 'settle' && (
            <Field label="Catatan (opsional)">
              <input className={inputClass} value={notes} onChange={(e) => setNotes(e.target.value)} placeholder="Misal: bayar tunai di toko" />
            </Field>
          )}
        </div>
      </Modal>

      <Modal open={mode === 'history'} onClose={close} size="lg" title="Riwayat kasbon" subtitle={selected?.store_name || selected?.user_name}>
        {history.length === 0 ? (
          <p className="py-6 text-center text-sm text-slate-500">Belum ada transaksi.</p>
        ) : (
          <ul className="divide-y divide-slate-100 dark:divide-slate-800">
            {history.map((h) => (
              <li key={h.id} className="flex items-center justify-between gap-3 py-3 text-sm">
                <div>
                  <StatusBadge status={h.type} />
                  <div className="mt-1 text-xs text-slate-500">{h.notes}</div>
                  <div className="text-[11px] text-slate-400">{formatDate(h.payment_date)}</div>
                </div>
                <div className="text-right">
                  <div className={`font-semibold ${h.type === 'CREDIT_USAGE' ? 'text-rose-600' : 'text-emerald-600'}`}>
                    {h.type === 'CREDIT_USAGE' ? '+' : '−'}{formatCurrency(h.amount)}
                  </div>
                  <div className="text-[11px] text-slate-400">Sisa hutang {formatCurrency(h.remaining_balance)}</div>
                </div>
              </li>
            ))}
          </ul>
        )}
      </Modal>
    </div>
  );
};
