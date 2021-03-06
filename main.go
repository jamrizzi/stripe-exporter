package main

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"gopkg.in/gin-gonic/gin.v1"
	"io/ioutil"
)

var chargesCount = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "stripe_charges_count",
	Help: "Number of charges processed over stripe",
})

var customersCount = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "stripe_customers_count",
	Help: "Number of stripe customers registered on stripe",
})

func countCharges() {
	params := &stripe.ChargeListParams{}
	i := charge.List(params)
	fmt.Println("****************************")
	count := 0
	for i.Next() {
		count++
	}
	fmt.Println(count)
}

func webhook(c *gin.Context) {
	body := getBody(c)
	event := body["type"].(string)
	switch event {
	case "charge.succeeded":
		chargesCount.Inc()
	case "customer.created":
		customersCount.Inc()
	}
	fmt.Println("Event: " + event)
	c.JSON(200, body)
}

func init() {
	prometheus.MustRegister(chargesCount)
	prometheus.MustRegister(customersCount)
}

func main() {
	r := gin.Default()
	stripe.Key = "sk_live_8UqNzhERSTGEFfEd4r3wUfGL"
	countCharges()
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.POST("/webhook", webhook)
	r.Run()
}

func getBody(c *gin.Context) map[string]interface{} {
	var bodyMap map[string]interface{}
	bodyRaw, _ := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(bodyRaw, &bodyMap)
	return bodyMap
}
