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
		postgres.ParameterValueStr(post.Type),
		postgres.ParameterValueStr(post.Data),
		postgres.ParameterValueStr(post.Visibility),
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
		fmt.Printf("Failed to populate created post's identifier, invalid kind returned from database: %v", id_val.Kind())
	}

	return nil
}

func DbReadById(id int) (Post, error) {
	db_url := getDbUrl()
	statement := "SELECT id, author_id, content, type, data, visibility FROM posts WHERE id=$1"
	params := []postgres.ParameterValue{
		postgres.ParameterValueInt32(int32(id)),
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

func DbReadAll(limit int, offset int) ([]Post, error) {
	db_url := getDbUrl()
	statement := "SELECT id, author_id, content, type, data, visibility FROM posts ORDER BY id LIMIT $1 OFFSET $2"
	params := []postgres.ParameterValue{
		postgres.ParameterValueInt64(int64(limit)),
		postgres.ParameterValueInt64(int64(offset)),
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
	statement := "UPDATE posts SET author_id=$1, content=$2, type=$3, data=$4, visibility=$5 WHERE id=$6"
	params := []postgres.ParameterValue{
		postgres.ParameterValueStr(post.AuthorID),
		postgres.ParameterValueStr(post.Content),
		postgres.ParameterValueStr(post.Type),
		postgres.ParameterValueStr(post.Data),
		postgres.ParameterValueStr(post.Visibility),
		postgres.ParameterValueInt32(int32(post.ID)),
	}

	_, err := postgres.Execute(db_url, statement, params)
	if err != nil {
		return fmt.Errorf("Error updating database: %s", err.Error())
	}

	return nil
}

func DbDelete(id int) error {
	db_url := getDbUrl()
	statement := "DELETE FROM posts WHERE id=$1"
	params := []postgres.ParameterValue{
		postgres.ParameterValueInt32(int32(id)),
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

	if val, err := assertValueKind(row, 0, postgres.DbValueKindInt32); err == nil {
		post.ID = int(val.GetInt32())
	} else {
		return post, err
	}

	if val, err := assertValueKind(row, 1, postgres.DbValueKindStr); err == nil {
		post.AuthorID = val.GetStr()
	} else {
		return post, err
	}

	if val, err := assertValueKind(row, 2, postgres.DbValueKindStr); err == nil {
		post.Content = val.GetStr()
	} else {
		return post, err
	}

	if val, err := assertValueKind(row, 3, postgres.DbValueKindStr); err == nil {
		post.Type = val.GetStr()
	} else {
		return post, err
	}

	if val, err := assertValueKind(row, 4, postgres.DbValueKindStr); err == nil {
		post.Data = val.GetStr()
	} else {
		return post, err
	}

	if val, err := assertValueKind(row, 5, postgres.DbValueKindStr); err == nil {
		post.Visibility = val.GetStr()
	} else {
		return post, err
	}

	return post, nil
}
