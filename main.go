package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	tdp := flag.Float64("tdp", 0, "tdp in watt (erforderlich)")
	hours := flag.Float64("hours", 0, "laufzeit in stunden (erforderlich)")
	price := flag.Float64("price", 0, "preis pro kwh in € (falls nicht angegeben, wird marktpreis verwendet)")
	monthly := flag.Bool("monthly", false, "monatskosten berechnen (ignoriert hours parameter)")
	device := flag.String("device", "", "gerätename für validierung")
	quiet := flag.Bool("quiet", false, "ruhiger modus - minimale ausgabe")
	help := flag.Bool("help", false, "hilfenachricht anzeigen")

	flag.Parse()

	if *help {
		showHelp()
		return
	}

	if *tdp <= 0 {
		fmt.Fprintf(os.Stderr, "fehler: tdp muss größer als 0 sein\n")
		flag.Usage()
		os.Exit(1)
	}

	if !*monthly && *hours <= 0 {
		fmt.Fprintf(os.Stderr, "fehler: stunden müssen größer als 0 sein (außer bei verwendung des -monthly flags)\n")
		flag.Usage()
		os.Exit(1)
	}

	var pricePerKWh float64
	if *price > 0 {
		pricePerKWh = *price
	} else {
		var err error
		pricePerKWh, err = GetCurrentElectricityPrice()
		if err != nil {
			log.Fatalf("fehler beim abrufen des strompreises: %v", err)
		}
	}

	var logger Logger
	if *quiet {
		logger = &SilentLogger{}
	} else {
		logger = &ConsoleLogger{}
	}

	var cost float64
	var err error

	if *monthly {
		cost, err = CalculateMonthlyCost(*tdp, pricePerKWh, logger)
	} else if *device != "" {
		if validationErr := ValidateDevice(*tdp, *device); validationErr != nil {
			log.Fatalf("gerätevalidierung fehlgeschlagen: %v", validationErr)
		}
		cost, err = CalculateElectricityCost(*tdp, *hours, pricePerKWh, logger)
	} else {
		cost, err = CalculateElectricityCost(*tdp, *hours, pricePerKWh, logger)
	}

	if err != nil {
		log.Fatalf("fehler beim berechnen der kosten: %v", err)
	}

	if *quiet {
		fmt.Printf("%.2f\n", cost)
	} else {
		if *monthly {
			fmt.Printf("monatskosten für %.0fw gerät: %s\n", *tdp, FormatCostEuro(cost))
		} else {
			fmt.Printf("kosten für %.0fw gerät mit %.1f stunden laufzeit: %s\n", *tdp, *hours, FormatCostEuro(cost))
		}
	}
}

func showHelp() {
	fmt.Println("stromkosten rechner")
	fmt.Println("==================")
	fmt.Println()
	fmt.Println("berechnet stromkosten basierend auf tdp, laufzeit und strompreis.")
	fmt.Println()
	fmt.Println("verwendung:")
	fmt.Println("  go run . -tdp <watt> -hours <stunden> [-price <preis>] [-device <name>] [-monthly] [-quiet]")
	fmt.Println()
	fmt.Println("beispiele:")
	fmt.Println("  go run . -tdp 300 -hours 5                      # gaming pc für 5 stunden")
	fmt.Println("  go run . -tdp 100 -monthly                      # server monatskosten")
	fmt.Println("  go run . -tdp 250 -hours 8 -price 0.25          # benutzerdefinierter preis")
	fmt.Println("  go run . -tdp 350 -hours 6 -device 'gaming pc'  # mit validierung")
	fmt.Println("  go run . -tdp 75 -hours 24 -quiet               # ruhiger modus")
	fmt.Println()
	fmt.Println("parameter:")
	flag.PrintDefaults()
}
