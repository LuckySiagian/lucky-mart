import React from 'react';
import { NavLink } from 'react-router-dom';
import { useAuthStore } from '../store/useAuthStore';
import {
  LayoutDashboard,
  Package,
  Boxes,
  ShoppingCart,
  CreditCard,
  Truck,
  Bell,
  Users,
  Settings,
  Store,
  Layers,
  FileSpreadsheet,
  ShieldCheck,
} from 'lucide-react';

interface SidebarItem {
  title: string;
  href: string;
  icon: React.ElementType;
  roles?: string[];
}

const menuItems: SidebarItem[] = [
  { title: 'Dashboard', href: '/dashboard', icon: LayoutDashboard },
  { title: 'Produk Catalog', href: '/products', icon: Package },
  { title: 'Kategori Produk', href: '/categories', icon: Layers },
  { title: 'Manajemen Stok', href: '/inventory', icon: Boxes },
  { title: 'Mutasi Stok Log', href: '/stock-movements', icon: FileSpreadsheet },
  { title: 'Pesanan (Orders)', href: '/orders', icon: ShoppingCart },
  { title: 'Pembayaran (Payments)', href: '/payments', icon: CreditCard },
  { title: 'B2B Credit Ledger', href: '/credit-ledgers', icon: Store, roles: ['SUPER_ADMIN', 'STORE_MANAGER'] },
  { title: 'Pengiriman & Resi', href: '/shipping', icon: Truck },
  { title: 'Notifikasi System', href: '/notifications', icon: Bell },
  { title: 'Pengguna & Role', href: '/users', icon: Users, roles: ['SUPER_ADMIN', 'STORE_MANAGER'] },
  { title: 'Audit Events Log', href: '/events', icon: ShieldCheck, roles: ['SUPER_ADMIN'] },
  { title: 'Pengaturan Sistem', href: '/settings', icon: Settings },
];

export const Sidebar: React.FC<{ isOpen: boolean; onClose: () => void }> = ({ isOpen, onClose }) => {
  const { user } = useAuthStore();
  const currentRole = user?.role || 'B2C_CUSTOMER';

  const filteredMenu = menuItems.filter(
    (item) => !item.roles || item.roles.includes(currentRole)
  );

  return (
    <>
      {/* Mobile Backdrop */}
      {isOpen && (
        <div
          className="fixed inset-0 z-40 bg-slate-900/50 lg:hidden"
          onClick={onClose}
        />
      )}

      <aside
        className={`fixed top-0 left-0 z-50 h-screen w-64 transform bg-white dark:bg-slate-900 border-r border-slate-200 dark:border-slate-800 transition-transform duration-200 ease-in-out lg:translate-x-0 ${
          isOpen ? 'translate-x-0' : '-translate-x-full'
        }`}
      >
        {/* Brand Header */}
        <div className="flex h-16 items-center px-6 border-b border-slate-200 dark:border-slate-800">
          <div className="flex items-center space-x-3">
            <img
              src="/logo.png"
              alt="LuckyMart Logo"
              className="h-10 w-auto object-contain"
            />
            <div>
              <h1 className="text-base font-bold text-slate-900 dark:text-white leading-none">
                LUCKYMART
              </h1>
              <span className="text-xs text-brand-600 font-semibold tracking-wider uppercase">
                Admin Panel
              </span>
            </div>
          </div>
        </div>

        {/* Navigation Menu */}
        <nav className="flex-1 overflow-y-auto p-4 space-y-1">
          {filteredMenu.map((item) => {
            const Icon = item.icon;
            return (
              <NavLink
                key={item.href}
                to={item.href}
                onClick={onClose}
                className={({ isActive }) =>
                  `flex items-center space-x-3 rounded-lg px-4 py-2.5 text-sm font-medium transition-colors ${
                    isActive
                      ? 'bg-brand-50 text-brand-700 dark:bg-brand-950/40 dark:text-brand-400 font-semibold'
                      : 'text-slate-600 hover:bg-slate-100 dark:text-slate-400 dark:hover:bg-slate-800'
                  }`
                }
              >
                <Icon className="h-5 w-5 flex-shrink-0" />
                <span>{item.title}</span>
              </NavLink>
            );
          })}
        </nav>

        {/* User Role Badge Footer */}
        <div className="p-4 border-t border-slate-200 dark:border-slate-800">
          <div className="rounded-lg bg-slate-100 dark:bg-slate-800 p-3">
            <div className="text-xs text-slate-500 dark:text-slate-400">Tersambung sebagai:</div>
            <div className="text-sm font-bold text-slate-900 dark:text-white truncate">
              {user?.full_name || 'Admin LuckyMart'}
            </div>
            <span className="inline-block mt-1 rounded bg-brand-100 dark:bg-brand-900/50 px-2 py-0.5 text-[10px] font-bold text-brand-700 dark:text-brand-300">
              {user?.role || 'SUPER_ADMIN'}
            </span>
          </div>
        </div>
      </aside>
    </>
  );
};
