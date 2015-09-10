
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE task (
    id int NOT NULL AUTO_INCREMENT,
    title varchar(128),
    description varchar(1024),
    priority int DEFAULT 0,
    created timestamp,
    updated timestamp,
    completed timestamp,
    isDeleted boolean DEFAULT false,
    isCompleted boolean DEFAULT false,
    PRIMARY KEY(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE task;
