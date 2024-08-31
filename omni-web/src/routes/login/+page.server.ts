import { redirect } from '@sveltejs/kit'

/** @type {import('./$types').PageServerLoad} */
export async function load({ cookies }) {
    return {
        sessionid: cookies.get('sessionid')
    }
}

/** @type {import('./$types').Actions} */
export const actions = {
    login: async ({ cookies, request }) => {
        cookies.set('sessionid', 'validsessionid', { path: '/' })
        redirect(303, '/')
    }
}