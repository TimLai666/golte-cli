<script>
    let count = 0

    // 添加動畫類別
    const addPopAnimation = () => {
        const btn = document.querySelector('.counter-btn')
        btn.classList.add('pop')
        setTimeout(() => btn.classList.remove('pop'), 200)
    }
</script>

<div class="counter-container">
    <span class="counter">{count}</span>
    <button 
        class="counter-btn"
        on:click={() => {
            count++
            addPopAnimation()
        }}
    >
        Click Me!
    </button>
</div>

<style>
    .counter-container {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 1rem;
        padding: 2rem;
        background: white;
        border-radius: 1rem;
        box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
    }

    .counter {
        font-size: 3rem;
        font-weight: bold;
        color: #3498db;
    }

    .counter-btn {
        background-color: #3498db;
        color: white;
        border: none;
        border-radius: 0.5rem;
        padding: 0.8rem 1.5rem;
        cursor: pointer;
        font-size: 1.1rem;
        transition: all 0.3s ease;
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    }

    .counter-btn:hover {
        background-color: #2980b9;
        transform: translateY(-2px);
        box-shadow: 0 6px 8px rgba(0, 0, 0, 0.15);
    }

    .counter-btn:active {
        transform: translateY(0);
    }

    /* 按鈕點擊動畫 */
    .pop {
        animation: pop 0.2s ease-out;
    }

</style>