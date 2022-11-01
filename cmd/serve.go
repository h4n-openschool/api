/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/h4n-openschool/server"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Classes server",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Read a configured address from the command line
		addr, err := cmd.Flags().GetString("addr")
		if err != nil {
			return err
		}

		// Create a new Gin router instance
		e := gin.Default()

		e.GET("/hello", func(ctx *gin.Context) {
			ctx.JSON(200, map[string]string{
				"hello": "world!",
			})
		})

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
}
