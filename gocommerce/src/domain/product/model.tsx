export interface ProductDTO {
  id: string;
  brand: string;
  name: string;
  description: string;
  code: string;
  price: Price;
  category: Category;
  image_url: string;
  stock: number;
}

export interface Product {
  id: string;
  brand: string;
  name: string;
  description: string;
  code: string;
  price: Price;
  category: Category;
  imageUrl: string;
  stock: number;
}

export interface Price {
  inclusive: number;
  exclusive: number;
  tax: number;
  currency: string;
}

export interface Category {
  id: string;
  name: string;
}

export interface CreateProduct {
  brand: string;
  name: string;
  description: string;
  code: string;
  price: number;
  tax_percentage: number;
  currency: string;
  category_id: string
  image_url: string;
}
