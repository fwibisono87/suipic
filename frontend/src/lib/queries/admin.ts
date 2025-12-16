import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
import { get } from 'svelte/store';
import { authToken } from '$lib/stores';
import { adminApi, type TCreatePhotographerRequest } from '$lib/api/admin';
import { queryKeys } from './keys';

export function useListPhotographers() {
	const token = get(authToken);

	return createQuery({
		queryKey: queryKeys.photographers.all,
		queryFn: () => adminApi.listPhotographers(token || ''),
		enabled: !!token
	});
}

export function useCreatePhotographer() {
	const queryClient = useQueryClient();
	const token = get(authToken);

	return createMutation({
		mutationFn: (data: TCreatePhotographerRequest) => 
			adminApi.createPhotographer(data, token || ''),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: queryKeys.photographers.all });
		}
	});
}
