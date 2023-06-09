import { P } from './src/p';

function main() {
  const promise = new Promise((resolve, reject) => {
    resolve(1);
    reject('hja');
  });

  // promise.then((v) => console.log(v));
  // promise.then((v) => console.log('next then', v));

  const p = new P((resolve, reject) => {
    resolve(1);
    // setTimeout(() => {
    //   reject('haha');
    // }, 3000);
  }, 1);

  p.then((v) => console.log('my', v));
  // p.catch((e) => console.log('rejected', e));
  // p.then((v) => console.log('my second then', v));

  // p.then((value) => { console.log("then 1", value); return 2 }).then((value) => {console.log(value);})
}

main();
