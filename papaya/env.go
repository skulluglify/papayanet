package papaya

import (
	"errors"
	"github.com/joho/godotenv"
)

type PnEnvImpl interface {
	Load() error
}

type PnDotEnv struct{}

func (env *PnDotEnv) Load() error {

	if err := godotenv.Load(); err != nil {

		return errors.New("cannot loaded `.env` file")
	}

	return nil
}
