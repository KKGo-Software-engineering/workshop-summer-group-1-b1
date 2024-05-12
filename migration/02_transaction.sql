-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "transaction" (
  id SERIAL PRIMARY KEY,
  spender_id INT NOT NULL,
  date TIMESTAMP WITH TIME ZONE,
  amount DECIMAL(10,2) DEFAULT 0,
  category VARCHAR(50) DEFAULT '',
  transaction_type VARCHAR(20) DEFAULT '',
  note VARCHAR(255) DEFAULT '',
  image_url VARCHAR(255) DEFAULT ''
);

INSERT INTO transaction ( date , spender_id, amount , category, transaction_type, note, image_url) VALUES
('2021-07-01 00:00:00', 1, 100, 'Food', 'expense', 'Lunch', ''),
('2021-07-02 00:00:00', 1, 200, 'Transport', 'income', 'Bus', ''),
('2021-07-03 00:00:00', 1, 300, 'Food', 'income', 'Dinner', ''),
('2021-07-04 00:00:00', 1, 400, 'Transport', 'income', 'Taxi', ''),
('2021-07-05 00:00:00', 1, 500, 'Food', 'income', 'Breakfast', ''),
('2021-07-06 00:00:00', 1, 600, 'Transport', 'income', 'Bus', ''),
('2021-07-07 00:00:00', 1, 700, 'Food', 'income', 'Lunch', ''),
('2021-07-08 00:00:00', 2, 800, 'Transport', 'income', 'Bus', ''),
('2021-07-09 00:00:00', 2, 900, 'Food', 'income', 'Dinner', ''),
('2021-07-10 00:00:00', 2, 1000, 'Transport', 'income', 'Taxi', ''),
('2021-07-11 00:00:00', 2, 1100, 'Food', 'income', 'Breakfast', ''),
('2021-07-12 00:00:00', 2, 1200, 'Transport', 'expense', 'Bus', ''),
('2021-07-13 00:00:00', 1, 1300, 'Food', 'expense', 'Lunch', ''),
('2021-07-14 00:00:00', 1, 1400, 'Transport', 'expense', 'Bus', ''),
('2021-07-15 00:00:00', 1, 1500, 'Food', 'expense', 'Dinner', ''),
('2021-07-16 00:00:00', 1, 1600, 'Transport', 'expense', 'Taxi', '');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "transaction";
-- +goose StatementEnd
