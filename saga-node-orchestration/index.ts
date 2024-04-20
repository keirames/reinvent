import express, { json } from 'express';
import { db } from './src/database';
import { sql } from 'kysely';
import { Kafka } from 'kafkajs';
import {
    Topic,
    listenToBillOrderEvent,
    listenToOrderPaidEvent,
    listenToOrderPlacedEvent,
    listenToOrderRejectedEvent,
    waitFor,
} from './utils';

const app = express();

async function initTopicUntilSuccess(kafka: Kafka) {
    console.log('init kafka topic');
    for (let i = 0; i < 50; i++) {
        const admin = kafka.admin();
        try {
            await admin.connect();
            console.log('admin connect successfully');

            try {
                await admin.createTopics({
                    topics: [
                        { topic: Topic.OrderPlaced, numPartitions: 5 },
                        { topic: Topic.BillOrder, numPartitions: 5 },
                        { topic: Topic.OrderPaid, numPartitions: 5 },
                        { topic: Topic.ShipOrder, numPartitions: 5 },
                        { topic: Topic.OrderDelivered, numPartitions: 5 },
                        { topic: Topic.OrderRejected, numPartitions: 5 },
                    ],
                });
            } catch (err) {
                console.error(err);
                console.log('Topics already exists');
                return;
            }

            await admin.disconnect();
            return;
        } catch (err) {
            console.error('admin error', err);
        }

        console.log('retrying init kafka topic...');
        await waitFor(1000);
    }

    throw new Error('fail to init kafka topic');
}

async function main() {
    const kafka = new Kafka({
        brokers: ['localhost:9094'],
    });
    await initTopicUntilSuccess(kafka);

    await listenToOrderPlacedEvent();

    await listenToBillOrderEvent();

    await listenToOrderPaidEvent();

    await listenToOrderRejectedEvent();

    await sql`create table if not exists orc_orders (
        id serial primary key,
        state text check (state in ('pending', 'paid', 'delivered', 'cancelled'))
    )`.execute(db);

    app.use(json());

    app.get('/', async (req, res) => {
        res.send('hello');
    });

    app.get('/orders', async (req, res) => {
        const orders = await db.selectFrom('orc_orders').selectAll().execute();
        res.send(orders);
    });

    app.get('/order', async (req, res) => {
        const v = await db
            .insertInto('orc_orders')
            .values({ state: 'pending' })
            .returning('id')
            .execute();
        console.log('insert order api return', v);

        console.log('produce order created event for orchestrator');
        const producer = kafka.producer();
        await producer.connect();
        await producer.send({
            topic: Topic.OrderPlaced,
            messages: [{ key: 'key1', value: String(v[0].id) }],
        });
        console.log('successfully');

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
