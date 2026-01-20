package data

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository holds all database access methods
type Repository struct {
	DB *pgxpool.Pool // Connection pool created in ./db.go
}

// NewRepository initializes a new instance of Repository struct
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{DB: db}
}

// GetAllTopics fetches all topics from the database
func (repo *Repository) GetAllTopics(sortBy, searchQuery string, page, limit int) (*TopicsResponse, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Offset for pagination
	offset := (page - 1) * limit

	// Determine ORDER BY based on sortBy
	var orderBy string
	switch sortBy {
	case "newest":
		orderBy = "t.created_at DESC"
	case "oldest":
		orderBy = "t.created_at ASC"
	case "most_posts":
		orderBy = "post_count DESC, t.created_at DESC"
	default:
		orderBy = "t.created_at DESC" // Default sorting
	}

	// WHERE clause for search query
	whereClause := ""
	args := []interface{}{}
	argPosition := 1

	if searchQuery != "" {
		whereClause = fmt.Sprintf(`
			WHERE (t.title ILIKE $%d OR t.description ILIKE $%d)`,
			argPosition,
			argPosition,
		)
		args = append(args, "%"+searchQuery+"%")
		argPosition++
	}

	// Get total count for pagination
	countQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT t.topic_id)
		FROM topics t
		%s`, whereClause,
	)

	var totalItems int
	err := repo.DB.QueryRow(ctx, countQuery, args...).Scan(&totalItems)
	if err != nil {
		return nil, fmt.Errorf("failed to get total topics count: %w", err)
	}

	// Get paginated topics
	args = append(args, limit, offset)
	query := fmt.Sprintf(`
		SELECT 
			t.topic_id, 
			t.title, 
			t.description, 
			t.created_by, 
			u.username, 
			t.created_at, 
			t.updated_at,
			COUNT(p.post_id) AS post_count
        FROM topics t
        JOIN users u ON t.created_by = u.user_id
        LEFT JOIN posts p ON t.topic_id = p.topic_id
		%s
        GROUP BY t.topic_id, u.username
        ORDER BY %s
		LIMIT $%d 
		OFFSET $%d`,
		whereClause,
		orderBy,
		argPosition,
		argPosition+1,
	)

	rows, err := repo.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query all topics failed: %w", err)
	}
	defer rows.Close()

	// Scan Results
	topics := []*Topic{} // Initialize empty slice of Topic pointers
	for rows.Next() {
		var topic Topic

		// Scan column values from current row into fields of the Topic struct (must match SELECT order)
		err := rows.Scan(
			&topic.TopicID,
			&topic.Title,
			&topic.Description,
			&topic.CreatedBy,
			&topic.Username,
			&topic.CreatedAt,
			&topic.UpdatedAt,
			&topic.PostCount,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning topic row: %w", err)
		}

		topics = append(topics, &topic) // Append a pointer to the Topic to the slice
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error encountered during row iteration: %w", err)
	}

	// Calculate pagination metadata
	totalPages := (totalItems + limit - 1) / limit // Divide and round up
	if totalPages == 0 {
		totalPages = 1
	}

	pagination := PaginationMetadata{
		CurrentPage: page,
		TotalPages:  totalPages,
		PageSize:    limit,
		TotalItems:  totalItems,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}

	return &TopicsResponse{
		Topics:     topics,
		Pagination: pagination,
	}, nil
}

func (repo *Repository) GetTopicByID(topicID int) (*Topic, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var topic Topic
	query := `
		SELECT t.topic_id, t.title, t.description, t.created_by, u.username, t.created_at, t.updated_at
        FROM topics t
        JOIN users u ON t.created_by = u.user_id
		WHERE t.topic_id = $1`

	err := repo.DB.QueryRow(ctx, query, topicID).Scan(
		&topic.TopicID,
		&topic.Title,
		&topic.Description,
		&topic.CreatedBy,
		&topic.Username,
		&topic.CreatedAt,
		&topic.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("topic not found: %d", topicID)
		}
		return nil, fmt.Errorf("query to find topic failed: %w", err)
	}

	return &topic, nil
}

// GetUserByUsername fetches user by their unique username
// Used to check if a user exists (during registration) and to retrieve credentials (during login)
func (repo *Repository) GetUserByUsername(username string) (*User, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var user User
	query := `
    	SELECT user_id, username, password_hash, created_at, updated_at
        FROM users
        WHERE username = $1`

	err := repo.DB.QueryRow(ctx, query, username).Scan(
		&user.UserID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found: %s", username)
		}
		return nil, fmt.Errorf("query to find user failed: %w", err)
	}

	// Return pointer to the found User
	return &user, nil
}

// CreateUser inserts a new user record into the database
// NOTE: Password MUST already be hashed (in service layer) before this function is called
func (repo *Repository) CreateUser(user *User) (*User, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query := `
        INSERT INTO users (username, password_hash, created_at, updated_at)
        VALUES ($1, $2, NOW(), NOW())
        RETURNING user_id, created_at, updated_at`

	err := repo.DB.QueryRow(
		ctx,
		query,
		user.Username,
		user.PasswordHash,
	).Scan(
		&user.UserID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// CreateTopic inserts a new topic into the database
func (repo *Repository) CreateTopic(title, description string, userID int) (*Topic, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query := `
		INSERT INTO topics (title, description, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING topic_id, title, description, created_by, created_at, updated_at
	`

	// Scan returned row into Topic struct
	var topic Topic
	err := repo.DB.QueryRow(
		ctx,
		query,
		title,
		description,
		userID,
	).Scan(
		&topic.TopicID,
		&topic.Title,
		&topic.Description,
		&topic.CreatedBy,
		&topic.CreatedAt,
		&topic.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create topic: %w", err)
	}

	// Fetch username of topic creator
	var username string
	userQuery := `
		SELECT username
		FROM users
		WHERE user_id = $1`

	err = repo.DB.QueryRow(ctx, userQuery, userID).Scan(&username)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch username: %w", err)
	}

	topic.Username = username

	return &topic, nil
}

// GetPostsByTopicID fetches all posts for a given topic ID
func (repo *Repository) GetPostsByTopicID(topicID int, userID *int, sortBy, searchQuery string, page, limit int) (*PostsResponse, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	offset := (page - 1) * limit

	var orderBy string
	switch sortBy {
	case "hot":
		// Hot Score: (votes) / (age in hours + 2)^1.5
		orderBy = `(
			COALESCE(p.vote_count, 0) - (EXTRACT(EPOCH FROM (NOW() - p.created_at)) / 86400)
		) DESC`
	case "newest":
		orderBy = "p.created_at DESC"
	case "oldest":
		orderBy = "p.created_at ASC"
	case "most_voted":
		orderBy = "p.vote_count DESC, p.created_at DESC"
	default:
		orderBy = `(
			COALESCE(p.vote_count, 0) - (EXTRACT(EPOCH FROM (NOW() - p.created_at)) / 86400)
		) DESC`
	}

	whereClause := "WHERE p.topic_id = $1"
	args := []interface{}{topicID}
	argPosition := 2

	if searchQuery != "" {
		whereClause = fmt.Sprintf(`
			AND (p.title ILIKE $%d OR p.content ILIKE $%d)`,
			argPosition,
			argPosition,
		)
		args = append(args, "%"+searchQuery+"%")
		argPosition++
	}

	// Add userID parameter
	args = append(args, userID)
	userIDPosition := argPosition
	argPosition++

	countQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT p.post_id)
		FROM posts p
		%s`,
		whereClause,
	)

	var totalItems int
	countArgs := args[:len(args)-1] // Exclude userID for count query
	err := repo.DB.QueryRow(ctx, countQuery, countArgs...).Scan(&totalItems)
	if err != nil {
		return nil, fmt.Errorf("failed to get total posts count: %w", err)
	}

	args = append(args, limit, offset)
	query := fmt.Sprintf(`
		SELECT 
			p.post_id, 
			p.topic_id, 
			t.title as topic_title, 
			p.title, 
			p.content, 
			p.created_by, 
			u.username, 
			p.created_at, 
			p.updated_at,
			p.vote_count,
			CASE 
				WHEN $%d::integer IS NOT NULL THEN (
					SELECT vote_type FROM votes 
					WHERE user_id = $%d AND post_id = p.post_id
				)
				ELSE NULL
			END AS user_vote
		FROM posts p
		JOIN users u ON p.created_by = u.user_id
		JOIN topics t ON p.topic_id = t.topic_id
		%s
		ORDER BY %s
		LIMIT $%d
		OFFSET $%d
		`,
		userIDPosition,
		userIDPosition,
		whereClause,
		orderBy,
		argPosition,
		argPosition+1,
	)

	rows, err := repo.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query posts: %w", err)
	}
	defer rows.Close()

	posts := []*Post{}
	for rows.Next() {
		var post Post

		err := rows.Scan(
			&post.PostID,
			&post.TopicID,
			&post.TopicTitle,
			&post.Title,
			&post.Content,
			&post.CreatedBy,
			&post.Username,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.VoteCount,
			&post.UserVote,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan post row: %w", err)
		}

		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error encountered during row iteration: %w", err)
	}

	totalPages := (totalItems + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}

	pagination := PaginationMetadata{
		CurrentPage: page,
		TotalPages:  totalPages,
		PageSize:    limit,
		TotalItems:  totalItems,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}

	return &PostsResponse{
		Posts:      posts,
		Pagination: pagination,
	}, nil
}

