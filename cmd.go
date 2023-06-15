package main

import (
	"context"
	"fmt"
	"github.com/Filecoin-Titan/titan-sdk-go"
	"github.com/Filecoin-Titan/titan-sdk-go/config"
	service "github.com/Filecoin-Titan/titan-sdk-go/titan"
	"github.com/Filecoin-Titan/titan-sdk-go/types"
	"github.com/cheggaaa/pb"
	"github.com/docker/go-units"
	"github.com/ipfs/go-cid"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
	"io"
	"net/http"
	"os"
	"time"
)

var downloadFileCmd = &cli.Command{
	Name:    "download",
	Aliases: []string{"d"},
	Usage:   "Get file from Titan network",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "block",
			Aliases: []string{"b"},
			Usage:   "pull data from block, it maybe slower",
		},
		&cli.Int64Flag{
			Name:        "size",
			Aliases:     []string{"s"},
			Usage:       "the bytes of per range request",
			DefaultText: "2MiB",
		},
		&cli.StringFlag{
			Name:    "cid",
			Aliases: []string{"c"},
			Usage:   "the asset CAR id you want to download",
		},
		&cli.StringFlag{
			Name:    "node",
			Aliases: []string{"n"},
			Usage:   "specify a node to pull data",
		},
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "the path you want to save the asset",
		},
		&cli.Int64Flag{
			Name:        "timeout",
			Aliases:     []string{"t"},
			DefaultText: "5",
			Usage:       "request timeout seconds",
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "Make the operation more talkative",
		},
	},
	Action: func(cctx *cli.Context) error {
		if cctx.String("cid") == "" {
			return fmt.Errorf("cid is required")
		}

		cid := cctx.String("cid")
		output := cctx.String("output")
		isBlock := cctx.Bool("block")
		verbose := cctx.Bool("verbose")
		timeout := cctx.Int64("timeout")

		address := os.Getenv("LOCATOR_API_INFO")
		opts := []config.Option{
			config.AddressOption(address),
			config.VerboseOption(verbose),
			config.TimeoutOption(time.Duration(timeout) * time.Second),
		}

		var decode bool

		if isBlock {
			opts = append(opts, config.TraversalModeOption(config.TraversalModeDFS))
		} else {
			decode = true
			opts = append(opts, config.TraversalModeOption(config.TraversalModeRange))
		}

		client, err := titan.New(opts...)
		if err != nil {
			return err
		}
		defer client.Close()

		return getFile(cctx.Context, client, cid, output, decode)
	},
}

var speedTestCmd = &cli.Command{
	Name:    "test",
	Aliases: []string{"t"},
	Usage:   "Test the bandwidth of nodes with specified resources",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "cid",
			Aliases: []string{"c"},
			Usage:   "the asset CAR id you want to test",
		},
		&cli.Int64Flag{
			Name:        "size",
			Aliases:     []string{"s"},
			Usage:       "the bytes of per range request",
			DefaultText: "2MiB",
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "Make the operation more talkative",
		},
	},
	Action: func(cctx *cli.Context) error {
		if cctx.String("cid") == "" {
			return fmt.Errorf("cid is required")
		}

		fmt.Println("Start testing node speed ...")

		id := cctx.String("cid")
		size := cctx.Int64("size")
		if size <= 0 {
			size = 2 << 20 // 2MiB
		}
		verbose := cctx.Bool("verbose")

		startTime := time.Now()
		address := os.Getenv("LOCATOR_API_INFO")
		opts := []config.Option{
			config.AddressOption(address),
			config.TraversalModeOption(config.TraversalModeRange),
			config.VerboseOption(verbose),
		}

		client, err := titan.New(opts...)
		if err != nil {
			return err
		}

		carId, _ := cid.Decode(id)
		service := client.GetTitanService()

		directlyNodes, natTraversalNode, err := service.GroupByEdges(cctx.Context, carId)
		if err != nil {
			return err
		}

		clients := make(map[string]*types.Client)
		for _, node := range directlyNodes {
			clients[node.NodeID] = &types.Client{
				Node:       node,
				HttpClient: service.GetDefaultClient(),
			}
		}

		ttfb1 := time.Since(startTime).Milliseconds()
		fmt.Printf("Direct Connection Nodes Time to First Byte(TTFB): %dms\n", ttfb1)

		_, err = service.Discover()
		if err != nil {
			return err
		}

		natTraversalClients, err := service.FilterAccessibleNodes(cctx.Context, natTraversalNode)
		if err != nil {
			return err
		}

		ttfb2 := time.Since(startTime).Milliseconds()
		fmt.Printf("Nat Traversal (Punching) Nodes Time to First Byte(TTFB): %dms\n", ttfb2)

		for k, v := range natTraversalClients {
			clients[k] = v
		}

		results, err := speedTest(cctx.Context, clients, id, size)
		if err != nil {
			return err
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Node", "IP", "Speed"})

		for _, result := range results {
			table.Append([]string{result.NodeID, result.IP, result.Speed})
		}

		table.Render()
		return nil
	},
}

