import { create } from 'zustand';
import { persist } from 'zustand/middleware';

export interface CartItem {
  id: string;
  product_id: string;
  product_name: string;
  product_unit_id: string;
  unit_name: string;
  quantity: number;
  unit_price: number;
  subtotal: number;
  notes?: string;
  available_stock: number;
  is_stock_sufficient: boolean;
  added_at?: number; // Epoch timestamp ms
}

export const CART_MAX_AGE_MS = 30 * 24 * 60 * 60 * 1000; // 30 Hari

const isExpired = (item: CartItem) => {
  if (!item.added_at) return false;
  return Date.now() - item.added_at > CART_MAX_AGE_MS;
};

const filterValidItems = (items: CartItem[] = []): CartItem[] => {
  return items.filter((item) => !isExpired(item));
};

interface CartState {
  items: CartItem[];
  totalItems: number;
  totalAmount: number;
  selectedItemIdsForCheckout: string[];
  setCart: (items: CartItem[], totalItems?: number, totalAmount?: number) => void;
  addItem: (item: CartItem) => void;
  updateQuantity: (itemId: string, newQuantity: number) => void;
  removeItem: (itemId: string) => void;
  removeMultipleItems: (itemIds: string[]) => void;
  clearCart: () => void;
  cleanExpiredItems: () => number;
  setSelectedForCheckout: (ids: string[]) => void;
}

export const useCartStore = create<CartState>()(
  persist(
    (set, get) => ({
      items: [],
      totalItems: 0,
      totalAmount: 0,
      selectedItemIdsForCheckout: [],
      setSelectedForCheckout: (ids) => set({ selectedItemIdsForCheckout: ids }),
      setCart: (rawItems, customTotalItems, customTotalAmount) => {
        const validItems = filterValidItems(rawItems.map(i => ({ ...i, added_at: i.added_at || Date.now() })));
        const totalItems = customTotalItems !== undefined ? customTotalItems : validItems.reduce((sum, i) => sum + i.quantity, 0);
        const totalAmount = customTotalAmount !== undefined ? customTotalAmount : validItems.reduce((sum, i) => sum + i.subtotal, 0);
        set({ items: validItems, totalItems, totalAmount });
      },
      addItem: (newItem) =>
        set((state) => {
          const now = Date.now();
          const itemWithDate = { ...newItem, added_at: newItem.added_at || now };
          const validExisting = filterValidItems(state.items);
          const existing = validExisting.find((i) => i.product_unit_id === itemWithDate.product_unit_id);
          
          let updatedItems: CartItem[] = [];
          if (existing) {
            updatedItems = validExisting.map((i) =>
              i.product_unit_id === itemWithDate.product_unit_id
                ? { 
                    ...i, 
                    quantity: i.quantity + itemWithDate.quantity, 
                    subtotal: (i.quantity + itemWithDate.quantity) * i.unit_price,
                    added_at: now // Perbarui timestamp saat jumlah ditambah
                  }
                : i
            );
          } else {
            updatedItems = [...validExisting, itemWithDate];
          }
          
          const totalItems = updatedItems.reduce((sum, i) => sum + i.quantity, 0);
          const totalAmount = updatedItems.reduce((sum, i) => sum + i.subtotal, 0);
          return { items: updatedItems, totalItems, totalAmount };
        }),
      updateQuantity: (itemId, newQuantity) =>
        set((state) => {
          const validItems = filterValidItems(state.items);
          let updatedItems: CartItem[] = [];
          if (newQuantity <= 0) {
            updatedItems = validItems.filter((i) => i.id !== itemId);
          } else {
            updatedItems = validItems.map((i) =>
              i.id === itemId
                ? { ...i, quantity: newQuantity, subtotal: newQuantity * i.unit_price }
                : i
            );
          }
          const totalItems = updatedItems.reduce((sum, i) => sum + i.quantity, 0);
          const totalAmount = updatedItems.reduce((sum, i) => sum + i.subtotal, 0);
          return { items: updatedItems, totalItems, totalAmount };
        }),
      removeItem: (itemId) =>
        set((state) => {
          const validItems = filterValidItems(state.items);
          const updatedItems = validItems.filter((i) => i.id !== itemId);
          const totalItems = updatedItems.reduce((sum, i) => sum + i.quantity, 0);
          const totalAmount = updatedItems.reduce((sum, i) => sum + i.subtotal, 0);
          return { items: updatedItems, totalItems, totalAmount };
        }),
      removeMultipleItems: (itemIds) =>
        set((state) => {
          const validItems = filterValidItems(state.items);
          const updatedItems = validItems.filter((i) => !itemIds.includes(i.id));
          const totalItems = updatedItems.reduce((sum, i) => sum + i.quantity, 0);
          const totalAmount = updatedItems.reduce((sum, i) => sum + i.subtotal, 0);
          return { items: updatedItems, totalItems, totalAmount };
        }),
      clearCart: () => set({ items: [], totalItems: 0, totalAmount: 0 }),
      cleanExpiredItems: () => {
        const currentItems = get().items;
        const validItems = filterValidItems(currentItems);
        const removedCount = currentItems.length - validItems.length;
        if (removedCount > 0) {
          const totalItems = validItems.reduce((sum, i) => sum + i.quantity, 0);
          const totalAmount = validItems.reduce((sum, i) => sum + i.subtotal, 0);
          set({ items: validItems, totalItems, totalAmount });
        }
        return removedCount;
      },
    }),
    { name: 'luckymart-customer-cart' }
  )
);
