package pgxdb

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/barbell-math/util/src/customerr"
	"github.com/barbell-math/util/src/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	ConnPoolComputer struct {
		ConnInfo *db.ArgparseVals
	}
)

func (c ConnPoolComputer) ComputeVals() (*pgxpool.Pool, error) {
	sb := strings.Builder{}
	sb.WriteString("postgres://")

	if c.ConnInfo.User != "" {
		sb.WriteString(c.ConnInfo.User)
		if c.ConnInfo.EnvPswdVar != "" {
			if pswd, ok := os.LookupEnv(c.ConnInfo.EnvPswdVar); ok {
				sb.WriteString(":")
				sb.WriteString(pswd)
			}
		}
		sb.WriteString("@")
	}

	if c.ConnInfo.NetLoc != "" {
		sb.WriteString(c.ConnInfo.NetLoc)
	}
	if c.ConnInfo.Port != 0 {
		sb.WriteString(":")
		sb.WriteString(strconv.Itoa(int(c.ConnInfo.Port)))
	}

	if c.ConnInfo.DBName != "" {
		sb.WriteString("/")
		sb.WriteString(c.ConnInfo.DBName)
	}

	rv, err := pgxpool.New(context.Background(), sb.String())
	if err != nil {
		return rv, customerr.AppendError(UnableToCreateConnPoolErr, err)
	}
	return rv, nil
}

func (c ConnPoolComputer) Reset() {
	// intentional noop - ConnPoolComputer has no state that needs to be reset
}
