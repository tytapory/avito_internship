CREATE OR REPLACE FUNCTION transfer_coins(sender_id_param INT, receiver_param VARCHAR(32), transfer_amount_param INT)
    RETURNS VOID AS $$
DECLARE
    sender_balance INT;
    receiver_balance INT;
    receiver_id_param INT;
BEGIN
    IF NOT EXISTS (SELECT 1 FROM users WHERE username = receiver_param) THEN
        RAISE EXCEPTION 'Получатель не существует: %', receiver_param;
    END IF;

    IF transfer_amount_param <= 0 THEN
        RAISE EXCEPTION 'Сумма перевода должна быть > 0';
    END IF;

    SELECT id INTO receiver_id_param FROM users WHERE username = receiver_param;

    IF receiver_id_param = sender_id_param THEN
        RAISE EXCEPTION 'Нельзя переводить средства самому себе';
    END IF;

    IF sender_id_param < receiver_id_param THEN
        SELECT balance INTO sender_balance FROM users WHERE id = sender_id_param FOR UPDATE;
        SELECT balance INTO receiver_balance FROM users WHERE id = receiver_id_param FOR UPDATE;
    ELSE
        SELECT balance INTO receiver_balance FROM users WHERE id = receiver_id_param FOR UPDATE;
        SELECT balance INTO sender_balance FROM users WHERE id = sender_id_param FOR UPDATE;
    END IF;

    IF sender_balance < transfer_amount_param THEN
        RAISE EXCEPTION 'Недостаточно средств на балансе отправителя';
    END IF;

    UPDATE users
    SET balance = balance - transfer_amount_param
    WHERE id = sender_id_param;

    UPDATE users
    SET balance = balance + transfer_amount_param
    WHERE id = receiver_id_param;

    INSERT INTO transactions (sender_id, receiver_id, amount)
    VALUES (sender_id_param, receiver_id_param, transfer_amount_param);
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION buy_item(user_id_param INT, item_name_param VARCHAR(32), item_amount_param INT)
    RETURNS VOID AS $$
DECLARE
    user_balance INT;
    item_price INT;
    item_exists BOOLEAN;
    item_id_param INT;
BEGIN
    SELECT EXISTS (SELECT 1 FROM items WHERE name = item_name_param) INTO item_exists;
    IF NOT item_exists THEN
        RAISE EXCEPTION 'Предмет не существует: %', item_name_param;
    END IF;

    IF item_amount_param <= 0 THEN
        RAISE EXCEPTION 'Количество покупаемых предметов должно быть > 0';
    END IF;

    SELECT balance INTO user_balance FROM users WHERE users.id = user_id_param FOR UPDATE;

    SELECT items.id INTO item_id_param FROM items WHERE name = item_name_param;

    SELECT price INTO item_price FROM items WHERE items.id = item_id_param;

    IF user_balance < item_amount_param * item_price THEN
        RAISE EXCEPTION 'Недостаточно средств на балансе пользователя';
    END IF;

    UPDATE users
    SET balance = balance - item_amount_param * item_price
    WHERE id = user_id_param;

    INSERT INTO user_items (user_id, item_id, amount)
    VALUES (user_id_param, item_id_param, item_amount_param)
    ON CONFLICT (user_id, item_id)
        DO UPDATE SET amount = user_items.amount + EXCLUDED.amount;

    INSERT INTO purchases (buyer_id, item_id, amount)
    VALUES (user_id_param, item_id_param, item_amount_param);
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_user_id_password_hash(username_param VARCHAR(32))
    RETURNS TABLE(id INT, password_hash CHAR(60)) AS $$
BEGIN
    RETURN QUERY
        SELECT users.id, users.password_hash
        FROM users
        WHERE users.username = username_param;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION register_user(username_param VARCHAR(32), password_hash_param CHAR(60))
    RETURNS INT AS $$
DECLARE
    user_id_param INT;
BEGIN
    INSERT INTO users (username, password_hash)
    VALUES (username_param, password_hash_param)
    RETURNING id INTO user_id_param;
    RETURN user_id_param;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_user_balance(user_id_param INT)
    RETURNS INT AS $$
DECLARE
    result INT;
BEGIN
    SELECT users.balance INTO result FROM users WHERE users.id = user_id_param;
    RETURN result;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_user_inventory(user_id_param INT)
    RETURNS TABLE(item_name VARCHAR(32), amount INT) AS $$
BEGIN
    RETURN QUERY
        SELECT items.name, user_items.amount FROM user_items
                                                      JOIN items ON user_items.item_id = items.id
        WHERE user_items.user_id = user_id_param;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_user_receive_history(user_id INT)
    RETURNS TABLE(user_from VARCHAR(32), amount INT) AS $$
BEGIN
    RETURN QUERY
        SELECT users.username, transactions.amount
        FROM transactions
                 JOIN users ON users.id = transactions.sender_id
        WHERE transactions.receiver_id = user_id;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_user_send_history(user_id INT)
    RETURNS TABLE(user_from VARCHAR(32), amount INT) AS $$
BEGIN
    RETURN QUERY
        SELECT users.username, transactions.amount
        FROM transactions
                 JOIN users ON users.id = transactions.receiver_id
        WHERE transactions.sender_id = user_id;
END;
$$ LANGUAGE plpgsql;