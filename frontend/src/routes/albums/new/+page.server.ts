import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ locals }) => {
	if (!locals.user || !locals.token) {
		throw redirect(303, '/login');
	}

	if (locals.user.role !== 'photographer' && locals.user.role !== 'admin') {
		throw redirect(303, '/');
	}

	return {
		user: locals.user
	};
};
