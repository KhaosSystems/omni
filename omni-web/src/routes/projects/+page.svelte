<script lang="ts">
    import {
        Badge,
        Button,
        Card,
        Input,
        ScrollArea,
        Dialog,
    } from "@khaossystems/matter";

    let { data } = $props();
    let projects = $state(data.projects);
</script>

<main class="p-5">
    <div class="flex flex-row w-full gap-4">
        {#if projects.length == 0}
            <div>No projects</div>
        {:else}
            {#each projects as project}
                <Card class="w-52 h-36">
                    <div>{JSON.stringify(project)}</div>
                </Card>
            {/each}
        {/if}
    </div>

    <form method="POST" action="?/createProject">
        <input type="text" name="title" />
        <button>Create</button>
    </form>
</main>
