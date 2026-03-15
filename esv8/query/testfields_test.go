package query_test

import "github.com/tomtwinkle/es-typed-go/estype"

// Field constants for test use, mirroring the output of cmd/estyped.

const FieldStatus estype.Field = "status"
const FieldCategory estype.Field = "category"
const FieldTitle estype.Field = "title"
const FieldTitleNgram estype.Field = "title.ngram"
const FieldTitleRaw estype.Field = "title.raw"
const FieldTags estype.Field = "tags"
const FieldItems estype.Field = "items"
const FieldItemsColor estype.Field = "items.color"
const FieldItemsStatus estype.Field = "items.status"
const FieldItemsIds estype.Field = "items.ids"
const FieldItemsDate estype.Field = "items.date"
const FieldItemsPrice estype.Field = "items.price"
const FieldItemsLocation estype.Field = "items.location"
const FieldDate estype.Field = "date"
const FieldPrice estype.Field = "price"
const FieldEnabled estype.Field = "enabled"
const FieldType estype.Field = "type"
const FieldName estype.Field = "name"
const FieldNameKeyword estype.Field = "name.keyword"
const FieldId estype.Field = "id"
const FieldLocation estype.Field = "location"
const FieldValue estype.Field = "value"
