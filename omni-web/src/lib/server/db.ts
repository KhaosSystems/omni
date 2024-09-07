import  * as krest from "$lib/khaos/krest";

export async function getTasks() {
    try {
        const tasks = await krest.getCollection("v1/tasks");
        return tasks;
    } catch (error) {
        console.error("Error fetching tasks:", error);
        return [];
    }
}
