package main

import (
	pkgadapter "knative.dev/eventing/pkg/adapter/v2"

	be "github.com/teapartycrypto/TeaParty/adams/pkg"
)

func main() {
	pkgadapter.Main("adams-adapter", be.EnvAccessorCtor, be.NewAdapter)
}
