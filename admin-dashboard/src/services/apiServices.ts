import type { AxiosResponse } from 'axios';
import { api } from '../lib/axios';
import type {
  Category,
  CreditHistory,
  CreditLedger,
  DashboardStats,
  EventLog,
  InventoryItem,
  Notification,
  Order,
  OrderStatus,
  Paginated,
  Payment,
  PaymentMethod,
  Product,
  Shipment,
  ShippingMethod,
  StockMovement,
  User,
  StoreInfo,
} from '../types';

// Every API response is { success, message, data }; these helpers return `data`.
const data = <T>(p: Promise<AxiosResponse<{ data: T }>>) => p.then((r) => r.data.data);

export interface ListParams {
  page?: number;
  limit?: number;
  search?: string;
  status?: string;
  [key: string]: string | number | undefined;
}

const clean = (params?: ListParams) =>
  Object.fromEntries(Object.entries(params || {}).filter(([, v]) => v !== undefined && v !== ''));

export const authService = {
  login: (credentials: { phone_number: string; password: string }) =>
    data<{ access_token: string; refresh_token: string; user: User }>(api.post('/auth/login', credentials)),
  me: () => data<User>(api.get('/auth/me')),
  updateProfile: (body: { full_name?: string; email?: string }) => data<User>(api.put('/auth/me', body)),
  changePassword: (body: { old_password: string; new_password: string }) => data(api.post('/auth/change-password', body)),
  logout: () => data(api.post('/auth/logout')),
};

export const dashboardService = {
  getStats: () => data<DashboardStats>(api.get('/admin/dashboard')),
};

export interface ProductUnitInput {
  unit_name: string;
  retail_price: number;
  wholesale_price?: number;
  min_wholesale_qty?: number;
  initial_stock: number;
}

export const productService = {
  getProducts: (params?: ListParams) => data<Paginated<Product>>(api.get('/admin/products', { params: clean(params) })),
  createProduct: (body: {
    name: string;
    sku?: string;
    category_id: string;
    description?: string;
    image_url?: string;
    min_stock_alert?: number;
    units: ProductUnitInput[];
  }) => data<Product>(api.post('/products', body)),
  updateProduct: (id: string, body: Record<string, unknown>) => data<Product>(api.put(`/products/${id}`, body)),
  deleteProduct: (id: string) => data(api.delete(`/products/${id}`)),
  getCategories: () => data<Category[]>(api.get('/categories')),
  createCategory: (body: { name: string; description?: string }) => data<Category>(api.post('/categories', body)),
  updateCategory: (id: string, body: { name: string; description?: string }) => data<Category>(api.put(`/categories/${id}`, body)),
  deleteCategory: (id: string) => data(api.delete(`/categories/${id}`)),
};

export const inventoryService = {
  getInventory: (params?: ListParams) => data<Paginated<InventoryItem>>(api.get('/inventory', { params: clean(params) })),
  getLowStock: () => data<InventoryItem[]>(api.get('/inventory/low-stock')),
  getMovements: (params?: ListParams) => data<Paginated<StockMovement>>(api.get('/inventory/movements', { params: clean(params) })),
  adjustStock: (body: {
    product_unit_id: string;
    adjustment_type: 'STOCK_IN' | 'STOCK_OUT' | 'ADJUSTMENT';
    quantity_changed: number;
    reference_number?: string;
    notes: string;
  }) => data<StockMovement>(api.post('/inventory/adjust', body)),
  updateMinStock: (id: string, min_stock_alert: number) => data<InventoryItem>(api.patch(`/inventory/${id}/min-stock`, { min_stock_alert })),
};

export const orderService = {
  getOrders: (params?: ListParams) => data<Paginated<Order>>(api.get('/orders', { params: clean(params) })),
  getOrderById: (id: string) => data<Order>(api.get(`/orders/${id}`)),
  updateOrderStatus: (id: string, status: OrderStatus, notes?: string) => data<Order>(api.patch(`/orders/${id}/status`, { status, notes })),
  cancelOrder: (id: string, reason: string) => data<Order>(api.put(`/orders/${id}/cancel`, { reason })),
};

