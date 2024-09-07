import { redirect } from '@sveltejs/kit'
import { createResource } from '$lib/khaos/krest.js'
import * as db from '$lib/server/db'


/** @type {import('./$types').PageServerLoad} */
export async function load({ cookies }) {
    const sessionid = cookies.get('sessionid')
    if (!sessionid || sessionid == 'invalidsessionid') {
        return redirect(303, '/login')
    }

    let tasks = []
    try {
        tasks = await db.getTasks()
    } catch (error) {
        console.error('Error fetching tasks:', error)
        tasks = []
    }

    return { tasks, sessionid }
}

// Docs: https://kit.svelte.dev/docs/form-actions
/** @type {import('./$types').Actions} */
export const actions = {
    addTask: async ({ cookies, request }) => {
        // For reference, see: https://kit.svelte.dev/docs/form-actions#anatomy-of-an-action
        const data = await request.formData()
        const title = data.get('title')
        const description = data.get('description')
        
        const task = { title, description }

        try {
            await createResource('v1/tasks', task)
        } catch (error) {
            console.error(error)
            return { status: 500 }
        }
    },
    logout: async ({ cookies }) => {
        console.log('Logging out...')
        cookies.set('sessionid', 'invalidsessionid', { path: '/' })
        redirect(303, '/login')
    }
}