import { configureStore } from "@reduxjs/toolkit";
import authReducer from "./auth/authSlice";
import topicsReducer from "./topics/topicsSlice";
import postsReducer from "./posts/postsSlice";
import commentsReducer from "./comments/commentsSlice";

export const store = configureStore({
    reducer: {
        auth: authReducer,
        topics: topicsReducer,
        posts: postsReducer,
        comments: commentsReducer,
    },
});

// Infer `RootState` and `AppDispatch` types from the store itself
export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;