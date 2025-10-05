import { Pool } from 'pg';

const pool = new Pool({
    user:
        process.env.PGUSER || 'postgres',
    host:
        process.env.PGHOST || 'localhost',
    database:
        process.env.PGDATABASE || 'postgres',
    password:
        process.env.PGPASSWORD || 'password',
    port:
        parseInt(process.env.PGPORT || '5432', 10),
    max: 20,
    idleTimeoutMillis: 30000,
    connectionTimeoutMillis: 2000,
});

export default pool;