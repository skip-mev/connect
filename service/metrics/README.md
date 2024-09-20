# Oracle Application / Service Metrics

The following metrics are registered to the Cosmos SDK metrics port (default :26660):

## `oracle_response_latency`

* **purpose**
    * This prometheus histogram measures the RTT time taken (per request) from the `metrics_client`'s request to the oracle's server's response.
    * Observations from this histogram are measured in nanoseconds
* **labels**
    * `chain_id`: the chain-id of this oracle deployment

## `oracle_responses`

* **purpose**
    * This prometheus counter measures the # of oracle responses that a `metrics_client` has received
* **labels**
    * `status` := (failure, success)
    * `chain_id`: the chain-id of this oracle deployment

## `ABCI_method_latency`

* **purpose**
    * This prometheus histogram measures the latency (per request) in seconds of ABCI method calls
    * The latency is measured over all connect-specific code, and ignores any down-stream dependencies
* **labels**
    * `method`: one of (ExtendVote, PrepareProposal, ProcessProposal, VerifyVoteExtension, FinalizeBlock), this is the ABCI method that this latency report resulted from
    * `chain_id`: the chain-id of this oracle deployment

## `ABCI_method_status`

* **purpose**
    * This prometheus counter measures the number of ABCI requests, and their associated statuses
    * Each observation is either a success, or failure, and is paginated by the failure type
* **labels**
    * `method`: one of (ExtendVote, PrepareProposal, ProcessProposal, VerifyVoteExtension, FinalizeBlock), this is the ABCI method that this latency report resulted from
    * `chain_id`: the chain-id of this oracle deployment
    * `status`: The status of the request, if it's a failure, the label is an indication of what logic failed

## `message_size`

* **purpose**
    * This prometheus histogram tracks the size of vote-extensions, and extended commits that connect is transmitting 
* **labels**
    * `chain_id`: the chain-id of this oracle deployment
    * `message_type`: the message-type whose size is being measured

## `oracle_prices`

* **purpose**
    * This prometheus gauge tracks the price written to state for each currency-pair
* **labels**
    * `chain_id`: the chain-id of this oracle deployment
    * `ticker`: the ticker for which the price was written to state

## `oracle_reports_per_validator`

* **purpose**
    * This prometheus gauge tracks the prices that each validator has reported for any block per ticker
* **labels**
    * `chain_id`: the chain-id of this oracle deployment
    * `ticker`: the ticker for which the price was written to state
    * `validator`: the consensus address of the validator that made the report

## `oracle_report_status_per_validator`

* **purpose**
    * This prometheus counter tracks the # of reports per validator and their status (absent: nil-vote, missing_price: the validator's vote was included but w/o a price, and with_price: validator's vote was included with a price)
* **labels**
    * `chain_id`: the chain-id of this oracle deployment
    * `ticker`: the ticker for which the price was written to state
    * `validator`: the consensus address of the validator that made the report
