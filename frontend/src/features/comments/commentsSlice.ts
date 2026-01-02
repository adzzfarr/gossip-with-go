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

// Create a new comment
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

        // Create a new comment
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
    },
});

export const { clearError: clearCommentsError } = commentsSlice.actions;
export default commentsSlice.reducer;