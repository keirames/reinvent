package bomber

import (
	"fmt"
	"io"
	"net/http"
)

func Run() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, "hello")
		if err != nil {
			fmt.Println(err)
		}
	})

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Bomber crashed!!!", err)
		return
	}
}
