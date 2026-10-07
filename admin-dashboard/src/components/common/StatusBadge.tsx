import React from 'react';

interface StatusBadgeProps {
  status: string;
}

export const StatusBadge: React.FC<StatusBadgeProps> = ({ status }) => {
  let badgeStyle = 'bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-300';

  const upper = status.toUpperCase();

  if (['PAID', 'DELIVERED', 'COMPLETED', 'CONFIRMED', 'SENT'].includes(upper)) {
    badgeStyle = 'bg-emerald-100 text-emerald-800 dark:bg-emerald-950/60 dark:text-emerald-400 border border-emerald-200 dark:border-emerald-900';
  } else if (['PENDING', 'WAITING_PAYMENT', 'WAITING_CONFIRMATION', 'PROCESSING', 'PACKING', 'IN_TRANSIT'].includes(upper)) {
    badgeStyle = 'bg-amber-100 text-amber-800 dark:bg-amber-950/60 dark:text-amber-400 border border-amber-200 dark:border-amber-900';
  } else if (['CANCELLED', 'FAILED', 'EXPIRED', 'REFUNDED'].includes(upper)) {
    badgeStyle = 'bg-rose-100 text-rose-800 dark:bg-rose-950/60 dark:text-rose-400 border border-rose-200 dark:border-rose-900';
  }

  return (
    <span className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-semibold ${badgeStyle}`}>
      {status}
    </span>
  );
};
