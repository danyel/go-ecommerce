import {useEffect, useState} from 'react';
import type {Product, ProductDTO} from '../../domain/product/model.tsx';
import {ChevronDown} from 'lucide-react';
import {useGlobalState} from '../../state/global-state.tsx';
import type {ShoppingBasket, UpdateShoppingBasketItem} from '../../domain/shopping-basket/model.tsx';
import Cookies from 'js-cookie';
import ApiClient from "../../domain/common/api-client.tsx";
import ProductMapper from "../../domain/product/mapper.tsx";

const ProductsPage = () => {
    const [products, setProducts] = useState<Product[]>([]);
    const [singleFetch, setSingleFetch] = useState<boolean>(false);
    const globalStateType = useGlobalState();
    const addToCart = async (product: Product) => {
        const updateShoppingBasketItem: UpdateShoppingBasketItem = {product_id: product.id, quantity: 1};
        let shoppingBasketId = globalStateType.shoppingBasket.id;
        if (!shoppingBasketId) {
            const newShoppingBasket = await ApiClient.POST<ShoppingBasket, ShoppingBasket>('/api/shopping-basket/v1/shopping-baskets', undefined);
            shoppingBasketId = newShoppingBasket.id;
        } else {
            const found = globalStateType.shoppingBasket.items?.find(e => e.product_id == product.id);
            if (found) {
                updateShoppingBasketItem.quantity += found.quantity;
            }
        }
        const updatedShoppingBasket = await ApiClient.PUT<ShoppingBasket, UpdateShoppingBasketItem>(`/api/shopping-basket/v1/shopping-baskets/${shoppingBasketId}`, updateShoppingBasketItem);
        Cookies.set('shopping_basket_id', shoppingBasketId);
        globalStateType.setShoppingBasket(updatedShoppingBasket);
        const freshProductData = await ApiClient.GET<Product>(`/api/product/v1/products/${product.id}`);
        setProducts(prevState => prevState.map(p => p.id === freshProductData.id ? {
            ...p,
            stock: freshProductData.stock
        } : p));
    };
    useEffect(() => {
        if (products.length == 0 && !singleFetch) {
            ApiClient.GET<ProductDTO[]>('/api/product/v1/products')
                .then(data => {
                    setProducts(data.map(v => ProductMapper.map(v)));
                    setSingleFetch(true);
                });
        }
    }, [products, singleFetch]);
    return (
        <main className='flex-1'>
            <div
                className='bg-white rounded-lg shadow-sm p-4 mb-6 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4'>
                <p className='text-gray-600'>
                    Showing <span className='font-semibold'>{products.length}</span> products
                </p>
                <div className='flex items-center gap-2'>
                    <label className='text-sm text-gray-600'>Sort by:</label>
                    <button className='flex items-center gap-2 px-4 py-2 border rounded-lg hover:bg-gray-50'>
                        Best Match <ChevronDown size={16}/>
                    </button>
                </div>
            </div>
            <div className='grid grid-cols-4 sm:grid-cols-2 xl:grid-cols-3 gap-6'>
                {products.map((product) => (
                    <div
                        key={product.id}
                        className='bg-white rounded-lg shadow-sm hover:shadow-md transition overflow-hidden'
                    >
                        <div className='relative'>
                            <img
                                src={product.imageUrl}
                                alt={product.name}
                                className='w-full h-48 object-cover'
                            />
                            {product.stock === 0 && (
                                <div className='absolute top-2 right-2 bg-red-500 text-white text-xs px-2 py-1 rounded'>
                                    Out of Stock
                                </div>
                            )}
                        </div>
                        <div className='p-4'>
                            <p className='text-xs text-gray-500 mb-1'>{product.category.name}</p>
                            <h3 className='font-semibold mb-2'>{product.name}</h3>
                            <div className='flex items-center justify-between'>
                                <p className='text-2xl font-bold text-blue-600'>
                                    {product.price.inclusive} €
                                </p>
                                <p className='text-xs font-bold text-green-300'>
                                    {product.stock} piece(s) left.
                                </p>
                                <button
                                    onClick={() => addToCart(product)}
                                    disabled={product.stock === 0}
                                    className={`px-4 py-2 rounded-lg font-medium transition ${product.stock !== 0
                                        ? 'bg-blue-600 text-white hover:bg-blue-700'
                                        : 'bg-gray-200 text-gray-500 cursor-not-allowed'
                                    }`}
                                >
                                    {product.stock !== 0 ? 'Add to Cart' : 'Unavailable'}
                                </button>
                            </div>
                        </div>
                    </div>
                ))}
            </div>
        </main>
    );
};

export default ProductsPage;
