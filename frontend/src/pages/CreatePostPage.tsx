import { useNavigate, useParams } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "../hooks/redux";
import { useState } from "react";
import { clearError, createPost } from "../features/posts/postsSlice";
import { Alert, Box, Button, CircularProgress, Container, Paper, TextField, Typography } from "@mui/material";
import { ArrowBack } from "@mui/icons-material";
import ForumBreadcrumbs from "../components/Breadcrumbs";

export default function CreatePostPage() {
    const { topicID } = useParams<{ topicID: string }>();
    const dispatch = useAppDispatch();
    const navigate = useNavigate();

    const { submitting, submitError } = useAppSelector(state => state.posts);

    const [title, setTitle] = useState('');
    const [content, setContent] = useState('');

    const handleTitleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setTitle(e.target.value);
        if (submitError) {
            dispatch(clearError());
        }
    }

    const handleContentChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setContent(e.target.value);
        if (submitError) {
            dispatch(clearError());
        }
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!topicID || !title.trim() || !content.trim()) return;

        const result = await dispatch(
            createPost({
                topicID: parseInt(topicID),
                title: title.trim(),
                content: content.trim(),
            })
        );

        if (createPost.fulfilled.match(result)) {
            // Redirect to topic posts page on successful creation
            navigate(`/topics/${topicID}`);
        }
    }

    return (
        <Container
            maxWidth="md"
            sx={{ 
                mt: 4, 
                mb: 4 
            }}
        >
            <ForumBreadcrumbs />
            
            <Button
                startIcon={<ArrowBack />}
                onClick={() => navigate(`/topics/${topicID}`)}
                variant="outlined"
                sx={{ mb: 3 }}
            >
                Back to Posts
            </Button>

            <Paper elevation={2} sx={{ p: 3 }}>
                <Typography 
                    variant="h4" 
                    component="h1" 
                    gutterBottom
                >
                    Create New Post
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
                    />

                    <TextField
                        fullWidth
                        label="Content"
                        value={content}
                        onChange={handleContentChange}
                        required
                        disabled={submitting}
                        multiline
                        rows={10}
                        sx={{ mb: 3 }}
                    />

                    <Box
                        sx={{
                            display: 'flex',
                            gap: 2,
                        }}
                    >
                        <Button
                            type="submit"
                            variant="contained"
                            disabled={submitting || !title.trim() || !content.trim()}
                        >
                            {submitting ? <CircularProgress size={24} /> : 'Create Post'}
                        </Button>

                        <Button
                            variant="outlined"
                            onClick={() => navigate(-1)}
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