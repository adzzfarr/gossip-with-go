import { useNavigate, useParams } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "../hooks/redux";
import { useEffect, useState, type SyntheticEvent } from "react";
import { getUserComments, getUserPosts, getUserProfileByID } from "../features/profilesSlice";
import { Alert, Box, Card, CardActionArea, CardContent, CircularProgress, Container, Divider, Paper, Tab, Tabs, Typography } from "@mui/material";
import { Article, Comment as CommentIcon, Person } from "@mui/icons-material";

export default function ProfilePage() {
    const { userID } = useParams<{ userID: string }>();
    const dispatch = useAppDispatch();
    const navigate = useNavigate();

    const { profile, posts, comments, loading, error } = useAppSelector(state => state.profiles);
    const { userID: currentUserID } = useAppSelector(state => state.auth);

    const [activeTab, setActiveTab] = useState(0);

    // If no userID in URL, show current user's profile
    const profileUserID = userID ? parseInt(userID) : currentUserID;

    useEffect(() => {
        if (profileUserID) {
            dispatch(getUserProfileByID(profileUserID));
            dispatch(getUserPosts(profileUserID));
            dispatch(getUserComments(profileUserID));
        }
    }, [profileUserID, dispatch]);

    const handleTabChange = (_event: SyntheticEvent, newValue: number) => {
        setActiveTab(newValue);
    }

    if (loading && !profile) {
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
        )
    }

    if (error) {
        return (
            <Container sx={{ mt: 4 }}>
                <Alert severity="error">{error}</Alert>
            </Container>
        )
    }
    
    if (!profile) {
        return (
            <Container sx={{ mt: 4 }}>
                <Alert severity="info">User not found.</Alert>
            </Container>
        )
    }
        
    const isOwnProfile = profile.userID === currentUserID;

    return (
        <Container
            sx={{ 
                mt: 4,
                mb: 4,
            }}
            maxWidth="md"
        >
            {/* Profile Header */}
            <Paper
                elevation={3}
                sx={{
                    p: 3,
                    mb: 3,
                }}
            >
                <Box
                    sx={{
                        display: 'flex',
                        alignItems: 'center',
                        gap: 2,
                        mb: 2
                    }}
                >
                    <Person sx={{ fontSize: 60, color: 'primary.main' }} />
                    <Box>
                        <Typography variant="h4" component="h1">
                            {profile.username}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                            Member since {new Date(profile.createdAt).toLocaleDateString()}
                        </Typography>
                        {isOwnProfile && (
                            <Typography variant="caption" color="primary">
                                (Your Profile)
                            </Typography>
                        )}
                    </Box>
                </Box>

                <Divider sx={{ my: 2 }} />

                {/* Stats */}
                <Box
                    sx={{
                        display: 'flex',
                        gap: 4,
                    }}
                >
                    <Box>
                        <Typography variant="h6" color="primary">
                            {posts.length}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                            Posts
                        </Typography>
                    </Box>
                    <Box>
                        <Typography variant="h6" color="primary">
                            {comments.length}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                            Comments
                        </Typography>
                    </Box>
                </Box>
            </Paper>

            {/* Activity Tabs */}
            <Paper elevation={2}>
                <Tabs
                    value={activeTab}
                    onChange={handleTabChange}
                    sx={{
                        borderBottom: 1,
                        borderColor: 'divider'
                    }}
                >
                    <Tab icon={<Article />} label={`Posts (${posts.length})`} />  
                    <Tab icon={<CommentIcon />} label={`Comments (${comments.length})`} />
                </Tabs> 

                <Box sx={{ p: 2 }}>
                    {/* Posts Tab */}
                    {activeTab === 0 && (
                        <Box
                            sx={{
                                display: 'flex',
                                flexDirection: 'column',
                                gap: 2,
                            }}
                        >
                            {posts.length === 0
                                ? (<Alert severity="info">No posts by this user.</Alert>)
                                : posts.map(post => (
                                    <Card key={post.postID} variant="outlined">
                                        <CardActionArea onClick={() => navigate(`/posts/${post.postID}`)}>
                                            <CardContent>
                                                <Typography variant="h6" gutterBottom>
                                                    {post.title}
                                                </Typography>

                                                <Typography
                                                    variant="body2"
                                                    color="text.secondary"
                                                    sx ={{
                                                        display: '-webkit-box',
                                                        WebkitLineClamp: 2,
                                                        WebkitBoxOrient: 'vertical',
                                                        overflow: 'hidden',
                                                        textOverflow: 'ellipsis',
                                                    }}
                                                >
                                                    {post.content}
                                                </Typography>

                                                <Box sx={{ display: 'flex', gap: 1, mt: 1, alignItems: 'center' }}>
                                                    <Typography variant="caption" color="text.secondary">
                                                        in '{post.topicTitle}'
                                                    </Typography>

                                                    <Typography variant="caption" color="text.secondary">
                                                        •
                                                    </Typography>

                                                    <Typography 
                                                        variant="caption" 
                                                        color="text.secondary"
                                                    >
                                                        {new Date(post.createdAt).toLocaleString()}
                                                    </Typography>
                                                </Box>
                                            </CardContent>
                                        </CardActionArea>
                                    </Card>
                                )) 
                            }  
                        </Box>
                    )}

                    {/* Comments Tab */}
                    {activeTab === 1 && (
                        <Box
                            sx={{
                                display: 'flex',
                                flexDirection: 'column',
                                gap: 2,
                            }}
                        >
                            {comments.length === 0
                                ? (<Alert severity="info">No comments by this user.</Alert>)
                                : (
                                    comments.map(comment => (
                                        <Card key={comment.commentID} variant="outlined">
                                            <CardActionArea onClick={() => navigate(`/posts/${comment.postID}`)}>
                                                <CardContent>
                                                    <Typography variant="body1" sx={{ whiteSpace: 'pre-wrap' }}>
                                                        {comment.content}
                                                    </Typography>

                                                    <Box sx={{ display: 'flex', gap: 1, mt: 1, alignItems: 'center' }}>
                                                        <Typography variant="caption" color="text.secondary">
                                                            on '{comment.postTitle}'  
                                                        </Typography>
                                                        <Typography variant="caption" color="text.secondary">
                                                            •
                                                        </Typography>
                                                        <Typography 
                                                            variant="caption"
                                                            color="text.secondary"
                                                        >
                                                            {new Date(comment.createdAt).toLocaleString()}
                                                        </Typography>
                                                    </Box>
                                                </CardContent>
                                            </CardActionArea>
                                        </Card>
                                    ))
                                )
                            }
                        </Box>
                    )}
                </Box>
            </Paper>
        </Container>
    );
}