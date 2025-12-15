import type { Actions } from './$types';
import { fail, redirect } from '@sveltejs/kit';

const API_URL = process.env.VITE_API_URL || 'http://localhost:8080/api';

export const actions = {
	default: async ({ request, cookies }) => {
		const data = await request.formData();
		const email = data.get('email') as string;
		const username = data.get('username') as string;
		const password = data.get('password') as string;
		const role = data.get('role') as string;

		if (!email || !username || !password) {
			return fail(400, { error: 'Email, username, and password are required' });
		}

		try {
			const response = await fetch(`${API_URL}/auth/register`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					email,
					username,
					password,
					role: role || 'client'
				})
			});

			if (!response.ok) {
				const error = await response.json().catch(() => ({ message: 'Registration failed' }));
				return fail(response.status, { error: error.message || 'Registration failed' });
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
