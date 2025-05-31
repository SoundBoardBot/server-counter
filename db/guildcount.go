package db

import (
	"context"
)

func GetGuildCount(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM guilds WHERE bot_in_server = true;`

	var count int
	err := Pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
