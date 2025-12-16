import { createQuery } from '@tanstack/svelte-query';
import { derived, writable, type Readable } from 'svelte/store';
import { goto } from '$app/navigation';
import { page } from '$app/stores';
import { get } from 'svelte/store';
import { searchApi, type TSearchParams, type TSearchResponse } from '$lib/api/search';
import { queryKeys } from './keys';

export type SearchFilters = {
	q: string;
	album?: number;
	dateFrom?: string;
	dateTo?: string;
	minStars?: number;
	maxStars?: number;
	state?: string;
};

export type PaginationState = {
	page: number;
	pageSize: number;
};

type SearchQueryOptions = {
	debounceMs?: number;
	enabled?: boolean;
};

export function useSearchQuery(options: SearchQueryOptions = {}) {
	const { debounceMs = 500, enabled = true } = options;

	const filters = writable<SearchFilters>({
		q: '',
		album: undefined,
		dateFrom: undefined,
		dateTo: undefined,
		minStars: undefined,
		maxStars: undefined,
		state: undefined,
	});

	const pagination = writable<PaginationState>({
		page: 1,
		pageSize: 50,
	});

	let debounceTimeout: ReturnType<typeof setTimeout> | null = null;
	const debouncedFilters = writable<SearchFilters>(get(filters));

	filters.subscribe((value) => {
		if (debounceTimeout) {
			clearTimeout(debounceTimeout);
		}
		debounceTimeout = setTimeout(() => {
			debouncedFilters.set(value);
			pagination.update((p) => ({ ...p, page: 1 }));
		}, debounceMs);
	});

	const searchParams: Readable<TSearchParams> = derived(
		[debouncedFilters, pagination],
		([$filters, $pagination]) => {
			const params: TSearchParams = {
				limit: $pagination.pageSize,
				offset: ($pagination.page - 1) * $pagination.pageSize,
			};

			if ($filters.q?.trim()) {
				params.q = $filters.q.trim();
			}

			if ($filters.album !== undefined) {
				params.album = $filters.album;
			}

			if ($filters.dateFrom) {
				params.dateFrom = new Date($filters.dateFrom).toISOString();
			}

			if ($filters.dateTo) {
				const date = new Date($filters.dateTo);
				date.setHours(23, 59, 59, 999);
				params.dateTo = date.toISOString();
			}

			if ($filters.minStars !== undefined && $filters.minStars > 0) {
				params.minStars = $filters.minStars;
			}

			if ($filters.maxStars !== undefined && $filters.maxStars < 5) {
				params.maxStars = $filters.maxStars;
			}

			if ($filters.state) {
				params.state = $filters.state;
			}

			return params;
		}
	);

	const query = createQuery(
		derived(searchParams, ($params) => ({
			queryKey: queryKeys.search.query($params),
			queryFn: () => searchApi.search($params),
			enabled: enabled && ($params.q !== undefined || $params.album !== undefined),
		}))
	);

	const totalPages = derived(query, ($query) => {
		if (!$query.data) return 0;
		const paginationState = get(pagination);
		return Math.ceil($query.data.total / paginationState.pageSize);
	});

	const hasActiveFilters = derived(filters, ($filters) => {
		return (
			$filters.album !== undefined ||
			$filters.dateFrom !== undefined ||
			$filters.dateTo !== undefined ||
			($filters.minStars !== undefined && $filters.minStars > 0) ||
			($filters.maxStars !== undefined && $filters.maxStars < 5) ||
			$filters.state !== undefined
		);
	});

	function setFilter<K extends keyof SearchFilters>(key: K, value: SearchFilters[K]) {
		filters.update((f) => ({ ...f, [key]: value }));
	}

	function setFilters(newFilters: Partial<SearchFilters>) {
		filters.update((f) => ({ ...f, ...newFilters }));
	}

	function clearFilters() {
		filters.set({
			q: get(filters).q,
			album: undefined,
			dateFrom: undefined,
			dateTo: undefined,
			minStars: undefined,
			maxStars: undefined,
			state: undefined,
		});
	}

	function setPage(page: number) {
		pagination.update((p) => ({ ...p, page }));
	}

	function setPageSize(pageSize: number) {
		pagination.update((p) => ({ ...p, pageSize, page: 1 }));
	}

	function nextPage() {
		const total = get(totalPages);
		pagination.update((p) => {
			if (p.page < total) {
				return { ...p, page: p.page + 1 };
			}
			return p;
		});
	}

	function previousPage() {
		pagination.update((p) => {
			if (p.page > 1) {
				return { ...p, page: p.page - 1 };
			}
			return p;
		});
	}

	function syncFromQueryParams() {
		const pageStore = get(page);
		const params = pageStore.url.searchParams;

		const newFilters: SearchFilters = {
			q: params.get('q') || '',
			album: params.has('album') ? Number(params.get('album')) : undefined,
			dateFrom: params.get('dateFrom') || undefined,
			dateTo: params.get('dateTo') || undefined,
			minStars: params.has('minStars') ? Number(params.get('minStars')) : undefined,
			maxStars: params.has('maxStars') ? Number(params.get('maxStars')) : undefined,
			state: params.get('state') || undefined,
		};

		const newPagination: PaginationState = {
			page: params.has('page') ? Number(params.get('page')) : 1,
			pageSize: params.has('pageSize') ? Number(params.get('pageSize')) : 50,
		};

		filters.set(newFilters);
		pagination.set(newPagination);
		debouncedFilters.set(newFilters);
	}

	function syncToQueryParams() {
		const currentFilters = get(filters);
		const currentPagination = get(pagination);
		const pageStore = get(page);
		const params = new URLSearchParams(pageStore.url.searchParams);

		if (currentFilters.q) {
			params.set('q', currentFilters.q);
		} else {
			params.delete('q');
		}

		if (currentFilters.album !== undefined) {
			params.set('album', currentFilters.album.toString());
		} else {
			params.delete('album');
		}

		if (currentFilters.dateFrom) {
			params.set('dateFrom', currentFilters.dateFrom);
		} else {
			params.delete('dateFrom');
		}

		if (currentFilters.dateTo) {
			params.set('dateTo', currentFilters.dateTo);
		} else {
			params.delete('dateTo');
		}

		if (currentFilters.minStars !== undefined && currentFilters.minStars > 0) {
			params.set('minStars', currentFilters.minStars.toString());
		} else {
			params.delete('minStars');
		}

		if (currentFilters.maxStars !== undefined && currentFilters.maxStars < 5) {
			params.set('maxStars', currentFilters.maxStars.toString());
		} else {
			params.delete('maxStars');
		}

		if (currentFilters.state) {
			params.set('state', currentFilters.state);
		} else {
			params.delete('state');
		}

		if (currentPagination.page > 1) {
			params.set('page', currentPagination.page.toString());
		} else {
			params.delete('page');
		}

		if (currentPagination.pageSize !== 50) {
			params.set('pageSize', currentPagination.pageSize.toString());
		} else {
			params.delete('pageSize');
		}

		const newUrl = `${pageStore.url.pathname}?${params.toString()}`;
		goto(newUrl, { replaceState: true, keepFocus: true, noScroll: true });
	}

	return {
		query,
		filters,
		pagination,
		totalPages,
		hasActiveFilters,
		setFilter,
		setFilters,
		clearFilters,
		setPage,
		setPageSize,
		nextPage,
		previousPage,
		syncFromQueryParams,
		syncToQueryParams,
	};
}
