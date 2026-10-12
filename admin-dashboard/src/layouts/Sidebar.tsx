import React, { useState } from 'react';
import { NavLink, useLocation } from 'react-router-dom';
import { useAuthStore } from '../store/useAuthStore';
import { ROLE_LABELS } from '../lib/utils';
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
  ChevronDown,
} from 'lucide-react';

interface SidebarItem {
  title: string;
  href: string;
  icon: React.ElementType;
  roles: string[];
  /** Less-used pages live under "Lainnya" to keep the menu short. */
  more?: boolean;
}

const ALL_STAFF = ['SUPER_ADMIN', 'STORE_MANAGER', 'WAREHOUSE_STAFF', 'DRIVER'];
const WAREHOUSE = ['SUPER_ADMIN', 'STORE_MANAGER', 'WAREHOUSE_STAFF'];
const MANAGERS = ['SUPER_ADMIN', 'STORE_MANAGER'];

export const menuItems: SidebarItem[] = [
  { title: 'Ringkasan', href: '/dashboard', icon: LayoutDashboard, roles: WAREHOUSE },
  { title: 'Pesanan', href: '/orders', icon: ShoppingCart, roles: ALL_STAFF },
  { title: 'Produk', href: '/products', icon: Package, roles: WAREHOUSE },
  { title: 'Stok Barang', href: '/inventory', icon: Boxes, roles: WAREHOUSE },
  { title: 'Kasbon Mitra', href: '/credit-ledgers', icon: Store, roles: MANAGERS },
  { title: 'Pelanggan & Staf', href: '/users', icon: Users, roles: MANAGERS },
  { title: 'Pengaturan', href: '/settings', icon: Settings, roles: ALL_STAFF },
  { title: 'Kategori', href: '/categories', icon: Layers, roles: WAREHOUSE, more: true },
  { title: 'Pembayaran', href: '/payments', icon: CreditCard, roles: MANAGERS, more: true },
  { title: 'Pengiriman', href: '/shipping', icon: Truck, roles: ALL_STAFF, more: true },
  { title: 'Mutasi Stok', href: '/stock-movements', icon: FileSpreadsheet, roles: WAREHOUSE, more: true },
  { title: 'Notifikasi', href: '/notifications', icon: Bell, roles: ALL_STAFF, more: true },
  { title: 'Log Aktivitas', href: '/events', icon: ShieldCheck, roles: MANAGERS, more: true },
];

export const Sidebar: React.FC<{ isOpen: boolean; onClose: () => void }> = ({ isOpen, onClose }) => {
  const { user } = useAuthStore();
  const currentRole = user?.role || '';

  const { pathname } = useLocation();
  const filteredMenu = menuItems.filter((item) => item.roles.includes(currentRole));
  const mainMenu = filteredMenu.filter((item) => !item.more);
  const moreMenu = filteredMenu.filter((item) => item.more);
  const [showMore, setShowMore] = useState(() => moreMenu.some((item) => pathname.startsWith(item.href)));

  const renderLink = (item: SidebarItem) => {
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
  };

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
        className={`fixed top-0 left-0 z-50 flex h-screen w-64 flex-col transform bg-white dark:bg-slate-900 border-r border-slate-200 dark:border-slate-800 transition-transform duration-200 ease-in-out lg:translate-x-0 ${
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
          {mainMenu.map(renderLink)}
          {moreMenu.length > 0 && (
            <>
              <button
                onClick={() => setShowMore((v) => !v)}
                className="mt-3 flex w-full items-center justify-between rounded-lg px-4 py-2 text-xs font-semibold uppercase tracking-wider text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800"
              >
                Lainnya
                <ChevronDown className={`h-4 w-4 transition-transform ${showMore ? 'rotate-180' : ''}`} />
              </button>
              {showMore && moreMenu.map(renderLink)}
            </>
          )}
        </nav>

        {/* User Role Badge Footer */}
        <div className="p-4 border-t border-slate-200 dark:border-slate-800">
          <div className="rounded-lg bg-slate-100 dark:bg-slate-800 p-3">
            <div className="text-xs text-slate-500 dark:text-slate-400">Tersambung sebagai:</div>
            <div className="text-sm font-bold text-slate-900 dark:text-white truncate">
              {user?.full_name}
            </div>
            <span className="inline-block mt-1 rounded bg-brand-100 dark:bg-brand-900/50 px-2 py-0.5 text-[10px] font-bold text-brand-700 dark:text-brand-300">
              {user ? ROLE_LABELS[user.role] : ''}
            </span>
          </div>
        </div>
      </aside>
    </>
  );
};
