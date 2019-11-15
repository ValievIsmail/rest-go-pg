package main

import (
	"context"
	"database/sql"
	"time"
)

func dbGetCommentByID(id int, db *sql.DB) (data []byte, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `select api.get_comment_by_id($1)`

	if err := db.QueryRowContext(ctx, query, id).Scan(&data); err != nil {
		return nil, err
	}

	return data, nil
}

func dbGetAllComments(db *sql.DB) (data []byte, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `select api.get_comments()`

	if err := db.QueryRowContext(ctx, query).Scan(&data); err != nil {
		return nil, err
	}

	return data, nil
}

func dbUpdateComment(cid int, msg string, db *sql.DB) (id int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `select api.update_comment($1, $2)`

	if err := db.QueryRowContext(ctx, query, cid, msg).Scan(&id); err != nil {
		return id, err
	}

	return id, nil
}

func dbCreateComment(userID int, msg string, db *sql.DB) (id int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `select api.create_comment($1, $2)`

	if err := db.QueryRowContext(ctx, query, userID, msg).Scan(&id); err != nil {
		return id, err
	}

	return id, nil
}

func dbDeleteComment(cid int, db *sql.DB) (id int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `select api.delete_comment($1)`

	if err := db.QueryRowContext(ctx, query, cid).Scan(&id); err != nil {
		return id, err
	}

	return id, nil
}

func dbGetAllUsers(db *sql.DB) (data []byte, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `select api.get_users()`

	if err := db.QueryRowContext(ctx, query).Scan(&data); err != nil {
		return nil, err
	}

	return data, nil
}

func dbCreateUser(name string, db *sql.DB) (id int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `select api.create_user($1)`

	if err := db.QueryRowContext(ctx, query, name).Scan(&id); err != nil {
		return id, err
	}

	return id, nil
}

func dbUpdateUser(uid int, name string, db *sql.DB) (id int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `select api.update_user($1, $2)`

	if err := db.QueryRowContext(ctx, query, uid, name).Scan(&id); err != nil {
		return id, err
	}

	return id, nil
}

func dbDeleteUser(uid int, db *sql.DB) (id int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `select api.delete_user($1)`

	if err := db.QueryRowContext(ctx, query, uid).Scan(&id); err != nil {
		return id, err
	}

	return id, nil
}
