import type { TPhoto } from '$lib/types';
import { authToken } from '$lib/stores';
import { get } from 'svelte/store';
import { config } from '$lib/config';

const API_URL = config.apiUrl;

class PhotosApiError extends Error {
	constructor(message: string) {
		super(message);
		this.name = 'PhotosApiError';
	}
}

const getAuthHeaders = () => {
	const token = get(authToken);
	return {
		...(token ? { Authorization: `Bearer ${token}` } : {})
	};
};

export const photosApi = {
	async listByAlbum(albumId: number): Promise<TPhoto[]> {
		const response = await fetch(`${API_URL}/albums/${albumId}/photos`, {
			method: 'GET',
			headers: getAuthHeaders()
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({ message: 'Failed to fetch photos' }));
			throw new PhotosApiError(error.message || 'Failed to fetch photos');
		}

		return response.json();
	},

	async create(albumId: number, file: File): Promise<TPhoto> {
		const formData = new FormData();
		formData.append('photo', file);

		const response = await fetch(`${API_URL}/albums/${albumId}/photos`, {
			method: 'POST',
			headers: getAuthHeaders(),
			body: formData
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({ message: 'Failed to upload photo' }));
			throw new PhotosApiError(error.message || 'Failed to upload photo');
		}

		return response.json();
	}
};
