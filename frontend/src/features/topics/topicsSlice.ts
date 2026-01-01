// Redux slice for topics state management
import { createAsyncThunk, createSlice, type PayloadAction } from '@reduxjs/toolkit';
import type { Topic } from '../../types';
import { apiClient } from '../../api/client';

interface TopicsState {
    topics: Topic[];
    currentTopic: Topic | null;
    loading: boolean;
    error: string | null;
}

const initialState: TopicsState = {
    topics: [],
    currentTopic: null,
    loading: false,
    error: null,
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

const topicsSlice = createSlice({
    name: 'topics',
    initialState,
    reducers: {
        clearError: (state) => {
            state.error = null;
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
    }
})

export const { clearError } = topicsSlice.actions;
export default topicsSlice.reducer;