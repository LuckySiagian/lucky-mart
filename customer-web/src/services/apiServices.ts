import type { AxiosResponse } from 'axios';
import { api } from '../lib/axios';
import type {
  Address,
  AuthResponse,
  Cart,
  Category,
  CheckoutSummary,
  CreditLedger,
  Notification,
  Order,
  Paginated,
  Payment,
  PaymentMethod,
  Product,
  Shipment,
  ShippingMethod,
  User,
  StoreInfo,
} from '../types';

// Every API response is { success, message, data }; these helpers return `data`.
const data = <T>(p: Promise<AxiosResponse<{ data: T }>>) => p.then((r) => r.data.data);

const clean = (params: Record<string, unknown>) =>
  Object.fromEntries(Object.entries(params).filter(([, v]) => v !== undefined && v !== ''));

export interface RegisterInput {
  full_name: string;
  phone_number: string;
  email?: string;
  password: string;
  is_b2b_registration: boolean;
  store_name?: string;
  store_address?: string;
}

export const authService = {
  login: (body: { phone_number: string; password: string }) => data<AuthResponse>(api.post('/auth/login', body)),
  register: (body: RegisterInput) => data<AuthResponse>(api.post('/auth/register', body)),
  me: () => data<User>(api.get('/auth/me')),
  updateProfile: (body: { full_name?: string; email?: string }) => data<User>(api.put('/auth/me', body)),
  changePassword: (body: { old_password: string; new_password: string }) => data(api.post('/auth/change-password', body)),
  forgotPassword: (phone_number: string) => api.post('/auth/forgot-password', { phone_number }).then((r) => r.data.message as string),
  logout: () => data(api.post('/auth/logout')),
};

export const productService = {
  getProducts: (params: { page?: number; limit?: number; category_id?: string; search?: string } = {}) =>
    data<Paginated<Product>>(api.get('/products', { params: clean(params) })),
  getProductById: (id: string) => data<Product>(api.get(`/products/${id}`)),
  getCategories: () => data<Category[]>(api.get('/categories')),
};

export const cartService = {
  getCart: () => data<Cart>(api.get('/cart')),
  addToCart: (body: { product_id: string; product_unit_id: string; quantity: number; notes?: string }) =>
    data<Cart>(api.post('/cart/items', body)),
  updateCartItem: (id: string, quantity: number) => data<Cart>(api.put(`/cart/items/${id}`, { quantity })),
  removeCartItem: (id: string) => data<Cart>(api.delete(`/cart/items/${id}`)),
  clearCart: () => data(api.delete('/cart')),
  checkoutSummary: (body: { cart_item_ids?: string[]; shipping_method_id?: string; city?: string; delivery_type?: DeliveryType }) =>
    data<CheckoutSummary>(api.post('/cart/checkout', body)),
};

export type DeliveryType = 'DELIVERY' | 'PICKUP';

export interface CreateOrderInput {
  payment_method: string;
  delivery_type?: DeliveryType;
  cart_item_ids?: string[];
  shipping_method_id?: string;
  address_id?: string;
  recipient_name?: string;
  phone?: string;
  full_address?: string;
  city?: string;
  postal_code?: string;
  notes?: string;
}

export const orderService = {
  createOrder: (body: CreateOrderInput) => data<Order>(api.post('/orders', body)),
  getOrders: (params: { page?: number; limit?: number; status?: string } = {}) =>
    data<Paginated<Order>>(api.get('/orders', { params: clean(params) })),
  getOrderById: (id: string) => data<Order>(api.get(`/orders/${id}`)),
  cancelOrder: (id: string, reason: string) => data<Order>(api.put(`/orders/${id}/cancel`, { reason })),
};

export const addressService = {
  getAddresses: () => data<Address[]>(api.get('/addresses')),
  createAddress: (body: Omit<Address, 'id'>) => data<Address>(api.post('/addresses', body)),
  updateAddress: (id: string, body: Omit<Address, 'id'>) => data<Address>(api.put(`/addresses/${id}`, body)),
  deleteAddress: (id: string) => data(api.delete(`/addresses/${id}`)),
  setDefault: (id: string) => data<Address>(api.patch(`/addresses/${id}/default`)),
};

export const shippingService = {
  getMethods: () => data<ShippingMethod[]>(api.get('/shipping/methods')),
  getShipmentByOrderId: (orderId: string) => data<Shipment>(api.get(`/shipping/orders/${orderId}`)),
};

export const paymentService = {
  getMethods: () => data<PaymentMethod[]>(api.get('/payment-methods')),
  uploadProof: (body: { payment_id: string; proof_image_url: string; notes?: string }) =>
    data<Payment>(api.post('/payments/upload-proof', body)),
};

export const creditService = {
  getMine: () => data<CreditLedger>(api.get('/credit-ledgers/me')),
};

export const notificationService = {
  getNotifications: (params: { page?: number; limit?: number } = {}) =>
    data<Paginated<Notification>>(api.get('/notifications', { params: clean(params) })),
  unreadCount: () => data<{ unread: number }>(api.get('/notifications/unread-count')),
  markRead: (id: string) => data(api.patch(`/notifications/${id}/read`)),
  markAllRead: () => data(api.patch('/notifications/read-all')),
};

export const storeService = {
  get: () => data<StoreInfo>(api.get('/store')),
};
