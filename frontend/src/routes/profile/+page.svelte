<script lang="ts">
	import { goto } from '$app/navigation';
	import { isAuthenticated, currentUser } from '$lib/stores';
	import { onMount } from 'svelte';
	import { formatDateTime } from '$lib/utils';

	onMount(() => {
		if (!$isAuthenticated) {
			goto('/login');
		}
	});
</script>

<svelte:head>
	<title>Profile - Suipic</title>
</svelte:head>

{#if $isAuthenticated && $currentUser}
	<div class="max-w-4xl mx-auto">
		<h1 class="text-4xl font-bold mb-6">My Profile</h1>

		<div class="card bg-base-100 shadow-xl">
			<div class="card-body">
				<div class="flex items-center gap-6">
					<div class="avatar placeholder">
						<div class="bg-neutral text-neutral-content rounded-full w-24">
							<span class="text-4xl">{$currentUser.username[0].toUpperCase()}</span>
						</div>
					</div>
					<div>
						<h2 class="text-2xl font-bold">{$currentUser.friendlyName || $currentUser.username}</h2>
						<p class="text-sm opacity-60">@{$currentUser.username}</p>
						<div class="badge badge-primary mt-2">{$currentUser.role}</div>
					</div>
				</div>

				<div class="divider"></div>

				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div class="form-control">
						<label class="label">
							<span class="label-text font-semibold">Email</span>
						</label>
						<p class="text-base">{$currentUser.email}</p>
					</div>

					<div class="form-control">
						<label class="label">
							<span class="label-text font-semibold">Username</span>
						</label>
						<p class="text-base">{$currentUser.username}</p>
					</div>

					<div class="form-control">
						<label class="label">
							<span class="label-text font-semibold">Member Since</span>
						</label>
						<p class="text-base">{formatDateTime($currentUser.createdAt)}</p>
					</div>

					<div class="form-control">
						<label class="label">
							<span class="label-text font-semibold">Last Updated</span>
						</label>
						<p class="text-base">{formatDateTime($currentUser.updatedAt)}</p>
					</div>
				</div>

				<div class="divider"></div>

				<div class="flex gap-2">
					<button class="btn btn-primary" disabled>Edit Profile</button>
					<button class="btn btn-outline" disabled>Change Password</button>
				</div>
			</div>
		</div>
	</div>
{/if}
