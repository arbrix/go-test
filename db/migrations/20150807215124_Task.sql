
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE task (
    id int NOT NULL AUTO_INCREMENT,
    title varchar(128),
    description varchar(1024),
    priority int DEFAULT 0,
    created datetime,
    updated datetime,
    completed datetime,
    isDeleted bool DEFAULT false,
    isCompleted bool DEFAULT false,
    PRIMARY KEY(id)
) ENGINE=InnoDB;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE task;