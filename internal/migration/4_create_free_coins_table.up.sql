CREATE TABLE free_coins (
	id INT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
	description VARCHAR(255) NOT NULL,
	amount INT NOT NULL,
	created_at TIMESTAMP NOT NULL
);