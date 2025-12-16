export const formatDateTime = (dateString: string): string => {
	const date = new Date(dateString);
	return new Intl.DateTimeFormat('en-US', {
		year: 'numeric',
		month: 'long',
		day: 'numeric',
		hour: '2-digit',
		minute: '2-digit'
	}).format(date);
};

export const formatDate = (dateString: string): string => {
	const date = new Date(dateString);
	return new Intl.DateTimeFormat('en-US', {
		year: 'numeric',
		month: 'long',
		day: 'numeric'
	}).format(date);
};

export const formatRelativeTime = (dateString: string): string => {
	const date = new Date(dateString);
	const now = new Date();
	const diffInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000);

	if (diffInSeconds < 60) return 'just now';
	
	const minutes = Math.floor(diffInSeconds / 60);
	if (diffInSeconds < 3600) {
		return `${minutes} ${minutes === 1 ? 'minute' : 'minutes'} ago`;
	}
	
	const hours = Math.floor(diffInSeconds / 3600);
	if (diffInSeconds < 86400) {
		return `${hours} ${hours === 1 ? 'hour' : 'hours'} ago`;
	}
	
	const days = Math.floor(diffInSeconds / 86400);
	if (diffInSeconds < 604800) {
		return `${days} ${days === 1 ? 'day' : 'days'} ago`;
	}
	
	return formatDate(dateString);
};
