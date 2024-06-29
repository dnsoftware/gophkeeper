package postgresql

import (
	"context"
	"fmt"
)

func (p *PgStorage) GetEntityCodes(ctx context.Context) (map[string]string, error) {

	query := `SELECT * FROM entity_codes ORDER BY etype`
	cRows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("GetEntityCodes error: %w", err)
	}
	defer cRows.Close()

	res := make(map[string]string)
	var etype, name string
	for cRows.Next() {
		err = cRows.Scan(&etype, &name)
		if err != nil {
			return nil, fmt.Errorf("GetEntityCodes fetch error: %w", err)
		}

		res[etype] = name
	}

	return res, nil
}
