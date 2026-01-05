import { useNavigate, useParams } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "../hooks/redux";
import { useEffect, useState } from "react";
import { clearError, fetchPostByID, updatePost } from "../features/postsSlice";
import { Alert, Box, Button, CircularProgress, Container, Paper, TextField, Typography } from "@mui/material";
import { ArrowBack } from "@mui/icons-material";

export default function EditPostPage() {
    const { postID } = useParams<{ postID: string }>();
    const dispatch = useAppDispatch();
    const navigate = useNavigate();

    const { currentPost, loading, submitting, submitError } = useAppSelector(state => state.posts);
    const { userID } = useAppSelector(state => state.auth);

    const [title, setTitle] = useState('');
    const [content, setContent] = useState('');

    // Fetch post on mount
    useEffect(() => {
        if (postID) {
            dispatch(fetchPostByID(parseInt(postID)));
        }
    }, [postID, dispatch]);

    // Populate form
    useEffect(() => {
        if (currentPost) {
            setTitle(currentPost.title);
            setContent(currentPost.content);
        }
    }, [currentPost]);

    // Check if user is author of post
    const isAuthor = currentPost && currentPost.createdBy === userID;

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

        if (!postID || !title.trim() || !content.trim()) return;

        const result = await dispatch(
            updatePost({
                postID: parseInt(postID),
                title: title.trim(),
                content: content.trim(),
            })
        );

        if (updatePost.fulfilled.match(result)) {
            // Redirect to post page on successful update
            navigate(`/posts/${postID}`);
        }
    }

    if (loading) {
        return (
            <Container
                sx={{
                    display: 'flex',
                    justifyContent: 'center',
                    alignItems: 'center',
                    minHeight: '60vh',
                }}
            >
                <CircularProgress />
            </Container>
        );
    }

    if (!currentPost) {
        return (
            <Container sx={{ mt: 4 }}>
                <Alert severity="error">Post not found.</Alert>
                <Button
                    startIcon={<ArrowBack />}
                    onClick={() => navigate(-1)}
                    sx={{ mt: 2 }}
                >
                    Back
                </Button>
            </Container>
        );
    }

    if (!isAuthor) {
        return (
            <Container sx={{ mt: 4 }}>
                <Alert severity="error">You are not authorised to edit this post.</Alert>
                <Button
                    startIcon={<ArrowBack />}
                    onClick={() => navigate(`/posts/${postID}`)}
                    sx={{ mt: 2 }}
                >
                    Back to Post
                </Button>
            </Container>
        );
    }

    return (
        <Container
            maxWidth="md"
            sx={{
                mt: 4,
                mb: 4,
            }}
        >
            <Button
                startIcon={<ArrowBack />}
                onClick={() => navigate(`/posts/${postID}`)}
                sx={{ mb: 3 }}
                variant="outlined"
            >
                Back to Post
            </Button>

            <Paper elevation={2} sx={{ p: 3 }}>
                <Typography
                    variant="h4"
                    component="h1"
                    gutterBottom
                >
                    Edit Post
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
                            {submitting ? <CircularProgress size={24} /> : 'Save Changes'}
                        </Button>

                        <Button
                            variant="outlined"
                            onClick={() => navigate(`/posts/${postID}`)}
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