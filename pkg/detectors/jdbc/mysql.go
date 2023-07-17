package jdbc

import (
	"context"
	"errors"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type mysqlJDBC struct {
	conn     string
	userPass string
	host     string
	database string
	params   string
}

func (s *mysqlJDBC) ping(ctx context.Context) pingResult {
	return ping(ctx, "mysql", isMySQLErrorDeterminate,
		s.conn,
		buildMySQLConnectionString(s.host, s.database, s.userPass, s.params),
		buildMySQLConnectionString(s.host, "", s.userPass, s.params))
}

func buildMySQLConnectionString(host, database, userPass, params string) string {
	conn := host + "/" + database
	if userPass != "" {
		conn = userPass + "@" + conn
	}
	if params != "" {
		conn = conn + "?" + params
	}
	return conn
}

func isMySQLErrorDeterminate(err error) bool {
	return true
}

func parseMySQL(subname string) (jdbc, error) {
	// expected form: [subprotocol:]//[user:password@]HOST[/DB][?key=val[&key=val]]
	hostAndDB, params, _ := strings.Cut(subname, "?")
	if !strings.HasPrefix(hostAndDB, "//") {
		return nil, errors.New("expected host to start with //")
	}
	userPassAndHostAndDB := strings.TrimPrefix(hostAndDB, "//")
	userPass, hostAndDB, found := strings.Cut(userPassAndHostAndDB, "@")
	if !found {
		hostAndDB = userPass
		userPass = ""
	}
	host, database, found := strings.Cut(hostAndDB, "/")
	if !found {
		return nil, errors.New("expected host and database to be separated by /")
	}
	return &mysqlJDBC{
		conn:     subname[2:],
		userPass: userPass,
		host:     host,
		database: database,
		params:   params,
	}, nil
}
