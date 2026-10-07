import { createBrowserRouter, Navigate } from 'react-router-dom';
import { MainLayout } from '../layouts/MainLayout';
import { ProtectedRoute } from '../layouts/ProtectedRoute';
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

export const router = createBrowserRouter([
  {
    path: '/login',
    element: <LoginPage />,
  },
  {
    element: <ProtectedRoute />,
    children: [
      {
        element: <MainLayout />,
        children: [
          { path: '/', element: <Navigate to="/dashboard" replace /> },
          { path: '/dashboard', element: <OverviewPage /> },
          { path: '/products', element: <ProductListPage /> },
          { path: '/categories', element: <CategoryListPage /> },
          { path: '/inventory', element: <InventoryPage /> },
          { path: '/stock-movements', element: <StockMovementPage /> },
          { path: '/orders', element: <OrderListPage /> },
          { path: '/payments', element: <PaymentListPage /> },
          { path: '/credit-ledgers', element: <CreditLedgerPage /> },
          { path: '/shipping', element: <ShipmentListPage /> },
          { path: '/notifications', element: <NotificationPage /> },
          { path: '/events', element: <AuditLogPage /> },
          { path: '/users', element: <UserListPage /> },
          { path: '/settings', element: <SettingsPage /> },
        ],
      },
    ],
  },
  {
    path: '*',
    element: <Navigate to="/dashboard" replace />,
  },
]);
