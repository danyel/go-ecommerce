import {Navigate, Route, Routes} from 'react-router-dom';
import ProductsPage from './page/product/products.tsx';
import HeaderComponent from './component/header/header.tsx';
import {useEffect, useState} from 'react';
import ShoppingBasketComponent from './component/shopping-basket/shopping-basket.tsx';
import Cookies from 'js-cookie';
import {useGlobalState} from './state/global-state.tsx';
import type {ShoppingBasket} from './domain/shopping-basket/model.tsx';
import ApiClient from "./domain/common/api-client.tsx";
import ProtectedRoute from './component/auth/protected-route.tsx';
import LoginPage from './page/auth/login.tsx';
import CategoriesPage from './page/catalog/categories.tsx';
import RouteDirectory from './page/generic/route-directory.tsx';

const App = () => {
    const [showShoppingCart, setShowShoppingCart] = useState(false);
    const globalStateType = useGlobalState();

    useEffect(() => {
        const shoppingBasketIdCookie = Cookies.get('shopping_basket_id');
        if (shoppingBasketIdCookie && !globalStateType.shoppingBasket.id) {
            ApiClient.GET<ShoppingBasket>(`/api/shopping-basket/v1/shopping-baskets/${shoppingBasketIdCookie}`)
                .then((shoppingBasket: ShoppingBasket) => {
                    globalStateType.setShoppingBasket(shoppingBasket);
                });
        }
    }, [globalStateType])
    return (
        <>
            <HeaderComponent showShoppingCart={showShoppingCart} setShoppingCartState={setShowShoppingCart}/>
            <div className='max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8'>
                <div className='flex flex-col lg:flex-row gap-8'>

                    <Routes>
                        <Route path='/' element={<Navigate to='/product/products' replace/>}/>
                        <Route path='/product/products' element={<ProductsPage/>}/>
                        <Route path='/management/categories' element={<CategoriesPage/>}/>
                        <Route path='/login' element={<LoginPage/>}/>
                        <Route path='/product-management/products' element={<RouteDirectory title='Product management'/>}/>
                        <Route path='/cms/translations' element={<RouteDirectory title='Translations'/>}/>
                        <Route element={<ProtectedRoute/>}>
                            <Route path='/checkout' element={<RouteDirectory title='Checkout'/>}/>
                        </Route>
                        <Route path='*' element={<RouteDirectory title='Go-Commerce'/>}/>
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