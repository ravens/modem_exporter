package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/maltegrosse/go-modemmanager"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const namespace = "modem"
const subsystem = ""

var (
	listenAddress = flag.String("web.listen-address", ":9898",
		"Address to listen on for telemetry")
	metricsPath = flag.String("web.telemetry-path", "/metrics",
		"Path under which to expose metrics")

	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "up"),
		"Was the last modem query successful",
		nil, nil,
	)

	roaming = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "roaming"),
		"Is the modem roaming",
		[]string{"imei", "icc", "imsi", "operatorid", "operator"}, nil,
	)

	operatorcode = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "operatorcode"),
		"Code of the operator currently used by the modem",
		[]string{"imei", "icc", "imsi", "operatorid", "operator"}, nil,
	)

	rssi = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "rssi"),
		"Level of signal reported by the modem",
		[]string{"imei", "icc", "imsi", "operatorid", "operator"}, nil,
	)

	rsrp = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "rsrp"),
		"Level of noise reported by the modem",
		[]string{"imei", "icc", "imsi", "operatorid", "operator"}, nil,
	)
)

type Exporter struct {
	mmgr modemmanager.ModemManager
}

func NewExporter(mmgr modemmanager.ModemManager) *Exporter {
	return &Exporter{
		mmgr: mmgr,
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- operatorcode
	ch <- rssi
	ch <- rsrp
	ch <- roaming
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	modems, err := e.mmgr.GetModems()
	if err != nil {
		log.Println(err)
		return
	}

	ch <- prometheus.MustNewConstMetric(
		up, prometheus.GaugeValue, 1,
	)

	for _, modem := range modems {
		sim, err := modem.GetSim()
		if err != nil {
			log.Println(err)
			continue
		}

		simIdent, err := sim.GetSimIdentifier()
		if err != nil {
			log.Println(err)
			continue
		}

		simImsi, err := sim.GetImsi()
		if err != nil {
			log.Println(err)
			continue
		}

		simOpIdent, err := sim.GetOperatorIdentifier()
		if err != nil {
			log.Println(err)
			continue
		}

		simOp, err := sim.GetOperatorName()
		if err != nil {
			log.Println(err)
			continue
		}

		modem3gpp, err := modem.Get3gpp()
		if err != nil {
			log.Println(err)
			continue
		}
		imei, err := modem3gpp.GetImei()
		if err != nil {
			log.Println(err)
			continue
		}

		regState, err := modem3gpp.GetRegistrationState()
		if err != nil {
			log.Println(err)
			continue
		}

		if regState.String() == "Roaming" {
			ch <- prometheus.MustNewConstMetric(
				roaming, prometheus.GaugeValue, 1, imei, simIdent, simImsi, simOpIdent, simOp,
			)
		} else {
			ch <- prometheus.MustNewConstMetric(
				roaming, prometheus.GaugeValue, 0, imei, simIdent, simImsi, simOpIdent, simOp,
			)
		}

		opCode, err := modem3gpp.GetOperatorCode()
		if err != nil {
			log.Println(err)
			continue
		}

		if s, err := strconv.ParseFloat(opCode, 64); err == nil {
			ch <- prometheus.MustNewConstMetric(
				operatorcode, prometheus.GaugeValue, s, imei, simIdent, simImsi, simOpIdent, simOp,
			)
		}

		modemSignal, err := modem.GetSignal()
		if err != nil {
			log.Println(err)
			continue
		}

		err = modemSignal.Setup(1)
		if err != nil {
			log.Println(err)
			continue
		}

		currentSignal, err := modemSignal.GetCurrentSignals()
		if err != nil {
			log.Println(err)
			continue
		}

		for _, sp := range currentSignal {
			ch <- prometheus.MustNewConstMetric(
				rssi, prometheus.GaugeValue, sp.Rssi, imei, simIdent, simImsi, simOpIdent, simOp,
			)

			ch <- prometheus.MustNewConstMetric(
				rsrp, prometheus.GaugeValue, sp.Rsrp, imei, simIdent, simImsi, simOpIdent, simOp,
			)
		}

		err = modemSignal.Setup(0)
		if err != nil {
			log.Println(err)
			continue
		}

	}

}

func main() {

	flag.Parse()

	mmgr, err := modemmanager.NewModemManager()
	if err != nil {
		log.Fatal(err.Error())
	}
	version, err := mmgr.GetVersion()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = mmgr.SetLogging(modemmanager.MMLoggingLevelError)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("Starting modem exporter using ModemManager v%s", version)

	exporter := NewExporter(mmgr)
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Modem Exporter</title></head>
             <body>
             <h1>Modem Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
