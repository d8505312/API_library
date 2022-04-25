package lib

import (
	"fmt"
	"time"

	_ "github.com/influxdata/influxdb1-client"
	client "github.com/influxdata/influxdb1-client/v2"
)

var Influx client.Client

func InitInflux() {
	Influx = ConnInflux(
		Config.GetString("influx.ip"),
		Config.GetString("influx.port"),
	)
}

//連接 InfluxDB
func ConnInflux(host, port string) client.Client {
	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: fmt.Sprintf("http://%v:%v", host, port),
	})
	if err != nil {
		panic("failed to connect influxDB, error " + err.Error())
	}
	return cli
}

func GetInflux() client.Client {
	return Influx
}

//查詢 InfluxDB
func QueryInflux(cli client.Client, cmd string) (resp []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: Config.GetString("influx.db"),
	}
	if response, err := cli.Query(q); err == nil {
		if response.Error() != nil {
			return resp, response.Error()
		}
		resp = response.Results
	} else {
		return resp, err
	}
	return resp, nil
}

//寫入InfluxDB
func InsertInflux(cli client.Client, measurement string, tags map[string]string, fields map[string]interface{}) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  Config.GetString("influx.db"),
		Precision: Config.GetString("influx.precision"), //精度，默認ns
	})
	if err != nil {
		return err
	}

	pt, err := client.NewPoint(measurement, tags, fields, time.Now())
	if err != nil {
		return err
	}

	bp.AddPoint(pt)
	err = cli.Write(bp)
	if err != nil {
		return err
	}

	return nil
}
