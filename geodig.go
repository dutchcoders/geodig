package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"github.com/oschwald/maxminddb-golang"
	"io"
	"net"
	"net/http"
	"os"
	"os/user"
	"path"
	"strings"
)

var reader *maxminddb.Reader
var verbose bool
var format string

type geoIPResult struct {
	Location struct {
		Longitude float64 `maxminddb:"longitude"`
		Latitude  float64 `maxminddb:"latitude"`
	} `maxminddb:"location"`
	City struct {
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"city"`
	Country struct {
		IsoCode string            `maxminddb:"iso_code"`
		Names   map[string]string `maxminddb:"names"`
	} `maxminddb:"country"`
}

// create ~/.geodig folder, and download geolite db
func init() {
	var err error
	var usr *user.User
	if usr, err = user.Current(); err != nil {
		panic(err)
	}

	geodig_path := path.Join(usr.HomeDir, ".geodig")

	_, err = os.Stat(geodig_path)

	switch {
	case err == nil:
		break
	case os.IsNotExist(err):
		if err = os.Mkdir(geodig_path, 0755); err != nil {
			panic(err)
		}
	default:
		panic(err)

	}

	db_path := path.Join(geodig_path, "GeoLite2-City.mmdb")

	_, err = os.Stat(db_path)
	if os.IsNotExist(err) {
		if verbose {
			fmt.Println("Downloading database")
		}

		if err = download(geoLiteURL, db_path); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if verbose {
		fmt.Printf("Using database %s.\n", db_path)
	}

	reader, err = maxminddb.Open(db_path)
}

func help() {
	fmt.Println("No ip addresses")
}

func download(url string, dest string) error {
	client := &http.Client{}

	req, err := http.NewRequest("GET", geoLiteURL, nil)
	if err != nil {
		return err
	}

	var resp *http.Response
	if resp, err = client.Do(req); err != nil {
		return err
	}

	defer resp.Body.Close()

	gzf, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer gzf.Close()

	f, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = io.Copy(f, gzf)
	if err != nil {
		return err
	}

	return nil
}

const geoLiteURL = "http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz"

func main() {
	// initial download
	// update
	// open from cache
	// resolve

	flag.StringVar(&format, "format", "(country) ((city))", "format")
	flag.BoolVar(&verbose, "verbose", false, "verbose")
	flag.Parse()

	format = strings.Replace(format, "\\n", "\n", -1)

	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	var args = flag.Args()

	if fi.Mode()&os.ModeNamedPipe > 0 {
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			args = append(args, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}

	if len(args) == 0 {
		help()
		os.Exit(1)
	}

	for _, arg := range args {
		var addr net.IP
		addr = net.ParseIP(arg)
		if addr == nil {
			fmt.Printf("%s is not a valid ip address", addr)
			continue
		}

		var result geoIPResult

		var err error
		if err = reader.Lookup(addr, &result); err != nil {
			panic(err)
		}

		var p string
		p = format
		p = strings.Replace(p, "(ip)", addr.String(), -1)
		p = strings.Replace(p, "(country)", result.Country.Names["en"], -1)
		p = strings.Replace(p, "(city)", result.City.Names["en"], -1)
		p = strings.Replace(p, "(lat)", fmt.Sprintf("%f", result.Location.Latitude), -1)
		p = strings.Replace(p, "(long)", fmt.Sprintf("%f", result.Location.Longitude), -1)

		fmt.Print(p)
	}
}
