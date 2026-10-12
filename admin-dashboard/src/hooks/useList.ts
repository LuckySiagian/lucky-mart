import { useCallback, useEffect, useRef, useState } from 'react';
import type { Paginated } from '../types';
import type { ListParams } from '../services/apiServices';
import { getErrorMessage } from '../lib/utils';

/**
 * Loads a paginated API list and keeps page / search / filters in state.
 * Search input is debounced so typing does not fire a request per key.
 */
export function useList<T>(fetcher: (params: ListParams) => Promise<Paginated<T>>, initial: ListParams = {}) {
  const [params, setParams] = useState<ListParams>({ page: 1, limit: 20, ...initial });
  const [items, setItems] = useState<T[]>([]);
  const [meta, setMeta] = useState<Paginated<T>['meta'] | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchInput, setSearchInput] = useState(initial.search?.toString() || '');
  const fetcherRef = useRef(fetcher);
  fetcherRef.current = fetcher;

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await fetcherRef.current(params);
      setItems(res.items || []);
      setMeta(res.meta);
    } catch (err) {
      setError(getErrorMessage(err, 'Gagal memuat data dari server.'));
    } finally {
      setLoading(false);
    }
  }, [params]);

  useEffect(() => {
    load();
  }, [load]);

  useEffect(() => {
    const t = setTimeout(() => {
      setParams((p) => ((p.search ?? '') === searchInput ? p : { ...p, search: searchInput, page: 1 }));
    }, 350);
    return () => clearTimeout(t);
  }, [searchInput]);

  const setFilter = (key: string, value: string | number | undefined) => setParams((p) => ({ ...p, [key]: value, page: 1 }));
  const setPage = (page: number) => setParams((p) => ({ ...p, page }));

  return { items, setItems, meta, loading, error, reload: load, params, setFilter, setPage, searchInput, setSearchInput };
}
