<script lang="ts">
    import Badge from "@khaossystems/matter/components/experimental/badge";
    import Button from "@khaossystems/matter/components/button";
    import Card from "@khaossystems/matter/components/card";
    import Input from "@khaossystems/matter/components/input";
    import Select from "@khaossystems/matter/components/select";
    import SquareCheck from "lucide-svelte/icons/square-check";
    import CheveronDown from "lucide-svelte/icons/chevron-down";

    let { data } = $props();
    let tasks = $state(data.tasks);
    let sessionid = $state(data.sessionid);

    let taskInTheMaking: any | null = $state(null);
</script>

<div><a href="/projects">Projects</a> / <a href={`/projects/${data.project.uuid}`}>{data.project.title}</a></div>
<h1 class="text-2xl mb-2">Backlog</h1>

<Card class="bg-neutral-900 full-w">
    <div class="mb-2">
        Board <span class="text-xs">({tasks.length || 0} tasks)</span>
    </div>
    {#each tasks as task, i}
        <div
            class="flex flex-row border-r border-l border-b p-1 items-center px-4 gap-2 bg-neutral-800"
            class:border-t={i == 0}
        >
            <div class="inline">
                <SquareCheck class="inline" size="16" />ALCH-0
            </div>
            <Input
                class="flex-grow"
                type="text"
                inputSize="xs"
                style="ghost"
                bind:value={task.title}
            />
            <Select class="text-xs px-1 py-0 h-fit" options={["TO DO", "IN PROGRESS", "DONE"]} value={task.status} onchange={(event: Event) => {
                console.log(event.target.value)
            }} />
            <div class="text-xs bg-orange-500 font-semibold rounded px-1 text-white">
                {task.status} <CheveronDown class="inline -translate-y-[1px]" size=14 strokeWidth=4 />
            </div>

            <Button size="xs" class="text-xs bg-green-500 hover:bg-green-400 active:bg-green-300 focus:bg-green-400 font-semibold text-white">
                Complete
            </Button>

            <form method="POST" action="?/deleteTask">
                <input type="hidden" name="uuid" value={task.uuid} />
                <Button type="submit" size="xs" class="text-xs bg-red-500 hover:bg-red-400 active:bg-red-300 focus:bg-red-400 font-semibold text-white"  >
                    Delete
                </Button>
            </form>
            
            <a href={`/projects/${data.project}/tasks/${task.uuid}`} class="text-xs">
                View
            </a>
        </div>
    {/each}
    {#if taskInTheMaking}
        <form method="POST" action="?/addTask">
            <input
                type="text"
                name="title"
                bind:value={taskInTheMaking.title}
            />
            <button>Create</button>
        </form>
    {:else}
        <Button
            style="ghost"
            size="xs"
            class="w-full"
            onclick={() => (taskInTheMaking = { title: "", done: false })}
            >+ Create task</Button
        >
    {/if}
</Card>
