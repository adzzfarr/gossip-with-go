// Redux slice for authentication state management
import { createSlice } from '@reduxjs/toolkit';

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
    userID: localStorage.getItem('userID') ? parseInt(localStorage.getItem('userID')!) : null,
    username: localStorage.getItem('username'),
    loading: false,
    error: null,
};

const authSlice = createSlice({
    name: 'auth',
    initialState,
    reducers: {
        // TODO: Add reducers for login, logout, register
    }
})

export default authSlice.reducer;