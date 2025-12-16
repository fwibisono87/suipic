export interface ExifDisplayData {
	make?: string;
	model?: string;
	lens?: string;
	focalLength?: string;
	aperture?: string;
	shutterSpeed?: string;
	iso?: string;
	dateTime?: string;
	dimensions?: string;
	latitude?: number;
	longitude?: number;
}

export async function extractExifData(file: File): Promise<ExifDisplayData | null> {
	try {
		const exifr = await import('exifr');
		const data = await exifr.parse(file, {
			tiff: true,
			exif: true,
			gps: true,
			iptc: false,
			icc: false
		});

		if (!data) return null;

		const exifDisplay: ExifDisplayData = {};

		if (data.Make) exifDisplay.make = data.Make;
		if (data.Model) exifDisplay.model = data.Model;
		if (data.LensModel) exifDisplay.lens = data.LensModel;

		if (data.FocalLength) {
			exifDisplay.focalLength = `${data.FocalLength}mm`;
		}

		if (data.FNumber) {
			exifDisplay.aperture = `f/${data.FNumber}`;
		}

		if (data.ExposureTime) {
			if (data.ExposureTime < 1) {
				exifDisplay.shutterSpeed = `1/${Math.round(1 / data.ExposureTime)}`;
			} else {
				exifDisplay.shutterSpeed = `${data.ExposureTime}s`;
			}
		}

		if (data.ISO) {
			exifDisplay.iso = `ISO ${data.ISO}`;
		}

		if (data.DateTimeOriginal) {
			exifDisplay.dateTime = new Date(data.DateTimeOriginal).toLocaleString();
		}

		if (data.ImageWidth && data.ImageHeight) {
			exifDisplay.dimensions = `${data.ImageWidth} × ${data.ImageHeight}`;
		}

		if (data.latitude !== undefined && data.longitude !== undefined) {
			exifDisplay.latitude = data.latitude;
			exifDisplay.longitude = data.longitude;
		}

		return exifDisplay;
	} catch (err) {
		console.error('EXIF extraction failed:', err);
		return null;
	}
}

export function formatExifValue(key: string, value: unknown): string {
	if (value === null || value === undefined) return '';
	
	switch (key) {
		case 'FocalLength':
			return `${value}mm`;
		case 'FNumber':
			return `f/${value}`;
		case 'ExposureTime':
			if (typeof value === 'number') {
				return value < 1 ? `1/${Math.round(1 / value)}` : `${value}s`;
			}
			return String(value);
		case 'ISO':
			return `ISO ${value}`;
		case 'DateTimeOriginal':
			if (value instanceof Date) {
				return value.toLocaleString();
			}
			return String(value);
		default:
			return String(value);
	}
}

export function formatGPSCoordinates(latitude?: number, longitude?: number): string | null {
	if (latitude === undefined || longitude === undefined) return null;
	
	const latDir = latitude >= 0 ? 'N' : 'S';
	const lonDir = longitude >= 0 ? 'E' : 'W';
	
	return `${Math.abs(latitude).toFixed(6)}° ${latDir}, ${Math.abs(longitude).toFixed(6)}° ${lonDir}`;
}
