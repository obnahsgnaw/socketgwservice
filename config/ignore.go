package config

import (
	"strings"
)

type ignorer func(meth string, uri string) bool

var (
	backendIgnorer = []ignorer{
		adminIgnored,
	}
	backendAppIgnorer = []ignorer{
		//
	}
	backendAuthIgnorer = []ignorer{
		//
	}
	backendPermIgnorer = []ignorer{
		//
	}
	frontendIgnorer = []ignorer{
		proxyIgnored,
	}
	frontendAppIgnorer = []ignorer{
		//
	}
	frontendAuthIgnorer = []ignorer{
		//
	}
	frontendPermIgnorer = []ignorer{
		//
	}
)

func adminIgnored(_, uri string) bool {
	return strings.HasPrefix(uri, VersionRoute(AdminRoute)) || strings.HasPrefix(uri, VersionProjectRoute(AdminRoute))
}

func proxyIgnored(_, uri string) bool {
	return strings.HasPrefix(uri, VersionRoute("proxy")) || strings.HasPrefix(uri, VersionProjectRoute("proxy"))
}
