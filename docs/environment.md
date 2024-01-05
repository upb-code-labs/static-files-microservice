# Environment

This document describes the required environment variables to run the micro-service.

| Name                   | Description                                                                                                                      | Example                       | Mandatory |
| ---------------------- | -------------------------------------------------------------------------------------------------------------------------------- | ----------------------------- | --------- |
| `ARCHIVES_VOLUME_PATH` | Path to the directory where the archives will be stored. It is supposed to be a Docker volume in production to persist the data. | `/var/lib/archives`           | Yes       |
| `TEMPLATES_PATH`       | Path to the directory where the templates are stored.                                                                            | `/var/lib/archives/templates` | No        |
