import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
import { get } from 'svelte/store';
import { authToken } from '$lib/stores';
import { photographerApi, type TCreateClientRequest } from '$lib/api/photographer';
import { queryKeys } from './keys';

export function useListClients() {
	const token = get(authToken);

	return createQuery({
		queryKey: queryKeys.clients.all,
		queryFn: () => photographerApi.listClients(token || ''),
		enabled: !!token
	});
}

export function useSearchClients(query: string) {
	const token = get(authToken);

	return createQuery({
		queryKey: queryKeys.clients.search(query),
		queryFn: () => photographerApi.searchClients(query, token || ''),
		enabled: !!token && query.length > 0
	});
}

export function useCreateOrLinkClient() {
	const queryClient = useQueryClient();
	const token = get(authToken);

	return createMutation({
		mutationFn: (data: TCreateClientRequest) => 
			photographerApi.createOrLinkClient(data, token || ''),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: queryKeys.clients.all });
		}
	});
}
