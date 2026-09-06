import type {Product, ProductDTO} from './model.tsx';

export default class ProductMapper {
    static map(productDto: ProductDTO): Product {
        return {
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
    }
}
