<script lang="ts">
	import { goto } from '$app/navigation';
	import { isAuthenticated, themeStore } from '$lib/stores';
	import { onMount } from 'svelte';

	let theme: 'light' | 'dark' = 'light';

	onMount(() => {
		if (!$isAuthenticated) {
			goto('/login');
		}

		themeStore.subscribe((value) => {
			theme = value;
		});
	});

	const handleThemeChange = (e: Event) => {
		const target = e.target as HTMLSelectElement;
		themeStore.setTheme(target.value as 'light' | 'dark');
	};
</script>

<svelte:head>
	<title>Settings - Suipic</title>
</svelte:head>

{#if $isAuthenticated}
	<div class="max-w-4xl mx-auto">
		<h1 class="text-4xl font-bold mb-6">Settings</h1>

		<div class="space-y-6">
			<div class="card bg-base-100 shadow-xl">
				<div class="card-body">
					<h2 class="card-title">Appearance</h2>
					<div class="form-control">
						<label class="label" for="theme">
							<span class="label-text">Theme</span>
						</label>
						<select
							id="theme"
							class="select select-bordered w-full max-w-xs"
							bind:value={theme}
							on:change={handleThemeChange}
						>
							<option value="light">Light</option>
							<option value="dark">Dark</option>
						</select>
					</div>
				</div>
			</div>

			<div class="card bg-base-100 shadow-xl">
				<div class="card-body">
					<h2 class="card-title">Notifications</h2>
					<div class="form-control">
						<label class="label cursor-pointer">
							<span class="label-text">Email notifications</span>
							<input type="checkbox" class="toggle toggle-primary" disabled />
						</label>
					</div>
					<div class="form-control">
						<label class="label cursor-pointer">
							<span class="label-text">Album share notifications</span>
							<input type="checkbox" class="toggle toggle-primary" disabled />
						</label>
					</div>
					<div class="form-control">
						<label class="label cursor-pointer">
							<span class="label-text">Comment notifications</span>
							<input type="checkbox" class="toggle toggle-primary" disabled />
						</label>
					</div>
				</div>
			</div>

			<div class="card bg-base-100 shadow-xl">
				<div class="card-body">
					<h2 class="card-title">Privacy</h2>
					<div class="form-control">
						<label class="label cursor-pointer">
							<span class="label-text">Public profile</span>
							<input type="checkbox" class="toggle toggle-primary" disabled />
						</label>
					</div>
					<div class="form-control">
						<label class="label cursor-pointer">
							<span class="label-text">Show email to other users</span>
							<input type="checkbox" class="toggle toggle-primary" disabled />
						</label>
					</div>
				</div>
			</div>

			<div class="card bg-base-100 shadow-xl border-error">
				<div class="card-body">
					<h2 class="card-title text-error">Danger Zone</h2>
					<p class="text-sm opacity-60">These actions cannot be undone</p>
					<div class="mt-4">
						<button class="btn btn-error btn-outline" disabled>Delete Account</button>
					</div>
				</div>
			</div>
		</div>
	</div>
{/if}
