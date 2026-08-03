export class Optional<T> {
    private value?: T;
    private ifPresentFunction: (value: T) => void = () => {
    };
    private orElseFunction: () => T = () => {
        return {} as T;
    };

    constructor(value?: T) {
        this.value = value;
    }

    static of<T>(value?: T): Optional<T> {
        return new Optional<T>(value);
    }

    ifPresent(param: (value: T) => void): Optional<T> {
        this.ifPresentFunction = param;
        return this;
    }

    orElse(orElseFunction: () => T): Optional<T> {
        this.orElseFunction = orElseFunction;
        return this;
    }

    execute() {
        if (!this.value) {
            this.value = this.orElseFunction();
        }

        if (this.value) {
            this.ifPresentFunction(this.value)
        }
    }
}