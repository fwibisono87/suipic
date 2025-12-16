<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import Icon from '@iconify/svelte';
	import type { TComment, TPhoto } from '$lib/types';
	import { commentsApi } from '$lib/api';
	import { currentUser } from '$lib/stores';
	import { formatRelativeTime } from '$lib/utils/format';
	import { LoadingSpinner } from '$lib/components';

	export let photo: TPhoto;
	export let photographerId: number;

	let comments: TComment[] = [];
	let isLoading = true;
	let error = '';
	let commentText = '';
	let isSubmitting = false;
	let replyToCommentId: number | null = null;
	let replyToUsername: string = '';
	let pollInterval: number | null = null;
	let showPreview = false;

	async function loadComments() {
		try {
			comments = await commentsApi.getByPhoto(photo.id);
			error = '';
		} catch (err) {
			error = (err as Error).message;
		} finally {
			isLoading = false;
		}
	}

	async function handleSubmit() {
		if (!commentText.trim() || isSubmitting) return;

		isSubmitting = true;
		error = '';

		try {
			await commentsApi.create(photo.id, commentText.trim(), replyToCommentId);
			commentText = '';
			replyToCommentId = null;
			replyToUsername = '';
			showPreview = false;
			await loadComments();
		} catch (err) {
			error = (err as Error).message;
		} finally {
			isSubmitting = false;
		}
	}

	function handleReply(comment: TComment) {
		replyToCommentId = comment.id;
		replyToUsername = comment.user.friendlyName || comment.user.username;
		const textarea = document.getElementById('comment-textarea');
		if (textarea) {
			textarea.focus();
		}
	}

	function cancelReply() {
		replyToCommentId = null;
		replyToUsername = '';
	}

	function getUserBadge(comment: TComment): { label: string; class: string } | null {
		if (comment.userId === photographerId) {
			return { label: 'Photographer', class: 'badge-primary' };
		}
		if (comment.user.role === 'client') {
			return { label: 'Client', class: 'badge-secondary' };
		}
		if (comment.user.role === 'admin') {
			return { label: 'Admin', class: 'badge-accent' };
		}
		return null;
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Enter' && (event.ctrlKey || event.metaKey)) {
			event.preventDefault();
			handleSubmit();
		}
	}

	onMount(() => {
		loadComments();
		pollInterval = window.setInterval(() => {
			loadComments(true);
		}, 10000);
	});

	onDestroy(() => {
		if (pollInterval) {
			clearInterval(pollInterval);
		}
	});

	$: previewText = commentText.trim();
	
	$: totalComments = comments.reduce((total, comment) => {
		return total + 1 + (comment.replies?.length || 0);
	}, 0);
</script>

