import type { TUser } from '$lib/types';

declare global {
	namespace App {
		interface Locals {
			user?: TUser;
			token?: string;
		}
		interface PageData {
			user?: TUser | null;
			isAuthenticated?: boolean;
		}
	}
}

export {};
