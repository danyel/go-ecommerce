import {Route, Routes} from 'react-router-dom';
import ProductsPage from "./page/product/products.tsx";
import HeaderComponent from "./component/header/header.tsx";
import {useState} from "react";
import ShoppingBasketComponent from "./component/shopping-basket/shopping-basket.tsx";

const App = () => {
    const [showShoppingCart, setShowShoppingCart] = useState(false);
    return (
        <>
            {/* Sidebar */}
            {/*<div className={`sidebar`}>*/}
            {/*    <div className="sidebar-header">*/}
            {/*        <h1 className="sidebar-title">Transit App</h1>*/}

            {/*        <X size={24}/>*/}
            {/*    </div>*/}

            {/*    <nav className="nav">*/}

            {/*    </nav>*/}
            {/*</div>*/}
            <HeaderComponent showShoppingCart={showShoppingCart} setShoppingCartState={setShowShoppingCart}/>
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                <div className="flex flex-col lg:flex-row gap-8">

                    <Routes>
                        <Route path={"/product/products"} element={<ProductsPage/>}/>
                    </Routes>
                </div>
            </div>
            {
                showShoppingCart && <ShoppingBasketComponent setShoppingCartState={setShowShoppingCart}
                                                             showShoppingCart={showShoppingCart}/>
            }
        </>
    );
};

export default App;