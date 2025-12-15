<script lang="ts">
	import { onMount } from 'svelte';
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';
	import { authStore, themeStore } from '$lib/stores';
	import { Navbar, Footer } from '$lib/components';
	import type { PageData } from './$types';
	import '../app.css';

	export let data: PageData;

	const queryClient = new QueryClient({
		defaultOptions: {
			queries: {
				staleTime: 1000 * 60 * 5,
				refetchOnWindowFocus: false
			}
		}
	});

	onMount(() => {
		authStore.loadFromStorage();
		themeStore.loadFromStorage();

		if (data.user && data.isAuthenticated) {
			const token = localStorage.getItem('suipic_token');
			if (token) {
				authStore.setAuth(data.user, token);
			}
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
