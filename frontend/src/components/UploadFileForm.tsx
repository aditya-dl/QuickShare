import React, { useState, useRef } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { SharedItem } from "@/types";
import { toast } from "sonner";
import { UploadIcon } from "lucide-react";
import { formatFileSize } from "@/lib/utils";

interface UploadFileFormProps {
    onFileUploaded: (newItem: SharedItem) => void;
}

export function UploadFileForm({ onFileUploaded }: UploadFileFormProps) {
    const [selectedFile, setSelectedFile] = useState<File | null>(null);
    const [isLoading, setIsLoading] = useState(false);
    const fileInputRef = useRef<HTMLInputElement>(null); // Ref to clear file input

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.files && e.target.files.length > 0) {
            setSelectedFile(e.target.files[0]);
        } else {
            setSelectedFile(null);
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!selectedFile) {
            toast.warning("Please select a file to upload.");
            return;
        }
        setIsLoading(true);

        const formData = new FormData();
        formData.append('file', selectedFile);
        // Optional: append name if you add an input field for custom name
        // formData.append('name', customName);

        try {
            const response = await fetch('/api/files', {
                method: 'POST',
                body: formData, // No content-type header needed for formdata
            });

            if (!response.ok) {
                const errorData = await response.json().catch(() => ({ error: `HTTP error: ${response.status}` }));
                throw new Error(errorData.error || 'Failed to upload file');
            }

            const newItem: SharedItem = await response.json();
            toast.success(`File "${newItem.name}" uploaded!`);
            onFileUploaded(newItem); // notify parent

            // clear the form state
            setSelectedFile(null);
            if (fileInputRef.current) {
                fileInputRef.current.value = "";
            }
        } catch (err) {
            console.error("File upload error: ", err);
            const errorMessage = err instanceof Error ? err.message : 'An unknown error occurred';
            toast.error("Error uploading file", { description: errorMessage });
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <Card>
            <CardHeader>
                <CardTitle>Upload File</CardTitle>
                <CardDescription>Share a file on your local network.</CardDescription>
            </CardHeader>
            <CardContent>
                <form onSubmit={handleSubmit} className="space-y-4">
                    <Input type="file" ref={fileInputRef} onChange={handleFileChange} disabled={isLoading} className="file:mr-4 file:py-2 file:px4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-primary file:text-primary-foreground hover:file:bg-primary/90"/>

                    {selectedFile && (
                        <p className="text-sm text-muted-foreground">
                            Selected: {selectedFile.name} ({formatFileSize(selectedFile.size)})
                        </p>
                    )}

                    <Button type="submit" disabled={isLoading || !selectedFile}>
                        <UploadIcon className="mr-2 h-4 w-4" />
                        {isLoading ? 'Uploading...' : 'Upload File'}
                    </Button>
                </form>
            </CardContent>
        </Card>
    );
}