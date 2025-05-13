// Custom hook for managing items state
import { useState, useEffect, useCallback } from "react";
import { SharedItem } from "@/types";
import { toast } from "sonner";

const API_BASE_URL = '/api'; // using next.js proxy

export function useItems() {
    const [items, setItems] = useState<SharedItem[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const fetchItems = useCallback(async () => {
        setIsLoading(true);
        setError(null);
        try {
            const response = await fetch(`${API_BASE_URL}/items`);
            if (!response.ok) {
                throw new Error(`Failed to fetch items: ${response.statusText}`);
            }
            const data: SharedItem[] = await response.json();
            setItems(data);
        } catch (err) {
            console.error("Fetch error: ", err);
            const errorMessage = err instanceof Error ? err.message : 'An unknown error occured';
            setError(errorMessage);
            toast.error("Error fetching items", { description: errorMessage });
        } finally {
            setIsLoading(false);
        }
    }, []);

    useEffect(() => {
        fetchItems();
    }, [fetchItems]);

    const addItem = useCallback((newItem: SharedItem) => {
        // Optimistic update (optional, can make UI feel faster)
        setItems((prevItems) => [newItem, ...prevItems]);
        // Fetch again to get the final list order and confirm
        fetchItems();
    }, [fetchItems]);

    const deleteItem = useCallback(async (id: string) => {
        // optimistic update
        setItems((prevItems) => prevItems.filter(item => item.id !== id));
        toast.info("Deleting item...");

        try {
            const response = await fetch(`${API_BASE_URL}/items/${id}`, {
                method: 'DELETE',
            });
            if (!response.ok) {
                // Revert optimistic update if delete fails
                fetchItems(); // Refetch to get actual state
                throw new Error(`Failed to delete item: ${response.status} ${response.statusText}`);
            }
            // On success, the item is already removed optimistically 
            toast.success("Item deleted successfully");
            // No need to call fetchItems() again if optimistic update worked
        } catch (err) {
            console.error("Delete error: ", err);
            const errorMessage = err instanceof Error ? err.message : 'An unknown error occured';
            setError(errorMessage); // Keep track of the error state
            toast.error("Error deleting item", { description: errorMessage });
            // Refetch to ensure UI matches backend state after error 
            fetchItems();
        }
    }, [fetchItems]); // Add fetchItems dependency

    return { items, isLoading, error, fetchItems, addItem, deleteItem };
}