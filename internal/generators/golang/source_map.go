package golang

import (
	"github.com/deronyan-llc/dwim/internal/common"
	"github.com/deronyan-llc/rdf/rdf"
)

var (
	SourceMap = make(map[rdf.Term]*common.SchemaContext)
)
