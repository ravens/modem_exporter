package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/maltegrosse/go-modemmanager"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const namespace = "modem"
const subsystem = ""

var (
	modemlabels = []string{"imei", "icc", "imsi", "operatorid", "operator", "v_operator", "rat"}

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
		modemlabels, nil,
	)

	operatorcode = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "operatorcode"),
		"Code of the operator currently used by the modem",
		modemlabels, nil,
	)

	rssi = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "rssi"),
		"Level of signal reported by the modem",
		modemlabels, nil,
	)

	rsrp = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "rsrp"),
		"Level of noise reported by the modem",
		modemlabels, nil,
	)

	registered = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "registered"),
		"Is the modem registered",
		modemlabels, nil,
	)

	connected = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "connected"),
		"Is the modem connected",
		modemlabels, nil,
	)

	cellid = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "cellid"),
		"CellID currently used by the modem",
		modemlabels, nil,
	)

	tac = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "tac"),
		"TAC currently used by the modem",
		modemlabels, nil,
	)

	lac = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "lac"),
		"LAC currently used by the modem",
		modemlabels, nil,
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

		state, err := modem.GetState()
		if err != nil {
			log.Println("cannot get modem state:" + err.Error())
			continue
		}

		if state.String() == "Disabled" {
			log.Println("modem disabled, trying to enable it")
			err = modem.Enable()
			if err != nil {
				log.Println(err)
				continue
			}
		}

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

		opName, err := modem3gpp.GetOperatorName()
		if err != nil {
			log.Println(err)
			continue
		}

		ratList, err := modem.GetAccessTechnologies()
		if err != nil {
			log.Println(err)
			continue
		}

		if len(ratList) != 1 {
			log.Println(err)
			continue
		}

		rat := strings.ToLower(ratList[0].String())

		state, err = modem.GetState()
		if err != nil {
			log.Println("cannot get modem state:" + err.Error())
			continue
		}

		// if we are registered, we should try to connect
		if state.String() == "Registered" {

			apn := os.Getenv("MODEM_EXPORTER_APN")

			if apn != "" {

				bearers, _ := modem.GetBearers()

				// delete all bearer - if registered but no bearer something is likely wrong
				for _, bearer := range bearers {
					bearer.Disconnect()
					err = modem.DeleteBearer(bearer)
					if err != nil {
						log.Println(err)
						continue
					}
				}

				modemSimple, err := modem.GetSimpleModem()
				if err != nil {
					log.Println(err)
				} else {
					property := modemmanager.SimpleProperties{Apn: apn}
					newBearer, err := modemSimple.Connect(property)
					if err != nil {
						log.Println(err)
					} else {
						fmt.Println("New Bearer: ", newBearer)
					}
				}

			}

		}

		if state.String() == "Registered" || state.String() == "Connected" {
			ch <- prometheus.MustNewConstMetric(
				registered, prometheus.GaugeValue, 1, imei, simIdent, simImsi, simOpIdent, simOp, opName, rat,
			)
		} else {
			ch <- prometheus.MustNewConstMetric(
				registered, prometheus.GaugeValue, 0, imei, simIdent, simImsi, simOpIdent, simOp, opName, rat,
			)
		}

		if state.String() == "Connected" {
			ch <- prometheus.MustNewConstMetric(
				connected, prometheus.GaugeValue, 1, imei, simIdent, simImsi, simOpIdent, simOp, opName, rat,
			)
		} else {
			ch <- prometheus.MustNewConstMetric(
				connected, prometheus.GaugeValue, 0, imei, simIdent, simImsi, simOpIdent, simOp, opName, rat,
			)
		}

		// failure reason
		_, err = modem.GetStateFailedReason()
		if err != nil {
			log.Println(err)
			continue
		}

		modemLocation, err := modem.GetLocation()
		if err != nil {
			log.Println(err)
			continue
		}

		mloc, err := modemLocation.GetCurrentLocation()
		if err != nil {
			log.Println(err)
			continue
		}

		cellID := mloc.ThreeGppLacCi.Ci

		if decCellID, err := strconv.ParseInt(cellID, 16, 64); err == nil {
			ch <- prometheus.MustNewConstMetric(
				cellid, prometheus.GaugeValue, float64(decCellID), imei, simIdent, simImsi, simOpIdent, simOp, opName, rat,
			)
		} else {
			log.Println(err)
		}

		lAC := mloc.ThreeGppLacCi.Lac
		if decLAC, err := strconv.ParseInt(lAC, 16, 64); err == nil {
			ch <- prometheus.MustNewConstMetric(
				lac, prometheus.GaugeValue, float64(decLAC), imei, simIdent, simImsi, simOpIdent, simOp, opName, rat,
			)
		} else {
			log.Println(err)
		}

		tAC := mloc.ThreeGppLacCi.Tac
		if decTAC, err := strconv.ParseInt(tAC, 16, 64); err == nil {
			ch <- prometheus.MustNewConstMetric(
				tac, prometheus.GaugeValue, float64(decTAC), imei, simIdent, simImsi, simOpIdent, simOp, opName, rat,
			)
		} else {
			log.Println(err)
		}

		regState, err := modem3gpp.GetRegistrationState()
		if err != nil {
			log.Println(err)
			continue
		}

		if regState.String() == "Roaming" {
			ch <- prometheus.MustNewConstMetric(
				roaming, prometheus.GaugeValue, 1, imei, simIdent, simImsi, simOpIdent, simOp, opName, rat,
			)
		} else {
			ch <- prometheus.MustNewConstMetric(
				roaming, prometheus.GaugeValue, 0, imei, simIdent, simImsi, simOpIdent, simOp, opName, rat,
			)
		}

		opCode, err := modem3gpp.GetOperatorCode()
		if err != nil {
			log.Println(err)
			continue
		}

		if s, err := strconv.ParseFloat(opCode, 64); err == nil {
			ch <- prometheus.MustNewConstMetric(
				operatorcode, prometheus.GaugeValue, s, imei, simIdent, simImsi, simOpIdent, simOp, opName, rat,
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

		time.Sleep(2 * time.Second)

		currentSignal, err := modemSignal.GetCurrentSignals()
		if err != nil {
			log.Println(err)
			continue
		}

		for _, sp := range currentSignal {
			ch <- prometheus.MustNewConstMetric(
				rssi, prometheus.GaugeValue, sp.Rssi, imei, simIdent, simImsi, simOpIdent, simOp, opName, rat,
			)

			ch <- prometheus.MustNewConstMetric(
				rsrp, prometheus.GaugeValue, sp.Rsrp, imei, simIdent, simImsi, simOpIdent, simOp, opName, rat,
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
