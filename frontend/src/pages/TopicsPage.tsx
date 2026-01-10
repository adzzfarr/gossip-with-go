import { useNavigate } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "../hooks/redux";
import { useEffect, useState } from "react";
import { fetchTopics } from "../features/topicsSlice";
import { Alert, Box, Button, Card, CardActionArea, CardContent, CircularProgress, Container, Grid, IconButton, InputAdornment, TextField, Typography } from "@mui/material";
import { Add, Clear, Search } from "@mui/icons-material";
import Username from "../components/Username";

export default function TopicsPage() {
    const dispatch = useAppDispatch();
    const navigate = useNavigate();
    const { topics, loading, error } = useAppSelector(state => state.topics);

    const [searchQuery, setSearchQuery] = useState('');



    // Filter topics based on search query
    const filteredTopics = topics.filter(
        topic => {
            const query = searchQuery.trim().toLowerCase();
            return (
                topic.title.toLowerCase().includes(query) ||
                topic.description.toLowerCase().includes(query) ||
                topic.username.toLowerCase().includes(query)
            );
        }
    );

    const handleClearSearch = () => {
        setSearchQuery('');
    }

    useEffect(() => {
        dispatch(fetchTopics());
    }, [dispatch]);

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
            <Container sx={{ mt: 4}}>
                <Alert severity="error">{error}</Alert>
            </Container>
        )
    }
        
    return (
        <Container
            sx={{ 
                mt: 4,
                mb: 4,
            }}
            maxWidth="lg"
        >
            {/* Header (Search Bar + New Topic) */}
            <Box
                sx={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    mb: 3,
                    gap: 2,
                    flexWrap: 'wrap',
                }}
            >
                <Typography 
                    variant="h4" 
                    component="h1" 
                    gutterBottom
                >
                    Discussion Topics
                </Typography>

                <Box
                    sx={{
                        display: 'flex',
                        gap: 2,
                        alignItems: 'center',
                        flex: 1,
                        maxWidth: 600,
                    }}
                >
                    {/* Search Bar */}
                    <TextField 
                        size="small"
                        fullWidth
                        placeholder="Search Topics..."
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

                    <Button
                        variant="contained"
                        startIcon={<Add />}
                        onClick={() => navigate('/topics/create')}
                        sx={{ whiteSpace: 'nowrap' }}
                    >
                        New Topic
                    </Button>
                </Box>
            </Box>
            
            {searchQuery && (
                <Typography
                    variant="body2"
                    color="text.secondary"
                    sx={{ mb: 2 }}
                >
                    Showing {filteredTopics.length} result{filteredTopics.length !== 1 ? 's' : ''} for "{searchQuery}"
                </Typography>
            )}

            {/* Topics List */}
            {topics.length === 0 
                ? (<Alert severity="info">
                    {
                        searchQuery
                            ? `No topics matching "${searchQuery}".`
                            : 'No topics available. Be the first to create one!'
                    }
                </Alert>) 
                : (
                    <Grid container spacing={3}>
                        {filteredTopics.map(
                            topic => (
                                <Grid
                                    size={{
                                        xs: 12,
                                        sm: 6,
                                        md: 4,
                                    }}
                                    key={topic.topicID}
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
                                            onClick={() => navigate(`/topics/${topic.topicID}`)}
                                            sx={{
                                                height: '100%',
                                                display: 'flex',
                                                flexDirection: 'column',
                                                alignItems: 'stretch',
                                            }}
                                        >
                                            <CardContent 
                                                sx={{
                                                    display: 'flex', 
                                                    flexDirection: 'column',
                                                    flex: 1,
                                                }}
                                            >
                                                <Typography
                                                        variant="h6"
                                                        component="h2"
                                                        fontWeight="bold"
                                                        gutterBottom
                                                    >
                                                        {topic.title}
                                                    </Typography>
                                                    
                                                <Typography 
                                                    variant="body2" 
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
                                                    {topic.description || 'No description provided.'}
                                                </Typography>

                                                <Box
                                                    sx={{
                                                        display: 'flex',
                                                        flexDirection: 'column',
                                                        gap: 0.5,
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
                                                            By
                                                        </Typography>
                                                        <Username
                                                            username={topic.username}
                                                            userID={topic.createdBy}
                                                            variant="caption"
                                                            color="text.secondary"
                                                        />
                                                    </Box>
                                                    <Typography variant="caption" color="text.secondary">
                                                        {new Date(topic.createdAt).toLocaleDateString()}
                                                    </Typography>
                                                </Box>
                                            </CardContent>
                                        </CardActionArea>
                                    </Card>
                                </Grid>
                                
                            )
                        )}
                    </Grid>
                    /*
                    <Grid container spacing={3}>
                        {topics.map(
                            (topic) => (
                                <Grid 
                                    size={{
                                        xs: 12,
                                        sm: 6,
                                        md: 4,
                                    }}
                                    key={topic.topicID}
                                >
                                    <Card>
                                        <CardActionArea onClick={() => navigate(`/topics/${topic.topicID}`)}>
                                            <CardContent>
                                                <Typography
                                                    variant="h5"
                                                    component="h2"
                                                    gutterBottom
                                                >
                                                    {topic.title}
                                                </Typography>
                                                {topic.description && (
                                                    <Typography variant="body2" color="text.secondary">
                                                        {topic.description}
                                                    </Typography>
                                                )}
                                            </CardContent>
                                        </CardActionArea>
                                    </Card>
                                </Grid>
                            )
                        )}
                    </Grid> 
                    */
                )
            }
        </Container>
    )
}