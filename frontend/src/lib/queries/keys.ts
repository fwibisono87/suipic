import type { TSearchParams } from '$lib/api/search';

export const queryKeys = {
	albums: {
		all: ['albums'] as const,
		detail: (albumId: number) => ['album', albumId] as const,
		users: (albumId: number) => ['albumUsers', albumId] as const,
	},
	photos: {
		all: ['photos'] as const,
		byAlbum: (albumId: number) => ['photos', albumId] as const,
		detail: (photoId: number) => ['photo', photoId] as const,
	},
	comments: {
		all: ['comments'] as const,
		byPhoto: (photoId: number) => ['comments', photoId] as const,
	},
	clients: {
		all: ['clients'] as const,
		search: (query: string) => ['clients', 'search', query] as const,
	},
	photographers: {
		all: ['photographers'] as const,
		detail: (photographerId: number) => ['photographer', photographerId] as const,
	},
	search: {
		all: ['search'] as const,
		query: (params: TSearchParams) => ['search', params] as const,
	},
} as const;
