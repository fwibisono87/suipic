<script lang="ts">
	export let size: 'xs' | 'sm' | 'md' | 'lg' = 'md';

	const sizeClasses = {
		xs: 'loading-xs',
		sm: 'loading-sm',
		md: 'loading-md',
		lg: 'loading-lg'
	};
</script>

<span class="loading loading-spinner {sizeClasses[size]}"></span>
