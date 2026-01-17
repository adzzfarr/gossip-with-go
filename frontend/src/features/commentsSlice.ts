// Redux slice for comments state management
import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import type { VoteResponse, Comment, PaginationMetadata, CommentsResponse } from '../types'; 
import { apiClient } from '../api/client';

interface CommentsState {
    comments: Comment[];
    loading: boolean;
    error: string | null;
    submitting: boolean;
    submitError: string | null;
    sortBy: string;
    searchQuery: string;
    pagination: PaginationMetadata | null;
    currentPage: number;
}

const initialState: CommentsState = {
    comments: [],
    loading: false,
    error: null,
    submitting: false,
    submitError: null,
    sortBy: 'newest',
    searchQuery: '',
    pagination: null,
    currentPage: 1,
}

// Fetch comments by postID
export const fetchCommentsByPostID = createAsyncThunk(
    'comments/fetchCommentsByPostID',
    async (
        { postID, sortBy = 'newest', page = 1, search = '' }: { 
            postID: number; 
            sortBy: string;
            page: number;
            search: string;
        }, 
        { rejectWithValue }
    ) => {
        try {
            const token = localStorage.getItem('token');

            const params = new URLSearchParams({
                sort: sortBy,
                page: page.toString(),
                limit: '20',
            });

            if (search.trim() !== '') {
                params.append('search', search.trim());
            }

            const response = await apiClient.get<CommentsResponse>(
                `/posts/${postID}/comments?${params.toString()}`,
                token ? {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    }   
                } : undefined
            );
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
            const token = localStorage.getItem('token');
            
            const response = await apiClient.post<Comment>(
                `/posts/${postID}/comments`, 
                { content },
                token ? {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    }   
                } : undefined
            );
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
            const token = localStorage.getItem('token');

            const response = await apiClient.put<Comment>(
                `/comments/${commentID}`, 
                { content },
                token ? {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    }   
                } : undefined
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
            const token = localStorage.getItem('token');

            await apiClient.delete(
                `/comments/${commentID}`,
                token ? {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    }   
                } : undefined
            );
            return commentID;
        } catch (error: any) {
            return rejectWithValue(error.response?.data?.error || 'Failed to delete comment');
        }
    }
)

// Vote on Comment
export const voteOnComment = createAsyncThunk(
    'comments/voteOnComment',
    async (
        { commentID, voteType } : {
            commentID: number;
            voteType: 1 | -1;
        },
        { rejectWithValue }
    ) => {
        try {
            const token = localStorage.getItem('token');

            const response = await apiClient.post<VoteResponse>(
                `/comments/${commentID}/vote`,
                { voteType },
                token ? {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    }   
                } : undefined
            );

            return {
                commentID,
                ...response.data,
            }
        } catch (error: any) {
            return rejectWithValue(error.response?.data?.error || 'Failed to vote on comment');
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
        setCommentsSortBy: (state, action) => {
            state.sortBy = action.payload;
            state.currentPage = 1;
        },
        setCommentsSearchQuery: (state, action) => {
            state.searchQuery = action.payload;
            state.currentPage = 1;
        },
        setCommentsCurrentPage: (state, action) => {
            state.currentPage = action.payload;
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
                state.comments = action.payload.comments;
                state.pagination = action.payload.pagination;
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

        // Vote on comment
        builder.addCase(
            voteOnComment.pending,
            (state) => {
                state.error = null;
            }
        );

        builder.addCase(
            voteOnComment.fulfilled,
            (state, action) => {
                const { commentID, newVoteCount, userVote } = action.payload;
                const comment = state.comments.find(c => c.commentID === commentID);
                if (comment) {
                    comment.voteCount = newVoteCount;
                    comment.userVote = userVote as 1 | -1 | null;
                }
            }
        );

        builder.addCase(
            voteOnComment.rejected,
            (state, action) => {
                state.error = action.payload as string;
            }
        );
    },
});

export const { clearError: clearCommentsError, setCommentsSortBy, setCommentsSearchQuery, setCommentsCurrentPage } = commentsSlice.actions;
export default commentsSlice.reducer;