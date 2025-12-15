<script lang="ts">
	import { goto } from '$app/navigation';
	import { authStore, isAuthenticated } from '$lib/stores';
	import { authApi } from '$lib/api';
	import { Alert, LoadingSpinner } from '$lib/components';
	import { validateEmail, validateRequired } from '$lib/utils';
	import type { TLoginRequest } from '$lib/types';
	import { onMount } from 'svelte';

	let username = '';
	let email = '';
	let password = '';
	let isLoading = false;
	let error = '';
	let useEmail = false;

	onMount(() => {
		if ($isAuthenticated) {
			goto('/');
		}
	});

	const handleSubmit = async (e: Event) => {
		e.preventDefault();
		error = '';

		if (!password || !validateRequired(password)) {
			error = 'Password is required';
			return;
		}

		if (useEmail) {
			if (!email || !validateEmail(email)) {
				error = 'Valid email is required';
				return;
			}
		} else {
			if (!username || !validateRequired(username)) {
				error = 'Username is required';
				return;
			}
		}

		isLoading = true;

		try {
			const loginData: TLoginRequest = {
				password
			};

			if (useEmail) {
				loginData.email = email;
			} else {
				loginData.username = username;
			}

			const response = await authApi.login(loginData);
			authStore.setAuth(response.user, response.token);
			goto('/');
		} catch (err: unknown) {
			error = (err as { message: string }).message || 'Login failed';
		} finally {
			isLoading = false;
		}
	};

	const toggleLoginMethod = () => {
		useEmail = !useEmail;
		error = '';
	};
</script>

<svelte:head>
	<title>Login - Suipic</title>
</svelte:head>

<div class="flex items-center justify-center min-h-[calc(100vh-300px)]">
	<div class="card w-full max-w-md bg-base-100 shadow-xl">
		<div class="card-body">
			<h2 class="card-title text-3xl font-bold text-center justify-center mb-4">Login</h2>

			{#if error}
				<Alert type="error" message={error} dismissible onDismiss={() => (error = '')} />
			{/if}

			<form on:submit={handleSubmit} class="space-y-4">
				{#if useEmail}
					<div class="form-control">
						<label class="label" for="email">
							<span class="label-text">Email</span>
						</label>
						<input
							type="email"
							id="email"
							bind:value={email}
							placeholder="your@email.com"
							class="input input-bordered w-full"
							disabled={isLoading}
							required
						/>
					</div>
				{:else}
					<div class="form-control">
						<label class="label" for="username">
							<span class="label-text">Username</span>
						</label>
						<input
							type="text"
							id="username"
							bind:value={username}
							placeholder="username"
							class="input input-bordered w-full"
							disabled={isLoading}
							required
						/>
					</div>
				{/if}

				<div class="form-control">
					<label class="label" for="password">
						<span class="label-text">Password</span>
					</label>
					<input
						type="password"
						id="password"
						bind:value={password}
						placeholder="••••••••"
						class="input input-bordered w-full"
						disabled={isLoading}
						required
					/>
					<label class="label">
						<a href="/forgot-password" class="label-text-alt link link-hover">Forgot password?</a>
					</label>
				</div>

				<div class="form-control mt-6">
					<button type="submit" class="btn btn-primary w-full" disabled={isLoading}>
						{#if isLoading}
							<LoadingSpinner size="sm" />
						{:else}
							Login
						{/if}
					</button>
				</div>
			</form>

			<div class="divider">OR</div>

			<button type="button" class="btn btn-outline w-full" on:click={toggleLoginMethod}>
				{useEmail ? 'Login with Username' : 'Login with Email'}
			</button>

			<div class="text-center mt-4">
				<p class="text-sm">
					Don't have an account?
					<a href="/register" class="link link-primary">Register</a>
				</p>
			</div>
		</div>
	</div>
</div>
