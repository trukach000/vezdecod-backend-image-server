package images

import (
	"backend-image-server/pkg/database"
	"backend-image-server/pkg/redisclient"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	ErrImageNotFound = errors.New("Image not found")
)

func GetRedisKey(hashP string, aspectRatio float64) string {
	return fmt.Sprintf("%s_%.2f", hashP, aspectRatio)
}

func SaveNewImage(
	ctx context.Context,
	hashP string,
	width int,
	height int,
	aspectRatio float64,
	data []byte,
) (string, error) {
	db, err := database.GetDatabaseFromContext(ctx)
	if err != nil {
		return "", err
	}

	red, err := redisclient.GetRedisFromContext(ctx)
	if err != nil {
		return "", err
	}

	newUuid := uuid.New()
	token := newUuid.String()

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
		token,
		data,
	)
	if err != nil {
		return "", err
	}

	value, err := json.Marshal(redisclient.ImgData{
		AspectRatio: aspectRatio,
		W:           int64(width),
		H:           int64(height),
		Token:       token,
	})
	if err != nil {
		return "", fmt.Errorf("Json marshal erorr during saving to the redis: %w", err)
	}

	_, err = red.Set(GetRedisKey(hashP, aspectRatio), value, 0).Result()
	if err != nil {
		return "", fmt.Errorf("Redis set erorr during saving to the redis: %w", err)
	}

	return token, nil
}

func ReplaceImage(
	ctx context.Context,
	token string,
	hashP string,
	width int,
	height int,
	aspectRatio float64,
	data []byte,
) error {

	db, err := database.GetDatabaseFromContext(ctx)
	if err != nil {
		return err
	}

	red, err := redisclient.GetRedisFromContext(ctx)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(
		ctx,
		`REPLACE INTO images
				(
					token,
					data
				)
				VALUES
					(
						?,
						?
					)
				`,
		token,
		data,
	)

	if err != nil {
		return err
	}

	value, err := json.Marshal(redisclient.ImgData{
		AspectRatio: aspectRatio,
		W:           int64(width),
		H:           int64(height),
		Token:       token,
	})

	logrus.Infof("VALUE: %+v", string(value))

	if err != nil {
		return fmt.Errorf("Json marshal erorr during replacing to the redis: %w", err)
	}

	_, err = red.Set(GetRedisKey(hashP, aspectRatio), value, 0).Result()
	if err != nil {
		return fmt.Errorf("Redis set erorr during replacing to the redis: %w", err)
	}

	return nil
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
