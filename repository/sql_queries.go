package repository

const (
	createUserQuery = `INSERT INTO users (nickname, email, password, phone_number,
	               		created_at, updated_at, login_date)
						VALUES ($1, $2, $3, $4, now(), now(), now()) 
						RETURNING *`

	updateUserQuery = `UPDATE users
						SET nickname = COALESCE(NULLIF($1, ''), nickname),
						    email = COALESCE(NULLIF($2, ''), email),
						    phone_number = COALESCE(NULLIF($3, ''), phone_number),
						    updated_at = now()
						WHERE user_id = $4
						RETURNING *
						`

	deleteUserQuery = `DELETE FROM users WHERE user_id = $1`

	getUserQuery = `SELECT user_id, nickname, email, phone_number, 
       				 created_at, updated_at, login_date  
					 FROM users 
					 WHERE user_id = $1`

	findUserByEmail = `SELECT user_id, nickname, email, phone_number, 
       			 		created_at, updated_at, login_date, password
				 		FROM users 
				 		WHERE email = $1`

	getUsers = `SELECT user_id, nickname, email, phone_number, 
       			 created_at, updated_at, login_date
				 FROM users 
				 ORDER BY COALESCE(NULLIF($1, ''), nickname) OFFSET $2 LIMIT $3`


	getTotalCount = `SELECT COUNT(user_id) FROM users 
						WHERE nickname ILIKE '%' || $1 || '%'`

	findUsers = `SELECT user_id, nickname, email, phone_number,
	              created_at, updated_at, login_date 
				  FROM users 
				  WHERE nickname ILIKE '%' || $1 || '%'
				  ORDER BY nickname
				  OFFSET $2 LIMIT $3
				  `

	getTotal = `SELECT COUNT(user_id) FROM users`

)