import { redirect } from '@sveltejs/kit'
import { createResource } from '$lib/khaos/krest.js'
import * as db from '$lib/server/db'


/** @type {import('./$types').PageServerLoad} */
export async function load({ cookies}) {
    const sessionid = cookies.get('sessionid')
    if (!sessionid || sessionid == 'invalidsessionid') {
        return redirect(303, '/login')
    }


    let projects = []
    try {
        projects = await db.getProjects()
    } catch (error) {
        console.error('Error fetching projects:', error)
        projects = []
    }

    return { projects, sessionid }
}

// Docs: https://kit.svelte.dev/docs/form-actions
/** @type {import('./$types').Actions} */
export const actions = {
    createProject: async ({ cookies, request }) => {
        // For reference, see: https://kit.svelte.dev/docs/form-actions#anatomy-of-an-action
        const data = await request.formData()
        const name = data.get('name')
        
        const project = { name }

        try {
            await createResource('v1/projects', project)
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