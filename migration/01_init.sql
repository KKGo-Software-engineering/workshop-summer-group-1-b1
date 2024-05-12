-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "spender" (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL
);

INSERT INTO spender ( name , email) VALUES
('John Doe', 'john001@gmail.com'),
('Jane Doe', 'jane002@gmail.com');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "spender";
-- +goose StatementEnd
