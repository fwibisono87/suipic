import { writable } from 'svelte/store';

const THEME_KEY = 'suipic_theme';

type Theme = 'light' | 'dark';

const createThemeStore = () => {
	const { subscribe, set } = writable<Theme>('light');

	return {
		subscribe,
		setTheme: (theme: Theme) => {
			if (typeof window !== 'undefined') {
				localStorage.setItem(THEME_KEY, theme);
				document.documentElement.setAttribute('data-theme', theme);
			}
			set(theme);
		},
		loadFromStorage: () => {
			if (typeof window !== 'undefined') {
				const theme = (localStorage.getItem(THEME_KEY) as Theme) || 'light';
				document.documentElement.setAttribute('data-theme', theme);
				set(theme);
			}
		},
		toggle: () => {
			if (typeof window !== 'undefined') {
				const current = localStorage.getItem(THEME_KEY) as Theme;
				const newTheme: Theme = current === 'dark' ? 'light' : 'dark';
				localStorage.setItem(THEME_KEY, newTheme);
				document.documentElement.setAttribute('data-theme', newTheme);
				set(newTheme);
			}
		}
	};
};

export const themeStore = createThemeStore();