<div class="comment-section space-y-6">
	<div class="flex items-center justify-between">
		<h3 class="text-xl font-semibold flex items-center gap-2">
			<Icon icon="mdi:comment-multiple" class="text-2xl" />
			Comments
			{#if totalComments > 0}
				<span class="badge badge-neutral">
					{totalComments} {totalComments === 1 ? 'comment' : 'comments'}
				</span>
			{/if}
			{#if isRefreshing}
				<span class="loading loading-spinner loading-xs opacity-50"></span>
			{/if}
		</h3>
	</div>

	{#if error && !isLoading}
		<div class="alert alert-error">
			<Icon icon="mdi:alert-circle" class="text-xl" />
			<span>{error}</span>
		</div>
	{/if}

	<div class="comment-form">
		<div class="space-y-2">
			{#if replyToCommentId}
				<div class="flex items-center gap-2 text-sm bg-base-200 px-3 py-2 rounded-lg">
					<Icon icon="mdi:reply" class="text-lg" />
					<span>Replying to <strong>{replyToUsername}</strong></span>
					<button
						class="btn btn-ghost btn-xs ml-auto"
						on:click={cancelReply}
						type="button"
					>
						<Icon icon="mdi:close" />
					</button>
				</div>
			{/if}

			<div class="form-control">
				<label class="label" for="comment-textarea">
					<span class="label-text font-medium">Add a comment</span>
					<button
						class="btn btn-ghost btn-xs"
						on:click={() => (showPreview = !showPreview)}
						type="button"
					>
						<Icon icon={showPreview ? 'mdi:eye-off' : 'mdi:eye'} class="text-lg" />
						{showPreview ? 'Hide' : 'Show'} Preview
					</button>
				</label>
				<textarea
					id="comment-textarea"
					class="textarea textarea-bordered h-24 resize-none"
					placeholder={replyToCommentId ? `Reply to ${replyToUsername}...` : 'Share your thoughts...'}
					bind:value={commentText}
					on:keydown={handleKeydown}
					disabled={isSubmitting}
				></textarea>
				<label class="label">
					<span class="label-text-alt text-xs opacity-60">
						Ctrl+Enter or Cmd+Enter to submit
					</span>
					<span class="label-text-alt text-xs opacity-60">
						{commentText.length} characters
					</span>
				</label>
			</div>

			{#if showPreview && previewText}
				<div class="preview bg-base-200 rounded-lg p-4">
					<div class="text-xs font-semibold opacity-60 mb-2">PREVIEW</div>
					<div class="prose prose-sm max-w-none">
						<p class="whitespace-pre-wrap">{previewText}</p>
					</div>
				</div>
			{/if}

			<div class="flex justify-end">
				<button
					class="btn btn-primary"
					disabled={!commentText.trim() || isSubmitting}
					on:click={handleSubmit}
				>
					{#if isSubmitting}
						<LoadingSpinner size="sm" />
					{:else}
						<Icon icon="mdi:send" class="text-xl" />
					{/if}
					{replyToCommentId ? 'Reply' : 'Comment'}
				</button>
			</div>
		</div>
	</div>

	<div class="divider"></div>

	<div class="comments-list space-y-4">
		{#if isLoading}
			<div class="flex justify-center py-8">
				<LoadingSpinner size="lg" />
			</div>
		{:else if comments.length === 0}
			<div class="text-center py-12 opacity-60">
				<Icon icon="mdi:comment-outline" class="text-6xl mx-auto mb-4" />
				<p class="text-lg">No comments yet</p>
				<p class="text-sm mt-1">Be the first to share your thoughts!</p>
			</div>
		{:else}
			{#each comments as comment (comment.id)}
				<div class="comment-thread">
					<div class="comment-item flex gap-3">
						<div class="avatar placeholder flex-shrink-0">
							<div class="bg-neutral text-neutral-content rounded-full w-10 h-10">
								<span class="text-sm">
									{(comment.user.friendlyName || comment.user.username).charAt(0).toUpperCase()}
								</span>
							</div>
						</div>

						<div class="flex-1 space-y-1">
							<div class="flex items-center gap-2 flex-wrap">
								<span class="font-semibold">
									{comment.user.friendlyName || comment.user.username}
								</span>
								{#if getUserBadge(comment)}
									<span class="badge badge-sm {getUserBadge(comment)?.class}">
										{getUserBadge(comment)?.label}
									</span>
								{/if}
								<span class="text-xs opacity-60">
									{formatRelativeTime(comment.createdAt)}
								</span>
							</div>

							<div class="prose prose-sm max-w-none">
								<p class="whitespace-pre-wrap m-0">{comment.text}</p>
							</div>

							<div class="flex gap-2 mt-2">
								<button
									class="btn btn-ghost btn-xs"
									on:click={() => handleReply(comment)}
								>
									<Icon icon="mdi:reply" class="text-base" />
									Reply
								</button>
							</div>

							{#if comment.replies && comment.replies.length > 0}
								<div class="replies space-y-3 mt-4 pl-4 border-l-2 border-base-300">
									{#each comment.replies as reply (reply.id)}
										<div class="reply-item flex gap-3">
											<div class="avatar placeholder flex-shrink-0">
												<div class="bg-neutral text-neutral-content rounded-full w-8 h-8">
													<span class="text-xs">
														{(reply.user.friendlyName || reply.user.username).charAt(0).toUpperCase()}
													</span>
												</div>
											</div>

											<div class="flex-1 space-y-1">
												<div class="flex items-center gap-2 flex-wrap">
													<span class="font-semibold text-sm">
														{reply.user.friendlyName || reply.user.username}
													</span>
													{#if getUserBadge(reply)}
														<span class="badge badge-xs {getUserBadge(reply)?.class}">
															{getUserBadge(reply)?.label}
														</span>
													{/if}
													<span class="text-xs opacity-60">
														{formatRelativeTime(reply.createdAt)}
													</span>
												</div>

												<div class="prose prose-sm max-w-none">
													<p class="whitespace-pre-wrap m-0 text-sm">{reply.text}</p>
												</div>

												<div class="flex gap-2 mt-1">
													<button
														class="btn btn-ghost btn-xs"
														on:click={() => handleReply(comment)}
													>
														<Icon icon="mdi:reply" class="text-base" />
														Reply
													</button>
												</div>
											</div>
										</div>
									{/each}
								</div>
							{/if}
						</div>
					</div>
				</div>
			{/each}
		{/if}
	</div>
</div>

<style>
	.comment-section {
		width: 100%;
	}

	.comment-item,
	.reply-item {
		animation: fadeIn 0.3s ease-out;
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
			transform: translateY(-10px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	.preview {
		border: 1px dashed currentColor;
		opacity: 0.8;
	}

	.replies {
		animation: slideDown 0.3s ease-out;
	}

	@keyframes slideDown {
		from {
			opacity: 0;
			max-height: 0;
		}
		to {
			opacity: 1;
			max-height: 1000px;
		}
	}
</style>
