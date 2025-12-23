import type { TSystemSettings, TUpdateSettingRequest } from '$lib/types';
import { authToken } from '$lib/stores';
import { get } from 'svelte/store';
import { config } from '$lib/config';

const API_URL = config.apiUrl;

class SettingsApiError extends Error {
	constructor(message: string) {
		super(message);
		this.name = 'SettingsApiError';
	}
}

const getAuthHeaders = () => {
	const token = get(authToken);
	return {
		'Content-Type': 'application/json',
		...(token ? { Authorization: `Bearer ${token}` } : {})
	};
};

export type TPublicSettings = {
	image_protection_enabled: boolean;
};

export const settingsApi = {
	async getSettings(): Promise<TSystemSettings[]> {
		const response = await fetch(`${API_URL}/admin/settings`, {
			method: 'GET',
			headers: getAuthHeaders()
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({ message: 'Failed to fetch settings' }));
			throw new SettingsApiError(error.message || 'Failed to fetch settings');
		}

		return response.json();
	},

	async updateSetting(settingKey: string, data: TUpdateSettingRequest): Promise<TSystemSettings> {
		const response = await fetch(`${API_URL}/admin/settings/${settingKey}`, {
			method: 'PUT',
			headers: getAuthHeaders(),
			body: JSON.stringify(data)
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({ message: 'Failed to update setting' }));
			throw new SettingsApiError(error.message || 'Failed to update setting');
		}

		return response.json();
	},

	async fetchPublicSettings(): Promise<TPublicSettings> {
		const response = await fetch(`${API_URL}/settings/public`, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			}
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({ message: 'Failed to fetch public settings' }));
			throw new SettingsApiError(error.message || 'Failed to fetch public settings');
		}

		return response.json();
	}
};
