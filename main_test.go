package main

import ("testing"
"fmt")

func TestMain(t *testing.T){
    searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		fmt.Println(err)
        
	}
   results := searcher.Search("Hamlet")
   fmt.Println(results) 

}