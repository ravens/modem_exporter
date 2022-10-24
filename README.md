# modem_exporter
Prometheus Exporter for modemmanager metrics. This is a Fork from [modem_exporter](https://github.com/ravens/modem_exporter) project. 

## architecture

This exporter is basically piggy-backing on ModemManager and some [golang bindings](github.com/maltegrosse/go-modemmanager) to build a minimal prometheus exporter to expose the state of an LTE modem as prometheus metrics.

## port and endpoint

Default port is 9898 and endpoint is http://EXPORTERENDPOINT:9898/metrics.

## Metrics exported

Right now are exposed:

``` bash
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
# HELP modem_connected Is the modem connected
# TYPE modem_connected gauge
modem_connected{icc="0000000000000000000",imei="00000000000000",imsi="00000000000000",operator="foobar",operatorid="00101",rat="lte",v_operator="VisitedNetwork"} 1
# HELP modem_operatorcode Code of the operator currently used by the modem
# TYPE modem_operatorcode gauge
modem_operatorcode{icc="0000000000000000000",imei="00000000000000",imsi="00000000000000",operator="foobar",operatorid="00101",rat="lte",v_operator="VisitedNetwork"} 00101
# HELP modem_lac LAC currently used by the modem
# TYPE modem_lac gauge
modem_lac{icc="0000000000000000000",imei="00000000000000",imsi="00000000000000",operator="foobar",operatorid="00101",rat="lte",v_operator="VisitedNetwork"} 65534
# TYPE modem_cellid gauge
modem_cellid{icc="0000000000000000000",imei="00000000000000",imsi="00000000000000",operator="foobar",operatorid="00101",rat="lte",v_operator="VisitedNetwork"} 1.28013588e+08

#GPS Metrics
# HELP modem_lat Latitude
# TYPE modem_lat gauge
modem_lat{icc="0000000000000000000",imei="00000000000000",imsi="00000000000000",operator="foobar",operatorid="00101",rat="lte",v_operator="VisitedNetwork"} 40.6219309 

# HELP modem_lon Longitude 
# TYPE modem_lon gauge
modem_lon{icc="0000000000000000000",imei="00000000000000",imsi="00000000000000",operator="foobar",operatorid="00101",rat="lte",v_operator="VisitedNetwork"} -5.596668

# HELP modem_alt Altitude 
# TYPE modem_alt gauge
modem_alt{icc="0000000000000000000",imei="00000000000000",imsi="00000000000000",operator="foobar",operatorid="00101",rat="lte",v_operator="VisitedNetwork"} 710.9


# Signal Metrics
# HELP modem_rsrq Reference Signal Received Quality
# TYPE modem_rsrq gauge
modem_rsrq{icc="0000000000000000000",imei="00000000000000",imsi="00000000000000",operator="foobar",operatorid="00101",rat="lte",v_operator="VisitedNetwork"} -12

# HELP modem_rssi Level of signal reported by the modem
# TYPE modem_rssi gauge
modem_rssi{icc="0000000000000000000",imei="00000000000000",imsi="00000000000000",operator="foobar",operatorid="00101",rat="lte",v_operator="VisitedNetwork"} -100

# HELP modem_sinr Signal-to-interference-plus-noise ratio
# TYPE modem_sinr gauge
modem_sinr{icc="0000000000000000000",imei="00000000000000",imsi="00000000000000",operator="foobar",operatorid="00101",rat="lte",v_operator="VisitedNetwork"} 0

# HELP modem_snr The LTE S/R ratio
# TYPE modem_snr gauge
modem_snr{icc="0000000000000000000",imei="00000000000000",imsi="00000000000000",operator="foobar",operatorid="00101",rat="lte",v_operator="VisitedNetwork"} 15.4

```

Only send GPS metrics if the GPS NMEA [TypeGGA](https://orolia.com/manuals/VSP/Content/NC_and_SS/Com/Topics/APPENDIX/NMEA_GGAmess.htm) is received.

## build

``` bash
docker-compose build
docker-compose up 
```
