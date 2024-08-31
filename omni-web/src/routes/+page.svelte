

<script lang="ts">
    let { data } = $props();
    let tasks = $state(data.tasks)
    let sessionid = $state(data.sessionid);

    let taskInTheMaking: any | null = $state(null);
</script>

sessionid: {sessionid}

<form method="POST" action="?/logout">
    <button>Logout</button>
</form>

<ol>
    {#each tasks as task}
        <li>
            <input type="checkbox" bind:checked={task.done}>
            <input type="text" bind:value={task.title}>
        </li>
    {/each}
    <li>
        {#if taskInTheMaking}
            <form method="POST" action="?/addTask">
                <input type="text" name="title" bind:value={taskInTheMaking.title}>
                <input type="text" name="tidescription	tle" bind:value={taskInTheMaking.description	}>
                <button>Create</button>
            </form>
        {:else}
            <button onclick={() => taskInTheMaking = { title: '', done: false }}>Add task</button>
        {/if}
    </li>
</ol>