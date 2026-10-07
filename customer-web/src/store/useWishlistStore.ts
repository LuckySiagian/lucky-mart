import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface WishlistState {
  wishlistIds: string[];
  toggleWishlist: (productId: string) => void;
  isInWishlist: (productId: string) => boolean;
}

export const useWishlistStore = create<WishlistState>()(
  persist(
    (set, get) => ({
      wishlistIds: [],
      toggleWishlist: (productId: string) => {
        const ids = get().wishlistIds;
        if (ids.includes(productId)) {
          set({ wishlistIds: ids.filter((id) => id !== productId) });
        } else {
          set({ wishlistIds: [...ids, productId] });
        }
      },
      isInWishlist: (productId: string) => get().wishlistIds.includes(productId),
    }),
    { name: 'luckymart-customer-wishlist' }
  )
);
