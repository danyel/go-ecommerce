import type {ShoppingBasket, UpdateShoppingBasketItem} from "./model.tsx";
import type {ServiceFactory} from "../common/factory.tsx";

export interface ShoppingBasketService {
    create(): Promise<ShoppingBasket>;

    findById(shoppingBasketId: string): Promise<ShoppingBasket>;

    update(shoppingBasketId: string, updateShoppingBasketItem: UpdateShoppingBasketItem): Promise<ShoppingBasket>;
}

export class ShoppingBasketServiceFactory implements ServiceFactory<ShoppingBasketService> {
    private static INSTANCE: ShoppingBasketService;
    static FACTORY: ServiceFactory<ShoppingBasketService> = new ShoppingBasketServiceFactory();

    private constructor() {
    }

    newService(): ShoppingBasketService {
        if (!ShoppingBasketServiceFactory.INSTANCE) {
            ShoppingBasketServiceFactory.INSTANCE = new InternalShoppingBasketService();
        }
        return ShoppingBasketServiceFactory.INSTANCE;
    }
}

class InternalShoppingBasketService implements ShoppingBasketService {
    async update(shoppingBasketId: string, updateShoppingBasketItem: UpdateShoppingBasketItem): Promise<ShoppingBasket> {
        return fetch(`/api/shopping-basket/v1/shopping-baskets/${shoppingBasketId}`, {
            method: "PUT",
            body: JSON.stringify(updateShoppingBasketItem),
            headers: {
                'Accept-Language': 'en',
                'Content-Type': 'application/json',
                Accept: 'application/json',
            }
        })
            .then(resp => {
                if (!resp.ok) {
                    return Promise.reject('Could not add item to shopping basket');
                }
                return resp.json();
            })
    }

    async create(): Promise<ShoppingBasket> {
        return fetch('/api/shopping-basket/v1/shopping-baskets', {
            method: "POST",
            headers: {
                'Accept-Language': 'en',
                'Content-Type': 'application/json',
                Accept: 'application/json',
            }
        })
            .then(resp => {
                if (!resp.ok) {
                    return Promise.reject('Could not create shopping basket');
                }
                return resp.json();
            })
    }

    async findById(shoppingBasketId: string): Promise<ShoppingBasket> {
        return fetch(`/api/shopping-basket/v1/shopping-baskets/${shoppingBasketId}`, {
            method: "GET",
            headers: {
                'Accept-Language': 'en',
                'Content-Type': 'application/json',
                Accept: 'application/json',
            }
        })
            .then(resp => {
                if (!resp.ok) {
                    throw Error('Could not fetch the shopping basket');
                }
                return resp.json();
            })
    }

}