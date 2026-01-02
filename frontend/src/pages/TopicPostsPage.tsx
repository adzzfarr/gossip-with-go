import { useNavigate, useParams } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "../hooks/redux";
import { useEffect } from "react";
import { fetchPostsByTopic } from "../features/posts/postsSlice";
import { Alert, Box, Button, Card, CardActionArea, CardContent, CircularProgress, Container, Typography } from "@mui/material";
import { ArrowBack } from "@mui/icons-material";

export default function TopicPostsPage() {
    const { topicID } = useParams<{ topicID: string }>();
    const dispatch = useAppDispatch(); 
    const navigate = useNavigate();
    const { posts, loading, error } = useAppSelector(state => state.posts);

    useEffect(() => {
        if (topicID) {
            dispatch(fetchPostsByTopic(parseInt(topicID)));
        }
    }, [dispatch, topicID]);

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

    if (error) {
        return (
            <Container sx={{ mt: 4 }}>
                <Alert severity="error">{error}</Alert>
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

    return (
        <Container
            sx={{ 
                mt: 4,
                mb: 4,
            }}
            maxWidth="lg"
        >
            <Box
                sx={{
                    mb: 3,
                    display: 'flex',
                    alignItems: 'center',
                    gap: 2,
                }}    
            >
                <Button
                    startIcon={<ArrowBack />}
                    onClick={() => navigate('/topics')}
                    variant="outlined"
                >
                    Back to Topics
                </Button>
                <Typography variant="h4" component="h1">
                    Posts
                </Typography>
            </Box>

            {posts.length === 0
                ? (<Alert severity="info">No posts available for this topic.</Alert>)
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
                                <Card key={post.postID}> 
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
                                            >
                                                {post.content.length > 200
                                                    ? post.content.substring(0, 200) + '...'
                                                    : post.content
                                                }
                                            </Typography>
                                            <Box
                                                sx={{
                                                    mt: 2,
                                                    display: 'flex',
                                                    gap: 2,
                                                }}
                                            >
                                                <Typography variant="caption" color="text.secondary">
                                                    By: {post.createdBy || 'Unknown'}
                                                </Typography>
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
                )
            }
        </Container>
    )
}