func pullData(ctx context.Context, c *http.Client, node *types.Edge, cid string, size int64) (string, error) {
	start := time.Now()
	format := "car"
	header := http.Header{}
	header.Add("Range", fmt.Sprintf("bytes=%d-%d", 0, size))
	_, data, err := service.PullData(ctx, c, node, cid, format, header)
	if err != nil {
		log.Infof("pull data: %v", err)
		return "", err
	}
	cost := time.Since(start).Seconds()
	speed := units.BytesSize(float64(len(data))/cost) + "/s"

	return speed, err
}

func speedTest(ctx context.Context, clients map[string]*types.Client, cid string, size int64) ([]*TestResult, error) {
	var result []*TestResult

	for _, client := range clients {
		if client.HttpClient == nil {
			continue
		}

		speed, err := pullData(ctx, client.HttpClient, client.Node, cid, size)
		if err != nil {
			log.Warnf("pull data failed (%s): %v", client.Node.Address, err)
			continue
		}

		result = append(result, &TestResult{
			NodeID: client.Node.NodeID,
			IP:     client.Node.Address,
			Speed:  speed,
		})
	}

	return result, nil
}

func getFile(ctx context.Context, c *titan.Client, id, output string, decode bool) error {
	size, reader, err := c.GetFile(ctx, id)
	if err != nil {
		return err
	}
	defer reader.Close()

	bar := pb.New64(size).SetUnits(pb.U_BYTES)
	bar.ShowSpeed = true
	barR := bar.NewProxyReader(reader)

	bar.Start()
	defer bar.Finish()

	if output == "" {
		io.Copy(io.Discard, barR)
		return nil
	}

	carFilePath := output + ".car"
	file, err := os.Create(carFilePath)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, barR)
	if err != nil {
		return err
	}

	if !decode {
		return nil
	}

	if err := decodeCARFile(carFilePath, output); err != nil {
		return err
	}

	return os.Remove(carFilePath)
}

var runCmd = &cli.Command{
	Name:  "run",
	Usage: "Start Titan speed test server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "address",
			Aliases: []string{"addr"},
		},
		&cli.StringFlag{
			Name:     "username",
			Aliases:  []string{"user"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "password",
			Aliases:  []string{"passwd"},
			Required: true,
		},
	},
	Action: func(cctx *cli.Context) error {
		user := cctx.String("username")
		passwd := cctx.String("password")
		addr := cctx.String("address")
		if addr == "" {
			addr = ":8898"
		}

		address := os.Getenv("LOCATOR_API_INFO")
		opts := []config.Option{
			config.AddressOption(address),
			config.TraversalModeOption(config.TraversalModeRange),
		}

		client, err := titan.New(opts...)
		if err != nil {
			return err
		}

		serv := NewServer(client)
		return serv.Run(addr, user, passwd)
	},
}
