import { config } from '$lib/config';

const API_URL = config.apiUrl;

export type TClient = {
	id: number;
	username: string;
	email: string;
	friendlyName: string;
	role: string;
	createdAt: string;
	isShared?: boolean;
};

export type TCreateClientRequest = {
	username: string;
	email?: string;
	password?: string;
	friendlyName?: string;
};

class PhotographerApiError extends Error {
	constructor(message: string) {
		super(message);
		this.name = 'PhotographerApiError';
	}
}

export const photographerApi = {
	async createOrLinkClient(data: TCreateClientRequest, token: string): Promise<TClient> {
		const response = await fetch(`${API_URL}/photographer/clients`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${token}`
			},
			body: JSON.stringify(data)
		});

		if (!response.ok) {
			const error = await response
				.json()
				.catch(() => ({ message: 'Failed to create or link client' }));
			throw new PhotographerApiError(error.error || error.message || 'Failed to create or link client');
		}

		return response.json();
	},

	async listClients(token: string): Promise<TClient[]> {
		const response = await fetch(`${API_URL}/photographer/clients`, {
			method: 'GET',
			headers: {
				Authorization: `Bearer ${token}`
			}
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({ message: 'Failed to fetch clients' }));
			throw new PhotographerApiError(error.error || error.message || 'Failed to fetch clients');
		}

		return response.json();
	},

	async searchClients(query: string, token: string): Promise<TClient[]> {
		const response = await fetch(`${API_URL}/photographer/clients/search?q=${encodeURIComponent(query)}`, {
			method: 'GET',
			headers: {
				Authorization: `Bearer ${token}`
			}
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({ message: 'Failed to search clients' }));
			throw new PhotographerApiError(error.error || error.message || 'Failed to search clients');
		}

		return response.json();
	}
};
