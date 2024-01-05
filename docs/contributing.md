# General organization contribution guidelines

Please, read the [General Contribution Guidelines](https://github.com/upb-code-labs/docs/blob/main/CONTRIBUTING.md) before contributing.

## Static files micro-service contribution guidelines

Please, read the following guidelines about the project architecture and remember to update the [Insomnia Collection](./http/Insomnia.json) and [OpenAPI specification](./openapi/spec.yaml) if you make any changes to the REST endpoints.

### Architecture

The static files micro-service's architecture is based on the MVC (Model View Controller) pattern in order to separate the different components of the service and make it easier to maintain.

The project structure is as follows:

- `config`: Contains the configuration files for the micro-service (e.g. `environment.go`).

- `models`: Contains the models of the micro-service. The models are responsible for handling the data of the micro-service. Note that even though the micro-service does not use a database, the models are still used to obtain the data (files) from the file system.

- `views`: Contains the views of the micro-service. Note that the views are not used to render graphical interfaces, but to receive the data from the REST endpoints and call the corresponding controllers to handle the requests.

- `controllers`: Contains the controllers of the micro-service. The controllers are responsible for handling the incoming REST requests and returning the corresponding responses.

  - `schemas.go`: Contains the schemas of the requests and responses of the micro-service.

- `utils`: Contains the utility functions of the static files micro-service.

Additionally, the `files` folder is used to store the `.zip` archives to be used locally for testing and development purposes.

### Local development

### Dependencies

The following dependencies are required to run the static files micro-service locally:

- [Go 1.21.5](https://golang.org/doc/install)
- [Podman](https://podman.io/getting-started/installation) (To build and test the container image)

Please, note that `Podman` is a drop-in replacement for `Docker`, so you can use `Docker` instead if you prefer.

Additionally, you may want to install the following dependencies to make your life easier:

- [Air](https://github.com/cosmtrek/air) (for live reloading)


### Running the Static Files micro-service locally

You can start the Static Files micro-service by running the following command:

```bash
air 
```

This will start the Static Files micro-service and will watch for changes in the source code and restart the service automatically.

Additionally, you may want to generate a `.air.toml` file and add the `files/` directory to the `exclude_dir` list in order to avoid restarting the service when new `.zip` archives are added to the directory, to do this, run the following command or refer to the [Air documentation](https://github.com/cosmtrek/air): 

```bash
air init
```