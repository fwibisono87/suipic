import type { Actions, PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async () => {
	throw redirect(303, '/login');
};

export const actions = {
	default: async ({ cookies }) => {
		cookies.delete('suipic_token', { path: '/' });
		cookies.delete('suipic_user', { path: '/' });
		throw redirect(303, '/login');
	}
} satisfies Actions;
