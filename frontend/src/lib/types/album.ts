export type TAlbum = {
	id: number;
	title: string;
	description: string | null;
	location: string | null;
	dateTaken: string | null;
	customFields: Record<string, unknown> | null;
	thumbnailPhotoId: number | null;
	photographerId: number;
	createdAt: string;
	updatedAt: string;
};

export type TCreateAlbumRequest = {
	title: string;
	description: string | null;
	location: string | null;
	dateTaken: string | null;
	customFields?: Record<string, unknown> | null;
};

export type TUpdateAlbumRequest = {
	title?: string;
	description?: string | null;
	location?: string | null;
	dateTaken?: string | null;
	customFields?: Record<string, unknown> | null;
	thumbnailPhotoId?: number | null;
};

export type TPhoto = {
	id: number;
	albumId: number;
	filename: string;
	title: string | null;
	description: string | null;
	dateTime: string | null;
	exifData: Record<string, unknown> | null;
	stars: number;
	pickRejectState: string | null;
	createdAt: string;
	updatedAt: string;
};
