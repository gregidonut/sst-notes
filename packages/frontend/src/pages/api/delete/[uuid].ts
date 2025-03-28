import type { APIRoute } from "astro";

export const DELETE: APIRoute = async ({ params, locals }) => {
    const { uuid } = params;

    if (!uuid) {
        return new Response(JSON.stringify({ error: "ID is required" }), {
            status: 400,
            headers: { "Content-Type": "application/json" },
        });
    }

    try {
        console.log(`Deleting resource with ID: ${uuid}`);

        const { getToken } = locals.auth();
        const token = await getToken();

        const ASTRO_API_URL = process.env.ASTRO_API_URL as string;
        await fetch(`${ASTRO_API_URL}/notes/${uuid}`, {
            method: "DELETE",
            headers: {
                Authorization: `Bearer ${token}`,
            },
        });
    } catch (error) {
        return new Response(JSON.stringify({ error: "Deletion failed" }), {
            status: 500,
            headers: { "Content-Type": "application/json" },
        });
    }
    return new Response(null, {
        status: 200,
    });
};
