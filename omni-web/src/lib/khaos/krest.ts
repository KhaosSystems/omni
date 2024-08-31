/**
 * Helpers functions for Khaos REST API specification.
 * See: https://www.notion.so/khaosgroup/Khaos-Collective-REST-API-Specification-2024-WIP-9b276e93b64c46ccb09d25e9757b3161
 */

export async function getCollection(endpoint: string) {
  const res = await fetch(`http://localhost:30090/${endpoint}`);

  if (res.ok) {
    const body = await res.json();
    return await body.results;
  } else {
    throw new Error(res.statusText);
  }
}

export async function createResource(endpoint: string, data: any) {
  const res = await fetch(`http://localhost:30090/${endpoint}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });

  if (res.ok) {
    const body = await res.json();
    return await body.results;
  } else {
    throw new Error(res.statusText);
  }
}