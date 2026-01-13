import { Grid, Card, CardActionArea, CardContent, Typography, Box } from "@mui/material";
import type { Post } from "../types";
import Username from "./Username";
import { useNavigate } from "react-router-dom";
import VoteButtons from "./VoteButtons";

function PostCard({ post }: { post: Post }) {
    const navigate = useNavigate();

    return (
        <Card 
            variant="outlined"
            sx={{
                height: '100%',
                display: 'flex',
                flexDirection: 'row',
            }}
        > 
            {/* Vote Buttons */}
            <Box
                sx={{
                    display: 'flex',
                    alignItems: 'flex-start',
                    p: 1.5,
                    pr: 0,
                    bgcolor: 'background.default',
                }}
                onClick={(e) => e.stopPropagation()}
            >
                <VoteButtons
                    postID={post.postID}
                    initialVoteCount={post.voteCount || 0}
                    initialUserVote={post.userVote}
                    size="small"
                    orientation="vertical"
                />
            </Box>

            {/* Post Content */}
            <CardActionArea 
                onClick={() => navigate(`/posts/${post.postID}`)}
                sx={{
                    display: 'flex',
                    flex: 1,
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
    );
}

export default function PostsList({ posts }: { posts: Post[] }) {
    return (
        <Grid container spacing={3}>
            {posts.map(
                (post) => (
                    <Grid
                        size={{
                            xs: 12,
                            sm: 6,
                            md: 4,
                        }}
                        key={post.postID}
                    >
                            <PostCard post={post} />
                    </Grid>   
                )
            )}
        </Grid>
    );
}