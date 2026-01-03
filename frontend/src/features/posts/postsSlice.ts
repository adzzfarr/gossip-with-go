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
export const fetchPostByID = createAsyncThunk(
    'posts/fetchPostByID',
    async (postID: number, { rejectWithValue }) => {
        try {
            const response = await apiClient.get<Post>(`/posts/${postID}`);
            return response.data;
        } catch (error: any) {
            return rejectWithValue(error.response?.data?.error || 'Failed to fetch post');
        }
    }
);

// Create new post
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

// Update post
export const updatePost = createAsyncThunk(
    'posts/updatePost',
    async (
        { postID, title, content }: { postID: number; title: string; content: string },
        { rejectWithValue }
    ) => {
        try {
            const response = await apiClient.put<Post>(
                `/posts/${postID}`,
                { title, content }
            )

            return response.data;
        } catch (error: any) {
            return rejectWithValue(error.response?.data?.error || 'Failed to update post');
        }
    }
);

// Delete post
export const deletePost = createAsyncThunk(
    'posts/deletePost',
    async (postID: number, { rejectWithValue }) => {
        try {
            await apiClient.delete(`/posts/${postID}`);
            return postID;
        } catch (error: any) {
            return rejectWithValue(error.response?.data?.error || 'Failed to delete post');
        }
    }
)

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
            fetchPostByID.pending,
            (state) => {
                state.loading = true;
                state.error = null;
            }
        );

        builder.addCase(
            fetchPostByID.fulfilled,
            (state, action) => {
                state.loading = false;
                state.currentPost = action.payload;
            }
        );

        builder.addCase(
            fetchPostByID.rejected,
            (state, action) => {
                state.loading = false;
                state.error = action.payload as string;
            }
        )

        // Create new post
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

        // Update post
        builder.addCase(
            updatePost.pending,
            (state) => {
                state.submitting = true;
                state.submitError = null;
            }
        );

        builder.addCase(
            updatePost.fulfilled,
            (state, action) => {
                state.submitting = false;

                // Update the post in posts list
                const index = state.posts.findIndex(post => post.postID === action.payload.postID);
                if (index !== -1) {
                    state.posts[index] = action.payload;
                }

                // Also update currentPost if it matches
                if (state.currentPost && state.currentPost.postID === action.payload.postID) {
                    state.currentPost = action.payload;
                }
            }
        );

        builder.addCase(
            updatePost.rejected,
            (state, action) => {
                state.submitting = false;
                state.submitError = action.payload as string;
            }
        );

        // Delete post
        builder.addCase(
            deletePost.pending,
            (state) => {
                state.submitting = true;
                state.submitError = null;
            }
        );

        builder.addCase(
            deletePost.fulfilled,
            (state, action) => {
                state.submitting = false;

                // Remove the deleted post from posts list
                state.posts = state.posts.filter(post => post.postID !== action.payload);

                // Clear currentPost if it was deleted
                if (state.currentPost && state.currentPost.postID === action.payload) {
                    state.currentPost = null;
                }
            }
        );

        builder.addCase(
            deletePost.rejected,
            (state, action) => {
                state.submitting = false;
                state.submitError = action.payload as string;
            }
        );
    }
})

export const { clearError } = postsSlice.actions;
export default postsSlice.reducer;