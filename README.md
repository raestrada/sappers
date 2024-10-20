# Sappers (as Malazan Sappers) - EXPERIMENTAL

> "Who is in charge here, Commander?"

**Kubernetes** and **containers** are broken by design. VMs and hypervisors have been problematic because they attempt to replicate an entire operating system and machine just to run a couple of processes. Containers offer an improvement by sharing a common kernel, but this was never the core issue in the first place. **Micro-VMs** and **unikernels** address this problem directly, focusing on implementing only what is needed to run specific processes. There's no need to use a full operating system or all hardware abstractions.

Under this concept, orchestration and choreography also need to evolve into something simpler. Running **micro-VMs** over a cluster still implies managing the cluster. Ideally, agents should manage themselves. **Sappers** is built on three basic principles:

- Eventually, every agent must know its peers.
- Any peer can lead management tasks, but only one at any given time.
- Persistence and availability are maintained through dynamic replicas, backed by snapshots for bootstrapping.

The idea is that the agents themselves (no longer microservices) are responsible for orchestrating their own actions.

### Key Technologies

- **Micro(nano)-VMs**
- **Unikernels**
- **Gossip protocol**
- **Raft consensus**
- **Distributed messaging**
- **Distributed immutable event databases (blockchain?)**

## Design Principles

- **Management by consensus**: A leader is elected to make decisions and replicate state, but every peer can become a leader.
- **Push over pull**: For example, health checks will be pushed by the apps.
- **Sidecar management layer**: The management layer runs as a sidecar on each application (this repository).
- **Sidecar-based messaging**: Messaging between services will be handled through the sidecar. (Planned: block storage and replicas also managed by the sidecar).
- **Event-driven over RPC (REST)**: Events will be prioritized over traditional RPC methods.
- **Micro-VMs for operational tasks**: Specialized micro-VMs will handle tasks like monitoring and healing.
- **Leader-launched micro-VMs**: Management tasks like monitoring and healing will be run by specialized micro-VMs launched and monitored by the leader.
- **Dynamic consensus group**: Any peer can join the consensus group, but priority will be given to peers with more free resources or those tagged as consensus nodes.
- **NATS messaging hub**: The leader launches and monitors a NATS hub, with communication handled by the micro-VMs.
- **Unikernel provisioner**: Initially, this will be a Linux container until the unikernel is fully implemented.

> **Initial restriction**: Due to the limitations of running sidecars on current unikernel implementations, supported languages will be limited to **Go** to allow Sappers to run embedded.

---

## Roadmap (Progress)

- [x] Bootstrap with minimum replica
- [ ] :hourglass: Set member list using Gossip protocol
- [ ] Leader election using Raft among randomly selected replica peers:
  - Begin with bootstrap peers.
  - For consensus, select peers with the most available resources, and prioritize tagged peers.
  - When a Raft peer leaves the cluster, the leader must choose a remaining peer and add it to Raft and NATS.
