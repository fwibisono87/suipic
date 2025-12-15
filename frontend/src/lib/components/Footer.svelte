<script lang="ts">
	import Icon from '@iconify/svelte';
</script>

<footer class="footer footer-center p-10 bg-base-200 text-base-content rounded mt-8">
	<nav class="grid grid-flow-col gap-4">
		<a href="/about" class="link link-hover">About</a>
		<a href="/contact" class="link link-hover">Contact</a>
		<a href="/privacy" class="link link-hover">Privacy</a>
		<a href="/terms" class="link link-hover">Terms</a>
	</nav>
	<nav>
		<div class="grid grid-flow-col gap-4">
			<a href="https://twitter.com" target="_blank" rel="noopener noreferrer" class="btn btn-ghost btn-circle">
				<Icon icon="mdi:twitter" class="text-2xl" />
			</a>
			<a href="https://github.com" target="_blank" rel="noopener noreferrer" class="btn btn-ghost btn-circle">
				<Icon icon="mdi:github" class="text-2xl" />
			</a>
			<a href="https://instagram.com" target="_blank" rel="noopener noreferrer" class="btn btn-ghost btn-circle">
				<Icon icon="mdi:instagram" class="text-2xl" />
			</a>
		</div>
	</nav>
	<aside>
		<p>Copyright © {new Date().getFullYear()} - Suipic Photo Manager</p>
	</aside>
</footer>
