import React from 'react';
import { Link } from 'react-router-dom';
import { Clock, MapPin, MessageCircle, Store, Tags, Truck } from 'lucide-react';
import { useStoreInfo } from '../../hooks/useStoreInfo';
import { whatsappLink } from '../../lib/utils';

const HIGHLIGHTS = [
  { icon: Truck, title: 'Diantar ke rumah', text: 'Pesan dari HP, barang diantar kurir toko.' },
  { icon: Store, title: 'Ambil sendiri di toko', text: 'Pesan dulu, ambil saat siap. Gratis ongkir.' },
  { icon: Tags, title: 'Harga grosir untuk warung', text: 'Daftar sebagai mitra toko untuk harga grosir & kasbon.' },
];

export const Footer: React.FC = () => {
  const store = useStoreInfo();
  const wa = whatsappLink(store?.whatsapp, 'Halo, saya mau tanya soal belanja di toko.');

  return (
    <footer className="mt-16 border-t border-slate-800 bg-slate-900 text-sm text-slate-400">
      <div className="mx-auto max-w-7xl border-b border-slate-800 px-4 py-10 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 gap-8 md:grid-cols-3">
          {HIGHLIGHTS.map((h) => (
            <div key={h.title} className="flex items-center gap-4">
              <div className="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-2xl bg-primary-600/20 text-primary-400">
                <h.icon className="h-6 w-6" />
              </div>
              <div>
                <h4 className="font-bold text-white">{h.title}</h4>
                <p className="mt-0.5 text-xs text-slate-400">{h.text}</p>
              </div>
            </div>
          ))}
        </div>
      </div>

      <div className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 gap-8 md:grid-cols-3">
          <div>
            <div className="mb-3 flex items-center gap-2">
              <div className="flex h-8 w-8 items-center justify-center rounded-xl bg-primary-600 text-white">
                <Store className="h-5 w-5" />
              </div>
              <span className="text-lg font-black text-white">{store?.name || 'Lucky Mart'}</span>
            </div>
            {store?.tagline && <p className="text-xs leading-relaxed">{store.tagline}</p>}
          </div>

          <div>
            <h5 className="mb-3 font-bold text-white">Belanja</h5>
            <ul className="space-y-2 text-xs">
              <li><Link to="/products" className="hover:text-primary-400">Semua produk</Link></li>
              <li><Link to="/orders" className="hover:text-primary-400">Pesanan saya</Link></li>
              <li><Link to="/register" className="hover:text-primary-400">Daftar akun / mitra warung</Link></li>
            </ul>
          </div>

          <div>
            <h5 className="mb-3 font-bold text-white">Hubungi Toko</h5>
            <div className="space-y-2.5 text-xs">
              {store?.address && (
                <div className="flex items-start gap-2">
                  <MapPin className="mt-0.5 h-4 w-4 flex-shrink-0 text-primary-400" />
                  {store.maps_url ? (
                    <a href={store.maps_url} target="_blank" rel="noreferrer" className="hover:text-primary-400">{store.address}</a>
                  ) : (
                    <span>{store.address}</span>
                  )}
                </div>
              )}
              {store?.open_hours && (
                <div className="flex items-center gap-2">
                  <Clock className="h-4 w-4 flex-shrink-0 text-primary-400" /> {store.open_hours}
                </div>
              )}
              {wa && (
                <a href={wa} target="_blank" rel="noreferrer" className="inline-flex items-center gap-2 rounded-xl bg-emerald-600 px-3 py-2 font-bold text-white hover:bg-emerald-700">
                  <MessageCircle className="h-4 w-4" /> Chat WhatsApp toko
                </a>
              )}
            </div>
          </div>
        </div>
      </div>

      <div className="border-t border-slate-800 bg-slate-950 py-4 text-center text-xs text-slate-500">
        © {new Date().getFullYear()} {store?.name || 'Lucky Mart'}
      </div>
    </footer>
  );
};
