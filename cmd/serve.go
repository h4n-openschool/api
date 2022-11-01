/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/h4n-openschool/classes/handlers"
	"github.com/h4n-openschool/classes/repos"
	"github.com/h4n-openschool/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Classes server",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Read a configured address from the command line
		addr := viper.GetString("addr")

		cr := repos.NewInMemoryClassRepository(50)

		// Create a new Gin router instance
		e := gin.Default()

		e.GET("/classes/:id", handlers.GetClass(&cr))

		// Create the server using the Gin handler
		s := server.Server{
			Addr:    addr,
			Handler: e.Handler(),
		}

		// Start the server itself
		return s.Listen()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Create a flag to hold the listen address for the server
	serveCmd.Flags().String("addr", "0.0.0.0:http", "The address to open a TCP listener on.")
	viper.BindPFlag("addr", serveCmd.Flags().Lookup("addr"))
}
