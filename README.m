# Sappers (as Malazan Sappers) - EXPERIMENTAL

> "Who is in charge here, Commander?"

Kubernetes and Containers are broken by design. VM and hypervisors had been an issue just because are trying
to clone a full operating system and machine to run only a couple of prcesses. Containers looks better on the sense
that they share a common kernel, but that was never the issue on first place. Micro-vms and unikernel face directly
the problem, vm and hypervisor at fine if we focus on implement just what its needed to run ou processes, the is not
need at all to use an operating system and all hardware abstracion.

Under this context, orchestration and coreography must also evolve to something simpler. Sappers its created about
just 2 basic principles:

- Eventually know who are your peers
- Any peer can lead a management role, but only one at the same point on time.

The idea behind this is that the same agents ([no micro-service anymore](https://medium.com/@rodrigo.estrada/micro-agents-the-evolution-of-micro-services-1397a1567767))
are in charge to orchestrate themselves.

The technologies behind are:
 - gossip protocol
 - RAFT consensus
 - Distributed messages
 - Distributed immutable events databases (blockchain?).
