package cmd

import (
	middleware "github.com/deepmap/oapi-codegen/pkg/gin-middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/h4n-openschool/classes/api"
	"github.com/h4n-openschool/classes/bus"
	"github.com/h4n-openschool/classes/handlers"
	classRepos "github.com/h4n-openschool/classes/repos/classes"
	"github.com/h4n-openschool/classes/server"
	"github.com/h4n-openschool/classes/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Classes server",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Read a configured address from the command line, environment variable
		// or any detected config file.
		addr := viper.GetString("addr")

		// Instantiate a new in-memory Class repository, generating 50 records.
		cr := classRepos.NewInMemoryClassRepository(50)

    // Create a new message bus instance
    b := bus.GetOrCreateBus()

    // Try to connect to the message bus or fail on error
    if err := b.Connect(); err != nil {
      return err
    }

		// Create a new Gin router instance
		e := gin.Default()

    // Configure CORS, allowing all origins
    corsConf := cors.DefaultConfig()
    corsConf.AllowAllOrigins = true
    e.Use(cors.New(corsConf))

    // Add request validation from codegen
    swagger, _ := api.GetSwagger()
    opts := middleware.Options{
      ErrorHandler: utils.ValidatorFunc,
    }
    e.Use(middleware.OapiRequestValidatorWithOptions(swagger, &opts))

    // Create Service Interface for codegen-based endpoint configuration
    si := handlers.OpenSchoolImpl{
      Repository: &cr,
      Bus: b,
    }

    // Register codegen handlers from implemented functions
    e = api.RegisterHandlers(e, &si)

		// Create an HTTP server instance using the Gin handler
		s := server.Server{
			Addr:    addr,
			Handler: e,
		}

		// Have the server start serving Gin requests on the listen address
		// configured above.
		return s.Listen()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Create a flag to hold the listen address for the server
	serveCmd.Flags().String("addr", "0.0.0.0:http", "The address to open a TCP listener on.")
  err := viper.BindPFlag("addr", serveCmd.Flags().Lookup("addr"))
  if err != nil {
    panic(err)
  }
}
