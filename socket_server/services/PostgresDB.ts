import { Pool } from 'pg'

const pool = new Pool({
	user: process.env.POSTGRES_USER || 'postgres',
	host: process.env.PGHOST || 'localhost',
	database: process.env.POSTGRES_DBNAME || 'postgres',
	password: process.env.POSTGRES_PASSWORD || 'password',
	port: parseInt(process.env.PGPORT || '5432', 10),
	max: 20,
	idleTimeoutMillis: 30000,
	connectionTimeoutMillis: 2000,
})

function parseJwt(token: string) {
	return JSON.parse(Buffer.from(token.split('.')[1], 'base64').toString())
}

function InsertIntoDB(
	message: string,
	chatId: string,
	token: string,
	createdAt: string
) {
	const { userId } = parseJwt(token)
	const query = `
		INSERT INTO messages (msg, chat_id, user_id, created_at)
		VALUES ($1, $2, $3, to_timestamp($4))
	`
	return pool.query(query, [message, chatId, userId, createdAt])
}

export default {
	InsertIntoDB,
}
