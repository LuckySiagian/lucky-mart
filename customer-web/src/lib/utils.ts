import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';
import type { AxiosError } from 'axios';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatCurrency(amount: number): string {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(amount || 0);
}

export function formatDate(dateString?: string | Date): string {
  if (!dateString) return '-';
  const date = new Date(dateString);
  if (Number.isNaN(date.getTime())) return '-';
  return new Intl.DateTimeFormat('id-ID', { dateStyle: 'medium', timeStyle: 'short' }).format(date);
}

/** Human-readable message from an API error (the backend sends `message`). */
export function getErrorMessage(err: unknown, fallback = 'Terjadi kesalahan. Coba lagi.'): string {
  const e = err as AxiosError<{ message?: string; errors?: { field: string; message: string }[] }>;
  const data = e?.response?.data;
  if (data?.errors && Array.isArray(data.errors) && data.errors[0]?.message) {
    return data.errors.map((x) => x.message).join(', ');
  }
  if (data?.message) return data.message;
  if (e?.code === 'ERR_NETWORK') return 'Tidak dapat terhubung ke server. Periksa koneksi internet Anda.';
  return fallback;
}

export const ORDER_STATUS: Record<string, { label: string; tone: string }> = {
  CREATED: { label: 'Pesanan Baru', tone: 'amber' },
  PENDING_PAYMENT: { label: 'Menunggu Pembayaran', tone: 'amber' },
  PAID: { label: 'Sudah Dibayar', tone: 'green' },
  CONFIRMED: { label: 'Pesanan Diterima', tone: 'green' },
  PROCESSING_PACKING: { label: 'Sedang Disiapkan', tone: 'sky' },
  OUT_FOR_DELIVERY: { label: 'Sedang Diantar', tone: 'sky' },
  COMPLETED: { label: 'Selesai', tone: 'green' },
  CANCELLED: { label: 'Dibatalkan', tone: 'red' },
};

/** Status for an order; a pickup order that is "out" is waiting at the store. */
export const orderStatusInfo = (status: string, deliveryType?: string) =>
  deliveryType === 'PICKUP' && status === 'OUT_FOR_DELIVERY'
    ? { label: 'Siap Diambil di Toko', tone: 'sky' }
    : ORDER_STATUS[status] || { label: status, tone: 'slate' };

export const PAYMENT_STATUS: Record<string, { label: string; tone: string }> = {
  PENDING: { label: 'Belum Dibayar', tone: 'amber' },
  WAITING_CONFIRMATION: { label: 'Sedang Diverifikasi', tone: 'sky' },
  PAID: { label: 'Lunas', tone: 'green' },
  FAILED: { label: 'Bukti Ditolak', tone: 'red' },
  CANCELLED: { label: 'Dibatalkan', tone: 'red' },
  REFUNDED: { label: 'Dikembalikan', tone: 'red' },
  PARTIALLY_REFUNDED: { label: 'Dikembalikan Sebagian', tone: 'red' },
};

export const SHIPMENT_STATUS: Record<string, string> = {
  PENDING: 'Menunggu',
  READY_TO_SHIP: 'Siap Dikirim',
  PICKED_UP: 'Diambil Kurir',
  IN_TRANSIT: 'Dalam Perjalanan',
  OUT_FOR_DELIVERY: 'Sedang Diantar',
  DELIVERED: 'Diterima',
  FAILED: 'Gagal Diantar',
  RETURNED: 'Dikembalikan',
};

export const TONE_CLASS: Record<string, string> = {
  amber: 'bg-amber-100 text-amber-800 dark:bg-amber-950/60 dark:text-amber-300',
  green: 'bg-emerald-100 text-emerald-800 dark:bg-emerald-950/60 dark:text-emerald-300',
  sky: 'bg-sky-100 text-sky-800 dark:bg-sky-950/60 dark:text-sky-300',
  red: 'bg-rose-100 text-rose-800 dark:bg-rose-950/60 dark:text-rose-300',
  slate: 'bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-300',
};

export const PASSWORD_HINT = 'Minimal 8 karakter, ada huruf besar, huruf kecil, angka, dan simbol (contoh: Rahasia#123).';

export const isPhone = (s: string) => /^(08|\+628)[0-9]{8,11}$/.test(s.trim());
