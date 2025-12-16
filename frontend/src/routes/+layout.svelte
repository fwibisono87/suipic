<script lang="ts">
	import { onMount } from 'svelte';
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';
	import { authStore, themeStore, settingsStore } from '$lib/stores';
	import { settingsApi } from '$lib/api/settings';
	import { Navbar, Footer } from '$lib/components';
	import type { PageData } from './$types';
	import '../app.css';

	export let data: PageData;

	const queryClient = new QueryClient({
		defaultOptions: {
			queries: {
				staleTime: 1000 * 60 * 5,
				gcTime: 1000 * 60 * 10,
				refetchOnWindowFocus: false,
				refetchOnReconnect: true,
				refetchOnMount: true,
				retry: (failureCount, error: any) => {
					if (error?.response?.status === 401 || error?.response?.status === 403) {
						return false;
					}
					if (error?.response?.status >= 400 && error?.response?.status < 500) {
						return false;
					}
					return failureCount < 3;
				},
				retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 30000),
				networkMode: 'online'
			},
			mutations: {
				retry: (failureCount, error: any) => {
					if (error?.response?.status === 401 || error?.response?.status === 403) {
						return false;
					}
					if (error?.response?.status >= 400 && error?.response?.status < 500) {
						return false;
					}
					return failureCount < 2;
				},
				retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 30000),
				networkMode: 'online',
				onError: (error: any) => {
					console.error('Mutation error:', error);
					if (error?.response?.status === 401) {
						authStore.clearAuth();
						window.location.href = '/login';
					}
				}
			}
		},
		queryCache: undefined,
		mutationCache: undefined
	});

	onMount(async () => {
		authStore.loadFromStorage();
		themeStore.loadFromStorage();

		if (data.user && data.isAuthenticated) {
			const token = localStorage.getItem('suipic_token');
			if (token) {
				authStore.setAuth(data.user, token);
			}
		}

		try {
			const publicSettings = await settingsApi.fetchPublicSettings();
			settingsStore.setImageProtection(publicSettings.image_protection_enabled);
		} catch (error) {
			console.error('Failed to load public settings:', error);
		}
	});
</script>

<QueryClientProvider client={queryClient}>
	<div class="min-h-screen flex flex-col">
		<Navbar />
		<main class="flex-1 container mx-auto px-4 py-8 max-w-7xl">
			<slot />
		</main>
		<Footer />
	</div>
</QueryClientProvider>
