import { useState, type FormEvent } from "react";
import type { Comment as CommentType } from "../types";
import { Box, Button, Card, CardContent, CircularProgress, Divider, TextField, Typography } from "@mui/material";
import Username from "./Username";
import { Delete, Edit, Reply } from "@mui/icons-material";
import VoteButtons from "./VoteButtons";

interface CommentCardProps {
    comment: CommentType;
    userID?: number;
    depth?: number;
    maxDepth?: number;
    isSubmitting: boolean;
    onEdit: (commentID: number, content: string) => void;
    onDelete: (commentID: number) => void;
    onReply: (parentCommentID: number, content: string) => void;
}

const MAX_DEPTH = 5;

export default function CommentCard({
    comment,
    userID,
    depth = 0,
    maxDepth = MAX_DEPTH,
    isSubmitting,
    onEdit,
    onDelete,
    onReply,
}: CommentCardProps) {
    const isAuthor = comment.createdBy === userID;
    const canReply = depth < maxDepth;

    // State for editing
    const [isEditing, setIsEditing] = useState(false);
    const [editContent, setEditContent] = useState(comment.content);

    // State for replying
    const [isReplying, setIsReplying] = useState(false);
    const [replyContent, setReplyContent] = useState('');

    // Editing handlers
    const handleEdit = () => {
        setIsEditing(true);
        setEditContent(comment.content);
    }

    const handleCancelEdit = () => {
        setIsEditing(false);
        setEditContent(comment.content);
    }

    const handleSaveEdit = async () => {
        if (!editContent.trim()) return;
        await onEdit(comment.commentID, editContent.trim());
        setIsEditing(false);
    }

    // Replying handlers
    const handleStartReply = () => {
        setIsReplying(true);
    }

    const handleCancelReply = () => {
        setIsReplying(false);
        setReplyContent('');
    }

    const handleSubmitReply = async (e: FormEvent) => {
        e.preventDefault();
        if (!replyContent.trim()) return;
        await onReply(comment.commentID, replyContent.trim());
        setIsReplying(false);
        setReplyContent('');
    }

    return (
        <Box
            sx={{
                ml: depth > 0 ? 3 : 0,
                borderLeft: depth > 0 ? '2px solid' : 'none',
                borderColor: depth > 0 ? 'divider' : 'transparent',
                pl: depth > 0 ? 2 : 0,
            }}
        >
            <Card
                variant="outlined"
                sx={{
                    mb: 2,
                }}
            >
                <CardContent>
                    {/* Comment Header */}
                    <Box
                        sx={{
                            display: 'flex',
                            justifyContent: 'space-between',
                            alignItems: 'center',
                            mb: 1,
                        }}
                    >
                        <Box
                            sx={{
                                display: 'flex',
                                alignItems: 'center',
                                gap: 2,
                            }}
                        >
                            <Username 
                                userID={comment.createdBy}
                                username={comment.username || "Unknown"}
                                variant="body2"
                                fontWeight="bold"
                            />

                            <Typography variant="body2" color="text.secondary">
                                {new Date(comment.createdAt).toLocaleString()}
                            </Typography>
                        </Box>

                        {isAuthor && !isEditing && (
                            <Box
                                sx={{
                                    display: 'flex',
                                    gap: 1,
                                }}
                            >
                                <Button
                                    size="small"
                                    startIcon={<Edit />}
                                    onClick={handleEdit}
                                >
                                    Edit
                                </Button>

                                <Button
                                    size="small"
                                    color="error"
                                    startIcon={<Delete />}
                                    onClick={() => onDelete(comment.commentID)}
                                    disabled={isSubmitting}
                                >
                                    Delete
                                </Button>
                            </Box>
                        )}
                    </Box>

                    {/* Edit and Delete Buttons */}


                    {/* Comment Content or Edit Form */}
                    {
                        isEditing ? (
                            <Box>
                                <TextField 
                                    fullWidth
                                    multiline
                                    minRows={2}
                                    value={editContent}
                                    onChange={(e) => setEditContent(e.target.value)}
                                    disabled={isSubmitting}
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
                                        onClick={handleSaveEdit}
                                        disabled={isSubmitting || !editContent.trim()}
                                    >
                                        {isSubmitting ? <CircularProgress size={24} /> : 'Save'}
                                    </Button>   
                                    
                                    <Button
                                        size="small"
                                        variant="outlined"
                                        onClick={handleCancelEdit}
                                        disabled={isSubmitting}
                                    >
                                        Cancel
                                    </Button>
                                </Box>
                            </Box>
                        ) : (
                            <>
                                <Typography
                                    variant="body1"
                                    sx={{ 
                                        whiteSpace: 'pre-wrap', 
                                        mb: 2,
                                    }}
                                >
                                    {comment.content}
                                </Typography>

                                <Divider sx={{ my: 1 }} />

                                {/* Vote and Reply Buttons */}
                                <Box
                                    sx={{
                                        display: 'flex',
                                        gap: 2,
                                        alignItems: 'center',
                                    }}
                                >
                                    <VoteButtons 
                                        commentID={comment.commentID}
                                        initialVoteCount={comment.voteCount}
                                        initialUserVote={comment.userVote}
                                        orientation="horizontal"
                                        size="small"
                                    />

                                    {canReply && userID && (
                                        <Button
                                            size="small"
                                            startIcon={<Reply />}
                                            onClick={handleStartReply}
                                            disabled={isSubmitting}                                        
                                        >
                                            Reply
                                        </Button>
                                    )}

                                    {!canReply && (
                                        <Typography variant="caption" color="text.secondary">
                                            Max reply depth reached
                                        </Typography>
                                    )}
                                </Box>
                            </>
                        )
                    }

                    {/* Reply Form */}
                    {isReplying && (
                        <Box
                            sx={{ mt: 2}}
                        >
                            <form onSubmit={handleSubmitReply}>
                                <TextField 
                                    fullWidth
                                    multiline
                                    minRows={2}
                                    placeholder={`Reply to ${comment.username || "this comment"}...`}
                                    value={replyContent}
                                    onChange={(e) => setReplyContent(e.target.value)}
                                    disabled={isSubmitting}
                                    sx={{ mb: 1 }}
                                />

                                <Box 
                                    sx={{ 
                                        display: 'flex', 
                                        gap: 1 
                                }}>
                                    <Button
                                        type="submit"
                                        variant="contained"
                                        size="small"
                                        disabled={isSubmitting || !replyContent.trim()}
                                    >
                                        {isSubmitting ? <CircularProgress size={24} /> : 'Reply'}
                                    </Button>

                                    <Button
                                        variant="outlined"
                                        size="small"
                                        onClick={handleCancelReply}
                                        disabled={isSubmitting}
                                    >
                                        Cancel
                                    </Button>
                                </Box>
                            </form>
                        </Box>
                    )}
                </CardContent>
            </Card>

            {/* Render Replies */}
            {comment.replies && comment.replies.length > 0 && (
                <Box sx={{ mt: 1 }}>
                    {comment.replies.map(
                        reply => (
                            <CommentCard 
                                key={reply.commentID}
                                comment={reply}
                                userID={userID}
                                depth={depth + 1}
                                maxDepth={maxDepth}
                                isSubmitting={isSubmitting}
                                onEdit={onEdit}
                                onDelete={onDelete}
                                onReply={onReply}
                            />
                        )
                    )}
                </Box>
            )}
        </Box>
    );
}