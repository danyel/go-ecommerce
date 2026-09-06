export interface ServiceFactory<T> {
    get(): T;
}

export class ServiceFactoryFactory {
}
