package instana

import (
	"context"
	"fmt"
	instana "github.com/instana/go-sensor"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"strings"
)

const parentSpanKey = "::instana::gorm::parentspan::"
const spanKey = "::instana::gorm::span::"

const (
	pq_prop_user = "user"
	pq_prop_host = "host"
	pq_prop_port = "port"
	pq_prop_db   = "database"
)

type adapter struct {
	driver   string
	peer     string
	user     string
	database string
}

func NewAdapter(db *gorm.DB, driver, connstring string) (func(context.Context, *gorm.DB) *gorm.DB, error) {
	if driver == "postgres" {
		driver = "postgresql"
	}

	peer, user, database, err := extractConnectionDetails(driver, connstring)
	if err != nil {
		return nil, err
	}

	a := adapter{driver, peer, user, database}
	db.Callback().Create().Before("gorm:create").Register("instana:create:before", a.before)
	db.Callback().Query().Before("gorm:query").Register("instana:query:before", a.before)
	db.Callback().Update().Before("gorm:update").Register("instana:update:before", a.before)
	db.Callback().Delete().Before("gorm:delete").Register("instana:delete:before", a.before)
	db.Callback().RowQuery().Before("gorm:rowquery").Register("instana:rowquery:before", a.before)

	db.Callback().Create().After("gorm:create").Register("instana:create:after", a.before)
	db.Callback().Query().After("gorm:query").Register("instana:query:after", a.before)
	db.Callback().Update().After("gorm:update").Register("instana:update:after", a.before)
	db.Callback().Delete().After("gorm:delete").Register("instana:delete:after", a.before)
	db.Callback().RowQuery().After("gorm:rowquery").Register("instana:rowquery:after", a.before)

	return a.adaptDB, nil
}

func (a adapter) adaptDB(ctx context.Context, db *gorm.DB) *gorm.DB {
	if ctx == nil {
		return db
	}
	parentSpan, ok := instana.SpanFromContext(ctx)
	if !ok {
		return db
	}
	return db.Set(parentSpanKey, parentSpan)
}

func (a adapter) beforeCreate(scope *gorm.Scope) {
	a.before(scope)
}

func (a adapter) afterCreate(scope *gorm.Scope) {
	a.after(scope)
}

func (a adapter) before(scope *gorm.Scope) {
	v, ok := scope.Get(parentSpanKey)
	if !ok {
		return
	}
	parentSpan, ok := v.(opentracing.Span)
	if !ok {
		return
	}

	span := parentSpan.Tracer().StartSpan(scope.SQL, opentracing.ChildOf(parentSpan.Context()))
	span.SetTag(string(ext.SpanKind), string(ext.SpanKindRPCClientEnum))
	span.SetTag(string(ext.DBType), a.driver)
	span.SetTag(string(ext.DBInstance), a.database)
	span.SetTag(string(ext.DBUser), a.user)
	span.SetTag(string(ext.DBStatement), scope.SQL)
	span.SetTag(string(ext.PeerAddress), a.peer)

	scope.Set(spanKey, span)
}

func (a adapter) after(scope *gorm.Scope) {
	v, ok := scope.Get(spanKey)
	if !ok {
		return
	}

	span, ok := v.(opentracing.Span)
	if !ok {
		return
	}

	if scope.HasError() {
		span.SetTag(string(ext.Error), scope.DB().Error)
	}
	span.Finish()
}

func extractConnectionDetails(driver, connstring string) (peer, user, database string, err error) {
	switch driver {
	case "postgres":
		parsed, err := pq.ParseURL(connstring)
		if err != nil {
			return "", "", "", err
		}
		split := strings.Split(parsed, " ")

		var host, port string
		for _, s := range split {
			if strings.HasPrefix(s, pq_prop_user) {
				user = extractValue(s)
			} else if strings.HasPrefix(s, pq_prop_db) {
				database = extractValue(s)
			} else if strings.HasPrefix(s, pq_prop_host) {
				host = extractValue(s)
			} else if strings.HasPrefix(s, pq_prop_port) {
				host = extractValue(s)
			}
		}

		if port != "" {
			peer = fmt.Sprintf("%s:%s", host, port)
		} else {
			peer = host
		}
	}
	return
}

func extractValue(v string) string {
	tokens := strings.Split(v, "=")
	if len(tokens) > 1 {
		return tokens[1]
	}
	return ""
}
