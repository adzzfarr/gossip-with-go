import { useNavigate, useParams } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "../hooks/redux";
import { useEffect, useState } from "react";
import { fetchPostsByTopic } from "../features/postsSlice";
import { Alert, Box, Button, Card, CardActionArea, CardContent, CircularProgress, Container, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle, Divider, Paper, Typography } from "@mui/material";
import { Add, ArrowBack, Delete, Edit } from "@mui/icons-material";
import ForumBreadcrumbs from "../components/Breadcrumbs";
import Username from "../components/Username";
import { deleteTopic } from "../features/topicsSlice";

export default function TopicPostsPage() {
    const { topicID } = useParams<{ topicID: string }>();
    const dispatch = useAppDispatch(); 
    const navigate = useNavigate();

    const { posts, loading: postsLoading, error: postsError } = useAppSelector(state => state.posts);
    const { topics, submitting: topicSubmitting, submitError } = useAppSelector(state => state.topics);
    const { userID } = useAppSelector(state => state.auth);

    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

    // Find current topic from topics list
    const topic = topics.find(t => t.topicID === parseInt(topicID || '0'));
    
    // Check if user is author of topic
    const isAuthor = topic && topic.createdBy === userID;

    useEffect(() => {
        if (topicID) {
            dispatch(fetchPostsByTopic(parseInt(topicID)));
        }
    }, [dispatch, topicID]);

    const handleDeleteTopic = async () => {
        if (!topic) return;

        const result = await dispatch(deleteTopic(topic.topicID));

        if (deleteTopic.fulfilled.match(result)) {
            setDeleteDialogOpen(false);
            navigate('/topics');
        }
    }


    if (postsLoading) {
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

    if (postsError) {
        return (
            <Container sx={{ mt: 4 }}>
                <Alert severity="error">{postsError}</Alert>
                <Button
                    startIcon={<ArrowBack />}
                    onClick={() => navigate('/topics')}
                    sx={{ mt: 2 }}    
                >
                    Back to Topics
                </Button>
            </Container>
        );
    }

    if (!topic) {
        return (
            <Container>
                <Alert severity="info">Topic not found.</Alert>
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
            sx={{ 
                mt: 4,
                mb: 4,
            }}
            maxWidth="lg"
        >
            <ForumBreadcrumbs />

            <Button
                startIcon={<ArrowBack />}
                variant="outlined"
                onClick={() => navigate('/topics')}
                sx={{ mb: 3 }}
            >
                Back to Topics
            </Button>
            
            {/* Topic Details */}
            <Paper
                elevation={2}
                sx={{ 
                    p: 2, 
                    mb: 3 
                }}
            >
                <Box
                    sx={{
                        display: 'flex',
                        justifyContent: 'space-between',
                        alignItems: 'flex-start',
                        mb: 2,
                    }}
                >
                    <Typography variant="h4" component="h1">
                        {topic.title}
                    </Typography>

                    {/* Edit and Delete Buttons */}
                    {isAuthor && (
                        <Box 
                            sx={{ 
                                display: 'flex', 
                                gap: 1
                            }}
                        >
                            <Button
                                startIcon={<Edit />}
                                variant="outlined"
                                onClick={() => navigate(`/topics/${topic.topicID}/edit`)}
                                size="small"
                            >
                                Edit
                            </Button>

                            <Button
                                startIcon={<Delete />}
                                variant="outlined"
                                color="error"
                                size="small"
                                disabled={topicSubmitting}
                                onClick={() => setDeleteDialogOpen(true)}
                            >
                                Delete
                            </Button>
                        </Box>
                    )}
                </Box>

                <Box
                    sx={{
                        display: 'flex',
                        alignItems: 'center',
                        gap: 2,
                        mb: 2,
                    }}
                >
                    <Box
                        sx={{
                            display: 'flex',
                            alignItems: 'center',
                            gap: 0.5
                        }}
                    >
                        <Typography variant="body2" color="text.secondary">
                            Created by
                        </Typography>
                        <Username
                            username={topic.username || "Unknown"}
                            userID={topic.createdBy}
                            variant="body2"
                            color="text.secondary"
                        />
                    </Box>

                    <Typography variant="body2" color="text.secondary">•</Typography>

                    <Typography variant="body2" color="text.secondary">
                        {new Date(topic.createdAt).toLocaleString()}
                    </Typography>
                </Box>

                <Divider sx={{ my: 2 }} />

                <Typography variant="body1" sx={{ whiteSpace: 'pre-wrap' }}>
                    {topic.description || 'No description provided.'}
                </Typography>
            </Paper>

            {/* Posts Section */}
            <Box>
                <Box
                    sx={{
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'space-between',
                        mb: 3,
                    }}    
                >
                    <Typography variant="h5" component="h2">
                        Posts
                    </Typography>

                    <Button
                        startIcon={<Add />}
                        variant="contained"
                        onClick={() => navigate(`/topics/${topicID}/create-post`)}
                    >
                        Create Post 
                    </Button>
                </Box>

                {posts.length === 0
                ? (
                    <Alert severity="info">
                        No posts available for this topic. Be the first to create one!
                    </Alert>
                )
                : (
                    <Box
                        sx={{
                            display: 'flex',
                            flexDirection: 'column',
                            gap: 2,
                        }}
                    >
                        {posts.map(
                            (post) => (
                                <Card key={post.postID} variant="outlined"> 
                                    <CardActionArea onClick={() => navigate(`/posts/${post.postID}`)}>
                                        <CardContent>
                                            <Typography
                                                variant="h6"
                                                component="h2"
                                                gutterBottom
                                            >
                                                {post.title}
                                            </Typography>

                                            <Typography
                                                variant="body2"
                                                color="text.secondary"
                                                sx={{
                                                    display: '-webkit-box',
                                                    WebkitLineClamp: 3,
                                                    WebkitBoxOrient: 'vertical',
                                                    overflow: 'hidden',
                                                    textOverflow: 'ellipsis',
                                                    mb: 2,

                                                }}
                                            >
                                                {post.content}
                                            </Typography>

                                            <Box
                                                sx={{
                                                    display: 'flex',
                                                    gap: 2,
                                                    alignItems: 'center',
                                                    mt: 2
                                                }}
                                            >
                                                <Box
                                                    sx={{
                                                        display: 'flex',
                                                        alignItems: 'center',
                                                        gap: 0.5
                                                    }}
                                                >
                                                    <Typography variant="caption" color="text.secondary">
                                                        Posted by
                                                    </Typography>

                                                    <Username
                                                        username={post.username || "Unknown"}
                                                        userID={post.createdBy}
                                                        variant="caption"
                                                        color="text.secondary"
                                                    />
                                                </Box>
                                                <Typography variant="caption" color="text.secondary">•</Typography>
                                                <Typography variant="caption" color="text.secondary">
                                                    {new Date(post.createdAt).toLocaleString()}
                                                </Typography>
                                            </Box>
                                        </CardContent>
                                    </CardActionArea>
                                </Card>
                            )
                        )}
                    </Box>
                )}
            </Box>

            {/* Delete Confirmation Dialog */}
            <Dialog
                open={deleteDialogOpen}
                onClose={() => setDeleteDialogOpen(false)}
            >
                <DialogTitle>Delete Topic</DialogTitle>
                <DialogContent>
                    {submitError && (
                        <Alert severity="error" sx={{ mb: 2 }}>{submitError}</Alert>
                    )}
                    <DialogContentText>
                        Are you sure you want to delete this topic? This will also delete all posts and comments within it. This action cannot be undone.
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setDeleteDialogOpen(false)}>
                        Cancel
                    </Button>
                    <Button 
                        onClick={handleDeleteTopic} 
                        color="error"
                        disabled={topicSubmitting}
                    >
                        {topicSubmitting ? <CircularProgress size={24} /> : 'Delete'}
                    </Button>
                </DialogActions>
            </Dialog>
        </Container>
    );
}