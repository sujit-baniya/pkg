package fts

import (
	"fmt"
	"github.com/goccy/go-reflect"
	"github.com/sujit-baniya/pkg/maps"
	"github.com/sujit-baniya/pkg/str"
	"github.com/sujit-baniya/xid"
	"regexp"
	"strings"
	"time"
)

var punctuationRegex = regexp.MustCompile(`[^\w|\s]`)

var stopWords = map[string]bool{
	"a":          true,
	"about":      true,
	"above":      true,
	"after":      true,
	"again":      true,
	"against":    true,
	"all":        true,
	"am":         true,
	"an":         true,
	"and":        true,
	"any":        true,
	"are":        true,
	"aren't":     true,
	"as":         true,
	"at":         true,
	"be":         true,
	"because":    true,
	"been":       true,
	"before":     true,
	"being":      true,
	"below":      true,
	"between":    true,
	"both":       true,
	"but":        true,
	"by":         true,
	"can't":      true,
	"cannot":     true,
	"could":      true,
	"couldn't":   true,
	"did":        true,
	"didn't":     true,
	"do":         true,
	"does":       true,
	"doesn't":    true,
	"doing":      true,
	"don't":      true,
	"down":       true,
	"during":     true,
	"each":       true,
	"few":        true,
	"for":        true,
	"from":       true,
	"further":    true,
	"had":        true,
	"hadn't":     true,
	"has":        true,
	"hasn't":     true,
	"have":       true,
	"haven't":    true,
	"having":     true,
	"he":         true,
	"he'd":       true,
	"he'll":      true,
	"he's":       true,
	"her":        true,
	"here":       true,
	"here's":     true,
	"hers":       true,
	"herself":    true,
	"him":        true,
	"himself":    true,
	"his":        true,
	"how":        true,
	"how's":      true,
	"i":          true,
	"i'd":        true,
	"i'll":       true,
	"i'm":        true,
	"i've":       true,
	"if":         true,
	"in":         true,
	"into":       true,
	"is":         true,
	"isn't":      true,
	"it":         true,
	"it's":       true,
	"its":        true,
	"itself":     true,
	"let's":      true,
	"me":         true,
	"more":       true,
	"most":       true,
	"mustn't":    true,
	"my":         true,
	"myself":     true,
	"no":         true,
	"nor":        true,
	"not":        true,
	"of":         true,
	"off":        true,
	"on":         true,
	"once":       true,
	"only":       true,
	"or":         true,
	"other":      true,
	"ought":      true,
	"our":        true,
	"ours":       true,
	"ourselves":  true,
	"out":        true,
	"over":       true,
	"own":        true,
	"same":       true,
	"shan't":     true,
	"she":        true,
	"she'd":      true,
	"she'll":     true,
	"she's":      true,
	"should":     true,
	"shouldn't":  true,
	"so":         true,
	"some":       true,
	"such":       true,
	"than":       true,
	"that":       true,
	"that's":     true,
	"the":        true,
	"their":      true,
	"theirs":     true,
	"them":       true,
	"themselves": true,
	"then":       true,
	"there":      true,
	"there's":    true,
	"these":      true,
	"they":       true,
	"they'd":     true,
	"they'll":    true,
	"they're":    true,
	"they've":    true,
	"this":       true,
	"those":      true,
	"through":    true,
	"to":         true,
	"too":        true,
	"under":      true,
	"until":      true,
	"up":         true,
	"very":       true,
	"was":        true,
	"wasn't":     true,
	"we":         true,
	"we'd":       true,
	"we'll":      true,
	"we're":      true,
	"we've":      true,
	"were":       true,
	"weren't":    true,
	"what":       true,
	"what's":     true,
	"when":       true,
	"when's":     true,
	"where":      true,
	"where's":    true,
	"which":      true,
	"while":      true,
	"who":        true,
	"who's":      true,
	"whom":       true,
	"why":        true,
	"why's":      true,
	"with":       true,
	"won't":      true,
	"would":      true,
	"wouldn't":   true,
	"you":        true,
	"you'd":      true,
	"you'll":     true,
	"you're":     true,
	"you've":     true,
	"your":       true,
	"yours":      true,
	"yourself":   true,
	"yourselves": true,
}

type SchemaProps any

type Record[Schema SchemaProps] struct {
	Id string
	S  Schema
}

type RecordInfo struct {
	recId string
	freq  uint32
}

type Option struct {
	Exact bool
	Size  int
}

type FTS[Schema SchemaProps] struct {
	key   string
	docs  maps.IMap[string, Schema]
	index maps.IMap[string, []RecordInfo]
	rules map[string]bool
}

var defaultSize = 20

func New[Schema SchemaProps](key string, rules ...map[string]bool) *FTS[Schema] {
	var r map[string]bool
	if len(rules) > 0 {
		r = rules[0]
	}
	return &FTS[Schema]{
		key:   key,
		docs:  maps.New[string, Schema](),
		index: maps.New[string, []RecordInfo](),
		rules: r,
	}
}

func (db *FTS[Schema]) Insert(doc Schema) (Record[Schema], error) {
	id := xid.New().String()
	db.docs.Set(id, doc)
	db.indexDocument(id, doc)
	return Record[Schema]{Id: id, S: doc}, nil
}

func (db *FTS[Schema]) IndexLen() uintptr {
	return db.index.Len()
}

