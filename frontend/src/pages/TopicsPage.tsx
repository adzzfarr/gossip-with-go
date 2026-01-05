import { useNavigate } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "../hooks/redux";
import { useEffect } from "react";
import { fetchTopics } from "../features/topics/topicsSlice";
import { Alert, Box, Button, Card, CardActionArea, CardContent, CircularProgress, Container, Grid, Typography } from "@mui/material";
import { Add } from "@mui/icons-material";

export default function TopicsPage() {
    const dispatch = useAppDispatch();
    const navigate = useNavigate();
    const { topics, loading, error } = useAppSelector(state => state.topics);

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

                <Button
                    variant="contained"
                    startIcon={<Add />}
                    onClick={() => navigate('/topics/create')}
                >
                    New Topic
                </Button>
            </Box>
            
            

            {topics.length === 0 
                ? (<Alert severity="info">No topics available.</Alert>) 
                : (
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
                )
            }
        </Container>
    )
}