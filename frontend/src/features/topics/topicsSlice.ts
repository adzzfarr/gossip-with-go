// Redux slice for topics state management
import { createSlice } from '@reduxjs/toolkit';
import type { Topic } from '../../types';

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

const topicsSlice = createSlice({
    name: 'topics',
    initialState,
    reducers: {
        // TODO: Add reducers for CRUD operations
    },
})

export default topicsSlice.reducer;