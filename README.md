# geodig
Command line tool for looking up Geolocation info for an ip address.

## Database
geodig uses the great MaxMind database.

## Build
```
$go build -o geodig geodig.go
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

## Contributions

Contributions are welcome.

## Creators

**Remco Verhoef**
- <https://twitter.com/remco_verhoef>
- <https://twitter.com/dutchcoders>

## Copyright and license

Code and documentation copyright 2011-2014 Remco Verhoef.
Code released under [the MIT license](LICENSE).
