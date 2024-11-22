<script lang="ts">
    export let title: string
    import Counter from '../components/Counter.svelte'
</script>

<div class="container">
    <h1 class="title">{title}</h1>
    <Counter />
</div>

<style>
    .container {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        height: 100vh;
        background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
    }

    .title {
        font-size: 2.5rem;
        color: #2c3e50;
        margin-bottom: 2rem;
        text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.1);
        transition: transform 0.3s ease;
    }

    .title:hover {
        transform: scale(1.05);
    }

    @keyframes pop {
        0% { transform: scale(1); }
        50% { transform: scale(1.1); }
        100% { transform: scale(1); }
    }
</style>
