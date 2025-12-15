export type TAlbum = {
	id: number;
	title: string;
	description: string | null;
	location: string | null;
	dateTaken: string | null;
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
};

export type TUpdateAlbumRequest = {
	title?: string;
	description?: string | null;
	location?: string | null;
	dateTaken?: string | null;
	thumbnailPhotoId?: number | null;
};
