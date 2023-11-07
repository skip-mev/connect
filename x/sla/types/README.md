# Price Feed Service Level Agreements (SLA)

## Overview

The price feed service level agreement (SLA) is a set of rules that govern how often validators must post price updates to the oracle. The specific set of rules is configurable by the application and chain developers and is enforced by the protocol. The SLA enforces that validators maintain a minimum uptime, with slashing conditions for downtime. 

## SLA Parameters

There are five key parameters that govern the SLA:

### MaximumViableWindow

This determines the previous number of blocks from the given height that are considered for the SLA. This parameter acts as a sliding window. 

For example, if the `maximumViableWindow` is set to 100, then the SLA will only consider the previous 100 blocks from the current block height.

### MinimumBlockUpdates

This determines the minimum number of blocks that the validator had to have voted on in the `maximumViableWindow` in order to be considered. This value must be strictly less than the `maximumViableWindow`. If a validator has not voted on at least `minimumBlockUpdates` blocks in the `maximumViableWindow`, they will not be considered for the SLA.

For example, if the `maximumViableWindow` is set to 100 and the `minimumBlockUpdates` is set to 50, then the SLA will only consider validators that have voted on at least 50 blocks in the previous 100 blocks.

### ExpectedUptime

Given the validator has voted on at least `minimumBlockUpdates` blocks within the `maximumViableWindow`, this determines the minimum percentage of the blocks that had to have included a price update in the validator's vote extension. 

For example, if the `expectedUptime` is set to 0.9, then the validator must have included price updates in their votes on at least 90% of the blocks they have voted on. If they underperform, they will be slashed.

### SlashConstant

This constant determines the how much the validator will be slashed if they deviate from the expected uptime. The formula for slashing is:

```go
slashPercentage := ((expectedUptime - actualUptime) / expectedUptime) * slashConstant
```

This formula will slash the validator proportionally to how much they deviate from the expected uptime.

### Frequency

Frequency defines how often the criteria of an SLA should be checked. This is a parameter that is set by the chain developer and/or chain governance. The frequency is set in terms of blocks. For example, if the frequency is set to 10, then the SLA will be checked every 10 blocks. This parameter should be less than the `maximumViableWindow` - otherwise the SLA will not be able to be enforced.

## Slashing

As described above, slashing is variable to how far the validator's uptime deviates from the expected uptime. Slashing is proportional to the each validator's power and therefore is relative. The larger the validator, the more they will be slashed. This is expected as larger validators have a larger say in the final aggregated price that is posted on chain to the `x/oracle` module.
