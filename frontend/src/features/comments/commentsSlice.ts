// Redux slice for comments state management
import { createSlice } from '@reduxjs/toolkit';
import type { Comment } from '../../types'; 

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

const commentsSlice = createSlice({
    name: 'comments',
    initialState,
    reducers: {
        // TODO: Add reducers for CRUD operations
    },
})

export default commentsSlice.reducer;