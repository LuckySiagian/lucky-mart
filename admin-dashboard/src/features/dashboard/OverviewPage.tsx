import React from 'react';
import { PageHeader } from '../../components/common/PageHeader';
import { StatusBadge } from '../../components/common/StatusBadge';
import { formatCurrency } from '../../lib/utils';
import { productService, orderService } from '../../services/apiServices';
import {
  DollarSign,
  ShoppingCart,
  TrendingUp,
  AlertTriangle,
  Users,
  CheckCircle2,
  Clock,
  ArrowUpRight,
} from 'lucide-react';
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  BarChart,
  Bar,
} from 'recharts';

const salesData = [
  { name: 'Senin', sales: 12500000, orders: 34 },
  { name: 'Selasa', sales: 18200000, orders: 48 },
  { name: 'Rabu', sales: 15400000, orders: 41 },
  { name: 'Kamis', sales: 22800000, orders: 62 },
  { name: 'Jumat', sales: 29500000, orders: 85 },
  { name: 'Sabtu', sales: 34100000, orders: 98 },
  { name: 'Minggu', sales: 28400000, orders: 74 },
];

const topProducts = [
  { name: 'Beras Premium Rajawali 10kg', category: 'Beras', sold: 420, stock: 85 },
  { name: 'Minyak Goreng SunCo 2L', category: 'Minyak Goreng', sold: 380, stock: 12 },
  { name: 'Gula Pasir Gulaku 1kg', category: 'Gula', sold: 290, stock: 45 },
  { name: 'Telur Ayam Negeri 1kg', category: 'Telur', sold: 210, stock: 8 },
];

