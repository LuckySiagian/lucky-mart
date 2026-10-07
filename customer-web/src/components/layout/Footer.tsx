import React from 'react';
import { Store, Phone, Mail, MapPin, ShieldCheck, Truck, CreditCard } from 'lucide-react';

export const Footer: React.FC = () => {
  return (
    <footer className="bg-slate-900 text-slate-400 text-sm mt-16 border-t border-slate-800">
      {/* Service Highlights Header */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10 border-b border-slate-800">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 text-center md:text-left">
          <div className="flex items-center space-x-4">
            <div className="h-12 w-12 rounded-2xl bg-primary-600/20 text-primary-400 flex items-center justify-center flex-shrink-0">
              <Truck className="h-6 w-6" />
            </div>
            <div>
              <h4 className="text-white font-bold">Pengiriman Cepat & Terjadwal</h4>
              <p className="text-xs text-slate-400 mt-0.5">Pengiriman sembako Same-Day & Cargo Grosir</p>
            </div>
          </div>

          <div className="flex items-center space-x-4">
            <div className="h-12 w-12 rounded-2xl bg-primary-600/20 text-primary-400 flex items-center justify-center flex-shrink-0">
              <ShieldCheck className="h-6 w-6" />
            </div>
            <div>
              <h4 className="text-white font-bold">100% Produk Original & Segar</h4>
              <p className="text-xs text-slate-400 mt-0.5">Langsung dari Pemasok Resmi & Produsen Sembako</p>
            </div>
          </div>

          <div className="flex items-center space-x-4">
            <div className="h-12 w-12 rounded-2xl bg-primary-600/20 text-primary-400 flex items-center justify-center flex-shrink-0">
              <CreditCard className="h-6 w-6" />
            </div>
            <div>
              <h4 className="text-white font-bold">B2B Kasbon Tempo 30 Hari</h4>
              <p className="text-xs text-slate-400 mt-0.5">Kemudahan transaksi khusus Mitra Toko Kelontong</p>
            </div>
          </div>
        </div>
      </div>

      {/* Main Footer Links */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
          <div>
            <div className="flex items-center space-x-2 mb-4">
              <div className="h-8 w-8 rounded-xl bg-primary-600 flex items-center justify-center text-white font-bold">
                <Store className="h-5 w-5" />
              </div>
              <span className="text-lg font-black text-white">LUCKYMART</span>
            </div>
            <p className="text-xs text-slate-400 leading-relaxed mb-4">
              Platform belanja sembako online terpercaya untuk kebutuhan rumah tangga dan pasokan grosir usaha kelontong keluarga.
            </p>
          </div>

          <div>
            <h5 className="text-white font-bold mb-3">Kategori Populer</h5>
            <ul className="space-y-2 text-xs">
              <li><a href="#" className="hover:text-primary-400">Beras & Palawija</a></li>
              <li><a href="#" className="hover:text-primary-400">Minyak Goreng & Margarin</a></li>
              <li><a href="#" className="hover:text-primary-400">Gula Pasir & Pemanis</a></li>
              <li><a href="#" className="hover:text-primary-400">Telur & Hasil Tani</a></li>
              <li><a href="#" className="hover:text-primary-400">Mie Instan & Bumbu Dapur</a></li>
            </ul>
          </div>

          <div>
            <h5 className="text-white font-bold mb-3">Layanan Pelanggan</h5>
            <ul className="space-y-2 text-xs">
              <li><a href="#" className="hover:text-primary-400">Cara Pemesanan</a></li>
              <li><a href="#" className="hover:text-primary-400">Pendaftaran B2B Mitra Warung</a></li>
              <li><a href="#" className="hover:text-primary-400">Syarat & Ketentuan Kasbon Tempo</a></li>
              <li><a href="#" className="hover:text-primary-400">Kebijakan Pengembalian Sembako</a></li>
            </ul>
          </div>

          <div>
            <h5 className="text-white font-bold mb-3">Hubungi Kami</h5>
            <div className="space-y-2 text-xs">
              <div className="flex items-center space-x-2">
                <Phone className="h-4 w-4 text-primary-400" />
                <span>+62 812-3456-7890</span>
              </div>
              <div className="flex items-center space-x-2">
                <Mail className="h-4 w-4 text-primary-400" />
                <span>cs@luckymart.com</span>
              </div>
              <div className="flex items-center space-x-2">
                <MapPin className="h-4 w-4 text-primary-400 flex-shrink-0" />
                <span>Gudang Pusat Sembako Cengkareng, Jakarta Barat</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="bg-slate-950 py-4 text-center text-xs text-slate-500 border-t border-slate-800">
        © {new Date().getFullYear()} LuckyMart. Hak Cipta Dilindungi Usaha Keluarga.
      </div>
    </footer>
  );
};
