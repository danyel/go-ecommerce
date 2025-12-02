export interface ShoppingBasket {
    id: string;
    items: ShoppingBasketItem[];
    total_price_inclusive: number;
    tax: number;
    total_price_exclusive: number
}

export interface ShoppingBasketItem {
    id: string;
    name: string;
    price: number;
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