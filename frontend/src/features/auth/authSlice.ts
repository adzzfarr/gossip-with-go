// Redux slice for authentication state management
import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import type { LoginCredentials, RegisterCredentials } from '../../types';
import { authAPI } from '../../api/auth';
import { decodeToken } from '../../api/client';

interface AuthState {
    isAuthenticated: boolean;
    token: string | null;
    userID: number | null;
    username: string | null;
    loading: boolean;
    error: string | null;
}

const initialState: AuthState = {
    isAuthenticated: !!localStorage.getItem('token'),
    token: localStorage.getItem('token'),
    userID: localStorage.getItem('userID') 
        ? parseInt(localStorage.getItem('userID')!) 
        : null,
    username: localStorage.getItem('username'),
    loading: false,
    error: null,
};

// Register thunk
export const registerUser = createAsyncThunk(
    'auth/registerUser',
    async (credentials: RegisterCredentials, { rejectWithValue }) => {
        try {
            // Register
            await authAPI.registerUser(credentials);

            // Auto login after successful registration
            const loginResponse = await authAPI.loginUser(credentials);
            const token = loginResponse.token;

            // Write to localStorage
            const decoded = decodeToken(token)
            if (decoded) {
                localStorage.setItem('token', token);
                localStorage.setItem('userID', decoded.userID.toString());
                localStorage.setItem('username', decoded.username);
            }

            // Return user data to Redux store
            return {
                token,
                userID: decoded?.userID || null,
                username: decoded?.username || null,
            }
        } catch (error: any) {
            return rejectWithValue(error.response?.data?.error || 'Registration failed');
        }
    }
)

// Login thunk
export const loginUser = createAsyncThunk(
    'auth/loginUser',
    async (credentials: LoginCredentials, { rejectWithValue }) => {
        try {
            // Login user
            const response = await authAPI.loginUser(credentials);
            const token = response.token;

            // Write to localStorage
            const decoded = decodeToken(token)
            if (decoded) {
                localStorage.setItem('token', token);
                localStorage.setItem('userID', decoded.userID.toString());
                localStorage.setItem('username', decoded.username);
            }

            // Return user data to Redux store
            return {
                token,
                userID: decoded?.userID || null,
                username: decoded?.username || null,
            }
        } catch (error: any) {
            return rejectWithValue(error.response?.data?.error || 'Login failed');
        }
    }
);

// Logout thunk
export const logoutUser = createAsyncThunk(
    '/auth/logoutUser',
    async () => {
        // Clear localStorage
        localStorage.removeItem('token');
        localStorage.removeItem('userID');
        localStorage.removeItem('username');
        return null;
    }
)

// Auth slice
const authSlice = createSlice({
    name: 'auth',
    initialState,
    reducers: {
        // Manual error clearing for instant UI update
        clearError: (state) => {
            state.error = null;
        },
    },
    extraReducers: (builder) => { // Handle thunk states (Pending, Fulfilled, Rejected)
        // Register Pending
        builder.addCase(
            registerUser.pending, 
            (state) => {
                state.loading = true;
                state.error = null;
            }
        );

        // Register Fulfilled
        builder.addCase(
            registerUser.fulfilled,
            (state, action) => {
                state.isAuthenticated = true;
                state.token = action.payload.token;
                state.userID = action.payload.userID;
                state.username = action.payload.username;
                state.loading = false;
                state.error = null;
            }
        );

        // Register Rejected
        builder.addCase(
            registerUser.rejected,
            (state, action) => {
                state.loading = false;
                state.error = action.payload as string;
            }
        );

        // Login Pending
        builder.addCase(
            loginUser.pending,
            (state) => {
                state.loading = true;
                state.error = null;
            }
        );

        // Login Fulfilled
        builder.addCase(
            loginUser.fulfilled,
            (state, action) => {
                state.isAuthenticated = true;
                state.token = action.payload.token;
                state.userID = action.payload.userID;
                state.username = action.payload.username;
                state.loading = false;
                state.error = null;
            }
        );

        // Login Rejected
        builder.addCase(
            loginUser.rejected,
            (state, action) => {
                state.loading = false;
                state.error = action.payload as string;
            }
        );

        // Logout Fulfilled (No backend call => no Pending or Rejected states)
        builder.addCase(
            logoutUser.fulfilled,
            (state) => {
                state.isAuthenticated = false;
                state.token = null;
                state.userID = null;
                state.username = null;
                state.loading = false;
                state.error = null;
            }
        );
    },
});

export const { clearError } = authSlice.actions;
export default authSlice.reducer;