import { redirect } from "@sveltejs/kit";

/** @type {import('./$types').PageServerLoad} */
export async function load({ cookies }) {
    const sessionid = cookies.get('sessionid');
    if (!sessionid || sessionid == 'invalidsessionid') {
        return redirect(303, '/login');
    } else {
        return redirect(303, '/projects');
    }
}