-- +goose Up
-- +goose StatementBegin
INSERT INTO owners (owner_name, owner_type)
SELECT
	'user-' || round(random() * 1000000),
	'user'
FROM
	generate_series(1, 1000);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE 
    owners,
    repos,
    stargazers,
    issues;
-- +goose StatementEnd