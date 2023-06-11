package users

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"mafia/internal/db"
)

type Config struct {
	Port int
	Host string

	DBConfig *db.Config
}

type Controller struct {
	cfg    *Config
	router *gin.Engine

	dbcontroller  *db.Manager
	pdfController *PDFController
}

func MakeController(cfg *Config) (*Controller, error) {
	router := gin.New()
	router.LoadHTMLGlob("users/templates/*")
	gin.SetMode(gin.ReleaseMode)
	queue := make(chan *Request, 100)
	ctl := &Controller{
		cfg:    cfg,
		router: router,
		pdfController: &PDFController{
			queue: queue,
			data:  map[string]*Data{},
		},
	}
	var err error
	ctl.dbcontroller, err = db.NewManager(cfg.DBConfig)
	ctl.pdfController.dbController = ctl.dbcontroller

	if err != nil {
		return nil, fmt.Errorf("failed to make manager")
	}
	router.GET("/user", ctl.getUser)
	router.PUT("/user", ctl.updateUser)
	router.POST("/user", ctl.newUser)
	router.DELETE("/user", ctl.deleteUser)
	router.GET("/users", ctl.getUsers)

	router.GET("/genPdf", ctl.pdfGen)
	router.GET("/genUserPdf", ctl.pdfGenUser)
	router.GET("/getPdf", ctl.pdfStats)

	return ctl, nil
}

func (c *Controller) StartServer() error {
	fmt.Println("server started")
	eg := errgroup.Group{}
	eg.Go(func() error {
		if err := c.router.Run(fmt.Sprintf("%s:%d", c.cfg.Host, c.cfg.Port)); err != nil {
			return fmt.Errorf("failed to run server: %v", err)
		}
		return nil
	})
	eg.Go(func() error {
		c.pdfController.Run()
		return nil
	})
	err := eg.Wait()
	return err
}
