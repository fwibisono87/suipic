import type { Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
	const token = event.cookies.get('suipic_token');
	const userStr = event.cookies.get('suipic_user');

	if (token && userStr) {
		try {
			const user = JSON.parse(userStr);
			event.locals.user = user;
			event.locals.token = token;
		} catch (e) {
			event.cookies.delete('suipic_token', { path: '/' });
			event.cookies.delete('suipic_user', { path: '/' });
		}
	}

	const response = await resolve(event);
	return response;
};
