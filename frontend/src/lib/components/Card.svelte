<script lang="ts">
	export let title = '';
	export let className = '';
</script>

<div class="card bg-base-100 shadow-xl {className}">
	<div class="card-body">
		{#if title}
			<h2 class="card-title">{title}</h2>
		{/if}
		<slot />
	</div>
</div>