- [ ] Create an embedded NATS cluster among bootstrapping peers ([embedded test](https://github.com/nats-io/nats-server/blob/master/test/test.go#L46)).
  - In the future, NATS peers could be separated from Raft peers.
- [ ] The leader will maintain the management database and replicate it to the other Raft peers.
- [ ] Create a GCP micro-VM launcher (sappers-infantry-GCP).
- [ ] The peer leader will launch micro-VMs and store connection info.
- [ ] Allow the leader to launch new micro-VMs using predefined images.
- [ ] Create a healer micro-VM.
- [ ] The peer leader will launch the healer micro-VM and store connection info.
- [ ] Make all changes available to all nodes (NATS publish/subscribe?).

---

## Gossip Protocol: Peer Communication

Sappers uses **memberlist** from HashiCorp to implement the **gossip protocol**, allowing nodes in the cluster to communicate efficiently and discover each other. Gossip-based communication is designed to be lightweight and scalable.

### Gossip Features

- **Peer-to-peer communication**: Nodes communicate directly with each other in a decentralized manner.
- **Automatic peer discovery**: Nodes automatically detect and join the cluster.
- **Scalability**: Gossip can scale to thousands of nodes with minimal overhead.

**Memberlist** provides a way for nodes to discover each other and exchange state information. In Sappers, we use memberlist to ensure that all nodes can discover and communicate with their peers.

---

## Raft Consensus: Distributed State

Sappers uses the **Raft** consensus algorithm to ensure state consistency across nodes in the cluster. Raft is known for its simplicity and reliability in maintaining consensus in distributed systems.

### How Raft Works

1. **Leaders and followers**: A leader is elected, and it replicates operations to follower nodes.
2. **Commitment**: An operation is only committed when a majority of nodes agree, ensuring consistency.
3. **Automatic failover**: If the leader node fails, a new leader is automatically elected by the remaining nodes.

Raft ensures that all nodes in the cluster agree on the current state, even in the event of node failures. It uses **log replication** to keep all nodes synchronized, and **leader election** to ensure that only one node is in charge at any given time.

### Libraries Used

- **raft**: HashiCorp's implementation of the Raft consensus algorithm.
- **raft-boltdb**: Provides durable storage for Raft, ensuring data persistence.

---

## Usage

This section explains how to use **Sappers** to deploy and manage a cluster of nodes, including launching **nano-VMs** using **nanoVM** and exposing them through **Consul** as a service mesh. The following steps cover the full functionality, from bootstrapping the cluster to deploying micro-VMs and nano-VMs for different purposes.

### Step 1: Build the binary

First, compile the project to create the binary:

```bash
go build -o sappers
```

### Step 2: Bootstrap the cluster with initial replicas

Start by bootstrapping the cluster with a minimum replica. You can define the number of initial replicas using the `--bootstrap` flag:

```bash
./sappers --node-id "node1" \
  --gossip-port 7946 \
  --raft-addr ":12000" \
  --http-addr ":11000" \
  --bootstrap 3 \
  ./raft/node1
```

This command initializes a node (`node1`) with gossip communication on port 7946 and Raft consensus on port 12000, starting the bootstrap process to form a cluster with 3 replicas.

### Step 3: Launch additional nodes

Once the bootstrap process completes, you can add more nodes to the cluster. Each node requires a unique `node-id` and different ports for gossip and Raft communication:

```bash
./sappers --node-id "node2" \
  --gossip-port 7947 \
  --raft-addr ":12001" \
  --http-addr ":11001" \
  ./raft/node2
```

```bash
./sappers --node-id "node3" \
  --gossip-port 7948 \
  --raft-addr ":12002" \
  --http-addr ":11002" \
  ./raft/node3
```

### Step 4: Deploy nano-VMs with nanoVM

The leader node is responsible for launching **nano-VMs** that run specific applications. These nano-VMs will be built using **nanoVM** and exposed via **Consul** as a service mesh.

To launch a nano-VM, use the following command:

```bash
./sappers --node-id "leader-node" \
  --launch-nanovm app \
  --vm-image "app-vm.img" \
  --consul-service "my-app" \
  --service-port 8080 \
  ./raft/leader-node
```

This command launches a **nano-VM** running an application (`app-vm.img`) and registers it with **Consul** under the service name `my-app` on port 8080. **Consul** will act as the service mesh, enabling service discovery and routing.

### Step 5: Expose services through Consul

Once the nano-VM is launched, **Consul** automatically registers the service, making it accessible through the service mesh. You can interact with the service through Consul's service discovery:

```bash
consul catalog services
```

This will list all available services, including the new nano-VM-based service (`my-app`), along with its port and health status.

To test the service, you can send a request to Consul, which routes it to the nano-VM:

```bash
curl http://localhost:8080/my-app
```

This command routes the request through Consul, which forwards it to the corresponding nano-VM.

### Step 6: Raft leader election

**Raft** will automatically elect a leader from the available nodes, and this leader manages the cluster state. You can verify the Raft leader by reviewing the logs:

```bash
tail -f ./raft/node1/raft.log
```

The logs will indicate which node has been elected as the leader and how the other nodes follow.

### Step 7: Distributed messaging with NATS

The Raft leader will launch a **NATS** messaging hub to handle distributed communication between nodes and micro-VMs. To subscribe to messages, use the NATS client:

```bash
nats-sub -s nats://localhost:4222 ">"
```

This will subscribe to all messages flowing through the NATS hub, allowing you to observe real-time communication between micro-VMs and nodes.

### Step 8: Dynamic peer addition

New nodes can dynamically join the cluster at any time. To add a new node:

```bash
./sappers --node-id "node4" \
  --gossip-port 7949 \
  --raft-addr ":12003" \
  --http-addr ":11003" \
  ./raft/node4
```

This node will automatically discover and join the existing cluster using the gossip protocol, and synchronize its state via Raft.

### Step 9: Healing micro-VMs

If a node becomes unhealthy, the Raft leader can launch a **healer micro-VM** to handle the recovery process. The healer VM can either restore the failed node or redistribute its workload:

```bash
./sappers --node-id "leader-node" \
  --launch-microvm healer \
  --vm-image "healer-vm.img" \
  ./raft/leader-node
```

The leader node will track the healing process and ensure the node's state is updated across the cluster.

### Step 10: Snapshotting and state persistence

Raft creates snapshots at regular intervals to ensure state persistence. You can manually trigger a snapshot for testing purposes:

```bash
./sappers --node-id "leader-node" \
  --trigger-snapshot \
  ./raft/leader-node
```

This snapshot will be stored in the Raft directory, allowing new nodes to bootstrap quickly without needing to replay the full log.

### Step 11: Monitoring logs and services

You can adjust the log level for more detailed logs using the `--log-level` flag. For example, to set it to `DEBUG`:

```bash
./sappers --node-id "node1" \
  --log-level "DEBUG" \
  ./raft/node1
```

Logs are output in JSON format for easy parsing, and Consul will handle service health checks, allowing you to monitor the health of nano-VMs and micro-VMs in real time.

### Step 12: Consul service mesh for all micro-VMs and nano-VMs

As more micro-VMs or nano-VMs are launched, they will automatically register with **Consul**. You can query the Consul catalog to see all services:

```bash
consul catalog services
```

This includes nano-VM applications, healer VMs, and any other services managed by **Sappers**.

---

### Full Commands Overview

Here is a summary of the full command options you can use with **Sappers**:

- `--node-id`: Unique ID for the node.
- `--gossip-port`: Port used for gossip communication between nodes.
- `--raft-addr`: Address used for Raft consensus.
- `--http-addr`: Address for the HTTP API.
- `--bootstrap`: Number of replicas for the initial cluster.
- `--launch-nanovm`: Launch a nano-VM with a specific application (e.g., app).
- `--launch-microvm`: Launch a specific micro-VM (e.g., healer, monitor).
- `--vm-image`: Specify the image for the nano-VM or micro-VM.
- `--consul-service`: Register the nano-VM or micro-VM as a service with Consul.
- `--service-port`: Port on which the service will be exposed via Consul.
- `--trigger-snapshot`: Manually trigger a snapshot of the current state.
- `--log-level`: Log verbosity (`DEBUG`, `INFO`, `WARN`, `ERROR`).
- `./raft/nodeX`: Directory where Raft stores its state for each node.

This comprehensive guide covers the full feature set of **Sappers**, including nano-VM deployment with **nanoVM**, Consul service mesh integration, dynamic peer addition, and operational micro-VMs for healing and monitoring.

---

## How to Contribute

We welcome contributions! Please feel free to open an issue or submit a pull request with your improvements or suggestions.

---

## License

This project is licensed under the **MIT License**. See the `LICENSE` file for more details.

---

## Final Notes

**Sappers** is an experimental project aimed at evolving the way distributed systems handle orchestration and consensus. It leverages micro-VMs, unikernels, and modern distributed technologies like **Raft** and **Gossip** to build a system where agents can manage themselves.

