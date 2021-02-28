package models

// const (
// 	createUserSQL     = `insert into users (username, first_name, last_name, email, password_hash, created_at) values ($1, $2, $3, $4, $5, $6) returning id`
// 	getUsersSQL       = `select id, username, first_name, last_name, email, password_hash, created_at, updated_at from users`
// 	getUserByID       = getUsersSQL + ` where id=$1`
// 	getUserByUsername = getUsersSQL + ` where username=$1`
// 	updateUserSQL     = `update users set (username, first_name, last_name, email, updated_at) = ($1, $2, $3, $4, $5) where id = $6`
// )

// type (
// 	Post struct {
// 		SequentialIdentifier
// 		UserID     string `json:"username"`
// 		CategoryID    string `json:"first_name"`
// 		Title     string `json:"last_name"`
// 		Content        string `json:"email"`
// 		Tags        string `json:"email"`
// 		PasswordHash string `json:"-"`
// 		Timestamps
// 	}
// )

// func (u *User) String() string {
// 	return fmt.Sprintf("User(%s)", u.Username)
// }
