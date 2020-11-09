package main

import (
	"fmt"
	"net/http"
)

const maxSize = 100 << 20 // 100M
const bufSize = 1 << 20   // 1M

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/upload", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			http.Error(writer, "use POST to upload files.", http.StatusMethodNotAllowed)
			return
		}
		if request.ContentLength > maxSize {
			http.Error(writer, "file too large. 100M limit.", http.StatusExpectationFailed)
			return
		}
		reader := http.MaxBytesReader(writer, request.Body, maxSize)
		buf := make([]byte, bufSize)
		bytesRead := int64(0)

		for {
			n, err := reader.Read(buf)
			fmt.Printf("Read %d bytes\n", n)

			if n > 0 {
				bytesRead += int64(n)
				percentFinished := (float64(bytesRead) / float64(request.ContentLength)) * 100.0
				nn, err := writer.Write([]byte(fmt.Sprintf("Read %d bytes (%.1f%%)\n", n, percentFinished)))
				if f, ok := writer.(http.Flusher); ok {
					f.Flush()
				}
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
