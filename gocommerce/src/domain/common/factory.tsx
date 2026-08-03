import {type ProductService, ProductServiceFactory} from "../product/service.tsx";
import {type ShoppingBasketService, ShoppingBasketServiceFactory} from "../shopping-basket/service.tsx";

export interface ServiceFactory<T> {
    newService(): T;
}

export class ServiceFactoryFactory {
    static PRODUCT_SERVICE_FACTORY: ServiceFactory<ProductService> = ProductServiceFactory.FACTORY;
    static SHOPPING_BASKET_SERVICE_FACTORY: ServiceFactory<ShoppingBasketService> = ShoppingBasketServiceFactory.FACTORY;
}
