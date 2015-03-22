# geodig
Command line tool for looking up Geolocation info for an ip address.

## Database
This product includes GeoLite data created by MaxMind, available from [http://www.maxmind.com/](http://www.maxmind.com).

## Demo

![](demo.gif)

## Build
```
$ go build -o geodig geodig.go
```

## Install using Homebrew

```
$ brew tap dutchcoders/homebrew-geodig
$ brew install geodig
```

## Examples

Get location for ip address
```
$geodig 192.30.252.131
United States (San Francisco)%
```

Ip addresses can be piped, for use with log files
```
$echo 192.30.252.131|geodig
United States (San Francisco)%
```

Analyzing log files
```
$curl http://#####.###/logs/access.log | awk '{print $1}' | sort | uniq | go run geodig.go --format "(country)\n"| sort | uniq
Afghanistan
Australia
Belarus
Bulgaria
Canada
China
Finland
France
Germany
India
Indonesia
Ireland
Israel
Netherlands
Poland
Romania
Russia
Rwanda
Spain
Thailand
Ukraine
United Kingdom
United States
```

Creating a shell alias 
```
$alias geodig='go run geodig.go --format "(country)\n"'
```

## Contributions

Contributions are welcome.

## Creators

**Remco Verhoef**
- <https://twitter.com/remco_verhoef>
- <https://twitter.com/dutchcoders>

## Copyright and license

Code and documentation copyright 2011-2014 Remco Verhoef.

Code released under [the MIT license](LICENSE).
