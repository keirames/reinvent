import { Kafka } from 'kafkajs';
import { db } from './src/database';

export const Topic = {
    BillOrder: 'bill_order',
    ShipOrder: 'ship_order',
    OrderPlaced: 'order_placed',
    OrderPaid: 'order_paid',
    OrderDelivered: 'order_delivered',
    OrderRejected: 'order_rejected',
};

export async function waitFor(ms: number) {
    return new Promise((resolve) => {
        setTimeout(() => {
            resolve(1);
        }, ms);
    });
}

export function getKafka(): Kafka {
    return new Kafka({
        brokers: ['localhost:9094'],
    });
}

export async function listenToBillOrderEvent() {
    const kafka = getKafka();
    const consumer = kafka.consumer({ groupId: 'payment' });
    await consumer.connect();
    await consumer.subscribe({
        topic: Topic.BillOrder,
    });
    await consumer.run({
        autoCommit: true,
        eachMessage: async ({
            topic,
            partition,
            message,
            heartbeat,
            pause,
        }) => {
            console.log('processing bill order...');
            await waitFor(2000);
            console.log({
                key: message?.key?.toString(),
                value: message?.value?.toString(),
                headers: message.headers,
            });
            // await consumer.commitOffsets([
            //     { topic: Topic.BillOrder, offset: message.offset, partition },
            // ]);

            const randNum = Math.floor(Math.random() * 100);
            if (randNum % 2 === 0) {
                console.log('pay order successfully!');
                console.log('produce OrderPaid event from Payment Service');

                const producer = kafka.producer();
                await producer.connect();
                await producer.send({
                    topic: Topic.OrderPaid,
                    messages: [{ value: message.value }],
                });
            } else {
                console.log('pay order fail, no money left!');
                console.log('produce OrderRejected event from Payment Service');

                const producer = kafka.producer();
                await producer.connect();
                await producer.send({
                    topic: Topic.OrderRejected,
                    messages: [{ value: message.value }],
                });
            }
        },
    });
}

export async function listenToOrderPlacedEvent() {
    const kafka = getKafka();
    const consumer = kafka.consumer({ groupId: 'orchestrator' });
    await consumer.connect();
    await consumer.subscribe({
        topic: Topic.OrderPlaced,
    });
    await consumer.run({
        autoCommit: true,
        eachMessage: async ({
            topic,
            partition,
            message,
            heartbeat,
            pause,
        }) => {
            console.log({
                key: message?.key?.toString(),
                value: message?.value?.toString(),
                headers: message.headers,
            });
            // await consumer.commitOffsets([
            //     { topic: Topic.OrderPlaced, offset: message.offset, partition },
            // ]);

            const producer = kafka.producer();
            await producer.connect();
            await producer.send({
                topic: Topic.BillOrder,
                messages: [{ value: message.value }],
            });
        },
    });
}

export async function listenToOrderPaidEvent() {
    const kafka = getKafka();
    const consumer = kafka.consumer({ groupId: 'orchestrator-order-paid' });
    await consumer.connect();
    await consumer.subscribe({
        topic: Topic.OrderPaid,
    });
    console.log('orchestrator subscribe to OrderPaid events');
    await consumer.run({
        autoCommit: true,
        eachMessage: async ({
            topic,
            partition,
            message,
            heartbeat,
            pause,
        }) => {
            console.log('listen to OrderPaid event got:');
            console.log({
                key: message?.key?.toString(),
                value: message?.value?.toString(),
                headers: message.headers,
            });

            await db
                .updateTable('orc_orders')
                .set({
                    state: 'paid',
                })
                .where('orc_orders.id', '=', Number(message!.value!.toString()))
                .execute();

            // await consumer.commitOffsets([
            //     { topic, offset: message.offset, partition },
            // ]);
        },
    });
}

export async function listenToOrderRejectedEvent() {
    const kafka = getKafka();
    const consumer = kafka.consumer({ groupId: 'orchestrator-order-rejected' });
    await consumer.connect();
    await consumer.subscribe({
        topic: Topic.OrderRejected,
    });
    console.log('orchestrator subscribe to OrderRejected events');
    await consumer.run({
        autoCommit: true,
        eachMessage: async ({
            topic,
            partition,
            message,
            heartbeat,
            pause,
        }) => {
            console.log('listen to OrderRejected event got:');
            console.log({
                key: message?.key?.toString(),
                value: message?.value?.toString(),
                headers: message.headers,
            });

            console.log('send command to Order Service update rejected order');
            await db
                .updateTable('orc_orders')
                .set({
                    state: 'cancelled',
                })
                .where('orc_orders.id', '=', Number(message!.value!.toString()))
                .execute();

            // await consumer.commitOffsets([
            //     {
            //         topic,
            //         offset: message.offset,
            //         partition,
            //     },
            // ]);
            console.log('committed OrderRejected event');
        },
    });
}
