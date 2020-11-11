package main

import (
    "crypto/tls"
	"crypto/x509"
	"time"
	"fmt"
	"flag"
	"strings"
	"strconv"
	"os"
	"io/ioutil"
	"gopkg.in/ldap.v2"
)

var (
	ldapSvcUser  string
	ldapSvcPass  string
	ldapHost string
	ldapPort int
	ldapProto string
	ldapBaseDn string
	ldapCACert string
	ldapTimeout time.Duration
	skipSslValdation bool
	env bool
	userName string
	help bool
)

func main() {
	var roots *x509.CertPool = nil
	var l *ldap.Conn
	var err error
	var tmpstr string
	var i int

	flag.StringVar(&ldapHost, "host", "", "LDAP host name")
	flag.StringVar(&ldapBaseDn, "basedn", "", "LDAP Base DN to search")
	flag.IntVar(&ldapPort, "port", 389, "LDAP port to connect to")
	flag.DurationVar(&ldapTimeout, "timeout", 60*time.Second, "Timeout value to dial LDAP host")
	flag.StringVar(&ldapCACert, "cacert", "", "CA cert file used to authenticate LDAP server")
	flag.StringVar(&ldapProto, "protocol", "tcp", "protocol used to connect to LDAP server, tcp/udp")
	flag.StringVar(&ldapSvcUser, "user", "", "username to authenticate LDAP server")
	flag.StringVar(&ldapSvcPass, "pass", "", "password to authenticate LDAP server. Use '\\$' if '$'  is present")
	flag.StringVar(&userName, "ldapuser", "", "LDAP user to be queried")
	flag.BoolVar(&skipSslValdation, "skip_ssl_validation", false, "Skip SSL validation")
	flag.BoolVar(&env, "env", false, `Read parameters from environment variables instead of commaond line arguments. Enviroment variables:
    LDAP_HOST - LDAP host name
    LDAP_PORT - LDAP port to connect to
    LDAP_BASEDN - LDAP Base DN to search
    LDAP_SVC_USER - username to authenticate LDAP server
    LDAP_SVC_PASS - password to authenticate LDAP server. Use '\$' if '$'  is present
    LDAP_CACERT - CA cert file used to authenticate LDAP server
    LDAP_PROTO - protocol used to connect to LDAP server, tcp/udp
    SKIP_SSL_VALIDATION - skip SSL validation
    LDAP_QUERY_USER - LDAP user to be queried`)
	flag.BoolVar(&help, "help", false, "Print usage messages")

	flag.Parse()
	if flag.NFlag() <= 0 {
		fmt.Printf("No arguments set.\n")
		flag.PrintDefaults()
		os.Exit(2)
	}

	if help {
		flag.PrintDefaults()
		os.Exit(2)
	}

	if env {
		fmt.Println("Loading parameters from envrionment variables")
		ldapHost, _ = os.LookupEnv("LDAP_HOST")
		tmpstr, _ = os.LookupEnv("LDAP_PORT")
		ldapPort, err = strconv.Atoi(tmpstr)
		if err != nil {
			fmt.Printf("Invalid port, set it to 389\n")
			ldapPort = 389
		}
		ldapBaseDn, _ = os.LookupEnv("LDAP_BASED_DN")
		tmpstr, _ = os.LookupEnv("LDAP_TIMEOUT")
		i, err = strconv.Atoi(tmpstr)
		if err != nil {
			fmt.Printf("Invalid timeout, set it to 60\n")
			ldapTimeout = 60
		}
		ldapTimeout = time.Duration(i)
		ldapCACert, _ = os.LookupEnv("LDAP_CACERT")
		ldapProto, _ = os.LookupEnv("LDAP_PROTO")
		ldapSvcUser, _ = os.LookupEnv("LDAP_SVC_USER")
		ldapSvcPass, _ = os.LookupEnv("LDAP_SVC_PASS")
		tmpstr, _ = os.LookupEnv("SKIP_SSL_VALIDATION")
		skipSslValdation, err = strconv.ParseBool(tmpstr)
		if err != nil {
			fmt.Printf("Invalid SKIP-SSL-VALIDATON, set it to true\n")
			skipSslValdation = true
		}

		userName, _ = os.LookupEnv("LDAP_QUERY_USER")
	}
	
	if ldapHost == "" {
		fmt.Printf("Invalid LDAP host\n")
		return
	}

	if ldapBaseDn == "" {
		fmt.Printf("Invalid LDAP BaseDN\n")
		return
	}

	if ldapSvcUser == "" || ldapSvcPass == "" {
		fmt.Printf("Invalid LDAP username or password\n")
		return
	}

	if ldapPort < 1 || ldapPort > 65535 {
		fmt.Printf("LDAP Port out of range (1-65535), set it to default 389\n")
		ldapPort = 389
	}

	if ldapTimeout < 0 || ldapTimeout > 3600 {
		fmt.Printf("Invalid LDAP Timeout, must be within 0-3600, set it to default 120 seconds\n")
		ldapTimeout = 120 * 1000000000
	} else {
		ldapTimeout *= 1000000000
	}

	if strings.ToLower(ldapProto) != "tcp" || strings.ToLower(ldapProto) != "udp" {
		fmt.Printf("Invalid LDAP protocol %s, set it to 'tcp'\n", ldapProto)
		ldapProto = "tcp"
	}

	fmt.Println("-----------------------------------------")
	fmt.Printf("ldap-host: %s\nldap-port: %d\nldap-cacert: %s\n", ldapHost, ldapPort, ldapCACert)
	fmt.Printf("ldap-user: %s\nldap-passwd: %s\n", ldapSvcUser, ldapSvcPass)
	fmt.Printf("ldap-basedn: %s\nldap-proto: %s\n", ldapBaseDn, ldapProto)
	fmt.Printf("ldap-timeout: %d seconds\nskip-ssl-validation: %s\n", int(ldapTimeout.Seconds()), strconv.FormatBool(skipSslValdation))
	fmt.Println("LDAP user to query: ", userName)
	fmt.Println("-----------------------------------------")

	// dial
	addr := fmt.Sprintf("%s:%d", ldapHost, ldapPort)
	if ldapCACert != "" {
		ca_buf, err := ioutil.ReadFile(ldapCACert)
		if err != nil {
			fmt.Printf("failed to read ldap CACert file. Error: %s\n", err)
			return
		}
		roots = x509.NewCertPool()
		ok := roots.AppendCertsFromPEM(ca_buf)
		if !ok {
			fmt.Printf("Failed to load CA certificate from file - %s\n", ldapCACert)
			return
		}
		fmt.Printf("Loaded CA certificate from file: %s\n", ldapCACert)
	}

	if ldapPort == 636 {
		l, err = ldap.DialTLS(ldapProto, addr, &tls.Config{
			ServerName: ldapHost,
			RootCAs:    roots,
			InsecureSkipVerify: skipSslValdation,
		})
	} else {
		l, err = ldap.Dial(ldapProto, addr)
	}
	
	if err != nil {
		fmt.Printf("DIAL failed: %s\n", err)
		return
	}

	l.SetTimeout(ldapTimeout)
	defer l.Close()
	fmt.Printf("DIAL successfully\n")

	//bind
	err = l.Bind(ldapSvcUser, ldapSvcPass)
	if err != nil {
		fmt.Printf("BIND failed: %s\n", err)
		return
	}
	fmt.Printf("BIND successfully\n")

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		ldapBaseDn,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=User)(cn=%s))", ldap.EscapeFilter(userName)),
		nil,
		//[]string{"dn", "uidNumber", "gidNumber"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		fmt.Printf("Search user failed: %s\n", err)
		return
	}

	if len(sr.Entries) == 0 {
		fmt.Printf("User '%s' doesn't exist\n", userName)
		return
	}
	fmt.Printf("Found %d entries for user %s\n", len(sr.Entries), userName)
	sr.PrettyPrint(4)
}