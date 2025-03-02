package api

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Storage interface {
	CreateUserDb(*User) error
	GetUsersDb() ([]*User, error)
	GetUserByEmailDb(string)(*User, error)
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
		create table if not exists users(
			id serial primary key,
			fullName varchar(200) not null,
			email varchar(200) not null unique,
			isAdmin bool default false,
			number int,
			password text not null,
			createdAt timestamp default now(),
			updatedAt timestamp default now()
		);
	`
	_, err := dbs.conn.Exec(context.Background(), query); 
	
	if err != nil {
		return fmt.Errorf("error while creating user table %w", err)
	}
	return nil
}

func (dbs *PostgreStore)CreateUserDb(user *User) error {
	query := `
		insert into users (fullName, email, isAdmin, number, password, createdAt, updatedAt) values (@fullName, @email, @isAdmin, @number, @password, @createdAt, @updatedAt);
	`
	args := pgx.NamedArgs{
		"fullName": user.FullName,
		"email": user.Email,
		"isAdmin": user.IsAdmin,
		"number": user.Number,
		"password": user.Password,
		"createdAt": time.Now().UTC(),
		"updatedAt": time.Now().UTC(),
	}

	_, err := dbs.conn.Exec(context.Background(), query, args); 
	
	if err != nil {
		return fmt.Errorf("error while creating user %v", err)
	}
	return nil
}

func (dbs *PostgreStore)GetUsersDb() ([]*User, error) {
	query := `select * from users;`

	rows, err := dbs.conn.Query(context.Background(), query)
	
	if err != nil {
		return nil, fmt.Errorf("error while fetching users: %v", err)
	}

	var usersList []*User 
	for rows.Next() {
		var user User 
		err := rows.Scan(&user.ID, &user.FullName, &user.Email, &user.IsAdmin, &user.Number, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error while scanning row: %v", err)
		}
		usersList = append(usersList, &user)
	}
	fmt.Println(usersList)
	return usersList, nil
}

func (dbs *PostgreStore)GetUserByEmailDb(email string) (*User, error) {
	query := `select * from users where email = @email;`
	args := pgx.NamedArgs{
		"email": email,
	}
	row := dbs.conn.QueryRow(context.Background(), query, args)
	fmt.Println(row)

	var user User
	err := row.Scan(&user.ID, &user.FullName, &user.Email, &user.IsAdmin, &user.Number, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}