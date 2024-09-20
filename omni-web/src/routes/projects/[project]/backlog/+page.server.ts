import { redirect } from '@sveltejs/kit'
import { getResource, createResource, deleteResource, updateResource } from '$lib/khaos/krest.js'
import * as db from '$lib/server/db'
import * as krest from '$lib/khaos/krest.js'


/** @type {import('./$types').PageServerLoad} */
export async function load({ cookies, params }) {
    const sessionid = cookies.get('sessionid')
    if (!sessionid || sessionid == 'invalidsessionid') {
        return redirect(303, '/login')
    }
 
    let tasks = []
    try {
        tasks = await krest.getCollection("v1/tasks", { expand: ["summary", "project_id"] })
        tasks = tasks.filter((task) => task.project_id == params.project)
    } catch (error) {
        console.error('Error fetching tasks:', error)
        tasks = []
    }

    let project = {}
    try {
        project = await getResource('v1/projects', params.project) 
    } catch (error) {
        console.error('Error fetching project:', error)
    }

    return { tasks, sessionid, project: project }
}

// Docs: https://kit.svelte.dev/docs/form-actions
/** @type {import('./$types').Actions} */
export const actions = {
    addTask: async ({ cookies, request }) => {
        // For reference, see: https://kit.svelte.dev/docs/form-actions#anatomy-of-an-action
        const data = await request.formData()
        const summary = data.get('summary')
        const description = data.get('description')
        const project_id = data.get('project_id')
        console.log(project_id)
        const task = { summary, description, project_id }

        try {
            await createResource('v1/tasks', task)
        } catch (error) {
            console.error(error)
            return { status: 500 }
        }
    },
    deleteTask: async ({ cookies, request }) => {
        const data = await request.formData()
        const uuid = data.get('uuid')?.toString() ?? ''
        try {
            await deleteResource(`v1/tasks`, uuid)
        } catch (error) {
            console.error(error)
            return { status: 500 }
        }
    },
    updateTaskSummary: async ({ cookies, request }) => {
        const data = await request.formData()
        const uuid = data.get('uuid')?.toString() ?? ''
        const summary = data.get('summary')?.toString() ?? ''
        const task = { summary }
        try {
            await updateResource(`v1/tasks`, uuid, task)
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