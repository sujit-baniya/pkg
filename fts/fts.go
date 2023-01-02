package fts

import (
	"fmt"
	"github.com/goccy/go-reflect"
	"github.com/google/uuid"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

var numCPU = runtime.NumCPU()

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

type MemDB[Schema SchemaProps] struct {
	docs  *HashMap[string, Schema]
	index *HashMap[string, []RecordInfo]
}

func New[Schema SchemaProps]() *MemDB[Schema] {
	return &MemDB[Schema]{
		docs:  NewMap[string, Schema](),
		index: NewMap[string, []RecordInfo](),
	}
}

func (db *MemDB[Schema]) Insert(doc Schema) (Record[Schema], error) {
	id := uuid.NewString()
	db.docs.Put(id, doc)

	fields := getIndexFields(doc)
	for _, field := range fields {
		db.indexField(id, field)
	}

	return Record[Schema]{Id: id, S: doc}, nil
}

func (db *MemDB[Schema]) InsertBatchSync(docs []Schema) []error {
	errs := make([]error, 0)

	for _, d := range docs {
		if _, err := db.Insert(d); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (db *MemDB[Schema]) IndexLen() int {
	return db.index.Len()
}

func (db *MemDB[Schema]) InsertBatchAsync(docs []Schema) []error {
	in := make(chan Schema)
	out := make(chan error)

	var wg sync.WaitGroup
	wg.Add(numCPU)

	for i := 0; i < numCPU; i++ {
		go func() {
			defer wg.Done()
			for d := range in {
				if _, err := db.Insert(d); err != nil {
					out <- err
				}
			}
		}()
	}
	go func() {
		for _, d := range docs {
			in <- d
		}
		close(in)
	}()
	go func() {
		wg.Wait()
		close(out)
	}()

	errs := make([]error, 0)
	for err := range out {
		errs = append(errs, err)
	}

	return errs
}

func (db *MemDB[Schema]) Update(id string, doc Schema) (Record[Schema], error) {
	prevDoc, ok := db.docs.Get(id)
	if !ok {
		return Record[Schema]{}, fmt.Errorf("document not found")
	}

	db.docs.Put(id, doc)

	fields := getIndexFields(prevDoc)
	for _, field := range fields {
		db.deindexField(id, field)
	}

	fields = getIndexFields(doc)
	for _, field := range fields {
		db.indexField(id, field)
	}

	return Record[Schema]{Id: id, S: doc}, nil
}

func (db *MemDB[Schema]) Delete(id string) error {
	doc, ok := db.docs.Get(id)
	if !ok {
		return fmt.Errorf("document not found")
	}

	db.docs.Del(id)

	fields := getIndexFields(doc)
	for _, field := range fields {
		db.deindexField(id, field)
	}

	return nil
}

func (db *MemDB[Schema]) SearchV2(query string) []Record[Schema] {
	records := make([]Record[Schema], 0)
	infos := make([]RecordInfo, 0)
	tokens := Tokenize(query)
	for _, token := range tokens {
		recordsInfos, _ := db.index.Get(token)
		if len(infos) == 0 {
			infos = append(infos, recordsInfos...)
		} else {
			infos = Simple(infos, recordsInfos)
		}
	}
	for _, info := range infos {
		doc, _ := db.docs.Get(info.recId)
		records = append(records, Record[Schema]{Id: info.recId, S: doc})
	}

	return records
}

func (db *MemDB[Schema]) Search(query string) []Record[Schema] {
	records := make([]Record[Schema], 0)
	infos := make([]RecordInfo, 0)
	tokens := Tokenize(query)
	for _, token := range tokens {
		recordsInfos, _ := db.index.Get(token)
		for _, info := range recordsInfos {
			if idx := findRecordInfo(infos, info.recId); idx >= 0 {
				infos[idx].freq += info.freq
			} else {
				infos = append(infos, info)
			}
		}
	}
	for _, info := range infos {
		doc, _ := db.docs.Get(info.recId)
		records = append(records, Record[Schema]{Id: info.recId, S: doc})
	}

	return records
}

func (db *MemDB[Schema]) indexField(id string, text string) {
	tokens := Tokenize(text)
	tokensCount := Count(tokens)

	for token, count := range tokensCount {
		recordsInfos, _ := db.index.GetOrInsert(token, []RecordInfo{})
		recordsInfos = append(recordsInfos, RecordInfo{id, count})
		db.index.Put(token, recordsInfos)
	}
}

func (db *MemDB[Schema]) deindexField(id string, text string) {
	tokens := Tokenize(text)

	for _, token := range tokens {
		if recordsInfos, ok := db.index.Get(token); ok {
			var newRecordsInfos []RecordInfo
			for _, info := range recordsInfos {
				if !strings.EqualFold(info.recId, id) {
					newRecordsInfos = append(newRecordsInfos, info)
				}
			}
			db.index.Put(token, newRecordsInfos)
		}
	}
}

func getIndexFields(obj any) []string {
	fields := make([]string, 0)
	val := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)

	for i := 0; i < val.NumField(); i++ {
		f := t.Field(i)
		if v, ok := f.Tag.Lookup("index"); ok && strings.EqualFold(v, "true") {
			fields = append(fields, val.Field(i).String())
		}
	}

	return fields
}

func findRecordInfo(infos []RecordInfo, id string) int {
	for idx, info := range infos {
		if strings.EqualFold(info.recId, id) {
			return idx
		}
	}
	return -1
}

func Tokenize(data string) []string {
	data = punctuationRegex.ReplaceAllString(data, "")
	data = strings.ToLower(data)
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
