package order

import (
	"context"
	"fmt"
	"main/internal/broker"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/segmentio/kafka-go"
)

func Create(conn *pgx.Conn) {
	go broker.KafkaCreateReader("order_paid", func(id string) {
		fmt.Println("executor got order_paid", id)
		var updateId int

		row := conn.QueryRow(
			context.Background(),
			"update orders set state = $1 where id = $2 returning id",
			"paid",
			id,
		)
		err := row.Scan(&updateId)
		if err != nil {
			fmt.Println("change order status into paid got problem")
			fmt.Println(err)
		} else {
			fmt.Println("change order status into paid success")
		}
	})

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "server healthy")
	})

	r.HandleFunc("/place-order", func(w http.ResponseWriter, r *http.Request) {
		var id int
		row := conn.QueryRow(
			context.Background(),
			"insert into orders (state) values ($1) returning id",
			"created",
		)
		err := row.Scan(&id)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = broker.KafkaWriteMessages("order_created", kafka.Message{
			Key:   []byte("random key"),
			Value: []byte(strconv.FormatInt(int64(id), 10)),
		})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "done:", id)
	})

	url := "0.0.0.0:3000"
	err := http.ListenAndServe(url, r)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Listen on %v", url)
}
