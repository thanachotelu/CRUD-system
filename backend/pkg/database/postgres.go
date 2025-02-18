package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type User struct {
	ID          string    `json:"user_id"`
	Username    string    `json:"username"`
	Password    string    `json:"password_hash"`
	Firstname   string    `json:"firstname"`
	Lastname    string    `json:"lastname"`
	Phonenumber string    `json:"phonenumber"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Status      string    `json:"status"`
}

type NewUser struct {
	Username    string `json:"username"`
	Password    string `json:"password_hash"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Phonenumber string `json:"phonenumber"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	Status      string `json:"status"`
}

type UserDatabase interface {
	GetAllUsers(ctx context.Context) ([]User, error)
	GetUserById(ctx context.Context, id string) (User, error)
	AddUser(ctx context.Context, user NewUser) (User, error)
	UpdateUser(ctx context.Context, id string, user User) error
	DeleteUser(ctx context.Context, id string) error
	Close() error
	Ping() error
}

type PostgresDB struct {
	*sqlx.DB
	dsn string
}

func NewPostgresDB(dataSourceName string) (*PostgresDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := sqlx.ConnectContext(ctx, "postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// กำหนดค่า connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	// ทดสอบการเชื่อมต่อ
	if err = db.PingContext(ctx); err != nil {
		db.Close() // ปิดการเชื่อมต่อถ้าไม่สามารถ ping ได้
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresDB{
		DB:  db,
		dsn: dataSourceName,
	}, nil
}

func (db *PostgresDB) Reconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newDB, err := sqlx.ConnectContext(ctx, "postgres", db.dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// ตั้งค่า connection pool
	newDB.SetMaxOpenConns(25)
	newDB.SetMaxIdleConns(10)
	newDB.SetConnMaxLifetime(5 * time.Minute)

	// ทดสอบการเชื่อมต่อ
	if err = newDB.PingContext(ctx); err != nil {
		newDB.Close() // ปิดการเชื่อมต่อใหม่ถ้าไม่สามารถ ping ได้
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// ปิดการเชื่อมต่อเดิม (ถ้ามี) และกำหนดการเชื่อมต่อใหม่
	if db.DB != nil {
		db.DB.Close()
	}
	db.DB = newDB

	return nil
}

func (db *PostgresDB) Close() error {
	return db.DB.Close()
}

func hashPassword(password string) (string, error) {
	// แฮชรหัสผ่านด้วย bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (pdb *PostgresDB) GetAllUsers(ctx context.Context) ([]User, error) {
	rows, err := pdb.DB.QueryContext(ctx, `
		SELECT user_id , username , password_hash , firstname ,
		lastname , phonenumber , email , role , created_at ,
		updated_at , status
		FROM users
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var userItem User
		if err := rows.Scan(
			&userItem.ID, &userItem.Username, &userItem.Password, &userItem.Firstname,
			&userItem.Lastname, &userItem.Phonenumber, &userItem.Email, &userItem.Role,
			&userItem.CreatedAt, &userItem.UpdatedAt, &userItem.Status,
		); err != nil {
			return nil, fmt.Errorf("failed to scan users: %w", err)
		}
		users = append(users, userItem)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	return users, nil
}

func (pdb *PostgresDB) GetUserById(ctx context.Context, id string) (User, error) {
	var userItem User

	err := pdb.DB.QueryRowContext(ctx, `
		SELECT user_id , username , password_hash , firstname , 
		lastname , phonenumber , email , role , created_at , 
		updated_at , status 
		FROM users where user_id = $1
	`, id).Scan(
		&userItem.ID, &userItem.Username, &userItem.Password, &userItem.Firstname,
		&userItem.Lastname, &userItem.Phonenumber, &userItem.Email, &userItem.Role,
		&userItem.CreatedAt, &userItem.UpdatedAt, &userItem.Status)

	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("User not found")
		}
		return User{}, fmt.Errorf("failed to get User: %w", err)
	}
	return userItem, nil
}

func (pdb *PostgresDB) AddUser(ctx context.Context, user NewUser) (User, error) {
	var createdUser User
	// แฮชรหัสผ่านก่อนที่จะ insert
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	err = pdb.DB.QueryRowContext(ctx, `
		INSERT INTO users (username, password_hash, firstname, lastname, phonenumber, email, role, status) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING user_id, username, password_hash, firstname, lastname, phonenumber, email, role, status, created_at, updated_at;
	`,
		user.Username, hashedPassword, user.Firstname, user.Lastname, user.Phonenumber, user.Email, user.Role, user.Status,
	).Scan(
		&createdUser.ID, &createdUser.Username, &createdUser.Password, &createdUser.Firstname,
		&createdUser.Lastname, &createdUser.Phonenumber, &createdUser.Email,
		&createdUser.Role, &createdUser.Status, &createdUser.CreatedAt, &createdUser.UpdatedAt,
	)
	if err != nil {
		return User{}, fmt.Errorf("failed to add user: %v", err)
	}
	return createdUser, nil
}

func (pdb *PostgresDB) UpdateUser(ctx context.Context, id string, user User) error {
	// ดึงข้อมูลเดิมจากฐานข้อมูล
	var existingUser User
	err := pdb.DB.QueryRowContext(ctx, `
		SELECT username, firstname, lastname, phonenumber, email, role, status
		FROM users WHERE user_id = $1
	`, id).Scan(
		&existingUser.Username, &existingUser.Firstname, &existingUser.Lastname,
		&existingUser.Phonenumber, &existingUser.Email, &existingUser.Role, &existingUser.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("User not found")
		}
		return fmt.Errorf("failed to get existing user: %w", err)
	}

	// ใช้ค่าที่ส่งมาอัปเดต เฉพาะฟิลด์ที่ไม่เป็นค่าเริ่มต้น (default)
	if user.Username != "" {
		existingUser.Username = user.Username
	}
	if user.Firstname != "" {
		existingUser.Firstname = user.Firstname
	}
	if user.Lastname != "" {
		existingUser.Lastname = user.Lastname
	}
	if user.Phonenumber != "" {
		existingUser.Phonenumber = user.Phonenumber
	}
	if user.Email != "" {
		existingUser.Email = user.Email
	}
	if user.Role != "" {
		existingUser.Role = user.Role
	}
	if user.Status != "" {
		existingUser.Status = user.Status
	}

	// อัปเดตข้อมูล
	query := `
		UPDATE users 
		SET username = $1, firstname = $2, lastname = $3, 
			phonenumber = $4, email = $5, role = $6, status = $7, updated_at = NOW()
		WHERE user_id = $8
	`
	_, err = pdb.DB.ExecContext(ctx, query,
		existingUser.Username, existingUser.Firstname, existingUser.Lastname,
		existingUser.Phonenumber, existingUser.Email, existingUser.Role, existingUser.Status, id,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (pdb *PostgresDB) DeleteUser(ctx context.Context, id string) error {
	query := `DELETE FROM users where user_id = $1`
	result, err := pdb.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check row affected %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("User not found")
	}
	return nil
}

type UserService struct {
	db UserDatabase
}

func NewUserService(db UserDatabase) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]User, error) {
	return s.db.GetAllUsers(ctx)
}

func (s *UserService) GetUserById(ctx context.Context, id string) (User, error) {
	return s.db.GetUserById(ctx, id)
}

func (s *UserService) AddUser(ctx context.Context, user NewUser) (User, error) {
	return s.db.AddUser(ctx, user)
}

func (s *UserService) UpdateUser(ctx context.Context, id string, user User) error {
	return s.db.UpdateUser(ctx, id, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.db.DeleteUser(ctx, id)
}
