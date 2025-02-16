package repository

const (
	createBlobQuery = `
		INSERT INTO "Blob" (id, user_id, content, created_at, updated_at)
		VALUES ($1, $2, $3, now(), now())
		RETURNING *`

	updateBlobQuery = `
		UPDATE "Blob"
		SET content = COALESCE(NULLIF($1, ''), content),
			updated_at = now()
		WHERE id = $2
		RETURNING *`

	getBlobByIDQuery = `
	SELECT 
		b.id AS blob_id, 
		b.user_id AS blob_user_id, 
		b.content AS blob_content, 
		b.created_at AS blob_created_at, 
		b.updated_at AS blob_updated_at,
		u.username AS user_username, 
		u.avatar_icon AS user_avatar_icon, 
		u.created_at AS user_created_at,
		c.id AS comment_id, 
		c.content AS comment_content, 
		c.created_at AS comment_created_at, 
		c.updated_at AS comment_updated_at, 
		c.user_id AS comment_user_id, 
		l.id AS like_id, 
		l.user_id AS like_user_id, 
		l.created_at AS like_created_at,
		i.id AS interest_id,
		i.name AS interest_name, 
		i.description AS interest_description, 
		i.created_at AS interest_created_at,
		i.updated_at AS interest_updated_at
	FROM "Blob" b
	LEFT JOIN "User" u ON b.user_id = u.id
	LEFT JOIN "Comment" c ON c.blob_id = b.id
	LEFT JOIN "Like" l ON l.blob_id = b.id
	LEFT JOIN "_BlobToInterest" bi ON bi.blob_id = b.id
	LEFT JOIN "Interest" i ON i.id = bi.interest_id
	WHERE b.id = $1;
	`

	deleteBlobQuery = `
		DELETE FROM "Blob"
		WHERE id = $1
		RETURNING id`

	listBlobsQuery = `
		SELECT * FROM listBlobs`

	getTotalBlob = `
		SELECT COUNT(id)
		FROM "Blob"`

	getTotalUser = `
		SELECT COUNT(id)
		FROM "User"`

	listUserByUsernameWithBlobsQuery = `
		SELECT
			u.id,
			u.name,
			u.email,
			u.email_verified,
			u.image,
			u.username,
			u.bio,
			u.avatar_icon,
			u.avatar_color,
			u.created_at,
			u.updated_at,
			b.id AS blob_id,
			b.content AS blob_content,
			b.created_at AS blob_created_at,
			b.updated_at AS blob_updated_at
		FROM
			"User" u
		LEFT JOIN
			"Blob" b ON b.user_id = u.id
		WHERE u.username = $1
		`

	getUserByEmail = `
		SELECT 
			u.id,
			u.name,
			u.email,
			u.email_verified,
			u.image,
			u.username,
			u.bio,
			u.avatar_icon,
			u.avatar_color,
			u.created_at,
			u.updated_at
		FROM "User" u
		WHERE u.email = $1
		`

	getUserByID = `
		SELECT
			u.id,
			u.name,
			u.email,
			u.email_verified,
			u.image,
			u.username,
			u.bio,
			u.avatar_icon,
			u.avatar_color,
			u.created_at,
			u.updated_at
		FROM "User" u
		WHERE u.id = $1
		`

	insertLikeQuery = `
		INSERT INTO "Like" (id, user_id, blob_id)
		VALUES ($1, $2, $3)
		RETURNING id;`

	deleteLikeQuery = `
		DELETE FROM "Like"
		WHERE id = $1
		AND user_id = $2
		AND blob_id = $3
		RETURNING id;
	`

	searchLikeQuery = `
		SELECT id
		FROM "Like"
		WHERE user_id = $1 AND blob_id = $2;
	`
	searchLikebyBlobIDQuery = `
		SELECT 
			l.id, 
			l.created_at, 
			l.user_id, 
			l.blob_id, 
			u.image, 
			u.username, 
			u.avatar_icon, 
			u.avatar_color
		FROM 
			"Like" l
		JOIN 
			"User" u ON l.user_id = u.id
		WHERE 
			l.blob_id = $1;
	`

	insertCommentQuery = `
		INSERT INTO "Comment" (id, content, user_id, blob_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
		`

	insertBlobInterest = `
		INSERT INTO "_BlobToInterest" (blob_id, interest_id)
		VALUES ($1, $2)
	`

	searchCommentsbyBlobIDQuery = `
	SELECT 
		c.id, 
		c.content, 
		c.created_at, 
		c.updated_at, 
		c.user_id, 
		u.image, 
		u.username, 
		u.avatar_icon, 
		u.avatar_color, 
		c.blob_id
	FROM 
		"Comment" c
	JOIN 
		"User" u ON c.user_id = u.id
	WHERE 
		c.blob_id = $1;
	`

	deleteCommentQuery = `
		DELETE FROM "Comment"
		WHERE id = $1 AND user_id = $2
		RETURNING id;
		`
)
