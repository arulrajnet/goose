-- +goose Up
-- +goose StatementBegin
INSERT INTO owners(owner_name, owner_type) 
    VALUES ('lucas', 'user'), ('spacey', 'organization');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM owners;
-- +goose StatementEnd
