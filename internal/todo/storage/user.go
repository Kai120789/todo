package storage

import (
	"context"
	"fmt"
	"todo/internal/todo/dto"
	"todo/internal/todo/models"
	"todo/internal/todo/utils/hash"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type UserStorage struct {
	db *pgxpool.Pool
}

type UserStorager interface {
	RegisterNewUser(body dto.PostUserDto) (*models.UserToken, error)
	AuthorizateUser(body dto.PostUserDto) (*models.UserToken, *uint, error)
	WriteRefreshToken(userId uint, refreshTokenValue string) error
	GetAuthUser(id uint) (*models.UserToken, error)
	UserLogout(id uint) error

	GetAllUsers() ([]models.TgUser, error)
	GetChatID(task *models.Task) (int64, error)
	AddChatID(tgName string, chatID int64) error
	GetMyTasks(tgName string, status int) ([]models.Task, int64, error)
	ChangeEndedTasksStatus() error
}

func NewUserStore(Conn *pgxpool.Pool, log *zap.Logger) *UserStorage {
	return &UserStorage{db: Conn}
}

// register new user
func (d *UserStorage) RegisterNewUser(body dto.PostUserDto) (*models.UserToken, error) {
	var id uint
	query := `INSERT INTO users (username, tg_name, password_hash) VALUES ($1, $2, $3) RETURNING id`
	err := d.db.QueryRow(context.Background(), query, body.Username, body.TgName, body.PasswordHash).Scan(&id)
	if err != nil {
		return nil, err
	}

	userRet, err := d.GetAuthUser(uint(id))
	if err != nil {
		return nil, err
	}

	return userRet, nil
}

// login user
func (d *UserStorage) AuthorizateUser(body dto.PostUserDto) (*models.UserToken, *uint, error) {
	var id uint
	var passwordHash string

	query := `SELECT id, password_hash FROM users WHERE username=$1`
	err := d.db.QueryRow(context.Background(), query, body.Username).Scan(&id, &passwordHash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil, fmt.Errorf("user not found")
		}
		return nil, nil, err
	}

	if !hash.CheckPasswordHash(body.PasswordHash, passwordHash) {
		return nil, nil, err
	}

	userRet, err := d.GetAuthUser(id)
	if err != nil {
		return nil, nil, err
	}

	return userRet, &id, nil
}

// get auth user
func (d *UserStorage) GetAuthUser(id uint) (*models.UserToken, error) {
	query := `SELECT * FROM user_token WHERE user_id=$1`
	row := d.db.QueryRow(context.Background(), query, id)

	var token models.UserToken
	err := row.Scan(&token.ID, &token.UserID, &token.RefreshToken)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &token, nil
}

// get tg_name by id
func (d *UserStorage) GetTgName(id uint) (string, error) {
	var tgName string

	query := `SELECT tg_name FROM users WHERE id = $1`
	err := d.db.QueryRow(context.Background(), query, id).Scan(&tgName)
	if err != nil {
		return "", fmt.Errorf("failed to get tg_name for user with id %d: %w", id, err)
	}

	return tgName, nil
}

// logout user
func (d *UserStorage) UserLogout(id uint) error {
	query := `DELETE FROM user_token WHERE user_id=$1`
	_, err := d.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	return nil
}

// add refresh token to db
func (d *UserStorage) WriteRefreshToken(userId uint, refreshTokenValue string) error {
	query := `INSERT INTO user_token (user_id, refresh_token) VALUES ($1, $2)`
	_, err := d.db.Exec(context.Background(), query, userId, refreshTokenValue)
	if err != nil {
		return err
	}

	return nil
}

func (d *UserStorage) Close() error {
	if d.db == nil {
		return nil
	}
	d.db.Close()
	return nil
}

func (d *UserStorage) AddChatID(tgName string, chatID int64) error {
	query := `UPDATE users SET chat_id = $1 WHERE tg_name = $2`
	_, err := d.db.Exec(context.Background(), query, chatID, tgName)
	if err != nil {
		fmt.Printf("Error registering user: %v\n", err)
		return err
	}

	return nil
}

func (d *UserStorage) GetMyTasks(tgName string, status int) ([]models.Task, *int64, error) {
	var id, chatID uint

	query := `SELECT id, chat_id FROM users WHERE tg_name=$1`
	err := d.db.QueryRow(context.Background(), query, tgName).Scan(&id, &chatID)
	fmt.Println("chatID from DB:", chatID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil, fmt.Errorf("user not found")
		}
		return nil, nil, err
	}

	query = `SELECT * FROM tasks WHERE user_id=$1 and status_id=$2 ORDER BY updated_at`
	rows, err := d.db.Query(context.Background(), query, id, status)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.BoardId, &task.StatusId, &task.UserId, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, nil, err
		}
		tasks = append(tasks, task)
	}

	intChatID := int64(chatID)

	return tasks, &intChatID, nil
}

func (d *UserStorage) GetAllUsers() ([]models.TgUser, error) {
	query := `SELECT id, tg_name, chat_id FROM users`
	rows, err := d.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.TgUser
	for rows.Next() {
		var user models.TgUser
		if err := rows.Scan(&user.ID, &user.TgName, &user.ChatID); err != nil {
			fmt.Println("Error during scanning row:", err)
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (d *UserStorage) ChangeEndedTasksStatus() error {
	query := `UPDATE tasks SET status_id = 3 WHERE status_id = 2`
	_, err := d.db.Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении статуса задач: %w", err)
	}

	return nil
}
