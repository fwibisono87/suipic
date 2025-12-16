import type { TPhoto } from '$lib/types';
import { authToken } from '$lib/stores';
import { get } from 'svelte/store';
import { config } from '$lib/config';

const API_URL = config.apiUrl;

class SearchApiError extends Error {
	constructor(message: string) {
		super(message);
		this.name = 'SearchApiError';
	}
}

const getAuthHeaders = () => {
	const token = get(authToken);
	return {
		...(token ? { Authorization: `Bearer ${token}` } : {})
	};
};

export type TSearchParams = {
	q?: string;
	album?: number;
	dateFrom?: string;
	dateTo?: string;
	minStars?: number;
	maxStars?: number;
	state?: string;
	limit?: number;
	offset?: number;
};

export type TSearchResponse = {
	total: number;
	photos: TPhoto[];
};

export const searchApi = {
	async search(params: TSearchParams): Promise<TSearchResponse> {
		const queryParams = new URLSearchParams();

		if (params.q) queryParams.append('q', params.q);
		if (params.album) queryParams.append('album', params.album.toString());
		if (params.dateFrom) queryParams.append('dateFrom', params.dateFrom);
		if (params.dateTo) queryParams.append('dateTo', params.dateTo);
		if (params.minStars !== undefined) queryParams.append('minStars', params.minStars.toString());
		if (params.maxStars !== undefined) queryParams.append('maxStars', params.maxStars.toString());
		if (params.state) queryParams.append('state', params.state);
		if (params.limit) queryParams.append('limit', params.limit.toString());
		if (params.offset) queryParams.append('offset', params.offset.toString());

		const response = await fetch(`${API_URL}/search?${queryParams.toString()}`, {
			method: 'GET',
			headers: getAuthHeaders()
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({ message: 'Search failed' }));
			throw new SearchApiError(error.message || 'Search failed');
		}

		return response.json();
	}
};
