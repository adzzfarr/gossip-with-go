import { useNavigate } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "../hooks/redux";
import { useEffect, useState } from "react";
import { fetchTopics, setTopicsCurrentPage, setTopicsSearchQuery, setTopicsSortBy } from "../features/topicsSlice";
import { Alert, Box, Button, CircularProgress, Container, IconButton, InputAdornment, TextField, Typography } from "@mui/material";
import { Add, Clear, Search } from "@mui/icons-material";
import TopicsList from "../components/TopicsList";
import SortDropdown from "../components/SortDropdown";
import Pagination from "../components/Pagination";

export default function TopicsPage() {
    const dispatch = useAppDispatch();
    const navigate = useNavigate();
    const { topics, loading, error, sortBy, searchQuery, pagination, currentPage } = useAppSelector(state => state.topics);
    const { isAuthenticated } = useAppSelector(state => state.auth);

    const [localSearchQuery, setLocalSearchQuery] = useState(searchQuery);

    const handleClearSearch = () => {
        setLocalSearchQuery('');
        dispatch(setTopicsSearchQuery(''));
    }

    const handleSearchChange = (value: string) => {
        setLocalSearchQuery(value);
    }

    const handleSortChange = (newSort: string) => {
        dispatch(setTopicsSortBy(newSort));
    }

    const handlePageChange = (newPage: number) => {
        dispatch(setTopicsCurrentPage(newPage));
        
        // Scroll to top on page change
        window.scrollTo({ top: 0, behavior: 'smooth' });
    }

    // Debounce search input
    useEffect(() => {
        const delayDebounce = setTimeout(() => {
            if (localSearchQuery !== searchQuery) {
                dispatch(setTopicsSearchQuery(localSearchQuery));
            }
            
        }, 500);

        return () => clearTimeout(delayDebounce);
    }, [localSearchQuery, searchQuery, dispatch]);

    useEffect(() => {
        dispatch(fetchTopics({
            sortBy,
            page: currentPage,
            search: searchQuery,
        }));
    }, [dispatch, sortBy, searchQuery, currentPage]);

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
            {/* Header */}
            <Box
                sx={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    mb: 3,
                }}
            >
                <Typography 
                    variant="h4" 
                    component="h1" 
                    gutterBottom
                >
                    Discussion Topics
                </Typography>

                {isAuthenticated && (
                    <Button
                        variant="contained"
                        startIcon={<Add />}
                        onClick={() => navigate('/topics/create')}
                        sx={{ whiteSpace: 'nowrap' }}

                    >
                        New Topic
                    </Button>
                )}
            </Box>

            {/* Search and Sort Controls */}
            <Box
                sx={{
                    display: 'flex',
                    gap: 2,
                    mb: 3,
                }}
            >
                <TextField 
                    size="small"
                    fullWidth
                    placeholder="Search Topics..."
                    value={localSearchQuery}
                    onChange={(e) => handleSearchChange(e.target.value)}
                    slotProps={{
                        input: {
                            'startAdornment': (
                                <InputAdornment position="start">
                                    <Search color="action"/>
                                </InputAdornment>
                            ),
                            'endAdornment': localSearchQuery && (
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

                <SortDropdown 
                    value={sortBy}
                    onChange={handleSortChange}
                    options={[
                        { value: 'newest', label: 'Newest' },
                        { value: 'oldest', label: 'Oldest' },
                        { value: 'most_posts', label: 'Most Posts' },
                    ]}
                />
            </Box>
            
            {localSearchQuery && pagination && (
                <Typography
                    variant="body2"
                    color="text.secondary"
                    sx={{ mb: 2 }}
                >
                    Showing {pagination.totalItems} result{pagination.totalItems !== 1 ? 's' : ''} for "{localSearchQuery}"
                </Typography>
            )}

            {/* Topics List */}
            {topics.length === 0 
                ? (<Alert severity="info">
                    {
                        localSearchQuery
                            ? `No topics matching "${localSearchQuery}".`
                            : 'No topics available. Be the first to create one!'
                    }
                </Alert>) 
                : (<TopicsList topics={topics} />)
            }

            {/* Pagination Controls */}
            {pagination && (
                <Pagination 
                    currentPage={pagination.currentPage}
                    totalPages={pagination.totalPages}
                    onPageChange={handlePageChange}
                    disabled={loading}
                />
            )}
        </Container>
    )
}