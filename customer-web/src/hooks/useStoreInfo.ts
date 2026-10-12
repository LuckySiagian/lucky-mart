import { useEffect, useState } from 'react';
import { storeService } from '../services/apiServices';
import type { StoreInfo } from '../types';

// Fetched once per page load and shared by every component that asks.
let cached: StoreInfo | null = null;
let pending: Promise<StoreInfo> | null = null;

/** Store details (address, WhatsApp, hours) set by the owner in the admin. */
export function useStoreInfo(): StoreInfo | null {
  const [info, setInfo] = useState<StoreInfo | null>(cached);
  useEffect(() => {
    if (cached) return;
    pending ??= storeService.get().then((s) => (cached = s));
    pending.then(setInfo).catch(() => {
      pending = null;
    });
  }, []);
  return info;
}