export const paymentService = {
  getPayments: (params?: ListParams) => data<Paginated<Payment>>(api.get('/payments', { params: clean(params) })),
  getPaymentById: (id: string) => data<Payment>(api.get(`/payments/${id}`)),
  approvePayment: (id: string, notes?: string) => data<Payment>(api.patch(`/payments/${id}/approve`, { notes })),
  rejectPayment: (id: string, notes: string) => data<Payment>(api.patch(`/payments/${id}/reject`, { notes })),
  refundPayment: (id: string, refund_amount: number, reason: string) =>
    data<Payment>(api.patch(`/payments/${id}/refund`, { refund_amount, reason })),
  getMethods: () => data<PaymentMethod[]>(api.get('/payment-methods/all')),
  updateMethod: (id: string, body: { name?: string; description?: string; is_active?: boolean }) =>
    data<PaymentMethod>(api.put(`/payment-methods/${id}`, body)),
};

export const creditService = {
  getLedgers: (params?: ListParams) => data<Paginated<CreditLedger>>(api.get('/credit-ledgers', { params: clean(params) })),
  createLedger: (user_id: string, credit_limit: number) => data<CreditLedger>(api.post('/credit-ledgers', { user_id, credit_limit })),
  adjustLimit: (id: string, new_credit_limit: number) => data<CreditLedger>(api.patch(`/credit-ledgers/${id}`, { new_credit_limit })),
  settle: (id: string, amount: number, notes?: string) => data<CreditLedger>(api.post(`/credit-ledgers/${id}/settle`, { amount, notes })),
  getHistory: (id: string, params?: ListParams) => data<Paginated<CreditHistory>>(api.get(`/credit-ledgers/${id}/history`, { params: clean(params) })),
};

export const shippingService = {
  getMethods: () => data<ShippingMethod[]>(api.get('/shipping/methods')),
  getShipments: (params?: ListParams) => data<Paginated<Shipment>>(api.get('/shipping/shipments', { params: clean(params) })),
  createShipment: (order_id: string, notes?: string) => data<Shipment>(api.post('/shipping/shipments', { order_id, notes })),
  updateStatus: (id: string, body: { status: string; location?: string; notes?: string }) =>
    data<Shipment>(api.patch(`/shipping/shipments/${id}/status`, body)),
};

export const notificationService = {
  getNotifications: (params?: ListParams) => data<Paginated<Notification>>(api.get('/notifications', { params: clean(params) })),
  unreadCount: () => data<{ unread: number }>(api.get('/notifications/unread-count')),
  markRead: (id: string) => data(api.patch(`/notifications/${id}/read`)),
  markAllRead: () => data(api.patch('/notifications/read-all')),
  remove: (id: string) => data(api.delete(`/notifications/${id}`)),
};

export const eventService = {
  getEvents: (params?: ListParams) => data<Paginated<EventLog>>(api.get('/events', { params: clean(params) })),
};

export const userService = {
  getUsers: (params?: ListParams) => data<Paginated<User>>(api.get('/users', { params: clean(params) })),
  createUser: (body: { full_name: string; phone_number: string; email?: string; password: string; role: string }) =>
    data<User>(api.post('/users', body)),
  updateUser: (id: string, body: { full_name?: string; email?: string; role?: string; is_active?: boolean }) =>
    data<User>(api.put(`/users/${id}`, body)),
  resetPassword: (id: string, new_password: string) => data(api.post(`/users/${id}/reset-password`, { new_password })),
  verifyB2B: (id: string, status: 'APPROVED' | 'REJECTED' | 'PENDING', credit_limit?: number) =>
    data<User>(api.patch(`/users/${id}/b2b`, { status, credit_limit })),
  deleteUser: (id: string) => data(api.delete(`/users/${id}`)),
};

export const storeService = {
  get: () => data<StoreInfo>(api.get('/store')),
  update: (body: StoreInfo) => data<StoreInfo>(api.put('/store', body)),
};
