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

export interface User {
  id: string;
  full_name: string;
  email?: string;
  phone_number: string;
  role: Role;
  b2b_profile?: {
    store_name: string;
    store_address: string;
    verification_status: 'PENDING' | 'APPROVED' | 'REJECTED';
    approved_credit_limit: number;
  };
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  user: User;
}

export interface Category {
  id: string;
  name: string;
  slug: string;
  description?: string;
}

export interface ProductUnit {
  id: string;
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
}

export interface CartItem {
  id: string;
  product_id: string;
  product_name: string;
  image_url?: string;
  sku: string;
  product_unit_id: string;
  unit_name: string;
  quantity: number;
  unit_price: number;
  subtotal: number;
  notes?: string;
  available_stock: number;
  is_stock_sufficient: boolean;
  is_product_active: boolean;
}

export interface Cart {
  id: string;
  total_items: number;
  total_amount: number;
  items: CartItem[] | null;
}

export interface CheckoutSummary {
  total_items: number;
  subtotal_amount: number;
  shipping_method_id?: string;
  estimated_shipping_fee: number;
  free_shipping_min: number;
  grand_total_amount: number;
  items: CartItem[];
  can_proceed_to_order: boolean;
  validation_errors?: string[];
}

export interface Address {
  id: string;
  title: string;
  recipient_name: string;
  phone: string;
  full_address: string;
  city: string;
  province: string;
  postal_code: string;
  is_default: boolean;
}

export interface PaymentMethod {
  id: string;
  code: string;
  name: string;
  type: 'COD' | 'MANUAL_TRANSFER' | 'B2B_CREDIT' | 'CASH' | 'GATEWAY';
  description?: string;
}

export interface ShippingMethod {
  id: string;
  name: string;
  description?: string;
  estimated_days: string;
  base_rate: number;
  courier_name: string;
}

export interface Payment {
  id: string;
  payment_number: string;
  method_code: string;
  method_name: string;
  method_type: string;
  grand_total: number;
  status: string;
  proof_image_url?: string;
  paid_at?: string;
  expired_at: string;
  notes?: string;
}

export interface Order {
  id: string;
  order_number: string;
  status: OrderStatus;
  can_cancel: boolean;
  payment_method: string;
  delivery_type: 'DELIVERY' | 'PICKUP';
  total_amount: number;
  shipping_fee: number;
  grand_total: number;
  total_items: number;
  due_date?: string;
  order_date: string;
  items: {
    id: string;
    product_id: string;
    product_name: string;
    unit_name: string;
    quantity: number;
    unit_price: number;
    subtotal: number;
  }[];
  address?: { recipient_name: string; phone: string; full_address: string; city: string; postal_code: string; notes?: string };
  invoice?: { invoice_number: string; status: string; due_date: string };
  payments?: Payment[];
  histories: { id: string; status: OrderStatus; notes?: string; created_at: string }[];
}

export interface Shipment {
  tracking_number: string;
  courier_name: string;
  status: string;
  trackings: { id: string; status: string; location: string; notes?: string; created_at: string }[];
}

export interface CreditLedger {
  credit_limit: number;
  used_credit: number;
  remaining_credit: number;
  outstanding_balance: number;
}

export interface Notification {
  id: string;
  title: string;
  message: string;
  is_read: boolean;
  created_at: string;
}

export interface StoreInfo {
  name: string;
  tagline: string;
  address: string;
  whatsapp: string;
  open_hours: string;
  maps_url: string;
  receipt_note: string;
  announcement: string;
  free_shipping_min: number;
}
