import React, { useState } from 'react';
import { PageHeader } from '../../components/common/PageHeader';
import { useAuthStore } from '../../store/useAuthStore';
import { Settings, Save, Shield, Store, CreditCard, Clock, Lock } from 'lucide-react';
import { toast } from 'sonner';

export const SettingsPage: React.FC = () => {
  const { user } = useAuthStore();
  const [activeTab, setActiveTab] = useState<'STORE' | 'PAYMENT' | 'B2B' | 'PROFILE'>('STORE');

  const [storeName, setStoreName] = useState('LuckyMart');
  const [storeAddress, setStoreAddress] = useState('Jl. Daan Mogot Km 14 No. 88, Cengkareng, Jakarta Barat');
  const [openHours, setOpenHours] = useState('07:00 - 20:00 WIB');
  const [csPhone, setCsPhone] = useState('0812-9900-1122');

  const [mandiriAccount, setMandiriAccount] = useState('137-00-99881-22 (a.n PT LuckyMart Indonesia)');
  const [qrisMerchantId, setQrisMerchantId] = useState('NMID-ID1029384756');

  const [defaultB2BLimit, setDefaultB2BLimit] = useState('50000000');
  const [tempoDays, setTempoDays] = useState('30');

  const handleSave = (e: React.FormEvent) => {
    e.preventDefault();
    toast.success('Pengaturan sistem & profil berhasil disimpan dan diperbarui!');
  };

  return (
    <div className="space-y-6">
      <PageHeader
        title="Pengaturan Sistem & Profil Toko"
        subtitle="Pusat konfigurasi profil usaha sembako, rekening pembayaran, aturan B2B Kasbon Tempo, dan akun admin."
      />

      {/* Tabs Navigation */}
      <div className="flex items-center space-x-2 border-b border-slate-200 dark:border-slate-800 pb-2">
        <button
          onClick={() => setActiveTab('STORE')}
          className={`flex items-center space-x-2 px-4 py-2.5 rounded-xl text-xs font-bold transition-all ${
            activeTab === 'STORE'
              ? 'bg-brand-600 text-white shadow-md shadow-brand-600/20'
              : 'bg-white dark:bg-slate-900 text-slate-600 dark:text-slate-300 hover:bg-slate-100'
          }`}
        >
          <Store className="h-4 w-4" />
          <span>Profil Toko & Usaha</span>
        </button>

        <button
          onClick={() => setActiveTab('PAYMENT')}
          className={`flex items-center space-x-2 px-4 py-2.5 rounded-xl text-xs font-bold transition-all ${
            activeTab === 'PAYMENT'
              ? 'bg-brand-600 text-white shadow-md shadow-brand-600/20'
              : 'bg-white dark:bg-slate-900 text-slate-600 dark:text-slate-300 hover:bg-slate-100'
          }`}
        >
          <CreditCard className="h-4 w-4" />
          <span>Rekening & Payment Gateway</span>
        </button>

        <button
          onClick={() => setActiveTab('B2B')}
          className={`flex items-center space-x-2 px-4 py-2.5 rounded-xl text-xs font-bold transition-all ${
            activeTab === 'B2B'
              ? 'bg-brand-600 text-white shadow-md shadow-brand-600/20'
              : 'bg-white dark:bg-slate-900 text-slate-600 dark:text-slate-300 hover:bg-slate-100'
          }`}
        >
          <Clock className="h-4 w-4" />
          <span>Kebijakan B2B Kasbon Tempo</span>
        </button>

        <button
          onClick={() => setActiveTab('PROFILE')}
          className={`flex items-center space-x-2 px-4 py-2.5 rounded-xl text-xs font-bold transition-all ${
            activeTab === 'PROFILE'
              ? 'bg-brand-600 text-white shadow-md shadow-brand-600/20'
              : 'bg-white dark:bg-slate-900 text-slate-600 dark:text-slate-300 hover:bg-slate-100'
          }`}
        >
          <Shield className="h-4 w-4" />
          <span>Akun Saya & Keamanan</span>
        </button>
      </div>

      <form onSubmit={handleSave} className="space-y-6 max-w-3xl">
        {activeTab === 'STORE' && (
          <div className="rounded-3xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 p-6 space-y-4 shadow-sm">
            <h3 className="text-sm font-bold text-slate-900 dark:text-white flex items-center space-x-2">
              <Store className="h-5 w-5 text-brand-600" />
              <span>Informasi Identitas Usaha Sembako</span>
            </h3>

            <div>
              <label className="block text-xs font-bold text-slate-700 dark:text-slate-300 mb-1">Nama Platform / Toko *</label>
              <input
                type="text"
                value={storeName}
                onChange={(e) => setStoreName(e.target.value)}
                className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-4 py-2.5 text-sm text-slate-900 dark:text-white font-bold"
              />
            </div>

            <div>
              <label className="block text-xs font-bold text-slate-700 dark:text-slate-300 mb-1">Alamat Gudang Utama *</label>
              <textarea
                rows={2}
                value={storeAddress}
                onChange={(e) => setStoreAddress(e.target.value)}
                className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-4 py-2.5 text-sm text-slate-900 dark:text-white resize-none"
              />
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label className="block text-xs font-bold text-slate-700 dark:text-slate-300 mb-1">Jam Operasional Pengiriman *</label>
                <input
                  type="text"
                  value={openHours}
                  onChange={(e) => setOpenHours(e.target.value)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-4 py-2.5 text-sm text-slate-900 dark:text-white"
                />
              </div>

              <div>
                <label className="block text-xs font-bold text-slate-700 dark:text-slate-300 mb-1">WhatsApp Customer Care CS *</label>
                <input
                  type="text"
                  value={csPhone}
                  onChange={(e) => setCsPhone(e.target.value)}
                  className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-4 py-2.5 text-sm text-slate-900 dark:text-white font-mono"
                />
              </div>
            </div>
          </div>
        )}

        {activeTab === 'PAYMENT' && (
          <div className="rounded-3xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 p-6 space-y-4 shadow-sm">
            <h3 className="text-sm font-bold text-slate-900 dark:text-white flex items-center space-x-2">
              <CreditCard className="h-5 w-5 text-brand-600" />
              <span>Rekening Pembayaran & Merchant QRIS</span>
            </h3>

            <div>
              <label className="block text-xs font-bold text-slate-700 dark:text-slate-300 mb-1">Rekening Bank Mandiri *</label>
              <input
                type="text"
                value={mandiriAccount}
                onChange={(e) => setMandiriAccount(e.target.value)}
                className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-4 py-2.5 text-sm text-slate-900 dark:text-white font-mono"
              />
            </div>

            <div>
              <label className="block text-xs font-bold text-slate-700 dark:text-slate-300 mb-1">Static QRIS Merchant ID *</label>
              <input
                type="text"
                value={qrisMerchantId}
                onChange={(e) => setQrisMerchantId(e.target.value)}
                className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-4 py-2.5 text-sm text-slate-900 dark:text-white font-mono"
              />
            </div>
          </div>
        )}

        {activeTab === 'B2B' && (
          <div className="rounded-3xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 p-6 space-y-4 shadow-sm">
            <h3 className="text-sm font-bold text-slate-900 dark:text-white flex items-center space-x-2">
              <Clock className="h-5 w-5 text-brand-600" />
              <span>Kebijakan Fasilitas B2B Kasbon Tempo</span>
            </h3>

            <div>
              <label className="block text-xs font-bold text-slate-700 dark:text-slate-300 mb-1">Default Plafon Limit Kredit Baru (Rp) *</label>
              <input
                type="number"
                value={defaultB2BLimit}
                onChange={(e) => setDefaultB2BLimit(e.target.value)}
                className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-4 py-2.5 text-sm text-slate-900 dark:text-white font-bold"
              />
            </div>

            <div>
              <label className="block text-xs font-bold text-slate-700 dark:text-slate-300 mb-1">Durasi Jatuh Tempo Standar (Hari) *</label>
              <input
                type="number"
                value={tempoDays}
                onChange={(e) => setTempoDays(e.target.value)}
                className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-4 py-2.5 text-sm text-slate-900 dark:text-white font-bold"
              />
            </div>
          </div>
        )}

        {activeTab === 'PROFILE' && (
          <div className="rounded-3xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 p-6 space-y-4 shadow-sm">
            <h3 className="text-sm font-bold text-slate-900 dark:text-white flex items-center space-x-2">
              <Shield className="h-5 w-5 text-brand-600" />
              <span>Profil Pengguna & Password Admin</span>
            </h3>

            <div>
              <label className="block text-xs font-bold text-slate-700 dark:text-slate-300 mb-1">Nama Lengkap Admin *</label>
              <input
                type="text"
                defaultValue={user?.full_name || 'Super Admin Sembako'}
                className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-4 py-2.5 text-sm text-slate-900 dark:text-white font-bold"
              />
            </div>

            <div>
              <label className="block text-xs font-bold text-slate-700 dark:text-slate-300 mb-1">Email Utama Login Admin *</label>
              <input
                type="email"
                defaultValue={user?.email || 'superadmin@sembako.com'}
                className="w-full rounded-xl bg-slate-100 dark:bg-slate-800 px-4 py-2.5 text-sm text-slate-900 dark:text-white font-mono"
              />
            </div>
          </div>
        )}

        <button
          type="submit"
          className="inline-flex items-center space-x-2 rounded-xl bg-brand-600 px-6 py-3 text-xs font-bold text-white hover:bg-brand-700 transition-colors shadow-lg shadow-brand-600/30"
        >
          <Save className="h-4 w-4" />
          <span>Simpan Perubahan Pengaturan</span>
        </button>
      </form>
    </div>
  );
};
