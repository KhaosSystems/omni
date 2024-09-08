<script lang="ts">
    import { Badge, Button, Card, Input } from "@khaossystems/matter";
    import SquareCheck from "lucide-svelte/icons/square-check";

    let { data } = $props();
    let tasks = $state(data.tasks);
    let sessionid = $state(data.sessionid);

    let taskInTheMaking: any | null = $state(null);
</script>

<h1 class="text-2xl mb-2">Backlog</h1>

<Card class="bg-neutral-900 full-w">
    <div class="mb-2">
        Board <span class="text-xs">({tasks.length || 0} tasks)</span>
    </div>
    {#each tasks as task, i}
        <div
            class="flex flex-row border-r border-l border-b p-1 items-center px-4 gap-2"
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
            <Badge>EDITOR TOOLS</Badge>
            <Badge>TO DO</Badge>
            <Badge>P1</Badge>
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
