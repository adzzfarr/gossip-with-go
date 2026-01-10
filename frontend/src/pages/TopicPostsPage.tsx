import { useNavigate, useParams } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "../hooks/redux";
import { useEffect, useState } from "react";
import { fetchPostsByTopic } from "../features/postsSlice";
import { Alert, Box, Button, Card, CardActionArea, CardContent, CircularProgress, Container, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle, Divider, Grid, IconButton, InputAdornment, Paper, TextField, Typography } from "@mui/material";
import { Add, ArrowBack, Clear, Delete, Edit, Search } from "@mui/icons-material";
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
    const [searchQuery, setSearchQuery] = useState('');

    // Find current topic from topics list
    const topic = topics.find(t => t.topicID === parseInt(topicID || '0'));
    
    // Check if user is author of topic
    const isAuthor = topic && topic.createdBy === userID;

    useEffect(() => {
        if (topicID) {
            dispatch(fetchPostsByTopic(parseInt(topicID)));
        }
    }, [dispatch, topicID]);

    const filteredPosts = posts.filter(
        post => {
            const query = searchQuery.trim().toLowerCase();
            return (
                post.title.toLowerCase().includes(query) ||
                post.content.toLowerCase().includes(query)
            );
        }
    )

    const handleClearSearch = () => {
        setSearchQuery('');
    }

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
                {/* Header (Search Bar + New Post) */}
                <Box
                    sx={{
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'space-between',
                        flexWrap: 'wrap',
                        mb: 3,
                        gap: 2,
                    }}    
                >
                    <Box
                        sx={{
                            display: 'flex',
                            alignItems: 'center',
                            gap: 2,
                            flex: 1,
                            maxWidth: 600,
                        }}
                    >
                        <Typography variant="h5" component="h2">
                            Posts
                        </Typography>
                        
                        {/* Search Bar */}
                        <TextField 
                            size="small"
                            fullWidth
                            placeholder="Search Posts..."
                            value={searchQuery}
                            onChange={(e) => setSearchQuery(e.target.value)}
                            slotProps={{
                                input: {
                                    'startAdornment': (
                                        <InputAdornment position="start">
                                            <Search color="action"/>
                                        </InputAdornment>
                                    ),
                                    'endAdornment': searchQuery && (
                                        <InputAdornment position="end">
                                            <IconButton
                                                size="small"
                                                onClick={handleClearSearch}
                                                title="Clear Search"
                                            >
                                                <Clear />
                                            </IconButton>
                                        </InputAdornment>
                                    )
                                }
                            }}
                        />
                    </Box>

                    <Button
                        startIcon={<Add />}
                        variant="contained"
                        onClick={() => navigate(`/topics/${topicID}/create-post`)}
                        sx={{ whiteSpace: 'nowrap' }}
                    >
                        Create Post 
                    </Button>
                </Box>

                {searchQuery && (
                    <Typography
                        variant="body2"
                        color="text.secondary"
                        sx={{ mb: 2 }}
                    >
                        Showing {filteredPosts.length} result{filteredPosts.length !== 1 ? 's' : ''} for "{searchQuery}"
                    </Typography>
                )}

                {filteredPosts.length === 0
                ? (
                    <Alert severity="info">
                        {
                            searchQuery
                                ? `No posts matching "${searchQuery}".`
                                : 'No posts available. Be the first to create one!'
                        }
                    </Alert>
                )
                : (
                    <Grid container spacing={3}>
                        {filteredPosts.map(
                            (post) => (
                                <Grid
                                    size={{
                                        xs: 12,
                                        sm: 6,
                                        md: 4,
                                    }}
                                    key={post.postID}
                                >
                                    <Card 
                                        variant="outlined"
                                        sx={{
                                            height: '100%',
                                            display: 'flex',
                                            flexDirection: 'column',
                                        }}
                                    > 
                                        <CardActionArea 
                                            onClick={() => navigate(`/posts/${post.postID}`)}
                                            sx={{
                                                height: '100%',
                                                display: 'flex',
                                                flexDirection: 'column',
                                                alignItems: 'stretch',
                                            }}    
                                        >
                                            <CardContent
                                                sx={{
                                                    flex: 1,
                                                    display: 'flex',
                                                    flexDirection: 'column',
                                                }}
                                            >
                                                <Typography
                                                    variant="h6"
                                                    component="h2"
                                                    fontWeight="bold"
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
                                                        flex: 1,
                                                    }}
                                                >
                                                    {post.content}
                                                </Typography>

                                                <Box
                                                    sx={{
                                                        display: 'flex',
                                                        gap: 0.5,
                                                        alignItems: 'center',
                                                        mt: 'auto',
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
                                </Grid>
                                
                            )
                        )}
                    </Grid>
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