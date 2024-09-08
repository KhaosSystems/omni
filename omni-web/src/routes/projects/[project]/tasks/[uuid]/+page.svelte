<script lang="ts">
    import Textarea from '@khaossystems/matter/components/textarea';
    import Avatar from '@khaossystems/matter/components/avatar';
    import Button from '@khaossystems/matter/components/button';
   
    type User = {
        name: string;
        image: string;
    };

    const users: User[] = [
        { name: 'Storm Søndergaard', image: 'https://upload.wikimedia.org/wikipedia/commons/thumb/c/ca/Original_Utah_Teapot.jpg/220px-Original_Utah_Teapot.jpg' },
        { name: 'Olivia H. Axelsen', image: 'https://upload.wikimedia.org/wikipedia/commons/thumb/c/ca/Original_Utah_Teapot.jpg/220px-Original_Utah_Teapot.jpg' },
    ];
</script>


{#snippet comment(user: User, content: string)}
<div class="flex felx-row flex-grow gap-4 z-20">
    <Avatar class="flex-none" image={user.image} fallback="SS"/>
    <div class="border rounded-lg flex-grow bg-matter-background">
        <header class="flex items-center border-b px-3 h-10">
            <div class="absolute w-2.5 h-2.5 border-t border-l -rotate-45 -translate-x-[17.5px] bg-matter-background"></div>
            <div>{user.name} <span class="text-xs">commented yesterday</span></div>
        </header>
        <main class="px-3 py-2">
            {content}
        </main>
    </div>    
</div>
{/snippet}

{#snippet action(user: User, content: string)}
<div class="flex flex-row pl-[calc(3.75rem+1px)] z-20">
    <div class="w-6 h-6 bg-purple-600 rounded-full mr-1"></div>
    <div><a class="underline" href="/user/123">{user.name}</a> {content}</div>
</div>
{/snippet}


<div class="max-w-8xl mx-auto">
    <div class="text-sm">Alchemist / Why Does Something Not Work</div>
    <div class="text-2xl text-neutral-100">[GAMEPLAY] Character Locomotion</div>
    <div><span class="text-green-600">OPEN</span> - Gameplay Task created by Storm 6/6/2026</div>

    <hr class="bt-4" />


    <div class="flex flex-row gap-8">
        <div class="flex flex-col flex-grow gap-4">
            <!-- History -->
            <div class="flex flex-col h-content gap-4 border-b py-4 relative">
                <div class="absolute border-r-2 h-[calc(100%-1rem)] translate-x-[4.5rem]"></div>
                {@render comment(users[0], "Soo we need the character to move.. and look like (s)he is moving.")}
                {@render comment(users[0], "Well @Olivia can make stuff move!?.")}
                {@render action(users[1], "added Walk Anim to 'blocked by'")}
                {@render comment(users[1], "I made som anim tasks.")}
                {@render action(users[1], "published temp anims")}
                {@render comment(users[1], "Should be in engine now.")}
            </div>

            <div class="flex felx-row flex-grow gap-4 z-20">
                <Avatar class="flex-none" image={users[0].image} fallback="SS"/>
                <div class="flex flex-col items-end flex-grow gap-2">
                    <div class="border rounded-lg flex-grow w-full bg-matter-background">
                        <header class="flex items-center border-b px-3 h-10">
                            <div class="absolute w-2.5 h-2.5 border-t border-l -rotate-45 -translate-x-[17.5px] bg-matter-background"></div>
                            <div>{users[0].name} <span class="text-xs">commented yesterday</span></div>
                        </header>
                        <main class="h-content">
                            <Textarea style="ghost" class="rounded p-4 w-full" placeholder="Add a comment..." />
                        </main>
                    </div>
                    <Button class="w-32">Comment</Button>    
                </div>
            </div>
        </div>
        <div class="w-96 py-4">
            Blocked By<br />
            Walk Anim
        </div>
    </div>
</div>