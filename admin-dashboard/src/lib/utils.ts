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

export function formatNumber(n: number): string {
  return new Intl.NumberFormat('id-ID').format(n || 0);
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
  if (data?.errors && Array.isArray(data.errors) && data.errors.length > 0 && data.errors[0]?.message) {
    return data.errors.map((x) => x.message).join(', ');
  }
  if (data?.message) return data.message;
  if (e?.code === 'ERR_NETWORK') return 'Tidak dapat terhubung ke server. Periksa koneksi atau pastikan backend berjalan.';
  return fallback;
}

export const ROLE_LABELS: Record<string, string> = {
  SUPER_ADMIN: 'Super Admin',
  STORE_MANAGER: 'Manajer Toko',
  WAREHOUSE_STAFF: 'Staf Gudang',
  DRIVER: 'Kurir',
  B2B_CUSTOMER: 'Mitra Toko (B2B)',
  B2C_CUSTOMER: 'Pelanggan',
};

export const STATUS_LABELS: Record<string, string> = {
  // Orders
  CREATED: 'Baru',
  PENDING_PAYMENT: 'Menunggu Bayar',
  PAID: 'Lunas',
  CONFIRMED: 'Pesanan Baru',
  PROCESSING_PACKING: 'Disiapkan',
  OUT_FOR_DELIVERY: 'Diantar',
  COMPLETED: 'Selesai',
  CANCELLED: 'Dibatalkan',
  // Payments
  PENDING: 'Belum Dibayar',
  WAITING_CONFIRMATION: 'Perlu Verifikasi',
  FAILED: 'Ditolak',
  REFUNDED: 'Dikembalikan',
  PARTIALLY_REFUNDED: 'Refund Sebagian',
  EXPIRED: 'Kedaluwarsa',
  UNPAID: 'Belum Lunas',
  // Shipments
  READY_TO_SHIP: 'Siap Kirim',
  PICKED_UP: 'Diambil Kurir',
  IN_TRANSIT: 'Dalam Perjalanan',
  DELIVERED: 'Terkirim',
  RETURNED: 'Dikembalikan',
  // B2B
  APPROVED: 'Disetujui',
  REJECTED: 'Ditolak',
  // Stock
  INITIAL_STOCK: 'Stok Awal',
  STOCK_IN: 'Barang Masuk',
  STOCK_OUT: 'Barang Keluar',
  ADJUSTMENT: 'Koreksi',
  // Kasbon
  CREDIT_USAGE: 'Pemakaian',
  PAYMENT_SETTLEMENT: 'Pelunasan',
  REVERSAL: 'Pengembalian',
};

export const statusLabel = (s?: string) => (s ? STATUS_LABELS[s] || s : '-');

export const PAYMENT_METHOD_LABELS: Record<string, string> = {
  COD: 'Bayar Tunai',
  MANUAL_BANK_TRANSFER: 'Transfer Bank',
  B2B_CREDIT_TEMPO: 'Kasbon',
};

/** Order status label; for a pickup order "out" means waiting at the store. */
export const orderStatusLabel = (status: string, deliveryType?: string) =>
  deliveryType === 'PICKUP' && status === 'OUT_FOR_DELIVERY' ? 'Siap Diambil' : statusLabel(status);

export const DELIVERY_LABELS: Record<string, string> = {
  DELIVERY: 'Diantar',
  PICKUP: 'Ambil di Toko',
};

/** wa.me link for an Indonesian phone number such as 0812... or +62812... */
export function whatsappLink(phone?: string, text?: string): string | null {
  const digits = (phone || '').replace(/\D/g, '');
  if (digits.length < 9) return null;
  const intl = digits.startsWith('0') ? `62${digits.slice(1)}` : digits;
  return `https://wa.me/${intl}${text ? `?text=${encodeURIComponent(text)}` : ''}`;
}

/** Payment method label, also for older codes such as TRANSFER_BCA. */
export const paymentMethodLabel = (code: string) =>
  PAYMENT_METHOD_LABELS[code] || (code.includes('TRANSFER') ? 'Transfer Bank' : code);
