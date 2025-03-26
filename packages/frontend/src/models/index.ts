/* Do not change, this code is generated from Golang structs */

export class Note {
    attachment: string;
    content: string;
    createdAt: string;
    noteId: string;
    userId: string;

    constructor(source: any = {}) {
        if ("string" === typeof source) source = JSON.parse(source);
        this.attachment = source["attachment"];
        this.content = source["content"];
        this.createdAt = source["createdAt"];
        this.noteId = source["noteId"];
        this.userId = source["userId"];
    }
}
