import { Pool } from 'pg';
import { Generated, Kysely, PostgresDialect } from 'kysely';

interface OrcOrdersTable {
    id: Generated<number>;
    state: 'pending' | 'paid' | 'delivered';
}

interface Database {
    orc_orders: OrcOrdersTable;
}

const dialect = new PostgresDialect({
    pool: new Pool({
        database: 'postgres',
        host: 'localhost',
        user: 'postgres',
        password: 'password',
        port: 5432,
        max: 10,
    }),
});

export const db = new Kysely<Database>({
    dialect,
});
