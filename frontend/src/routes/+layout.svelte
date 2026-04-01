<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { ModeWatcher } from 'mode-watcher';
	import { resolve } from '$app/paths';
	import { Toaster } from 'svelte-sonner';

	let { children } = $props();
	import { Button } from '$lib/components/ui/button';
	import { ToggleTheme } from '$lib/components/ui/toggle-theme';

	import { auth } from '$lib/stores/auth.svelte';
</script>

<svelte:head>
	<link rel="icon" href={favicon} />

	<title>Todo-list</title>
</svelte:head>

<ModeWatcher />

<Toaster closeButton richColors />
<!-- usage: toast('something'), toast.success('something'), toast.error('something') -->

<div class="mx-auto flex min-h-screen max-w-6xl flex-col gap-4">
	<header class="flex items-center justify-between gap-4 p-4">
		<h1 class="text-4xl font-bold">Todo-list</h1>

		<nav class="flex list-none justify-evenly gap-4">
			<li>
				<Button disabled={!auth.token} href={resolve('/')}>home</Button>
			</li>

			{#if !auth.token}
				<li>
					<Button href={resolve('/login')}>login</Button>
				</li>
				<li>
					<Button href={resolve('/register')}>register</Button>
				</li>
			{:else}
				<li>
					<Button href={resolve('/logout')}>logout</Button>
				</li>
			{/if}

			<li>
				<ToggleTheme />
			</li>
		</nav>
	</header>

	<main class="flex-1 p-4">
		{@render children()}
	</main>

	<footer class="p-4">
		<center>
			<p>
				Made by
				<a
					class="text-sky-500 transition-colors hover:text-sky-600"
					href="https://github.com/minhhoccode111">minhhoccode111</a
				>.
			</p>
		</center>
	</footer>
</div>
