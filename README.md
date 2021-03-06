# Sappers (as Malazan Sappers) - EXPERIMENTAL

> "Who is in charge here, Commander?"

Kubernetes and Containers are broken by design. VM and hypervisors had been an issue just because are trying
to clone a full operating system and machine to run only a couple of processes. Containers looks better on the sense
that they share a common kernel, but that was never the issue on first place. Micro-vms and unikernel face directly
the problem, vm and hypervisor at fine if we focus on implement just what its needed to run specific processes, there is not
need at all to use an operating system and all hardware abstraction.

Under this context, orchestration and choreography must also evolve to something simpler. Run micro-vm over a cluster
means that even using micro-vms you will still managing a cluster, ideally, the  agents should be managed
by themselves. Sappers its created under just 3 basic principles:

- Eventually know who are your peers
- Any peer can lead a management role, but only one at the same point on time.
- Persistence and availability using dynamic replicas backed by snapshot for bootstrapping

The idea behind this is that the same agents ([no micro-service anymore](https://medium.com/@rodrigo.estrada/micro-agents-the-evolution-of-micro-services-1397a1567767))
are in charge to orchestrate themselves.

The technologies behind are:
 - micro(nano)-vms
 - unikernel
 - gossip protocol
 - RAFT consensus
 - Distributed messages
 - Distributed immutable events databases (blockchain?).

 ## Design principles

 - Management its done by a consensus group. The leader its only taking decisions and replicating the state
 - Push will be prefered over pull. for example, health checks will be pushed by apps
 - The managment layer will run as sidecar on any app (the current repo ;) )
 - The messaging will be handling using the sidecar (**planned provide block storage using micro-services and replicas using the sidecar**)
 - Event will be prefered always over RPC (REST)
 - Operations task will be handled by small specialized micro-vms (operator micro-vms)
 - Managment task like monitor and healing are done by specialized (operator) micro-vms launched and monitored by the leader
 - Any peers can be part of the consensus group, but will have preference peers with more free resources or consensus tagged peers
 - The communication will be provided by a NATS hub launched by the leader and monitored by a specialized micro-vm
 - Unikernel provisioner initially will be a Linux Container until research how to convert to Unikernel.

 > **Initial Restriction:** Due to limitations to run Sidecars on the current unikernel implementations, the language support will be limited to GO to allow run Sappers embedded.

 ## Progress (RoadMap)

 - [x] Bootstrap with minimun replica
 - [ ] :hourglass: Set member list using gossip protocol
 - [ ] Choose leader using RAFT between randomly number of replica peers:
  - Start with bootstrap peers
  - For consensus, choose peers with more free resources and gave preferences to tagged ones
  - When a RAFT consensus peer leave the cluster, the leader must peak one of the remaining peers and join to RAFT and NATS
 - [ ] Create embedded NATS cluster between bootstraping peers ([embedded test](https://github.com/nats-io/nats-server/blob/master/test/test.go#L46))
  - On the future, the NATS peers could be separate from the raft peers
 - [ ] Leader will keep managment database updated and will replicate to the rest of RAFT peers
 - [ ] Create GCP micro-vm launcher (sappers-infantry-GCP)
 - [ ] Peer leader will launch micro-vm launcher and will store connection inf
 - [ ] Let leader launch new micro-vms using images
 - [ ] Create healer micro vm
 - [ ] Peer leader will launch healer micro vm and will store connection info
 - [ ] Make avalaible all changes to all nodes (NATS publish/subscribe?)
