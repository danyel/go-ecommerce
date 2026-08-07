import {Route, Routes} from 'react-router-dom';
import ProductsPage from "./page/product/products.tsx";
import HeaderComponent from "./component/header/header.tsx";
import {useEffect, useState} from "react";
import ShoppingBasketComponent from "./component/shopping-basket/shopping-basket.tsx";
import Cookies from 'js-cookie';
import {useGlobalState} from "./state/global-state.tsx";
import type {ShoppingBasket} from "./domain/shopping-basket/model.tsx";
import {ServiceFactoryFactory} from "./domain/common/factory.tsx";

const App = () => {
    const [showShoppingCart, setShowShoppingCart] = useState(false);
    const globalStateType = useGlobalState();

    useEffect(() => {
        const shoppingBasketIdCookie = Cookies.get("shopping_basket_id");
        if (shoppingBasketIdCookie && !globalStateType.shoppingBasket.id) {
           ServiceFactoryFactory.SHOPPING_BASKET_SERVICE_FACTORY.newService()
               .findById({id: shoppingBasketIdCookie})
                .then((shoppingBasket: ShoppingBasket) => {
                    globalStateType.setShoppingBasket(shoppingBasket);
                });
        }
    }, [globalStateType])
    return (
        <>
            <HeaderComponent showShoppingCart={showShoppingCart} setShoppingCartState={setShowShoppingCart}/>
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                <div className="flex flex-col lg:flex-row gap-8">

                    <Routes>
                        <Route path={"/product/products"}
                               element={<ProductsPage/>}/>
                    </Routes>
                </div>
            </div>
            {
                showShoppingCart &&
                <ShoppingBasketComponent changeShoppingBasketVisibility={setShowShoppingCart}
                                         showShoppingBasket={showShoppingCart}/>
            }
        </>
    );
};

export default App;