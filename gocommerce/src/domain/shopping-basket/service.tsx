import type {ShoppingBasket, ShoppingBasketId, UpdateShoppingBasketItem} from "./model.tsx";
import type {ServiceFactory} from "../common/factory.tsx";

export interface ShoppingBasketService {
    create(): Promise<ShoppingBasketId>;

    findById(shoppingBasketId: ShoppingBasketId): Promise<ShoppingBasket>;

    update(shoppingBasketId: ShoppingBasketId, updateShoppingBasketItem: UpdateShoppingBasketItem): Promise<ShoppingBasket>;
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
    async update(shoppingBasketId: ShoppingBasketId, updateShoppingBasketItem: UpdateShoppingBasketItem): Promise<ShoppingBasket> {
        return fetch(`/api/shopping-basket/v1/shopping-baskets/${shoppingBasketId.id}`, {
            method: "POST",
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

    async create(): Promise<ShoppingBasketId> {
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

    async findById(shoppingBasketId: ShoppingBasketId): Promise<ShoppingBasket> {
        return fetch(`/api/shopping-basket/v1/shopping-baskets/${shoppingBasketId.id}`, {
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