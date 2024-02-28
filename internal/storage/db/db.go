package db

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/ShvetsovYura/metrics-collector/internal/storage/metric"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBStore struct {
	pool *pgxpool.Pool
}

func NewDBPool(ctx context.Context, connString string) (*DBStore, error) {
	connPool, err := pgxpool.New(ctx, connString)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}
	createErr := createTables(ctx, connPool)
	if createErr != nil {
		return nil, createErr
	}

	return &DBStore{pool: connPool}, nil
}

func createTables(ctx context.Context, connectionPool *pgxpool.Pool) error {
	_, err := connectionPool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS counter
		(
			id  bigserial not null,
			name TEXT NOT NULL,
			value bigint NOT NULL,
			updated_at timestamp with time zone NOT NULL DEFAULT now(),
			CONSTRAINT counter_pkey PRIMARY KEY (id),
			CONSTRAINT counter_metric_name UNIQUE (name)
		);
		CREATE TABLE IF NOT EXISTS gauge
		(
			id bigserial NOT NULL,
			name TEXT NOT NULL,
			value double precision NOT NULL,
			updated_at timestamp with time zone NOT NULL DEFAULT now(),
			CONSTRAINT gauge_pkey PRIMARY KEY (id),
			CONSTRAINT gauge_metric_name UNIQUE (name)
		);
	`)
	return err
}

func (db *DBStore) SetGauge(ctx context.Context, name string, value float64) error {
	tag, err := db.pool.Exec(ctx,
		`
		insert into gauge (name, value) values($1, $2)
		on conflict (name) do update set value = $2
		`, name, value)

	if err != nil {
		return err
	}
	logger.Log.Info(tag)
	return nil
}

func (db *DBStore) SetCounter(ctx context.Context, name string, value int64) error {
	stmt, args, _ := sq.Insert("counter").
		Columns("name", "value").
		Values(name, value).
		Suffix("on conflict (name) do update set value=EXCLUDED.value + counter.value").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	_, err := db.pool.Exec(ctx, stmt, args...)
	return err
}

func (db *DBStore) GetCounter(ctx context.Context, name_ string) (metric.Counter, error) {
	stmt, args, _ := sq.Select("name", "value").From("counter").Where(sq.Eq{"name": name_}).PlaceholderFormat(sq.Dollar).ToSql()
	row := db.pool.QueryRow(ctx, stmt, args...)
	var name string
	var value float64
	err := row.Scan(&name, &value)
	if err != nil {
		return metric.Counter(0), err
	}
	return metric.Counter(value), nil
}

func (db *DBStore) GetGauge(ctx context.Context, name_ string) (metric.Gauge, error) {
	stmt, args, _ := sq.Select("name", "value").From("gauge").Where(sq.Eq{"name": name_}).PlaceholderFormat(sq.Dollar).ToSql()
	row := db.pool.QueryRow(ctx, stmt, args...)
	var name string
	var value float64
	err := row.Scan(&name, &value)
	if err != nil {
		return metric.Gauge(0), err
	}
	return metric.Gauge(value), nil
}

func (db *DBStore) GetGauges(ctx context.Context) (map[string]metric.Gauge, error) {
	stmt, _, err := sq.Select("name", "value").From("gauge").ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := db.pool.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	var gauges = make(map[string]metric.Gauge, 100)

	for rows.Next() {
		var name string
		var value float64

		err := rows.Scan(&name, &value)
		if err != nil {
			return nil, err
		}

		gauges[name] = metric.Gauge(value)
	}
	return gauges, nil
}

func (db *DBStore) GetCounters(ctx context.Context) (map[string]metric.Counter, error) {
	stmt, _, err := sq.Select("name", "value").From("couter").ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := db.pool.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	var counters = make(map[string]metric.Counter, 1)

	for rows.Next() {
		var name string
		var value int64

		err := rows.Scan(&name, &value)
		if err != nil {
			return nil, err
		}

		counters[name] = metric.Counter(value)
	}

	return counters, nil
}

func (db *DBStore) ToList(ctx context.Context) ([]string, error) {
	var list []string
	g, err := db.GetGauges(ctx)
	if err != nil {
		return nil, err
	}
	for _, v := range g {
		list = append(list, v.ToString())
	}
	c, err := db.GetCounters(ctx)
	if err != nil {
		return nil, err
	}
	for _, v := range c {
		list = append(list, v.ToString())
	}
	return list, nil
}

func (db *DBStore) Ping(ctx context.Context) error {
	err := db.pool.Ping(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (db *DBStore) SaveGaugesBatch(ctx context.Context, gauges map[string]metric.Gauge) error {
	logger.Log.Info("save metrics in DBStorage GAUGES")
	stmt := "insert into gauge(name, value) values(@name, @value) on conflict (name) do update set value=@value"
	batch := &pgx.Batch{}
	for k, v := range gauges {
		args := pgx.NamedArgs{
			"name":  k,
			"value": v.GetRawValue(),
		}
		batch.Queue(stmt, args)
	}
	results := db.pool.SendBatch(ctx, batch)
	defer results.Close()
	_, err := results.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (db *DBStore) SaveCountersBatch(ctx context.Context, counters map[string]metric.Counter) error {
	logger.Log.Info("save metrics in DBStorage COUNTERS")

	insertStmt := sq.Insert("counter").Columns("name", "value").
		Suffix("on conflict (name) do update set value = counter.value + EXCLUDED.value").
		PlaceholderFormat(sq.Dollar)

	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	for k, v := range counters {
		stmt, args, err := insertStmt.Values(k, *v.GetRawValue()).ToSql()
		if err != nil {
			return err
		}
		_, err = db.pool.Exec(ctx, stmt, args...)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (db *DBStore) Save() error {
	return nil
}

func (db *DBStore) Restore(ctx context.Context) error {
	return nil
}
