export enum EUserRole {
	ADMIN = 'admin',
	PHOTOGRAPHER = 'photographer',
	CLIENT = 'client'
}

export type TUser = {
	id: number;
	username: string;
	email: string;
	friendlyName: string;
	role: EUserRole;
	createdAt: string;
	updatedAt: string;
};

export type TLoginRequest = {
	username?: string;
	email?: string;
	password: string;
};

export type TRegisterRequest = {
	email: string;
	username: string;
	password: string;
	role: EUserRole;
};

export type TAuthResponse = {
	user: TUser;
	token: string;
};

export type TRefreshRequest = {
	token: string;
};
