// Redux slice for comments state management
import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import type { Comment } from '../../types'; 
import { apiClient } from '../../api/client';

interface CommentsState {
    comments: Comment[];
    loading: boolean;
    error: string | null;
    submitting: boolean;
    submitError: string | null;
}

const initialState: CommentsState = {
    comments: [],
    loading: false,
    error: null,
    submitting: false,
    submitError: null,
}

// Fetch comments by postID
export const fetchCommentsByPostID = createAsyncThunk(
    'comments/fetchCommentsByPostID',
    async (postID: number, { rejectWithValue }) => {
        try {
            const response = await apiClient.get<Comment[]>(`/posts/${postID}/comments`);
            return response.data;
        } catch (error: any) {
            return rejectWithValue(error.response?.data?.error || 'Failed to fetch comments');
        }
    }
);

// Create new comment
export const createComment = createAsyncThunk(
    'comments/createComment',
    async (
        { postID, content }: {
            postID: number; 
            content: string
        }, 
        { rejectWithValue }
    ) => {
        try {
            const response = await apiClient.post<Comment>(`/posts/${postID}/comments`, { content });
            return response.data;
        } catch (error: any) {
            return rejectWithValue(error.response?.data?.error || 'Failed to create comment');
        }
    }
);

// Update comment
export const updateComment = createAsyncThunk(
    'comments/updateComment',
    async (
        { commentID, content } : { commentID: number; content: string },
        { rejectWithValue }
    ) => {
        try {
            const response = await apiClient.put<Comment>(
                `/comments/${commentID}`, 
                { content }
            );
            return response.data;
        } catch (error: any) {
            return rejectWithValue(error.response?.data?.error || 'Failed to update comment');
        }
    }
);

// Delete comment
export const deleteComment = createAsyncThunk(
    'comments/deleteComment',
    async (commentID: number, { rejectWithValue }) => {
        try {
            await apiClient.delete(`/comments/${commentID}`);
            return commentID;
        } catch (error: any) {
            return rejectWithValue(error.response?.data?.error || 'Failed to delete comment');
        }
    }
)

const commentsSlice = createSlice({
    name: 'comments',
    initialState,
    reducers: {
        clearError: (state) => {
            state.error = null;
            state.submitError = null;
        },
    },
    extraReducers: (builder) => {
        // Fetch comments by postID
        builder.addCase(
            fetchCommentsByPostID.pending,
            (state) => {
                state.loading = true;
                state.error = null;
            }
        );

        builder.addCase(
            fetchCommentsByPostID.fulfilled,
            (state, action) => {
                state.loading = false;
                state.comments = action.payload;
            }
        );

        builder.addCase(
            fetchCommentsByPostID.rejected,
            (state, action) => {
                state.loading = false;
                state.error = action.payload as string;
            }
        );

        // Create new comment
        builder.addCase(
            createComment.pending,
            (state) => {
                state.submitting = true;
                state.submitError = null;
            }
        );

        builder.addCase(
            createComment.fulfilled,
            (state, action) => {
                state.submitting = false;
                state.comments.push(action.payload);
            }
        );

        builder.addCase(
            createComment.rejected,
            (state, action) => {
                state.submitting = false;
                state.submitError = action.payload as string;
            }
        );

        // Update comment
        builder.addCase(
            updateComment.pending,
            (state) => {
                state.submitting = true;
                state.submitError = null;
            }
        );

        builder.addCase(
            updateComment.fulfilled,
            (state, action) => {
                state.submitting = false;

                const index = state.comments.findIndex(c => c.commentID === action.payload.commentID);
                
                if (index !== -1) {
                    state.comments[index] = action.payload;
                } 
            }
        );

        builder.addCase(
            updateComment.rejected,
            (state, action) => {
                state.submitting = false;
                state.submitError = action.payload as string;
            }
        );

        // Delete comment
        builder.addCase(
            deleteComment.pending,
            (state) => {
                state.submitting = true;
                state.submitError = null;
            }
        );

        builder.addCase(
            deleteComment.fulfilled,
            (state, action) => {
                state.submitting = false;
                state.comments = state.comments.filter(comment => comment.commentID !== action.payload);
            }
        );

        builder.addCase(
            deleteComment.rejected,
            (state, action) => {
                state.submitting = false;
                state.submitError = action.payload as string;
            }
        );
    },
});

export const { clearError: clearCommentsError } = commentsSlice.actions;
export default commentsSlice.reducer;