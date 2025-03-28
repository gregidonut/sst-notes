import type { APIRoute } from "astro";

export const POST: APIRoute = async ({ request, locals }) => {
    const formData = await request.formData();
    const content = formData.get("content");
    const attachment = formData.get("attachment");
    if (!content) {
        return new Response(
            JSON.stringify({
                message: "Missing required fields",
            }),
            { status: 400 },
        );
    }

    const { getToken } = locals.auth();
    const token = await getToken();
    const ASTRO_API_URL = process.env.ASTRO_API_URL as string;
    let resp;
    try {
        console.log("sent fetch");
        resp = await fetch(`${ASTRO_API_URL}/notes`, {
            method: "POST",
            headers: {
                Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify({
                content: content,
                attachment: attachment ?? attachment,
            }),
        });
        console.log("got fetch back");
    } catch (e) {
        console.error(e);
    }
    if (!resp || !resp.ok) {
        return new Response(JSON.stringify({ message: "Failed to fetch" }), {
            status: resp!.status,
        });
    }

    return new Response(null, {
        status: 302,
        headers: {
            Location: "/notes",
        },
    });
};
