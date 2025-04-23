# FraudBlocker

FraudBlocker is a Go-based phone number blocking system designed to help identify and block fraudulent phone numbers. The system provides a web interface to check and manage blocked phone numbers.

## Project Structure

- `config/` - Configuration files and settings
- `db/` - Database related code and migrations
- `handlers/` - HTTP handlers for the web server
- `models/` - Data models used in the application
- `static/` - Static assets like CSS and JavaScript files
- `templates/` - HTML templates for rendering web pages
- `utils/` - Utility functions

## Configuration

The application uses a `config.yaml` file for database and server configuration. Update this file with your database credentials and server port.

Example:

```yaml
database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: postgres
  sslmode: disable

server:
  port: 8080
```

## Building and Running

Make sure you have Go installed on your system.

To build the application:

```bash
go build -o fraudblocker main.go
```

To run the application:

```bash
./fraudblocker
```

The server will start on the port specified in `config.yaml`.

## Usage

Open your browser and navigate to `http://localhost:<port>` (default is 8080) to access the web interface. You can check phone numbers and manage the blocked list.

## Database

The application uses a PostgreSQL database. The schema is managed via SQL migrations located in `db/migrations/`.

## License