func (db *FTS[Schema]) DocumentLen() uintptr {
	return db.docs.Len()
}

func (db *FTS[Schema]) InsertBatch(docs []Schema) []error {
	errs := make([]error, 0)
	for _, d := range docs {
		if _, err := db.Insert(d); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (db *FTS[Schema]) Update(id string, doc Schema) (Record[Schema], error) {
	prevDoc, ok := db.docs.Get(id)
	if !ok {
		return Record[Schema]{}, fmt.Errorf("document not found")
	}
	db.deIndexDocument(id, prevDoc)
	db.docs.Set(id, doc)
	db.indexDocument(id, doc)
	return Record[Schema]{Id: id, S: doc}, nil
}

func (db *FTS[Schema]) Delete(id string) error {
	doc, ok := db.docs.Get(id)
	if !ok {
		return fmt.Errorf("document not found")
	}
	db.deIndexDocument(id, doc)
	db.docs.Del(id)
	return nil
}

func (db *FTS[Schema]) Search(query string, params ...Option) []Record[Schema] {
	option := Option{Size: defaultSize, Exact: true}
	if len(params) > 0 {
		option = params[0]
	}
	recordsIds := make(map[string]int)
	records := make([]Record[Schema], 0)
	tokens := Tokenize(query)
	for _, token := range tokens {
		infos, _ := db.index.Get(token)
		for _, info := range infos {
			recordsIds[info.recId] += 1
		}
	}
	i := 0
	for id, tokensCount := range recordsIds {
		if !option.Exact || tokensCount == len(tokens) {
			if option.Size == 0 || (option.Size > 0 && i < option.Size) {
				i++
				doc, _ := db.docs.Get(id)
				records = append(records, Record[Schema]{Id: id, S: doc})
			}
		}
	}

	return records
}

func (db *FTS[Schema]) SearchExact(query string, size ...int) []Record[Schema] {
	s := defaultSize
	if len(size) > 0 {
		s = size[0]
	}
	return db.Search(query, Option{Size: s, Exact: true})
}

func (db *FTS[Schema]) indexDocument(id string, doc Schema) {
	text := strings.Join(db.getIndexFields(doc), " ")
	tokens := Tokenize(text)
	tokensCount := Count(tokens)

	for token, count := range tokensCount {
		recordsInfos, _ := db.index.GetOrSet(token, []RecordInfo{})
		recordsInfos = append(recordsInfos, RecordInfo{id, count})
		db.index.Set(token, recordsInfos)
	}
}

func (db *FTS[Schema]) deIndexDocument(id string, doc Schema) {
	text := strings.Join(db.getIndexFields(doc), " ")
	tokens := Tokenize(text)

	for _, token := range tokens {
		if recordsInfos, ok := db.index.Get(token); ok {
			var newRecordsInfos []RecordInfo
			for _, info := range recordsInfos {
				if info.recId != id {
					newRecordsInfos = append(newRecordsInfos, info)
				}
			}
			db.index.Set(token, newRecordsInfos)
		}
	}
}

func (db *FTS[Schema]) getIndexFields(obj any) (fields []string) {
	switch v := obj.(type) {
	case string, bool, time.Time, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		fields = append(fields, fmt.Sprintf("%v", v))
	case map[string]any:
		rules := make(map[string]bool)
		if db.rules != nil {
			rules = db.rules
		}
		for field, val := range v {
			if len(rules) > 0 {
				if canIndex, ok := rules[field]; ok && canIndex {
					fields = append(fields, fmt.Sprintf("%v", val))
				}
			} else {
				fields = append(fields, fmt.Sprintf("%v", val))
			}
		}
	default:
		val := reflect.ValueOf(obj)
		t := reflect.TypeOf(obj)
		hasIndexField := false
		for i := 0; i < val.NumField(); i++ {
			f := t.Field(i)
			if v, ok := f.Tag.Lookup("index"); ok && str.EqualFold(v, "true") {
				hasIndexField = true
				fields = append(fields, val.Field(i).String())
			}
		}
		if !hasIndexField {
			for i := 0; i < val.NumField(); i++ {
				fields = append(fields, val.Field(i).String())
			}
		}
	}
	return
}

func Tokenize(data string) []string {
	data = punctuationRegex.ReplaceAllString(data, "")
	data = str.ToLower(data)
	arr := strings.Fields(data)
	noStopWords := removeStopWords(arr)
	return uniqueSlice(noStopWords)
}

func Count(tokens []string) map[string]uint32 {
	dict := make(map[string]uint32)
	for _, token := range tokens {
		dict[token]++
	}
	return dict
}

func removeStopWords(tokens []string) []string {
	var newSlice []string
	for _, value := range tokens {
		_, ok := stopWords[value]
		if !ok {
			newSlice = append(newSlice, value)
		}
	}

	return newSlice
}

func uniqueSlice(tokens []string) []string {
	tokenHash := make(map[string]bool)
	var newSlice []string

	for _, token := range tokens {
		if _, ok := tokenHash[token]; !ok {
			tokenHash[token] = true
			newSlice = append(newSlice, token)
		}
	}

	return newSlice
}

func copyAppend(slices [][]byte) []byte {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	tmp := make([]byte, totalLen)
	var i int
	for _, s := range slices {
		i += copy(tmp[i:], s)
	}
	return tmp
}
