# Data & Akun Demo Lucky Mart

Data demo dibuat oleh `go run ./cmd/seed-demo` (jalankan dari folder `Lucky Mart`
saat database menyala). Program ini **hanya untuk development**: menolak jalan
kalau `APP_ENV=production`, dan tidak mengubah apa pun kalau data demo sudah ada.

> ⚠️ Semua akun di bawah memakai kata sandi yang sama dan tertulis di sini.
> Jangan pernah menjalankan seed-demo di server yang dipakai pelanggan sungguhan.

## Akun

Kata sandi semua akun demo: **`LuckyDemo#2026`**

| Peran | Nama | Nomor HP | Login di |
|---|---|---|---|
| Manajer Toko | Maria Simanjuntak | 089900000301 | Admin (port 3002) |
| Staf Gudang | Andi Pratama | 089900000302 | Admin |
| Kurir | Joko Susilo | 089900000303 | Admin |
| Pelanggan | Sari Wulandari | 089900000101 | Toko (port 3001) |
| Pelanggan | Budi Santoso | 089900000102 | Toko |
| Pelanggan | Dewi Lestari | 089900000103 | Toko |
| Pelanggan | Ahmad Fauzi | 089900000104 | Toko |
| Pelanggan | Rina Marlina | 089900000105 | Toko |
| Pelanggan | Yohana Sitompul | 089900000106 | Toko |
| Mitra warung (kasbon Rp3 jt) | Tini Rahayu — Warung Bu Tini | 089900000201 | Toko |
| Mitra warung (kasbon Rp5 jt) | Harun Nasution — Toko Sinar Jaya | 089900000202 | Toko |
| Mitra warung (belum diverifikasi) | Rosmawati — Kedai Mak Ros | 089900000203 | Toko |

Super Admin tetap akun bawaan: `081234567890` dengan kata sandi dari
`SEED_ADMIN_PASSWORD` (atau kata sandi lama kalau akunnya sudah ada).

## Isi data

- **39 produk** di 11 kategori (beras, minyak, gula, telur, tepung, mie, bumbu &
  bahan segar, kopi/teh/susu, minuman & air, gas, kebutuhan rumah), banyak yang
  punya satuan eceran **dan** grosir (pcs/dus, pouch/dus, kg/tray).
  Harga mengikuti harga toko kelontong Oktober 2026 (acuan data harga pangan
  PIHPS/Siskaperbapo dan harga rak Alfamart/Indomaret), contoh: gula pasir
  Gulaku 1 kg Rp20.000, Bimoli 2 L Rp42.500, telur Rp29.500/kg,
  Indomie goreng Rp3.500. **Sesuaikan dengan harga toko sendiri.**
- Beberapa barang sengaja **stoknya menipis** (Bright Gas, Dancow, cabai) untuk
  mencoba peringatan stok.
- **±35 pesanan** selama 2 minggu terakhir: diantar & ambil sendiri, bayar
  tunai / transfer / kasbon, plus 2 pesanan batal.
- **Pekerjaan hari ini** (menu Pesanan):
  - 2 pesanan baru (1 diantar, 1 ambil di toko)
  - 1 pesanan sedang disiapkan (kasbon Toko Sinar Jaya)
  - 1 sedang diantar, 1 siap diambil di toko
  - 1 bukti transfer menunggu dicek, 1 pesanan belum dibayar
- **Kasbon**: Warung Bu Tini punya kasbon jatuh tempo (lewat 30 hari) dan sudah
  mencicil Rp500.000; Toko Sinar Jaya punya kasbon berjalan.
- **Barang masuk** dari supplier tercatat di Mutasi Stok.
- **Info toko**: alamat contoh, jam buka, pengumuman, gratis ongkir mulai
  Rp300.000. Nomor WhatsApp sengaja dikosongkan — isi dengan nomor asli toko di
  Admin → Pengaturan → Info toko.
- **Metode bayar**: rekening transfer berisi nomor CONTOH; ganti dengan rekening
  asli di Admin → Pengaturan → Metode pembayaran.

## Yang bisa dicoba

1. **Pelanggan** (Sari): tambah ke keranjang → lihat "Belanja RpX lagi untuk
   gratis ongkir" → checkout diantar / ambil sendiri → bayar transfer dengan
   foto bukti.
2. **Mitra warung** (Bu Tini): beli dalam jumlah grosir → harga grosir otomatis →
   bayar kasbon → lihat sisa limit di halaman akun.
3. **Manajer** (Maria): Pesanan → cek bukti transfer → "Uang Sudah Masuk" →
   "Siapkan Barang" → cetak Daftar Barang & Nota → "Antar Sekarang" / "Siap
   Diambil" → selesai. Coba juga Kasbon Mitra (catat pelunasan), Pelanggan &
   Staf (verifikasi Kedai Mak Ros), Stok Barang (barang menipis), Pengaturan.
4. **Kurir** (Joko): hanya bisa menandai pesanan selesai.
5. **Staf gudang** (Andi): kelola produk & stok, tidak bisa melihat keuangan.
