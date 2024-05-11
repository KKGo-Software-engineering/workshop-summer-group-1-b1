-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "transaction" (
  id SERIAL PRIMARY KEY,
  date TIMESTAMP WITH TIME ZONE,
  amount DECIMAL(10,2) DEFAULT 0,
  category VARCHAR(50) DEFAULT '',
  transaction_type VARCHAR(20) DEFAULT '',
  note VARCHAR(255) DEFAULT '',
  image_url VARCHAR(255) DEFAULT ''
);

INSERT INTO transaction ( date , amount , category, transaction_type, note, image_url) VALUES
( '2020-01-01', 100, 'Food', 'expense', 'Lunch', ''),
( '2020-01-02', 200, 'Transport', 'expense', 'Bus', ''),
( '2020-01-03', 300, 'Food', 'expense', 'Dinner', ''),
( '2020-01-04', 400, 'Transport', 'expense', 'Train', ''),
( '2020-01-05', 500, 'Food', 'expense', 'Breakfast', ''),
( '2020-01-06', 600, 'Transport', 'expense', 'Bus', ''),
( '2020-01-07', 700, 'Food', 'expense', 'Lunch', ''),
('2020-01-08', 800, 'Transport', 'expense', 'Train', ''),
('2020-01-09', 900, 'Food', 'expense', 'Dinner', ''),
('2020-01-10', 1000, 'Transport', 'expense', 'Bus', '');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "transaction";
-- +goose StatementEnd
