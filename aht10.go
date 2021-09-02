package main

import (
	"log"
	"time"
	"errors"

  "github.com/d2r2/go-i2c"
  "github.com/prometheus/client_golang/prometheus"
)

const (
	aht10_wait = 200 * time.Millisecond
	aht10_init = 500 * time.Millisecond
	sensor_aht10_name = "AHT10"
)

type AHT10Exporter struct {
	metricsTemperature     *prometheus.Desc
	metricsHumidity        *prometheus.Desc
	metricsStatus          *prometheus.Desc
	id                      int
	addr                    uint8
	i2c                     *i2c.I2C
}

var _ prometheus.Collector = &AHT10Exporter{}

// null init - could have option to set
func NewAHT10Exporter(id int, addr uint8) *AHT10Exporter {
	e := &AHT10Exporter{id: id, addr: addr}
	e.init()
	return e
}

func (e *AHT10Exporter) init() {
	i2c, err := i2c.NewI2C(e.addr, e.id)
  if err != nil { log.Fatal(err) }
	e.i2c = i2c
	_, err = e.i2c.WriteBytes([]byte{0xE1, 0x28, 0x00})
	time.Sleep(aht10_init)
	if err != nil { log.Fatal(err) }
}

func (e *AHT10Exporter) read(cond bool) []float32 {
  ret := []float32{0.0, 0.0, 0.0}

	val := []byte{}
	err := errors.New("")
	if (cond) {
    valb := byte(0)
    _, err = e.i2c.ReadRegU8(valb)
    if err != nil { log.Fatal(err) }
    _, err = e.i2c.WriteBytes([]byte{0xAC, 0x33, 0x00})
    if err != nil { log.Fatal(err) }
		time.Sleep(aht10_wait)
	}
  val, _, err = e.i2c.ReadRegBytes(0, 8)

  ret[0] = (float32(val[3] & 0x0F) * 256 + float32(val[4])) * 256 + float32(val[5])
  ret[0] = (ret[0] * 200 / 1048576) - 50
  ret[1] = (float32(val[1]) * 256 + float32(val[2])) * 16 + (float32(val[3] & 0xF0) / 16)
  ret[1] = ret[1] * 100 / 1048576
	ret[2] = float32(val[0])
  return ret
}

func (e *AHT10Exporter) Describe(ch chan<- *prometheus.Desc) {
	labels := []string{"AHT10"}
	e.metricsTemperature = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, sensor_aht10_name, "temperature"),
		"Temperature on sensor",
		labels,
		nil,
	)
	ch <- e.metricsTemperature
	e.metricsHumidity = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, sensor_aht10_name, "humidity"),
		"Humidity on sensor",
		labels,
		nil,
	)
	ch <- e.metricsHumidity
	e.metricsStatus = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, sensor_aht10_name, "status"),
		"Status on sensor",
		labels,
		nil,
	)
	ch <- e.metricsStatus
}

func (e *AHT10Exporter) Collect(ch chan<- prometheus.Metric) {
	val := e.read(true)
	if (val[1] == 0.0) { val = e.read(false) }
	ch <- prometheus.MustNewConstMetric(
		e.metricsTemperature,
		prometheus.GaugeValue,
		float64(val[0]),
		"A",
	)
	if (val[1] != 0.0) {
  	ch <- prometheus.MustNewConstMetric(
  		e.metricsHumidity,
	  	prometheus.GaugeValue,
		  float64(val[1]),
  		"A",
	  )
	}
	ch <- prometheus.MustNewConstMetric(
		e.metricsStatus,
		prometheus.GaugeValue,
		float64(val[2]),
		"A",
	)	
}

