package db

import (
	"context"

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
	createErr := createTables(connPool)
	if err != nil {
		return nil, createErr
	}

	return &DBStore{pool: connPool}, nil
}

func createTables(connectionPool *pgxpool.Pool) error {
	_, err := connectionPool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS public.counter
		(
			id  serial not null,
			name character varying(255) NOT NULL,
			value integer NOT NULL,
			updated_at timestamp with time zone NOT NULL DEFAULT now(),
			CONSTRAINT counter_pkey PRIMARY KEY (id)
		);
		CREATE TABLE IF NOT EXISTS public.gauge
		(
			id serial NOT NULL,
			name character varying(255) NOT NULL,
			value double precision NOT NULL,
			updated_at timestamp with time zone NOT NULL DEFAULT now(),
			CONSTRAINT gauge_pkey PRIMARY KEY (id)
		);

		DO
		$do$
		BEGIN
		IF NOT EXISTS (
			SELECT * FROM pg_constraint WHERE conname='gauge_metric_name' AND contype = 'u'
		) THEN
			ALTER TABLE IF EXISTS gauge ADD CONSTRAINT gauge_metric_name UNIQUE (name);
		END IF;
		END
		$do$;

		DO
		$do$
		BEGIN
		IF NOT EXISTS (
			SELECT * FROM pg_constraint WHERE conname='counter_metric_name' AND contype = 'u'
		) THEN
			ALTER TABLE IF EXISTS counter ADD CONSTRAINT counter_metric_name UNIQUE (name);
		END IF;
		END
		$do$;
	`)
	return err
}

func (db *DBStore) SetGauge(name string, value float64) error {
	tag, err := db.pool.Exec(context.Background(),
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

func (db *DBStore) SetCounter(name string, value int64) error {
	_, err := db.pool.Exec(context.Background(),
		`
	insert into counter(name, value) values($1, $2)
	on conflict (name) do update set value = $2	
	`, name, value)
	return err
}

func (db *DBStore) GetCounter(name_ string) (metric.Counter, error) {
	row := db.pool.QueryRow(context.Background(), "select name, value from counter where name=$1", name_)
	var name string
	var value float64
	err := row.Scan(&name, &value)
	if err != nil {
		return metric.Counter(0), err
	}
	return metric.Counter(value), nil
}

func (db *DBStore) GetGauge(name_ string) (metric.Gauge, error) {
	row := db.pool.QueryRow(context.Background(), "select name, value from gauge where name=$1", name_)
	var name string
	var value float64
	err := row.Scan(&name, &value)
	if err != nil {
		return metric.Gauge(0), err
	}
	return metric.Gauge(value), nil
}

func (db *DBStore) GetGauges() map[string]metric.Gauge {
	rows, _ := db.pool.Query(context.Background(), "select name, value from gauge")
	var gauges = make(map[string]metric.Gauge, 100)

	for rows.Next() {
		var name string
		var value float64

		_ = rows.Scan(&name, &value)

		gauges[name] = metric.Gauge(value)
	}
	return gauges
}

func (db *DBStore) GetCounters() map[string]metric.Counter {
	rows, _ := db.pool.Query(context.Background(), "select name, value from counter")
	var counters = make(map[string]metric.Counter, 1)

	for rows.Next() {
		var name string
		var value int64

		_ = rows.Scan(&name, &value)

		counters[name] = metric.Counter(value)
	}

	return counters
}

func (db *DBStore) ToList() []string {
	var list []string

	for _, g := range db.GetGauges() {
		list = append(list, g.ToString())
	}
	for _, c := range db.GetCounters() {
		list = append(list, c.ToString())
	}
	return list
}

func (db *DBStore) Ping() error {
	err := db.pool.Ping(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (db *DBStore) Save() error {
	return nil
}

func (db *DBStore) Restore() error {
	return nil
}

func (db *DBStore) SaveGaugesBatch(gauges map[string]metric.Gauge) {
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
	results := db.pool.SendBatch(context.Background(), batch)
	defer results.Close()
	results.Exec()
}
func (db *DBStore) SaveCountersBatch(counters map[string]metric.Counter) {
	logger.Log.Info("save metrics in DBStorage COUNTERS")
	stmt := "insert into counter(name, value) values(@name, @value) on conflict (name) do update set value=@value"
	batch := &pgx.Batch{}
	for k, v := range counters {
		args := pgx.NamedArgs{
			"name":  k,
			"value": v.GetRawValue(),
		}
		batch.Queue(stmt, args)
	}
	results := db.pool.SendBatch(context.Background(), batch)
	defer results.Close()
	results.Exec()
}
