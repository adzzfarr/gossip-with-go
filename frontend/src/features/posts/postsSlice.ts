// Redux slice for topics state management
import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import type { Post } from '../../types';
import { apiClient } from '../../api/client';

interface PostsState {
    posts: Post[];
    currentPost: Post | null;
    loading: boolean;
    error: string | null;
}

const initialState: PostsState = {
    posts: [],
    currentPost: null,
    loading: false,
    error: null,
}

// Fetch posts by topic
export const fetchPostsByTopic = createAsyncThunk(
    'posts/fetchPostsByTopic',
    async (topicID: number, { rejectWithValue }) => {
        try {
            const response = await apiClient.get<Post[]>(`/topics/${topicID}/posts`);
            return response.data;
        } catch (error: any) {
            return rejectWithValue(error.response?.data?.error || 'Failed to fetch posts');
        }
    }
)

const postsSlice = createSlice({
    name: 'posts',
    initialState,
    reducers: {
        clearError: (state) => {
            state.error = null;
        }
    },
    extraReducers: (builder) => {
        builder.addCase(
            fetchPostsByTopic.pending,
            (state) => {
                state.loading = true;
                state.error = null;
            }
        );

        builder.addCase(
            fetchPostsByTopic.fulfilled,
            (state, action) => {
                state.loading = false;
                state.posts = action.payload;
            }
        );

        builder.addCase(
            fetchPostsByTopic.rejected,
            (state, action) => {
                state.loading = false;
                state.error = action.payload as string;
            }
        )
    }
})

export const { clearError } = postsSlice.actions;
export default postsSlice.reducer;