import { writable, derived } from 'svelte/store';
import type { TUser } from '$lib/types';

const TOKEN_KEY = 'suipic_token';
const USER_KEY = 'suipic_user';

type AuthState = {
	user: TUser | null;
	token: string | null;
	isLoading: boolean;
};

const setCookie = (name: string, value: string, days = 7) => {
	if (typeof document !== 'undefined') {
		const expires = new Date();
		expires.setTime(expires.getTime() + days * 24 * 60 * 60 * 1000);
		document.cookie = `${name}=${value};expires=${expires.toUTCString()};path=/;SameSite=Strict`;
	}
};

const deleteCookie = (name: string) => {
	if (typeof document !== 'undefined') {
		document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:00 UTC;path=/;`;
	}
};

const createAuthStore = () => {
	const { subscribe, set, update } = writable<AuthState>({
		user: null,
		token: null,
		isLoading: true
	});

	return {
		subscribe,
		setAuth: (user: TUser, token: string) => {
			if (typeof window !== 'undefined') {
				localStorage.setItem(TOKEN_KEY, token);
				localStorage.setItem(USER_KEY, JSON.stringify(user));
				setCookie(TOKEN_KEY, token);
				setCookie(USER_KEY, JSON.stringify(user));
			}
			update((state) => ({ ...state, user, token, isLoading: false }));
		},
		clearAuth: () => {
			if (typeof window !== 'undefined') {
				localStorage.removeItem(TOKEN_KEY);
				localStorage.removeItem(USER_KEY);
				deleteCookie(TOKEN_KEY);
				deleteCookie(USER_KEY);
			}
			set({ user: null, token: null, isLoading: false });
		},
		loadFromStorage: () => {
			if (typeof window !== 'undefined') {
				const token = localStorage.getItem(TOKEN_KEY);
				const userStr = localStorage.getItem(USER_KEY);
				if (token && userStr) {
					try {
						const user = JSON.parse(userStr) as TUser;
						setCookie(TOKEN_KEY, token);
						setCookie(USER_KEY, userStr);
						update((state) => ({ ...state, user, token, isLoading: false }));
					} catch (e) {
						localStorage.removeItem(TOKEN_KEY);
						localStorage.removeItem(USER_KEY);
						deleteCookie(TOKEN_KEY);
						deleteCookie(USER_KEY);
						update((state) => ({ ...state, isLoading: false }));
					}
				} else {
					update((state) => ({ ...state, isLoading: false }));
				}
			}
		},
		updateUser: (user: TUser) => {
			if (typeof window !== 'undefined') {
				const userStr = JSON.stringify(user);
				localStorage.setItem(USER_KEY, userStr);
				setCookie(USER_KEY, userStr);
			}
			update((state) => ({ ...state, user }));
		}
	};
};

export const authStore = createAuthStore();

export const isAuthenticated = derived(authStore, ($auth) => $auth.user !== null && $auth.token !== null);

export const currentUser = derived(authStore, ($auth) => $auth.user);

export const isLoading = derived(authStore, ($auth) => $auth.isLoading);

export const authToken = derived(authStore, ($auth) => $auth.token);
