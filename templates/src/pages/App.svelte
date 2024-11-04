<script>
    export let title
    let count = 0
</script>

<div class="container">
    <h1>{title}</h1>
    <span>{count}</span>
    <button on:click={() => count++}>Increment</button>
</div>

<style>
    .container {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        height: 100vh;
    }

    h1 {
        font-size: 2rem;
    }

    button {
        background-color: #007bff;
        color: white;
        border: none;
        border-radius: 0.5rem;
        padding: 0.5rem 1rem;
        cursor: pointer;
        font-size: 1rem;
    }

    button:hover {
        background-color: #0056b3;
    }

    span {
        font-size: 1.5rem;
    }
</style>
