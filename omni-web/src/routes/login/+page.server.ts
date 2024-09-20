import { redirect } from "@sveltejs/kit";

/** @type {import('./$types').PageServerLoad} */
export async function load({ cookies }) {
  return {
    sessionid: cookies.get("sessionid"),
  };
}

/** @type {import('./$types').Actions} */
export const actions = {
  login: async ({ cookies, request }) => {
    // Reference: https://kit.svelte.dev/docs/form-actions
    const data = await request.formData();
    const username = data.get("username");
    const password = data.get("password");

    if (username !== "admin" || password !== "password") {
      return { status: 401 };
    }

    cookies.set("sessionid", "validsessionid", { path: "/" });
    redirect(303, "/");
  },
};
