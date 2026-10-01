export interface ShoppingBasket {
    id: string;
    items: ShoppingBasketItem[];
    total_price: Price;
}

export interface ShoppingBasketItem {
    name: string;
    base_price: Price;
    total_price: Price;
    image_url: string;
    quantity: number;
    product_id: string;
}

export interface ShoppingBasketId {
    id: string;
}

export interface UpdateShoppingBasketItem {
    product_id: string;
    quantity: number;
}

export interface Price {
    inclusive: number;
    tax: number;
    exclusive: number;
}