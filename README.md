# opensmtpd-go-filter-prepend
OpenSMTPD filter which prepends a string on email's subject if not present

## Usage
* build the filter for your target platform  
`env GOOS=openbsd GOARCH=amd64 go build filter-prepend.go`
* make OpenSMTPD use the filter
```
filter prepend proc-exec "filter-prepend-go --prefix='[*EXT*]'"
listen on em0 tls pki "*" filter { senderscore, rspamd, prepend }
```
* default prefix is `[*EXT*]` if not specified on the CLI  
Be warned that a too generic prefix could match legitimate wording and thus it won't be added to the Subject

## Known limitations
1. if Subject is too long and spans between multiple datalines, only the begining will be analyzed for prefix token
2. RFC2047 Subject won't be decoded if 1/ happens inside an encoded string
3. for RFC2047 Subject, the rewrote Subject line will be a mix of the prefix in clear ASCII and the encoded original subject
