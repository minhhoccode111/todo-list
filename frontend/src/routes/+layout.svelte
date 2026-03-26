<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { Toaster } from 'svelte-sonner';

	let { children } = $props();
	import { auth } from '$lib/stores/auth.svelte';
	import { Button } from '$lib/components/ui/button';
	import { goto } from '$app/navigation';

	$effect(() => {
		if (!auth) {
			goto('/login');
		}
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />

	<title>Todo-list</title>
</svelte:head>

<Toaster />
<!-- usage: toast('something'), toast.success('something'), toast.error('something') -->

<div id="wrapper" class="bg-gray-800 text-yellow-100">
	<div class="mx-auto flex min-h-screen max-w-6xl flex-col gap-4 border border-red-500">
		<header class="flex justify-between gap-4 p-4">
			<h1>
				<a href="/">Todo-list</a>
			</h1>

			<nav class="flex list-none justify-evenly gap-4">
				<li>
					<Button href="/">home</Button>
				</li>
				{#if !auth}
					<li>
						<Button href="/login">login</Button>
					</li>
					<li>
						<Button href="/register">register</Button>
					</li>
				{:else}
					<li>
						<Button href="/logout">logout</Button>
					</li>
				{/if}
			</nav>
		</header>

		<main class="flex-1 p-4">
			{@render children()}
		</main>

		<footer class="p-4">
			<center>
				<p>Made by minhhoccode111.</p>
			</center>
		</footer>
	</div>
</div>
