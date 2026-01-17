// Redux slice for topics state management
import { createAsyncThunk, createSlice, type PayloadAction } from '@reduxjs/toolkit';
import type { PaginationMetadata, Topic, TopicsResponse } from '../types';
import { apiClient } from '../api/client';

interface TopicsState {
    topics: Topic[];
    loading: boolean;
    error: string | null;
    submitting: boolean;
    submitError: string | null;
    sortBy: string;
    searchQuery: string;
    pagination: PaginationMetadata | null;
    currentPage: number;
}

const initialState: TopicsState = {
    topics: [],
    loading: false,
    error: null,
    submitting: false,
    submitError: null,
    sortBy: 'newest',
    searchQuery: '',
    pagination: null,
    currentPage: 1,
}

// Fetch topics
export const fetchTopics = createAsyncThunk(
    `topics/fetchTopics`,
    async (
        { sortBy = 'newest' , page = 1, search = '' }: {
            sortBy?: string;
            page?: number;
            search?: string;
        }, 
        { rejectWithValue }
    ) => {
        try {
            const params = new URLSearchParams({
                sort: sortBy,
                page: page.toString(),
                limit: '9',
            });

            if (search.trim() !== '') {
                params.append('search', search.trim());
            }

            const response = await apiClient.get<TopicsResponse>(`/topics?${params.toString()}`);
            return response.data;
        } catch (error: any) {
            return rejectWithValue(error.response?.data?.error || 'Failed to fetch topics');
        }
    }
)

export const fetchTopicByID = createAsyncThunk(
    'topics/fetchTopicByID',
    async (topicID: number, { rejectWithValue }) => {
        try {
            const response = await apiClient.get<Topic>(`/topics/${topicID}`);
            return response.data;
        } catch (error: any) {
            return rejectWithValue(error.response?.data?.error || 'Failed to fetch topic');
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

export const updateTopic = createAsyncThunk(
    'topics/updateTopic',
    async (
        { topicID, title, description } : {
            topicID: number;
            title: string;
            description: string;
        },
        { rejectWithValue }
    ) => {
        try {
            const response = await apiClient.put<Topic>(
                `/topics/${topicID}`,
                { title, description }
            );

            return response.data;
        } catch (error: any) {
            return rejectWithValue(error.response?.data?.error || 'Failed to update topic');
        }
    }
);

// Delete Topic
export const deleteTopic = createAsyncThunk(
    'topics/deleteTopic',
    async (topicID: number, { rejectWithValue }) => {
        try {
            await apiClient.delete(`/topics/${topicID}`);
            return topicID;
        } catch (error: any) {
            return rejectWithValue(error.response?.data?.error || 'Failed to delete topic');
        }
    }
)

const topicsSlice = createSlice({
    name: 'topics',
    initialState,
    reducers: {
        clearError: (state) => {
            state.error = null;
            state.submitError = null;
        },
        setTopicsSortBy: (state, action) => {
            state.sortBy = action.payload;
            state.currentPage = 1; // Reset to first page on sort change
        },
        setTopicsSearchQuery: (state, action) => {
            state.searchQuery = action.payload;
            state.currentPage = 1; // Reset to first page on new search
        },
        setTopicsCurrentPage: (state, action) => {
            state.currentPage = action.payload;
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
            (state, action) => {
                state.loading = false;
                state.topics = action.payload.topics;
                state.pagination = action.payload.pagination || null;
            }
        );

        builder.addCase(
            fetchTopics.rejected, 
            (state, action) => {
                state.loading = false;
                state.error = action.payload as string;
            }
        );

        // Fetch topic by ID
        builder.addCase(
            fetchTopicByID.pending,
            (state) => {
                state.loading = true;
                state.error = null;
            }
        );

        builder.addCase(
            fetchTopicByID.fulfilled,
            (state, action: PayloadAction<Topic>) => {
                state.loading = false;
                const index = state.topics.findIndex(topic => topic.topicID === action.payload.topicID);
                if (index === -1) {
                    state.topics.push(action.payload);
                } else {
                    state.topics[index] = action.payload;
                }
            }
        );

        builder.addCase(
            fetchTopicByID.rejected,
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

        // Update topic
        builder.addCase(
            updateTopic.pending,
            (state) => {
                state.submitting = true;
                state.submitError = null;
            }
        );

        builder.addCase(
            updateTopic.fulfilled,
            (state, action: PayloadAction<Topic>) => {
                state.submitting = false;
                const index = state.topics.findIndex(topic => topic.topicID === action.payload.topicID);
                if (index !== -1) {
                    state.topics[index] = action.payload;
                }
            }
        );

        builder.addCase(
            updateTopic.rejected,
            (state, action) => {
                state.submitting = false;
                state.submitError = action.payload as string;
            }
        );

        // Delete topic
        builder.addCase(
            deleteTopic.pending,
            (state) => {
                state.submitting = true;
                state.submitError = null;
            }
        );

        builder.addCase(
            deleteTopic.fulfilled,
            (state, action: PayloadAction<number>) => {
                state.submitting = false;
                // Remove deleted topic from state
                state.topics = state.topics.filter(topic => topic.topicID !== action.payload);
            }
        );

        builder.addCase(
            deleteTopic.rejected,
            (state, action) => {
                state.submitting = false;
                state.submitError = action.payload as string;
            }
        );
    }
})

export const { clearError, setTopicsSortBy, setTopicsSearchQuery, setTopicsCurrentPage } = topicsSlice.actions;
export default topicsSlice.reducer;