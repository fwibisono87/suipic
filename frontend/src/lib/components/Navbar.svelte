<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import Icon from '@iconify/svelte';
	import { authStore, currentUser, isAuthenticated, themeStore } from '$lib/stores';
	import { authApi } from '$lib/api';
	import { EUserRole } from '$lib/types';

	const handleLogout = async () => {
		try {
			const token = localStorage.getItem('suipic_token');
			if (token) {
				await authApi.logout(token);
			}
		} catch (err) {
			console.error('Logout error:', err);
		} finally {
			authStore.clearAuth();
			goto('/login');
		}
	};

	const toggleTheme = () => {
		themeStore.toggle();
	};
</script>

<nav class="navbar bg-base-200 shadow-lg mb-4">
	<div class="navbar-start">
		<div class="dropdown">
			<label tabindex="0" class="btn btn-ghost lg:hidden">
				<Icon icon="mdi:menu" class="text-2xl" />
			</label>
			<ul
				tabindex="0"
				class="menu menu-sm dropdown-content mt-3 z-[1] p-2 shadow bg-base-100 rounded-box w-52"
			>
				{#if $isAuthenticated}
					<li><a href="/">Home</a></li>
					<li><a href="/albums">Albums</a></li>
					{#if $currentUser?.role === EUserRole.PHOTOGRAPHER || $currentUser?.role === EUserRole.ADMIN}
						<li><a href="/albums/new">New Album</a></li>
					{/if}
					{#if $currentUser?.role === EUserRole.ADMIN}
						<li><a href="/admin">Admin</a></li>
					{/if}
				{/if}
			</ul>
		</div>
		<a href="/" class="btn btn-ghost text-xl">
			<Icon icon="mdi:camera" class="text-2xl" />
			Suipic
		</a>
	</div>

	<div class="navbar-center hidden lg:flex">
		{#if $isAuthenticated}
			<ul class="menu menu-horizontal px-1">
				<li><a href="/" class:active={$page.url.pathname === '/'}>Home</a></li>
				<li><a href="/albums" class:active={$page.url.pathname.startsWith('/albums')}>Albums</a></li>
				<li><a href="/search" class:active={$page.url.pathname === '/search'}>Search</a></li>
				{#if $currentUser?.role === EUserRole.PHOTOGRAPHER || $currentUser?.role === EUserRole.ADMIN}
					<li><a href="/albums/new" class:active={$page.url.pathname === '/albums/new'}>New Album</a></li>
				{/if}
				{#if $currentUser?.role === EUserRole.ADMIN}
					<li><a href="/admin" class:active={$page.url.pathname === '/admin'}>Admin</a></li>
				{/if}
			</ul>
		{/if}
	</div>

	<div class="navbar-end gap-2">
		{#if $isAuthenticated}
			<a href="/search" class="btn btn-ghost btn-circle" aria-label="Search">
				<Icon icon="mdi:magnify" class="text-2xl" />
			</a>
		{/if}
		
		<button class="btn btn-ghost btn-circle" on:click={toggleTheme}>
			<Icon icon="mdi:theme-light-dark" class="text-2xl" />
		</button>

		{#if $isAuthenticated}
			<div class="dropdown dropdown-end">
				<label tabindex="0" class="btn btn-ghost btn-circle avatar placeholder">
					<div class="bg-neutral text-neutral-content rounded-full w-10">
						<span class="text-xl">{$currentUser?.username?.charAt(0).toUpperCase()}</span>
					</div>
				</label>
				<ul
					tabindex="0"
					class="menu menu-sm dropdown-content mt-3 z-[1] p-2 shadow bg-base-100 rounded-box w-52"
				>
					<li class="menu-title">
						<span>{$currentUser?.username}</span>
						<span class="text-xs opacity-60">{$currentUser?.role}</span>
					</li>
					<li><a href="/profile">Profile</a></li>
					<li><a href="/settings">Settings</a></li>
					<li><button on:click={handleLogout}>Logout</button></li>
				</ul>
			</div>
		{:else}
			<a href="/login" class="btn btn-primary btn-sm">Login</a>
			<a href="/register" class="btn btn-outline btn-sm">Register</a>
		{/if}
	</div>
</nav>
