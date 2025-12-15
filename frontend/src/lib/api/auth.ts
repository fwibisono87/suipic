import type { TLoginRequest, TRegisterRequest, TAuthResponse, TUser, TRefreshRequest } from '$lib/types';
import { config } from '$lib/config';

const API_URL = config.apiUrl;

class AuthApiError extends Error {
	constructor(message: string) {
		super(message);
		this.name = 'AuthApiError';
	}
}

export const authApi = {
	async login(data: TLoginRequest): Promise<TAuthResponse> {
		const response = await fetch(`${API_URL}/auth/login`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(data)
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({ message: 'Login failed' }));
			throw new AuthApiError(error.message || 'Login failed');
		}

		return response.json();
	},

	async register(data: TRegisterRequest): Promise<TAuthResponse> {
		const response = await fetch(`${API_URL}/auth/register`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(data)
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({ message: 'Registration failed' }));
			throw new AuthApiError(error.message || 'Registration failed');
		}

		return response.json();
	},

	async refresh(data: TRefreshRequest): Promise<TAuthResponse> {
		const response = await fetch(`${API_URL}/auth/refresh`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(data)
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({ message: 'Token refresh failed' }));
			throw new AuthApiError(error.message || 'Token refresh failed');
		}

		return response.json();
	},

	async logout(token: string): Promise<void> {
		const response = await fetch(`${API_URL}/auth/logout`, {
			method: 'POST',
			headers: {
				'Authorization': `Bearer ${token}`
			}
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({ message: 'Logout failed' }));
			throw new AuthApiError(error.message || 'Logout failed');
		}
	},

	async me(token: string): Promise<TUser> {
		const response = await fetch(`${API_URL}/auth/me`, {
			method: 'GET',
			headers: {
				'Authorization': `Bearer ${token}`
			}
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({ message: 'Failed to fetch user' }));
			throw new AuthApiError(error.message || 'Failed to fetch user');
		}

		return response.json();
	}
};
