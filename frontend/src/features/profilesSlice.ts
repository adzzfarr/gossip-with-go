import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { apiClient } from "../api/client";

interface User {
    userID: number;
    username: string;
    createdAt: string;
    updatedAt: string;
}

interface Post {
    postID: number;
    topicID: number;
    title: string;
    content: string;
    createdBy: number;
    username: string;
    createdAt: string;
    updatedAt: string;
}

interface Comment {
    commentID: number;
    postID: number;
    content: string;
    createdBy: number;
    username: string;
    createdAt: string;
    updatedAt: string;
}

interface ProfileState {
    profile: User | null;
    posts: Post[];
    comments: Comment[];
    loading: boolean;
    error: string | null;
}

const initialState: ProfileState = {
    profile: null,
    posts: [],
    comments: [],
    loading: false,
    error: null,
}

// Fetch user profile by ID
export const getUserProfileByID = createAsyncThunk(
    'profiles/fetchUserProfile',
    async (userID: number, { rejectWithValue }) => {
        try {
            const response = await apiClient.get<User>(`/users/${userID}`);
            return response.data;
        } catch(error: any) {
            return rejectWithValue(error.response?.data?.error || 'Failed to fetch user profile');
        }
    }
)

// Fetch posts by userID
export const getUserPosts = createAsyncThunk(
    'profiles/fetchUserPosts',
    async (userID: number, { rejectWithValue }) => {
        try {
            const response = await apiClient.get<Post[]>(`/users/${userID}/posts`);
            return response.data;
        } catch(error: any) {
            return rejectWithValue(error.response?.data?.error || 'Failed to fetch user posts');
        }
    }
)

// Fetch comments by userID
export const getUserComments = createAsyncThunk(
    'profiles/fetchUserComments',
    async (userID: number, { rejectWithValue }) => {
        try {
            const response = await apiClient.get<Comment[]>(`/users/${userID}/comments`);
            return response.data;
        } catch(error: any) {
            return rejectWithValue(error.response?.data?.error || 'Failed to fetch user comments');
        }
    }
)

const profilesSlice = createSlice({
    name: 'profiles',
    initialState,
    reducers: {
        clearError: (state) => {
            state.error = null;
        }
    },
    extraReducers: (builder) => {
        // Get user profile
        builder.addCase(getUserProfileByID.pending, (state) => {
            state.loading = true;
            state.error = null;
        });
        builder.addCase(getUserProfileByID.fulfilled, (state, action) => {
            state.loading = false;
            state.profile = action.payload;
        });
        builder.addCase(getUserProfileByID.rejected, (state, action) => {
            state.loading = false;
            state.error = action.payload as string;
        });

        // Get user posts
        builder.addCase(getUserPosts.pending, (state) => {
            state.loading = true;
            state.error = null;
        });
        builder.addCase(getUserPosts.fulfilled, (state, action) => {
            state.loading = false;
            state.posts = action.payload;
        });
        builder.addCase(getUserPosts.rejected, (state, action) => {
            state.loading = false;
            state.error = action.payload as string;
        });

        // Get user comments
        builder.addCase(getUserComments.pending, (state) => {
            state.loading = true;
            state.error = null;
        });
        builder.addCase(getUserComments.fulfilled, (state, action) => {
            state.loading = false;
            state.comments = action.payload;
        });
        builder.addCase(getUserComments.rejected, (state, action) => {
            state.loading = false;
            state.error = action.payload as string;
        });
    }
})

export const { clearError } = profilesSlice.actions;
export default profilesSlice.reducer;