import { useNavigate, useParams } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "../hooks/redux";
import { fetchPostById } from "../features/posts/postsSlice";
import { clearCommentsError, createComment, fetchCommentsByPostID } from "../features/comments/commentsSlice";
import { useEffect, useState } from "react";
import { Alert, Box, Button, Card, CardContent, CircularProgress, Container, Divider, Paper, TextField, Typography } from "@mui/material";
import { ArrowBack } from "@mui/icons-material";
import ForumBreadcrumbs from "../components/Breadcrumbs";

export default function PostPage() {
    const { postID } = useParams<{ postID: string}>();
    const dispatch = useAppDispatch();
    const navigate = useNavigate();

    const { currentPost, loading: postLoading, error: postError } = useAppSelector(state => state.posts);
    const { comments, loading: commentsLoading, error: commentsError, submitting, submitError } = useAppSelector(state => state.comments);

    const [commentContent, setCommentContent] = useState('');

    useEffect(() => {
        if (postID) {
            dispatch(fetchPostById(parseInt(postID)));
            dispatch(fetchCommentsByPostID(parseInt(postID)));
        }
    }, [postID, dispatch]);

    const handleSubmitComment = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!postID || !commentContent.trim()) return;

        const result = await dispatch(
            createComment({
                postID: parseInt(postID),
                content: commentContent.trim(),
            })
        );

        if (createComment.fulfilled.match(result)) {
            // Clear comment input on successful submission
            setCommentContent('');
        }
    }

    const handleCommentChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setCommentContent(e.target.value);
        if (submitError) {
            dispatch(clearCommentsError());
        }
    }

    if (postLoading || commentsLoading) {
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

    if (postError) {
        return (
            <Container sx={{ mt: 4 }}>
                <Alert severity="error">{postError}</Alert>
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

    if (!currentPost) {
        return (
            <Container>
                <Alert severity="info">Post not found.</Alert>
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

    return (
        <Container
            maxWidth="md"
            sx={{
                mt: 4,
                mb: 4,
            }}
        >
            <ForumBreadcrumbs />
            
            <Button
                startIcon={<ArrowBack />}
                onClick={() => navigate(-1)}
                sx={{ mb: 3 }}
            >
                Back
            </Button>

            {/* Post Content */}
            <Paper
                elevation={2}
                sx={{
                    p: 3,
                    mb: 4,
                }}
            >
                <Typography 
                    variant="h4" 
                    component="h1" 
                    gutterBottom
                >
                    {currentPost.title}
                </Typography>

                <Box
                    sx={{
                        display: 'flex',
                        gap: 2,
                        mb: 2,
                    }}
                >
                    <Typography variant="body2" color="text.secondary">
                        By: {currentPost.createdBy || 'Unknown'} 
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        {new Date(currentPost.createdAt).toLocaleString()}
                    </Typography>
                </Box>

                <Divider sx={{ my: 2 }} />

                <Typography variant="body1" sx={{ whiteSpace: 'pre-wrap' }}>
                    {currentPost.content}
                </Typography>
            </Paper>

            {/* Comments Section */}
            <Box>
                <Typography 
                    variant="h5" 
                    component="h2" 
                    gutterBottom
                >
                    Comments
                </Typography>

                {/* Add Comment Form */}
                <Paper
                    elevation={1}
                    sx={{ 
                        p: 2, 
                        mb: 3 
                    }}
                >
                    <form onSubmit={handleSubmitComment}>
                        <TextField
                            fullWidth
                            multiline
                            minRows={3}
                            placeholder="Write a comment..."
                            value={commentContent}
                            onChange={handleCommentChange}
                            disabled={submitting}
                            sx={{ mb: 2 }}
                        />
                        
                        {submitError && (
                            <Alert severity="error" sx={{ mb: 2 }}>{submitError}</Alert>
                        )}

                        <Button
                            type="submit"
                            variant="contained"
                            disabled={submitting || !commentContent.trim()}
                        >
                            {submitting ? <CircularProgress size={24} /> : 'Post Comment'}
                        </Button>
                    </form>
                </Paper>

                {commentsError && (
                    <Alert severity="error" sx={{ mb: 2 }}>
                        {commentsError}
                    </Alert>
                )}

                {comments.length === 0
                    ? (<Alert severity="info">No comments yet. Be the first to comment!</Alert>)
                    : (
                        <Box
                            sx={{
                                display: 'flex',
                                flexDirection: 'column',
                                gap: 2,
                            }}
                        >
                            {comments.map((comment) => (
                                <Card key={comment.commentID} variant="outlined">
                                    <CardContent>
                                        <Box
                                            sx={{
                                                display: 'flex',
                                                gap: 2,
                                                mb: 1,
                                            }}
                                        >
                                            <Typography variant="body2" fontWeight="bold">
                                                {comment.createdBy || 'Unknown'}
                                            </Typography>
                                            <Typography variant="body2" color="text.secondary">
                                                {new Date(comment.createdAt).toLocaleString()}
                                            </Typography>
                                        </Box>
                                        <Typography variant="body1" sx={{ whiteSpace: 'pre-wrap' }}>
                                            {comment.content}
                                        </Typography>
                                    </CardContent>
                                </Card>
                            ))}
                        </Box>
                    )
                }
            </Box>
        </Container>
    );
}