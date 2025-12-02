import {useEffect, useState} from "react";
import {Minus, Plus, X} from "lucide-react";
import type {
    ShoppingBasket,
    ShoppingBasketItem,
    UpdateShoppingBasketItem
} from "../../domain/shopping-basket/model.tsx";
import {useGlobalState} from "../../state/global-state.tsx";


export interface ShoppingBasketComponentProperties {
    changeShoppingBasketVisibility: (state: boolean) => void;
    showShoppingBasket: boolean;
}

const ShoppingBasketComponent = (props: ShoppingBasketComponentProperties) => {
    const globalStateType = useGlobalState();
    const [shoppingBasket, setShoppingBasket] = useState<ShoppingBasket>(globalStateType.shoppingBasket);

    useEffect(() => {
        // the shopping basket id can only be filled in when we created a shopping basket
        if (!globalStateType.shoppingBasket.id || globalStateType.shoppingBasket.id !== "") {
            setShoppingBasket(globalStateType.shoppingBasket);
        }
    }, [globalStateType.shoppingBasket.id]);

    const addItem = (shoppingBasketItem: ShoppingBasketItem) => {
        updateShoppingBasketItem({
            product_id: shoppingBasketItem.product_id,
            quantity: shoppingBasketItem.quantity + 1
        });
    };

    function updateShoppingBasketItem(updateShoppingBasketItem: UpdateShoppingBasketItem) {
        fetch(`/api/shopping-basket/v1/shopping-baskets/${globalStateType.shoppingBasket.id}`, {
            method: "POST",
            body: JSON.stringify(updateShoppingBasketItem),
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
                setShoppingBasket(data);
            });
    }

    function removeItem(shoppingBasketItem: ShoppingBasketItem) {
        updateShoppingBasketItem({
            product_id: shoppingBasketItem.product_id,
            quantity: shoppingBasketItem.quantity - 1
        });
    }

    function onQuantityChange(shoppingBasketItem: ShoppingBasketItem, quantity: number) {
        if (quantity >= 0) {
            setShoppingBasket(prev => ({
                ...prev,
                items: prev.items.map(i =>
                    i.id === shoppingBasketItem.id ? {...i, quantity} : i
                )
            }));
            updateShoppingBasketItem({
                product_id: shoppingBasketItem.product_id,
                quantity: quantity
            });
        }
    }

    return (
        <div
            className="fixed inset-0 bg-neutral-400/20 backdrop-blur-lg z-50"/* onClick={() => props.setShoppingCartState(false)}*/>
            <div className="absolute right-0 top-0 h-full w-full max-w-4xl bg-white shadow-xl"
                 onClick={(e) => e.stopPropagation()}>
                <div className="flex flex-col h-full">
                    {/* header part */}
                    <div className="flex items-center justify-between p-6 border-b">
                        <h2 className="text-xl font-semibold text-blue-600">Shopping basket</h2>
                        <button onClick={() => props.changeShoppingBasketVisibility(false)}>
                            <X size={24}/>
                        </button>
                    </div>

                    {/*  content of the shopping basket  */}
                    {
                        !shoppingBasket || !shoppingBasket?.items || shoppingBasket?.items?.length == 0 ?
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
                                            shoppingBasket.items.map((shoppingBasketItem: ShoppingBasketItem, index: number) => (
                                                <div key={index} className="flex gap-4 border-b pb-4">
                                                    <img src={shoppingBasketItem.image_url}
                                                         alt={shoppingBasketItem.name}
                                                         className="w-20 h-20 object-cover rounded"/>
                                                    <div className="flex-1">
                                                        <h3 className="font-medium">{shoppingBasketItem.name}</h3>
                                                        <p className="font-semibold mt-1">{shoppingBasketItem.price}</p>
                                                    </div>
                                                    {/*    add or remove items form the shopping basket*/}
                                                    <div className="flex gap-2">
                                                        <button className="text-red-500 hover:text-red-700">
                                                            <Minus onClick={() => removeItem(shoppingBasketItem)}/>
                                                        </button>
                                                        <input size={3} value={shoppingBasketItem.quantity}
                                                               onInput={(e) => {
                                                                   const newValue = Number((e.target as HTMLInputElement).value);
                                                                   onQuantityChange(shoppingBasketItem, newValue);
                                                               }}/>
                                                        <button className="text-red-500 hover:text-red-700">
                                                            <Plus onClick={() => addItem(shoppingBasketItem)}/>
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
                                                className="text-lg font-semibold">{shoppingBasket.total_price_exclusive.toFixed(2)} €</span>
                                            <span
                                                className="text-lg font-semibold">{shoppingBasket.tax.toFixed(2)} €</span>
                                            <span
                                                className="text-lg font-semibold">{shoppingBasket.total_price_inclusive.toFixed(2)} €</span>
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