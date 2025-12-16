export type TSystemSettings = {
	id: number;
	settingKey: string;
	settingValue: string;
	createdAt: string;
	updatedAt: string;
};

export type TUpdateSettingRequest = {
	settingValue: string;
};
