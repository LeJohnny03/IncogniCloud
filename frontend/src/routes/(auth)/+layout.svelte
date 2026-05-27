<script lang="ts">
    // Die 'children' Prop enthält die jeweilige Unterseite (hier: +page.svelte vom login)
    let { children } = $props();
</script>

{@render children()}