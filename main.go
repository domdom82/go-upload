package main

import (
	"fmt"
	"net/http"
)

const maxSize = 10 << 20  // 10M
const bufSize = 100 << 10 // 100K

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/upload", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			http.Error(writer, "use POST to upload files.", http.StatusMethodNotAllowed)
			return
		}
		if request.ContentLength > maxSize {
			http.Error(writer, "file too large. 10M limit.", http.StatusExpectationFailed)
			return
		}
		reader := http.MaxBytesReader(writer, request.Body, maxSize)
		buf := make([]byte, bufSize)

		for {
			n, err := reader.Read(buf)
			fmt.Printf("Read %d bytes\n", n)

			if n > 0 {
				nn, err := writer.Write([]byte(fmt.Sprintf("Read %d bytes\n", n)))
				fmt.Printf("Wrote %d bytes\n", nn)

				if err != nil {
					fmt.Println("ERROR (write): ", err)
					break
				}
			}

			if err != nil {
				break
			}
		}
	})

	_ = http.ListenAndServe(":8080", mux)
}
