<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import Icon from '@iconify/svelte';
	import { isAuthenticated, currentUser } from '$lib/stores';
	import { EUserRole } from '$lib/types';
	import { Card } from '$lib/components';

	onMount(() => {
		if (!$isAuthenticated) {
			goto('/login');
			return;
		}

		if ($currentUser?.role !== EUserRole.ADMIN) {
			goto('/');
		}
	});
</script>

<svelte:head>
	<title>Admin - Suipic</title>
</svelte:head>

{#if $isAuthenticated && $currentUser?.role === EUserRole.ADMIN}
	<div class="space-y-6">
		<h1 class="text-4xl font-bold">Admin Dashboard</h1>

		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
			<div class="stats shadow">
				<div class="stat">
					<div class="stat-title">Total Users</div>
					<div class="stat-value text-primary">-</div>
					<div class="stat-desc">All registered users</div>
				</div>
			</div>

			<div class="stats shadow">
				<div class="stat">
					<div class="stat-title">Photographers</div>
					<div class="stat-value text-secondary">-</div>
					<div class="stat-desc">Active photographers</div>
				</div>
			</div>

			<div class="stats shadow">
				<div class="stat">
					<div class="stat-title">Total Albums</div>
					<div class="stat-value text-accent">-</div>
					<div class="stat-desc">Across all users</div>
				</div>
			</div>

			<div class="stats shadow">
				<div class="stat">
					<div class="stat-title">Total Photos</div>
					<div class="stat-value">-</div>
					<div class="stat-desc">All uploaded photos</div>
				</div>
			</div>
		</div>

		<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
			<Card title="User Management">
				<p class="text-sm opacity-60 mb-4">Manage users and their permissions</p>
				<button class="btn btn-primary btn-sm" disabled>
					<Icon icon="mdi:account-multiple" class="text-lg" />
					View Users
				</button>
			</Card>

			<Card title="Album Management">
				<p class="text-sm opacity-60 mb-4">Manage all albums across the platform</p>
				<button class="btn btn-primary btn-sm" disabled>
					<Icon icon="mdi:image-multiple" class="text-lg" />
					View Albums
				</button>
			</Card>

			<Card title="System Settings">
				<p class="text-sm opacity-60 mb-4">Configure system-wide settings</p>
				<button class="btn btn-primary btn-sm" disabled>
					<Icon icon="mdi:cog" class="text-lg" />
					Settings
				</button>
			</Card>

			<Card title="Storage">
				<p class="text-sm opacity-60 mb-4">Manage storage and backups</p>
				<div class="space-y-2">
					<div class="flex items-center justify-between text-sm">
						<span>Used Storage</span>
						<span class="font-semibold">-</span>
					</div>
					<progress class="progress progress-primary w-full" value="0" max="100"></progress>
				</div>
			</Card>
		</div>

		<Card title="Recent Activity">
			<div class="overflow-x-auto">
				<table class="table table-zebra w-full">
					<thead>
						<tr>
							<th>User</th>
							<th>Action</th>
							<th>Time</th>
						</tr>
					</thead>
					<tbody>
						<tr>
							<td colspan="3" class="text-center opacity-60">No recent activity</td>
						</tr>
					</tbody>
				</table>
			</div>
		</Card>
	</div>
{/if}