// GetPostByID fetches a specific post by its ID
func (repo *Repository) GetPostByID(postID int, userID *int) (*Post, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var post Post
	query := `
		SELECT 
			p.post_id, 
			p.topic_id,
			t.title as topic_title, 
			p.title, 
			p.content, 
			p.created_by, 
			u.username, 
			p.created_at, 
			p.updated_at,
			p.vote_count,
			CASE
				WHEN $2::integer IS NOT NULL THEN (
					SELECT vote_type FROM votes
					WHERE user_id = $2 AND post_id = p.post_id
				)
				ELSE NULL
			END AS user_vote
		FROM posts p
		JOIN users u ON p.created_by = u.user_id
		JOIN topics t ON p.topic_id = t.topic_id
		WHERE p.post_id = $1`

	err := repo.DB.QueryRow(ctx, query, postID, userID).Scan(
		&post.PostID,
		&post.TopicID,
		&post.TopicTitle,
		&post.Title,
		&post.Content,
		&post.CreatedBy,
		&post.Username,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.VoteCount,
		&post.UserVote,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("post not found with ID: %d", postID)
		}
		return nil, fmt.Errorf("query to find post failed: %w", err)
	}

	return &post, nil
}

// buildCommentTree constructs a tree of comments from a flat list
func buildCommentTree(comments []*Comment) []*Comment {
	// Map for quick lookup by commentID
	commentMap := make(map[int]*Comment)
	topLevelComments := []*Comment{}

	// Initialise all comments in map
	for _, comment := range comments {
		comment.Replies = []*Comment{}
		commentMap[comment.CommentID] = comment
	}

	// Assign children to parents
	for _, comment := range comments {
		if comment.ParentCommentID == nil {
			// Top-level comment
			comment.Depth = 0
			topLevelComments = append(topLevelComments, comment)
		} else {
			// Reply to another comments
			parentComment, exists := commentMap[*comment.ParentCommentID]

			if exists {
				comment.Depth = parentComment.Depth + 1
				parentComment.Replies = append(parentComment.Replies, comment)
			} else {
				// Orphaned comment, treat as top-level
				comment.Depth = 0
				topLevelComments = append(topLevelComments, comment)
			}
		}
	}

	return topLevelComments
}

