import type { Actions } from './$types';
import { fail, redirect } from '@sveltejs/kit';

const API_URL = process.env.VITE_API_URL || 'http://localhost:8080/api';

export const actions = {
	default: async ({ request, cookies }) => {
		const data = await request.formData();
		const username = data.get('username') as string;
		const email = data.get('email') as string;
		const password = data.get('password') as string;

		if (!password) {
			return fail(400, { error: 'Password is required' });
		}

		if (!username && !email) {
			return fail(400, { error: 'Username or email is required' });
		}

		try {
			const response = await fetch(`${API_URL}/auth/login`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					username: username || undefined,
					email: email || undefined,
					password
				})
			});

			if (!response.ok) {
				const error = await response.json().catch(() => ({ message: 'Login failed' }));
				return fail(response.status, { error: error.message || 'Login failed' });
			}

			const result = await response.json();

			cookies.set('suipic_token', result.token, {
				path: '/',
				httpOnly: false,
				sameSite: 'strict',
				maxAge: 60 * 60 * 24 * 7
			});

			cookies.set('suipic_user', JSON.stringify(result.user), {
				path: '/',
				httpOnly: false,
				sameSite: 'strict',
				maxAge: 60 * 60 * 24 * 7
			});

			throw redirect(303, '/');
		} catch (error) {
			if (error instanceof Response) {
				throw error;
			}
			return fail(500, { error: 'An unexpected error occurred' });
		}
	}
} satisfies Actions;
