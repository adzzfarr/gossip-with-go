import { useNavigate } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "../hooks/redux";
import { useEffect, useState } from "react";
import { fetchTopics } from "../features/topicsSlice";
import { Alert, Box, Button, CircularProgress, Container, IconButton, InputAdornment, TextField, Typography } from "@mui/material";
import { Add, Clear, Search } from "@mui/icons-material";
import TopicsList from "../components/TopicsList";

export default function TopicsPage() {
    const dispatch = useAppDispatch();
    const navigate = useNavigate();
    const { topics, loading, error } = useAppSelector(state => state.topics);

    const [searchQuery, setSearchQuery] = useState('');

    const handleClearSearch = () => {
        setSearchQuery('');
    }

    useEffect(() => {
        dispatch(fetchTopics());
    }, [dispatch]);

    const filteredTopics = topics.filter(
        topic => {
            const query = searchQuery.trim().toLowerCase();
            return topic.title.toLowerCase().includes(query);
        }
    );

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
            {filteredTopics.length === 0 
                ? (<Alert severity="info">
                    {
                        searchQuery
                            ? `No topics matching "${searchQuery}".`
                            : 'No topics available. Be the first to create one!'
                    }
                </Alert>) 
                : (<TopicsList topics={filteredTopics} />)
            }
        </Container>
    )
}