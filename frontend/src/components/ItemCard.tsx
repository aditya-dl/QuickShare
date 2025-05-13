import { SharedItem } from "@/types";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { formatFileSize, formatRelativeDate } from "@/lib/utils";
import { FileTextIcon, FileIcon, Trash2Icon, DownloadIcon, EyeIcon } from "lucide-react";
import React, { useState } from "react";
import { ViewSnippetModal } from "./ViewSnippetModal"; // Import the modal

interface ItemCardProps {
    item: SharedItem;
    onDelete: (id: string) => void;
}

export function ItemCard({ item, onDelete }: ItemCardProps) {
    const [isModalOpen, setIsModalOpen] = useState(false);

    const handleCardClick = () => {
        if (item.type === 'text') {
            setIsModalOpen(true);
        } else if (item.type === 'file') {
            // Direct download link
            window.location.href = `/api/files/${item.id}/download`;
        }
    };

    const handleDeleteClick = (e: React.MouseEvent) => {
        e.stopPropagation(); // Prevent card click event when deleting
        onDelete(item.id);
    };

    return (
        <>
            <Card className="hover:shadow-md transition-shadow cursor-pointer group" onClick={handleCardClick}>
                <CardHeader className="pb-2 flex flex-row items-start justify-between space-y-0">
                    <div className="flex-1 overflow-hidden"> {/* Container for text */}
                        <CardTitle className="text-lg font-medium truncate" title={item.name}> {/* Tooltip for long names */}
                            {item.name}
                        </CardTitle>
                        <CardDescription className="text-xs text-muted-foreground pt-1">
                            {formatRelativeDate(item.createdAt)}
                            {item.type === 'file' && item.size != null && ` • ${formatFileSize(item.size)}`}
                            {item.type === 'file' && item.fileName && ` • ${item.fileName}`}
                        </CardDescription>
                    </div>
                    {/* Icon indicating type */}
                    <div className="pl-2">
                        {item.type === 'text' ? (
                            <FileTextIcon className="h-5 w-5 text-muted-foreground" />
                        ) : (
                            <FileIcon className="h-5 w-5 text-muted-foreground" />
                        )}
                    </div>
                </CardHeader>
                <CardContent className="relative pt-2">
                    {/* Action buttons - maybe show on hover or keep minimal? */}
                    {/* Let's show delete always for now, download/view implied by click */}
                    <div className="absolute bottom-2 right-2 flex items-center space-x-1 opacity-50 group-hover:opacity-100 transition-opacity">
                        {/* Optional: Explicit View/Download icons if needed */}
                        {/* {item.type === 'text' && (
                            <Button variant="ghost" size="icon" className="h-7 w-7" aria-label="View Snippet" onClick={(e) => { e.stopPropagation(); setIsModalOpen(true); }}>
                                <EyeIcon className="h-4 w-4" />
                            </Button>
                        )}
                        {item.type === 'file' && (
                            <a href={`/api/files/${item.id}/download`} download={item.fileName} onClick={(e) => e.stopPropagation()} aria-label="Download File">
                                <Button variant="ghost" size="icon" className="h-7 w-7">
                                    <DownloadIcon className="h-4 w-4" />
                                </Button>
                            </a>
                        )} */}

                        {/* Delete Button */}
                        <Button variant="ghost" size="icon" className="h-7 w-7 text-destructive hover:text-destructive hover:bg-destructive/10" onClick={handleDeleteClick} aria-label="Delete Item">
                            <Trash2Icon className="h-4 w-4" />
                        </Button>
                    </div>
                </CardContent>
            </Card>

            {/* Modal for viewing text snippets */}
            {item.type === 'text' && (
                <ViewSnippetModal item={item} isOpen={isModalOpen} onOpenChange={setIsModalOpen} />
            )}
        </>
    );
}