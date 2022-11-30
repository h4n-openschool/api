package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/h4n-openschool/api/api"
	"github.com/h4n-openschool/api/bus"
	"github.com/h4n-openschool/api/handlers"
	classRepos "github.com/h4n-openschool/api/repos/classes"
	studentRepos "github.com/h4n-openschool/api/repos/students"
	teacherRepos "github.com/h4n-openschool/api/repos/teachers"
	"github.com/h4n-openschool/api/server"
	"github.com/h4n-openschool/api/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Classes server",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Read a configured address from the command line, environment variable
		// or any detected config file.
		addr := viper.GetString("addr")

		logger, _ := zap.NewDevelopment()

		// Instantiate a new in-memory Teacher repository, generating 50 records.
		tr := teacherRepos.NewInMemoryTeacherRepository(10)

		// Instantiate a new in-memory Class repository, generating 50 records.
		cr := classRepos.NewInMemoryClassRepository(10)

		// Instantiate a new in-memory Student repository, generating 250 records.
		sr := studentRepos.NewInMemoryStudentRepository(4 * 30)

		// Create a new Gin router instance with the required middleware already
		// bootstrapped.
		gin.SetMode(gin.ReleaseMode)
		e := utils.ApplyMiddleware(gin.New(), logger)

		// Create Service Interface for codegen-based endpoint configuration
		si := handlers.OpenSchoolImpl{
			ClassRepository:   &cr,
			TeacherRepository: &tr,
			StudentRepository: &sr,
			Logger:            logger,
		}

		// Register codegen handlers from implemented functions
		e = api.RegisterHandlers(e, &si)

		// Create an HTTP server instance using the Gin handler
		s := server.Server{
			Addr:    addr,
			Handler: e,
			Logger:  logger,
		}

		// Have the server start serving Gin requests on the listen address
		// configured above.
		return s.Listen()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	var err error

	// Create a flag to hold the listen address for the server
	serveCmd.Flags().String("addr", "0.0.0.0:http", "The address to open a TCP listener on.")
	err = viper.BindPFlag("addr", serveCmd.Flags().Lookup("addr"))
	if err != nil {
		panic(err)
	}

	// Create a flag to configure the AMQP connection string

	serveCmd.Flags().String("amqp.dsn", "amqp://guest:guest@localhost:5672", "The DSN of the AMQP service to connect to for the event bus.")
	err = viper.BindPFlag("amqp.dsn", serveCmd.Flags().Lookup("amqp.dsn"))
	if err != nil {
		panic(err)
	}
}
