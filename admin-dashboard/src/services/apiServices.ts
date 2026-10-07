import { api } from '../lib/axios';

export const authService = {
  login: async (credentials: any) => {
    const res = await api.post('/auth/login', credentials);
    return res.data;
  },
  getProfile: async () => {
    const res = await api.get('/users/me');
    return res.data;
  },
  logout: async () => {
    const res = await api.post('/auth/logout');
    return res.data;
  },
};

export const productService = {
  getProducts: async (page = 1, limit = 10, category_id?: string) => {
    const res = await api.get('/products', { params: { page, limit, category_id } });
    return res.data;
  },
  getProductById: async (id: string) => {
    const res = await api.get(`/products/${id}`);
    return res.data;
  },
  createProduct: async (data: any) => {
    const res = await api.post('/products', data);
    return res.data;
  },
  updateProduct: async (id: string, data: any) => {
    const res = await api.put(`/products/${id}`, data);
    return res.data;
  },
  deleteProduct: async (id: string) => {
    const res = await api.delete(`/products/${id}`);
    return res.data;
  },
  getCategories: async () => {
    const res = await api.get('/categories');
    return res.data;
  },
  createCategory: async (data: { name: string; description?: string }) => {
    const res = await api.post('/categories', data);
    return res.data;
  },
  updateCategory: async (id: string, data: { name: string; description?: string }) => {
    const res = await api.put(`/categories/${id}`, data);
    return res.data;
  },
  deleteCategory: async (id: string) => {
    const res = await api.delete(`/categories/${id}`);
    return res.data;
  },
};

export const inventoryService = {
  getLowStock: async () => {
    const res = await api.get('/inventory/low-stock');
    return res.data;
  },
  adjustStock: async (data: any) => {
    const res = await api.post('/inventory/adjust', data);
    return res.data;
  },
};

export const orderService = {
  getOrders: async (page = 1, limit = 50) => {
    const res = await api.get('/orders', { params: { page, limit } });
    return res.data;
  },
  getOrderById: async (id: string) => {
    const res = await api.get(`/orders/${id}`);
    return res.data;
  },
  updateOrderStatus: async (id: string, status: string, notes?: string) => {
    const res = await api.put(`/orders/${id}/status`, { status, notes });
    return res.data;
  },
  cancelOrder: async (id: string, reason: string) => {
    const res = await api.put(`/orders/${id}/cancel`, { reason });
    return res.data;
  },
};

export const paymentService = {
  getPayments: async (page = 1, limit = 50) => {
    const res = await api.get('/payments', { params: { page, limit } });
    return res.data;
  },
  getWaitingConfirmation: async (page = 1, limit = 50) => {
    const res = await api.get('/payments/waiting-confirmation', { params: { page, limit } });
    return res.data;
  },
  approvePayment: async (id: string, notes?: string) => {
    const res = await api.patch(`/payments/${id}/approve`, { notes });
    return res.data;
  },
  rejectPayment: async (id: string, notes?: string) => {
    const res = await api.patch(`/payments/${id}/reject`, { notes });
    return res.data;
  },
  getCreditLedgers: async (page = 1, limit = 50) => {
    const res = await api.get('/credit-ledgers', { params: { page, limit } });
    return res.data;
  },
  adjustCreditLimit: async (id: string, new_credit_limit: number) => {
    const res = await api.patch(`/credit-ledgers/${id}`, { new_credit_limit });
    return res.data;
  },
};

export const shippingService = {
  getShippingMethods: async () => {
    const res = await api.get('/shipping/methods');
    return res.data;
  },
  getShipmentByOrderId: async (orderId: string) => {
    const res = await api.get(`/shipping/orders/${orderId}`);
    return res.data;
  },
  updateShipmentStatus: async (shipmentId: string, data: any) => {
    const res = await api.patch(`/shipping/shipments/${shipmentId}/status`, data);
    return res.data;
  },
};

export const notificationService = {
  getNotifications: async (page = 1, limit = 50) => {
    const res = await api.get('/notifications', { params: { page, limit } });
    return res.data;
  },
  markAllRead: async () => {
    const res = await api.patch('/notifications/read-all');
    return res.data;
  },
  getEventLogs: async (page = 1, limit = 50) => {
    const res = await api.get('/events', { params: { page, limit } });
    return res.data;
  },
};

export const userService = {
  getUsers: async (page = 1, limit = 50) => {
    const res = await api.get('/users', { params: { page, limit } });
    return res.data;
  },
  updateUserRole: async (id: string, role: string) => {
    const res = await api.put(`/users/${id}/role`, { role });
    return res.data;
  },
};

