package cmd

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/h4n-openschool/classes/bus"
	classRepos "github.com/h4n-openschool/classes/repos/classes"
	"github.com/h4n-openschool/classes/handlers/classes"
	"github.com/h4n-openschool/server"
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

		// Instantiate a new Class repository
		// NOTE: For the purposes of mocking, the InMemoryClassRepository will
		// generate a number of classes at random. In this case, we instruct it
		// to generate 50 classes at once.
		cr := classRepos.NewInMemoryClassRepository(50)

    // Create a new message bus instance
    b := bus.GetOrCreateBus()

    // Try to connect to the message bus or fail on error
    if err := b.Connect(); err != nil {
      return err
    }

		// Create a new Gin router instance
		e := gin.Default()

    corsConf := cors.DefaultConfig()
    corsConf.AllowAllOrigins = true
    e.Use(cors.New(corsConf))

		e.GET("/debug", func(ctx *gin.Context) {
			// This handler will return arrays of all mocked data so that the
			// developer can use it for testing.
			//
			// In a real-world deployment with a database configured, this
			// endpoint would be removed as the developer would be able to just
			// query the database locally.
			ctx.JSON(200, gin.H{
				"classes.items": cr.Items,
			})
		})

		// Register routes
		classEndpoints := e.Group("/classes")
		{
			classEndpoints.GET("/", classes.GetClasses(&cr))   // list all classes
			classEndpoints.GET("/:id", classes.GetClass(&cr))  // get class by id
			classEndpoints.POST("/", classes.CreateClass(&cr)) // create new class
		}

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
	viper.BindPFlag("addr", serveCmd.Flags().Lookup("addr"))
}
