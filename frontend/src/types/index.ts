export type ItemType = "text" | "file";

export interface SharedItem {
    id: string;
    name: string;
    type: ItemType;
    createdAt: string;
    content?: string;
    fileName?: string;
    contentType?: string;
    size?: string;
}