<script lang="ts">
	import { page } from '$app/stores';
</script>

<svelte:head>
	<title>Error - Suipic</title>
</svelte:head>

<div class="flex items-center justify-center min-h-[calc(100vh-300px)]">
	<div class="text-center">
		<h1 class="text-9xl font-bold text-error">{$page.status || 500}</h1>
		<p class="text-2xl mt-4">
			{$page.error?.message || 'An unexpected error occurred'}
		</p>
		<a href="/" class="btn btn-primary mt-8">Go Home</a>
	</div>
</div>
