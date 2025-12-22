import { apiClient } from './client';

import type {
    LoginCredentials,
    LoginResponse,
    RegisterCredentials,
    RegisterResponse,
} from '../types';

export const authAPI = {
    // Register new user
    registerUser: async (credentials: RegisterCredentials): Promise<RegisterResponse> => {
        const response = await apiClient.post<RegisterResponse>('/users', credentials);
        return response.data;
    },

    // Login existing user
    loginUser: async (credentials: LoginCredentials): Promise<LoginResponse> => {
        const response = await apiClient.post<LoginResponse>('/login', credentials);
        return response.data;
    },
};