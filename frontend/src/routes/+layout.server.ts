import type { LayoutServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

const publicRoutes = ['/login', '/register', '/about', '/contact', '/privacy', '/terms'];
const adminRoutes = ['/admin'];
const photographerRoutes = ['/albums/new'];

export const load: LayoutServerLoad = async ({ cookies, url }) => {
	const token = cookies.get('suipic_token');
	const userStr = cookies.get('suipic_user');
	
	let user = null;
	if (userStr) {
		try {
			user = JSON.parse(userStr);
		} catch (e) {
			cookies.delete('suipic_token', { path: '/' });
			cookies.delete('suipic_user', { path: '/' });
		}
	}

	const isAuthenticated = !!(token && user);
	const pathname = url.pathname;

	if (!isAuthenticated && !publicRoutes.includes(pathname) && pathname !== '/') {
		throw redirect(303, '/login');
	}

	if (isAuthenticated) {
		if (adminRoutes.some(route => pathname.startsWith(route)) && user.role !== 'admin') {
			throw redirect(303, '/');
		}

		if (photographerRoutes.some(route => pathname.startsWith(route)) && 
			user.role !== 'photographer' && user.role !== 'admin') {
			throw redirect(303, '/');
		}
	}

	return {
		user,
		isAuthenticated
	};
};
