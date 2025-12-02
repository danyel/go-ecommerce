import {ShoppingBasket} from "lucide-react";
import {useGlobalState} from "../../state/global-state.tsx";

export interface HeaderComponentProperties {
    showShoppingCart: boolean;
    setShoppingCartState: (state: boolean) => void;
}

const HeaderComponent = (props: HeaderComponentProperties) => {
    const globalStateType = useGlobalState();
    return (
        <header className="bg-white shadow-sm sticky top-0 z-50">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="flex items-center justify-between h-16">
                    {/* logo */}
                    <div className="flex items-center">
                        <h1 className="text-2xl font-bold text-blue-600">Go-Commerce</h1>
                    </div>
                    {/* shopping cart */}
                    <button onClick={() => {
                        props.setShoppingCartState(!props.showShoppingCart)
                    }} className="relative p-2 hover:bg-gray-100 rounded-lg">
                        <ShoppingBasket size={24}/>
                        <span
                            className="absolute top-0 right-0 bg-red-500 text-white text-xs rounded-full h-5 w-5 flex items-center justify-center">
                            {/*todo the shopping cart total items*/}
                            {
                                globalStateType.shoppingBasket.id && globalStateType.shoppingBasket.items?.length > 0 ? globalStateType.shoppingBasket.items?.reduce((accumulator, currentItem) => {
                                    return accumulator + currentItem.quantity;
                                }, 0) : 0
                            }
                        </span>
                    </button>
                </div>
            </div>
        </header>
    );
};

export default HeaderComponent;