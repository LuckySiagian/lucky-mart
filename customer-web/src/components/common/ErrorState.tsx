import React from 'react';
import { AlertCircle, RotateCcw } from 'lucide-react';

export const ErrorState: React.FC<{ message?: string; onRetry?: () => void }> = ({
  message = 'Gagal memuat data dari server backend.',
  onRetry,
}) => {
  return (
    <div className="flex flex-col items-center justify-center rounded-3xl border border-rose-200 dark:border-rose-950 bg-rose-50/50 dark:bg-rose-950/20 p-8 text-center my-6">
      <AlertCircle className="h-10 w-10 text-rose-500 mb-3" />
      <h4 className="text-base font-bold text-rose-900 dark:text-rose-300 mb-1">Terjadi Kesalahan Koneksi</h4>
      <p className="text-xs text-rose-700 dark:text-rose-400 mb-4">{message}</p>
      {onRetry && (
        <button
          onClick={onRetry}
          className="inline-flex items-center space-x-2 rounded-xl bg-rose-600 px-4 py-2 text-xs font-bold text-white hover:bg-rose-700 transition-colors"
        >
          <RotateCcw className="h-3.5 w-3.5" />
          <span>Coba Muat Ulang</span>
        </button>
      )}
    </div>
  );
};
