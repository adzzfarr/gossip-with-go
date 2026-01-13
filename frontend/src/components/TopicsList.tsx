import { Box, Card, CardActionArea, CardContent, Grid, Typography } from "@mui/material";
import type { Topic } from "../types";
import Username from "./Username";
import { useNavigate } from "react-router-dom";

function TopicCard({ topic }: { topic: Topic }) {
    const navigate = useNavigate();

    return (
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
    )
}

export default function TopicsList({ topics }: { topics: Topic[] }) {
    return (
        <Grid container spacing={3}>
            {topics.map(
                topic => (
                    <Grid
                        size={{
                            xs: 12,
                            sm: 6,
                            md: 4,
                        }}
                        key={topic.topicID}
                    >
                        <TopicCard topic={topic} />
                    </Grid>

                )
            )}
        </Grid>
    );
}