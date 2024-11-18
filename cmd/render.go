/*
Copyright © 2024 Daniel Lawrence <go@dansysadm.com>
*/
package cmd

import (
	"fmt"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/daniellawrence/micro-hiera/lib"
	"github.com/spf13/cobra"
)

var (
	outFile string = "/dev/stdout"
)

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:    "render",
	Short:  "A brief description of your command",
	PreRun: toggleDebug,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			log.Fatalf("Expected 2 or more input files, got: %d, %s\n", len(args), args)
			return
		}

		m := lib.Merger{}
		mergedObj := m.MergeFiles(args)
		mergedBytes, _ := yaml.Marshal(mergedObj)

		for _, v := range m.Voliations {
			logrus.WithFields(logrus.Fields{}).Log(v.Level, v.Error())
		}

		if len(m.Voliations) > 0 {

			fmt.Printf("\n%-30s %s\n", "violation", "count")
			fmt.Printf("--------------------------------------------------\n")
			for name, count := range m.CountViolationByType() {
				fmt.Printf("%-30s %-d\n", name, count)
			}
		} else {
			fmt.Println(string(mergedBytes))
		}

	},
}

func init() {
	rootCmd.AddCommand(renderCmd)
	renderCmd.Flags().StringVarP(&outFile, "out", "o", outFile, "file that will contain the newly merged data")
}
