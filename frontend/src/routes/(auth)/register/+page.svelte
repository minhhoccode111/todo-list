<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle
	} from '$lib/components/ui/card';
	import { api, getApiError } from '$lib';
	import { goto } from '$app/navigation';
	import { setAuth } from '$lib/stores/auth.svelte';
	import { toast } from 'svelte-sonner';
	import { resolve } from '$app/paths';

	const minlength = 1;
	const maxlength = 255;

	let name = $state('');
	let email = $state('');
	let password = $state('');
	let errors = $state<{ name?: string; email?: string; password?: string }>({});
	let loading = $state(false);

	const validate = (): boolean => {
		errors = {};

		if (!name.trim()) {
			errors.name = 'Name is required';
		} else if (name.length > 255) {
			errors.name = 'Name must be 255 characters or less';
		} else if (!/^[\p{L}0-9]+$/u.test(name)) {
			errors.name = 'Name can only contain letters and digits';
		}

		if (!email.trim()) {
			errors.email = 'Email is required';
		} else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
			errors.email = 'Invalid email format';
		} else if (email.length > 255) {
			errors.email = 'Email must be 255 characters or less';
		}

		if (!password) {
			errors.password = 'Password is required';
		} else if (password.length > 255) {
			errors.password = 'Password must be 255 characters or less';
		} else if (!/[A-Z]/.test(password)) {
			errors.password = 'Password must contain at least one uppercase letter';
		} else if (!/[a-z]/.test(password)) {
			errors.password = 'Password must contain at least one lowercase letter';
		} else if (!/\d/.test(password)) {
			errors.password = 'Password must contain at least one digit';
		} else if (!/[!@#~$%^&*()+|_{}<>?,./-]/.test(password)) {
			errors.password = 'Password must contain at least one special character';
		}

		return Object.keys(errors).length === 0;
	};

	const handleSubmit = async (e: Event) => {
		e.preventDefault();

		if (!validate()) return;

		loading = true;

		try {
			const res = await api.register.register({ name, email, password });
			setAuth(res.data.token);
			goto(resolve('/'));
		} catch (e) {
			const err = await getApiError(e);
			toast.error(err.message || 'Register failed');
		} finally {
			loading = false;
		}
	};
</script>

<div>
	<Card class="mx-auto w-95">
		<CardHeader>
			<CardTitle>Register</CardTitle>
			<CardDescription>Enter your credentials to sign up.</CardDescription>
		</CardHeader>
		<CardContent>
			<form onsubmit={handleSubmit} class="space-y-4">
				<div class="space-y-2">
					<Label for="name">Name</Label>
					<Input
						id="name"
						type="text"
						bind:value={name}
						placeholder="John"
						{minlength}
						{maxlength}
					/>
					{#if errors.name}
						<p class="text-sm text-destructive">{errors.name}</p>
					{/if}
				</div>

				<div class="space-y-2">
					<Label for="email">Email</Label>
					<Input
						id="email"
						type="email"
						bind:value={email}
						placeholder="you@example.com"
						{minlength}
						{maxlength}
					/>
					{#if errors.email}
						<p class="text-sm text-destructive">{errors.email}</p>
					{/if}
				</div>

				<div class="space-y-2">
					<Label for="password">Password</Label>
					<Input
						id="password"
						type="password"
						bind:value={password}
						placeholder="password"
						{minlength}
						{maxlength}
					/>
					{#if errors.password}
						<p class="text-sm text-destructive">{errors.password}</p>
					{/if}
				</div>

				<Button disabled={loading} type="submit" class="w-full">
					{loading ? 'Signing Up...' : 'Sign Up'}
				</Button>

				<p>
					Already have an account? <a class="text-sky-500" href={resolve('/login')}>Sign in</a>.
				</p>
			</form>
		</CardContent>
	</Card>
</div>
