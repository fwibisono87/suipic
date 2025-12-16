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
	},

	async createBatch(albumId: number, files: File[]): Promise<TPhoto[]> {
		const photos: TPhoto[] = [];
		const errors: string[] = [];

		for (const file of files) {
			try {
				const photo = await this.create(albumId, file);
				photos.push(photo);
			} catch (err) {
				errors.push(`${file.name}: ${(err as Error).message}`);
			}
		}

		if (errors.length > 0) {
			throw new PhotosApiError(`Some uploads failed: ${errors.join(', ')}`);
		}

		return photos;
	},

	async update(photoId: number, updates: Partial<{ title: string | null; pickRejectState: string | null; stars: number }>): Promise<TPhoto> {
		const response = await fetch(`${API_URL}/photos/${photoId}`, {
			method: 'PUT',
			headers: {
				...getAuthHeaders(),
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(updates)
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({ message: 'Failed to update photo' }));
			throw new PhotosApiError(error.message || 'Failed to update photo');
		}

		return response.json();
	}
};
