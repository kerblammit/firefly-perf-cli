# FireFly Performance CLI

FireFly Performance CLI is a HTTP load testing tool that generates a constant request rate against a [FireFly](https://github.com/hyperledger/firefly)
network and measure performance. This is used to confirm confidence that [FireFly](https://github.com/hyperledger/firefly)
can perform under normal conditions for an extended period of time.

## Items Subject to Testing

- Broadcasts (`POST /messages/broadcasts`)
- Private Messaging (`POST /messages/private`)
- Mint Tokens (`POST /tokens/mint`)
  - Fungible vs. Non-Fungible Token Toggle
- Blobs
- Contract Invocation (`POST /contracts/invoke`)
  - Ethereum vs. Fabric

## Run

The test configuration is structured around running `ffperf` as either a single process or in a distributed fashion as
multiple processes.

The tool has 2 basic modes of operation:

1. Run against a local FireFly stack
   - In this mode the `ffperf` tool loads information about the FireFly endpoint(s) to test by reading from a FireFly `stack.json` file on the local system. The location of the `stack.json` file is configured in the `instances.yaml` file by setting the `stackJSONPath` option.
2. Run against a remote Firefly node
   - In this mode the `ffperf` tool connects to a FireFly instance running on a different system. Since there won't be a FireFly `stack.json` on the system where `ffperf` is running the nodes to test must be configured in the `instances.yaml` file by settings the `Nodes` option.

### Local FireFly stack

See the [`Getting Started`](GettingStarted.md) guide for help running tests against a local stack.

In the test configuration you define one or more test _instances_ for a single `ffperf` process to run. An instance then
describes running one or more test _cases_ with a dedicated number of goroutine _workers_ against a _sender_ org and
a _recipient_ org. The test configuration consumes a file reference to the stack JSON configuration produced by the
[`ff` CLI](https://github.com/firefly-cli) (or can be defined manually) to understand the network topology, so that
sender's and recipient's just refer to indices within the stack.

As a result, running the CLI consists of providing an `instances.yaml` file describe the test configuration
and an instance index or name indicating which instance the process should run:

```shell
ffperf run -c /path/to/instances.yaml -i 0
```

See [`example-instances.yaml`](config/example-instances.yaml) for examples of how to define multiple instances
and multiple test cases per instance with all the various options.

### Remote FireFly node

See the [`Getting Started with Remote Nodes`](GettingStartedRemoteNode.md) guide for help running tests against a remote FireFly node.

In the test configuration you define one or more test _instances_ for a single `ffperf` process to run. An instance then
describes running one or more test _cases_ with a dedicated number of goroutine _workers_. Instead of setting a _sender_ org and
_recipient_ org (because there is no local FireFly `stack.json` to read) the instance must be configured to use a `Node` that has
been defined in `instances.yaml`.

Currently the types of test that can be run against a remote node are limited to those that only invoke a single endpoint. This makes
it most suitable for test types `token_mint`, `custom_ethereum_contract` and `custom_fabric_contract` since these don't need
responses to be received from other members of the FireFly network.

As a result, running the CLI consists of providing an `instances.yaml` file describe the test configuration
and an instance index or name indicating which instance the process should run:

```shell
ffperf run -c /path/to/instances.yaml -i 0
```

See [`example-remote-node-instances.yaml`](config/example-remote-node-instances.yaml) for examples of how to define nodes manually
and configure test instances to use them.

## Command line options

```
Executes a instance within a performance test suite to generate synthetic load across multiple FireFly nodes within a network

Usage:
  ffperf run [flags]

Flags:
  -c, --config string          Path to performance config that describes the network and test instances
  -d, --daemon                 Run in long-lived, daemon mode. Any provided test length is ignored.
      --delinquent string      Action to take when delinquent messages are detected. Valid options: [exit log] (default "exit")
  -h, --help                   help for run
  -i, --instance-idx int       Index of the instance within performance config to run against the network (default -1)
  -n, --instance-name string   Instance within performance config to run against the network
)
```

## Metrics

The `ffperf` tool registers the following metrics for prometheus to consume:

- ffperf_runner_received_events_total
- ffperf_runner_incomplete_events_total
- ffperf_runner_unexpected_events_total
- ffperf_runner_sent_mints_total
- ffperf_runner_sent_mint_errors_total
- ffperf_runner_deliquent_msgs_total
- ffperf_runner_perf_test_duration_seconds

## Distributed Deployment

See the [`ffperf` Helm chart](charts/ffperf) for running multiple instances of `ffperf` using Kubernetes.
