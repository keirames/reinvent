// promise can only in one of three state
const states = {
    PENDING: 'pending',
    FULFILLED: 'fulfilled',
    REJECTED: 'rejected',
};

const isThenable = (maybePromise: any) =>
    maybePromise && typeof maybePromise.then === 'function';

export class MyPromise {
    private state: any;
    private value: any;
    private reason: any;
    private thenQueue: any[];
    private finallyQueue: any[];

    constructor(computation?: any) {
        this.state = states.PENDING;

        this.value = undefined;
        this.reason = undefined;

        this.thenQueue = [];
        this.finallyQueue = [];

        if (typeof computation === 'function') {
            setTimeout(() => {
                try {
                    computation(
                        this.onFulfilled.bind(this),
                        this.onRejected.bind(this),
                    );
                } catch (err) {
                    // comeback later
                }
            });
        }
    }

    // fulfilledFunction is a function that transform the value
    // if promise is fulfilled.
    // catchFunction is transform a reason if promise is rejected.
    then(fulfilledFunction?: any, catchFunction?: any) {
        // create a new promise in pending state
        // call this promise "controlled" because we is the one control
        // when this promise change state, not user's computation function
        const controlledPromise = new MyPromise();
        this.thenQueue.push([
            controlledPromise,
            fulfilledFunction,
            catchFunction,
        ]);

        // the user call .then after promise is settled
        if (this.state === states.FULFILLED) {
            this.propagateFulfilled();
        } else if (this.state === states.REJECTED) {
            this.propagateRejected();
        }

        return controlledPromise;
    }

    catch(catchFunction: any) {
        return this.then(undefined, catchFunction);
    }

    // only run when parent promise resolved or rejected
    // sideEffectFunction don't need any params
    finally(sideEffectFunction: any) {
        if (this.state !== states.PENDING) {
            sideEffectFunction();

            if (this.state === states.FULFILLED) {
                return new MyPromise((resolve: any) => resolve(this.value));
            } else {
                return new MyPromise((_: any, reject: any) =>
                    reject(this.reason),
                );
            }
        }

        // if a promise still pending
        const controlledPromise = new MyPromise();
        this.finallyQueue.push([controlledPromise, sideEffectFunction]);

        return controlledPromise;
    }

    private propagateFulfilled() {
        for (const item of this.thenQueue) {
            const [controlledPromise, fulfilledFunction, catchFunction] = item;

            if (typeof fulfilledFunction === undefined) {
                // specification declare if fulfilledFunction is undefined
                // we simply pass current value to it
                controlledPromise.onFulfilled(this.value);
                continue;
            }

            // .then(v => v + 1) || .then(v => new Promise())
            const valueOrPromise = fulfilledFunction(this.value);

            if (isThenable(valueOrPromise)) {
                // a promise
                valueOrPromise.then(
                    (value: any) => controlledPromise.onFulfilled(value),
                    (reason: any) => controlledPromise.onRejected(reason),
                );
            } else {
                // a single value
                controlledPromise.onFulfilled(valueOrPromise);
            }
        }

        this.finallyQueue.forEach(([controlledPromise, sideEffectFunction]) => {
            sideEffectFunction();
            controlledPromise.onFulfilled(this.value);
        });

        // clear all processed promise
        this.thenQueue = [];
        this.finallyQueue = [];
    }

    // catchFunction used to recovery, it converse throw error into a normal promise
    private propagateRejected() {
        for (const item of this.thenQueue) {
            const [controlledPromise, _, catchFunction] = item;

            if (typeof catchFunction === 'function') {
                const valueOrPromise = catchFunction(this.reason);

                if (isThenable(valueOrPromise)) {
                    valueOrPromise.then(
                        (value: any) => controlledPromise.onFulfilled(value),
                        (reason: any) => controlledPromise.onRejected(reason),
                    );
                } else {
                    controlledPromise.onFulfilled(valueOrPromise);
                }
            } else {
                controlledPromise.onRejected(this.reason);
            }
        }

        this.finallyQueue.forEach(([controlledPromise, sideEffectFunction]) => {
            sideEffectFunction();
            controlledPromise.onRejected(this.reason);
        });

        this.thenQueue = [];
        this.finallyQueue = [];
    }

    // is resolve function
    private onFulfilled(value: any) {
        if (this.state === states.PENDING) {
            this.state = states.FULFILLED;
            this.value = value;

            this.propagateFulfilled();
        }
    }

    // is reject function
    private onRejected(reason: any) {
        if (this.state === states.PENDING) {
            this.state = states.REJECTED;
            this.reason = reason;

            this.propagateRejected();
        }
    }
}
