# internal-api
API used by the BarnBridge App. 

## Configuration
For all the available config options please see the [configuration file](./config-generated.yml). A sample configuration
file can be found [here](./config-sample.yml).

The tool also includes rich help functions. To access them, use the `--help` flag.

```shell
./internal-api --help
```

## Running
The API can be run using the following command:
```shell
./internal-api run
```

## Implementation details
Each product is grouped into its own set of routes and implementations.

Example: for governance related endpoints, see the [governance](./governance) module.

