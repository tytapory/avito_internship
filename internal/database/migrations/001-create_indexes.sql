CREATE INDEX idx_username ON users (username);

CREATE INDEX idx_sender_id ON transactions (sender_id);
CREATE INDEX idx_receiver_id ON transactions (receiver_id);

CREATE INDEX idx_buyer_id ON purchases (buyer_id);
CREATE INDEX idx_item_id ON purchases (item_id);

CREATE INDEX idx_user_item_inventory ON user_items (user_id, item_id);


