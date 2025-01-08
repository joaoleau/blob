package repository

const (
	createBlobQuery = `INSERT INTO blobs (blob_id, user_id, description, title, created_at, updated_at)
					   VALUES ($1, $2, $3, $4, now(), now())
					   RETURNING *`

	updateBlobQuery = `UPDATE blobs
					   SET description = COALESCE(NULLIF($1, ''), description),
						   title = COALESCE(NULLIF($2, ''), title),
						   updated_at = now()
					   WHERE blob_id = $3
					   RETURNING *`

	getBlobByIDQuery = `SELECT blob_id, user_id, description, title, created_at, updated_at
					    FROM blobs
					    WHERE blob_id = $1`

	listBlobsQuery = `SELECT blob_id, user_id, description, title, created_at, updated_at
					  FROM blobs
					  ORDER BY created_at DESC
					  LIMIT $1 OFFSET $2`

	deleteBlobQuery = `DELETE FROM blobs
					   WHERE blob_id = $1
					   RETURNING blob_id`

	getTotalBlob = `SELECT COUNT(blob_id) FROM blob`
)
