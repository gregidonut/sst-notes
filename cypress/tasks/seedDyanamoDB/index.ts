export default async function main(url: string, token: string) {
  try {
    console.log("sent fetch");
    const resp = await fetch(`${url}/seedNotes`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    console.log("got fetch back");

    if (!resp || !resp.ok) {
      return new Response(JSON.stringify({ message: "Failed to fetch" }), {
        status: 500,
      });
    }
  } catch (e) {
    console.error(e);
  }
  return new Response(JSON.stringify({ status: "successfully seeded db" }), {
    status: 200,
  });
}
