# Oracle Application / Service Metrics

## `oracle_response_latency`

* **purpose**
    * This prometheus histogram measures the RTT time taken (per request) from the `metrics_client`'s request to the oracle's server's response.
    * Observations from this histogram are measured in nano-seconds

## `oracle_responses`

* **purpose**
    * This prometheus counter measures the the # of oracle responses that a `metrics_client` has received
* **labels**
    * `status` := (failure, success)

## `oracle_ABCI_method_latency`

* **purpose**
    * This prometheus histogram measures the latency (per request) in seconds of ABCI method calls
    * The latency is measured over all of the slinky-specific code, and ignores any down-stream dependencies
* **labels**
    * `method`: one of (ExtendVote, PrepareProposal, ProcessProposal, VerifyVoteExtension, FinalizeBlock), this is the ABCI method that this latency report resulted from
