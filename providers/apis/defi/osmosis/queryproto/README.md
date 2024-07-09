# Osmosis x/poolmanager Query Proto

This package contains the generated Go files for the `x/poolmanager` proto.  This is a 
temporary solution to expose only the query stubs but not import any main Osmosis code
to avoid module `replace` clauses.  In the future, this provider may be spun out as a 
separate gRPC service provider as its own module.
