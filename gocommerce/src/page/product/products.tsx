import {useEffect, useState} from "react";
import type {Product} from "../../domain/product/model.tsx";
import {ChevronDown} from "lucide-react";
import {useGlobalState} from "../../state/global-state.tsx";
import type {UpdateShoppingBasketItem, ShoppingBasket, ShoppingBasketId} from "../../domain/shopping-basket/model.tsx";
import Cookies from "js-cookie";

const ProductsPage = () => {
    const [products, setProducts] = useState<Product[]>([]);
    const globalStateType = useGlobalState();
    useEffect(() => {
        fetch("/api/product/v1/products", {
            method: "GET",
            headers: {
                'Accept-Language': 'en',
                'Content-Type': 'application/json',
                Accept: 'application/json',
            }
        })
            .then(res => {
                if (!res.ok) {
                    throw Error('Error while fetching the products');
                }

                return res.json();
            })
            .then(data => setProducts(data));
        }, []);

    function addToCart(product: Product) {
        const addItem: UpdateShoppingBasketItem = {product_id: product.id, quantity: 1};
        if (!globalStateType.shoppingBasket.id) {
            fetch('/api/shopping-basket/v1/shopping-baskets', {
                method: "POST",
                headers: {
                    'Accept-Language': 'en',
                    'Content-Type': 'application/json',
                    Accept: 'application/json',
                }
            })
                .then(resp => {
                    if (!resp.ok) {
                        throw Error('Could not create shopping basket');
                    }
                    return resp.json();
                })
                .then((shoppingBasketId: ShoppingBasketId) => {
                    fetch(`/api/shopping-basket/v1/shopping-baskets/${shoppingBasketId.id}`, {
                        method: "POST",
                        body: JSON.stringify(addItem),
                        headers: {
                            'Accept-Language': 'en',
                            'Content-Type': 'application/json',
                            Accept: 'application/json',
                        }
                    })
                        .then(resp => {
                            if (!resp.ok) {
                                throw Error('Could not add item to shopping basket');
                            }
                            return resp.json();
                        })
                        .then((data: ShoppingBasket) => {
                            Cookies.set("shopping_basket_id", shoppingBasketId.id);
                            globalStateType.setShoppingBasket(data);
                        });
                });
        } else {
            fetch(`/api/shopping-basket/v1/shopping-baskets/${globalStateType.shoppingBasket.id}`, {
                method: "POST",
                body: JSON.stringify(addItem),
                headers: {
                    'Accept-Language': 'en',
                    'Content-Type': 'application/json',
                    Accept: 'application/json',
                }
            })
                .then(resp => {
                    if (!resp.ok) {
                        throw Error('Could not add item to shopping basket');
                    }
                    return resp.json();
                })
                .then((data: ShoppingBasket) => {
                    globalStateType.setShoppingBasket(data);
                });
        }

    }

    return (
        <main className="flex-1">
            {/* Sorting Bar */}
            <div
                className="bg-white rounded-lg shadow-sm p-4 mb-6 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                <p className="text-gray-600">
                    Showing <span className="font-semibold">{products.length}</span> products
                </p>
                <div className="flex items-center gap-2">
                    <label className="text-sm text-gray-600">Sort by:</label>
                    <button className="flex items-center gap-2 px-4 py-2 border rounded-lg hover:bg-gray-50">
                        Best Match <ChevronDown size={16}/>
                    </button>
                </div>
            </div>

            {/* Products Grid */}
            <div className="grid grid-cols-4 sm:grid-cols-2 xl:grid-cols-3 gap-6">
                {products.map((product) => (
                    <div
                        key={product.id}
                        className="bg-white rounded-lg shadow-sm hover:shadow-md transition overflow-hidden"
                    >
                        <div className="relative">
                            <img
                                src={product.image_url}
                                alt={product.name}
                                className="w-full h-48 object-cover"
                            />
                            {product.stock === 0 && (
                                <div className="absolute top-2 right-2 bg-red-500 text-white text-xs px-2 py-1 rounded">
                                    Out of Stock
                                </div>
                            )}
                        </div>

                        <div className="p-4">
                            <p className="text-xs text-gray-500 mb-1">{product.category.name}</p>
                            <h3 className="font-semibold mb-2">{product.name}</h3>

                            {/*<div className="flex items-center gap-2 mb-3">*/}
                            {/*<div className="flex items-center">*/}
                            {/*    <Star size={16} className="text-yellow-400 fill-yellow-400" />*/}
                            {/*    <span className="text-sm ml-1">{product.rating}</span>*/}
                            {/*</div>*/}
                            {/*<span className="text-xs text-gray-500">({product.reviews} reviews)</span>*/}
                            {/*</div>*/}

                            <div className="flex items-center justify-between">
                                <p className="text-2xl font-bold text-blue-600">
                                    {product.price} €
                                </p>
                                <button
                                    onClick={() => addToCart(product)}
                                    disabled={product.stock === 0}
                                    className={`px-4 py-2 rounded-lg font-medium transition ${
                                        product.stock !== 0
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