"use client"; 

import { useItems } from "@/hooks/useItems";
import { CreateSnippetFrom } from "@/components/CreateSnippetForm";
import { UploadFileForm } from "@/components/UploadFileForm";
import { ItemCard } from "@/components/ItemCard";
import { Skeleton } from "@/components/ui/skeleton";
import { Toaster } from "sonner";
import { Card, CardContent, CardHeader } from "@/components/ui/card";

export default function HomePage() {
  // Use the custom hook to manage items
  const { items, isLoading, error, addItem, deleteItem } = useItems();

  return (
    <main className="container mx-auto px-4 py-8">
      <Toaster position="top-center" richColors /> 
      <h1 className="text-3xl font-bold mb-6 text-center">QuickShare</h1>
      <p className="text-center text-muted-foreground mb-8">
        Share text and files across your local network easily.
      </p>

      {/* Form Side-By-Side (or stacked on small screens) */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-10">
        <CreateSnippetFrom onSnippetCreated={addItem} />
        <UploadFileForm onFileUploaded={addItem} />
      </div>

      {/* Items List */}
      <h2 className="text-2xl font-semibold mb-4 border-b pb-2">Shared Items</h2>

      {isLoading && (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {/* Skeleton Loaders */}
          {[...Array(3)].map((_, i) => (
            <Card key={i}>
              <CardHeader>
                <Skeleton className="h-6 w-3/4 mb-2" />
                <Skeleton className="h-4 w-1/2" />
              </CardHeader>
              <CardContent>
                <Skeleton className="h-8 w-1/4 ml-auto" /> {/* Simulate button area */}
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {error && (
        <div className="text-center text-red-600 bg-red-100 p-4 rounded-md">
          Error loading items: {error}
        </div>
      )}

      {!isLoading && !error && items.length === 0 && (
        <p className="text-center text-muted-foreground mt-6">
          No items shared yet. Use the forms above to add snippets or files.
        </p>
      )}

      {!isLoading && !error && items.length > 0 && (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {items.map((item) => (
            <ItemCard key={item.id} item={item} onDelete={deleteItem} />
          ))}
        </div>
      )}
    </main>
  );
}