import { create } from 'zustand';
import { cartService } from '../services/apiServices';
import type { Cart, CartItem } from '../types';

/**
 * Mirror of the customer's server-side cart. The server is the source of
 * truth (stock, price and wholesale rules are checked there); this store only
 * keeps the latest copy for the navbar badge and the cart / checkout pages.
 */
interface CartState {
  items: CartItem[];
  totalItems: number;
  totalAmount: number;
  loaded: boolean;
  selectedIds: string[];
  apply: (cart: Cart) => void;
  refresh: () => Promise<void>;
  setSelected: (ids: string[]) => void;
  reset: () => void;
}

export const useCartStore = create<CartState>()((set) => ({
  items: [],
  totalItems: 0,
  totalAmount: 0,
  loaded: false,
  selectedIds: [],
  apply: (cart) =>
    set((s) => {
      const items = cart.items || [];
      const ids = new Set(items.map((i) => i.id));
      return {
        items,
        totalItems: items.reduce((n, i) => n + i.quantity, 0),
        totalAmount: cart.total_amount,
        loaded: true,
        selectedIds: s.selectedIds.filter((id) => ids.has(id)),
      };
    }),
  refresh: async () => {
    const cart = await cartService.getCart();
    useCartStore.getState().apply(cart);
  },
  setSelected: (ids) => set({ selectedIds: ids }),
  reset: () => set({ items: [], totalItems: 0, totalAmount: 0, loaded: false, selectedIds: [] }),
}));
