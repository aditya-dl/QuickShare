import React, { useState } from "react";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { SharedItem } from "@/types";
import { toast } from "sonner";

interface CreateSnippetFromProps {
    onSnippetCreated: (newItem: SharedItem) => void;
}

export function CreateSnippetFrom({ onSnippetCreated }: CreateSnippetFromProps) {
    const [content, setContent] = useState('');
    const [isLoading, setIsLoading] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!content.trim()) {
            toast.warning("Snippet content cannot be empty.");
            return;
        }
        setIsLoading(true);

        try {
            const response = await fetch('/api/snippets', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ content: content }),
            });

            if (!response.ok) {
                const errorData = await response.json().catch(() => ({ error: `HTTP error: ${response.status}` }));
                throw new Error(errorData.error || 'Failed to create snippet');
            }

            const newItem: SharedItem = await response.json();
            toast.success(`Snippet "${newItem.name}" created!`);
            onSnippetCreated(newItem); // notify parent component
            setContent(''); // clear the form
        } catch (err) {
            console.error("Snippet creation error: ", err);
            const errorMessage = err instanceof Error ? err.message : 'An unknown error occurred';
            toast.error("Error creating snippet", { description: errorMessage });
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <Card>
            <CardHeader>
                <CardTitle>Create New Snippet</CardTitle>
            </CardHeader>
            <CardContent>
                <form onSubmit={handleSubmit}>
                    <Textarea placeholder="Paste or type your text snippet here..." value={content} onChange={(e) => setContent(e.target.value)} rows={5} className="mb-4" disabled={isLoading} />
                    <Button type="submit" disabled={isLoading || !content.trim()}>
                        {isLoading ? 'Creating...' : 'Create Snippet'}
                    </Button>
                </form>
            </CardContent>
        </Card>
    );
}