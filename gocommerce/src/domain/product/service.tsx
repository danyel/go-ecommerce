import type { Product, ProductDTO } from "./model.tsx";
import type { ServiceFactory } from "../common/factory.tsx";

export interface ProductService {
  fetchAll(): Promise<Product[]>;
}

export class ProductServiceFactory implements ServiceFactory<ProductService> {
  private static INSTANCE: ProductService;
  static FACTORY: ServiceFactory<ProductService> = new ProductServiceFactory();

  private constructor() {
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

    return await res.json()
      .then((productDto: ProductDTO[]) => {
        return productDto.map(productDto => {
          const product: Product = {
            id: productDto.id,
            brand: productDto.brand,
            name: productDto.name,
            description: productDto.description,
            code: productDto.code,
            price: {
              inclusive: productDto.price.inclusive,
              exclusive: productDto.price.exclusive,
              currency: productDto.price.currency,
              tax: productDto.price.tax
            },
            category: {
              id: productDto.category.id,
              name: productDto.category.name
            },
            imageUrl: productDto.image_url,
            stock: productDto.stock
          } as Product;
          console.log('Product', product);
          return product;
        })
      });
  }

}
