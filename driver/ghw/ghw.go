package ghw

import (
	"github.com/jaypipes/ghw"
	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/driver/definition"
)

// Use GHW based Driver
func UseDriver(driver *definition.Driver){

	if driver.GetBoardModel == nil || driver.GetBoardVendor == nil {

		product, err := ghw.Product()
		if err == nil {

			// Get Board Model
			if driver.GetBoardModel == nil {
				productName := product.Name
				driver.GetBoardModel = func() string {
					return productName
				}
			}

			// Get Board Vendor
			if driver.GetBoardVendor == nil {
				productVendor := product.Vendor
				driver.GetBoardVendor = func() string {
					return productVendor
				}
			}
		}
	}

}

