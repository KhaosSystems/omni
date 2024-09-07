

<script lang="ts">
    import { Avatar, Button } from "@khaossystems/matter";


    let { data } = $props();
    let tasks = $state(data.tasks)
    let sessionid = $state(data.sessionid);

    let taskInTheMaking: any | null = $state(null);
</script>

<header class="flex flex-row border-b items-center px-4 py-2 gap-3" >
    <div class="text-2xl font-semibold text-neutral-300">OMNI</div>
    <div class="flex-grow"></div>
    <form method="POST" action="?/logout">
        <Button size="sm" type="submit">Logout</Button>
    </form>
    <Avatar fallback="ss" class="bg-purple-600"/>
</header>
<main class="p-8">
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
</main>
