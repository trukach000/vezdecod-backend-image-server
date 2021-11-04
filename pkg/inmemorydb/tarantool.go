package inmemorydb

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/tarantool/go-tarantool"
)

func InitTarantool(
	host string,
	port string,
	user string,
	pass string,
) (*tarantool.Connection, error) {
	server := fmt.Sprintf("%s:%s", host, port)
	opts := tarantool.Opts{
		Timeout:       500 * time.Millisecond,
		Reconnect:     1 * time.Second,
		MaxReconnects: 3,
		User:          user,
		Pass:          pass,
	}

	logrus.Infof("trying to connect to tarantool with credentials: %s:%s", user, pass)

	client, err := tarantool.Connect(server, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the tarantool: %w", err)
	}

	_, err = client.Ping()
	if err != nil {
		return nil, fmt.Errorf("can't ping tarantool: %w", err)
	}

	return client, nil
}
