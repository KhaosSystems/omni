<script lang="ts">
    import { Button, Input } from '@khaossystems/matter'

    let { data } = $props();
    let sessionid = $state(data.sessionid);
</script>

<main class="flex items-center justify-center h-screen">
    <form class="flex flex-col w-96 border rounded p-8 gap-4" method="POST" action="?/login">
        <Input type="text" name="username" placeholder="Username" />
        <Input type="password" name="password" placeholder="Password" />
        <Button type="submit">Login</Button>
        sessionid: {sessionid}
    </form>
</main>

