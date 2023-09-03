package xweb

import "github.com/caumette-co/x/xweb/params"

var StringsParamExtractors = []params.StringsParamExtractor{
	params.PathExtractor{},
	params.QueryExtractor{},
	params.HeaderExtractor{},
}

var ValuesParamExtractors = []params.ValueParamExtractor{
	params.ContextExtractor{},
}
