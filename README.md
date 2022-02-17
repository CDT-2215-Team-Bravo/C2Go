# C2Go

A basic tool that runs commands remotely

## Commands (For controller.go)
All communication happends over port 8085

``` connect [IP address] ```

Connects to [IP address] (Default port is 8085). Commands entered after connecting are ran on the [IP address] target and output is returned back

``` flood [IP address] [count] ```

Sends [count] messages to [IP address]

``` pingpong [IP address 1] [IP address 2] ```

Tells [IP address 1] to send a message to [IP address 2]. The two machines will continuously send messages back and forth to each other.
