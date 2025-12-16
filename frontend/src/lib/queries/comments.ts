import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
import { commentsApi } from '$lib/api/comments';
import { queryKeys } from './keys';
import type { TComment } from '$lib/types';

export interface UseCommentsQueryOptions {
	photoId: number;
	enabled?: boolean;
	refetchInterval?: number | false;
}

export interface CreateCommentData {
	photoId: number;
	text: string;
	parentCommentId?: number | null;
}

export function useCommentsQuery(options: UseCommentsQueryOptions) {
	const { photoId, enabled = true, refetchInterval = 5000 } = options;

	return createQuery({
		queryKey: queryKeys.comments.byPhoto(photoId),
		queryFn: () => commentsApi.getByPhoto(photoId),
		enabled,
		refetchInterval,
		refetchIntervalInBackground: true,
		staleTime: 3000,
	});
}

export function useCreateCommentMutation() {
	const queryClient = useQueryClient();

	return createMutation({
		mutationFn: ({ photoId, text, parentCommentId }: CreateCommentData) =>
			commentsApi.create(photoId, text, parentCommentId),
		onSuccess: (newComment: TComment) => {
			queryClient.invalidateQueries({
				queryKey: queryKeys.comments.byPhoto(newComment.photoId),
			});
			queryClient.invalidateQueries({
				queryKey: queryKeys.photos.detail(newComment.photoId),
			});
		},
	});
}
