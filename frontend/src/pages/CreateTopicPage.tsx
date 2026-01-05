import { useNavigate } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "../hooks/redux";
import React, { useState } from "react";
import { clearError, createTopic } from "../features/topicsSlice"
import { Alert, Box, Button, CircularProgress, Container, Paper, TextField, Typography } from "@mui/material";
import { ArrowBack } from "@mui/icons-material";

export default function CreateTopicPage() {
    const dispatch = useAppDispatch();
    const navigate = useNavigate();

    const { submitting, submitError } = useAppSelector(state => state.topics);

    const [title, setTitle] = useState('');
    const [description, setDescription] = useState('');

    const handleTitleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setTitle(e.target.value);
        if (submitError) {
            dispatch(clearError());
        }
    }

    const handleDescriptionChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setDescription(e.target.value);
        if (submitError) {
            dispatch(clearError());
        }
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!title.trim() || !description.trim()) return;

        const result = await dispatch(
            createTopic({
                title: title.trim(),
                description: description.trim(),
            })
        );

        if (createTopic.fulfilled.match(result)) {
            // Redirect to topics page on successful creation
            navigate(`/topics`);
        }
    };

    return (
        <Container
            maxWidth="md"
            sx={{
                mt: 4, 
                mb: 4
            }}
        >
            <Button
                startIcon={<ArrowBack />}
                onClick={() => navigate('/topics')}
                variant="outlined"
                sx={{ mb: 3 }}
            >
                Back to Topics
            </Button>

            <Paper elevation={2} sx={{ p: 3 }}>
                <Typography 
                    variant="h4" 
                    component="h1"
                    gutterBottom>
                    New Topic
                </Typography>

                {submitError && (
                    <Alert severity="error" sx={{ mb: 2 }}>{submitError}</Alert>
                )}

                <form onSubmit={handleSubmit}>
                    <TextField
                        fullWidth
                        label="Title"
                        value={title}
                        onChange={handleTitleChange}
                        required
                        disabled={submitting}
                        sx={{ mb: 2 }}
                        placeholder="e.g. General Discussion"
                    />

                    <TextField 
                        fullWidth
                        label="Description"
                        value={description}
                        onChange={handleDescriptionChange}
                        required
                        disabled={submitting}
                        multiline
                        rows={4}
                        sx={{ mb: 3 }}
                        placeholder="Describe what this topic is about..."
                    />

                    <Box
                        sx={{ 
                            display: 'flex', 
                            gap: 2
                        }}    
                    >
                        <Button 
                            type="submit" 
                            variant="contained" 
                            disabled={submitting || !title.trim() || !description.trim()}
                        >
                                {submitting ? <CircularProgress size={24} /> : 'Create Topic'}
                        </Button>

                        <Button
                            variant="outlined"
                            onClick={() => navigate('/topics')}
                            disabled={submitting}
                        >
                            Cancel
                        </Button>
                    </Box>
                </form>
            </Paper>
        </Container>
    );
}