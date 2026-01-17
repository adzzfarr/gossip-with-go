import { useNavigate, useParams } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "../hooks/redux";
import { useEffect, useState } from "react";
import { fetchPostsByTopic, setPostsSortBy } from "../features/postsSlice";
import { Alert, Box, Button, Chip, CircularProgress, Container, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle, Divider, IconButton, InputAdornment, Paper, TextField, Typography } from "@mui/material";
import { Add, ArrowBack, Clear, Delete, Edit, Search } from "@mui/icons-material";
import ForumBreadcrumbs from "../components/Breadcrumbs";
import Username from "../components/Username";
import { deleteTopic, fetchTopicByID } from "../features/topicsSlice";
import PostsList from "../components/PostsList";
import SortDropdown from "../components/SortDropdown";

export default function TopicPostsPage() {
    const { topicID } = useParams<{ topicID: string }>();
    const dispatch = useAppDispatch(); 
    const navigate = useNavigate();

    const { posts, loading: postsLoading, error: postsError, sortBy } = useAppSelector(state => state.posts);
    const { topics, loading: topicLoading, submitting: topicSubmitting, submitError } = useAppSelector(state => state.topics);
    const { userID } = useAppSelector(state => state.auth);

    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
    
    const [searchQuery, setSearchQuery] = useState('');

    // Find current topic from topics list
    const topic = topics.find(t => t.topicID === parseInt(topicID || '0'));
    
    // Check if user is author of topic
    const isAuthor = topic && topic.createdBy === userID;

    useEffect(() => {
        if (topicID) {
            const id = parseInt(topicID);

            dispatch(fetchTopicByID(id));
            dispatch(fetchPostsByTopic({ topicID: id, sortBy }));
        }
    }, [dispatch, topicID, sortBy]);

    const handleSortChange = (newSort: string) => {
        dispatch(setPostsSortBy(newSort));
        if (topicID) {
            dispatch(fetchPostsByTopic({ topicID: parseInt(topicID), sortBy: newSort }));
        }
    }

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


    if (topicLoading || postsLoading) {
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
                    <Box
                        sx={{
                            display: 'flex',
                            alignItems: 'center',
                            gap: 2,
                        }}
                    >
                        <Typography variant="h4" component="h1">
                            {topic.title}
                        </Typography>
                        <Chip
                            label={`💬 ${posts.length} ${posts.length === 1 ? 'post' : 'posts'}`}
                            size="small"
                            variant="outlined"
                            sx={{
                                p: 0.5,
                            }}
                        />
                    </Box>

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
                {/* Search and Sort Controls + New Post */}
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

                        {/* Sort Options */}
                        <SortDropdown 
                            value={sortBy}
                            onChange={handleSortChange}
                            options={[
                                { value: 'hot', label: 'Hot' },
                                { value: 'newest', label: 'Newest' },
                                { value: 'oldest', label: 'Oldest' },
                                { value: 'most_voted', label: 'Most Voted' },
                            ]}
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
                : (<PostsList posts={filteredPosts} />)
                }
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