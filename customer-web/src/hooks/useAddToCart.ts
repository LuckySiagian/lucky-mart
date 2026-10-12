import { useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { toast } from 'sonner';
import { cartService } from '../services/apiServices';
import { useAuthStore } from '../store/useAuthStore';
import { useCartStore } from '../store/useCartStore';
import { getErrorMessage } from '../lib/utils';
import type { Product, ProductUnit } from '../types';

/**
 * Adds a product unit to the server cart. Guests are sent to the login page
 * and come back here afterwards. With buyNow, only that item is selected and
 * the customer goes straight to checkout.
 */
export function useAddToCart() {
  const navigate = useNavigate();
  const location = useLocation();
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);
  const [adding, setAdding] = useState(false);

  const add = async (product: Product, unit: ProductUnit, quantity = 1, buyNow = false) => {
    if (!isAuthenticated) {
      toast.info('Masuk dulu untuk mulai belanja.');
      navigate(`/login?next=${encodeURIComponent(location.pathname + location.search)}`);
      return;
    }
    setAdding(true);
    try {
      const cart = await cartService.addToCart({ product_id: product.id, product_unit_id: unit.id, quantity });
      const store = useCartStore.getState();
      store.apply(cart);
      if (buyNow) {
        const item = (cart.items || []).find((i) => i.product_unit_id === unit.id);
        store.setSelected(item ? [item.id] : []);
        navigate('/checkout');
      } else {
        toast.success(`${product.name} (${unit.unit_name}) masuk keranjang`);
      }
    } catch (err) {
      toast.error(getErrorMessage(err, 'Gagal menambahkan ke keranjang.'));
    } finally {
      setAdding(false);
    }
  };

  return { add, adding };
}
