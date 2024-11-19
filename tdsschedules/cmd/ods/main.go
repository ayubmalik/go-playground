// Command ods finds origin destination pairs
package main

import "tdsschedules"

func main() {

	tdsClient := tdsschedules.NewTDSClient()

	tdsClient.FindStops()
}
