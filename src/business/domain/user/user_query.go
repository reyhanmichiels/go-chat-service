package user

const (
	insertUser = `
		INSERT INTO user
		(
			fk_role_id,
		 	name,
		 	email,
		 	password,
		 	created_at,
		 	created_by
		)
		VALUES
		(
			:fk_role_id,
		 	:name,
		 	:email,
		 	:password,
		 	:created_at,
		 	:created_by
		)
	`

	readUser = `
		SELECT
		    id,
			fk_role_id,
		 	name,
		 	email,
		 	password,
			status,
			flag,
			meta,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			user
	`

	countUser = `
		SELECT
			COUNT(*)
		FROM
			user
	`

	updateUser = `
		UPDATE
			user
	`
)
