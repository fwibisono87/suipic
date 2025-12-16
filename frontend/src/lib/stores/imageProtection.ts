import { writable } from 'svelte/store';

const IMAGE_PROTECTION_KEY = 'suipic_image_protection';

const createImageProtectionStore = () => {
	const { subscribe, set } = writable<boolean>(false);

	return {
		subscribe,
		setEnabled: (enabled: boolean) => {
			if (typeof window !== 'undefined') {
				localStorage.setItem(IMAGE_PROTECTION_KEY, String(enabled));
			}
			set(enabled);
		},
		loadFromStorage: () => {
			if (typeof window !== 'undefined') {
				const enabled = localStorage.getItem(IMAGE_PROTECTION_KEY) === 'true';
				set(enabled);
			}
		},
		toggle: () => {
			if (typeof window !== 'undefined') {
				const current = localStorage.getItem(IMAGE_PROTECTION_KEY) === 'true';
				const newValue = !current;
				localStorage.setItem(IMAGE_PROTECTION_KEY, String(newValue));
				set(newValue);
			}
		}
	};
};

export const imageProtectionEnabled = createImageProtectionStore();
