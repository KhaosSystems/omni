import { getCollection } from "$lib/khaos/krest";

export async function getTasks() {
    return await getCollection("v1/tasks");
}
