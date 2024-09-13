<script lang="ts">
    import Card from "@khaossystems/matter/components/card";

    let { data } = $props();
    let projects = $state(data.projects);
</script>

<main class="p-5">
    <h2 class="mb-2 text-lg">Recent Projects</h2>
    <div class="flex flex-row w-full gap-4">
        {#if projects.length == 0}
            <div>No projects</div>
        {:else}
            {#each projects as project}
                <a href={`/projects/${project.uuid}`}>
                    <Card class="w-52 h-36">
                        <div class="text-lg">{project.name || project.title}</div>
                        <div class="text-xs">{project.uuid}</div>
                    </Card>
                </a>
            {/each}
        {/if}
    </div>

    <form method="POST" action="?/createProject">
        <input type="text" name="name" />
        <button>Create</button>
    </form>
</main>
