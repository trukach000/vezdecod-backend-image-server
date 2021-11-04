package images

import (
	"backend-image-server/pkg/database"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrImageNotFound = errors.New("Image not found")
)

func SaveNewImage(ctx context.Context, data []byte) (string, error) {
	db, err := database.GetDatabaseFromContext(ctx)
	if err != nil {
		return "", err
	}

	newUuid := uuid.New()

	_, err = db.ExecContext(
		ctx,
		`INSERT INTO images
				(
					token,
					data
				)
				VALUES
					(
						?,
						?
					);
				`,
		newUuid.String(),
		data,
	)
	if err != nil {
		return "", err
	}

	return newUuid.String(), nil
}

func GetImageByToken(
	ctx context.Context,
	token string,
) ([]byte, error) {

	var res []byte

	db, err := database.GetDatabaseFromContext(ctx)
	if err != nil {
		return nil, err
	}

	err = db.QueryRowContext(
		ctx,
		`	SELECT 
				data 
			FROM 
				images 
			WHERE
				token = ?
		`,
		token,
	).Scan(&res)

	if err == sql.ErrNoRows {
		return nil, ErrImageNotFound
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}
