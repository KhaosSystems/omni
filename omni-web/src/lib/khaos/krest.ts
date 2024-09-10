/**
 * Helpers functions for Khaos REST API specification.
 * See: https://www.notion.so/khaosgroup/Khaos-Collective-REST-API-Specification-2024-WIP-9b276e93b64c46ccb09d25e9757b3161
 */

export async function getCollection(endpoint: string) {
  const res = await fetch(`http://localhost:30090/${endpoint}`);

  if (res.ok) {
    const body = await res.json();
    // Some APIs may (SHOULD NOT) return undefined instead of an empty array.
    // This is due to the fact that the go developers where high on something when they wrote the API...
    // We need to deal with this here.
    if (body.results === undefined) {
      console.warn(`API returned undefined instead of an empty array for ${endpoint}. This is non-spec compliant. Find the dev and bully them.`);
      return [];
    }

    return await body.results;
  } else {
    throw new Error(res.statusText);
  }
}

export async function getResource(endpoint: string, uuid: string) {
  const res = await fetch(`http://localhost:30090/${endpoint}/${uuid}`);
  if (res.ok) {
    var resource = await res.json();
    // Strip metadata (starting with @) from the resource.
    resource = Object.fromEntries(Object.entries(resource).filter(([key, value]) => !key.startsWith("@")));

    return resource;
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

export async function deleteResource(endpoint: string, uuid: string) {
  const res = await fetch(`http://localhost:30090/${endpoint}/${uuid}`, {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
    },
  });

  if (res.ok) {
    const body = await res.json();
    return await body.results;
  } else {
    throw new Error(res.statusText);
  }
}