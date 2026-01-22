# Gossip With Go: A Full Stack Web Forum

**Created by:** Adzfar

A modern, feature-rich forum application built with Go, React, and PostgreSQL. Users can create topics, engage in discussions, comment on posts, and vote on content with real-time vote counts. This project was created as part of my application for NUS CVWO.

---

## Contents
- [Features](#features)
- [Technology Stack](#technology-stack)
- [Project Structure](#project-structure)
- [Local Development Setup](#local-development-setup)
  - [Prerequisites](#prerequisites)
  - [Backend Setup](#2-backend-setup)
  - [Frontend Setup](#3-frontend-setup)
  - [Database Schema](#4-database-schema)
- [Testing the Application](#testing-the-application)
  - [Seeding Test Data](#seeding-test-data)
  - [Manual API Testing](#manual-api-testing)
- [Design Decisions & Implementation](#design-decisions--implementation-details)
  - [Frontend Architecture](#frontend-architecture)
  - [Backend Architecture](#backend-architecture)
- [AI Usage Disclosure](#ai-usage-disclosure)
- [API Documentation](#api-documentation)
  - [Authentication Endpoints](#authentication-endpoints)
  - [Topic Endpoints](#topic-endpoints)
  - [Post Endpoints](#post-endpoints)
  - [Comment Endpoints](#comment-endpoints)
  - [Vote Endpoints](#vote-endpoints)
  - [User Profile Endpoints](#user-profile-endpoints)
- [Troubleshooting](#troubleshooting)
  - [Backend Issues](#backend-issues)
  - [Frontend Issues](#frontend-issues)
  - [Database Issues](#database-issues)
- [Production Deployment](#production-deployment)
- [Future Enhancements](#future-enhancements)
- [License](#license)
- [Author](#author)
  
---

## Features

- **User Authentication**: Secure JWT-based registration and login with protected routes
- **Topic Management**: Create, edit, and delete discussion topics
- **Posts & Comments**: Threaded discussions with nested comment support
- **Voting System**: Upvote/downvote posts and comments with automatic vote count updates
- **User Profiles**: View user activity including posts and comments
- **Dark/Light Theme**: Toggle between themes with persistent preference
- **Responsive Design**: Mobile-friendly interface using Material-UI
- **Sorting & Pagination**: Sort content by newest, top-rated, or most discussed
- **RESTful API**: Clean, documented backend with comprehensive error handling

---

## Technology Stack

### Backend
- **Go 1.25** - Primary backend language
- **Gin** - HTTP web framework for routing and middleware
- **PostgreSQL** - Relational database with pgx/v5 driver
- **JWT** - JSON Web Tokens for stateless authentication
- **golang-migrate** - Database migration management
- **bcrypt** - Password hashing

### Frontend
- **React** - UI library
- **TypeScript** - Type-safe JavaScript
- **Redux Toolkit** - Centralized state management
- **Material UI (MUI)** - Component library and theming
- **Vite** - Fast build tool and dev server
- **React Router v6** - Client-side routing
- **Axios** - HTTP client for API requests

### Deployment
- **Backend**: Render
- **Frontend**: Vercel 
- **Database**: Render PostgreSQL 

---

## Project Structure

```
gossip-with-go/
├── backend/
│   ├── cmd/
│   │   ├── server/         # Main application entry point
│   │   │   └── main.go
│   │   └── seed/           # Database seeding utility
│   │       └── seed.go
│   ├── internal/
│   │   ├── api/            # HTTP handlers, middleware, tests
│   │   │   ├── *_handler.go
│   │   │   ├── auth_middleware.go
│   │   │   └── api_test.go
│   │   ├── data/           # Database layer (repository pattern)
│   │   │   ├── db.go
│   │   │   ├── repository.go
│   │   │   ├── models.go
│   │   │   └── repository_test.go
│   │   └── service/        # Business logic layer
│   │       ├── *_service.go
│   │       ├── jwt.go
│   │       └── login_test.go
│   ├── migrations/         # SQL migration files
│   │   ├── 01_initial_schema.up.sql
│   │   ├── 02_votes.up.sql
│   │   └── 03_comment_replies.up.sql
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── components/     # Reusable React components
│   │   │   ├── Layout.tsx
│   │   │   ├── AppBar.tsx
│   │   │   ├── TopicCard.tsx
│   │   │   ├── PostCard.tsx
│   │   │   ├── CommentCard.tsx
│   │   │   ├── VoteButtons.tsx
│   │   │   └── ...
│   │   ├── pages/          # Page-level components
│   │   │   ├── TopicsPage.tsx
│   │   │   ├── TopicPostsPage.tsx
│   │   │   ├── PostPage.tsx
│   │   │   ├── ProfilePage.tsx
│   │   │   └── ...
│   │   ├── features/       # Redux Toolkit slices
│   │   │   ├── store.ts
│   │   │   ├── authSlice.ts
│   │   │   ├── topicsSlice.ts
│   │   │   ├── postsSlice.ts
│   │   │   ├── commentsSlice.ts
│   │   │   └── profilesSlice.ts
│   │   ├── api/            # API client configuration
│   │   │   ├── client.ts
│   │   │   └── auth.ts
│   │   ├── contexts/       # React contexts
│   │   │   └── ThemeContext.tsx
│   │   ├── hooks/          # Custom hooks
│   │   │   └── redux.ts
│   │   ├── types/          # TypeScript type definitions
│   │   │   └── index.ts
│   │   ├── App.tsx
│   │   ├── main.tsx
│   │   └── theme.ts
│   ├── vercel.json         # Vercel deployment config
│   └── package.json
└── README.md
```

---

## Local Development Setup

### Prerequisites

- **Go 1.25+** 
- **Node.js 18+** and npm 
- **PostgreSQL 14+** 
- **golang-migrate CLI** (

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/gossip-with-go.git
cd gossip-with-go
```

### 2. Backend Setup

```bash
cd backend

# Install Go dependencies
go mod download

# Set environment variables
export DATABASE_URL="postgresql://user:password@localhost:5432/gossip_forum?sslmode=disable"
export JWT_SECRET="your-super-secret-jwt-key-change-this"
export PORT=8080

# Run database migrations
migrate -path migrations -database "${DATABASE_URL}" up

# Start the backend server
go run cmd/server/main.go
```

**Backend will be running at:** `http://localhost:8080`

**API Health Check:** `http://localhost:8080/api/v1/health`

### 3. Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Set environment variable (create .env file)
echo "VITE_API_URL=http://localhost:8080/api/v1" > .env

# Start development server
npm run dev
```

**Frontend will be running at:** `http://localhost:5173`

### 4. Database Schema

The application uses PostgreSQL with three migrations:

1. **Initial Schema** - Users, topics, posts, comments tables with foreign keys
2. **Votes System** - Votes table with triggers for automatic vote_count updates
3. **Comment Replies** - Parent-child relationships for nested comments

Migrations automatically create:
- Tables with proper constraints and indexes
- Foreign key relationships with CASCADE deletes
- Database triggers for vote counting
- Sequences for auto-incrementing IDs

---

## Testing the Application

### Seeding Test Data

You can populate the database with realistic test data:

```bash
# Option 1: Run seed script directly (if database is accessible)
cd backend
go run cmd/seed/seed.go

# Option 2: Use the temporary seed API endpoint (in production)
1. In backend/cmd/server/main.go,
   uncomment line 84: seedHandler := api.NewSeedHandler(dbPool)
        and line 145: router.POST("/seed", seedHandler.SeedDatabase)
2. Run: curl -X POST https://gossip-with-go.onrender.com/api/v1/seed 
```

**Seed data includes:**
- 10 test users (password: `Password123`)
  - `alice_wonder`, `bob_builder`, `charlie_brown`, etc.
- 15 diverse topics (Technology, Gaming, Cooking, etc.)
- 300 posts across topics
- 7,500 comments with varied content
- Thousands of votes (upvotes and downvotes)

### Manual API Testing

**Register a new user:**
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Password123"
  }'
```

**Login:**
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Password123"
  }'
```

**Get all topics:**
```bash
curl http://localhost:8080/api/v1/topics
```

---

## Design Decisions & Implementation Details

### Frontend Architecture

I implemented a **Redux-based architecture** with TypeScript for type safety and scalability:

**State Management with Redux Toolkit:**
- Created separate **slices** for each domain (auth, topics, posts, comments, profiles)
- Used **createAsyncThunk** for handling API calls with loading/error states
- Implemented **typed hooks** (`useAppDispatch`, `useAppSelector`) for type-safe Redux usage
- Centralized API configuration in `api/client.ts` with automatic token injection

**Component Design:**
- Built **reusable components** (VoteButtons, CommentCard, PostCard, TopicCard)
- Implemented **Layout component** with persistent navigation and theme toggle
- Created **custom breadcrumbs** for navigation hierarchy
- Designed **Material-UI based** UI with consistent theming

**Theme System:**
- Implemented **ThemeContext** for dark/light mode toggling
- Configured **Material-UI theme** with custom color palettes
- Stored **theme preference** in localStorage for persistence
- Applied theme to all MUI components automatically

**Routing & Navigation:**
- Set up **React Router v6** with protected routes
- Implemented **navigation guards** for authenticated-only pages
- Created **breadcrumb navigation** that reflects current location
- Configured **Vercel rewrites** for SPA routing (no 404 on refresh)

### Backend Architecture

I designed a **clean, layered architecture** following best practices:

**Three-Layer Architecture:**

1. **API Layer** (`internal/api/`)
   - HTTP handlers for each resource (topics, posts, comments, votes, users)
   - JWT authentication middleware with token validation
   - Request validation and error handling
   - Response formatting with consistent JSON structure
   - Optional authentication middleware for guest access

2. **Service Layer** (`internal/service/`)
   - Business logic implementation
   - Data validation and sanitization
   - JWT token generation and validation
   - Password hashing with bcrypt
   - Cross-service communication (e.g., vote service updates post/comment services)

3. **Data Layer** (`internal/data/`)
   - Repository pattern for database operations
   - PostgreSQL connection pooling with pgx/v5
   - Parameterized queries to prevent SQL injection
   - Transaction support for complex operations
   - Clean separation from business logic

**Authentication & Security:**
- **JWT tokens** with 24-hour expiration
- **Bcrypt password hashing** (cost factor 10)
- **CORS middleware** configured for Vercel domains
- **Protected routes** require valid JWT in Authorization header
- **Optional auth routes** allow guests to browse (posts, topics) but authenticated users see their votes
- **SQL injection prevention** via parameterized queries

**Database Design:**
- **Automatic vote counting** via PostgreSQL triggers (efficient, no manual updates)
- **Foreign key constraints** with CASCADE deletes maintain referential integrity
- **Composite unique constraints** prevent duplicate votes (user_id + post_id/comment_id)
- **Indexed columns** on frequently queried fields (created_by, post_id, topic_id)
- **Timestamps** on all records (created_at, updated_at) for audit trails
- **Check constraints** ensure vote_type is either 1 (upvote) or -1 (downvote)

**Testing:**
- Unit tests for service layer (`login_test.go`)
- Integration tests for API handlers (`api_test.go`)
- Repository tests for database operations (`repository_test.go`)

---

## AI Usage Disclosure

I acknowledge using AI assistance (GitHub Copilot with Claude Sonnet 4.5) for specific supporting tasks:

### AI-Assisted Components (~15-20% of project)

1. **Unit Test Scaffolding**: AI helped generate test case structures and mock data for backend unit tests
2. **Seed Data Content**: AI generated realistic forum content (usernames, topic descriptions, post titles, comment text) for the seed script
3. **Deployment Troubleshooting**: AI assisted with debugging specific deployment errors:
   - SSL connection configuration for Render PostgreSQL
   - CORS middleware setup for Vercel frontend
   - Vercel SPA routing configuration (`vercel.json`)
   - Environment variable configuration
4. **Documentation Formatting**: AI helped structure and format this README file

### Independently Developed Core Features (~80-85% of project)

All major features, architecture decisions, and implementations were designed and coded by me:

**Backend:**
- Three-layer architecture design (API → Service → Data)
- Complete database schema with relationships and triggers
- JWT authentication system with middleware
- All API endpoints and handlers (topics, posts, comments, votes, users)
- Business logic in service layer
- Repository pattern implementation
- Vote counting trigger logic
- Migration files (3 migrations)
- CORS and security configuration

**Frontend:**
- Redux Toolkit state management architecture
- All Redux slices (auth, topics, posts, comments, profiles)
- TypeScript type definitions
- Complete component library (15+ components)
- All page components with routing
- Material UI theming and customization
- Dark/light theme implementation
- API client configuration with interceptors
- Form handling and validation
- Responsive design and layout

**Integration & Features:**
- Frontend-backend API integration
- Voting system (frontend + backend)
- Comment threading logic
- User profile pages
- Sorting and pagination
- Protected routes and auth guards
- Error handling and loading states

AI was largely used as a productivity tool and for when I couldn't figure out specific errors during deployment, but all architectural decisions, feature implementations, database design, and core application logic were entirely my own work.

---

## API Documentation

### Base URL
```
Production: https://gossip-with-go.onrender.com/api/v1
Local: http://localhost:8080/api/v1
```

### Authentication Endpoints

#### Register User
```http
POST /register
Content-Type: application/json

{
  "username": "string",
  "email": "string",
  "password": "string"
}

Response: 201 Created
{
  "user_id": 1,
  "username": "string",
  "email": "string",
  "created_at": "2026-01-22T10:00:00Z"
}
```

#### Login
```http
POST /login
Content-Type: application/json

{
  "username": "string",
  "password": "string"
}

Response: 200 OK
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "user_id": 1,
    "username": "string",
    "email": "string"
  }
}
```

### Topic Endpoints

#### Get All Topics
```http
GET /topics

Response: 200 OK
[
  {
    "topic_id": 1,
    "title": "Technology Trends 2026",
    "description": "Discuss latest innovations...",
    "created_by": 1,
    "username": "alice_wonder",
    "post_count": 20,
    "created_at": "2026-01-01T10:00:00Z",
    "updated_at": "2026-01-22T10:00:00Z"
  }
]
```

#### Get Topic by ID
```http
GET /topics/:topicID

Response: 200 OK
{
  "topic_id": 1,
  "title": "Technology Trends 2026",
  "description": "Discuss latest innovations...",
  "created_by": 1,
  "username": "alice_wonder",
  "post_count": 20,
  "created_at": "2026-01-01T10:00:00Z"
}
```

#### Create Topic (Auth Required)
```http
POST /topics
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "string",
  "description": "string"
}

Response: 201 Created
{
  "topic_id": 2,
  "title": "string",
  "description": "string",
  "created_by": 5,
  "created_at": "2026-01-22T10:00:00Z"
}
```

#### Update Topic (Auth Required)
```http
PUT /topics/:topicID
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "Updated Title",
  "description": "Updated description"
}
```

#### Delete Topic (Auth Required)
```http
DELETE /topics/:topicID
Authorization: Bearer {token}

Response: 200 OK
{
  "message": "Topic deleted successfully"
}
```

### Post Endpoints

#### Get Posts by Topic (Optional Auth)
```http
GET /topics/:topicID/posts
Authorization: Bearer {token}  // Optional

Response: 200 OK
[
  {
    "post_id": 1,
    "topic_id": 1,
    "title": "Getting Started Guide",
    "content": "This is a comprehensive guide...",
    "vote_count": 15,
    "comment_count": 8,
    "created_by": 2,
    "username": "bob_builder",
    "user_vote": 1,  // Present if authenticated (1, -1, or 0)
    "created_at": "2026-01-20T10:00:00Z",
    "updated_at": "2026-01-22T10:00:00Z"
  }
]
```

#### Get Post by ID (Optional Auth)
```http
GET /posts/:postID
Authorization: Bearer {token}  // Optional

Response: 200 OK
{
  "post_id": 1,
  "topic_id": 1,
  "title": "Getting Started Guide",
  "content": "This is a comprehensive guide...",
  "vote_count": 15,
  "comment_count": 8,
  "created_by": 2,
  "username": "bob_builder",
  "user_vote": 1,  // Present if authenticated
  "created_at": "2026-01-20T10:00:00Z"
}
```

#### Create Post (Auth Required)
```http
POST /topics/:topicID/posts
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "My New Post",
  "content": "Post content here..."
}

Response: 201 Created
{
  "post_id": 101,
  "topic_id": 1,
  "title": "My New Post",
  "content": "Post content here...",
  "vote_count": 0,
  "created_by": 5,
  "created_at": "2026-01-22T10:00:00Z"
}
```

#### Update Post (Auth Required)
```http
PUT /posts/:postID
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "Updated Title",
  "content": "Updated content"
}
```

#### Delete Post (Auth Required)
```http
DELETE /posts/:postID
Authorization: Bearer {token}

Response: 200 OK
{
  "message": "Post deleted successfully"
}
```

### Comment Endpoints

#### Get Comments by Post (Optional Auth)
```http
GET /posts/:postID/comments
Authorization: Bearer {token}  // Optional

Response: 200 OK
[
  {
    "comment_id": 1,
    "post_id": 1,
    "parent_id": null,
    "content": "Great post! This really helped me.",
    "vote_count": 5,
    "created_by": 3,
    "username": "charlie_brown",
    "user_vote": 1,  // Present if authenticated
    "created_at": "2026-01-21T10:00:00Z",
    "updated_at": "2026-01-22T10:00:00Z"
  },
  {
    "comment_id": 2,
    "post_id": 1,
    "parent_id": 1,  // Reply to comment 1
    "content": "I agree!",
    "vote_count": 2,
    "created_by": 4,
    "username": "diana_prince",
    "user_vote": 0,
    "created_at": "2026-01-21T11:00:00Z"
  }
]
```

#### Create Comment (Auth Required)
```http
POST /posts/:postID/comments
Authorization: Bearer {token}
Content-Type: application/json

{
  "content": "This is my comment",
  "parent_id": null  // Optional: set to comment_id for nested reply
}

Response: 201 Created
{
  "comment_id": 150,
  "post_id": 1,
  "parent_id": null,
  "content": "This is my comment",
  "vote_count": 0,
  "created_by": 5,
  "created_at": "2026-01-22T10:00:00Z"
}
```

#### Update Comment (Auth Required)
```http
PUT /comments/:commentID
Authorization: Bearer {token}
Content-Type: application/json

{
  "content": "Updated comment content"
}
```

#### Delete Comment (Auth Required)
```http
DELETE /comments/:commentID
Authorization: Bearer {token}

Response: 200 OK
{
  "message": "Comment deleted successfully"
}
```

### Vote Endpoints

#### Vote on Post (Auth Required)
```http
POST /posts/:postID/vote
Authorization: Bearer {token}
Content-Type: application/json

{
  "vote_type": 1  // 1 for upvote, -1 for downvote
}

Response: 200 OK
{
  "message": "Vote recorded",
  "vote_count": 16  // Updated vote count
}
```

#### Remove Vote from Post (Auth Required)
```http
DELETE /posts/:postID/vote
Authorization: Bearer {token}

Response: 200 OK
{
  "message": "Vote removed",
  "vote_count": 15
}
```

#### Vote on Comment (Auth Required)
```http
POST /comments/:commentID/vote
Authorization: Bearer {token}
Content-Type: application/json

{
  "vote_type": -1  // 1 for upvote, -1 for downvote
}
```

#### Remove Vote from Comment (Auth Required)
```http
DELETE /comments/:commentID/vote
Authorization: Bearer {token}
```

### User Profile Endpoints

#### Get User by ID (Auth Required)
```http
GET /users/:id
Authorization: Bearer {token}

Response: 200 OK
{
  "user_id": 1,
  "username": "alice_wonder",
  "email": "alice_wonder@example.com",
  "created_at": "2026-01-01T10:00:00Z"
}
```

#### Get User Posts (Auth Required)
```http
GET /users/:id/posts
Authorization: Bearer {token}

Response: 200 OK
[
  {
    "post_id": 1,
    "title": "My Post",
    "vote_count": 10,
    "comment_count": 5,
    "created_at": "2026-01-20T10:00:00Z"
  }
]
```

#### Get User Comments (Auth Required)
```http
GET /users/:id/comments
Authorization: Bearer {token}

Response: 200 OK
[
  {
    "comment_id": 1,
    "post_id": 5,
    "content": "My comment",
    "vote_count": 3,
    "created_at": "2026-01-21T10:00:00Z"
  }
]
```

---

## Troubleshooting

### Backend Issues

**Database connection fails:**
```bash
# Check if PostgreSQL is running
pg_isready

# Verify DATABASE_URL format
echo $DATABASE_URL

# For local: postgresql://user:pass@localhost:5432/db?sslmode=disable
# For Render: postgresql://user:pass@host.render.com:5432/db?sslmode=require
```

**Migration errors:**
```bash
# Check current version
migrate -path migrations -database "${DATABASE_URL}" version

# Roll back one migration
migrate -path migrations -database "${DATABASE_URL}" down 1

# Force to specific version (if stuck)
migrate -path migrations -database "${DATABASE_URL}" force 2
```

**Port already in use:**
```bash
# Find and kill process on port 8080
lsof -ti:8080 | xargs kill -9

# Or use a different port
export PORT=8081
go run cmd/server/main.go
```

**JWT token issues:**
```bash
# Ensure JWT_SECRET is set
echo $JWT_SECRET

# Generate a strong secret
openssl rand -base64 32
```

### Frontend Issues

**CORS errors in browser console:**
- Verify backend CORS middleware allows your frontend origin
- Check that `VITE_API_URL` in `.env` matches backend URL
- For local dev: `http://localhost:8080/api/v1`
- For production: `https://your-backend.onrender.com/api/v1`

**Build errors:**
```bash
# Clear cache and reinstall
rm -rf node_modules package-lock.json
npm install

# Clear Vite cache
rm -rf node_modules/.vite
npm run dev
```

**404 errors on page refresh (production):**
- Ensure `vercel.json` exists with rewrite rules:
```json
{
  "rewrites": [
    { "source": "/(.*)", "destination": "/index.html" }
  ]
}
```

**Redux DevTools not working:**
```bash
# Install Redux DevTools extension in browser
# Chrome: https://chrome.google.com/webstore (search "Redux DevTools")
# Firefox: https://addons.mozilla.org (search "Redux DevTools")
```

**TypeScript errors:**
```bash
# Check for type errors
npm run build

# Update TypeScript if needed
npm install -D typescript@latest
```

### Database Issues

**Seed script timeout:**
- Database might be inaccessible from local machine
- Use the seed API endpoint instead:
```bash
curl -X POST https://your-backend.onrender.com/api/v1/seed \
  -H "X-Seed-Secret: seed-my-database-2026"
```

**Vote counts not updating:**
- Check if triggers exist:
```sql
-- In PostgreSQL shell
\df update_post_vote_count
\df update_comment_vote_count
```

**Duplicate vote errors:**
- This is expected - unique constraint prevents duplicate votes
- Frontend should disable vote button after voting

---

## Production Deployment

### Live Application

```
gossip-with-go-seven.vercel.app
```

### Environment Variables

#### Backend (Render)
```
DATABASE_URL=postgresql://user:password@host/db?sslmode=require
JWT_SECRET=your-production-secret-key-here
PORT=8080  # Automatically set by Render
```

#### Frontend (Vercel)
```
VITE_API_URL=https://gossip-with-go.onrender.com/api/v1
```

### Deployment Process

**Backend (Render):**
1. Push to `main` branch
2. Render automatically builds: `go build -o server cmd/server/main.go`
3. Migrations run automatically on startup (via `runMigrations()` in main.go)
4. Server starts: `./server`

**Frontend (Vercel):**
1. Push to `main` branch
2. Vercel automatically builds: `npm run build`
3. Static files deployed to CDN
4. SPA routing configured via `vercel.json`

---

## Future Enhancements

Potential features for future development:

- [ ] Comment pinning 
- [ ] Real-time updates using WebSockets
- [ ] Rich text editor with Markdown support
- [ ] Image/file uploads for posts
- [ ] User avatars and profile customization
- [ ] Email notifications for replies
- [ ] Admin dashboard and moderation tools
- [ ] Rate limiting and spam protection
- [ ] Bookmarks/favorites system
- [ ] Topic tags and filtering
- [ ] Mobile app (React Native)

---

## License

This project is part of an academic assignment and is for educational purposes only.

---

## Author

**Adzfar**

- GitHub: [@adzzfarr](https://github.com/adzzfarr)
- Project Repository: [gossip-with-go](https://github.com/adzzfarr/gossip-with-go)
