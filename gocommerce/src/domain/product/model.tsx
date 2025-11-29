export interface Product {
    id: string;
    brand: string;
    name: string;
    description: string;
    code: string;
    price: number;
    category: Category;
    image_url: string;
    stock: number;
}

export interface Category {
    id: string;
    name: string;
}