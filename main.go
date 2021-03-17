package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
)

func main() {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", handleSearch(searcher))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("Listening on port %s...", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

type Searcher struct {
	CompleteWorks string
	SuffixArray   *suffixarray.Index
}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()["q"]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}
		results := searcher.Search(query[0])
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err := enc.Encode(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("encoding failure"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	}
}

func (s *Searcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.CompleteWorks = string(dat)
	s.SuffixArray = suffixarray.New(dat)
	return nil
}

func (s *Searcher) Search(query string) []string {
	idxs := s.SuffixArray.Lookup([]byte(query), -1)
	
	var sortidxs sort.IntSlice
	sortidxs = idxs
	sortidxs.Sort()
	
	results := []string{}
	if sortidxs == nil{
		norows := fmt.Sprintf("%s is not present. The word searched is case sensitive.",query) 
		results = append(results,norows)
	}else{ 
		for _, idx := range sortidxs {
			
			sentence := s.CompleteWorks[idx-251:idx+251]
			formatted := getformatedsentence(sentence)
			results = append(results, formatted)
		}
	}
	return results
}


func getformatedsentence(sentence string)string{
	 const space = ' '
	 const newline = '\n'
	 
	 if sentence[0] != space{
		for i:=1;i<len(sentence);i++{
			if sentence[i] == space || sentence[i] == newline{
				sentence = sentence[i+1:]
				break
			}
		} 
	 }
	 for i:=len(sentence)-1;i>=0;i--{
		if sentence[i] == space || sentence[i] == newline{
			sentence = sentence[:i]
			break
		}
	} 
	
	return sentence
	 
}
