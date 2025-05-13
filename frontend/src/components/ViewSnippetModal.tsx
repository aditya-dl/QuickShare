import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { SharedItem } from "@/types";
import { CopyIcon, CheckIcon } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";

interface ViewSnippetModalProps {
    item: SharedItem | null;
    isOpen: boolean;
    onOpenChange: (open: boolean) => void;
}

export function ViewSnippetModal({ item, isOpen, onOpenChange }: ViewSnippetModalProps) {
    const [hasCopied, setHasCopied] = useState(false);

    const handleCopy = () => {
        if (item?.content) {
            navigator.clipboard.writeText(item.content)
                .then(() => {
                    setHasCopied(true);
                    toast.success("Snippet copied to clipboard!");
                    setTimeout(() => setHasCopied(false), 2000); // Reset icon after 2 seconds
                })
                .catch(err => {
                    console.error("Copy failed: ", err);
                    toast.error("Failed to copy snippet.");
                });
        }
    };

    if (!item || item.type !== 'text') {
        return null; // Don't render if no item or item is not text
    }

    return (
        <Dialog open={isOpen} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-[600px]">{/* Adjust max width */}
                <DialogHeader>
                    <DialogTitle className="truncate pr-10">{item.name}</DialogTitle> {/* Add pr for copy button */}
                    <DialogDescription>
                        Text Snippet
                    </DialogDescription>
                </DialogHeader>

                {/* Place copy button inside header or near content */}
                <Button
                    variant="ghost"
                    size="icon"
                    onClick={handleCopy}
                    className="absolute right-4 top-4" // Position top-right
                    aria-label="Copy snippet content"
                >
                    {hasCopied ? <CheckIcon className="h-4 w-4 text-green-500" /> : <CopyIcon className="h-4 w-4" />}
                </Button>

                <ScrollArea className="max-h-[60vh] my-4"> {/* Make content scrollable */}
                    <pre className="text-sm whitespace-pre-wrap break-words p-1">
                        {/* Using <pre> preserves whitespace and line breaks */}
                        {item.content}
                    </pre>
                </ScrollArea>

                {/* Optional: Footer for close button if needed, though ShadCN includes 'X' */}
                {/* <DialogFooter>
                    <Button variant="outline" onClick={() => onOpenChange(fase)}>Close</Button>
                    </DialogFooter> */}
            </DialogContent>
        </Dialog>
    );
}