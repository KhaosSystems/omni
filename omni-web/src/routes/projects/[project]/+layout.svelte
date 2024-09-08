<script lang="ts">
    import Archive from "lucide-svelte/icons/archive";
    import { page } from '$app/stores';

    let { data } = $props();
</script>

<aside class="w-60 flex-none border-r p-2">
    <div class="flex flex-row gap-2 px-3 py-7">
        <div class="bg-purple-600 w-8 h-8 rounded"></div>
        <div class="flex flex-col">
            <div class="text-xs font-semibold">Project Name</div>
            <div class="text-xs">Khaos Systems</div>
        </div>
    </div>

    <div>
        <div class="text-sm font-semibold px-3">PLANNING</div>
        <ul>
            <li 
                class="flex items-center text-sm px-3 py-2 hover:bg-neutral-900 data-[active=true]:bg-neutral-900 rounded-lg"
                data-active={$page.url.pathname == '/projects/text/dashboard'}
            >
                <a href={`/projects/${data.project}/dashboard`}><Archive size="20" class="inline mr-2" />Dashboard</a>
            </li>
            <li
                class="flex items-center text-sm px-3 py-2 hover:bg-neutral-900 data-[active=true]:bg-neutral-900 rounded-lg"
                data-active={$page.url.pathname == '/projects/text/backlog'}
            >
                <a href={`/projects/${data.project}/backlog`}><Archive size="20" class="inline mr-2" />Backlog</a>
            </li>
        </ul>
    </div>
</aside>

<main class="p-8 w-full">
    <slot />
</main>