# dYdX Market Params API

## Overview

The dYdX Market Params API provides a list of all market parameters for the dYdX protocol. Specifically, this includes all markets that the dYdX protocol supports, the exchanges that need to provide prices along with the relevant translations, and a set of operations that convert all markets into a common set of prices. The side-car utilizes this API to update the price providers with all relevant market parameters, such that it can support the dYdX protocol out of the box.
