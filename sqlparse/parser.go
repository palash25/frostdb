package sqlparse

import (
	"fmt"

	"github.com/pingcap/tidb/parser"

	"github.com/polarsignals/frostdb/query"
)

type Parser struct {
	p *parser.Parser
}

func NewParser() *Parser {
	return &Parser{p: parser.New()}
}

// ExperimentalParse uses the provided query builder to build a FrostDB query
// specified using the provided SQL.
// TODO(asubiotto): This API will change over time. Currently,
// queryEngine.ScanTable is provided as a starting point and no table needs to
// be specified in the SQL statement. Additionally, the idea is to change to
// creating logical plans directly (rather than through a builder).
func (p *Parser) ExperimentalParse(builder query.Builder, dynColNames []string, sql string) (query.Builder, error) {
	asts, _, err := p.p.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}

	if len(asts) != 1 {
		return nil, fmt.Errorf("cannot handle multiple asts, found %d", len(asts))
	}

	v := newASTVisitor(builder, dynColNames)
	asts[0].Accept(v)
	if v.err != nil {
		return nil, v.err
	}

	return v.builder, nil
}
