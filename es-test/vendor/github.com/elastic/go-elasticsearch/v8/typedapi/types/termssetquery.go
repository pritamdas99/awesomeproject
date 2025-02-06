// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.


// Code generated from the elasticsearch-specification DO NOT EDIT.
// https://github.com/elastic/elasticsearch-specification/tree/4316fc1aa18bb04678b156f23b22c9d3f996f9c9


package types

// TermsSetQuery type.
//
// https://github.com/elastic/elasticsearch-specification/blob/4316fc1aa18bb04678b156f23b22c9d3f996f9c9/specification/_types/query_dsl/term.ts#L139-L143
type TermsSetQuery struct {
	Boost                    *float32 `json:"boost,omitempty"`
	MinimumShouldMatchField  *Field   `json:"minimum_should_match_field,omitempty"`
	MinimumShouldMatchScript *Script  `json:"minimum_should_match_script,omitempty"`
	QueryName_               *string  `json:"_name,omitempty"`
	Terms                    []string `json:"terms"`
}

// TermsSetQueryBuilder holds TermsSetQuery struct and provides a builder API.
type TermsSetQueryBuilder struct {
	v *TermsSetQuery
}

// NewTermsSetQuery provides a builder for the TermsSetQuery struct.
func NewTermsSetQueryBuilder() *TermsSetQueryBuilder {
	r := TermsSetQueryBuilder{
		&TermsSetQuery{},
	}

	return &r
}

// Build finalize the chain and returns the TermsSetQuery struct
func (rb *TermsSetQueryBuilder) Build() TermsSetQuery {
	return *rb.v
}

func (rb *TermsSetQueryBuilder) Boost(boost float32) *TermsSetQueryBuilder {
	rb.v.Boost = &boost
	return rb
}

func (rb *TermsSetQueryBuilder) MinimumShouldMatchField(minimumshouldmatchfield Field) *TermsSetQueryBuilder {
	rb.v.MinimumShouldMatchField = &minimumshouldmatchfield
	return rb
}

func (rb *TermsSetQueryBuilder) MinimumShouldMatchScript(minimumshouldmatchscript *ScriptBuilder) *TermsSetQueryBuilder {
	v := minimumshouldmatchscript.Build()
	rb.v.MinimumShouldMatchScript = &v
	return rb
}

func (rb *TermsSetQueryBuilder) QueryName_(queryname_ string) *TermsSetQueryBuilder {
	rb.v.QueryName_ = &queryname_
	return rb
}

func (rb *TermsSetQueryBuilder) Terms(terms ...string) *TermsSetQueryBuilder {
	rb.v.Terms = terms
	return rb
}
