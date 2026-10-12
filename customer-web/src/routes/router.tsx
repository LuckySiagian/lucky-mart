import { createBrowserRouter, Navigate } from 'react-router-dom';
import { MainLayout } from '../layouts/MainLayout';
import { ProtectedRoute } from '../layouts/ProtectedRoute';
import { HomePage } from '../features/home/HomePage';
import { ProductListPage } from '../features/products/ProductListPage';
import { ProductDetailPage } from '../features/products/ProductDetailPage';
import { CartPage } from '../features/cart/CartPage';
import { CheckoutPage } from '../features/checkout/CheckoutPage';
import { OrderListPage } from '../features/orders/OrderListPage';
import { OrderDetailPage } from '../features/orders/OrderDetailPage';
import { ProfilePage } from '../features/profile/ProfilePage';
import { NotificationPage } from '../features/notifications/NotificationPage';
import { LoginPage } from '../features/auth/LoginPage';
import { RegisterPage } from '../features/auth/RegisterPage';

export const router = createBrowserRouter([
  { path: '/login', element: <LoginPage /> },
  { path: '/register', element: <RegisterPage /> },
  {
    element: <MainLayout />,
    children: [
      { path: '/', element: <HomePage /> },
      { path: '/products', element: <ProductListPage /> },
      { path: '/products/:id', element: <ProductDetailPage /> },
      { path: '/wishlist', element: <ProductListPage wishlistOnly /> },
      {
        element: <ProtectedRoute />,
        children: [
          { path: '/cart', element: <CartPage /> },
          { path: '/checkout', element: <CheckoutPage /> },
          { path: '/orders', element: <OrderListPage /> },
          { path: '/orders/:id', element: <OrderDetailPage /> },
          { path: '/profile', element: <ProfilePage /> },
          { path: '/notifications', element: <NotificationPage /> },
        ],
      },
      { path: '*', element: <Navigate to="/" replace /> },
    ],
  },
]);
