
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE user (
    `id` int NOT NULL AUTO_INCREMENT,
    `email` varchar(255),
    `pass` varchar(64),
    `name` varchar(127),
    `created` timestamp,
    `updated` timestamp,
    `deleted` timestamp,
    PRIMARY KEY(`id`),
    UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE user;
