import type { TAlbum, TCreateAlbumRequest, TUpdateAlbumRequest } from '$lib/types';
import { authToken } from '$lib/stores';
import { get } from 'svelte/store';
import { config } from '$lib/config';

const API_URL = config.apiUrl;

class AlbumsApiError extends Error {
	constructor(message: string) {
		super(message);
		this.name = 'AlbumsApiError';
	}
}

const getAuthHeaders = () => {
	const token = get(authToken);
	return {
		'Content-Type': 'application/json',
		...(token ? { Authorization: `Bearer ${token}` } : {})
	};
};

export const albumsApi = {
	async list(): Promise<TAlbum[]> {
		const response = await fetch(`${API_URL}/albums`, {
			method: 'GET',
			headers: getAuthHeaders()
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({ message: 'Failed to fetch albums' }));
			throw new AlbumsApiError(error.message || 'Failed to fetch albums');
		}

		return response.json();
	},

	async get(id: number): Promise<TAlbum> {
		const response = await fetch(`${API_URL}/albums/${id}`, {
			method: 'GET',
			headers: getAuthHeaders()
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({ message: 'Failed to fetch album' }));
			throw new AlbumsApiError(error.message || 'Failed to fetch album');
		}

		return response.json();
	},

	async create(data: TCreateAlbumRequest): Promise<TAlbum> {
		const response = await fetch(`${API_URL}/albums`, {
			method: 'POST',
			headers: getAuthHeaders(),
			body: JSON.stringify(data)
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({ message: 'Failed to create album' }));
			throw new AlbumsApiError(error.message || 'Failed to create album');
		}

		return response.json();
	},

	async update(id: number, data: TUpdateAlbumRequest): Promise<TAlbum> {
		const response = await fetch(`${API_URL}/albums/${id}`, {
			method: 'PUT',
			headers: getAuthHeaders(),
			body: JSON.stringify(data)
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({ message: 'Failed to update album' }));
			throw new AlbumsApiError(error.message || 'Failed to update album');
		}

		return response.json();
	},

	async delete(id: number): Promise<void> {
		const response = await fetch(`${API_URL}/albums/${id}`, {
			method: 'DELETE',
			headers: getAuthHeaders()
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({ message: 'Failed to delete album' }));
			throw new AlbumsApiError(error.message || 'Failed to delete album');
		}
	}
};
