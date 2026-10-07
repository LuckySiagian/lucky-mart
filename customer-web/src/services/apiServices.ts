import { api } from '../lib/axios';

export const authService = {
  login: async (credentials: any) => {
    const res = await api.post('/auth/login', credentials);
    return res.data;
  },
  register: async (data: any) => {
    const res = await api.post('/auth/register', data);
    return res.data;
  },
  getProfile: async () => {
    const res = await api.get('/users/me');
    return res.data;
  },
};

export const productService = {
  getProducts: async (page = 1, limit = 12, category_id?: string) => {
    const res = await api.get('/products', { params: { page, limit, category_id } });
    return res.data;
  },
  getProductById: async (id: string) => {
    const res = await api.get(`/products/${id}`);
    return res.data;
  },
  getCategories: async () => {
    const res = await api.get('/categories');
    return res.data;
  },
};

export const cartService = {
  getCart: async () => {
    const res = await api.get('/cart');
    return res.data;
  },
  addToCart: async (data: { product_id: string; product_unit_id: string; quantity: number; notes?: string }) => {
    const res = await api.post('/cart/items', data);
    return res.data;
  },
  updateCartItem: async (id: string, data: { quantity: number; notes?: string }) => {
    const res = await api.put(`/cart/items/${id}`, data);
    return res.data;
  },
  removeCartItem: async (id: string) => {
    const res = await api.delete(`/cart/items/${id}`);
    return res.data;
  },
  clearCart: async () => {
    const res = await api.delete('/cart');
    return res.data;
  },
  getCheckoutSummary: async () => {
    const res = await api.post('/cart/checkout');
    return res.data;
  },
};

export const orderService = {
  createOrder: async (data: any) => {
    const res = await api.post('/orders', data);
    return res.data;
  },
  getOrders: async (page = 1, limit = 10) => {
    const res = await api.get('/orders', { params: { page, limit } });
    return res.data;
  },
  getOrderById: async (id: string) => {
    const res = await api.get(`/orders/${id}`);
    return res.data;
  },
  cancelOrder: async (id: string, reason: string) => {
    const res = await api.put(`/orders/${id}/cancel`, { reason });
    return res.data;
  },
};

export const addressService = {
  getAddresses: async () => {
    const res = await api.get('/addresses');
    return res.data;
  },
  createAddress: async (data: any) => {
    const res = await api.post('/addresses', data);
    return res.data;
  },
  setDefaultAddress: async (id: string) => {
    const res = await api.patch(`/addresses/${id}/default`);
    return res.data;
  },
};

export const shippingService = {
  getShippingMethods: async () => {
    const res = await api.get('/shipping/methods');
    return res.data;
  },
  calculateRate: async (data: any) => {
    const res = await api.post('/shipping/calculate', data);
    return res.data;
  },
  getShipmentByOrderId: async (orderId: string) => {
    const res = await api.get(`/shipping/orders/${orderId}`);
    return res.data;
  },
};

export const paymentService = {
  getPaymentMethods: async () => {
    const res = await api.get('/payment-methods');
    return res.data;
  },
  createPayment: async (data: any) => {
    const res = await api.post('/payments', data);
    return res.data;
  },
  uploadProof: async (data: { payment_id: string; proof_image_url: string; notes?: string }) => {
    const res = await api.post('/payments/upload-proof', data);
    return res.data;
  },
};

export const notificationService = {
  getNotifications: async () => {
    const res = await api.get('/notifications');
    return res.data;
  },
  markAllRead: async () => {
    const res = await api.patch('/notifications/read-all');
    return res.data;
  },
};
