[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=mildis_opensmtpd-go-filter-prepend&metric=alert_status)](https://sonarcloud.io/dashboard?id=mildis_opensmtpd-go-filter-prepend)

# opensmtpd-go-filter-prepend
OpenSMTPD filter which prepends a string on email's subject if not present.  
Works with OpenSMTPD 6.6 and 6.7.

## Usage
* build the filter for your target platform  
`env GOOS=openbsd GOARCH=amd64 go build filter-prepend.go`
* make OpenSMTPD use the filter
```
filter prepend proc-exec "filter-prepend-go --prefix='[*EXT*]' --extraprefix='[EXT]'"
listen on em0 tls pki "*" filter { senderscore, rspamd, prepend }
```
* default prefix is `[*EXT*]` if not specified on the CLI  
* default extraprefix is `[EXT]` if not specified on the CLI  
Be warned that a too generic prefix could match legitimate wording and thus it won't be added to the Subject
* option `--encode` forces prefix encoding whether the subject is encoded or not.

## Known limitations
1. if Subject is too long and spans between multiple datalines, only the begining will be analyzed for prefix token
2. RFC2047 Subject won't be decoded if 1/ happens inside an encoded string
