// Redux slice for topics state management
import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import type { Post } from '../../types';
import { apiClient } from '../../api/client';

interface PostsState {
    posts: Post[];
    currentPost: Post | null;
    loading: boolean;
    error: string | null;
    submitting: boolean;
    submitError: string | null;
}

const initialState: PostsState = {
    posts: [],
    currentPost: null,
    loading: false,
    error: null,
    submitting: false,
    submitError: null,
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

// Fetch single post by ID
export const fetchPostById = createAsyncThunk(
    'posts/fetchPostById',
    async (postID: number, { rejectWithValue }) => {
        try {
            const response = await apiClient.get<Post>(`/posts/${postID}`);
            return response.data;
        } catch (error: any) {
            return rejectWithValue(error.response?.data?.error || 'Failed to fetch post');
        }
    }
);

// Create a new post
export const createPost = createAsyncThunk(
    'posts/createPost',
    async (
        { topicID, title, content } : {
            topicID: number;
            title: string;
            content: string;
        },
        { rejectWithValue }
     ) => {
        try {
            const response = await apiClient.post<Post>(`/topics/${topicID}/posts`, { title, content });
            return response.data;
        } catch (error: any) {
            return rejectWithValue(error.response?.data?.error || 'Failed to create post');
        }
    }
);

const postsSlice = createSlice({
    name: 'posts',
    initialState,
    reducers: {
        clearError: (state) => {
            state.error = null;
            state.submitError = null;
        }
    },
    extraReducers: (builder) => {
        // Fetch posts by topic
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

        // Fetch single post by ID
        builder.addCase(
            fetchPostById.pending,
            (state) => {
                state.loading = true;
                state.error = null;
            }
        );

        builder.addCase(
            fetchPostById.fulfilled,
            (state, action) => {
                state.loading = false;
                state.currentPost = action.payload;
            }
        );

        builder.addCase(
            fetchPostById.rejected,
            (state, action) => {
                state.loading = false;
                state.error = action.payload as string;
            }
        )

        // Create a new post
        builder.addCase(
            createPost.pending,
            (state) => {
                state.submitting = true;
                state.submitError = null;
            }
        );

        builder.addCase(
            createPost.fulfilled,
            (state, action) => {
                state.submitting = false;
                // Add new post to the top of posts list
                state.posts.unshift(action.payload);
            }
        );

        builder.addCase(
            createPost.rejected,
            (state, action) => {
                state.submitting = false;
                state.submitError = action.payload as string;
            }
        );
    }
})

export const { clearError } = postsSlice.actions;
export default postsSlice.reducer;