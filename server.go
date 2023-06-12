package main

import (
	"github.com/Filecoin-Titan/titan-sdk-go"
	"github.com/gin-gonic/gin"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
	"net/http"
	"strconv"
)

var log = logging.Logger("main")

type Server struct {
	titan *titan.Client
}

func NewServer(t *titan.Client) *Server {
	return &Server{
		titan: t,
	}
}

func (s *Server) Run(addr, username, passwd string) error {
	router := gin.Default()
	authorized := router.Group("/v1", gin.BasicAuth(gin.Accounts{
		username: passwd,
	},
	))
	authorized.GET("/speed", s.NodeSpeedHandler)
	return router.Run(addr)
}

type TestResult struct {
	NodeID string `json:"node_id"`
	IP     string `json:"ip"`
	Speed  string `json:"speed"`
	RTT    string `json:"RTT"`
}

func (s *Server) NodeSpeedHandler(c *gin.Context) {
	id := c.Query("cid")
	sizeQ := c.Query("size")
	size, err := strconv.ParseInt(sizeQ, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	carId, _ := cid.Decode(id)
	service := s.titan.GetTitanService()

	clients, err := service.GetAccessibleEdges(c.Request.Context(), carId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	results, err := speedTest(c.Request.Context(), clients, id, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"results": results,
	})
}
