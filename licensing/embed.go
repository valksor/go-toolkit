package licensing

import (
	_ "embed"
)

// licenseText contains the BSD 3-Clause license text from the Valksor project.
//
//go:embed LICENSE
var licenseText string

// GetProjectLicense returns the embedded project license text.
func GetProjectLicense() string {
	return licenseText
}
