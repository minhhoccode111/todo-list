<script lang="ts">
	import { api, getApiError } from '$lib';
	import type { ResponseSession } from '$lib/types/api';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Spinner } from '$lib/components/ui/spinner';
	import { Separator } from '$lib/components/ui/separator';
	import { goto } from '$app/navigation';
	import { clearAuth } from '$lib/stores/auth.svelte';
	import { toast } from 'svelte-sonner';
	import { resolve } from '$app/paths';
	import { SvelteSet } from 'svelte/reactivity';

	let sessions = $state<ResponseSession[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let loggingOut = new SvelteSet<number>();
	let loggingOutAll = $state(false);

	async function fetchSessions() {
		loading = true;
		error = null;
		try {
			const res = await api.sessions.listSessions();
			sessions = res.data;
		} catch (err) {
			error = (await getApiError(err)).message;
		} finally {
			loading = false;
		}
	}

	async function deleteSession(id: number) {
		loggingOut.add(id);
		try {
			await api.sessions.deleteSession(id);
			sessions = sessions.filter((s) => s.id !== id);
			toast.success('Session logged out');
		} catch (err) {
			toast.error((await getApiError(err)).message);
		} finally {
			loggingOut.delete(id);
		}
	}

	async function logoutAll() {
		loggingOutAll = true;
		try {
			await api.logout.logoutAll();
			clearAuth();
			goto(resolve('/login'));
		} catch (err) {
			toast.error((await getApiError(err)).message);
		} finally {
			loggingOutAll = false;
		}
	}

	function formatDate(dateStr: string) {
		return new Date(dateStr).toLocaleString();
	}

	fetchSessions();
</script>

<div class="mx-auto max-w-2xl">
	<Card>
		<CardHeader>
			<CardTitle>Active Sessions</CardTitle>
		</CardHeader>
		<CardContent class="space-y-4">
			{#if loading}
				<div class="flex justify-center p-8">
					<Spinner />
				</div>
			{:else if error}
				<p class="text-destructive text-center">{error}</p>
			{:else if sessions.length === 0}
				<p class="text-muted-foreground text-center">No active sessions</p>
			{:else}
				{#each sessions as session (session.id)}
					<div class="flex items-start justify-between gap-4 rounded-lg border p-4">
						<div class="space-y-1">
							<div class="flex items-center gap-2">
								<span class="font-medium">{session.device || 'Unknown device'}</span>
								{#if session.is_current}
									<Badge variant="default">Current session</Badge>
								{/if}
							</div>
							<p class="text-muted-foreground text-sm">
								Created: {formatDate(session.created_at)}
							</p>
							<p class="text-muted-foreground text-sm">
								Expires: {formatDate(session.expires_at)}
							</p>
						</div>
						<Button
							variant="destructive"
							size="sm"
							disabled={session.is_current || loggingOut.has(session.id)}
							onclick={() => deleteSession(session.id)}
						>
							{loggingOut.has(session.id) ? 'Logging out...' : 'Logout'}
						</Button>
					</div>
				{/each}

				<Separator />

				<Button
					variant="destructive"
					class="w-full"
					disabled={loggingOutAll}
					onclick={logoutAll}
				>
					{loggingOutAll ? 'Logging out all...' : 'Logout All'}
				</Button>
			{/if}
		</CardContent>
	</Card>
</div>
