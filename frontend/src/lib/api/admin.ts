import { config } from '$lib/config';
import type { TUser } from '$lib/types';

const API_URL = config.apiUrl;

export type TCreatePhotographerRequest = {
	email: string;
	username: string;
};

export type TCreatePhotographerResponse = {
	user: TUser;
	password: string;
};

export type TListPhotographersResponse = {
	photographers: TUser[];
};

class AdminApiError extends Error {
	constructor(message: string) {
		super(message);
		this.name = 'AdminApiError';
	}
}

export const adminApi = {
	async createPhotographer(
		data: TCreatePhotographerRequest,
		token: string
	): Promise<TCreatePhotographerResponse> {
		const response = await fetch(`${API_URL}/admin/photographers`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${token}`
			},
			body: JSON.stringify(data)
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({ message: 'Failed to create photographer' }));
			throw new AdminApiError(error.error || error.message || 'Failed to create photographer');
		}

		return response.json();
	},

	async listPhotographers(token: string): Promise<TListPhotographersResponse> {
		const response = await fetch(`${API_URL}/admin/photographers`, {
			method: 'GET',
			headers: {
				Authorization: `Bearer ${token}`
			}
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({ message: 'Failed to fetch photographers' }));
			throw new AdminApiError(error.error || error.message || 'Failed to fetch photographers');
		}

		return response.json();
	}
};
