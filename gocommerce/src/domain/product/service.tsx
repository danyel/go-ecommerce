import type {Product} from "./model.tsx";
import type {ServiceFactory} from "../common/factory.tsx";

export interface ProductService {
    fetchAll(): Promise<Product[]>;
}

export class ProductServiceFactory implements ServiceFactory<ProductService> {
    private static INSTANCE: ProductService;
    static FACTORY: ServiceFactory<ProductService> = new ProductServiceFactory();

private    constructor() {
    }

    newService(): ProductService {
        if (!ProductServiceFactory.INSTANCE) {
            ProductServiceFactory.INSTANCE = new InternalProductService();
        }
        return ProductServiceFactory.INSTANCE;
    }
}

class InternalProductService implements ProductService {
    async fetchAll(): Promise<Product[]> {
        const res = await fetch("/api/product/v1/products", {
            method: "GET",
            headers: {
                'Accept-Language': 'en',
                'Content-Type': 'application/json',
                Accept: 'application/json',
            }
        });

        if (!res.ok) {
            // noinspection ExceptionCaughtLocallyJS
            return Promise.reject('Error while fetching the products');
        }

        return await res.json();
    }

}