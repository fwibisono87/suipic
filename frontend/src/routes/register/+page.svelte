<script lang="ts">
	import { goto } from '$app/navigation';
	import { authStore, isAuthenticated } from '$lib/stores';
	import { authApi } from '$lib/api';
	import { Alert, LoadingSpinner } from '$lib/components';
	import { validateEmail, validateUsername, validatePassword } from '$lib/utils';
	import { EUserRole, type TRegisterRequest } from '$lib/types';
	import { onMount } from 'svelte';

	let email = '';
	let username = '';
	let password = '';
	let confirmPassword = '';
	let role: EUserRole = EUserRole.CLIENT;
	let isLoading = false;
	let error = '';

	onMount(() => {
		if ($isAuthenticated) {
			goto('/');
		}
	});

	const handleSubmit = async (e: Event) => {
		e.preventDefault();
		error = '';

		if (!email || !validateEmail(email)) {
			error = 'Please enter a valid email address';
			return;
		}

		if (!username || !validateUsername(username)) {
			error = 'Username must be between 3 and 50 characters';
			return;
		}

		if (!password || !validatePassword(password)) {
			error = 'Password must be at least 6 characters';
			return;
		}

		if (password !== confirmPassword) {
			error = 'Passwords do not match';
			return;
		}

		isLoading = true;

		try {
			const registerData: TRegisterRequest = {
				email,
				username,
				password,
				role
			};

			const response = await authApi.register(registerData);
			authStore.setAuth(response.user, response.token);
			goto('/');
		} catch (err: unknown) {
			error = (err as { message: string }).message || 'Registration failed';
		} finally {
			isLoading = false;
		}
	};
</script>

<svelte:head>
	<title>Register - Suipic</title>
</svelte:head>

<div class="flex items-center justify-center min-h-[calc(100vh-300px)]">
	<div class="card w-full max-w-md bg-base-100 shadow-xl">
		<div class="card-body">
			<h2 class="card-title text-3xl font-bold text-center justify-center mb-4">Register</h2>

			{#if error}
				<Alert type="error" message={error} dismissible onDismiss={() => (error = '')} />
			{/if}

			<form on:submit={handleSubmit} class="space-y-4">
				<div class="form-control">
					<label class="label" for="email">
						<span class="label-text">Email</span>
					</label>
					<input
						type="email"
						id="email"
						name="email"
						bind:value={email}
						placeholder="your@email.com"
						class="input input-bordered w-full"
						disabled={isLoading}
						required
					/>
				</div>

				<div class="form-control">
					<label class="label" for="username">
						<span class="label-text">Username</span>
					</label>
					<input
						type="text"
						id="username"
						name="username"
						bind:value={username}
						placeholder="username"
						class="input input-bordered w-full"
						disabled={isLoading}
						required
					/>
					<label class="label">
						<span class="label-text-alt">3-50 characters</span>
					</label>
				</div>

				<div class="form-control">
					<label class="label" for="password">
						<span class="label-text">Password</span>
					</label>
					<input
						type="password"
						id="password"
						name="password"
						bind:value={password}
						placeholder="••••••••"
						class="input input-bordered w-full"
						disabled={isLoading}
						required
					/>
					<label class="label">
						<span class="label-text-alt">At least 6 characters</span>
					</label>
				</div>

				<div class="form-control">
					<label class="label" for="confirmPassword">
						<span class="label-text">Confirm Password</span>
					</label>
					<input
						type="password"
						id="confirmPassword"
						name="confirmPassword"
						bind:value={confirmPassword}
						placeholder="••••••••"
						class="input input-bordered w-full"
						disabled={isLoading}
						required
					/>
				</div>

				<div class="form-control">
					<label class="label" for="role">
						<span class="label-text">I am a...</span>
					</label>
					<select
						id="role"
						name="role"
						bind:value={role}
						class="select select-bordered w-full"
						disabled={isLoading}
					>
						<option value={EUserRole.CLIENT}>Client (view photos)</option>
						<option value={EUserRole.PHOTOGRAPHER}>Photographer (upload photos)</option>
					</select>
				</div>

				<div class="form-control mt-6">
					<button type="submit" class="btn btn-primary w-full" disabled={isLoading}>
						{#if isLoading}
							<LoadingSpinner size="sm" />
						{:else}
							Register
						{/if}
					</button>
				</div>
			</form>

			<div class="text-center mt-4">
				<p class="text-sm">
					Already have an account?
					<a href="/login" class="link link-primary">Login</a>
				</p>
			</div>
		</div>
	</div>
</div>
