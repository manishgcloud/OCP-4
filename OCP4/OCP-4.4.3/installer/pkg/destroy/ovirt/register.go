package ovirt

import "github.com/openshift/installer/pkg/destroy/providers"

func init() {
	providers.Registry["ovirt"] = New
}
