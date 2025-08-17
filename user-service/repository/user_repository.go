package repository

import (
	"database/sql"
	"errors"
	"user-service/model"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByEmail(email string) (*model.User, error)
	GetUserByID(id string) (*model.User, error)
	UpdateUserVerification(id string, verified bool) error
	UpdateUserProfile(user *model.User) error
	ListUsers() ([]model.User, error)
	DeleteUser(id string) error
	UpdateUserRole(id string, role string) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) CreateUser(user *model.User) error {
	query := `INSERT INTO users (id, name, email, password_hash, phone, is_verified, role, created_at)
	          VALUES ($1,$2,$3,$4,$5,$6,$7,NOW())`
	_, err := r.db.Exec(query, user.ID, user.Name, user.Email, user.PasswordHash, user.Phone, user.IsVerified, user.Role)
	return err
}

func (r *userRepository) GetUserByEmail(email string) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, name, email, password_hash, phone, is_verified, role, created_at, updated_at FROM users WHERE email=$1`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Phone,
		&user.IsVerified, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUserByID(id string) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, name, email, password_hash, phone, is_verified, role, created_at, updated_at FROM users WHERE id=$1`
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Phone,
		&user.IsVerified, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) UpdateUserVerification(id string, verified bool) error {
	query := `UPDATE users SET is_verified=$1, updated_at=NOW() WHERE id=$2`
	res, err := r.db.Exec(query, verified, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("no rows updated")
	}
	return nil
}

func (r *userRepository) UpdateUserProfile(user *model.User) error {
	query := `UPDATE users SET name=$1, phone=$2, updated_at=NOW() WHERE id=$3`
	_, err := r.db.Exec(query, user.Name, user.Phone, user.ID)
	return err
}

func (r *userRepository) ListUsers() ([]model.User, error) {
	rows, err := r.db.Query(`SELECT id, name, email, phone, is_verified, role, created_at, updated_at FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.IsVerified, &u.Role, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *userRepository) DeleteUser(id string) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id=$1`, id)
	return err
}

func (r *userRepository) UpdateUserRole(id string, role string) error {
	_, err := r.db.Exec(`UPDATE users SET role=$1, updated_at=NOW() WHERE id=$2`, role, id)
	return err
}
