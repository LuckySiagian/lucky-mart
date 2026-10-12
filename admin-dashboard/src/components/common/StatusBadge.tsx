import React from 'react';
import { statusLabel } from '../../lib/utils';

const GREEN = ['PAID', 'DELIVERED', 'COMPLETED', 'SENT', 'APPROVED', 'STOCK_IN', 'INITIAL_STOCK', 'PAYMENT_SETTLEMENT', 'REVERSAL'];
const AMBER = ['CONFIRMED', 'PENDING', 'PENDING_PAYMENT', 'WAITING_CONFIRMATION', 'PROCESSING_PACKING', 'READY_TO_SHIP', 'CREATED', 'UNPAID', 'ADJUSTMENT', 'CREDIT_USAGE'];
const BLUE = ['OUT_FOR_DELIVERY', 'IN_TRANSIT', 'PICKED_UP'];
const RED = ['CANCELLED', 'FAILED', 'EXPIRED', 'REFUNDED', 'PARTIALLY_REFUNDED', 'REJECTED', 'RETURNED', 'STOCK_OUT'];

export const StatusBadge: React.FC<{ status: string; label?: string }> = ({ status, label }) => {
  const upper = (status || '').toUpperCase();
  let style = 'bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-300 border border-slate-200 dark:border-slate-700';
  if (GREEN.includes(upper)) style = 'bg-emerald-100 text-emerald-800 dark:bg-emerald-950/60 dark:text-emerald-400 border border-emerald-200 dark:border-emerald-900';
  else if (AMBER.includes(upper)) style = 'bg-amber-100 text-amber-800 dark:bg-amber-950/60 dark:text-amber-400 border border-amber-200 dark:border-amber-900';
  else if (BLUE.includes(upper)) style = 'bg-sky-100 text-sky-800 dark:bg-sky-950/60 dark:text-sky-400 border border-sky-200 dark:border-sky-900';
  else if (RED.includes(upper)) style = 'bg-rose-100 text-rose-800 dark:bg-rose-950/60 dark:text-rose-400 border border-rose-200 dark:border-rose-900';

  return (
    <span className={`inline-flex items-center whitespace-nowrap rounded-full px-2.5 py-0.5 text-xs font-semibold ${style}`}>
      {label || statusLabel(upper)}
    </span>
  );
};
