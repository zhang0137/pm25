// get latest PM2.5 from http://pm25.in
// Author skyblue.
// -- I hope the sky is blue and the air is clean.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aybabtme/color"
	"github.com/shxsun/pm25/model"
	"github.com/shxsun/pm25/servd"
)

var (
	daemon = flag.Bool("daemon", false, "start as daemon")
	addr   = flag.String("addr", ":8077", "listen address(if deamoned) or dial addr")
	token  = flag.String("token", "5j1znBVAsnSf5xQyNQyq", "token required by http://pm25.in")
	dbname = flag.String("dbname", "pm25", "database name")
	dbuser = flag.String("dbuser", "root", "database username")
	dbpass = flag.String("dbpass", "toor", "database password")
)

var colorLevel = []color.Paint{
	color.GreenPaint,
	color.YellowPaint,
	color.RedPaint,
	color.PurplePaint,
	color.PurplePaint,
}

var faceLevel = []string{
	"^O^",
	"-_-",
	"-_!",
	"-_-!",
	"-_-!!",
}

func progress(tot int, cur int, paint color.Paint) string {
	brush := color.NewBrush("", paint)
	return "[" + brush(strings.Repeat("#", cur)) + strings.Repeat("-", tot-cur) + "]"
}

func cli(loc string) (err error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/%s", *addr, flag.Arg(0)))
	if err != nil {
		return
	}
	record := &model.Record{}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, record)
	if err != nil {
		return
	}
	l := record.Aqi / 100
	if l > 5 {
		l = 5
	}
	brush := color.NewBrush("", colorLevel[l])
	stars := (record.Aqi + 9) / 10
	bar := progress(50, stars, colorLevel[l])
	fmt.Printf("%-5s %s\n", brush(faceLevel[l]), bar)

	fmt.Printf("%#v\n", *record)
	return
}

func main() {
	flag.Parse()
	if *daemon {
		servd.Token = *token
		servd.DBName = *dbname
		servd.DBUser = *dbuser
		servd.DBPass = *dbpass
		err := servd.Run(*addr, time.Minute*30)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		if flag.NArg() != 1 {
			flag.Usage()
			fmt.Printf("[EXAMPLE]\n%s beijing   # will get beijing pm2.5\n", os.Args[0])
			return
		}
		if err := cli(flag.Arg(0)); err != nil {
			log.Fatal(err)
		}
	}
}
