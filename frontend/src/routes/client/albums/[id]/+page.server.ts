import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ locals, params }) => {
	if (!locals.user || !locals.token) {
		throw redirect(303, '/login');
	}

	if (locals.user.role !== 'client') {
		throw redirect(303, '/albums');
	}

	return {
		user: locals.user,
		albumId: params.id
	};
};
