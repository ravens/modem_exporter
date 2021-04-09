# modem_exporter
Prometheus Exporter for modemmanager metrics. 

## architecture

This exporter is basically piggy-backing on ModemManager and some [golang bindings](github.com/maltegrosse/go-modemmanager) to build a minimal prometheus exporter to expose the state of an LTE modem as prometheus metrics.

## port and endpoint

Default port is 9898 and endpoint is http://EXPORTERENDPOINT:9898/metrics.

## Metrics exported

Right now are exposed:
```
# HELP modem_up Was the last modem query successful
# TYPE modem_up gauge
modem_up 1
# HELP modem_tac TAC currently used by the modem
# TYPE modem_tac gauge
modem_tac{icc="0000000000000000000",imei="00000000000000",imsi="00000000000000",operator="foobar",operatorid="00101",rat="lte",v_operator="VisitedNetwork"} 50620
# HELP modem_roaming Is the modem roaming
# TYPE modem_roaming gauge
modem_roaming{icc="0000000000000000000",imei="00000000000000",imsi="00000000000000",operator="foobar",operatorid="00101",rat="lte",v_operator="VisitedNetwork"} 1
# HELP modem_registered Is the modem registered
# TYPE modem_registered gauge
modem_registered{icc="0000000000000000000",imei="00000000000000",imsi="00000000000000",operator="foobar",operatorid="00101",rat="lte",v_operator="VisitedNetwork"} 1
# HELP modem_operatorcode Code of the operator currently used by the modem
# TYPE modem_operatorcode gauge
modem_operatorcode{icc="0000000000000000000",imei="00000000000000",imsi="00000000000000",operator="foobar",operatorid="00101",rat="lte",v_operator="VisitedNetwork"} 00101
# HELP modem_lac LAC currently used by the modem
# TYPE modem_lac gauge
modem_lac{icc="0000000000000000000",imei="00000000000000",imsi="00000000000000",operator="foobar",operatorid="00101",rat="lte",v_operator="VisitedNetwork"} 65534
# TYPE modem_cellid gauge
modem_cellid{icc="0000000000000000000",imei="00000000000000",imsi="00000000000000",operator="foobar",operatorid="00101",rat="lte",v_operator="VisitedNetwork"} 1.28013588e+08
```

Right now CellID, TAC and LAC are exposed as Gauge. Perhaps it should be labels. I am not really sure.

## Hardware BOM

Using Amazon:
- M2 to USB2 converter https://www.amazon.es/gp/product/B06ZYXH76M - 32 euros
- Dell M2 4G LT X7 Qualcomm Modem https://www.amazon.es/gp/product/B07QMQRZCN - 18 euros

A linux machine with Ubuntu 20.04 to host it should be enough. 

## build

```
docker-compose build
docker-compose up 
```

