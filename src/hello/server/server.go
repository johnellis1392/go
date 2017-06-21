package hello

import (
	"fmt"
	"net/http"
)

func concurrently(integers []int) []int {
	ch := make(chan int)
	responses := []int{}
	for _, i := range integers {
		go func(i int) {
			ch <- i * i
		}(i)
	}

	for {
		result := <-ch
		responses = append(responses, result)
		fmt.Print(responses)
		fmt.Print("\n")
		if len(responses) == len(integers) {
			return responses
		}
	}
}

func handler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprint(response, "Hello, World!")
}

func main() {
	fmt.Print("Hello, World!\n")
	concurrently([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
