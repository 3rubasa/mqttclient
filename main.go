package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker("tcp://10.42.0.1:1883").SetClientID("gotrivial")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetUsername("dima")
	opts.SetPassword("pass1234")
	//opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/lights/on", func(rw http.ResponseWriter, r *http.Request) {
		text := `{"state":"ON"}`
		token := c.Publish("zigbee2mqtt/0x00124b0023429241/set", 0, false, text)
		token.Wait()
		rw.Write([]byte("OK"))
	})

	mux.HandleFunc("/lights/off", func(rw http.ResponseWriter, r *http.Request) {
		text := `{"state":"OFF"}`
		token := c.Publish("zigbee2mqtt/0x00124b0023429241/set", 0, false, text)
		token.Wait()
		rw.Write([]byte("OK"))
	})

	mux.HandleFunc("/restart", func(rw http.ResponseWriter, r *http.Request) {
		syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
		rw.Write([]byte("OK"))
	})

	srv := &http.Server{
		Addr:    ":8889",
		Handler: mux,
	}

	err := srv.ListenAndServe()
	if err != nil {
		fmt.Printf("Server exited with an error: %s \n", err.Error())
	}

	c.Disconnect(250)

	time.Sleep(1 * time.Second)
}
