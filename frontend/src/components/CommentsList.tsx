import { Alert, Box } from "@mui/material";
import type { Comment as CommentType } from "../types";
import CommentCard from "./CommentCard";

interface CommentsListProps {
    comments: CommentType[];
    userID?: number;
    isSubmitting: boolean;
    onEdit: (commentID: number, content: string) => void;
    onDelete: (commentID: number) => void;
    onReply: (parentCommentID: number, content: string) => void;
}

export default function CommentsList({
    comments,
    userID,
    isSubmitting,
    onEdit,
    onDelete,
    onReply
}: CommentsListProps) {
    if (comments.length === 0) {
        return (
            <Alert severity="info">No comments yet.</Alert>
        );
    }

    return (
        <Box
            sx={{
                display: 'flex',
                flexDirection: 'column',
                gap: 2,
            }}
        >
            {comments.map(
                comment => (
                    <CommentCard 
                        key={comment.commentID}
                        comment={comment}
                        userID={userID}
                        isSubmitting={isSubmitting}
                        onEdit={onEdit}
                        onDelete={onDelete}
                        onReply={onReply}
                    />
                )
            )}
        </Box>
    );
}