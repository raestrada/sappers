package consensus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"go.uber.org/zap"

	"github.com/raestrada/sappers/consensus/service"
	"github.com/raestrada/sappers/consensus/store"
)

// Consensus ...
type Consensus struct {
	raftDir  string
	inMem    bool
	httpAddr string
	raftAddr string
	joinAddr string
	nodeID   string
}

// New ...
func New(raftDir, httpAddr, raftAddr, joinAddr, nodeID string, inMem bool) Consensus {
	return Consensus{
		raftDir:  raftDir,
		inMem:    inMem,
		httpAddr: httpAddr,
		raftAddr: raftAddr,
		joinAddr: joinAddr,
		nodeID:   nodeID,
	}
}

// Init ...
func (c *Consensus) Init(ctx context.Context) {
	funcDesc := "Consensus - Init"

	// Ensure Raft storage exists.
	if c.raftDir == "" {
		zap.L().Fatal(
			funcDesc,
			zap.String("type", "No Raft storage directory specified"),
		)
	}
	os.MkdirAll(c.raftDir, 0700)

	s := store.New(c.inMem)
	s.RaftDir = c.raftDir
	s.RaftBind = c.raftAddr
	if err := s.Open(c.joinAddr == "", c.nodeID); err != nil {
		zap.L().Fatal(
			funcDesc,
			zap.String("type", "failed to open store"),
			zap.String("msg", err.Error()),
		)
	}

	h := service.New(c.httpAddr, s)
	if err := h.Start(); err != nil {
		zap.L().Fatal(
			funcDesc,
			zap.String("type", "failed to start HTTP service"),
			zap.String("msg", err.Error()),
		)
	}

	// If join was specified, make the join request.
	if c.joinAddr != "" {
		if err := c.join(c.joinAddr, c.raftAddr, c.nodeID); err != nil {
			zap.L().Fatal(
				funcDesc,
				zap.String("type", fmt.Sprintf("failed to join node at %s", c.joinAddr)),
				zap.String("msg", err.Error()),
			)
		}
	}

	zap.L().Info(
		funcDesc,
		zap.String("msg", "hraftd started successfully"),
	)

	select {
	case <-ctx.Done():
		zap.L().Info(
			funcDesc,
			zap.String("msg", "hraftd exiting"),
		)
		return
	}
}

func (c *Consensus) join(joinAddr, raftAddr, nodeID string) error {
	funcDesc := "Consensus - Join"

	b, err := json.Marshal(map[string]string{"addr": raftAddr, "id": nodeID})
	if err != nil {
		zap.L().Error(
			funcDesc,
			zap.String("type", "failing un-marshaling"),
			zap.String("msg", err.Error()),
		)

		return err
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/join", joinAddr), "application-type/json", bytes.NewReader(b))
	if err != nil {
		zap.L().Error(
			funcDesc,
			zap.String("type", "post to join addr %s"),
			zap.String("msg", err.Error()),
		)
		return err
	}
	defer resp.Body.Close()

	return nil
}
