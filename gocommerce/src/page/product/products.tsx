import {useEffect, useMemo, useState} from 'react';
import {useSearchParams} from 'react-router-dom';
import type {Product, ProductDTO} from '../../domain/product/model.tsx';
import {ChevronDown} from 'lucide-react';
import {useGlobalState} from '../../state/global-state.tsx';
import type {
    ShoppingBasket,
    ShoppingBasketItem,
    UpdateShoppingBasketItem
} from '../../domain/shopping-basket/model.tsx';
import Cookies from 'js-cookie';
import ApiClient from "../../domain/common/api-client.tsx";
import ProductMapper from "../../domain/product/mapper.tsx";
import type {Category} from '../../domain/product/model.tsx';
import {ApiError} from '../../domain/common/api-client.tsx';
import type {Page} from "../../domain/common/page.tsx";

const ProductsPage = () => {
    const [products, setProducts] = useState<Product[]>([]);
    const [singleFetch, setSingleFetch] = useState<boolean>(false);
    const [categories, setCategories] = useState<Category[]>([]);
    const [page] = useState(1);
    const [searchParams, setSearchParams] = useSearchParams();
    const selectedCategory = searchParams.get('category') ?? '';
    const globalStateType = useGlobalState();
    const visibleProducts = useMemo(() => selectedCategory ? products.filter(product => product.category.id === selectedCategory) : products, [products, selectedCategory]);
    const addToCart = async (product: Product) => {
        const updateShoppingBasketItem: UpdateShoppingBasketItem = {product_id: product.id, quantity: 1};
        let shoppingBasketId = globalStateType.shoppingBasket.id;
        if (shoppingBasketId) {
            const found = globalStateType.shoppingBasket.items?.find(e => e.product_id == product.id);
            if (found) {
                updateShoppingBasketItem.quantity = found.quantity + 1;
            }
        } else {
            try {
                const newShoppingBasket = await ApiClient.POST<ShoppingBasket, ShoppingBasket>('/api/shopping-basket/v1/shopping-baskets', undefined);
                shoppingBasketId = newShoppingBasket.id;
            } catch (error) {
                if (error instanceof ApiError && error.status === 401) window.location.assign('/login');
                return;
            }
        }
        console.log('Current shopping basket:', globalStateType.shoppingBasket);
        let updatedShoppingBasket: ShoppingBasket;
        try {
            updatedShoppingBasket = await ApiClient.PUT<ShoppingBasket, UpdateShoppingBasketItem>(`/api/shopping-basket/v1/shopping-baskets/${shoppingBasketId}`, updateShoppingBasketItem);
        } catch (error) {
            if (error instanceof ApiError && error.status === 401) window.location.assign('/login');
            return;
        }
        console.log('updatedShoppingBasket', updatedShoppingBasket);
        Cookies.set('shopping_basket_id', shoppingBasketId);
        globalStateType.setShoppingBasket(updatedShoppingBasket);
        const freshProductData = await ApiClient.GET<Product>(`/api/product/v1/products/${product.id}`);
        setProducts(prevState => prevState.map(p => p.id === freshProductData.id ? {
            ...p,
            stock: freshProductData.stock
        } : p));
    };
    useEffect(() => {
        if (!singleFetch) {
            ApiClient.GET<Page<ProductDTO>>(`/api/product/v1/products?page=${page}&page_size=24`)
                .then(data => {
                    setProducts(data.items.map(v => ProductMapper.map(v)));
                    setSingleFetch(true);
                });
        }
    }, [page, singleFetch]);
    useEffect(() => {
        ApiClient.GET<Page<Category>>('/api/management/v1/categories?page=1&page_size=50').then(data => setCategories(data.items)).catch(() => setCategories([]));
    }, []);
    useEffect(() => {
        globalStateType.shoppingBasket?.items?.forEach((shoppingBasketItem: ShoppingBasketItem) =>
            setProducts((prevState: Product[]) =>
                prevState.map((product: Product) => {
                    if (product.id === shoppingBasketItem.id) {
                        return {
                            ...product,
                            stock: shoppingBasketItem.remaining
                        };
                    }
                    return product;
                })
            )
        );
    }, [globalStateType.shoppingBasket]);

    return (
        <main className='flex-1'>
            <div
                className='bg-white rounded-lg shadow-sm p-4 mb-6 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4'>
                <p className='text-gray-600'>
                    Showing <span className='font-semibold'>{visibleProducts.length}</span> products
                </p>
                <div className='flex items-center gap-2'>
                        <label className='text-sm text-gray-600' htmlFor='category-filter'>Category:</label>
                        <select id='category-filter' value={selectedCategory} onChange={event => {
                            if (event.target.value) setSearchParams({category: event.target.value});
                            else setSearchParams({});
                        }} className='border rounded-lg px-3 py-2'>
                            <option value=''>All categories</option>
                            {categories.map(category => <option key={category.id} value={category.id}>{category.name}</option>)}
                        </select>
                        <ChevronDown size={16}/>
                    </div>
            </div>
            <div className='grid grid-cols-4 sm:grid-cols-2 xl:grid-cols-3 gap-6'>
                {visibleProducts.map((product) => (
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
                                <div
                                    className='absolute top-2 right-2 bg-red-500 text-white text-xs px-2 py-1 rounded'>
                                    Out of Stock
                                </div>
                            )}
                        </div>
                        <div className='p-4'>
                            <p className='text-xs text-gray-500 mb-1'>{product.category.name}</p>
                            <h3 className='font-semibold mb-2'>{product.brand}</h3>
                            <div className='flex items-center justify-between'>{product.description}</div>
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
            {visibleProducts.length === 0 && <p className='text-gray-500 text-center py-8'>No products found for this category.</p>}
        </main>
    );
};

export default ProductsPage;
