import React from 'react';
import { orderStatusInfo, PAYMENT_STATUS, TONE_CLASS } from '../../lib/utils';

export const Badge: React.FC<{ tone?: string; children: React.ReactNode }> = ({ tone = 'slate', children }) => (
  <span className={`inline-flex items-center whitespace-nowrap rounded-full px-2.5 py-0.5 text-xs font-bold ${TONE_CLASS[tone] || TONE_CLASS.slate}`}>
    {children}
  </span>
);

export const OrderStatusBadge: React.FC<{ status: string; deliveryType?: string }> = ({ status, deliveryType }) => {
  const s = orderStatusInfo(status, deliveryType);
  return <Badge tone={s.tone}>{s.label}</Badge>;
};

export const PaymentStatusBadge: React.FC<{ status: string }> = ({ status }) => {
  const s = PAYMENT_STATUS[status] || { label: status, tone: 'slate' };
  return <Badge tone={s.tone}>{s.label}</Badge>;
};

export const inputClass =
  'w-full rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 px-3.5 py-2.5 text-sm text-slate-900 dark:text-white placeholder-slate-400 focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/30';

export const btnPrimary =
  'inline-flex items-center justify-center gap-1.5 rounded-2xl bg-primary-600 px-5 py-2.5 text-sm font-bold text-white shadow-md shadow-primary-600/20 hover:bg-primary-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors';

export const btnSecondary =
  'inline-flex items-center justify-center gap-1.5 rounded-2xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 px-5 py-2.5 text-sm font-bold text-slate-700 dark:text-slate-200 hover:bg-slate-50 dark:hover:bg-slate-700 disabled:opacity-50 transition-colors';

export const Field: React.FC<{ label: string; children: React.ReactNode; hint?: string; className?: string }> = ({ label, children, hint, className }) => (
  <label className={`block ${className || ''}`}>
    <span className="mb-1 block text-xs font-semibold text-slate-500 dark:text-slate-400">{label}</span>
    {children}
    {hint && <span className="mt-1 block text-[11px] text-slate-400">{hint}</span>}
  </label>
);

export const Panel: React.FC<{ title?: React.ReactNode; children: React.ReactNode; className?: string }> = ({ title, children, className = '' }) => (
  <section className={`rounded-3xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 p-5 sm:p-6 shadow-sm ${className}`}>
    {title && <h2 className="mb-4 flex items-center gap-2 text-base font-bold text-slate-900 dark:text-white">{title}</h2>}
    {children}
  </section>
);