export const OverviewPage: React.FC = () => {
  const [totalOmzet, setTotalOmzet] = React.useState(0);
  const [totalOrders, setTotalOrders] = React.useState(0);
  const [lowStockCount, setLowStockCount] = React.useState(0);
  const [totalProducts, setTotalProducts] = React.useState(0);
  const [topProducts, setTopProducts] = React.useState<any[]>([]);

  React.useEffect(() => {
    const fetchStats = async () => {
      try {
        const prodRes = await productService.getProducts(1, 50);
        const prods = Array.isArray(prodRes.data)
          ? prodRes.data
          : (prodRes.data?.items || []);
        setTotalProducts(prods.length);

        const lowStock = prods.filter((p: any) => {
          const units = (p.units && p.units.length > 0) ? p.units : (p.product_units && p.product_units.length > 0) ? p.product_units : [];
          return units.some((u: any) => u.stock_quantity < 10);
        });
        setLowStockCount(lowStock.length);
        setTopProducts(prods.slice(0, 4));
      } catch (e) {}

      try {
        const orderRes = await orderService.getOrders(1, 50);
        const orders = Array.isArray(orderRes.data)
          ? orderRes.data
          : (orderRes.data?.items || []);
        setTotalOrders(orders.length);

        const omzet = orders.reduce((sum: number, o: any) => sum + (o.total_amount || 0), 0);
        setTotalOmzet(omzet);
      } catch (e) {}
    };

    fetchStats();
  }, []);

  const salesData = [
    { name: 'Senin', sales: 0, orders: 0 },
    { name: 'Selasa', sales: 0, orders: 0 },
    { name: 'Rabu', sales: 0, orders: 0 },
    { name: 'Kamis', sales: 0, orders: 0 },
    { name: 'Jumat', sales: 0, orders: 0 },
    { name: 'Sabtu', sales: 0, orders: 0 },
    { name: 'Minggu', sales: 0, orders: 0 },
  ];

  return (
    <div className="space-y-6">
      <PageHeader
        title="Dashboard Overview"
        subtitle="Ringkasan transaksi penjualan, stok gudang, dan performa toko LuckyMart."
      />

      {/* Summary KPI Cards Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 p-5 shadow-sm">
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold uppercase text-slate-500">Omzet Hari Ini</span>
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-emerald-100 text-emerald-600 dark:bg-emerald-950/60 dark:text-emerald-400">
              <DollarSign className="h-5 w-5" />
            </div>
          </div>
          <div className="mt-3 text-2xl font-black text-slate-900 dark:text-white">
            {formatCurrency(totalOmzet)}
          </div>
          <div className="mt-2 flex items-center text-xs font-semibold text-slate-500">
            <span>Total akumulasi dari transaksi real</span>
          </div>
        </div>

        <div className="rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 p-5 shadow-sm">
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold uppercase text-slate-500">Pesanan Hari Ini</span>
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-blue-100 text-blue-600 dark:bg-blue-950/60 dark:text-blue-400">
              <ShoppingCart className="h-5 w-5" />
            </div>
          </div>
          <div className="mt-3 text-2xl font-black text-slate-900 dark:text-white">{totalOrders} Order</div>
          <div className="mt-2 flex items-center text-xs font-semibold text-slate-500">
            <Clock className="h-3.5 w-3.5 mr-1 text-amber-500" />
            <span>Total pesanan masuk</span>
          </div>
        </div>

        <div className="rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 p-5 shadow-sm">
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold uppercase text-slate-500">Peringatan Stok Rendah</span>
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-amber-100 text-amber-600 dark:bg-amber-950/60 dark:text-amber-400">
              <AlertTriangle className="h-5 w-5" />
            </div>
          </div>
          <div className="mt-3 text-2xl font-black text-amber-600 dark:text-amber-400">{lowStockCount} Produk</div>
          <div className="mt-2 text-xs text-slate-500">Stok di bawah batas minimum</div>
        </div>

        <div className="rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 p-5 shadow-sm">
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold uppercase text-slate-500">Total Katalog Produk</span>
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-purple-100 text-purple-600 dark:bg-purple-950/60 dark:text-purple-400">
              <Users className="h-5 w-5" />
            </div>
          </div>
          <div className="mt-3 text-2xl font-black text-slate-900 dark:text-white">{totalProducts}</div>
          <div className="mt-2 text-xs font-semibold text-purple-600">Produk siap di-order</div>
        </div>
      </div>

      {/* Sales Trend Chart */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 p-6 shadow-sm">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h2 className="text-base font-bold text-slate-900 dark:text-white">Grafik Penjualan Mingguan</h2>
              <p className="text-xs text-slate-500">Tren nominal transaksi penjualan (Rp)</p>
            </div>
          </div>
          <div className="h-72 w-full">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={salesData}>
                <defs>
                  <linearGradient id="colorSales" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#16a34a" stopOpacity={0.4} />
                    <stop offset="95%" stopColor="#16a34a" stopOpacity={0} />
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#334155" opacity={0.2} />
                <XAxis dataKey="name" stroke="#64748b" fontSize={12} />
                <YAxis stroke="#64748b" fontSize={12} tickFormatter={(val) => `${val / 1000000}JT`} />
                <Tooltip
                  formatter={(val: number) => [formatCurrency(val), 'Penjualan']}
                  contentStyle={{ backgroundColor: '#0f172a', borderColor: '#334155', borderRadius: '12px', color: '#fff' }}
                />
                <Area type="monotone" dataKey="sales" stroke="#16a34a" strokeWidth={3} fillOpacity={1} fill="url(#colorSales)" />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Top Products Box */}
        <div className="rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 p-6 shadow-sm">
          <h2 className="text-base font-bold text-slate-900 dark:text-white mb-4">Produk Terlaris</h2>
          {topProducts.length === 0 ? (
            <div className="text-center py-10 text-xs text-slate-400">Belum ada data produk di katalog.</div>
          ) : (
            <div className="space-y-4">
              {topProducts.map((p, idx) => {
                const units = (p.units && p.units.length > 0) ? p.units : (p.product_units && p.product_units.length > 0) ? p.product_units : [];
                const stock = units[0]?.stock_quantity || 0;
                return (
                  <div key={idx} className="flex items-center justify-between border-b border-slate-100 dark:border-slate-800 pb-3 last:border-0 last:pb-0">
                    <div>
                      <div className="text-sm font-semibold text-slate-900 dark:text-white truncate max-w-[180px]">{p.name}</div>
                      <div className="text-xs text-slate-500">{p.sku}</div>
                    </div>
                    <div className="text-right">
                      <div className="text-sm font-bold text-brand-600">Stok: {stock}</div>
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
