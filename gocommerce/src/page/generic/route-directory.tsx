const RouteDirectory = ({title}: {title: string}) => <main className='flex-1 bg-white rounded-lg shadow-sm p-6'>
    <h2 className='text-2xl font-semibold mb-2'>{title}</h2>
    <p className='text-gray-600'>This route is available through the backend API. Use the product browser and category links for public shopping, or the authenticated checkout flow for basket changes.</p>
</main>;

export default RouteDirectory;
