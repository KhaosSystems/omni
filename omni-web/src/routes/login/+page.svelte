<script lang="ts">
    let { data } = $props();
    let sessionid = $state(data.sessionid);
</script>

<form method="POST" action="?/login">
    <input type="text" name="username">
    <input type="password" name="password">
    <button type="submit">Login</button>
</form>


sessionid: {sessionid}
