package handlers

import (
	"database/sql"
	"forum/structs"
	"strconv"
	"strings"
)

func GetAllPosts(db *sql.DB, tagID int, user structs.User) ([]structs.Post, error) {
	if tagID == 98 {
		return getPostsCreatedByUser(db, user.Username)
	} else if tagID == 99 {
		return getPostsLikedByUser(db, user.ID)
	}

	query := `
		SELECT p.postID, u.username, p.title, p.description, p.creationDate, GROUP_CONCAT(t.name), p.likes, p.dislikes
		FROM posts p
		JOIN users u ON p.username = u.username
		LEFT JOIN post_tags pt ON p.postID = pt.postID
		LEFT JOIN tags t ON pt.tagID = t.id
		GROUP BY p.postID
	`

	if tagID > 0 {
		query += `
			HAVING GROUP_CONCAT(t.id) LIKE '%` + strconv.Itoa(tagID) + `%'
		`
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []structs.Post{}
	for rows.Next() {
		var post structs.Post
		var tags string
		if err = rows.Scan(&post.ID, &post.Username, &post.Title, &post.Description, &post.CreationDate, &tags, &post.Likes, &post.Dislikes); err != nil {
			return nil, err
		}
		post.Tags = strings.Split(tags, ",")
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func getPostsCreatedByUser(db *sql.DB, username string) ([]structs.Post, error) {
	query := `
		SELECT p.postID, u.username, p.title, p.description, p.creationDate, GROUP_CONCAT(t.name), p.likes, p.dislikes
		FROM posts p
		JOIN users u ON p.username = u.username
		LEFT JOIN post_tags pt ON p.postID = pt.postID
		LEFT JOIN tags t ON pt.tagID = t.id
		WHERE p.username = ?
		GROUP BY p.postID
	`

	rows, err := db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []structs.Post{}
	for rows.Next() {
		var post structs.Post
		var tags string
		if err = rows.Scan(&post.ID, &post.Username, &post.Title, &post.Description, &post.CreationDate, &tags, &post.Likes, &post.Dislikes); err != nil {
			return nil, err
		}
		post.Tags = strings.Split(tags, ",")
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func getPostsLikedByUser(db *sql.DB, userID string) ([]structs.Post, error) {
	query := `
		SELECT p.postID, u.username, p.title, p.description, p.creationDate, GROUP_CONCAT(t.name), p.likes, p.dislikes
		FROM posts p
		JOIN users u ON p.username = u.username
		LEFT JOIN post_tags pt ON p.postID = pt.postID
		LEFT JOIN tags t ON pt.tagID = t.id
		JOIN post_reactions pr ON p.postID = pr.post_id
		WHERE pr.user_id = ? AND pr.reaction_type = 1
		GROUP BY p.postID
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []structs.Post{}
	for rows.Next() {
		var post structs.Post
		var tags string
		if err = rows.Scan(&post.ID, &post.Username, &post.Title, &post.Description, &post.CreationDate, &tags, &post.Likes, &post.Dislikes); err != nil {
			return nil, err
		}
		post.Tags = strings.Split(tags, ",")
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
