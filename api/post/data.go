package main

import (
	"fmt"

	"github.com/fermyon/spin/sdk/go/postgres"
)

// Database Operations

func DbInsert(post *Post) error {
	db_url := getDbUrl()
	statement := "INSERT INTO posts (author_id, content, type, data, visibility) VALUES ($1, $2, $3, $4, $5)"
	params := []postgres.ParameterValue{
		postgres.ParameterValueStr(post.AuthorID),
		postgres.ParameterValueStr(post.Content),
		postgres.ParameterValueStr(post.Type.String()),
		postgres.ParameterValueStr(post.Data),
		postgres.ParameterValueStr(post.Visibility.String()),
	}

	_, err := postgres.Execute(db_url, statement, params)
	if err != nil {
		return fmt.Errorf("Error inserting into database: %s", err.Error())
	}

	// this is a gross hack that will surely bite me later
	rowset, err := postgres.Query(db_url, "SELECT lastval()", []postgres.ParameterValue{})
	if err != nil || len(rowset.Rows) != 1 || len(rowset.Rows[0]) != 1 {
		return fmt.Errorf("Error querying id from database: %s", err.Error())
	}

	id_val := rowset.Rows[0][0]
	if id_val.Kind() == postgres.DbValueKindInt64 {
		post.ID = int(id_val.GetInt64())
	} else {
		fmt.Printf("Failed to populate created post's identifier, invalid kind returned from database: %v\n", id_val.Kind())
	}

	return nil
}

func DbReadById(id int, authorId string) (Post, error) {
	db_url := getDbUrl()
	statement := "SELECT id, author_id, content, type, data, visibility FROM posts WHERE id=$1 and author_id=$2"
	params := []postgres.ParameterValue{
		postgres.ParameterValueInt32(int32(id)),
		postgres.ParameterValueStr(authorId),
	}

	rowset, err := postgres.Query(db_url, statement, params)
	if err != nil {
		return Post{}, fmt.Errorf("Error reading from database: %s", err.Error())
	}

	if rowset.Rows == nil || len(rowset.Rows) == 0 {
		return Post{}, nil
	} else {
		return fromRow(rowset.Rows[0])
	}
}

func DbReadAll(limit int, offset int, authorId string) ([]Post, error) {
	db_url := getDbUrl()
	statement := "SELECT id, author_id, content, type, data, visibility FROM posts WHERE author_id=$3 ORDER BY id LIMIT $1 OFFSET $2"
	params := []postgres.ParameterValue{
		postgres.ParameterValueInt64(int64(limit)),
		postgres.ParameterValueInt64(int64(offset)),
		postgres.ParameterValueStr(authorId),
	}
	rowset, err := postgres.Query(db_url, statement, params)
	if err != nil {
		return []Post{}, fmt.Errorf("Error reading from database: %s", err.Error())
	}

	posts := make([]Post, len(rowset.Rows))
	for i, row := range rowset.Rows {
		if post, err := fromRow(row); err != nil {
			return []Post{}, err
		} else {
			posts[i] = post
		}
	}

	return posts, nil
}

func DbUpdate(post Post) error {
	db_url := getDbUrl()
	statement := "UPDATE posts SET content=$1, type=$2, data=$3, visibility=$4 WHERE id=$5"
	params := []postgres.ParameterValue{
		postgres.ParameterValueStr(post.Content),
		postgres.ParameterValueStr(post.Type.String()),
		postgres.ParameterValueStr(post.Data),
		postgres.ParameterValueStr(post.Visibility.String()),
		postgres.ParameterValueInt32(int32(post.ID)),
	}

	_, err := postgres.Execute(db_url, statement, params)
	if err != nil {
		return fmt.Errorf("Error updating database: %s", err.Error())
	}

	return nil
}

func DbDelete(id int, authorId string) error {
	db_url := getDbUrl()
	statement := "DELETE FROM posts WHERE id=$1 and author_id=$2"
	params := []postgres.ParameterValue{
		postgres.ParameterValueInt32(int32(id)),
		postgres.ParameterValueStr(authorId),
	}

	_, err := postgres.Execute(db_url, statement, params)
	return err
}

// Database Helper Functions

func assertValueKind(row []postgres.DbValue, col int, expected postgres.DbValueKind) (postgres.DbValue, error) {
	if row[col].Kind() != expected {
		return postgres.DbValue{}, fmt.Errorf("Expected column %v to be %v kind but got %v\n", col, expected, row[col].Kind())
	}
	return row[col], nil
}

func fromRow(row []postgres.DbValue) (Post, error) {
	var post Post

	if val, err := assertValueKind(row, 0, postgres.DbValueKindInt32); err != nil {
		return post, err
	} else {
		post.ID = int(val.GetInt32())
	}

	if val, err := assertValueKind(row, 1, postgres.DbValueKindStr); err != nil {
		return post, err
	} else {
		post.AuthorID = val.GetStr()
	}

	if val, err := assertValueKind(row, 2, postgres.DbValueKindStr); err != nil {
		return post, err
	} else {
		post.Content = val.GetStr()
	}

	if val, err := assertValueKind(row, 3, postgres.DbValueKindStr); err != nil {
		return post, err
	} else if val, err := ParsePostType(val.GetStr()); err != nil {
		return post, err
	} else {
		post.Type = val
	}

	if val, err := assertValueKind(row, 4, postgres.DbValueKindStr); err != nil {
		return post, err
	} else {
		post.Data = val.GetStr()
	}

	if val, err := assertValueKind(row, 5, postgres.DbValueKindStr); err != nil {
		return post, err
	} else if val, err := ParsePostVisibility(val.GetStr()); err != nil {
		return post, err
	} else {
		post.Visibility = val
	}

	return post, nil
}
