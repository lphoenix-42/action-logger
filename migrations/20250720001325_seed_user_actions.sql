-- +goose Up
-- +goose StatementBegin
insert into user_actions (user_id, action_type, timestamp, details) values
(101, 1, NOW() - INTERVAL '1 day', '{"item": "Book", "price": 12.99}'),
(102, 2, '2025-07-01 10:00:00+00', '{"item": "Pen", "reason": "Broken"}'),
(103, 1, NOW() - INTERVAL '6 hours', '{"item": "Book", "discount": true}');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
delete from user_actions
where
    (user_id = 101 and action_type = 1) or
    (user_id = 102 and action_type = 2) or
    (user_id = 103 and action_type = 1);
-- +goose StatementEnd
