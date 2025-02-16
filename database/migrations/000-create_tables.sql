--Таблица для хранения базовой информации о пользователях
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(32) NOT NULL UNIQUE,
    password_hash CHAR(60) NOT NULL,
    balance INT DEFAULT 1000 CHECK (balance >= 0)
);

--Таблица со всеми доступными вещами к покупке
CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    price INT NOT NULL CHECK (price >= 0)
);

--Таблица для хранения логов переводов между пользователями
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    sender_id INT REFERENCES users(id),
    receiver_id INT REFERENCES users(id),
    amount INT NOT NULL CHECK (amount >= 0)
);

--Таблица для хранения логов покупок
CREATE TABLE purchases (
    id SERIAL PRIMARY KEY,
    purchase_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    buyer_id INT REFERENCES users(id),
    item_id INT REFERENCES items (id),
    amount INT NOT NULL CHECK (amount >= 0)
);

--Таблица для хранения вещей в инвентарях пользователей
CREATE TABLE user_items (
    user_id INT REFERENCES users(id),
    item_id INT REFERENCES items(id),
    amount INT NOT NULL CHECK (amount >= 0),
    PRIMARY KEY (user_id, item_id)
);
