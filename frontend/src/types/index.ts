// All types match Go struct JSON tags

// Data models
export interface User {
    userID: number;
    username: string;
    createdAt: string;
    updatedAt: string;
}

export interface Topic {
    topicID: number;
    title: string;
    description: string;
    createdBy: number;
    createdAt: string;
    updatedAt: string;
}

export interface Post {
    postID: number;
    topicID: number;
    title: string;
    content: string;
    createdBy: number;
    createdAt: string;
    updatedAt: string;
}

export interface Comment {
    commentID: number;
    postID: number;
    content: string;
    createdBy: number;
    createdAt: string;
    updatedAt: string;
}

// Auth types for login/register 
export interface RegisterCredentials {
    username: string;
    password: string;
} 

export interface RegisterResponse { // returned by RegisterUser in UserHandler
    message: string;
    user: User;
}

export interface LoginCredentials { 
    username: string;
    password: string;
}

export interface LoginResponse { // returned by LoginUser in LoginHandler
    message: string;
    token: string;
}

// Generic API error response
export interface APIError {
    error: string;
}

// Request types for creating and updating
export interface CreateTopicRequest {
    title: string;
    description: string;
}

export interface UpdateTopicRequest {
    title: string;
    description: string;
}

export interface CreatePostRequest {
    title: string;
    content: string;
}

export interface UpdatePostRequest {
    title: string;
    content: string;
}

export interface CreateCommentRequest {
    content: string;
}

export interface UpdateCommentRequest {
    content: string;
}
