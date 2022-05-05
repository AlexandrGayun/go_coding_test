CREATE TABLE IF NOT EXISTS blocks(
    number bigint PRIMARY KEY,
    transactions_count int,
    total_amount float8,
    created_at timestamp with time zone
);
