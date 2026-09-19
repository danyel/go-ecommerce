import {useEffect, useState} from 'react';
import {Link} from 'react-router-dom';
import ApiClient from '../../domain/common/api-client.tsx';
import type {Category} from '../../domain/product/model.tsx';
import type {Page} from "../../domain/common/page.tsx";

const CategoriesPage = () => {
    const [categories, setCategories] = useState<Category[]>([]);
    useEffect(() => {
        ApiClient.GET<Page<Category>>('/api/management/v1/categories?page=1&page_size=50').then(data => setCategories(data.items)).catch(() => setCategories([]));
    }, []);
    return <main className='flex-1 bg-white rounded-lg shadow-sm p-6'>
        <h2 className='text-2xl font-semibold mb-4'>Categories</h2>
        <div className='grid sm:grid-cols-2 lg:grid-cols-3 gap-3'>
            {categories.map(category => <Link key={category.id} to={`/product/products?category=${category.id}`} className='border rounded p-4 hover:border-blue-500'>
                <span className='font-medium'>{category.name}</span>
                {category.children?.length ? <span className='block text-sm text-gray-500'>{category.children.map(child => child.name).join(', ')}</span> : null}
            </Link>)}
        </div>
    </main>;
};

export default CategoriesPage;
