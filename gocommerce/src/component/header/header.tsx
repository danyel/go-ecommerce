import {Link} from 'react-router-dom';
import {LogIn, LogOut, ShoppingBasket} from 'lucide-react';
import {useGlobalState} from '../../state/global-state.tsx';
import {useAuth} from '../../state/use-auth.ts';
import {useEffect, useState} from 'react';
import ApiClient from '../../domain/common/api-client.tsx';
import type {Category} from '../../domain/product/model.tsx';

export interface HeaderComponentProperties {
    showShoppingCart: boolean;
    setShoppingCartState: (state: boolean) => void;
}

const HeaderComponent = (props: HeaderComponentProperties) => {
    const globalStateType = useGlobalState();
    const auth = useAuth();
    const [categories, setCategories] = useState<Category[]>([]);
    useEffect(() => {
        ApiClient.GET<Category[]>('/api/management/v1/categories').then(setCategories).catch(() => setCategories([]));
    }, []);
    return (
        <header className='bg-white shadow-sm sticky top-0 z-50'>
            <div className='max-w-7xl mx-auto px-4 sm:px-6 lg:px-8'>
                <div className='flex items-center justify-between h-16'>
                    {/* logo */}
                    <div className='flex items-center gap-6'>
                        <Link to='/product/products' className='text-2xl font-bold text-blue-600'>Go-Commerce</Link>
                        <nav className='hidden md:flex gap-4 text-sm'>
                            <Link to='/product/products' className='hover:text-blue-600'>Products</Link>
                            <Link to='/management/categories' className='hover:text-blue-600'>Categories</Link>
                            <Link to='/product-management/products' className='hover:text-blue-600'>Product management</Link>
                            {categories.slice(0, 4).map(category => <Link key={category.id} to={`/product/products?category=${category.id}`} className='text-gray-500 hover:text-blue-600'>{category.name}</Link>)}
                        </nav>
                    </div>
                    <div className='flex items-center gap-2'>
                        {auth.isAuthenticated ? <button onClick={auth.logout} className='p-2 hover:bg-gray-100 rounded-lg' title='Sign out'><LogOut size={18}/></button> : <Link to='/login' className='p-2 hover:bg-gray-100 rounded-lg' title='Sign in'><LogIn size={18}/></Link>}
                        <button onClick={() => {
                            props.setShoppingCartState(!props.showShoppingCart)
                        }} className='relative p-2 hover:bg-gray-100 rounded-lg'>
                            <ShoppingBasket size={24}/>
                            <span
                                className='absolute top-0 right-0 bg-red-500 text-white text-xs rounded-full h-5 w-5 flex items-center justify-center'>
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
            </div>
        </header>
    );
};

export default HeaderComponent;