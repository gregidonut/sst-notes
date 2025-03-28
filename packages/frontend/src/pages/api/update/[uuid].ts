import type { APIRoute } from "astro";

export const PUT: APIRoute = async function ({ request, params, locals }) {
    const formData = await request.formData();
    const content = formData.get("content");
    const { uuid } = params;
    try {
        console.log(`updating resource with ID: ${uuid}`);

        const { getToken } = locals.auth();
        const token = await getToken();

        const ASTRO_API_URL = process.env.ASTRO_API_URL as string;
        await fetch(`${ASTRO_API_URL}/notes/${uuid}`, {
            method: "PUT",
            headers: {
                Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify({
                content: content,
            }),
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
