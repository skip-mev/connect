# Oracle Application / Service Metrics

## `oracle_response_latency`

* **purpose**
    * This prometheus histogram measures the RTT time taken (per request) from the `metrics_client`'s request to the oracle's server's response.
    * Observations from this histogram are measured in nano-seconds
* **labels**
    * `chain_id`: the chain-id of this oracle deployment

## `oracle_responses`

* **purpose**
    * This prometheus counter measures the the # of oracle responses that a `metrics_client` has received
* **labels**
    * `status` := (failure, success)
    * `chain_id`: the chain-id of this oracle deployment

## `oracle_ABCI_method_latency`

* **purpose**
    * This prometheus histogram measures the latency (per request) in seconds of ABCI method calls
    * The latency is measured over all of the slinky-specific code, and ignores any down-stream dependencies
* **labels**
    * `method`: one of (ExtendVote, PrepareProposal, ProcessProposal, VerifyVoteExtension, FinalizeBlock), this is the ABCI method that this latency report resulted from
    * `chain_id`: the chain-id of this oracle deployment

## `oracle_ABCI_method_status`

* **purpose**
    * This prometheus counter measures the number of ABCI requests, and their associated statuses
    * Each observation is either a success, or failure, and is paginated by the failure type
* **labels**
    * `method`: one of (ExtendVote, PrepareProposal, ProcessProposal, VerifyVoteExtension, FinalizeBlock), this is the ABCI method that this latency report resulted from
    * `chain_id`: the chain-id of this oracle deployment
    * `status`: The status of the request, if it's a failure, the label is an indication of what logic failed

## `oracle_message_size`

* **purpose**
    * This prometheus histogram tracks the size of vote-extensions, and extended commits that slinky is transmitting 
* **labels**
    * `chain_id`: the chain-id of this oracle deployment
    * `message_type`: the message-type whose size is being measured

## `oracle_prices`

* **purpose**
    * This prometheus gauge tracks the price written to state for each currency-pair
* **labels**
    * `chain_id`: the chain-id of this oracle deployment
    * `ticker`: the ticker for which the price was written to state
