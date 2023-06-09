type OnSuccess = (v: number) => void;

type OnFail = (reason?: any) => void;

type State = 'pending' | 'fulfilled' | 'rejected';

export class P {
  private no: number = 0;
  private thenSubscribers: ((v: unknown) => void)[] = [];
  private catchSubscribers: ((reason?: any) => void)[] = [];
  private state: State = 'pending';
  private value: number = 0;

  constructor(cb: (resolve: OnSuccess, reject: OnFail) => void) {
    try {
      // This make onSuccess and onFail lost it's this keyword
      // need manually bind func or could use arrow function
      cb(this.onSuccess, this.onFail);
    } catch (e) {
      console.log(e);
      this.onFail();
    }
  }

  private notify() {
    if (this.state === 'fulfilled') {
      for (const s of this.thenSubscribers) {
        s(this.value);
      }

      this.thenSubscribers = [];
    }

    if (this.state === 'rejected') {
      for (const s of this.catchSubscribers) {
        s(this.value);
      }

      this.catchSubscribers = [];
    }
  }

  onSuccess = (v: number | any) => {
    queueMicrotask(() => {
      if (this.state !== 'pending') {
        return;
      }

      if (v instanceof P) {
        v.then(this.onSuccess, this.onFail);
        return;
      }

      this.state = 'fulfilled';
      this.value = v;

      this.notify();
    });
  };

  onFail: OnFail = (v) => {
    queueMicrotask(() => {
      if (this.state !== 'pending') {
        return;
      }

      if (v instanceof P) {
        v.then(this.onSuccess, this.onFail);
        return;
      }

      this.state = 'rejected';
      this.value = v;

      this.notify();
    });
  };

  then(thenCb?: any, catchCb?: (reason?: any) => any) {
    return new P((resolve, reject) => {
      this.thenSubscribers.push((result) => {
        if (thenCb === undefined) {
          resolve(result as number);
          return;
        }

        try {
          resolve(thenCb(result));
        } catch (e) {
          reject(e);
        }
      });

      this.catchSubscribers.push((result) => {
        if (catchCb === undefined) {
          reject(result as number);
          return;
        }

        try {
          resolve(catchCb(result));
        } catch (e) {
          reject(e);
        }
      });

      // Why ?
      // If promise resolved immediately, call resolve func instantly
      // onSuccess will run before .then func.
      this.notify();
    });
  }

  catch(cb: (reason?: any) => void) {
    this.then(undefined, cb);
  }

  finally = (cb: any) => {
    return this.then(
      () => {
        cb();
        return;
      },
      (result) => {
        cb();
        throw result;
      },
    );
  };

  static resolve(v: any) {
    return new P((resolve, reject) => {
      resolve(v);
    });
  }

  static reject(v: any) {
    return new P((resolve, reject) => {
      reject(v);
    });
  }

  static all(promises: P[]) {
    const result: any = [];
    let completedPromises = 0;

    return new P((resolve, reject) => {
      for (let i = 0; i < promises.length; i++) {
        const promise = promises[i];
        promise
          .then((v: any) => {
            completedPromises++;
            result[i] = v;

            if (completedPromises === promises.length) {
              resolve(result);
            }
          })
          .then(reject);
      }
    });
  }

  static allSettled(promises: P[]) {
    const result: any = [];
    let completedPromises = 0;

    return new P((resolve, reject) => {
      for (let i = 0; i < promises.length; i++) {
        const promise = promises[i];
        promise
          .then((v: any) => {
            result[i] = { status: 'fulfilled', value: v };
          })
          .then((r: any) => {
            result[i] = { status: 'rejected', value: r };
          })
          .finally(() => {
            completedPromises++;

            if (completedPromises === promises.length) {
              resolve(result);
            }
          });
      }
    });
  }

  static race(promises: P[]) {
    return new P((resolve, reject) => {
      promises.forEach((p) => p.then(resolve).catch(reject));
    });
  }
}
