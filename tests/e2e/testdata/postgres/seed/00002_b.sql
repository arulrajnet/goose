-- +goose Up

-- insert 150 more owners
INSERT INTO owners (owner_name, owner_type)
SELECT
	'seed-user-' || round(random() * 1000000),
	(SELECT('{user,organization}'::owner_type []) [MOD(i, 2)+1])
FROM
	generate_series(1, 150) s (i);

-- +goose Down
DELETE FROM owners where owner_name LIKE 'seed-user-%' AND owner_id > 100;
SELECT setval('owners_owner_id_seq', max(owner_id)) FROM owners;
