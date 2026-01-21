import { Alert, Button, CircularProgress, Paper, TextField } from "@mui/material";
import { useState, type FormEvent } from "react";

interface CommentFormProps {
    onSubmit: (content: string) => void;
    isSubmitting: boolean;
    error?: string | null;
    placeholder?: string;
}

export default function CommentForm({
    onSubmit,
    isSubmitting,
    error,
    placeholder = "Write a comment..."
}: CommentFormProps) {
    const [content, setContent] = useState('');

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();

        if (!content.trim()) return;

        await onSubmit(content.trim());
        setContent('');
    };

    return (
        <Paper
            elevation={1}
            sx={{ p: 2 }}
        >
            <form onSubmit={handleSubmit}>
                <TextField 
                    fullWidth
                    multiline
                    minRows={3}
                    placeholder={placeholder}
                    value={content}
                    onChange={(e) => setContent(e.target.value)}
                    disabled={isSubmitting}
                    sx={{ mb: 2 }}
                />

                {error && (
                    <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>
                )}

                <Button 
                    type="submit" 
                    disabled={isSubmitting}
                >
                    {isSubmitting ? <CircularProgress size={24} /> : 'Post Comment'}
                </Button>
            </form>
        </Paper>
    )
}