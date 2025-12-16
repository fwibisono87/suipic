import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
import { photosApi } from '$lib/api';
import { queryKeys } from './keys';
import type { TPhoto } from '$lib/types';

export const usePhotosQuery = (albumId: number, enabled = true) => {
	return createQuery({
		queryKey: queryKeys.photos.byAlbum(albumId),
		queryFn: () => photosApi.listByAlbum(albumId),
		enabled
	});
};

export const usePhotoQuery = (photoId: number, enabled = true) => {
	return createQuery({
		queryKey: queryKeys.photos.detail(photoId),
		queryFn: () => photosApi.get(photoId),
		enabled
	});
};

export const useCreatePhotoMutation = (albumId: number) => {
	const queryClient = useQueryClient();

	return createMutation({
		mutationFn: (file: File) => photosApi.create(albumId, file),
		onMutate: async () => {
			await queryClient.cancelQueries({ queryKey: queryKeys.photos.byAlbum(albumId) });
			const previousPhotos = queryClient.getQueryData<TPhoto[]>(queryKeys.photos.byAlbum(albumId));
			return { previousPhotos };
		},
		onSuccess: (newPhoto) => {
			queryClient.setQueryData<TPhoto[]>(
				queryKeys.photos.byAlbum(albumId),
				(old) => old ? [...old, newPhoto] : [newPhoto]
			);
			queryClient.invalidateQueries({ queryKey: queryKeys.photos.byAlbum(albumId) });
		},
		onError: (_err, _newPhoto, context) => {
			if (context?.previousPhotos) {
				queryClient.setQueryData(queryKeys.photos.byAlbum(albumId), context.previousPhotos);
			}
		}
	});
};

export const useUpdatePhotoMutation = () => {
	const queryClient = useQueryClient();

	return createMutation({
		mutationFn: ({ 
			photoId, 
			updates 
		}: { 
			photoId: number; 
			updates: Partial<{ title: string | null; pickRejectState: string | null; stars: number }> 
		}) => photosApi.update(photoId, updates),
		onMutate: async ({ photoId, updates }) => {
			await queryClient.cancelQueries({ queryKey: queryKeys.photos.detail(photoId) });
			const previousPhoto = queryClient.getQueryData<TPhoto>(queryKeys.photos.detail(photoId));
			
			if (previousPhoto) {
				queryClient.setQueryData<TPhoto>(
					queryKeys.photos.detail(photoId),
					{ ...previousPhoto, ...updates }
				);
			}

			const albumQueries = queryClient.getQueriesData<TPhoto[]>({ 
				queryKey: queryKeys.photos.all 
			});
			
			const previousAlbumPhotos: Array<{ queryKey: unknown[]; data: TPhoto[] }> = [];
			
			albumQueries.forEach(([queryKey, photos]) => {
				if (photos) {
					previousAlbumPhotos.push({ queryKey: queryKey as unknown[], data: photos });
					const updatedPhotos = photos.map(p => 
						p.id === photoId ? { ...p, ...updates } : p
					);
					queryClient.setQueryData(queryKey, updatedPhotos);
				}
			});

			return { previousPhoto, previousAlbumPhotos };
		},
		onSuccess: (updatedPhoto) => {
			queryClient.setQueryData(queryKeys.photos.detail(updatedPhoto.id), updatedPhoto);
			
			queryClient.invalidateQueries({ 
				queryKey: queryKeys.photos.byAlbum(updatedPhoto.albumId) 
			});
		},
		onError: (_err, { photoId }, context) => {
			if (context?.previousPhoto) {
				queryClient.setQueryData(queryKeys.photos.detail(photoId), context.previousPhoto);
			}
			
			if (context?.previousAlbumPhotos) {
				context.previousAlbumPhotos.forEach(({ queryKey, data }) => {
					queryClient.setQueryData(queryKey, data);
				});
			}
		}
	});
};

export const useBatchUploadMutation = (albumId: number) => {
	const queryClient = useQueryClient();

	return createMutation({
		mutationFn: (files: File[]) => photosApi.createBatch(albumId, files),
		onMutate: async () => {
			await queryClient.cancelQueries({ queryKey: queryKeys.photos.byAlbum(albumId) });
			const previousPhotos = queryClient.getQueryData<TPhoto[]>(queryKeys.photos.byAlbum(albumId));
			return { previousPhotos };
		},
		onSuccess: (newPhotos) => {
			queryClient.setQueryData<TPhoto[]>(
				queryKeys.photos.byAlbum(albumId),
				(old) => old ? [...old, ...newPhotos] : newPhotos
			);
			queryClient.invalidateQueries({ queryKey: queryKeys.photos.byAlbum(albumId) });
		},
		onError: (_err, _files, context) => {
			if (context?.previousPhotos) {
				queryClient.setQueryData(queryKeys.photos.byAlbum(albumId), context.previousPhotos);
			}
		}
	});
};
