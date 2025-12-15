import { goto } from '$app/navigation';
import { get } from 'svelte/store';
import { isAuthenticated, currentUser } from '$lib/stores';
import { EUserRole } from '$lib/types';

export const requireAuth = () => {
	const authenticated = get(isAuthenticated);
	if (!authenticated) {
		goto('/login');
		return false;
	}
	return true;
};

export const requireRole = (role: EUserRole) => {
	if (!requireAuth()) return false;
	
	const user = get(currentUser);
	if (!user || user.role !== role) {
		goto('/');
		return false;
	}
	return true;
};

export const requireRoles = (roles: EUserRole[]) => {
	if (!requireAuth()) return false;
	
	const user = get(currentUser);
	if (!user || !roles.includes(user.role)) {
		goto('/');
		return false;
	}
	return true;
};
