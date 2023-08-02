export class Queue {
    constructor() {
        this._queue = [];
        this._front = 0;
    }

    enqueue(data) {
        this._queue.push(data);
    }

    dequeue() {
        if (this._front === this._queue.length) {
            return null;
        }
        const data = this._queue[this._front];
        this._front++;
        return data;
    }
}
