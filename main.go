package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type MessageResponse struct {
	Response     string
	Url          string
	TimeResponse time.Duration
}

func requestApi(url string, data chan MessageResponse) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	timeReqInit := time.Now()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	timeReqResp := time.Since(timeReqInit)
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	data <- MessageResponse{
		Response:     string(b),
		Url:          url,
		TimeResponse: timeReqResp,
	}
}

func main() {
	chanReqViaCep := make(chan MessageResponse)
	chanReqBrasilApi := make(chan MessageResponse)

	go requestApi("http://viacep.com.br/ws/04921080/json/", chanReqViaCep)
	go requestApi("https://brasilapi.com.br/api/cep/v1/04921080", chanReqBrasilApi)

	select {
	case reqBrasilApi := <-chanReqBrasilApi:
		slog.Info(fmt.Sprintf("%+v", reqBrasilApi))

	case reqViaCep := <-chanReqViaCep:
		slog.Info(fmt.Sprintf("%+v", reqViaCep))

	case <-time.After(time.Second * 1):
		slog.Error("timeout")
	}
}
