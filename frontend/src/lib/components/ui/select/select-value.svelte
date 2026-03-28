<script lang="ts">
	type Props = {
		placeholder?: string;
		value?: string;
		class?: string;
	};

	let { placeholder = '', value = '', class: className }: Props = $props();
</script>

<span class={className}>{value || placeholder}</span>
