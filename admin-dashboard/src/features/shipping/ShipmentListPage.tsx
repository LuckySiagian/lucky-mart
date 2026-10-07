import React, { useState } from 'react';
import { PageHeader } from '../../components/common/PageHeader';
import { Truck, MapPin, Search, UserCheck, CheckCircle2, Clock, Navigation, X, RefreshCw, Plus, Trash2 } from 'lucide-react';
import { toast } from 'sonner';

interface Shipment {
  id: string;
  trackingNumber: string;
  orderNumber: string;
  customerName: string;
  address: string;
  courierName: string;
  driverName: string;
  driverPhone: string;
  vehicle: string;
  lastLocation: string;
  status: 'READY_TO_SHIP' | 'IN_TRANSIT' | 'DELIVERED';
  timeline: { time: string; location: string; note: string }[];
}

export const ShipmentListPage: React.FC = () => {
  const [shipments, setShipments] = useState<Shipment[]>([]);
  const [search, setSearch] = useState('');

  // Timeline Modal State
  const [selectedShipment, setSelectedShipment] = useState<Shipment | null>(null);

  // Create Shipment Modal State
  const [isAddOpen, setIsAddOpen] = useState(false);
  const [orderNum, setOrderNum] = useState('');
  const [custName, setCustName] = useState('');
  const [address, setAddress] = useState('');
  const [driverName, setDriverName] = useState('Pak Joko Widodo (Driver Toko)');
  const [vehicle, setVehicle] = useState('PickUp L300 (B 9821 TKG)');

  const handleUpdateStatus = (id: string, newStatus: 'IN_TRANSIT' | 'DELIVERED') => {
    setShipments((prev) =>
      prev.map((s) => {
        if (s.id === id) {
          const updatedTimeline = [
            ...s.timeline,
            {
              time: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
              location: newStatus === 'DELIVERED' ? 'Tujuan Pembeli' : 'Lokasi Rute Pengiriman',
              note: newStatus === 'DELIVERED' ? 'Barang sembako berhasil diterima pelanggan' : 'Kurir bergerak ke lokasi berikutnya',
            },
          ];
          return {
            ...s,
            status: newStatus,
            lastLocation: newStatus === 'DELIVERED' ? 'Sudah Diterima Pelanggan' : s.lastLocation,
            timeline: updatedTimeline,
          };
        }
        return s;
      })
    );
    toast.success(`Status pengiriman resi "${id}" berhasil diperbarui!`);
    setSelectedShipment(null);
  };

  const handleDeleteShipment = (id: string, trackingNum: string) => {
    if (window.confirm(`Hapus data pengiriman resi "${trackingNum}"?`)) {
      setShipments((prev) => prev.filter((s) => s.id !== id));
      toast.success(`Data pengiriman resi "${trackingNum}" berhasil dihapus!`);
      if (selectedShipment?.id === id) setSelectedShipment(null);
    }
  };

  const handleCreateShipment = (e: React.FormEvent) => {
    e.preventDefault();
    if (!custName.trim() || !address.trim()) {
      toast.error('Nama pelanggan dan alamat pengiriman wajib diisi.');
      return;
    }

    const resiNum = `RESI-${new Date().toISOString().slice(0, 10).replace(/-/g, '')}-${Math.floor(1000 + Math.random() * 9000)}`;
    const newShip: Shipment = {
      id: `SHIP-${Date.now()}`,
      trackingNumber: resiNum,
      orderNumber: orderNum.trim() || `ORD-${Math.floor(1000 + Math.random() * 9000)}`,
      customerName: custName.trim(),
      address: address.trim(),
      courierName: 'Kurir Toko Internal',
      driverName: driverName,
      driverPhone: '0812-9988-7766',
      vehicle: vehicle,
      lastLocation: 'Gudang Pusat Cengkareng (Penugasan Kurir)',
      status: 'IN_TRANSIT',
      timeline: [
        {
          time: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
          location: 'Gudang Pusat Cengkareng',
          note: 'Barang dikemas dan diserahkan ke driver armada',
        },
      ],
    };

    setShipments((prev) => [newShip, ...prev]);
    toast.success(`Resi pengiriman baru "${resiNum}" berhasil dibuat!`);
    setIsAddOpen(false);
    setCustName('');
    setAddress('');
  };

  const filtered = shipments.filter(
    (s) =>
      s.trackingNumber.toLowerCase().includes(search.toLowerCase()) ||
      s.driverName.toLowerCase().includes(search.toLowerCase()) ||
      s.customerName.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className="space-y-6">
      <PageHeader
        title="Pengiriman & Pelacakan Resi (Shipments)"
        subtitle="Pantau lokasi pengiriman barang sembako, penugasan kurir toko & driver armada, serta update status resi."
      >
        <button
          onClick={() => setIsAddOpen(true)}
          className="inline-flex items-center space-x-2 rounded-xl bg-brand-600 px-4 py-2.5 text-sm font-semibold text-white hover:bg-brand-700 transition-colors shadow-lg shadow-brand-600/20"
        >
          <Plus className="h-4 w-4" />
          <span>Buat Resi Pengiriman Baru</span>
        </button>
      </PageHeader>

      {/* Filter & Search Bar */}
      <div className="flex flex-col sm:flex-row items-center justify-between gap-4 bg-white dark:bg-slate-900 p-4 rounded-2xl border border-slate-200 dark:border-slate-800">
        <div className="relative w-full sm:w-80">
          <Search className="absolute left-3.5 top-3 h-4 w-4 text-slate-400" />
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Cari no. resi, nama kurir/driver..."
            className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 pl-10 pr-4 py-2 text-sm text-slate-900 dark:text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-brand-500"
          />
        </div>
        <div className="text-xs text-slate-500 font-semibold">
          Dalam Perjalanan: <span className="text-brand-600 font-bold">{shipments.filter((s) => s.status === 'IN_TRANSIT').length} Resi Active</span>
        </div>
      </div>

      {/* Shipment Table */}
      <div className="overflow-x-auto rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 shadow-sm">
        <table className="w-full text-left text-sm text-slate-600 dark:text-slate-400">
          <thead className="bg-slate-50 dark:bg-slate-800/50 text-xs uppercase text-slate-500 font-bold border-b border-slate-200 dark:border-slate-800">
            <tr>
              <th className="px-6 py-4">Nomor Resi & Order</th>
              <th className="px-6 py-4">Kurir & Driver Pembawa</th>
              <th className="px-6 py-4">Tujuan Pengiriman</th>
              <th className="px-6 py-4">Lokasi Pelacakan Terakhir</th>
              <th className="px-6 py-4">Status</th>
              <th className="px-6 py-4 text-right">Aksi Kelola (CRUD)</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-200 dark:divide-slate-800">
            {filtered.map((s) => (
              <tr key={s.id} className="hover:bg-slate-50/80 dark:hover:bg-slate-800/30">
                <td className="px-6 py-4">
                  <div className="font-mono font-bold text-slate-900 dark:text-white flex items-center space-x-2">
                    <Truck className="h-4 w-4 text-brand-600" />
                    <span>{s.trackingNumber}</span>
                  </div>
                  <div className="text-xs font-mono text-slate-400">{s.orderNumber}</div>
                </td>
                <td className="px-6 py-4">
                  <div className="font-bold text-slate-900 dark:text-white flex items-center space-x-1.5">
                    <UserCheck className="h-4 w-4 text-brand-600" />
                    <span>{s.driverName}</span>
                  </div>
                  <div className="text-xs text-slate-400">{s.vehicle} ({s.driverPhone})</div>
                </td>
                <td className="px-6 py-4 text-xs font-medium text-slate-700 dark:text-slate-300">
                  <div className="font-bold text-slate-900 dark:text-white">{s.customerName}</div>
                  <div className="text-slate-400 max-w-xs truncate">{s.address}</div>
                </td>
                <td className="px-6 py-4 text-xs font-semibold">
                  <div className="flex items-center space-x-1.5 text-brand-600 dark:text-brand-400">
                    <MapPin className="h-3.5 w-3.5 flex-shrink-0 text-rose-500" />
                    <span>{s.lastLocation}</span>
                  </div>
                </td>
                <td className="px-6 py-4">
                  {s.status === 'IN_TRANSIT' && (
                    <span className="inline-flex items-center space-x-1 rounded-full bg-blue-100 dark:bg-blue-950/50 px-2.5 py-0.5 text-xs font-bold text-blue-700 dark:text-blue-400">
                      <Navigation className="h-3.5 w-3.5" />
                      <span>DALAM PERJALANAN</span>
                    </span>
                  )}
                  {s.status === 'DELIVERED' && (
                    <span className="inline-flex items-center space-x-1 rounded-full bg-emerald-100 dark:bg-emerald-950/50 px-2.5 py-0.5 text-xs font-bold text-emerald-700 dark:text-emerald-400">
                      <CheckCircle2 className="h-3.5 w-3.5" />
                      <span>TERKIRIM (DELIVERED)</span>
                    </span>
                  )}
                </td>
                <td className="px-6 py-4 text-right space-x-2">
                  <button
                    onClick={() => setSelectedShipment(s)}
                    className="inline-flex items-center space-x-1 rounded-lg bg-brand-600 px-3 py-1.5 text-xs font-bold text-white hover:bg-brand-700 shadow-md shadow-brand-600/20"
                  >
                    <MapPin className="h-3.5 w-3.5" />
                    <span>Lacak Rute</span>
                  </button>
                  <button
                    onClick={() => handleDeleteShipment(s.id, s.trackingNumber)}
                    className="inline-flex items-center space-x-1 rounded-lg bg-rose-100 text-rose-700 dark:bg-rose-950/50 dark:text-rose-400 px-2.5 py-1.5 text-xs font-bold hover:bg-rose-200"
                  >
                    <Trash2 className="h-3.5 w-3.5" />
                    <span>Hapus</span>
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Modal Live Tracking Timeline */}
      {selectedShipment && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 backdrop-blur-sm p-4">
          <div className="relative w-full max-w-lg rounded-3xl bg-white dark:bg-slate-900 p-6 shadow-2xl border border-slate-200 dark:border-slate-800 space-y-6">
            <div className="flex items-center justify-between border-b border-slate-200 dark:border-slate-800 pb-4">
              <div className="flex items-center space-x-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-brand-100 text-brand-600">
                  <Truck className="h-5 w-5" />
                </div>
                <div>
                  <h3 className="text-lg font-bold text-slate-900 dark:text-white">Pelacakan Rute Resi</h3>
                  <p className="text-xs text-slate-500">{selectedShipment.trackingNumber}</p>
                </div>
              </div>
              <button onClick={() => setSelectedShipment(null)} className="rounded-xl p-2 text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800">
                <X className="h-5 w-5" />
              </button>
            </div>

            <div className="rounded-2xl bg-slate-100 dark:bg-slate-800 p-4 space-y-2 text-xs">
              <div className="flex justify-between">
                <span className="text-slate-400">Driver Toko:</span>
                <span className="font-bold text-slate-900 dark:text-white">{selectedShipment.driverName}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-slate-400">Armada Kendaraan:</span>
                <span className="font-bold text-slate-900 dark:text-white">{selectedShipment.vehicle}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-slate-400">Penerima Sembako:</span>
                <span className="font-bold text-slate-900 dark:text-white">{selectedShipment.customerName} ({selectedShipment.address})</span>
              </div>
            </div>

            {/* Timeline */}
            <div className="space-y-4 relative pl-4 border-l-2 border-brand-500">
              {selectedShipment.timeline.map((item, idx) => (
                <div key={idx} className="relative">
                  <div className="absolute -left-[21px] top-1 h-3 w-3 rounded-full bg-brand-600 ring-4 ring-white dark:ring-slate-900" />
                  <div className="text-xs font-bold text-slate-900 dark:text-white flex items-center space-x-2">
                    <span className="text-brand-600 font-mono">{item.time}</span>
                    <span>• {item.location}</span>
                  </div>
                  <div className="text-xs text-slate-500 mt-0.5">{item.note}</div>
                </div>
              ))}
            </div>

            <div className="flex items-center justify-between pt-4 border-t border-slate-200 dark:border-slate-800">
              {selectedShipment.status === 'IN_TRANSIT' ? (
                <button
                  type="button"
                  onClick={() => handleUpdateStatus(selectedShipment.id, 'DELIVERED')}
                  className="w-full rounded-xl bg-emerald-600 px-5 py-2.5 text-xs font-bold text-white hover:bg-emerald-700 shadow-lg shadow-emerald-600/20"
                >
                  Tandai Barang Sudah Sampai (DELIVERED)
                </button>
              ) : (
                <button
                  type="button"
                  onClick={() => setSelectedShipment(null)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-4 py-2.5 text-xs font-bold text-slate-700 dark:text-slate-200"
                >
                  Tutup Pelacakan
                </button>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Modal Buat Resi Pengiriman Baru */}
      {isAddOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 backdrop-blur-sm p-4">
          <div className="relative w-full max-w-md rounded-3xl bg-white dark:bg-slate-900 p-6 shadow-2xl border border-slate-200 dark:border-slate-800 space-y-6">
            <div className="flex items-center justify-between border-b border-slate-200 dark:border-slate-800 pb-4">
              <div className="flex items-center space-x-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-brand-100 text-brand-600">
                  <Plus className="h-5 w-5" />
                </div>
                <div>
                  <h3 className="text-lg font-bold text-slate-900 dark:text-white">Buat Resi Pengiriman Baru</h3>
                  <p className="text-xs text-slate-500">Penugasan driver armada toko untuk pengantaran barang sembako.</p>
                </div>
              </div>
              <button onClick={() => setIsAddOpen(false)} className="rounded-xl p-2 text-slate-400 hover:bg-slate-100">
                <X className="h-5 w-5" />
              </button>
            </div>

            <form onSubmit={handleCreateShipment} className="space-y-4">
              <div>
                <label className="block text-xs font-bold text-slate-500 mb-1">Nomor Order Terkait</label>
                <input
                  type="text"
                  placeholder="ORD-20260727-xxxx"
                  value={orderNum}
                  onChange={(e) => setOrderNum(e.target.value)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white font-mono"
                />
              </div>

              <div>
                <label className="block text-xs font-bold text-slate-500 mb-1">Nama Penerima / Pelanggan *</label>
                <input
                  type="text"
                  required
                  placeholder="Contoh: Budi Santoso"
                  value={custName}
                  onChange={(e) => setCustName(e.target.value)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white font-bold"
                />
              </div>

              <div>
                <label className="block text-xs font-bold text-slate-500 mb-1">Alamat Tujuan Pengiriman *</label>
                <textarea
                  rows={2}
                  required
                  placeholder="Alamat lengkap tujuan pengantaran..."
                  value={address}
                  onChange={(e) => setAddress(e.target.value)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-sm text-slate-900 dark:text-white resize-none"
                />
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs font-bold text-slate-500 mb-1">Driver Armada</label>
                  <select
                    value={driverName}
                    onChange={(e) => setDriverName(e.target.value)}
                    className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-xs text-slate-900 dark:text-white font-bold"
                  >
                    <option value="Pak Joko Widodo (Driver Toko)">Pak Joko Widodo</option>
                    <option value="Mas Bambang S. (Driver Truk)">Mas Bambang S.</option>
                    <option value="Kurir Express Same-Day">Kurir Express Same-Day</option>
                  </select>
                </div>

                <div>
                  <label className="block text-xs font-bold text-slate-500 mb-1">Armada Kendaraan</label>
                  <select
                    value={vehicle}
                    onChange={(e) => setVehicle(e.target.value)}
                    className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-3.5 py-2 text-xs text-slate-900 dark:text-white font-bold"
                  >
                    <option value="PickUp L300 (B 9821 TKG)">PickUp L300 (B 9821 TKG)</option>
                    <option value="Truk Engkel Box (B 9012 SMB)">Truk Engkel (B 9012 SMB)</option>
                    <option value="Motor Kurir Toko (B 4410 SMB)">Motor Kurir (B 4410 SMB)</option>
                  </select>
                </div>
              </div>

              <div className="flex items-center justify-end space-x-3 pt-4 border-t border-slate-200 dark:border-slate-800">
                <button
                  type="button"
                  onClick={() => setIsAddOpen(false)}
                  className="rounded-xl px-4 py-2.5 text-xs font-bold text-slate-600 hover:bg-slate-100"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  className="rounded-xl bg-brand-600 px-5 py-2.5 text-xs font-bold text-white hover:bg-brand-700 shadow-lg shadow-brand-600/20"
                >
                  Simpan & Terbitkan Resi
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
