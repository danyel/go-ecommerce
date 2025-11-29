import {useEffect, useState} from "react";
import {Minus, Plus, X} from "lucide-react";

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
}

export interface ShoppingBasketComponentProperties {
    id?: string;
    setShoppingCartState: (state: boolean) => void;
    showShoppingCart: boolean;
}

const ShoppingBasketComponent = (props: ShoppingBasketComponentProperties) => {
    const [shoppingBasket, setShoppingBasket] = useState<ShoppingBasket>();
    useEffect(() => {
        if (props.id) {
            fetch(`/api/shopping-basket/v1/shopping-baskets/${props.id}`)
                .then(res => {
                    if (!res.ok) {
                        throw Error('Error while fetching the shopping basket');
                    }
                    return res.json();
                })
                .then(data => setShoppingBasket(data));
        }
    }, [props.id]);

    return (
        <div className="fixed inset-0 bg-black bg-opacity=50 z-50" onClick={() => props.setShoppingCartState(false)}>
            <div className="absolute right-0 top-0 h-full w-full max-w-4xl bg-white shadow-xl"
                 onClick={(e) => e.stopPropagation()}>
                <div className="flex flex-col h-full">
                    {/* header part */}
                    <div className="flex items-center justify-between p-6 border-b">
                        <h2 className="text-xl font-semibold">Shopping basket</h2>
                        <button onClick={() => props.setShoppingCartState(false)}>
                            <X size={24}/>
                        </button>
                    </div>

                    {/*  content of the shopping basket  */}
                    {
                        !shoppingBasket || shoppingBasket?.items.length == 0 ?
                            // shopping basket is empty
                            (
                                <div className="flex-1 overflow-y-auto p-6">
                                    <p className="text-gray-500 text-center mt-8">Your cart is empty</p>
                                </div>
                            ) :
                            // display cart items
                            (
                                <div className="flex-1 overflow-y-auto p-6">
                                    <div className="space-y-4">
                                        {
                                            shoppingBasket.items.map((item: ShoppingBasketItem, index: number) => (
                                                <div key={index} className="flex gap-4 border-b pb-4">
                                                    <img src={item.image_url} alt={item.name}
                                                         className="w-20 h-20 object-cover rounded"/>
                                                    <div className="flex-1">
                                                        <h3 className="font-medium">{item.name}</h3>
                                                        <p className="font-semibold mt-1">{item.price}</p>
                                                    </div>
                                                    {/*    add or remove items form the shopping basket*/}
                                                    <div className="flex gap-2">
                                                        <button className="text-red-500 hover:text-red-700">
                                                            <Minus/>
                                                        </button>
                                                        <input size={1} value={item.amount}/>
                                                        <button className="text-red-500 hover:text-red-700">
                                                            <Plus/>
                                                        </button>
                                                    </div>
                                                </div>
                                            ))
                                        }
                                    </div>
                                    {/* checkout part */}
                                    <div className="border-t  p-6">
                                        <div className="flex justify-between mb-4">
                                            <span
                                                className="text-lg font-semibold">{shoppingBasket.total_price_exclusive}</span>
                                            <span className="text-lg font-semibold">{shoppingBasket.tax}</span>
                                            <span
                                                className="text-lg font-semibold">{shoppingBasket.total_price_inclusive}</span>
                                        </div>
                                        {/* todo checkout action */}
                                        <button
                                            className="w-full bg-blue-600 text-white py-3 rounded-lg font-semibold hover:bg-blue-700 transition">
                                            Proceed to checkout
                                        </button>
                                    </div>
                                </div>
                            )
                    }
                </div>
            </div>
        </div>
    );
}

export default ShoppingBasketComponent;