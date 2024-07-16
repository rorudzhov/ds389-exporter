package src

import "github.com/go-ldap/ldap/v3"

type LDAPCollector struct {
	LdapURL      string
	BindDN       string
	BindPassword string
	conn         *ldap.Conn
}

// Collect connects to a 389ds, performs a ldap search BaseDN "cn=monitor" and parses the result, returning as hashmap
func (collector *LDAPCollector) Collect() (map[string]string, error) {

	// Connection to LDAP
	conn, err := ldap.DialURL(collector.LdapURL)
	if err != nil {
		return nil, err
	}

	// Simple bind
	err = conn.Bind(collector.BindDN, collector.BindPassword)
	if err != nil {
		return nil, err
	}

	// Result HashMap
	result := make(map[string]string)

	// Looping array for extract attributes
	sResult, err := conn.Search(&ldap.SearchRequest{BaseDN: "cn=monitor", Filter: "(cn=*)", Scope: ldap.ScopeWholeSubtree})
	if err != nil {
		return nil, err
	}

	if sResult != nil {
		for _, entry := range sResult.Entries {
			for _, attr := range entry.Attributes {
				for _, value := range attr.Values {
					result[attr.Name] = value
				}
			}
		}
	}

	err = conn.Close()
	if err != nil {
		return result, err
	}

	return result, err
}
