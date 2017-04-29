package main

import (
	"fmt"
	"encoding/json"
	"github.com/labstack/echo"
	"net/http"
	"bytes"
	"github.com/op/go-logging"
	"os"
	"gopkg.in/alecthomas/kingpin.v2"
)

var local = kingpin.Arg("local", "Running locally?").Bool()

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{pid} %{longfunc} ▶ %{level:.4s} %{id:03x} %{message}`,
)
var logg = logging.MustGetLogger("example")

var serviceLocations = map[string]string{"sum":"n2","prod":"n1"}
var FuncMap  = map[string]func([]float64)error{"sum":SendSumjob,"prod":SendProdjob}

type Job struct{
	Data []float64 `json:"data"`
	Service string `json:"service"`
}

type Result struct{
	Out float64 `json:"out"`
}

func main(){

	kingpin.Parse()

	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend1Formatter := logging.NewBackendFormatter(backend1, format)
	logging.SetBackend(backend1, backend1Formatter)

	Listen()
}

func Recivejob(c echo.Context) error {
	logg.Info("Recieved Job: public")
	j := Job{}
	err := c.Bind(&j)
	if err != nil {
		logg.Error(err)
		return err
	}
	logg.Info("Sending job: public")
	service := FuncMap[j.Service]
	service(j.Data)
	return nil
}

func Reciveres(c echo.Context) error {
	logg.Info("Recieved result: public")
	r := Result{}
	err := c.Bind(&r)
	if err != nil {
		logg.Error(err)
		return err
	}
	logg.Info("Result: ",r.Out)
	return nil
}


func SendSumjob(d []float64) error {
	j := Job{Data:d}
	data, err := json.Marshal(j)
	if err != nil {
		logg.Error(err)
		return err
	}
	sumhost := "8000"
	_, err = http.Post(fmt.Sprintf("http://localhost:%s/",sumhost),"application/json",bytes.NewBuffer(data))
	if err != nil {
		logg.Error(err)
		return err
	}
	return nil
}

func SendProdjob(d []float64) error {
	j := Job{Data:d}
	data, err := json.Marshal(j)
	if err != nil {
		logg.Error(err)
		return err
	}
	prodhost := "9000"
	_, err = http.Post(fmt.Sprintf("http://localhost:%s/",prodhost),"application/json",bytes.NewBuffer(data))
	if err != nil {
		logg.Error(err)
		return err
	}
	return nil
}

func Listen() {
	logg.Info("Starting to Listen: public")
	e := echo.New()

	e.POST("/", Recivejob)
	e.POST("result/", Reciveres)
	if *local {
		e.Start(":7000")
	} else {
		e.Start(":8000")
	}
}
