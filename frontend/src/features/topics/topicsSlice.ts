// Redux slice for topics state management
import { createAsyncThunk, createSlice, type PayloadAction } from '@reduxjs/toolkit';
import type { Topic } from '../../types';
import { apiClient } from '../../api/client';

interface TopicsState {
    topics: Topic[];
    loading: boolean;
    error: string | null;
    submitting: boolean;
    submitError: string | null;
}

const initialState: TopicsState = {
    topics: [],
    loading: false,
    error: null,
    submitting: false,
    submitError: null,
}

// Fetch topics
export const fetchTopics = createAsyncThunk(
    `topics/fetchTopics`,
    async (_, { rejectWithValue }) => {
        try {
            const response = await apiClient.get<Topic[]>('/topics');
            return response.data;
        } catch (error: any) {
            return rejectWithValue(error.response?.data?.error || 'Failed to fetch topics');
        }
    }
)

// Create topic
export const createTopic = createAsyncThunk(
    'topics/createTopic',
    async (
        { title, description } : {
            title: string;
            description: string;
        },  
        { rejectWithValue }
    ) => {
        try {
            const response = await apiClient.post<Topic>(
                '/topics',
                { title, description }
            )

            return response.data;
        } catch (error: any) {
            return rejectWithValue(error.response?.data?.error || 'Failed to create topic');
        }
    }
);

const topicsSlice = createSlice({
    name: 'topics',
    initialState,
    reducers: {
        clearError: (state) => {
            state.error = null;
            state.submitError = null;
        },
    },
    extraReducers: (builder) => {
        builder.addCase(
            fetchTopics.pending, 
            (state) => {
                state.loading = true;
                state.error = null;
            }
        );

        builder.addCase(
            fetchTopics.fulfilled, 
            (state, action: PayloadAction<Topic[]>) => {
                state.loading = false;
                state.topics = action.payload;
            }
        );

        builder.addCase(
            fetchTopics.rejected, 
            (state, action) => {
                state.loading = false;
                state.error = action.payload as string;
            }
        );

        // Create topic
        builder.addCase(
            createTopic.pending,
            (state) => {
                state.submitting = true;
                state.submitError = null;
            }
        );

        builder.addCase(
            createTopic.fulfilled,
            (state, action: PayloadAction<Topic>) => {
                state.submitting = false;
                state.topics.unshift(action.payload); // Add new topic to the beginning
            }
        );

        builder.addCase(
            createTopic.rejected,
            (state, action) => {
                state.submitting = false;
                state.submitError = action.payload as string;
            }
        );
    }
})

export const { clearError } = topicsSlice.actions;
export default topicsSlice.reducer;