// Redux slice for topics state management
import { createSlice } from '@reduxjs/toolkit';
import type { Post } from '../../types';

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

const postsSlice = createSlice({
    name: 'posts',
    initialState,
    reducers: {
        // TODO: Add reducers for CRUD operations
    },
})

export default postsSlice.reducer;