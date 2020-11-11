# ldaptest
This a tool to test connectivity of a remote LDAP server. It can also retrieve records for one ldap user

## Install
Setup GOPATH env variable first.
```
$ git clone https://github.com/scottgai/ldaptest.git
$ cd ldaptest
$ go get
$ go build
```

## Run

The tool could read arguments either from command line or environment variables.
If LDAP port is set to 636 (hard coded) a secure connection over ldaps will be initiated).

```
$./ldaptest --help
Usage of ./ldaptest:
  -basedn string
        LDAP Base DN to search
  -cacert string
        CA cert file used to authenticate LDAP server
  -env
        Read parameters from environment variables instead of commaond line arguments. Enviroment variables:
            LDAP_HOST - LDAP host name
            LDAP_PORT - LDAP port to connect to
            LDAP_BASEDN - LDAP Base DN to search
            LDAP_SVC_USER - username to authenticate LDAP server
            LDAP_SVC_PASS - password to authenticate LDAP server. Use '\$' if '$'  is present
            LDAP_CACERT - CA cert file used to authenticate LDAP server
            LDAP_PROTO - protocol used to connect to LDAP server, tcp/udp
            SKIP_SSL_VALIDATION - skip SSL validation
            LDAP_QUERY_USER - LDAP user to be queried
  -help
        Print usage messages
  -host string
        LDAP host name
  -ldapuser string
        LDAP user to be queried
  -pass string
        password to authenticate LDAP server. Use '\$' if '$'  is present
  -port int
        LDAP port to connect to (default 389)
  -protocol string
        protocol used to connect to LDAP server, tcp/udp (default "tcp")
  -skip_ssl_validation
        Skip SSL validation
  -timeout duration
        Timeout value to dial LDAP host (default 1m0s)
  -user string
        username to authenticate LDAP server
```

## Examples
- Reading arguments from command line
```
$ ./ldaptest --host ldap.example.com --port 636 --basedn "ou=users,dc=example,dc=com" --timeout 120s --user "cn=bind,ou=Users,dc=example,dc=com" --pass "abcd\$1234" --ldapuser "James Smith" -cacert ./ldap-ca.cert 

Invalid LDAP Timeout, must be within 0-3600, set it to default 120 seconds
Invalid LDAP protocol tcp, set it to 'tcp'
-----------------------------------------
ldap-host: ldap.example.com
ldap-port: 636
ldap-cacert: ./ldap-ca.cert
ldap-user: cn=bind,ou=Users,dc=example,dc=com
ldap-passwd: abcd$1234
ldap-basedn: ou=users,dc=example,dc=com
ldap-proto: tcp
ldap-timeout: 120 seconds
skip-ssl-validation: false
LDAP user to query:  James Smith
-----------------------------------------
Loaded CA certificate from file: ./ldap-ca.cert
DIAL successfully
BIND successfully
Found 1 entries for user Sandra Walls
    DN: cn=James Smith,ou=Users,dc=example,dc=com
      givenName: [James]
      sn: [Smith]
      cn: [James Smith]
      uid: [swalls]
      uidNumber: [2001]
      gidNumber: [101]
      homeDirectory: [/home/users/jsmith]
      loginShell: [/bin/bash]
      objectClass: [inetOrgPerson posixAccount top User]
```

- Reading argumets from environment variables
```
$ export LDAP_HOST="ldap.example.com" LDAP_PORT=636 LDAP_BASED_DN="ou=users,dc=example,dc=com" LDAP_PROTO="tcp" LDAP_TIMEOUT=180 LDAP_SVC_USER="cn=bind,ou=Users,dc=example,dc=com" LDAP_SVC_PASS="abcd\$1234" SKIP_SSL_VALIDATION=true LDAP_CACERT=./lab-ca.cert LDAP_QUERY_USER="James Smith"

$./ldaptest --env
Loading parameters from envrionment variables
Invalid LDAP protocol tcp, set it to 'tcp'
-----------------------------------------
ldap-host: ldap.example.com
ldap-port: 636
ldap-cacert: ./lab-ca.cert
ldap-user: cn=bind,ou=Users,dc=example,dc=com
ldap-passwd: abcde$1234
ldap-basedn: ou=users,dc=example,dc=com
ldap-proto: tcp
ldap-timeout: 180 seconds
skip-ssl-validation: true
LDAP user to query:  James Smith
-----------------------------------------
Loaded CA certificate from file: ./lab-ca.cert
DIAL successfully
BIND successfully
Found 1 entries for user Sandra Walls
    DN: cn=James Smith,ou=Users,dc=example,dc=com
      givenName: [James]
      sn: [Smith]
      cn: [James Smith]
      uid: [swalls]
      uidNumber: [2001]
      gidNumber: [101]
      homeDirectory: [/home/users/jsmith]
      loginShell: [/bin/bash]
      objectClass: [inetOrgPerson posixAccount top User]
```
