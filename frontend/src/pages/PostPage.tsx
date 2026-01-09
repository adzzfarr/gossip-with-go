import type { Comment as CommentType } from "../types/index";
import { useNavigate, useParams } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "../hooks/redux";
import { deletePost, fetchPostByID } from "../features/postsSlice";
import { clearCommentsError, createComment, deleteComment, fetchCommentsByPostID, updateComment } from "../features/commentsSlice";
import { useEffect, useState } from "react";
import { Alert, Box, Button, Card, CardContent, CircularProgress, Container, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle, Divider, Paper, TextField, Typography } from "@mui/material";
import { ArrowBack, Delete, Edit } from "@mui/icons-material";
import ForumBreadcrumbs from "../components/Breadcrumbs";
import Username from "../components/Username";

export default function PostPage() {
    const { postID } = useParams<{ postID: string}>();
    const dispatch = useAppDispatch();
    const navigate = useNavigate();

    const { currentPost, loading: postLoading, error: postError, submitting: postSubmitting } = useAppSelector(state => state.posts);
    const { comments, loading: commentsLoading, error: commentsError, submitting: commentSubmitting, submitError: commentSubmitError } = useAppSelector(state => state.comments);
    const { userID } = useAppSelector(state => state.auth);

    useEffect(() => {
        if (postID) {
            dispatch(fetchPostByID(parseInt(postID)));
            dispatch(fetchCommentsByPostID(parseInt(postID)));
        }
    }, [postID, dispatch]);

    // Check if user is author 
    const isAuthor = currentPost && currentPost.createdBy === userID;

    // For adding new comment
    const [commentContent, setCommentContent] = useState('');    

    const handleCommentChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setCommentContent(e.target.value);
        if (commentSubmitError) {
            dispatch(clearCommentsError());
        }
    }
    
    const handleSubmitComment = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!postID || !commentContent.trim()) return;

        const result = await dispatch(
            createComment({
                postID: parseInt(postID),
                content: commentContent.trim(),
            })
        );

        if (createComment.fulfilled.match(result)) {
            // Clear comment input on successful submission
            setCommentContent('');
        }
    }

    // For editing existing comments
    const [commentToEditID, setCommentToEditID] = useState<number | null>(null);
    const [editedCommentContent, setEditedCommentContent] = useState('');

    const handleEditCommentClicked = (comment: CommentType) => {
        setCommentToEditID(comment.commentID);
        setEditedCommentContent(comment.content);
    }

    const handleCancelEdit = () => {
        setCommentToEditID(null);
        setEditedCommentContent('');
    }

    const handleSaveEditedComment = async (commentID: number) => {
        if (!editedCommentContent.trim()) return;

        const result = await dispatch(
            updateComment({
                commentID,
                content: editedCommentContent.trim(),
            })
        );

        if (updateComment.fulfilled.match(result)) {
            // Clear edit state on successful update
            setCommentToEditID(null);
            setEditedCommentContent('');
        }
    }

    // For deleting existing comments
    const [commentToDelete, setCommentToDelete] = useState<number | null>(null);
    const [deleteCommentDialogOpen, setDeleteCommentDialogOpen] = useState(false);

    const handleDeleteCommentClicked = (commentID: number) => {
        setCommentToDelete(commentID);
        setDeleteCommentDialogOpen(true);
    }

    const handleDeleteComment = async () => {
        if (!commentToDelete) return;

        await dispatch(deleteComment(commentToDelete));
        setDeleteCommentDialogOpen(false);
        setCommentToDelete(null);
    }

    // For deleting post
    const [deletePostDialogOpen, setDeletePostDialogOpen] = useState(false); 

    const handleDeletePost = async () => {
        if (!postID) return;

        const result = await dispatch(deletePost(parseInt(postID)));

        if (deletePost.fulfilled.match(result)) {
            if (currentPost?.topicID) {
                // Redirect to topic posts page after deletion
                navigate(`/topics/${currentPost.topicID}`);
            } else {
                navigate('/topics');
            }
        }

        setDeletePostDialogOpen(false);
    }

    if (postLoading || commentsLoading) {
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

    if (postError) {
        return (
            <Container sx={{ mt: 4 }}>
                <Alert severity="error">{postError}</Alert>
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

    if (!currentPost) {
        return (
            <Container>
                <Alert severity="info">Post not found.</Alert>
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
            maxWidth="md"
            sx={{
                mt: 4,
                mb: 4,
            }}
        >
            <ForumBreadcrumbs />

            <Button
                startIcon={<ArrowBack />}
                onClick={() => navigate(-1)}
                sx={{ mb: 3 }}
                variant="outlined"
            >
                Back
            </Button>

            {/* Post Content */}
            <Paper
                elevation={2}
                sx={{
                    p: 3,
                    mb: 4,
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
                    <Typography 
                        variant="h4" 
                        component="h1" 
                        gutterBottom
                    >
                        {currentPost.title}
                    </Typography>

                    {/* Edit and Delete Buttons for author */}
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
                                onClick={() => navigate(`/posts/${currentPost.postID}/edit`)}
                                size="small"
                            >
                                Edit
                            </Button>

                            <Button
                                startIcon={<Delete />}
                                variant="outlined"
                                color="error"
                                onClick={() => setDeletePostDialogOpen(true)}
                                disabled={postSubmitting}
                                size="small"
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
                        gap: 0.5,
                    }}
                >
                    <Username
                        username={currentPost.username || 'Unknown'}
                        userID={currentPost.createdBy}
                        variant="body2"
                        color="text.secondary"
                    />
                    
                    <Typography variant="body2" color="text.secondary">
                        {new Date(currentPost.createdAt).toLocaleString()}
                    </Typography>
                </Box>

                <Divider sx={{ my: 2 }} />

                <Typography variant="body1" sx={{ whiteSpace: 'pre-wrap' }}>
                    {currentPost.content}
                </Typography>
            </Paper>

            {/* Comments Section */}
            <Box>
                <Typography 
                    variant="h5" 
                    component="h2" 
                    gutterBottom
                >
                    Comments
                </Typography>

                {/* Add Comment Form */}
                <Paper
                    elevation={1}
                    sx={{ 
                        p: 2, 
                        mb: 3 
                    }}
                >
                    <form onSubmit={handleSubmitComment}>
                        <TextField
                            fullWidth
                            multiline
                            minRows={3}
                            placeholder="Write a comment..."
                            value={commentContent}
                            onChange={handleCommentChange}
                            disabled={commentSubmitting}
                            sx={{ mb: 2 }}
                        />
                        
                        {commentSubmitError && (
                            <Alert severity="error" sx={{ mb: 2 }}>{commentSubmitError}</Alert>
                        )}

                        <Button
                            type="submit"
                            variant="contained"
                            disabled={commentSubmitting || !commentContent.trim()}
                        >
                            {commentSubmitting ? <CircularProgress size={24} /> : 'Post Comment'}
                        </Button>
                    </form>
                </Paper>

                {commentsError && (
                    <Alert severity="error" sx={{ mb: 2 }}>
                        {commentsError}
                    </Alert>
                )}

                {comments.length === 0
                    ? (<Alert severity="info">No comments yet. Be the first to comment!</Alert>)
                    : (
                        <Box
                            sx={{
                                display: 'flex',
                                flexDirection: 'column',
                                gap: 2,
                            }}
                        >
                            {comments.map((comment) => {
                                const isCommentAuthor = comment.createdBy === userID;
                                const isEditing = commentToEditID === comment.commentID;

                                return (
                                    <Card key={comment.commentID} variant="outlined">
                                        <CardContent>
                                            <Box
                                                sx={{
                                                    display: 'flex',
                                                    justifyContent: 'space-between',
                                                    alignItems: 'flex-start',
                                                    mb: 1,
                                                }}
                                            >
                                                <Box
                                                    sx={{
                                                        display: 'flex',
                                                        gap: 2,
                                                        alignItems: 'center',
                                                    }}
                                                >
                                                    <Username
                                                        username={comment.username || 'Unknown'}
                                                        userID={comment.createdBy}
                                                        variant="body2"
                                                        fontWeight="bold"
                                                    />

                                                    <Typography variant="body2" color="text.secondary">
                                                        {new Date(comment.createdAt).toLocaleString()}
                                                    </Typography>
                                                </Box>

                                                {/* Edit and Delete Buttons for comment author */}
                                                {isCommentAuthor && !isEditing && (
                                                    <Box
                                                        sx={{
                                                            display: 'flex',
                                                            gap: 1,
                                                        }}
                                                    >
                                                        <Button
                                                            size="small"
                                                            startIcon={<Edit />}
                                                            onClick={() => handleEditCommentClicked(comment)}
                                                        >
                                                            Edit
                                                        </Button>
                                                        <Button
                                                            size="small"
                                                            color="error"
                                                            startIcon={<Delete />}
                                                            onClick={() => handleDeleteCommentClicked(comment.commentID)}
                                                            disabled={commentSubmitting}
                                                        >
                                                            Delete
                                                        </Button>
                                                    </Box>
                                                )}
                                            </Box>

                                            {isEditing 
                                                ? (
                                                    <Box>
                                                        <TextField
                                                            fullWidth
                                                            multiline
                                                            minRows={2}
                                                            value={editedCommentContent}
                                                            onChange={(e) => setEditedCommentContent(e.target.value)}
                                                            disabled={commentSubmitting}
                                                            sx={{ mb: 1 }}
                                                        />
                                                        <Box
                                                            sx={{
                                                                display: 'flex',
                                                                gap: 1,
                                                            }}
                                                        >
                                                            <Button
                                                                size="small"
                                                                variant="contained"
                                                                onClick={() => handleSaveEditedComment(comment.commentID)}
                                                                disabled={commentSubmitting || !editedCommentContent.trim()}
                                                            >
                                                                {commentSubmitting ? <CircularProgress size={24} /> : 'Save'}
                                                            </Button>
                                                            <Button
                                                                size="small"
                                                                variant="outlined"
                                                                onClick={handleCancelEdit}
                                                                disabled={commentSubmitting}
                                                            >
                                                                Cancel
                                                            </Button>
                                                        </Box>
                                                    </Box>
                                                )
                                                : (
                                                    <Typography variant="body1" sx={{ whiteSpace: 'pre-wrap' }}>
                                                        {comment.content}
                                                    </Typography>
                                                )
                                            }
                                        </CardContent>
                                    </Card>
                                );
                            })}
                        </Box>
                    )
                }
            </Box>

            {/* Delete Comment Confirmation Dialog */}
            <Dialog
                open={deleteCommentDialogOpen}
                onClose={() => setDeleteCommentDialogOpen(false)}
            >
                <DialogTitle>Delete Comment?</DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        Are you sure you want to delete this comment? This action cannot be undone.
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setDeleteCommentDialogOpen(false)}>
                        Cancel
                    </Button>
                    <Button
                        onClick={handleDeleteComment}
                        color="error"
                        disabled={commentSubmitting}
                    >
                        {commentSubmitting ? <CircularProgress size={24} /> : 'Delete'}
                    </Button>
                </DialogActions>
            </Dialog>

            {/* Delete Post Confirmation Dialog */}
            <Dialog
                open={deletePostDialogOpen}
                onClose={() => setDeletePostDialogOpen(false)}
            >
                <DialogTitle>Delete Post?</DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        Are you sure you want to delete this post? This action cannot be undone.
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button 
                        onClick={() => setDeletePostDialogOpen(false)} 
                    >
                        Cancel
                    </Button>
                    <Button
                        onClick={handleDeletePost}
                        color="error"
                        disabled={postSubmitting}
                    >
                        {postSubmitting ? <CircularProgress size={24} /> : 'Delete'}
                    </Button>
                </DialogActions>
            </Dialog>
        </Container>
    );
}