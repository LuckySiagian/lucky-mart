import React from 'react';
import { AlertTriangle, RotateCcw } from 'lucide-react';

interface ErrorStateProps {
  title?: string;
  message?: string;
  onRetry?: () => void;
}

export const ErrorState: React.FC<ErrorStateProps> = ({
  title = 'Gagal Memuat Data',
  message = 'Terjadi kesalahan saat berkomunikasi dengan server API backend.',
  onRetry,
}) => {
  return (
    <div className="flex flex-col items-center justify-center rounded-2xl border border-rose-200 dark:border-rose-950 bg-rose-50/50 dark:bg-rose-950/20 p-8 text-center my-6">
      <div className="flex h-14 w-14 items-center justify-center rounded-xl bg-rose-100 dark:bg-rose-900/50 text-rose-600 mb-3">
        <AlertTriangle className="h-7 w-7" />
      </div>
      <h3 className="text-base font-bold text-rose-900 dark:text-rose-300 mb-1">{title}</h3>
      <p className="text-xs text-rose-700 dark:text-rose-400 max-w-md mb-4">{message}</p>
      {onRetry && (
        <button
          onClick={onRetry}
          className="inline-flex items-center space-x-2 rounded-lg bg-rose-600 px-3.5 py-1.5 text-xs font-semibold text-white hover:bg-rose-700 transition-colors"
        >
          <RotateCcw className="h-3.5 w-3.5" />
          <span>Coba Lagi</span>
        </button>
      )}
    </div>
  );
};
