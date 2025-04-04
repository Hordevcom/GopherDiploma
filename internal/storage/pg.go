package storage

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/Hordevcom/GopherDiploma/internal/config"
	"github.com/Hordevcom/GopherDiploma/internal/middleware/logging"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
	"go.uber.org/zap"
)

type PGDB struct {
	logger logging.Logger
	DB     *pgxpool.Pool
}

func (p *PGDB) GetUserPassword(ctx context.Context, username string) string {
	var password string

	query := `SELECT user_password FROM users WHERE username = $1`
	row := p.DB.QueryRow(ctx, query, username)
	row.Scan(&password)

	return password
}

func (p *PGDB) CheckUsernameLogin(ctx context.Context, username string) bool {
	var user string

	query := `SELECT username FROM users WHERE username = $1`
	row := p.DB.QueryRow(ctx, query, username)
	row.Scan(&user)

	return user != ""
}

func (p *PGDB) AddUserToDB(ctx context.Context, username, password string) error {
	fmt.Printf("user: %v, password: %v", username, password)
	var insertedUser string
	query := `INSERT INTO users (username, user_password)
				VALUES ($1, $2) ON CONFLICT (username) DO NOTHING
				RETURNING username`

	//result, err := p.DB.Exec(ctx, query, username, password)

	err := p.DB.QueryRow(ctx, query, username, password).Scan(&insertedUser)

	if err != nil {
		return err
	}

	fmt.Printf("New user in base: %v", insertedUser)

	// if rows := result.RowsAffected(); rows == 0 {
	// 	fmt.Println("!!!!!!No data added!!!!!!")
	// 	return errors.New("no update")
	// }

	return nil
}

func NewPGDB(conf config.Config, logger logging.Logger) *PGDB {
	db, err := pgxpool.New(context.Background(), conf.DatabaseDsn)

	if err != nil {
		logger.Logger.Errorw("Problem with connecting to db: ", err)
		return nil
	}

	err = db.Ping(context.Background())

	if err != nil {
		logger.Logger.Errorw("Problem with ping to db: ", err)
		return nil
	}

	return &PGDB{logger: logger, DB: db}
}

func InitMigrations(logger zap.SugaredLogger, conf config.Config) {
	logger.Infow("Start migrations")
	db, err := sql.Open("pgx", conf.DatabaseDsn)

	if err != nil {
		logger.Fatalw("Error with connection to DB: ", err)
	}

	defer db.Close()

	_, filename, _, _ := runtime.Caller(0)
	migrationsPath := filepath.Join(filepath.Dir(filename), "migrations")

	err = goose.Up(db, migrationsPath)
	if err != nil {
		logger.Fatalw("Error with migrations: ", err)
	}
}
