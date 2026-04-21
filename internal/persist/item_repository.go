package persist

import (
	"context"

	"GrandExchange/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ItemRepository interface {
	GetItem(context.Context, string) ([]models.SearchItem, error)
	SaveItems(context.Context, []models.SearchItem) error
}

type PostgresItemRepository struct {
	conn *pgxpool.Pool
}

func NewItemRepository(db *DatabaseContext) ItemRepository {
	return &PostgresItemRepository{conn: db.Conn}
}

func (r *PostgresItemRepository) GetItem(ctx context.Context, name string) ([]models.SearchItem, error) {
	sql := `SELECT id, name, description, type, type_icon, icon, icon_large, members, 
        current_trend, current_price, today_trend, today_price 
        FROM search_items WHERE name ILIKE $1`

	rows, err := r.conn.Query(ctx, sql, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.SearchItem
	for rows.Next() {
		var item models.SearchItem
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Type,
			&item.TypeIcon,
			&item.Icon,
			&item.IconLarge,
			&item.Members,
			&item.Current.Trend,
			&item.Current.Price,
			&item.Today.Trend,
			&item.Today.Price,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *PostgresItemRepository) SaveItems(ctx context.Context, items []models.SearchItem) error {
	sql := `INSERT INTO search_items (id, name, description, type, type_icon, icon, 
  		icon_large, members, current_trend, current_price, today_trend, today_price)
		    VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		    ON CONFLICT (id) DO UPDATE SET
			  current_trend = EXCLUDED.current_trend,
			  current_price = EXCLUDED.current_price,
			  today_trend   = EXCLUDED.today_trend,
			  today_price   = EXCLUDED.today_price,
			  updated_at    = NOW()`

	batch := &pgx.Batch{}

	for _, item := range items {
		batch.Queue(sql, item.ID, item.Name, item.Description, item.Type, item.TypeIcon,
			item.Icon, item.IconLarge, item.Members,
			item.Current.Trend, string(item.Current.Price),
			item.Today.Trend, string(item.Today.Price))
	}
	return r.conn.SendBatch(ctx, batch).Close()
}
