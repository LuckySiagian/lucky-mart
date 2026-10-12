import React from 'react';
import { ChevronLeft, ChevronRight } from 'lucide-react';

interface PaginationProps {
  meta: { page: number; total_pages: number; total_records: number } | null;
  onPage: (page: number) => void;
}

export const Pagination: React.FC<PaginationProps> = ({ meta, onPage }) => {
  if (!meta || meta.total_pages <= 1) {
    return meta ? <p className="px-1 pt-3 text-xs text-slate-500">{meta.total_records} data</p> : null;
  }
  return (
    <div className="flex items-center justify-between gap-3 pt-3 text-xs text-slate-500">
      <span>
        Halaman {meta.page} dari {meta.total_pages} · {meta.total_records} data
      </span>
      <div className="flex gap-1.5">
        <button
          disabled={meta.page <= 1}
          onClick={() => onPage(meta.page - 1)}
          className="inline-flex items-center rounded-lg border border-slate-200 dark:border-slate-700 px-2.5 py-1.5 font-semibold disabled:opacity-40 hover:bg-slate-100 dark:hover:bg-slate-800"
        >
          <ChevronLeft className="h-4 w-4" /> Sebelumnya
        </button>
        <button
          disabled={meta.page >= meta.total_pages}
          onClick={() => onPage(meta.page + 1)}
          className="inline-flex items-center rounded-lg border border-slate-200 dark:border-slate-700 px-2.5 py-1.5 font-semibold disabled:opacity-40 hover:bg-slate-100 dark:hover:bg-slate-800"
        >
          Berikutnya <ChevronRight className="h-4 w-4" />
        </button>
      </div>
    </div>
  );
};

export const SearchInput: React.FC<{ value: string; onChange: (v: string) => void; placeholder?: string }> = ({ value, onChange, placeholder }) => (
  <input
    type="search"
    value={value}
    onChange={(e) => onChange(e.target.value)}
    placeholder={placeholder || 'Cari...'}
    className="w-full sm:w-72 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 px-3.5 py-2.5 text-sm text-slate-900 dark:text-white placeholder-slate-400 focus:border-brand-500 focus:outline-none focus:ring-2 focus:ring-brand-500/30"
  />
);

export const FilterTabs: React.FC<{
  value: string;
  onChange: (v: string) => void;
  options: { value: string; label: string }[];
}> = ({ value, onChange, options }) => (
  <div className="flex gap-1.5 overflow-x-auto pb-1">
    {options.map((o) => (
      <button
        key={o.value}
        onClick={() => onChange(o.value)}
        className={`whitespace-nowrap rounded-xl px-3 py-2 text-xs font-semibold transition-colors ${
          value === o.value
            ? 'bg-brand-600 text-white shadow-sm'
            : 'bg-white dark:bg-slate-900 text-slate-600 dark:text-slate-300 border border-slate-200 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800'
        }`}
      >
        {o.label}
      </button>
    ))}
  </div>
);

export const Card: React.FC<{ children: React.ReactNode; className?: string }> = ({ children, className = '' }) => (
  <div className={`rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 shadow-sm ${className}`}>{children}</div>
);

export const Table: React.FC<{ head: React.ReactNode; children: React.ReactNode }> = ({ head, children }) => (
  <div className="overflow-x-auto">
    <table className="w-full min-w-[720px] text-left text-sm">
      <thead className="border-b border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-800/40 text-xs font-semibold uppercase tracking-wide text-slate-500">
        {head}
      </thead>
      <tbody className="divide-y divide-slate-100 dark:divide-slate-800">{children}</tbody>
    </table>
  </div>
);
