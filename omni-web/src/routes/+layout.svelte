<script lang="ts">
    import { Button, Avatar, Dialog, Toast } from "@khaossystems/matter";

    import '../app.css'
</script>

<Dialog />
<Toast />

<div class="flex flex-col h-full">
    <header class="flex flex-row border-b items-center px-4 py-1 gap-3">
        <div class="text-2xl font-semibold text-neutral-300">OMNI</div>
        <div class="flex-grow"></div>
        <form method="POST" action="?/logout">
            <Button size="sm" type="submit">Logout</Button>
        </form>
        <Avatar fallback="ss" class="bg-purple-600" />
    </header>
    <div class="flex flex-row flex-grow">
        <slot/>
    </div>
</div>
