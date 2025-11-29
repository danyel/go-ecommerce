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
    amount: number;
    product_id: string;
}

export interface ShoppingBasketId {
    id: string;
}
export interface AddItem {
    product_id: string;
}