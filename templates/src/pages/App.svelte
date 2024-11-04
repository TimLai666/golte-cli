<script lang="ts">
    export let title: string;
    let count = 0;
</script>

<h1>{title}</h1>
<span>{count}</span>
<button on:click={() => count++}>Increment</button>