import { createBrowserRouter, Navigate } from 'react-router-dom';
import { MainLayout } from '../layouts/MainLayout';
import { ProtectedRoute, homeFor } from '../layouts/ProtectedRoute';
import { useAuthStore } from '../store/useAuthStore';
import { LoginPage } from '../features/auth/LoginPage';
import { OverviewPage } from '../features/dashboard/OverviewPage';
import { ProductListPage } from '../features/products/ProductListPage';
import { CategoryListPage } from '../features/products/CategoryListPage';
import { InventoryPage } from '../features/inventory/InventoryPage';
import { StockMovementPage } from '../features/inventory/StockMovementPage';
import { OrderListPage } from '../features/orders/OrderListPage';
import { PaymentListPage } from '../features/payments/PaymentListPage';
import { CreditLedgerPage } from '../features/payments/CreditLedgerPage';
import { ShipmentListPage } from '../features/shipping/ShipmentListPage';
import { NotificationPage } from '../features/notifications/NotificationPage';
import { UserListPage } from '../features/users/UserListPage';
import { AuditLogPage } from '../features/audit/AuditLogPage';
import { SettingsPage } from '../features/settings/SettingsPage';
import type { Role } from '../types';

const WAREHOUSE: Role[] = ['SUPER_ADMIN', 'STORE_MANAGER', 'WAREHOUSE_STAFF'];
const MANAGERS: Role[] = ['SUPER_ADMIN', 'STORE_MANAGER'];

const Home = () => <Navigate to={homeFor(useAuthStore.getState().user?.role)} replace />;

export const router = createBrowserRouter([
  { path: '/login', element: <LoginPage /> },
  {
    element: <ProtectedRoute />,
    children: [
      {
        element: <MainLayout />,
        children: [
          { path: '/', element: <Home /> },
          // All staff, including couriers
          { path: '/orders', element: <OrderListPage /> },
          { path: '/shipping', element: <ShipmentListPage /> },
          { path: '/notifications', element: <NotificationPage /> },
          { path: '/settings', element: <SettingsPage /> },
          {
            element: <ProtectedRoute allowedRoles={WAREHOUSE} />,
            children: [
              { path: '/dashboard', element: <OverviewPage /> },
              { path: '/products', element: <ProductListPage /> },
              { path: '/categories', element: <CategoryListPage /> },
              { path: '/inventory', element: <InventoryPage /> },
              { path: '/stock-movements', element: <StockMovementPage /> },
            ],
          },
          {
            element: <ProtectedRoute allowedRoles={MANAGERS} />,
            children: [
              { path: '/payments', element: <PaymentListPage /> },
              { path: '/credit-ledgers', element: <CreditLedgerPage /> },
              { path: '/users', element: <UserListPage /> },
              { path: '/events', element: <AuditLogPage /> },
            ],
          },
        ],
      },
    ],
  },
  { path: '*', element: <Home /> },
]);
