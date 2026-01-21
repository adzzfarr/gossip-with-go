import { useNavigate, useParams } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "../hooks/redux";
import { deletePost, fetchPostByID } from "../features/postsSlice";
import { createComment, createCommentReply, deleteComment, fetchCommentsByPostID, setCommentsCurrentPage, setCommentsSearchQuery, setCommentsSortBy, updateComment } from "../features/commentsSlice";
import { useEffect, useState } from "react";
import { Alert, Box, Button, CircularProgress, Container, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle, Divider, IconButton, InputAdornment, Paper, TextField, Typography } from "@mui/material";
import { ArrowBack, Clear, Delete, Edit, Search } from "@mui/icons-material";
import ForumBreadcrumbs from "../components/Breadcrumbs";
import Username from "../components/Username";
import VoteButtons from "../components/VoteButtons";
import SortDropdown from "../components/SortDropdown";
import Pagination from "../components/Pagination";
import CommentsList from "../components/CommentsList";
import CommentForm from "../components/CommentForm";
import type { Comment } from "../types";

export default function PostPage() {
    const { postID } = useParams<{ postID: string}>();
    const dispatch = useAppDispatch();
    const navigate = useNavigate();

    const { currentPost, loading: postLoading, error: postError, submitting: postSubmitting } = useAppSelector(state => state.posts);
    const { comments, loading: commentsLoading, error: commentsError, submitting: commentSubmitting, submitError: commentSubmitError, sortBy, searchQuery, pagination, currentPage } = useAppSelector(state => state.comments);
    const { userID, isAuthenticated } = useAppSelector(state => state.auth);

    const [localSearchQuery, setLocalSearchQuery] = useState(searchQuery);

    const handleClearSearch = () => {
        setLocalSearchQuery('');
        dispatch(setCommentsSearchQuery(''));
    }

    const handleSearchChange = (value: string) => {
        setLocalSearchQuery(value);
    }

    // Debounce search input
    useEffect(() => {
        const delayDebounce = setTimeout(() => {
            if (localSearchQuery !== searchQuery) {
                dispatch(setCommentsSearchQuery(localSearchQuery));
            }
            
        }, 500);

        return () => clearTimeout(delayDebounce);
    }, [localSearchQuery, searchQuery, dispatch]);

    // Fetch post on mount
    useEffect(() => {
        if (postID) {
            dispatch(fetchPostByID(parseInt(postID)));
        }
    }, [postID, dispatch]);

    // Fetch comments when postID, sortBy, searchQuery, or currentPage changes
    useEffect(() => {
        if (postID) {
            dispatch(fetchCommentsByPostID({
                postID: parseInt(postID),
                sortBy,
                search: searchQuery,
                page: currentPage,
            }));
        }
    }, [postID, sortBy, searchQuery, currentPage, dispatch]);

    const handleCommentSortChange = (newSort: string) => {
        dispatch(setCommentsSortBy(newSort));
    }

    const handlePageChange = (newPage: number) => {
        dispatch(setCommentsCurrentPage(newPage));
        
        // Scroll to top of comments section on page change
        const commentsSection = document.getElementById('comments-section');
        if (commentsSection) {
            commentsSection.scrollIntoView({ behavior: 'smooth' });
        }
    }

    // Check if user is author 
    const isAuthor = currentPost && currentPost.createdBy === userID;  

    // For adding new comments
    const handleCreateComment = async (content: string) => {
        if (!postID || !content.trim()) return;

        const result = await dispatch(
            createComment({
                postID: parseInt(postID),
                content: content.trim(),
            })
        );

        if (createComment.fulfilled.match(result)) {
            // Refresh comments to show new comment
            dispatch(fetchCommentsByPostID({
                postID: parseInt(postID),
                sortBy,
                search: searchQuery,
                page: currentPage,
            }));
        }
    }

    // For editing existing comments
    const handleEdit = async (commentID: number, content: string) => {
        const result = await dispatch(
            updateComment({
                commentID,
                content,
            })
        );

        if (updateComment.fulfilled.match(result)) {
            // Refresh to show edit
            dispatch(fetchCommentsByPostID({
                postID: parseInt(postID || '0'),
                sortBy,
                search: searchQuery,
                page: currentPage,
            }));
        }
    }

    // For deleting existing comments
    const [commentToDelete, setCommentToDelete] = useState<number | null>(null);
    const [deleteCommentDialogOpen, setDeleteCommentDialogOpen] = useState(false);

    const countReplies = (comment: Comment): number => {
        if (!comment.replies || comment.replies.length === 0) return 0;

        return comment.replies.length + comment.replies.reduce(
            (sum, reply) => sum + countReplies(reply), 0,
        )
    }

    const findCommentByID = (commentsList: Comment[], commentID: number): Comment | null => {
        for (const comment of commentsList) {
            if (comment.commentID === commentID) return comment;

            if (comment.replies) {
                const foundInReplies = findCommentByID(comment.replies, commentID);
                if (foundInReplies) return foundInReplies;
            }
        }
        return null;
    }

    const handleDeleteCommentClicked = (commentID: number) => {
        setCommentToDelete(commentID);
        setDeleteCommentDialogOpen(true);
    }

    const handleConfirmDeleteComment = async () => {
        if (!commentToDelete || !postID) return;

        const result = await dispatch(
            deleteComment(commentToDelete)
        );

        if (deleteComment.fulfilled.match(result)) {
            // Refresh comments to show deletion
            dispatch(fetchCommentsByPostID({
                postID: parseInt(postID),
                sortBy,
                search: searchQuery,
                page: currentPage,
            }));
        }

        setDeleteCommentDialogOpen(false);
        setCommentToDelete(null);
    }

    // For replying to comments
    const handleReply = async (parentCommentID: number, content: string) => {
        if (!postID || !content.trim()) return;

        const result = await dispatch(
            createCommentReply({
                postID: parseInt(postID),
                parentCommentID,
                content: content.trim(),
            })
        );

        if (createCommentReply.fulfilled.match(result)) {
            // Refresh comments to show new comment
            dispatch(fetchCommentsByPostID({
                postID: parseInt(postID),
                sortBy,
                search: searchQuery,
                page: currentPage,
            }));
        };
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
                <Box>
                    {/* Header */}
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

                    {/* Metadata */}
                    <Box
                        sx={{
                            display: 'flex',
                            alignItems: 'center',
                            gap: 0.5,
                            mb: 2
                        }}
                    >
                        <Typography variant="body2" color="text.secondary">
                            Posted by
                        </Typography>
                        <Username
                            username={currentPost.username || 'Unknown'}
                            userID={currentPost.createdBy}
                            variant="body2"
                            color="text.secondary"
                        />
                        <Typography variant="body2" color="text.secondary">•</Typography>
                        <Typography variant="body2" color="text.secondary">
                            {new Date(currentPost.createdAt).toLocaleString()}
                        </Typography>
                    </Box>

                    <Divider sx={{ my: 2 }} />

                    {/* Post Content */}
                    <Typography variant="body1" sx={{ whiteSpace: 'pre-wrap' }}>
                        {currentPost.content}
                    </Typography>

                    <Divider sx={{ my: 2 }} />

                    {/* Vote Buttons */}
                    <Box
                        sx={{
                            display: 'flex',
                            alignItems: 'center',
                            gap: 1,
                        }}
                    >
                        <VoteButtons 
                            postID={currentPost.postID}
                            initialVoteCount={currentPost.voteCount || 0}
                            initialUserVote={currentPost.userVote}
                            orientation="horizontal"
                            size="small"
                        />
                    </Box>
                </Box>
            </Paper>

            {/* Comments Section */}
            <Box>
                {/* Comments Header */}
                <Typography variant="h5" component="h2">
                    Comments ({pagination?.totalItems || comments.length})
                </Typography>

                {/* Search and Sort Controls */}
                <Box
                    sx={{
                        display: 'flex',
                        alignItems: 'center',
                        gap: 2,
                        mb: 3,
                    }}
                >
                    <TextField
                        size="small"
                        fullWidth
                        placeholder="Search Comments..."
                        value={localSearchQuery}
                        onChange={(e) => handleSearchChange(e.target.value)}
                        slotProps={{
                            input: {
                                'startAdornment': (
                                    <InputAdornment position="start">
                                        <Search color="action" />
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
                        onChange={handleCommentSortChange}
                        options={[
                            { value: 'hot', label: 'Hot' },
                            { value: 'most_voted', label: 'Most Voted' },
                            { value: 'newest', label: 'Newest' },
                            { value: 'oldest', label: 'Oldest' },
                        ]}
                    />
                </Box>

                {searchQuery && pagination && (
                    <Typography
                        variant="body2"
                        color="text.secondary"
                        sx={{ mb: 2 }}
                    >
                        Showing {pagination.totalItems} result{pagination.totalItems !== 1 ? 's' : ''} for "{searchQuery}"
                    </Typography>
                )}

                <Divider sx={{ mb: 3 }} />

                {commentsError && (
                    <Alert severity="error" sx={{ mb: 2 }}>
                        {commentsError}
                    </Alert>
                )}

                {/* Comments List */}
                <CommentsList 
                    comments={comments}
                    userID={userID ?? undefined}
                    isSubmitting={commentSubmitting}
                    onEdit={handleEdit}
                    onDelete={handleDeleteCommentClicked}
                    onReply={handleReply}
                />

                {/* Add Comment Form */}
                {isAuthenticated && (
                    <Box>
                        <CommentForm 
                            onSubmit={handleCreateComment}
                            isSubmitting={commentSubmitting}
                            error={commentSubmitError}
                        />
                    </Box>
                )}
            </Box>

            {/* Pagination Controls */}
            {pagination && (
                <Pagination 
                    currentPage={pagination.currentPage}
                    totalPages={pagination.totalPages}
                    onPageChange={handlePageChange}
                    disabled={commentsLoading}
                />
            )}

            {/* Delete Comment Confirmation Dialog */}
            <Dialog
                open={deleteCommentDialogOpen}
                onClose={() => setDeleteCommentDialogOpen(false)}
            >
                <DialogTitle>Delete Comment?</DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        {commentToDelete && (
                            () => {
                                const comment = findCommentByID(comments, commentToDelete);
                                const replyCount = comment ? countReplies(comment) : 0;
                                return (
                                    <>
                                        replyCount == 0 
                                            ? <>Are you sure you want to delete this comment? This action cannot be undone.</>
                                            : <>Are you sure you want to delete this comment? This will also delete {replyCount} repl{replyCount === 1 ? 'y' : 'ies'}. This action cannot be undone.</>
                                    </>
                                );
                            }
                        )()}
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setDeleteCommentDialogOpen(false)}>
                        Cancel
                    </Button>
                    <Button
                        onClick={handleConfirmDeleteComment}
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