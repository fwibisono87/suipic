import type { TComment } from '$lib/types';
import { authToken } from '$lib/stores';
import { get } from 'svelte/store';
import { config } from '$lib/config';

const API_URL = config.apiUrl;

class CommentsApiError extends Error {
	constructor(message: string) {
		super(message);
		this.name = 'CommentsApiError';
	}
}

const getAuthHeaders = () => {
	const token = get(authToken);
	return {
		...(token ? { Authorization: `Bearer ${token}` } : {})
	};
};

export const commentsApi = {
	async getByPhoto(photoId: number): Promise<TComment[]> {
		const response = await fetch(`${API_URL}/photos/${photoId}/comments`, {
			method: 'GET',
			headers: getAuthHeaders()
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({ message: 'Failed to fetch comments' }));
			throw new CommentsApiError(error.message || 'Failed to fetch comments');
		}

		return response.json();
	},

	async create(photoId: number, text: string, parentCommentId?: number | null): Promise<TComment> {
		const response = await fetch(`${API_URL}/photos/${photoId}/comments`, {
			method: 'POST',
			headers: {
				...getAuthHeaders(),
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				text,
				...(parentCommentId ? { parentCommentId } : {})
			})
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({ message: 'Failed to create comment' }));
			throw new CommentsApiError(error.message || 'Failed to create comment');
		}

		return response.json();
	}
};
