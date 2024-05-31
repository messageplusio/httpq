# httpq: Webhook as a Service 🚀

## Description 📋
The `httpq` repository is a project that provides a "Webhook as a service" solution. It allows you to create custom HTTP endpoints that can be used to receive and process webhook payloads from various sources.

## Key Features 🌟
- **Webhook Endpoint Creation**: Easily create custom HTTP endpoints to receive webhook payloads.
- **Payload Handling**: Handle and process incoming webhook payloads, allowing you to define custom logic to handle the data.
- **Retries and Persistence**: Supports retrying failed webhook deliveries and persisting the webhook payloads for later processing.
- **Scalability**: Designed to be scalable and can handle a high volume of webhook traffic.

## Getting Started 🛠️

1. **Clone the Repository** 📂: Start by cloning the `httpq` repository to your local machine.
   ```sh
   git clone https://github.com/messageplusio/httpq.git
   cd httpq
   ```
3. **Build the Project 🔨**: Navigate to the project directory and use the provided Makefile to build the project. Run make to build the binary.
   ```sh
   make
   ```
3. **Configure the Environment ⚙️**: Rename the .envrc.example file to .envrc and update the environment variables with your specific configuration.
   ```sh
   mv .envrc.example .envrc
   # Edit .envrc with your configuration```
4. **Run the Service 🚀**: Use the httpq binary to start the webhook service. Run ./httpq to launch the application.
   ```
   ./httpq
   ```
## Usage 📝
Once the httpq service is running, you can create custom HTTP endpoints to receive webhook payloads. The project provides a simple API to interact with the service and manage your webhook configurations.
For more detailed information on usage and customization, please refer to the project's documentation in the GitHub repository.

### Example Configuration 🔧
Here is an example of how you might configure your .envrc file:

```sh
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=username
DB_PASSWORD=password
DB_NAME=httpq
```

### Example API Usage 🌐
To create a new webhook endpoint:

```sh
curl -X POST http://localhost:8080/webhooks \
   -H 'Content-Type: application/json' \
   -d '{
         "name": "example-webhook",
         "url": "http://example.com/webhook-handler",
         "method": "POST"
       }'
```
## Contributing 🤝
We welcome contributions! Please see our contributing guidelines for more details.
