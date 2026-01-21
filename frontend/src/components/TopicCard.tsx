import { Box, Card, CardActionArea, CardContent, Typography } from "@mui/material";
import Username from "./Username";
import { useNavigate } from "react-router-dom";
import type { Topic } from "../types";

export function TopicCard({ topic }: { topic: Topic }) {
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
                                Started by
                            </Typography>
                            <Username
                                username={topic.username}
                                userID={topic.createdBy}
                                variant="caption"
                                color="text.secondary"
                            />
                            <Typography variant="caption" color="text.secondary">•</Typography>
                            <Typography variant="caption" color="text.secondary">
                                {new Date(topic.createdAt).toLocaleDateString()}
                            </Typography>
                        </Box>
                    </Box>

                    <Typography 
                        variant="caption" 
                        color="text.secondary"
                    >
                        {topic.postCount} {topic.postCount === 1 ? 'post' : 'posts'}
                    </Typography>
                </CardContent>
            </CardActionArea>
        </Card>
    )
}