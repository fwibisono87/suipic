export * from './auth';
export * from './album';

export type TComment = {
	id: number;
	photoId: number;
	userId: number;
	parentCommentId?: number | null;
	text: string;
	createdAt: string;
	updatedAt: string;
	user: {
		id: number;
		username: string;
		email: string;
		friendlyName: string;
		role: string;
		createdAt: string;
		updatedAt: string;
	};
	replies?: TComment[];
};
