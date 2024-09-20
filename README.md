
![Logo](https://ik.imagekit.io/pj3r6oe9k/prasorganic-high-resolution-logo-transparent.svg?updatedAt=1726835541390)

# Prasorganic User Service

Prasorganic User Service is one of the components in the Prasorganic Microservice architecture Built with Go (Golang). This service supports user data management operations via RESTful API and gRPC.

## Tech Stack

[![My Skills](https://skillicons.dev/icons?i=go,docker,postgresql,redis,bash,git&theme=light)](https://skillicons.dev)

## Features

- **User Management:** Supports operations for creating, retrieving, and updating user data.

- **RESTful API:** Provides a RESTful API using Fiber with various middleware for managing requests and responses.

- **gRPC:** Utilizes gRPC for inter-service communication, equipped with interceptors for handling requests and responses.

- **Image Storage:** Integrates with ImageKit for storing and managing user profile photos.

- **Database:** Uses PostgreSQL for data storage with database migration support.

- **Caching:** Redis is used for data caching.

- **Logging:** Logs are recorded using Logrus.

- **Error Handling:** Implements error handling to ensure proper detection and handling of errors, minimizing the impact on both the client and server.

- **System Resilience:** Uses a Circuit Breaker to enhance application resilience and fault tolerance, protecting the system from cascading failures.

- **Configuration and Security:** Employs Viper and HashiCorp Vault for integrated configuration and security management.

- **Testing:** Implements unit testing using Testify.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

This project makes use of third-party packages and tools. The licenses for these
dependencies can be found in the `LICENSES` directory.

## Dependencies and Their Licenses

- `Docker:` Licensed under the Apache License 2.0. For more information, see the [Docker License](https://github.com/docker/docs/blob/main/LICENSE).

- `Docker Compose:` Licensed under the Apache License 2.0. For more information, see the [Docker Compose License](https://github.com/docker/compose/blob/main/LICENSE).

- `PostgreSQL:` Licensed under PostgreSQL License. For more information,see the [PostgreSQL License](https://www.postgresql.org/about/licence/).

- `Redis:` Follows a dual-licensing model with RSALv2 and SSPLv1. For more information, see the [Redis License](https://redis.io/legal/licenses/).

- `RedisInsight:` Licensed under the RedisInsight License. For more information, see the [RedisInsight License](https://github.com/RedisInsight/RedisInsight/blob/main/LICENSE).

- `GNU Make:` Licensed under the GNU General Public License v3.0. For more information, see the [GNU Make License](https://www.gnu.org/licenses/gpl.html).

- `GNU Bash:` Licensed under the GNU General Public License v3.0. For more information, see the [Bash License](https://www.gnu.org/licenses/gpl-3.0.html).

- `Git:` Licensed under the GNU General Public License version 2.0. For more information, see the [Git License](https://opensource.org/license/GPL-2.0).

## Thanks 👍
Thank you for viewing my project.