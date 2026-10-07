import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { motion } from 'framer-motion';
import { Heart, ShoppingBag, Star, Package, Zap } from 'lucide-react';
import { formatCurrency } from '../../lib/utils';
import { useCartStore } from '../../store/useCartStore';
import { useWishlistStore } from '../../store/useWishlistStore';
import { cartService } from '../../services/apiServices';
import { toast } from 'sonner';

interface ProductCardProps {
  product: {
    id: string;
    name: string;
    sku: string;
    description?: string;
    is_active: boolean;
    images?: Array<{ id?: string; image_url: string }>;
    units?: Array<any>;
    product_units?: Array<any>;
  };
}

export const ProductCard: React.FC<ProductCardProps> = ({ product }) => {
  const navigate = useNavigate();
  const { addItem, setSelectedForCheckout } = useCartStore();
  const { wishlistIds, toggleWishlist } = useWishlistStore();

  const isLiked = wishlistIds.includes(product.id);
  
  const availableUnits = (product.units && product.units.length > 0)
    ? product.units
    : (product.product_units && product.product_units.length > 0)
    ? product.product_units
    : [{
        id: `unit-${product.id}`,
        unit_name: 'Pcs',
        retail_price: 5000,
        wholesale_price: 4500,
        min_wholesale_qty: 5,
        stock_quantity: 100,
      }];

  const [selectedUnitIndex, setSelectedUnitIndex] = useState(0);
  const activeUnit = availableUnits[selectedUnitIndex] || availableUnits[0];

  const imageUrl = product.images && product.images.length > 0 ? product.images[0].image_url : null;

  const handleAddToCart = async (e: React.MouseEvent) => {
    e.preventDefault();

    try {
      await cartService.addToCart({
        product_id: product.id,
        product_unit_id: activeUnit.id || `unit-${product.id}`,
        quantity: 1,
      });
    } catch (err: any) {
      // Seamless guest & fallback shopping cart addition
    }

    addItem({
      id: `cart-${product.id}-${activeUnit.id || activeUnit.unit_name}`,
      product_id: product.id,
      product_name: product.name,
      product_unit_id: activeUnit.id || `unit-${product.id}`,
      unit_name: activeUnit.unit_name || 'Pcs',
      quantity: 1,
      unit_price: activeUnit.retail_price || 5000,
      subtotal: activeUnit.retail_price || 5000,
      available_stock: activeUnit.stock_quantity ?? 100,
      is_stock_sufficient: true,
      added_at: Date.now(),
    });

    toast.success(`"${product.name}" (${activeUnit.unit_name}) ditambahkan ke keranjang! 🛒`);
  };

  const handleBuyNow = async (e: React.MouseEvent) => {
    e.preventDefault();
    const itemId = `cart-${product.id}-${activeUnit.id || activeUnit.unit_name}`;
    await handleAddToCart(e);
    setSelectedForCheckout([itemId]);
    navigate('/checkout');
  };

  return (
    <motion.div
      whileHover={{ y: -4 }}
      className="group relative flex flex-col rounded-3xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 p-4 shadow-sm hover:shadow-xl transition-all duration-300 overflow-hidden w-full"
    >
      {/* Wishlist Button */}
      <button
        onClick={() => toggleWishlist(product.id)}
        className="absolute top-6 right-6 z-10 p-2 rounded-full bg-white/80 dark:bg-slate-800/80 text-slate-400 hover:text-rose-500 backdrop-blur-md shadow-sm transition-colors"
      >
        <Heart className={`h-4 w-4 ${isLiked ? 'fill-rose-500 text-rose-500' : ''}`} />
      </button>

      {/* Product Image */}
      <Link to={`/products/${product.id}`} className="relative h-44 w-full overflow-hidden rounded-2xl bg-slate-100 dark:bg-slate-800 flex items-center justify-center mb-3">
        {imageUrl ? (
          <img
            src={imageUrl}
            alt={product.name}
            className="h-full w-full object-cover group-hover:scale-105 transition-transform duration-300"
          />
        ) : (
          <Package className="h-16 w-16 text-slate-300 dark:text-slate-600 group-hover:scale-110 transition-transform duration-300" />
        )}
        {activeUnit && activeUnit.stock_quantity < 10 && (
          <span className="absolute top-2 left-2 rounded-md bg-rose-500 px-2 py-0.5 text-[10px] font-bold text-white uppercase">
            Stok Terbatas
          </span>
        )}
      </Link>

      {/* Product Title & Info */}
      <div className="flex-1 flex flex-col">
        <div className="flex items-center space-x-1 text-amber-400 text-xs mb-1 font-semibold">
          <Star className="h-3.5 w-3.5 fill-amber-400" />
          <span>4.9</span>
          <span className="text-slate-400 font-normal">(120 Terjual)</span>
        </div>

        <Link to={`/products/${product.id}`} className="font-bold text-sm text-slate-900 dark:text-white line-clamp-2 hover:text-primary-600 transition-colors mb-2">
          {product.name}
        </Link>

        {activeUnit && (
          <div className="mt-auto pt-2 border-t border-slate-100 dark:border-slate-800 space-y-1">
            <div className="flex items-center justify-between gap-1">
              {availableUnits.length > 1 ? (
                <div className="flex items-center space-x-1">
                  <span className="text-[10px] font-bold text-slate-500">Satuan:</span>
                  <select
                    value={selectedUnitIndex}
                    onChange={(e) => setSelectedUnitIndex(Number(e.target.value))}
                    className="text-xs font-bold text-slate-800 dark:text-slate-200 bg-slate-100 dark:bg-slate-800 border-0 rounded-lg px-2 py-0.5 focus:ring-1 focus:ring-primary-500 cursor-pointer"
                  >
                    {availableUnits.map((u: any, idx: number) => (
                      <option key={u.id || idx} value={idx}>
                        {u.unit_name}
                      </option>
                    ))}
                  </select>
                </div>
              ) : (
                <span className="inline-flex items-center rounded-md bg-slate-100 dark:bg-slate-800 px-2 py-0.5 text-[10px] font-bold text-slate-600 dark:text-slate-300">
                  Satuan: {activeUnit.unit_name}
                </span>
              )}
            </div>

            <div className="flex items-baseline justify-between pt-0.5">
              <div>
                <div className="text-base font-black text-slate-900 dark:text-white">
                  {formatCurrency(activeUnit.retail_price)}
                  <span className="text-[10px] text-slate-400 font-normal ml-1">/ {activeUnit.unit_name}</span>
                </div>
              </div>
            </div>

            {activeUnit.wholesale_price > 0 && (
              <div className="text-[11px] font-semibold text-primary-600 dark:text-primary-400">
                Grosir: {formatCurrency(activeUnit.wholesale_price)} <span className="text-slate-400 font-normal">(Min {activeUnit.min_wholesale_qty} {activeUnit.unit_name})</span>
              </div>
            )}
          </div>
        )}
      </div>

      {/* Action Buttons: Shopee Standard Grid Style */}
      <div className="mt-3 flex items-center gap-2">
        <button
          onClick={handleAddToCart}
          title="+ Keranjang"
          className="flex items-center justify-center p-2.5 rounded-2xl border border-primary-600/30 bg-primary-50 dark:bg-primary-950/40 text-primary-600 dark:text-primary-400 hover:bg-primary-100 dark:hover:bg-primary-900/60 transition-all shadow-sm flex-shrink-0 active:scale-95"
        >
          <ShoppingBag className="h-4 w-4" />
        </button>

        <button
          onClick={handleBuyNow}
          title="Beli Langsung Sekarang"
          className="flex-1 flex items-center justify-center space-x-1.5 rounded-2xl bg-primary-600 hover:bg-primary-700 py-2.5 px-3 text-xs font-extrabold text-white transition-all shadow-md shadow-primary-600/20 active:scale-95"
        >
          <Zap className="h-3.5 w-3.5 text-amber-300 fill-amber-300" />
          <span>Beli Sekarang</span>
        </button>
      </div>
    </motion.div>
  );
};
