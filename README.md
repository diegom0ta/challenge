# Challenge

This project includes a Go application with PostgreSQL database and Flyway migrations.

## Database Setup

The project uses Docker Compose to set up PostgreSQL and run Flyway migrations.

### Prerequisites

- Docker and Docker Compose installed
- Make (optional, for convenience commands)

### Getting Started

1. **Start the database and run migrations:**
   ```bash
   make up
   ```
   Or manually:
   ```bash
   docker-compose up -d
   ```

2. **View logs:**
   ```bash
   make logs
   ```

3. **Connect to PostgreSQL:**
   ```bash
   make psql
   ```

4. **Stop services:**
   ```bash
   make down
   ```

5. **Clean up (removes volumes):**
   ```bash
   make clean
   ```

### Database Schema

The B3 table structure:
- `id`: Serial primary key
- `data_negocio`: Timestamp of the trading date
- `codigo_instrumento`: Instrument code (VARCHAR 50)
- `preco_negocio`: Trading price (DECIMAL 15,2)
- `quantidade_negociada`: Quantity traded (INTEGER)
- `hora_fechamento`: Closing time (TIMESTAMP)
- `created_at`: Record creation timestamp
- `updated_at`: Record update timestamp

### Environment Variables

Copy `.env.example` to `.env` and modify as needed:
```bash
cp .env.example .env
```

### Migration Files

Migration files are located in the `migrations/` directory and follow Flyway naming convention:
- `V1__Create_b3_table.sql` - Initial table creation
