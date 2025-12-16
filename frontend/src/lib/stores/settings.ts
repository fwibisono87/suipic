import { writable } from 'svelte/store';

type SettingsState = {
	imageProtectionEnabled: boolean;
	isLoaded: boolean;
};

const createSettingsStore = () => {
	const { subscribe, set, update } = writable<SettingsState>({
		imageProtectionEnabled: false,
		isLoaded: false
	});

	return {
		subscribe,
		setImageProtection: (enabled: boolean) => {
			update((state) => ({ ...state, imageProtectionEnabled: enabled, isLoaded: true }));
		},
		reset: () => {
			set({ imageProtectionEnabled: false, isLoaded: false });
		}
	};
};

export const settingsStore = createSettingsStore();
