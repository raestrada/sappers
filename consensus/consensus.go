package consensus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/raestrada/sappers/consensus/service"
	"github.com/raestrada/sappers/consensus/store"
	"github.com/raestrada/sappers/members"
	"go.uber.org/zap"
)

// Consensus gestiona el consenso de Raft.
type Consensus struct {
	raftDir      string
	inMem        bool
	httpAddr     string
	raftAddr     string
	joinAddr     string
	nodeID       string
	memberList   members.MemberList
	knownMembers map[string]struct{} // Rastrea los miembros conocidos para evitar uniones duplicadas
}

// ConsensusFactory es una fábrica para crear instancias de Consensus.
type ConsensusFactory struct{}

// Create crea una nueva instancia de Consensus usando la configuración global.
func (f ConsensusFactory) Create(memberList members.MemberList) *Consensus {
	cfg := config.GetConfig() // Obtener la configuración desde el singleton

	return &Consensus{
		raftDir:      cfg.RaftDir,
		inMem:        false, // Ajusta según sea necesario
		httpAddr:     cfg.HTTPAddr,
		raftAddr:     cfg.RaftAddr,
		nodeID:       cfg.NodeID,
		memberList:   memberList,
		knownMembers: make(map[string]struct{}), // Inicializar el mapa de miembros conocidos
	}
}

// Init inicializa el nodo de Raft y empieza el servicio de consenso.
func (c *Consensus) Init(ctx context.Context) {
	funcDesc := "Consensus - Init"

	// Asegurarse de que el directorio de Raft exista
	if c.raftDir == "" {
		zap.L().Fatal(funcDesc, zap.String("type", "No Raft storage directory specified"))
	}
	os.MkdirAll(c.raftDir, 0700)

	// Inicializar el almacén de Raft
	s := store.New(c.inMem)
	s.RaftDir = c.raftDir
	s.RaftBind = c.raftAddr

	// Abrir el almacén de Raft, ya sea como un nuevo clúster o uniéndose a uno existente
	if err := s.Open(c.joinAddr == "", c.nodeID); err != nil {
		zap.L().Fatal(funcDesc, zap.String("type", "failed to open store"), zap.Error(err))
	}

	// Iniciar el servicio HTTP para gestionar Raft
	h := service.New(c.httpAddr, s)
	if err := h.Start(); err != nil {
		zap.L().Fatal(funcDesc, zap.String("type", "failed to start HTTP service"), zap.Error(err))
	}

	zap.L().Info(funcDesc, zap.String("msg", "Raft node started successfully"), zap.String("nodeID", c.nodeID))

	// Comenzar a monitorizar gossip para detectar nuevos miembros
	go c.monitorGossipForNewMembers(ctx)

	// Esperar la señal de apagado
	select {
	case <-ctx.Done():
		zap.L().Info(funcDesc, zap.String("msg", "Shutting down Raft node"))
		return
	}
}

// monitorGossipForNewMembers verifica continuamente nuevos miembros y trata de agregarlos al clúster de Raft.
func (c *Consensus) monitorGossipForNewMembers(ctx context.Context) {
	funcDesc := "Consensus - monitorGossipForNewMembers"

	for {
		select {
		case <-ctx.Done():
			zap.L().Info(funcDesc, zap.String("msg", "Stopping gossip monitoring"))
			return

		case <-time.After(10 * time.Second): // Verifica cada 10 segundos
			members := c.memberList.Get()

			for _, member := range members {
				// Omitir el nodo actual y los miembros ya conocidos
				if member.Addr == c.raftAddr || c.isKnownMember(member.Addr) {
					continue
				}

				// Intentar agregar el nuevo miembro al clúster de Raft
				zap.L().Info(funcDesc, zap.String("msg", fmt.Sprintf("New member detected: %s. Attempting to join Raft cluster.", member.Addr)))
				if err := c.joinCluster(member.Addr, c.raftAddr, c.nodeID); err != nil {
					zap.L().Error(funcDesc, zap.String("msg", "Failed to join new member to Raft cluster"), zap.Error(err))
				} else {
					// Marcar el miembro como conocido
					c.markMemberAsKnown(member.Addr)
				}
			}
		}
	}
}

// isKnownMember verifica si un miembro ya es conocido.
func (c *Consensus) isKnownMember(addr string) bool {
	_, exists := c.knownMembers[addr]
	return exists
}

// markMemberAsKnown agrega la dirección de un miembro al mapa de miembros conocidos.
func (c *Consensus) markMemberAsKnown(addr string) {
	c.knownMembers[addr] = struct{}{}
}

// joinCluster envía una solicitud de unión a un nodo de Raft existente.
func (c *Consensus) joinCluster(joinAddr, raftAddr, nodeID string) error {
	funcDesc := "Consensus - JoinCluster"

	// Preparar el payload para la solicitud de unión
	payload := map[string]string{"addr": raftAddr, "id": nodeID}
	b, err := json.Marshal(payload)
	if err != nil {
		zap.L().Error(funcDesc, zap.String("type", "failed to marshal JSON"), zap.Error(err))
		return err
	}

	// Enviar la solicitud de unión
	url := fmt.Sprintf("http://%s/join", joinAddr)
	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		zap.L().Error(funcDesc, zap.String("type", fmt.Sprintf("post to join addr %s failed", joinAddr)), zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	// Comprobar si la solicitud de unión fue exitosa
	if resp.StatusCode != http.StatusOK {
		zap.L().Error(funcDesc, zap.String("type", "join request failed"), zap.String("status", resp.Status))
		return fmt.Errorf("failed to join cluster, status: %s", resp.Status)
	}

	zap.L().Info(funcDesc, zap.String("msg", "Joined cluster successfully"), zap.String("joinAddr", joinAddr))
	return nil
}
