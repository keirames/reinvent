import express, { json } from 'express';
import { db } from './src/database';
import { sql } from 'kysely';

const app = express();

async function main() {
    await sql`create table if not exists orc_orders (
        id serial primary key,
        state text check (state in ('pending', 'paid', 'delivered'))
    )`.execute(db);

    app.use(json());

    app.get('/', async (req, res) => {
        res.send('hello');
    });

    app.get('/order', async (req, res) => {
        const v = await db
            .insertInto('orc_orders')
            .values({ state: 'pending' })
            .returning('id')
            .execute();
        console.log('insert order api return', v);

        res.send(v);
    });

    let port = 3000;
    app.listen(port, () => {
        console.log(`Listening on port ${port}`);
    });
}

main().catch((err) => {
    console.error(err);
    console.log('App start failed!');
});