// GetCommentsByPostID fetches all comments for a given post ID
func (repo *Repository) GetCommentsByPostID(postID int, userID *int, sortBy, searchQuery string, page, limit int) (*CommentsResponse, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	offset := (page - 1) * limit

	var orderBy string
	switch sortBy {
	case "hot":
		// Hot Score: (votes) / (age in hours + 2)^1.5
		orderBy = `(
			COALESCE(c.vote_count, 0) - (EXTRACT(EPOCH FROM (NOW() - c.created_at)) / 86400)
		) DESC`
	case "newest":
		orderBy = "c.created_at DESC"
	case "oldest":
		orderBy = "c.created_at ASC"
	case "most_voted":
		orderBy = "c.vote_count DESC, c.created_at DESC"
	default:
		orderBy = `(
			COALESCE(c.vote_count, 0) - (EXTRACT(EPOCH FROM (NOW() - c.created_at)) / 86400)
		) DESC`
	}

	whereClause := "WHERE c.post_id = $1"
	args := []interface{}{postID}
	argPosition := 2

	if searchQuery != "" {
		whereClause += fmt.Sprintf(`
			AND (c.content ILIKE $%d)`,
			argPosition,
		)
		args = append(args, "%"+searchQuery+"%")
		argPosition++
	}

	// Add userID parameter
	args = append(args, userID)
	userIDPosition := argPosition
	argPosition++

	countQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT c.comment_id)
		FROM comments c
		%s`,
		whereClause,
	)

	var totalItems int
	countArgs := args[:len(args)-1] // Exclude userID for count query
	err := repo.DB.QueryRow(ctx, countQuery, countArgs...).Scan(&totalItems)
	if err != nil {
		return nil, fmt.Errorf("failed to get total comments count: %w", err)
	}

	args = append(args, limit, offset)
	query := fmt.Sprintf(`
		SELECT 
			c.comment_id, 
			c.post_id, 
			c.parent_comment_id,
			c.content, 
			c.created_by, 
			u.username, 
			c.created_at, 
			c.updated_at,
			c.vote_count,
			CASE
				WHEN $%d::integer IS NOT NULL THEN (
					SELECT vote_type FROM votes
					WHERE user_id = $%d AND comment_id = c.comment_id
				)
				ELSE NULL
			END AS user_vote
		FROM comments c
		JOIN users u ON c.created_by = u.user_id
		%s
		ORDER BY %s
		LIMIT $%d
		OFFSET $%d
		`,
		userIDPosition,
		userIDPosition,
		whereClause,
		orderBy,
		argPosition,
		argPosition+1,
	)

	rows, err := repo.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}
	defer rows.Close()

	// Initialise flat list
	comments := []*Comment{}
	for rows.Next() {
		var comment Comment

		err := rows.Scan(
			&comment.CommentID,
			&comment.PostID,
			&comment.ParentCommentID,
			&comment.Content,
			&comment.CreatedBy,
			&comment.Username,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.VoteCount,
			&comment.UserVote,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan comment row: %w", err)
		}

		comments = append(comments, &comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error encountered during row iteration: %w", err)
	}

	// Build tree structure from flat list
	comments = buildCommentTree(comments)

	totalPages := (totalItems + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}

	pagination := PaginationMetadata{
		CurrentPage: page,
		TotalPages:  totalPages,
		PageSize:    limit,
		TotalItems:  totalItems,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}

	return &CommentsResponse{
		Comments:   comments,
		Pagination: pagination,
	}, nil
}

// GetCommentByID fetches a specific comment by its ID
func (repo *Repository) GetCommentByID(commentID int, userID *int) (*Comment, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var comment Comment
	query := `
		SELECT
			c.comment_id,
			c.post_id,
			c.parent_comment_id,
			p.title as post_title,
			c.content,
			c.created_by,
			u.username,
			c.created_at,
			c.updated_at,
			c.vote_count,
			CASE
				WHEN $2::integer IS NOT NULL THEN (
					SELECT vote_type FROM votes
					WHERE user_id = $2 AND comment_id = c.comment_id
				)
				ELSE NULL
			END AS user_vote
		FROM comments c
		JOIN users u ON c.created_by = u.user_id
		JOIN posts p ON c.post_id = p.post_id
		WHERE c.comment_id = $1`

	err := repo.DB.QueryRow(ctx, query, commentID, userID).Scan(
		&comment.CommentID,
		&comment.PostID,
		&comment.ParentCommentID,
		&comment.PostTitle,
		&comment.Content,
		&comment.CreatedBy,
		&comment.Username,
		&comment.CreatedAt,
		&comment.UpdatedAt,
		&comment.VoteCount,
		&comment.UserVote,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("comment not found with ID: %d", commentID)
		}
		return nil, fmt.Errorf("query to find comment failed: %w", err)
	}

	// Initialise empty replies
	comment.Replies = []*Comment{}

	return &comment, nil
}

// CreatePost inserts a new post into the database
func (repo *Repository) CreatePost(topicID int, title, content string, userID int) (*Post, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query := `
		INSERT INTO posts (topic_id, title, content, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING post_id, topic_id, title, content, created_by, created_at, updated_at`

	var post Post
	err := repo.DB.QueryRow(
		ctx,
		query,
		topicID,
		title,
		content,
		userID,
	).Scan(
		&post.PostID,
		&post.TopicID,
		&post.Title,
		&post.Content,
		&post.CreatedBy,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	// Fetch username of post creator
	var username string
	userQuery := `
		SELECT username
		FROM users
		WHERE user_id = $1
	`

	err = repo.DB.QueryRow(ctx, userQuery, userID).Scan(&username)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch username: %w", err)
	}

	post.Username = username

	return &post, nil
}

// CreateComment inserts a new comment into the database
func (repo *Repository) CreateComment(postID int, content string, userID int, parentCommentID *int) (*Comment, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query := `
		INSERT INTO comments (post_id, parent_comment_id, content, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING comment_id, post_id, parent_comment_id, content, created_by, created_at, updated_at`

	var comment Comment
	err := repo.DB.QueryRow(
		ctx,
		query,
		postID,
		parentCommentID,
		content,
		userID,
	).Scan(
		&comment.CommentID,
		&comment.PostID,
		&comment.ParentCommentID,
		&comment.Content,
		&comment.CreatedBy,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	// Fetch username of comment creator
	var username string
	userQuery := `
		SELECT username
		FROM users
		WHERE user_id = $1
	`

	err = repo.DB.QueryRow(ctx, userQuery, userID).Scan(&username)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch username: %w", err)
	}

	comment.Username = username
	comment.VoteCount = 0
	comment.UserVote = nil
	comment.Replies = []*Comment{}

	return &comment, nil
}

// UpdateTopic updates an existing topic's title and description
func (repo *Repository) UpdateTopic(topicID int, title, description string, userID int) (*Topic, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Verify that topic exists and was created by the user
	var creatorID int

	checkQuery := `
		SELECT created_by
		FROM topics
		WHERE topic_id = $1`

	err := repo.DB.QueryRow(
		ctx,
		checkQuery,
		topicID,
	).Scan(&creatorID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("topic with ID %d not found", topicID)
		}

		return nil, fmt.Errorf("failed to verify topic ownership: %w", err)
	}

	if creatorID != userID {
		return nil, fmt.Errorf("user %d is not authorized to update topic %d", userID, topicID)
	}

	// Update topic
	query := `
		UPDATE topics
		SET title = $1, description = $2, updated_at = NOW()
		WHERE topic_id = $3 AND created_by = $4
		RETURNING 
			topic_id, 
			title, 
			description, 
			created_by, 
			(SELECT username FROM users WHERE user_id = $4) AS username,
			created_at, 
			updated_at`

	var updatedTopic Topic
	err = repo.DB.QueryRow(
		ctx,
		query,
		title,
		description,
		topicID,
		userID,
	).Scan(
		&updatedTopic.TopicID,
		&updatedTopic.Title,
		&updatedTopic.Description,
		&updatedTopic.CreatedBy,
		&updatedTopic.Username,
		&updatedTopic.CreatedAt,
		&updatedTopic.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update topic: %w", err)
	}

	return &updatedTopic, nil
}

// UpdatePost updates an existing post's title and content
func (repo *Repository) UpdatePost(postID int, title, content string, userID int) (*Post, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Verify that post exists and was created by the user
	var creatorID int

	checkQuery := `
		SELECT created_by
		FROM posts
		WHERE post_id = $1`

	err := repo.DB.QueryRow(
		ctx,
		checkQuery,
		postID,
	).Scan(&creatorID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("post with ID %d not found", postID)
		}

		return nil, fmt.Errorf("failed to verify post ownership: %w", err)
	}

	if creatorID != userID {
		return nil, fmt.Errorf("user %d is not authorized to update post %d", userID, postID)
	}

	// Update post
	query := `
		UPDATE posts
		SET title = $1, content = $2, updated_at = NOW()
		WHERE post_id = $3 AND created_by = $4
		RETURNING 
			post_id, 
			topic_id, 
			title, 
			content, 
			created_by, 
			(SELECT username FROM users WHERE user_id = $4) AS username,
			created_at, 
			updated_at`

	var updatedPost Post
	err = repo.DB.QueryRow(
		ctx,
		query,
		title,
		content,
		postID,
		userID,
	).Scan(
		&updatedPost.PostID,
		&updatedPost.TopicID,
		&updatedPost.Title,
		&updatedPost.Content,
		&updatedPost.CreatedBy,
		&updatedPost.Username,
		&updatedPost.CreatedAt,
		&updatedPost.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	return &updatedPost, nil
}

// UpdateComment updates an existing comment's content
func (repo *Repository) UpdateComment(commentID int, content string, userID int) (*Comment, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Verify that comment exists and was created by the user
	var creatorID int

	checkQuery := `
		SELECT created_by
		FROM comments
		WHERE comment_id = $1`

	err := repo.DB.QueryRow(
		ctx,
		checkQuery,
		commentID,
	).Scan(&creatorID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("comment with ID %d not found", commentID)
		}

		return nil, fmt.Errorf("failed to verify comment ownership: %w", err)
	}

	if creatorID != userID {
		return nil, fmt.Errorf("user %d is not authorized to update comment %d", userID, commentID)
	}

	// Update comment
	query := `
		UPDATE comments
		SET content = $1, updated_at = NOW()
		WHERE comment_id = $2 AND created_by = $3
		RETURNING 
			comment_id, 
			post_id, 
			content, 
			created_by, 
			(SELECT username FROM users WHERE user_id = $3) AS username,
			created_at, 
			updated_at`

	var updatedComment Comment
	err = repo.DB.QueryRow(
		ctx,
		query,
		content,
		commentID,
		userID,
	).Scan(
		&updatedComment.CommentID,
		&updatedComment.PostID,
		&updatedComment.Content,
		&updatedComment.CreatedBy,
		&updatedComment.Username,
		&updatedComment.CreatedAt,
		&updatedComment.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	return &updatedComment, nil
}

// DeleteComment deletes an existing comment
func (repo *Repository) DeleteComment(commentID, userID int) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Verify that comment exists and was created by the user
	var creatorID int

	checkQuery := `
		SELECT created_by
		FROM comments
		WHERE comment_id = $1`

	err := repo.DB.QueryRow(
		ctx,
		checkQuery,
		commentID,
	).Scan(&creatorID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("comment with ID %d not found", commentID)
		}

		return fmt.Errorf("failed to verify comment ownership: %w", err)
	}

	if creatorID != userID {
		return fmt.Errorf("user %d is not authorized to delete comment %d", userID, commentID)
	}

	// Delete comment
	query := `
		DELETE FROM comments
		WHERE comment_id = $1 AND created_by = $2`

	commandTag, err := repo.DB.Exec(
		ctx,
		query,
		commentID,
		userID,
	)

	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("comment with ID %d not found or not owned by user %d", commentID, userID)
	}

	return nil
}

// DeletePost deletes an existing post and its comments
func (repo *Repository) DeletePost(postID, userID int) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Verify that post exists and was created by the user
	var creatorID int

	checkQuery := `
		SELECT created_by
		FROM posts
		WHERE post_id = $1`

	err := repo.DB.QueryRow(
		ctx,
		checkQuery,
		postID,
	).Scan(&creatorID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("post with ID %d not found", postID)
		}

		return fmt.Errorf("failed to verify post ownership: %w", err)
	}

	if creatorID != userID {
		return fmt.Errorf("user %d is not authorized to delete post %d", userID, postID)
	}

	// Delete post (comments deleted automatically via CASCADE)
	query := `
		DELETE FROM posts
		WHERE post_id = $1 AND created_by = $2`

	_, err = repo.DB.Exec(
		ctx,
		query,
		postID,
		userID,
	)

	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

// DeleteTopic deletes an existing topic, including its posts and their comments
func (repo *Repository) DeleteTopic(topicID, userID int) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Verify that topic exists and was created by the user
	var creatorID int

	checkQuery := `
		SELECT created_by
		FROM topics
		WHERE topic_id = $1`

	err := repo.DB.QueryRow(
		ctx,
		checkQuery,
		topicID,
	).Scan(&creatorID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("topic with ID %d not found", topicID)
		}

		return fmt.Errorf("failed to verify topic ownership: %w", err)
	}

	if creatorID != userID {
		return fmt.Errorf("user %d is not authorized to delete topic %d", userID, topicID)
	}

	// Delete topic (all posts and their comments delete automatically via CASCADE)
	query := `
		DELETE FROM topics
		WHERE topic_id = $1 AND created_by = $2`

	_, err = repo.DB.Exec(
		ctx,
		query,
		topicID,
		userID,
	)

	if err != nil {
		return fmt.Errorf("failed to delete topic: %w", err)
	}

	return nil
}

// GetUserByID fetches user by their unique user ID
// Used internally when we need to get user details by ID
func (repo *Repository) GetUserByID(userID int) (*User, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var user User
	query := `
		SELECT user_id, username, created_at, updated_at
		FROM users
		WHERE user_id = $1`

	err := repo.DB.QueryRow(ctx, query, userID).Scan(
		&user.UserID,
		&user.Username,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user with ID %d not found", userID)
		}

		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return &user, nil
}

// GetUserPosts fetches all posts created by a specific user
func (repo *Repository) GetUserPosts(userID int) ([]*Post, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query := `
		SELECT p.post_id, p.topic_id, t.title as topic_title, p.title, p.content, p.created_by, u.username, p.created_at, p.updated_at
		FROM posts p
		JOIN topics t ON p.topic_id = t.topic_id
		JOIN users u ON p.created_by = u.user_id
		WHERE p.created_by = $1
		ORDER BY p.created_at DESC`

	rows, err := repo.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user posts: %w", err)
	}
	defer rows.Close()

	posts := []*Post{}
	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.PostID,
			&post.TopicID,
			&post.TopicTitle,
			&post.Title,
			&post.Content,
			&post.CreatedBy,
			&post.Username,
			&post.CreatedAt,
			&post.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan user post: %w", err)
		}
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error encountered during row iteration: %w", err)
	}

	return posts, nil
}

// GetUserComments fetches all comments created by a specific user
func (repo *Repository) GetUserComments(userID int) ([]*Comment, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query := `
		SELECT c.comment_id, c.post_id, p.title as post_title, c.content, c.created_by, u.username, c.created_at, c.updated_at
		FROM comments c
		JOIN users u ON c.created_by = u.user_id
		JOIN posts p ON c.post_id = p.post_id
		WHERE c.created_by = $1
		ORDER BY c.created_at DESC`

	rows, err := repo.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user comments: %w", err)
	}
	defer rows.Close()

	comments := []*Comment{}
	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&comment.CommentID,
			&comment.PostID,
			&comment.PostTitle,
			&comment.Content,
			&comment.CreatedBy,
			&comment.Username,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan user comment: %w", err)
		}
		comments = append(comments, &comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error encountered during row iteration: %w", err)
	}

	return comments, nil
}

// VotePost creates/updates a vote on a post
func (repo *Repository) VotePost(userID, postID, voteType int) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Validate voteType
	if voteType != 1 && voteType != -1 {
		return fmt.Errorf("invalid vote type: %d", voteType)
	}

	query := `
		INSERT INTO votes (user_id, post_id, vote_type, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		ON CONFLICT (user_id, post_id) 
		DO UPDATE SET vote_type = $3, updated_at = NOW()`

	_, err := repo.DB.Exec(ctx, query, userID, postID, voteType)
	if err != nil {
		return fmt.Errorf("failed to vote on post: %w", err)
	}

	return nil
}

// RemovePostVote removes a user's vote from a post
func (repo *Repository) RemovePostVote(userID, postID int) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query := `
		DELETE FROM votes
		WHERE user_id = $1 AND post_id = $2`

	commandTag, err := repo.DB.Exec(ctx, query, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to remove vote from post: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("no vote found for user %d on post %d", userID, postID)
	}

	return nil
}

// VoteComment creates/updates a vote on a comment
func (repo *Repository) VoteComment(userID, commentID, voteType int) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Validate voteType
	if voteType != 1 && voteType != -1 {
		return fmt.Errorf("invalid vote type: %d", voteType)
	}

	query := `
		INSERT INTO votes (user_id, comment_id, vote_type, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		ON CONFLICT (user_id, comment_id) 
		DO UPDATE SET vote_type = $3, updated_at = NOW()`

	_, err := repo.DB.Exec(ctx, query, userID, commentID, voteType)
	if err != nil {
		return fmt.Errorf("failed to vote on comment: %w", err)
	}

	return nil
}

// RemoveCommentVote removes a user's vote from a comment
func (repo *Repository) RemoveCommentVote(userID, commentID int) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query := `
		DELETE FROM votes
		WHERE user_id = $1 AND comment_id = $2`

	commandTag, err := repo.DB.Exec(ctx, query, userID, commentID)
	if err != nil {
		return fmt.Errorf("failed to remove vote from comment: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("no vote found for user %d on comment %d", userID, commentID)
	}

	return nil
}

// GetUserVoteOnPost gets current user's vote on a specific post, if any
func (repo *Repository) GetUserVoteOnPost(userID, postID int) (*int, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var voteType int
	query := `
		SELECT vote_type
		FROM votes
		WHERE user_id = $1 AND post_id = $2`

	err := repo.DB.QueryRow(ctx, query, userID, postID).Scan(&voteType)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // No vote found
		}
		return nil, fmt.Errorf("failed to get user vote on post: %w", err)
	}

	return &voteType, nil
}

// GetUserVoteOnComment gets current user's vote on a specific comment, if any
func (repo *Repository) GetUserVoteOnComment(userID, commentID int) (*int, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var voteType int
	query := `
		SELECT vote_type
		FROM votes
		WHERE user_id = $1 AND comment_id = $2`

	err := repo.DB.QueryRow(ctx, query, userID, commentID).Scan(&voteType)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user vote on comment: %w", err)
	}

	return &voteType, nil
}
