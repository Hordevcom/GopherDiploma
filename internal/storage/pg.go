package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/Hordevcom/GopherDiploma/internal/config"
	"github.com/Hordevcom/GopherDiploma/internal/middleware/logging"
	"github.com/Hordevcom/GopherDiploma/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PGDB struct {
	logger logging.Logger
	DB     *pgxpool.Pool
}

func (p *PGDB) SetUserWithdrawn(ctx context.Context, orderNum string, newBalance, withdrawn float32) {
}

func (p *PGDB) UpdateUserBalance(ctx context.Context, user string, accrual, withdrawn float32) error {
	query := `UPDATE users SET accrual = $1, withdrawn = $2
			WHERE username = $3`

	result, err := p.DB.Exec(ctx, query, int(accrual*100), int(withdrawn*100), user)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("0 rows affected")
	}

	return nil
}

func (p *PGDB) GetUserBalance(ctx context.Context, user string) (models.UserBalance, error) {
	var balance models.UserBalance
	query := `SELECT accrual, withdrawn 
		FROM users WHERE username = $1`
	row := p.DB.QueryRow(ctx, query, user)
	err := row.Scan(&balance.Current, &balance.Withdrawn)

	if err != nil {
		return models.UserBalance{}, err
	}

	balance.Current = balance.Current * 100

	return balance, nil
}

func (p *PGDB) UpdateStatus(ctx context.Context, newStatus, order, user string) error {
	query := `UPDATE orders SET status = $1
			WHERE number = $2`

	result, err := p.DB.Exec(ctx, query, newStatus, order)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("0 rows affected")
	}

	// query = `UPDATE users SET accrual = $1
	// 		WHERE username = $2`

	// result, err = p.DB.Exec(ctx, query, int(accrual*100), user)

	// if err != nil {
	// 	return err
	// }

	// if result.RowsAffected() == 0 {
	// 	return fmt.Errorf("0 rows affected")
	// }

	return nil
}

func (p *PGDB) GetUserOrders(ctx context.Context, user string) ([]models.Order, error) {
	query := `SELECT number, status, uploaded_at 
		FROM orders WHERE username = $1`
	rows, err := p.DB.Query(ctx, query, user)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		// var accrual sql.NullInt64

		err := rows.Scan(&o.Number, &o.Status, &o.UploadAt)
		if err != nil {
			return nil, err
		}

		// if accrual.Valid {
		// 	val := int(accrual.Int64)
		// 	o.Accrual = val
		// } else {
		// 	o.Accrual = 0
		// }
		orders = append(orders, o)
	}

	return orders, nil
}

func (p *PGDB) GetOrderAndUser(ctx context.Context, order string) (string, string, error) {
	var userOrder string
	var username string

	query := `SELECT number, username FROM orders WHERE number = $1`
	row := p.DB.QueryRow(ctx, query, order)
	err := row.Scan(&userOrder, &username)

	return userOrder, username, err
}

func (p *PGDB) AddOrderToDB(ctx context.Context, order string, username string) error {
	query := `INSERT INTO orders (number, uploaded_at, username)
				VALUES ($1, $2, $3) ON CONFLICT (number) DO NOTHING`

	_, err := p.DB.Exec(ctx, query, order, time.Now(), username)

	return err
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
	var insertedUser string
	query := `INSERT INTO users (username, user_password)
				VALUES ($1, $2) ON CONFLICT (username) DO NOTHING
				RETURNING username`

	err := p.DB.QueryRow(ctx, query, username, password).Scan(&insertedUser)

	if err != nil {
		return err
	}

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
