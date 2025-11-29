import {useEffect, useState} from "react";
import type {Product} from "../../domain/product/model.tsx";

const ProductsPage = () => {
    const [products, setProducts] = useState<Product[]>([]);
    useEffect(() => {
        fetch("/api/product/v1/products")
            .then(res => {
                if (!res.ok) {
                    throw Error('Error while fetching the products');
                }

                return res.json();
            })
            .then(data => setProducts(data));
    }, [])
    return (
        <table>
            <thead>
            <tr>
                <th>ID</th>
                <th>Name</th>
                <th>Description</th>
                <th>Price (€)</th>
            </tr>
            </thead>
           <tbody>
           {products.map(product => {
               return (
                   <tr key={product.id}>
                       <td>{product.id}</td>
                       <td>{product.name}</td>
                       <td>{product.description}</td>
                       <td>{product.price}</td>
                   </tr>
               )
           })}
           </tbody>
        </table>
    );
};

export default ProductsPage;