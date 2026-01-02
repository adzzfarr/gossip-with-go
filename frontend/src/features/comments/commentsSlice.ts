// Redux slice for comments state management
import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import type { Comment } from '../../types'; 
import { apiClient } from '../../api/client';

interface CommentsState {
    comments: Comment[];
    loading: boolean;
    error: string | null;
}

const initialState: CommentsState = {
    comments: [],
    loading: false,
    error: null,
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

const commentsSlice = createSlice({
    name: 'comments',
    initialState,
    reducers: {
        clearError: (state) => {
            state.error = null;
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
    },
});

export const { clearError: clearCommentsError } = commentsSlice.actions;
export default commentsSlice.reducer;