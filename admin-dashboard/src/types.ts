// Shapes returned by the Lucky Mart API (see internal/dto in the backend).

export type Role = 'SUPER_ADMIN' | 'STORE_MANAGER' | 'WAREHOUSE_STAFF' | 'DRIVER' | 'B2B_CUSTOMER' | 'B2C_CUSTOMER';

export type OrderStatus =
  | 'CREATED'
  | 'PENDING_PAYMENT'
  | 'PAID'
  | 'CONFIRMED'
  | 'PROCESSING_PACKING'
  | 'OUT_FOR_DELIVERY'
  | 'COMPLETED'
  | 'CANCELLED';

export interface Paginated<T> {
  items: T[];
  meta: { page: number; limit: number; total_records: number; total_pages: number };
}

export interface B2BProfile {
  id: string;
  store_name: string;
  store_address: string;
  tax_id_npwp?: string;
  verification_status: 'PENDING' | 'APPROVED' | 'REJECTED';
  approved_credit_limit: number;
}

export interface User {
  id: string;
  full_name: string;
  email?: string;
  phone_number: string;
  role: Role;
  is_active: boolean;
  last_login_at?: string;
  created_at: string;
  b2b_profile?: B2BProfile;
}

export interface Category {
  id: string;
  name: string;
  slug: string;
  description?: string;
}

export interface ProductUnit {
  id: string;
  unit_id: string;
  unit_name: string;
  conversion_factor: number;
  retail_price: number;
  wholesale_price: number;
  min_wholesale_qty: number;
  stock_quantity: number;
}

export interface Product {
  id: string;
  sku: string;
  name: string;
  description?: string;
  is_active: boolean;
  category_id: string;
  category: Category;
  units: ProductUnit[];
  images: { id: string; image_url: string; is_primary: boolean }[];
  image_url?: string;
  total_stock: number;
  created_at: string;
}

export interface InventoryItem {
  id: string;
  product_id: string;
  product_unit_id: string;
  warehouse_name: string;
  product_name: string;
  sku: string;
  category_name?: string;
  unit_name: string;
  retail_price: number;
  current_stock: number;
  min_stock_alert: number;
  is_low_stock: boolean;
  is_product_active: boolean;
  updated_at: string;
}

export interface StockMovement {
  id: string;
  inventory_id: string;
  product_name?: string;
  sku?: string;
  unit_name?: string;
  movement_type: 'INITIAL_STOCK' | 'STOCK_IN' | 'STOCK_OUT' | 'ADJUSTMENT';
  reference_number: string;
  quantity_before: number;
  quantity_changed: number;
  quantity_after: number;
  notes?: string;
  created_at: string;
}

export interface PaymentLog {
  id: string;
  action: string;
  old_status: string;
  new_status: string;
  notes?: string;
  created_at: string;
}

export interface Payment {
  id: string;
  payment_number: string;
  order_id: string;
  order_number?: string;
  order_status?: OrderStatus;
  invoice_id: string;
  user_id: string;
  customer_name?: string;
  customer_phone?: string;
  method_code: string;
  method_name: string;
  method_type: string;
  amount: number;
  grand_total: number;
  status: string;
  proof_image_url?: string;
  paid_at?: string;
  expired_at: string;
  notes?: string;
  logs: PaymentLog[];
  created_at: string;
}

export interface OrderItem {
  id: string;
  product_id: string;
  product_name: string;
  sku: string;
  unit_name: string;
  quantity: number;
  unit_price: number;
  subtotal: number;
  notes?: string;
}

export interface Order {
  id: string;
  order_number: string;
  user_id: string;
  customer?: { id: string; full_name: string; phone_number: string; role: Role };
  status: OrderStatus;
  next_statuses: OrderStatus[];
  can_cancel: boolean;
  payment_method: string;
  delivery_type: 'DELIVERY' | 'PICKUP';
  total_amount: number;
  shipping_fee: number;
  grand_total: number;
  total_items: number;
  due_date?: string;
  order_date: string;
  items: OrderItem[];
  address?: { recipient_name: string; phone: string; full_address: string; city: string; postal_code: string; notes?: string };
  invoice?: { id: string; invoice_number: string; status: string; total_amount: number; due_date: string; paid_at?: string };
  payments?: Payment[];
  histories: { id: string; status: OrderStatus; notes?: string; created_at: string }[];
}

export interface Shipment {
  id: string;
  tracking_number: string;
  order_id: string;
  order_number?: string;
  order_status?: OrderStatus;
  customer_name?: string;
  recipient_name?: string;
  recipient_phone?: string;
  address?: string;
  courier_name: string;
  method_name: string;
  status: string;
  shipping_fee: number;
  shipped_at?: string;
  delivered_at?: string;
  trackings: { id: string; status: string; location: string; notes?: string; created_at: string }[];
  created_at: string;
}

export interface CreditLedger {
  id: string;
  user_id: string;
  user_name: string;
  user_phone?: string;
  store_name?: string;
  credit_limit: number;
  used_credit: number;
  remaining_credit: number;
  outstanding_balance: number;
  updated_at: string;
}

export interface CreditHistory {
  id: string;
  order_id?: string;
  type: 'CREDIT_USAGE' | 'PAYMENT_SETTLEMENT' | 'REVERSAL';
  amount: number;
  payment_date: string;
  remaining_balance: number;
  notes?: string;
}

export interface Notification {
  id: string;
  title: string;
  message: string;
  type: string;
  is_read: boolean;
  created_at: string;
}

export interface EventLog {
  id: string;
  event_name: string;
  aggregate: string;
  aggregate_id: string;
  payload: string;
  status: string;
  created_at: string;
}

export interface PaymentMethod {
  id: string;
  code: string;
  name: string;
  type: string;
  is_active: boolean;
  description?: string;
}

export interface ShippingMethod {
  id: string;
  code: string;
  name: string;
  description?: string;
  estimated_days: string;
  base_rate: number;
  courier_name: string;
}

export interface DashboardStats {
  today_revenue: number;
  today_orders: number;
  month_revenue: number;
  month_orders: number;
  active_products: number;
  total_customers: number;
  kasbon_outstanding: number;
  pending: {
    orders_to_confirm: number;
    orders_to_pack: number;
    orders_packing: number;
    orders_delivering: number;
    awaiting_payment: number;
    payments_to_verify: number;
    b2b_to_verify: number;
    low_stock_items: number;
  };
  sales_last_7_days: { date: string; revenue: number; orders: number }[];
  top_products: { product_name: string; quantity: number; revenue: number }[];
  recent_orders: { id: string; order_number: string; customer_name: string; grand_total: number; status: OrderStatus; order_date: string }[];
}
