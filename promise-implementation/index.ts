import { MyPromise } from './src/my-promise';

function main() {
    const promise = new MyPromise((resolve: any, reject: any) => {
        setTimeout(() => {
            reject(69);
        }, 1000);
    }).catch((r: any) => {
        console.log(`throw error=${r}`);
        return 'recovered';
    });

    const firstThen = promise.then((value: any) => {
        console.log(`got value=${value}`);
        return value + 1;
    });

    const secondThen = promise.then((value: any) => {
        console.log(`got value=${value}`);
        return value + 1;
    });
}

main();
