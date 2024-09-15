CREATE TABLE IF NOT EXISTS prices (
    id INT8 PRIMARY KEY DEFAULT unique_rowid(),
    update_time TIMESTAMP NOT NULL,
    instrument STRING NOT NULL,
    price FLOAT8 NOT NULL,
    quoted_instrument STRING NOT NULL
);
