import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
import type { CreateQueryResult, CreateMutationResult } from '@tanstack/svelte-query';
import { albumsApi } from '$lib/api';
import { queryKeys } from './keys';
import type {
	TAlbum,
	TCreateAlbumRequest,
	TUpdateAlbumRequest
} from '$lib/types';

export function useAlbumsQuery() {
	return createQuery({
		queryKey: queryKeys.albums.all,
		queryFn: () => albumsApi.list()
	});
}

export function useAlbumQuery(albumId: number, enabled = true) {
	return createQuery({
		queryKey: queryKeys.albums.detail(albumId),
		queryFn: () => albumsApi.get(albumId),
		enabled
	});
}

export function useAlbumUsersQuery(albumId: number, enabled = true) {
	return createQuery({
		queryKey: queryKeys.albums.users(albumId),
		queryFn: () => albumsApi.getUsers(albumId),
		enabled
	});
}

export function useCreateAlbumMutation() {
	const queryClient = useQueryClient();

	return createMutation({
		mutationFn: (data: TCreateAlbumRequest) => albumsApi.create(data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: queryKeys.albums.all });
		}
	});
}

export function useUpdateAlbumMutation(albumId: number) {
	const queryClient = useQueryClient();

	return createMutation({
		mutationFn: (data: TUpdateAlbumRequest) => albumsApi.update(albumId, data),
		onSuccess: (updatedAlbum) => {
			queryClient.invalidateQueries({ queryKey: queryKeys.albums.all });
			queryClient.invalidateQueries({ queryKey: queryKeys.albums.detail(albumId) });
			queryClient.setQueryData(queryKeys.albums.detail(albumId), updatedAlbum);
		}
	});
}

export function useDeleteAlbumMutation() {
	const queryClient = useQueryClient();

	return createMutation({
		mutationFn: (albumId: number) => albumsApi.delete(albumId),
		onSuccess: (_, albumId) => {
			queryClient.invalidateQueries({ queryKey: queryKeys.albums.all });
			queryClient.removeQueries({ queryKey: queryKeys.albums.detail(albumId) });
			queryClient.removeQueries({ queryKey: queryKeys.albums.users(albumId) });
		}
	});
}

export function useAssignUsersMutation(albumId: number) {
	const queryClient = useQueryClient();

	return createMutation({
		mutationFn: (userIds: number[]) => albumsApi.assignUsers(albumId, userIds),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: queryKeys.albums.users(albumId) });
		}
	});
}
