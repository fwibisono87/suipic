import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
import { settingsApi } from '$lib/api';
import { queryKeys } from './keys';
import type { TUpdateSettingRequest } from '$lib/types';

export function useSettingsQuery() {
	return createQuery({
		queryKey: queryKeys.settings.all,
		queryFn: () => settingsApi.getSettings()
	});
}

export function useUpdateSettingMutation() {
	const queryClient = useQueryClient();

	return createMutation({
		mutationFn: ({ settingKey, data }: { settingKey: string; data: TUpdateSettingRequest }) =>
			settingsApi.updateSetting(settingKey, data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: queryKeys.settings.all });
		}
	});
}
