package api

import (
	"context"
	"demo/utils"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type Storage interface {
	CreateUserDb(*User) error
	GetUsersDb() ([]*User, error)
	GetUserByEmailDb(string)(*User, error)
	GetUserByIdDb(string)(*User, error)
	DeleteUserByIDDb(string)error
	UpdateUserByIDDb(*UpdateUserReq, string)error
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

func (dbs *PostgreStore)GetUserByIdDb(accNum string) (*User, error) {
	query := `select * from users where id = @id;`
	id, err := strconv.Atoi(accNum)
	if err != nil {
		return nil, fmt.Errorf("wrong id type is passed", err)
	}
	args := pgx.NamedArgs{
		"id": id,
	}
	row := dbs.conn.QueryRow(context.Background(), query, args)
	fmt.Println(row)

	var user User
	errr := row.Scan(&user.ID, &user.FullName, &user.Email, &user.IsAdmin, &user.Number, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if errr != nil {
		return nil, errr
	}

	return &user, nil
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

func (dbs *PostgreStore)DeleteUserByIDDb(idd string) error {
	query := `delete from users where id = @id`
	id, err := strconv.Atoi(idd)
	if err != nil {
		return fmt.Errorf("provide proper id for deletion")
	}

	args := pgx.NamedArgs{
		"id": id,
	}
	if _, errr := dbs.conn.Exec(context.Background(), query, args); errr != nil {
		return fmt.Errorf("error while deleting in db %w", errr)
	}
	return nil
}

func (dbs *PostgreStore)UpdateUserByIDDb(upReq *UpdateUserReq, id string) error {
	idd , err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("incorrect id value (update)", err)
	}
	updates := make(map[string]any)
	var cols []string
	args := pgx.NamedArgs{}

	if upReq.FullName != nil {
		updates["fullName"] = *upReq.FullName
	}
	if upReq.Email != nil {
		updates["email"] = *upReq.Email
	}
	if upReq.Number != nil {
		updates["number"] = *upReq.Number
	}
	if upReq.Password != nil {
		hashed, err := utils.HashPassword(*upReq.Password)
		if err != nil {
			fmt.Errorf("error in hashing pass", err)
		}
		updates["password"] = hashed
	}
	
	for key, value := range updates {
		cols = append(cols, fmt.Sprintf("%s = @%s", key, key))
		args[key] = value
	}
	args["id"] = idd

	if len(cols) == 0 {
		return nil
	}

	queryStr := fmt.Sprintf("update users set %s where id = @id", strings.Join(cols, ", "))

	if _, err := dbs.conn.Exec(context.Background(), queryStr, args); err != nil {
		return fmt.Errorf("error updating id: %w", err)
	}
	return nil
}