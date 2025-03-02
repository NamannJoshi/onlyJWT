package api

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Storage interface {
	CreateUserDb(*User) error
}

type PostgreStore struct {
	conn *pgx.Conn
}

func NewPostgreStore() (*PostgreStore, error) {
	connStr := "postgres://postgres:mukeshakun@localhost:5432/test"

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	return &PostgreStore{
		conn: conn,
	}, nil
}

func (dbs *PostgreStore)Init() error {
	return dbs.CreateUserTable()
}

func (dbs *PostgreStore)CreateUserTable() error {
	query := `
		create table if not exists users (
			id serial primary key,
			fullName varchar(200) not null,
			email varchar(200) not null unique,
			isAdmin bool default false,
			number int,
			password string not null
			createdAt timestamp not null,
			updateAt timestamp not null
		);
	`
	if _, err := dbs.conn.Exec(context.Background(), query); err != nil {
		return fmt.Errorf("error while creating user table %v", err)
	}
	return nil
}

func (dbs *PostgreStore)CreateUserDb(user *User) error {
	// query := `
	// 	insert into users (fullName, email, isAdmin, number, createdAt, updatedAt) values (@fullName, @email, @isAdmin, @number, @createdAt, @updatedAt);
	// `

	// if _, err := dbs.conn.Exec(context.Background(), query); err != nil {
	// 	return fmt.Errorf("error while creating user %v", err)
	// }
	return fmt.Errorf("sab sahi hora hai")
}