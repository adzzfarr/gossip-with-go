import { configureStore } from "@reduxjs/toolkit";
import authReducer from "./authSlice";
import topicsReducer from "./topicsSlice";
import postsReducer from "./postsSlice";
import commentsReducer from "./commentsSlice";
import profilesReducer from "./profilesSlice";

export const store = configureStore({
    reducer: {
        auth: authReducer,
        topics: topicsReducer,
        posts: postsReducer,
        comments: commentsReducer,
        profiles: profilesReducer,
    },
});

// Infer `RootState` and `AppDispatch` types from the store itself
export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;